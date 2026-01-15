package gateway

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/gateway/codec"
)

// MQTTClientConfig MQTT 客户端配置
type MQTTClientConfig struct {
	Broker           string        `yaml:"broker"`             // EMQX Broker 地址，如 tcp://localhost:1883
	ClientID         string        `yaml:"client_id"`          // 客户端 ID
	Username         string        `yaml:"username"`           // 用户名
	Password         string        `yaml:"password"`           // 密码
	KeepAlive        time.Duration `yaml:"keep_alive"`         // 心跳间隔
	ConnectTimeout   time.Duration `yaml:"connect_timeout"`    // 连接超时
	AutoReconnect    bool          `yaml:"auto_reconnect"`     // 自动重连
	CleanSession     bool          `yaml:"clean_session"`      // 清理会话
	SubscribeTopics  []string      `yaml:"subscribe_topics"`   // 订阅主题列表
	DefaultCodecType string        `yaml:"default_codec_type"` // 默认编解码器类型
	TopicPrefix      string        `yaml:"topic_prefix"`       // 主题前缀，默认 /sys
}

// MQTTClient Paho MQTT 客户端封装
// 用于对接 EMQX Broker，处理设备上行消息与下行指令
type MQTTClient struct {
	config        *MQTTClientConfig
	client        mqtt.Client
	messageBus    core.MessageBus
	codecRegistry *codec.CodecRegistry
	mu            sync.RWMutex
	running       bool
}

// NewMQTTClient 创建 MQTT 客户端
func NewMQTTClient(config *MQTTClientConfig, messageBus core.MessageBus, codecRegistry *codec.CodecRegistry) *MQTTClient {
	if config.TopicPrefix == "" {
		config.TopicPrefix = "/sys"
	}
	return &MQTTClient{
		config:        config,
		messageBus:    messageBus,
		codecRegistry: codecRegistry,
	}
}

// Start 启动 MQTT 客户端
func (c *MQTTClient) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	opts := mqtt.NewClientOptions().
		AddBroker(c.config.Broker).
		SetClientID(c.config.ClientID).
		SetUsername(c.config.Username).
		SetPassword(c.config.Password).
		SetKeepAlive(c.config.KeepAlive).
		SetConnectTimeout(c.config.ConnectTimeout).
		SetAutoReconnect(c.config.AutoReconnect).
		SetCleanSession(c.config.CleanSession).
		SetOnConnectHandler(c.onConnect).
		SetConnectionLostHandler(c.onConnectionLost).
		SetDefaultPublishHandler(c.onMessage)

	c.client = mqtt.NewClient(opts)

	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("mqtt connect failed: %w", token.Error())
	}

	c.running = true
	log.Printf("[MQTTClient] Connected to broker: %s", c.config.Broker)

	return nil
}

// Stop 停止 MQTT 客户端
func (c *MQTTClient) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	}

	c.running = false
	log.Printf("[MQTTClient] Disconnected from broker")
}

// Publish 发布消息到指定主题
func (c *MQTTClient) Publish(topic string, qos byte, payload []byte) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("mqtt client not connected")
	}

	token := c.client.Publish(topic, qos, false, payload)
	token.Wait()
	return token.Error()
}

// onConnect 连接成功回调
func (c *MQTTClient) onConnect(client mqtt.Client) {
	log.Printf("[MQTTClient] Connected, subscribing to topics...")

	for _, topic := range c.config.SubscribeTopics {
		token := client.Subscribe(topic, 1, nil)
		if token.Wait() && token.Error() != nil {
			log.Printf("[MQTTClient] Subscribe to %s failed: %v", topic, token.Error())
		} else {
			log.Printf("[MQTTClient] Subscribed to topic: %s", topic)
		}
	}
}

// onConnectionLost 连接丢失回调
func (c *MQTTClient) onConnectionLost(client mqtt.Client, err error) {
	log.Printf("[MQTTClient] Connection lost: %v", err)
}

// onMessage 接收消息回调（上行消息处理）
func (c *MQTTClient) onMessage(client mqtt.Client, msg mqtt.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[MQTTClient] Panic recovered in onMessage: %v", r)
		}
	}()

	topic := msg.Topic()
	payload := msg.Payload()

	// log.Printf("[MQTTClient] Received message on topic: %s", topic)

	// 1. 解析 Topic 获取 productKey 和 deviceName
	// Topic 格式: {TopicPrefix}/{productKey}/{deviceName}/thing/event/property/post
	_, _, err := c.parseTopicDeviceInfo(topic)
	if err != nil {
		// log.Printf("[MQTTClient] Parse topic failed: %v", err)
		return
	}

	// 2. 使用编解码器解码消息
	codecType := c.config.DefaultCodecType
	if codecType == "" {
		codecType = "Alink"
	}
	deviceCodec := c.codecRegistry.Get(codecType)
	if deviceCodec == nil {
		log.Printf("[MQTTClient] Codec not found: %s", codecType)
		return
	}

	message, err := deviceCodec.Decode(payload)
	if err != nil {
		log.Printf("[MQTTClient] Decode message failed: %v", err)
		return
	}

	// 3. 补充消息元数据
	message.ReportTime = time.Now()

	// log.Printf("[MQTTClient] Device message: productKey=%s, deviceName=%s, method=%s",
	// 	productKey, deviceName, message.Method)

	// 4. 发布到内部消息总线
	c.messageBus.Post(core.DeviceMessageTopic, message)
}

// parseTopicDeviceInfo 从 Topic 中解析设备信息
// Topic 格式: {TopicPrefix}/{productKey}/{deviceName}/...
func (c *MQTTClient) parseTopicDeviceInfo(topic string) (productKey, deviceName string, err error) {
	// 去除前缀
	prefix := c.config.TopicPrefix
	if !strings.HasPrefix(topic, prefix) {
		return "", "", fmt.Errorf("invalid topic prefix: %s", topic)
	}

	relativeTopic := strings.TrimPrefix(topic, prefix)
	// relativeTopic: /{productKey}/{deviceName}/...

	parts := strings.Split(relativeTopic, "/")
	// parts[0] = "", parts[1] = productKey, parts[2] = deviceName

	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid topic format: %s", topic)
	}

	return parts[1], parts[2], nil
}

// SendDownstreamMessage 发送下行指令到设备
func (c *MQTTClient) SendDownstreamMessage(productKey, deviceName string, message *core.IotDeviceMessage) error {
	// 阿里 Alink 协议下行 Topic 逻辑 (严格对齐 Java IotMqttTopicUtils.buildTopicByMethod)
	// 逻辑：{TopicPrefix}/{productKey}/{deviceName}/ + strings.ReplaceAll(method, ".", "/") + (isReply ? "_reply" : "")
	topicSuffix := strings.ReplaceAll(message.Method, ".", "/")
	if message.Code != nil {
		topicSuffix += "_reply"
	}

	prefix := c.config.TopicPrefix
	if prefix == "" {
		prefix = "/sys"
	}

	topic := fmt.Sprintf("%s/%s/%s/%s", prefix, productKey, deviceName, topicSuffix)

	// 编码消息
	codecType := c.config.DefaultCodecType
	if codecType == "" {
		codecType = "Alink"
	}
	deviceCodec := c.codecRegistry.Get(codecType)
	if deviceCodec == nil {
		return fmt.Errorf("codec not found: %s", codecType)
	}

	payload, err := deviceCodec.Encode(message)
	if err != nil {
		return fmt.Errorf("encode message failed: %w", err)
	}

	log.Printf("[MQTTClient] Publishing downstream message to topic: %s", topic)
	return c.Publish(topic, 1, payload)
}
