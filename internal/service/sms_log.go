package service

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"

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

// GetSmsLogPage 获得短信日志分页
func (s *SmsLogService) GetSmsLogPage(ctx context.Context, req *req.SmsLogPageReq) (*core.PageResult[*resp.SmsLogRespVO], error) {
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

	return &core.PageResult[*resp.SmsLogRespVO]{
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
