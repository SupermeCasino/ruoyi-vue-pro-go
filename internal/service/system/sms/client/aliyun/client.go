package aliyun

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system/sms/client"
)

type SmsClient struct {
	channel   *model.SystemSmsChannel
	apiClient *dysmsapi.Client
}

func NewSmsClient(channel *model.SystemSmsChannel) (client.SmsClient, error) {
	// 阿里云通常不需要指定 region，或者默认为 cn-hangzhou
	apiClient, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", channel.ApiKey, channel.ApiSecret)
	if err != nil {
		return nil, err
	}
	return &SmsClient{
		channel:   channel,
		apiClient: apiClient,
	}, nil
}

func (c *SmsClient) GetCode() string {
	return consts.SMSChannelCodeAliyun
}

func (c *SmsClient) SendSms(ctx context.Context, mobile string, apiTemplateId string, templateParams []client.KeyValue) (*client.SmsSendResp, error) {
	request := dysmsapi.CreateSendSmsRequest()
	request.PhoneNumbers = mobile
	request.SignName = c.channel.Signature
	request.TemplateCode = apiTemplateId

	// 将 KeyValue 切片转换为 map 以供阿里云 SDK 使用
	paramsMap := make(map[string]any)
	for _, kv := range templateParams {
		paramsMap[kv.Key] = kv.Value
	}

	paramBytes, err := json.Marshal(paramsMap)
	if err != nil {
		return nil, fmt.Errorf("序列化短信参数失败: %w", err)
	}
	request.TemplateParam = string(paramBytes)

	response, err := c.apiClient.SendSms(request)
	if err != nil {
		return nil, err
	}

	return &client.SmsSendResp{
		ApiSendCode:  response.Code,
		ApiSendMsg:   response.Message,
		ApiRequestId: response.RequestId,
		ApiSerialNo:  response.BizId,
	}, nil
}
