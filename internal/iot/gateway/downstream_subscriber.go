package gateway

import (
	"log"

	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
)

// DownstreamSubscriber 下行指令订阅者
// 订阅内部消息总线的下行指令，通过 MQTT 客户端推送到设备
type DownstreamSubscriber struct {
	mqttClient *MQTTClient
	serverID   string
}

// NewDownstreamSubscriber 创建下行指令订阅者
func NewDownstreamSubscriber(mqttClient *MQTTClient, serverID string) *DownstreamSubscriber {
	return &DownstreamSubscriber{
		mqttClient: mqttClient,
		serverID:   serverID,
	}
}

// Topic 返回订阅的主题（下行指令主题，基于 ServerID）
func (s *DownstreamSubscriber) Topic() string {
	return "iot_gateway_downstream_" + s.serverID
}

// Group 返回订阅者分组
func (s *DownstreamSubscriber) Group() string {
	return "iot_gateway_downstream"
}

// OnMessage 处理下行指令消息
func (s *DownstreamSubscriber) OnMessage(message any) {
	msg, ok := message.(*DownstreamCommand)
	if !ok {
		log.Printf("[DownstreamSubscriber] Invalid message type")
		return
	}

	log.Printf("[DownstreamSubscriber] Sending downstream command to device: %s.%s",
		msg.ProductKey, msg.DeviceName)

	err := s.mqttClient.SendDownstreamMessage(msg.ProductKey, msg.DeviceName, msg.Message)
	if err != nil {
		log.Printf("[DownstreamSubscriber] Send downstream message failed: %v", err)
	}
}

// DownstreamCommand 下行指令结构
type DownstreamCommand struct {
	ProductKey string
	DeviceName string
	Message    *core.IotDeviceMessage
}

var _ core.MessageSubscriber = (*DownstreamSubscriber)(nil)
