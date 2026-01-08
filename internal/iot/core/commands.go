package core

// DownstreamCommand 下行指令结构
type DownstreamCommand struct {
	ProductKey string
	DeviceName string
	Message    *IotDeviceMessage
}
