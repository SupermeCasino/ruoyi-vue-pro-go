package service

import (
	"context"
	"errors"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/sms/client"
	bzErr "github.com/wxlbd/ruoyi-mall-go/pkg/errors"

	"go.uber.org/zap"
)

type SmsSendService struct {
	q           *query.Query
	templateSvc *SmsTemplateService
	smsLogSvc   *SmsLogService
	factory     *SmsClientFactory
}

func NewSmsSendService(
	q *query.Query,
	templateSvc *SmsTemplateService,
	smsLogSvc *SmsLogService,
	factory *SmsClientFactory,
) *SmsSendService {
	return &SmsSendService{
		q:           q,
		templateSvc: templateSvc,
		smsLogSvc:   smsLogSvc,
		factory:     factory,
	}
}

// SendSingleSmsToAdmin 发送单条短信给 Admin 用户
func (s *SmsSendService) SendSingleSmsToAdmin(ctx context.Context, mobile string, userId int64, templateCode string, templateParams map[string]interface{}) (int64, error) {
	// 如果 mobile 为空，查询 Admin 用户手机号 (此处暂略，假设 mobile 必传或调用者已处理)
	return s.SendSingleSms(ctx, mobile, userId, 1, templateCode, templateParams) // 1: Admin
}

// SendSingleSmsToMember 发送单条短信给 Member 用户
func (s *SmsSendService) SendSingleSmsToMember(ctx context.Context, mobile string, userId int64, templateCode string, templateParams map[string]interface{}) (int64, error) {
	// 如果 mobile 为空，查询 Member 用户手机号
	return s.SendSingleSms(ctx, mobile, userId, 2, templateCode, templateParams) // 2: Member
}

// SendSingleSms 发送单条短信
func (s *SmsSendService) SendSingleSms(ctx context.Context, mobile string, userId int64, userType int32, templateCode string, templateParams map[string]interface{}) (int64, error) {
	// 1. 校验参数
	if mobile == "" {
		return 0, bzErr.NewBizError(400, "手机号不能为空")
	}

	// 2. 获得短信模板
	t := s.q.SystemSmsTemplate
	template, err := t.WithContext(ctx).Where(t.Code.Eq(templateCode)).First()
	if err != nil {
		return 0, bzErr.NewBizError(1004003002, "短信模板不存在")
	}

	// 3. 获得短信渠道
	c := s.q.SystemSmsChannel
	channel, err := c.WithContext(ctx).Where(c.ID.Eq(template.ChannelId)).First()
	if err != nil {
		return 0, bzErr.NewBizError(1004003001, "短信渠道不存在")
	}

	// 4. 创建发送日志
	log := &model.SystemSmsLog{
		ChannelId:       channel.ID,
		ChannelCode:     channel.Code,
		TemplateId:      template.ID,
		TemplateCode:    template.Code,
		TemplateType:    template.Type,
		TemplateContent: s.templateSvc.FormatSmsTemplateContent(template.Content, templateParams),
		TemplateParams:  templateParams,
		ApiTemplateId:   template.ApiTemplateId,
		Mobile:          mobile,
		UserId:          userId,
		UserType:        userType,
		SendStatus:      0, // Sending
		SendTime:        nil,
		ReceiveStatus:   0, // Waiting
	}
	logId, err := s.smsLogSvc.CreateSmsLog(ctx, log)
	if err != nil {
		zap.L().Error("Create SMS log failed", zap.Error(err))
		return 0, err
	}

	// 5. 调用 Client 发送
	smsClient := s.factory.GetClient(channel.ID)
	if smsClient == nil {
		s.factory.CreateOrUpdateClient(channel)
		smsClient = s.factory.GetClient(channel.ID)
	}

	if smsClient == nil {
		s.updateLogSendFail(ctx, log, errors.New("短信客户端初始化失败"))
		return logId, errors.New("短信客户端初始化失败")
	}

	// 执行发送
	sendResp, err := smsClient.SendSms(ctx, logId, mobile, template.ApiTemplateId, templateParams)

	// 6. 更新日志
	if err != nil {
		s.updateLogSendFail(ctx, log, err)
		return logId, err
	}
	s.updateLogSendSuccess(ctx, log, sendResp)

	return logId, nil
}

func (s *SmsSendService) updateLogSendFail(ctx context.Context, log *model.SystemSmsLog, err error) {
	log.SendStatus = 2 // Fail
	now := time.Now()
	log.SendTime = &now
	log.ApiSendMsg = err.Error()
	_ = s.smsLogSvc.UpdateSmsLog(ctx, log)
}

func (s *SmsSendService) updateLogSendSuccess(ctx context.Context, log *model.SystemSmsLog, sendResp *client.SmsSendResp) {
	log.SendStatus = 1 // Success
	now := time.Now()
	log.SendTime = &now
	if sendResp != nil {
		log.ApiSendCode = sendResp.ApiSendCode
		log.ApiSendMsg = sendResp.ApiSendMsg
		log.ApiRequestId = sendResp.ApiRequestId
		log.ApiSerialNo = sendResp.ApiSerialNo
	}
	_ = s.smsLogSvc.UpdateSmsLog(ctx, log)
}
