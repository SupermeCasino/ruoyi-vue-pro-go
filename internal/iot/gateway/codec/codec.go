package codec

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
)

// DeviceMessageCodec 设备消息编解码器接口
type DeviceMessageCodec interface {
	// Type 返回编解码器类型标识
	Type() string

	// Encode 编码设备消息为字节流
	Encode(message *core.IotDeviceMessage) ([]byte, error)

	// Decode 解码字节流为设备消息
	Decode(bytes []byte) (*core.IotDeviceMessage, error)
}

// CodecRegistry 编解码器注册表
type CodecRegistry struct {
	codecs map[string]DeviceMessageCodec
}

// NewCodecRegistry 创建编解码器注册表
func NewCodecRegistry() *CodecRegistry {
	return &CodecRegistry{
		codecs: make(map[string]DeviceMessageCodec),
	}
}

// Register 注册编解码器
func (r *CodecRegistry) Register(codec DeviceMessageCodec) {
	r.codecs[codec.Type()] = codec
}

// Get 获取指定类型的编解码器
func (r *CodecRegistry) Get(codecType string) DeviceMessageCodec {
	return r.codecs[codecType]
}

// DefaultRegistry 创建默认的编解码器注册表（包含 Alink）
func DefaultRegistry() *CodecRegistry {
	registry := NewCodecRegistry()
	registry.Register(NewAlinkCodec())
	return registry
}
