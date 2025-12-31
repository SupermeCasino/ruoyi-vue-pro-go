package system

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type SocialClientService struct {
	q *query.Query
}

func NewSocialClientService(q *query.Query) *SocialClientService {
	return &SocialClientService{q: q}
}

// CreateSocialClient 创建社交客户端
func (s *SocialClientService) CreateSocialClient(ctx context.Context, r *req.SocialClientSaveReq) (int64, error) {
	client := &model.SocialClient{
		Name:         r.Name,
		SocialType:   r.SocialType,
		UserType:     r.UserType,
		ClientId:     r.ClientId,
		ClientSecret: r.ClientSecret,
		AgentId:      r.AgentId,
		Status:       r.Status,
	}
	if err := s.q.SocialClient.WithContext(ctx).Create(client); err != nil {
		return 0, err
	}
	return client.ID, nil
}

// UpdateSocialClient 更新社交客户端
func (s *SocialClientService) UpdateSocialClient(ctx context.Context, r *req.SocialClientSaveReq) error {
	if r.ID == nil {
		return errors.New("ID不能为空")
	}
	_, err := s.q.SocialClient.WithContext(ctx).Where(s.q.SocialClient.ID.Eq(*r.ID)).Updates(map[string]interface{}{
		"name":          r.Name,
		"social_type":   r.SocialType,
		"user_type":     r.UserType,
		"client_id":     r.ClientId,
		"client_secret": r.ClientSecret,
		"agent_id":      r.AgentId,
		"status":        r.Status,
	})
	return err
}

// DeleteSocialClient 删除社交客户端
func (s *SocialClientService) DeleteSocialClient(ctx context.Context, id int64) error {
	_, err := s.q.SocialClient.WithContext(ctx).Where(s.q.SocialClient.ID.Eq(id)).Delete()
	return err
}

// GetSocialClient 获取社交客户端
func (s *SocialClientService) GetSocialClient(ctx context.Context, id int64) (*model.SocialClient, error) {
	return s.q.SocialClient.WithContext(ctx).Where(s.q.SocialClient.ID.Eq(id)).First()
}

// GetSocialClientPage 获取社交客户端分页
func (s *SocialClientService) GetSocialClientPage(ctx context.Context, r *req.SocialClientPageReq) (*pagination.PageResult[*model.SocialClient], error) {
	q := s.q.SocialClient.WithContext(ctx)

	if r.Name != "" {
		q = q.Where(s.q.SocialClient.Name.Like("%" + r.Name + "%"))
	}
	if r.SocialType != nil {
		q = q.Where(s.q.SocialClient.SocialType.Eq(*r.SocialType))
	}
	if r.UserType != nil {
		q = q.Where(s.q.SocialClient.UserType.Eq(*r.UserType))
	}
	if r.ClientId != "" {
		q = q.Where(s.q.SocialClient.ClientId.Like("%" + r.ClientId + "%"))
	}

	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.SocialClient.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.SocialClient]{
		List:  list,
		Total: total,
	}, nil
}

// SendWxaSubscribeMessage 发送小程序订阅消息 (Skeleton)
func (s *SocialClientService) SendWxaSubscribeMessage(ctx context.Context, r *req.SocialWxaSubscribeMessageSendReq) error {
	// TODO: 集成真实的微信小程序 API
	// 详见 Java: SocialClientApiImpl.sendWxaSubscribeMessage
	return nil
}
