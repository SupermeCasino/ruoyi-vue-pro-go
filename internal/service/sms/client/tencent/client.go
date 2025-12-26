package tencent

import (
	"context"
	"fmt"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/sms/client"
)

type SmsClient struct {
	channel   *model.SystemSmsChannel
	apiClient *sms.Client
	sdkAppId  string
}

func NewSmsClient(channel *model.SystemSmsChannel) (client.SmsClient, error) {
	// 格式为 [secretId sdkAppId]
	apiKey := strings.TrimSpace(channel.ApiKey)
	parts := strings.Split(apiKey, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("腾讯云短信 apiKey 配置格式错误，请配置为 [secretId sdkAppId]")
	}
	secretId := parts[0]
	sdkAppId := parts[1]

	credential := common.NewCredential(secretId, channel.ApiSecret)
	cpf := profile.NewClientProfile()
	// 默认使用广州地域，或者可以从配置中读取（如果模型支持）
	apiClient, err := sms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return nil, err
	}

	return &SmsClient{
		channel:   channel,
		apiClient: apiClient,
		sdkAppId:  sdkAppId,
	}, nil
}

func (c *SmsClient) GetCode() string {
	return consts.SMSChannelCodeTencent
}

func (c *SmsClient) SendSms(ctx context.Context, mobile string, apiTemplateId string, templateParams []client.KeyValue) (*client.SmsSendResp, error) {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(c.sdkAppId)
	request.SignName = common.StringPtr(c.channel.Signature)
	request.TemplateId = common.StringPtr(apiTemplateId)
	request.PhoneNumberSet = common.StringPtrs([]string{mobile})

	// 腾讯云短信参数是按顺序传递的值数组
	// 我们直接利用传入的有序 KeyValue 切片
	var params []string
	for _, kv := range templateParams {
		params = append(params, fmt.Sprint(kv.Value))
	}
	request.TemplateParamSet = common.StringPtrs(params)

	response, err := c.apiClient.SendSms(request)
	if err != nil {
		return nil, err
	}

	if len(response.Response.SendStatusSet) == 0 {
		return nil, fmt.Errorf("腾讯云短信响应状态为空")
	}

	status := response.Response.SendStatusSet[0]
	return &client.SmsSendResp{
		ApiSendCode:  *status.Code,
		ApiSendMsg:   *status.Message,
		ApiRequestId: *response.Response.RequestId,
		ApiSerialNo:  *status.SerialNo,
	}, nil
}
