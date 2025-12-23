package service

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type SmsChannelService struct {
	q *query.Query
}

func NewSmsChannelService(q *query.Query) *SmsChannelService {
	return &SmsChannelService{
		q: q,
	}
}

// CreateSmsChannel 创建短信渠道
func (s *SmsChannelService) CreateSmsChannel(ctx context.Context, req *req.SmsChannelSaveReq) (int64, error) {

	channel := &model.SystemSmsChannel{
		Signature:   req.Signature,
		Code:        req.Code,
		Status:      *req.Status,
		Remark:      req.Remark,
		ApiKey:      req.ApiKey,
		ApiSecret:   req.ApiSecret,
		CallbackUrl: req.CallbackUrl,
	}
	err := s.q.SystemSmsChannel.WithContext(ctx).Create(channel)
	return channel.ID, err
}

// UpdateSmsChannel 更新短信渠道
func (s *SmsChannelService) UpdateSmsChannel(ctx context.Context, req *req.SmsChannelSaveReq) error {
	c := s.q.SystemSmsChannel
	_, err := c.WithContext(ctx).Where(c.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("短信渠道不存在")
	}

	// 使用 map 更新确保零值能被正确处理（如 Status = 0 时）
	updates := map[string]any{
		"signature":    req.Signature,
		"code":         req.Code,
		"status":       *req.Status,
		"remark":       req.Remark,
		"api_key":      req.ApiKey,
		"api_secret":   req.ApiSecret,
		"callback_url": req.CallbackUrl,
	}
	_, err = c.WithContext(ctx).Where(c.ID.Eq(req.ID)).Updates(updates)
	return err
}

// DeleteSmsChannel 删除短信渠道
func (s *SmsChannelService) DeleteSmsChannel(ctx context.Context, id int64) error {
	c := s.q.SystemSmsChannel
	_, err := c.WithContext(ctx).Where(c.ID.Eq(id)).Delete()
	return err
}

// GetSmsChannel 获得短信渠道
func (s *SmsChannelService) GetSmsChannel(ctx context.Context, id int64) (*resp.SmsChannelRespVO, error) {
	c := s.q.SystemSmsChannel
	item, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetSmsChannelPage 获得短信渠道分页
func (s *SmsChannelService) GetSmsChannelPage(ctx context.Context, req *req.SmsChannelPageReq) (*pagination.PageResult[*resp.SmsChannelRespVO], error) {
	c := s.q.SystemSmsChannel
	qb := c.WithContext(ctx)

	if req.Signature != "" {
		qb = qb.Where(c.Signature.Like("%" + req.Signature + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(c.Status.Eq(*req.Status))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(c.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*resp.SmsChannelRespVO]{
		List:  lo.Map(list, func(item *model.SystemSmsChannel, _ int) *resp.SmsChannelRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

// GetSimpleSmsChannelList 获得短信渠道精简列表
func (s *SmsChannelService) GetSimpleSmsChannelList(ctx context.Context) ([]*resp.SmsChannelSimpleRespVO, error) {
	c := s.q.SystemSmsChannel
	list, err := c.WithContext(ctx).Order(c.ID.Asc()).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *model.SystemSmsChannel, _ int) *resp.SmsChannelSimpleRespVO {
		return &resp.SmsChannelSimpleRespVO{
			ID:        item.ID,
			Signature: item.Signature,
			Code:      item.Code,
		}
	}), nil
}

func (s *SmsChannelService) convertResp(item *model.SystemSmsChannel) *resp.SmsChannelRespVO {
	return &resp.SmsChannelRespVO{
		ID:          item.ID,
		Signature:   item.Signature,
		Code:        item.Code,
		Status:      item.Status,
		Remark:      item.Remark,
		ApiKey:      item.ApiKey,
		ApiSecret:   item.ApiSecret,
		CallbackUrl: item.CallbackUrl,
		CreateTime:  item.CreateTime,
	}
}
