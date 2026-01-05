package gateway

import (
	"log"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
)

// DeviceMessageSubscriber 设备消息业务订阅者
// 消费内部消息总线的上行报文，执行业务逻辑（属性入库、状态同步等）
type DeviceMessageSubscriber struct {
	// TODO: 添加依赖注入
	// deviceService 设备服务（用于状态更新）
	// propertyService 属性服务（用于属性入库）
	// 这些依赖将在后续集成时注入
}

// NewDeviceMessageSubscriber 创建设备消息业务订阅者
func NewDeviceMessageSubscriber() *DeviceMessageSubscriber {
	return &DeviceMessageSubscriber{}
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
		log.Printf("[DeviceMessageSubscriber] Not an upstream message, skipping")
		return
	}

	log.Printf("[DeviceMessageSubscriber] Processing device message: method=%s, deviceId=%d",
		msg.Method, msg.DeviceID)

	// 2. 更新设备最后上报时间
	s.updateDeviceReportTime(msg.DeviceID, msg.ReportTime)

	// 3. 根据消息类型处理业务逻辑
	switch {
	case isPropertyPost(msg.Method):
		s.handlePropertyPost(msg)
	case isEventPost(msg.Method):
		s.handleEventPost(msg)
	case isStateUpdate(msg.Method):
		s.handleStateUpdate(msg)
	default:
		log.Printf("[DeviceMessageSubscriber] Unknown method: %s", msg.Method)
	}
}

// updateDeviceReportTime 更新设备上报时间
func (s *DeviceMessageSubscriber) updateDeviceReportTime(deviceID int64, reportTime time.Time) {
	// TODO: 调用 DevicePropertyService 更新设备上报时间
	log.Printf("[DeviceMessageSubscriber] Update device report time: deviceId=%d, time=%v",
		deviceID, reportTime)
}

// handlePropertyPost 处理属性上报
func (s *DeviceMessageSubscriber) handlePropertyPost(msg *core.IotDeviceMessage) {
	log.Printf("[DeviceMessageSubscriber] Property post: deviceId=%d, params=%v",
		msg.DeviceID, msg.Params)
	// TODO: 调用 DevicePropertyService 保存属性值
}

// handleEventPost 处理事件上报
func (s *DeviceMessageSubscriber) handleEventPost(msg *core.IotDeviceMessage) {
	log.Printf("[DeviceMessageSubscriber] Event post: deviceId=%d, params=%v",
		msg.DeviceID, msg.Params)
	// TODO: 调用事件处理逻辑
}

// handleStateUpdate 处理状态更新
func (s *DeviceMessageSubscriber) handleStateUpdate(msg *core.IotDeviceMessage) {
	log.Printf("[DeviceMessageSubscriber] State update: deviceId=%d, params=%v",
		msg.DeviceID, msg.Params)
	// TODO: 调用 DeviceService 更新设备状态
}

// isPropertyPost 判断是否为属性上报
func isPropertyPost(method string) bool {
	return method == "thing.event.property.post"
}

// isEventPost 判断是否为事件上报
func isEventPost(method string) bool {
	return method == "thing.event.post" ||
		(len(method) > 12 && method[:12] == "thing.event." && method != "thing.event.property.post")
}

// isStateUpdate 判断是否为状态更新
func isStateUpdate(method string) bool {
	return method == "thing.lifecycle.state.update"
}

var _ core.MessageSubscriber = (*DeviceMessageSubscriber)(nil)
