package gateway

import (
	"context"
	"log"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
)

// DeviceMessageSubscriber 设备消息业务订阅者
// 消费内部消息总线的上行报文，执行业务逻辑（属性入库、状态同步等）
type DeviceMessageSubscriber struct {
	deviceService        *iotsvc.DeviceService
	deviceMessageService *iotsvc.DeviceMessageService
}

// NewDeviceMessageSubscriber 创建设备消息业务订阅者
func NewDeviceMessageSubscriber(
	deviceService *iotsvc.DeviceService,
	deviceMessageService *iotsvc.DeviceMessageService,
) *DeviceMessageSubscriber {
	return &DeviceMessageSubscriber{
		deviceService:        deviceService,
		deviceMessageService: deviceMessageService,
	}
}

// Topic 返回订阅的主题
func (s *DeviceMessageSubscriber) Topic() string {
	return core.DeviceMessageTopic
}

// Group 返回订阅者分组
func (s *DeviceMessageSubscriber) Group() string {
	return "iot_device_message_consumer"
}

// OnMessage 处理上行设备消息
func (s *DeviceMessageSubscriber) OnMessage(message any) {
	msg, ok := message.(*core.IotDeviceMessage)
	if !ok {
		log.Printf("[DeviceMessageSubscriber] Invalid message type")
		return
	}

	// 1. 检查是否为上行消息
	if !msg.IsUpstreamMessage() {
		return
	}

	ctx := context.Background()

	// 2. 更新设备最后上报时间
	s.updateDeviceReportTime(ctx, msg.DeviceID, msg.ReportTime)

	// 3. 获取设备信息
	device, err := s.deviceService.Get(ctx, msg.DeviceID)
	if err != nil {
		log.Printf("[DeviceMessageSubscriber] Get device failed: %v", err)
		return
	}
	if device == nil {
		log.Printf("[DeviceMessageSubscriber] Device not found: %d", msg.DeviceID)
		return
	}

	// 4. 交由 DeviceMessageService 处理业务逻辑（属性保存、状态更新、日志记录、回复）
	s.deviceMessageService.HandleUpstreamDeviceMessage(ctx, msg, device)
}

// updateDeviceReportTime 更新设备上报时间
func (s *DeviceMessageSubscriber) updateDeviceReportTime(ctx context.Context, deviceID int64, reportTime time.Time) {
	if err := s.deviceService.UpdateDeviceActiveTime(ctx, deviceID, reportTime); err != nil {
		log.Printf("[DeviceMessageSubscriber] Update device active time failed: %v", err)
	}
}

var _ core.MessageSubscriber = (*DeviceMessageSubscriber)(nil)
