package service

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type SmsLogService struct {
	q *query.Query
}

func NewSmsLogService(q *query.Query) *SmsLogService {
	return &SmsLogService{
		q: q,
	}
}

// CreateSmsLog 创建短信日志
func (s *SmsLogService) CreateSmsLog(ctx context.Context, item *model.SystemSmsLog) (int64, error) {
	err := s.q.SystemSmsLog.WithContext(ctx).Create(item)
	return item.ID, err
}

// CreateSmsLogWithStatus 创建短信日志（根据 isSend 参数设置不同的状态）
func (s *SmsLogService) CreateSmsLogWithStatus(ctx context.Context, mobile string, userId int64, userType int32, isSend bool,
	template *model.SystemSmsTemplate, content string, templateParams map[string]interface{}) (int64, error) {

	// 根据是否需要发送设置不同的状态
	sendStatus := model.SmsSendStatusInit
	if !isSend {
		sendStatus = model.SmsSendStatusIgnore
	}

	log := &model.SystemSmsLog{
		ChannelId:       template.ChannelId,
		ChannelCode:     template.ChannelCode,
		TemplateId:      template.ID,
		TemplateCode:    template.Code,
		TemplateType:    template.Type,
		TemplateContent: content,
		TemplateParams:  templateParams,
		ApiTemplateId:   template.ApiTemplateId,
		Mobile:          mobile,
		UserId:          userId,
		UserType:        userType,
		SendStatus:      sendStatus,
		SendTime:        nil,
		ReceiveStatus:   model.SmsReceiveStatusInit,
	}
	err := s.q.SystemSmsLog.WithContext(ctx).Create(log)
	return log.ID, err
}

// UpdateSmsLog 更新短信日志
func (s *SmsLogService) UpdateSmsLog(ctx context.Context, item *model.SystemSmsLog) error {
	_, err := s.q.SystemSmsLog.WithContext(ctx).Where(s.q.SystemSmsLog.ID.Eq(item.ID)).Updates(item)
	return err
}

// UpdateSmsLogFields 更新短信日志指定字段
func (s *SmsLogService) UpdateSmsLogFields(ctx context.Context, logId int64, updates map[string]interface{}) error {
	_, err := s.q.SystemSmsLog.WithContext(ctx).Where(s.q.SystemSmsLog.ID.Eq(logId)).Updates(updates)
	return err
}

// GetSmsLogPage 获得短信日志分页
func (s *SmsLogService) GetSmsLogPage(ctx context.Context, req *req.SmsLogPageReq) (*pagination.PageResult[*resp.SmsLogRespVO], error) {
	l := s.q.SystemSmsLog
	qb := l.WithContext(ctx)

	if req.ChannelId != nil {
		qb = qb.Where(l.ChannelId.Eq(*req.ChannelId))
	}
	if req.TemplateId != nil {
		qb = qb.Where(l.TemplateId.Eq(*req.TemplateId))
	}
	if req.Mobile != "" {
		qb = qb.Where(l.Mobile.Like("%" + req.Mobile + "%"))
	}
	if req.SendStatus != nil {
		qb = qb.Where(l.SendStatus.Eq(*req.SendStatus))
	}
	if req.ReceiveStatus != nil {
		qb = qb.Where(l.ReceiveStatus.Eq(*req.ReceiveStatus))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(l.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*resp.SmsLogRespVO]{
		List:  lo.Map(list, func(item *model.SystemSmsLog, _ int) *resp.SmsLogRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *SmsLogService) convertResp(item *model.SystemSmsLog) *resp.SmsLogRespVO {
	return &resp.SmsLogRespVO{
		ID:              item.ID,
		ChannelId:       item.ChannelId,
		ChannelCode:     item.ChannelCode,
		TemplateId:      item.TemplateId,
		TemplateCode:    item.TemplateCode,
		TemplateType:    item.TemplateType,
		TemplateContent: item.TemplateContent,
		TemplateParams:  item.TemplateParams,
		ApiTemplateId:   item.ApiTemplateId,
		Mobile:          item.Mobile,
		UserId:          item.UserId,
		UserType:        item.UserType,
		SendStatus:      item.SendStatus,
		SendTime:        item.SendTime,
		ApiSendCode:     item.ApiSendCode,
		ApiSendMsg:      item.ApiSendMsg,
		ApiRequestId:    item.ApiRequestId,
		ApiSerialNo:     item.ApiSerialNo,
		ReceiveStatus:   item.ReceiveStatus,
		ReceiveTime:     item.ReceiveTime,
		ApiReceiveCode:  item.ApiReceiveCode,
		ApiReceiveMsg:   item.ApiReceiveMsg,
		CreateTime:      item.CreateTime,
	}
}
