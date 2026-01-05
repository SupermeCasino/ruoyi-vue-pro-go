package core

import (
	"time"
)

// DeviceMessageTopic 设备消息主题 (与 Java IotDeviceMessage.MESSAGE_BUS_DEVICE_MESSAGE_TOPIC 对齐)
const DeviceMessageTopic = "iot_device_message"

// IotDeviceMessage IoT 设备消息结构体
// 与 Java cn.iocoder.yudao.module.iot.core.mq.message.IotDeviceMessage 对齐
type IotDeviceMessage struct {
	// ID 消息唯一 ID
	ID string `json:"id"`

	// RequestID 请求 ID（用于请求-响应模式）
	RequestID string `json:"requestId,omitempty"`

	// Method 请求方法 (如 thing.event.property.post)
	Method string `json:"method"`

	// Params 请求参数
	Params map[string]interface{} `json:"params,omitempty"`

	// Data 响应数据
	Data any `json:"data,omitempty"`

	// Code 响应错误码
	Code int `json:"code,omitempty"`

	// Msg 响应提示
	Msg string `json:"msg,omitempty"`

	// DeviceID 设备 ID
	DeviceID int64 `json:"deviceId,omitempty"`

	// TenantID 租户 ID
	TenantID int64 `json:"tenantId,omitempty"`

	// ServerID 设备连接的网关服务器 ID
	ServerID string `json:"serverId,omitempty"`

	// ReportTime 上报时间
	ReportTime time.Time `json:"reportTime,omitempty"`
}

// BuildStateUpdateOnline 构建设备上线状态消息
func BuildStateUpdateOnline() *IotDeviceMessage {
	return &IotDeviceMessage{
		Method: "thing.lifecycle.state.update",
		Params: map[string]any{"state": "online"},
	}
}

// BuildStateOffline 构建设备离线状态消息
func BuildStateOffline() *IotDeviceMessage {
	return &IotDeviceMessage{
		Method: "thing.lifecycle.state.update",
		Params: map[string]any{"state": "offline"},
	}
}

// IsUpstreamMessage 判断是否为上行消息
func (m *IotDeviceMessage) IsUpstreamMessage() bool {
	// 上行消息特征：无 Code (非响应)
	return m.Code == 0 && m.Data == nil
}
