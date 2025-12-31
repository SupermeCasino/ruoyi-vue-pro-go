package client

import "context"

// SmsSendResp 短信发送结果
type SmsSendResp struct {
	ApiSendCode  string
	ApiSendMsg   string
	ApiRequestId string
	ApiSerialNo  string
}

// KeyValue 键值对
type KeyValue struct {
	Key   string
	Value any
}

// SmsClient 短信客户端接口
type SmsClient interface {
	// GetCode 获得渠道编码
	GetCode() string
	// SendSms 发送消息
	SendSms(ctx context.Context, mobile string, apiTemplateId string, templateParams []KeyValue) (*SmsSendResp, error)
}
