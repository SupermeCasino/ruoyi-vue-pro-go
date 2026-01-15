package gateway

import (
	"context"
	"log"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/gateway/codec"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
)

// IotGatewayBootstrapper IoT 网关启动器
// 统一管理网关组件的生命周期
type IotGatewayBootstrapper struct {
	config               *MQTTClientConfig
	messageBus           core.MessageBus
	codecRegistry        *codec.CodecRegistry
	mqttClient           *MQTTClient
	connectionManager    *ConnectionManager
	subscribers          []core.MessageSubscriber
	deviceService        *iotsvc.DeviceService
	deviceMessageService *iotsvc.DeviceMessageService
}

// NewIotGatewayBootstrapper 创建网关启动器
func NewIotGatewayBootstrapper(
	config *MQTTClientConfig,
	messageBus core.MessageBus,
	codecRegistry *codec.CodecRegistry,
	deviceService *iotsvc.DeviceService,
	deviceMessageService *iotsvc.DeviceMessageService,
) *IotGatewayBootstrapper {
	return &IotGatewayBootstrapper{
		config:               config,
		messageBus:           messageBus,
		codecRegistry:        codecRegistry,
		deviceService:        deviceService,
		deviceMessageService: deviceMessageService,
	}
}

// Start 启动网关
func (b *IotGatewayBootstrapper) Start(ctx context.Context) error {
	log.Println("[IotGatewayBootstrapper] Starting IoT Gateway...")

	// 1. 创建连接管理器
	serverID := "iot-gateway-" + time.Now().Format("20060102150405")
	b.connectionManager = NewConnectionManager(b.messageBus, serverID)

	// 2. 启动心跳检查器 (60s 超时，30s 检查间隔)
	b.connectionManager.StartHeartbeatChecker(60*time.Second, 30*time.Second)

	// 3. 创建 MQTT 客户端
	b.mqttClient = NewMQTTClient(b.config, b.messageBus, b.codecRegistry)

	// 4. 注册消息订阅者
	deviceMessageSub := NewDeviceMessageSubscriber(b.deviceService, b.deviceMessageService)
	downstreamSub := NewDownstreamSubscriber(b.mqttClient, serverID)

	b.subscribers = []core.MessageSubscriber{deviceMessageSub, downstreamSub}
	for _, sub := range b.subscribers {
		b.messageBus.Register(sub)
		log.Printf("[IotGatewayBootstrapper] Registered subscriber: topic=%s, group=%s",
			sub.Topic(), sub.Group())
	}

	// 启动消息总线
	b.messageBus.Start()

	// 5. 启动 MQTT 客户端
	if err := b.mqttClient.Start(ctx); err != nil {
		return err
	}

	log.Println("[IotGatewayBootstrapper] IoT Gateway started successfully")
	return nil
}

// Stop 停止网关
func (b *IotGatewayBootstrapper) Stop() {
	log.Println("[IotGatewayBootstrapper] Stopping IoT Gateway...")

	if b.mqttClient != nil {
		b.mqttClient.Stop()
	}

	if b.messageBus != nil {
		b.messageBus.Stop()
	}

	log.Println("[IotGatewayBootstrapper] IoT Gateway stopped")
}

// GetConnectionManager 获取连接管理器
func (b *IotGatewayBootstrapper) GetConnectionManager() *ConnectionManager {
	return b.connectionManager
}

// GetMQTTClient 获取 MQTT 客户端
func (b *IotGatewayBootstrapper) GetMQTTClient() *MQTTClient {
	return b.mqttClient
}
