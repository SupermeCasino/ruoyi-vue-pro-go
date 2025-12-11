package member

import (
	"backend-go/internal/api/req"
	"backend-go/internal/model/member"
	"backend-go/internal/repo/query"
	"context"
)

type MemberConfigService struct {
	q *query.Query
}

func NewMemberConfigService(q *query.Query) *MemberConfigService {
	return &MemberConfigService{q: q}
}

// SaveConfig 保存会员配置
func (s *MemberConfigService) SaveConfig(ctx context.Context, r *req.MemberConfigSaveReq) error {
	config, err := s.GetConfig(ctx)
	if err != nil {
		return err
	}

	newConfig := &member.MemberConfig{
		PointTradeDeductEnable:    r.PointTradeDeductEnable,
		PointTradeDeductUnitPrice: r.PointTradeDeductUnitPrice,
		PointTradeDeductMaxPrice:  r.PointTradeDeductMaxPrice,
		PointTradeGivePoint:       r.PointTradeGivePoint,
		ID:                        0,
	}

	if config != nil {
		newConfig.ID = config.ID
		_, err := s.q.MemberConfig.WithContext(ctx).Where(s.q.MemberConfig.ID.Eq(config.ID)).Updates(newConfig)
		return err
	}

	return s.q.MemberConfig.WithContext(ctx).Create(newConfig)
}

// GetConfig 获得会员配置
func (s *MemberConfigService) GetConfig(ctx context.Context) (*member.MemberConfig, error) {
	config, err := s.q.MemberConfig.WithContext(ctx).First()
	if err != nil {
		// If record not found, return nil is acceptable for singleton config, or return default.
		// Usually GORM returns ErrRecordNotFound.
		// For config, we might want to return nil without error if not initialized.
		// But let's follow the Java logic: `CollectionUtils.getFirst(list)` implies it might be null.
		return nil, nil
	}
	return config, nil
}
