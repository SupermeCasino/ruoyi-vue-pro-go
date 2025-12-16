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

// UpdateSmsLog 更新短信日志
func (s *SmsLogService) UpdateSmsLog(ctx context.Context, item *model.SystemSmsLog) error {
	_, err := s.q.SystemSmsLog.WithContext(ctx).Where(s.q.SystemSmsLog.ID.Eq(item.ID)).Updates(item)
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
		CreateTime:      item.CreatedAt,
	}
}
