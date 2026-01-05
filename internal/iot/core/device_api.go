package core

import (
	"context"
)

// IotDeviceCommonApi IoT 设备通用 API 接口
// 用于 Gateway 层调用 Biz 层的设备认证与信息查询
type IotDeviceCommonApi interface {
	// AuthDevice 设备认证
	AuthDevice(ctx context.Context, req *DeviceAuthReq) (*DeviceAuthResp, error)

	// GetDevice 获取设备信息
	GetDevice(ctx context.Context, req *DeviceGetReq) (*DeviceResp, error)
}

// DeviceAuthReq 设备认证请求
type DeviceAuthReq struct {
	ClientID string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// DeviceAuthResp 设备认证响应
type DeviceAuthResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// DeviceGetReq 获取设备信息请求
type DeviceGetReq struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// DeviceResp 设备信息响应
type DeviceResp struct {
	ID           int64  `json:"id"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
	State        int    `json:"state"`
	CodecType    string `json:"codecType"`
	TenantID     int64  `json:"tenantId"`
}
