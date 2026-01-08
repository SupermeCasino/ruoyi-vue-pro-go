package codec

import (
	"encoding/json"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
)

// AlinkVersion Alink 协议版本
const AlinkVersion = "1.0"

// AlinkMessage 阿里云 Alink 协议消息结构
type AlinkMessage struct {
	// ID 消息 ID
	ID string `json:"id"`

	// Version 版本号
	Version string `json:"version"`

	// Method 请求方法
	Method string `json:"method,omitempty"`

	// Params 请求参数
	Params any `json:"params,omitempty"`

	// Data 响应数据
	Data any `json:"data,omitempty"`

	// Code 响应错误码
	// Code 响应错误码 (*int 对齐 Java Integer)
	Code *int `json:"code,omitempty"`

	// Msg 响应提示
	Msg string `json:"msg,omitempty"`
}

// AlinkCodec 阿里云 Alink 协议编解码器
type AlinkCodec struct{}

// NewAlinkCodec 创建 Alink 编解码器
func NewAlinkCodec() *AlinkCodec {
	return &AlinkCodec{}
}

// Type 返回编解码器类型
func (c *AlinkCodec) Type() string {
	return "Alink"
}

// Encode 编码设备消息为 Alink 格式字节流
func (c *AlinkCodec) Encode(message *core.IotDeviceMessage) ([]byte, error) {
	alinkMsg := &AlinkMessage{
		ID:      message.RequestID,
		Version: AlinkVersion,
		Method:  message.Method,
		Params:  message.Params,
		Data:    message.Data,
		Code:    message.Code,
		Msg:     message.Msg,
	}
	return json.Marshal(alinkMsg)
}

// Decode 解码 Alink 格式字节流为设备消息
func (c *AlinkCodec) Decode(bytes []byte) (*core.IotDeviceMessage, error) {
	var alinkMsg AlinkMessage
	if err := json.Unmarshal(bytes, &alinkMsg); err != nil {
		return nil, err
	}

	if alinkMsg.Version != AlinkVersion {
		return nil, errors.New("unsupported Alink version: " + alinkMsg.Version)
	}

	// 转换 params 为 map
	var params map[string]interface{}
	if alinkMsg.Params != nil {
		switch p := alinkMsg.Params.(type) {
		case map[string]interface{}:
			params = p
		default:
			// 尝试再次解析
			b, _ := json.Marshal(alinkMsg.Params)
			json.Unmarshal(b, &params)
		}
	}

	return &core.IotDeviceMessage{
		ID:        alinkMsg.ID,
		RequestID: alinkMsg.ID,
		Method:    alinkMsg.Method,
		Params:    params,
		Data:      alinkMsg.Data,
		Code:      alinkMsg.Code,
		Msg:       alinkMsg.Msg,
	}, nil
}
