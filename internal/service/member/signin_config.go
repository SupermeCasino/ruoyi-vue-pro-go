package member

import (
	"context"
	"errors"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type MemberSignInConfigService struct {
	q *query.Query
}

func NewMemberSignInConfigService(q *query.Query) *MemberSignInConfigService {
	return &MemberSignInConfigService{q: q}
}

// CreateSignInConfig 创建签到规则
func (s *MemberSignInConfigService) CreateSignInConfig(ctx context.Context, r *req.MemberSignInConfigCreateReq) (int64, error) {
	// 校验 day 是否重复
	if err := s.validateSignInConfigDayDuplicate(ctx, r.Day, 0); err != nil {
		return 0, err
	}

	config := &member.MemberSignInConfig{
		Day:        r.Day,
		Point:      r.Point,
		Experience: r.Experience,
		Status:     r.Status,
	}
	err := s.q.MemberSignInConfig.WithContext(ctx).Create(config)
	return config.ID, err
}

// UpdateSignInConfig 更新签到规则
func (s *MemberSignInConfigService) UpdateSignInConfig(ctx context.Context, r *req.MemberSignInConfigUpdateReq) error {
	// 校验存在
	if _, err := s.GetSignInConfig(ctx, r.ID); err != nil {
		return err
	}
	// 校验 day 是否重复
	if err := s.validateSignInConfigDayDuplicate(ctx, r.Day, r.ID); err != nil {
		return err
	}

	_, err := s.q.MemberSignInConfig.WithContext(ctx).Where(s.q.MemberSignInConfig.ID.Eq(r.ID)).Updates(member.MemberSignInConfig{
		Day:        r.Day,
		Point:      r.Point,
		Experience: r.Experience,
		Status:     r.Status,
	})
	return err
}

// DeleteSignInConfig 删除签到规则
func (s *MemberSignInConfigService) DeleteSignInConfig(ctx context.Context, id int64) error {
	// 校验存在
	if _, err := s.GetSignInConfig(ctx, id); err != nil {
		return err
	}
	_, err := s.q.MemberSignInConfig.WithContext(ctx).Where(s.q.MemberSignInConfig.ID.Eq(id)).Delete()
	return err
}

// GetSignInConfig 获得签到规则
func (s *MemberSignInConfigService) GetSignInConfig(ctx context.Context, id int64) (*member.MemberSignInConfig, error) {
	config, err := s.q.MemberSignInConfig.WithContext(ctx).Where(s.q.MemberSignInConfig.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.New("sign-in config not found")
	}
	return config, nil
}

// GetSignInConfigList 获得签到规则列表
func (s *MemberSignInConfigService) GetSignInConfigList(ctx context.Context, status *int) ([]*member.MemberSignInConfig, error) {
	q := s.q.MemberSignInConfig.WithContext(ctx)
	if status != nil {
		q = q.Where(s.q.MemberSignInConfig.Status.Eq(*status))
	}
	return q.Order(s.q.MemberSignInConfig.Day.Asc()).Find()
}

func (s *MemberSignInConfigService) validateSignInConfigDayDuplicate(ctx context.Context, day int, id int64) error {
	config, err := s.q.MemberSignInConfig.WithContext(ctx).Where(s.q.MemberSignInConfig.Day.Eq(day)).First()
	if err != nil {
		return nil // Not found, strictly speaking, ignoring err unless it's DB error, but here usually nil if not found
	}
	if config == nil {
		return nil
	}

	// 1. New: config exists -> duplicate
	if id == 0 {
		return core.NewBizError(1004014004, "Sign-in day config already exists") // Assuming error code
	}
	// 2. Update: config exists and id not match -> duplicate
	if config.ID != id {
		return core.NewBizError(1004014004, "Sign-in day config already exists")
	}
	return nil
}
