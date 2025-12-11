package repo

import (
	"backend-go/internal/api/resp"
	"backend-go/internal/model/member"
	"backend-go/internal/repo/query"
	"backend-go/internal/service"
	"context"
	"time"

	"gorm.io/gorm"
)

// ============ MemberStatisticsRepository 实现 ============

// MemberStatisticsRepositoryImpl 会员统计 Repository 实现
type MemberStatisticsRepositoryImpl struct {
	q  *query.Query
	db *gorm.DB
}

// NewMemberStatisticsRepository 创建会员统计 Repository
func NewMemberStatisticsRepository(q *query.Query, db *gorm.DB) service.MemberStatisticsRepository {
	return &MemberStatisticsRepositoryImpl{q: q, db: db}
}

// GetMemberSummary 获得会员统计摘要
func (r *MemberStatisticsRepositoryImpl) GetMemberSummary(ctx context.Context) (*resp.MemberSummaryRespVO, error) {
	var totalCount int64
	err := r.db.WithContext(ctx).Model(&member.MemberUser{}).
		Where("deleted = ?", false).
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	// 今日注册
	now := time.Now()
	todayBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayBegin.Add(24 * time.Hour)
	var todayCount int64
	err = r.db.WithContext(ctx).Model(&member.MemberUser{}).
		Where("deleted = ?", false).
		Where("create_time BETWEEN ? AND ?", todayBegin, todayEnd).
		Count(&todayCount).Error
	if err != nil {
		return nil, err
	}

	return &resp.MemberSummaryRespVO{
		TotalUserCount: totalCount,
		RegisterCount:  todayCount,
	}, nil
}

// GetMemberAreaStatisticsList 按照省份获得会员统计列表
func (r *MemberStatisticsRepositoryImpl) GetMemberAreaStatisticsList(ctx context.Context) ([]*resp.MemberAreaStatisticsRespVO, error) {
	var results []struct {
		AreaID int   `gorm:"column:area_id"`
		Count  int64 `gorm:"column:count"`
	}

	err := r.db.WithContext(ctx).Model(&member.MemberUser{}).
		Select("area_id, COUNT(*) as count").
		Where("deleted = ?", false).
		Group("area_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make([]*resp.MemberAreaStatisticsRespVO, 0, len(results))
	for _, r := range results {
		stats = append(stats, &resp.MemberAreaStatisticsRespVO{
			AreaName:  "", // 需要从地区表获取，暂空
			UserCount: r.Count,
		})
	}
	return stats, nil
}

// GetMemberSexStatisticsList 按照性别获得会员统计列表
func (r *MemberStatisticsRepositoryImpl) GetMemberSexStatisticsList(ctx context.Context) ([]*resp.MemberSexStatisticsRespVO, error) {
	var results []struct {
		Sex   int   `gorm:"column:sex"`
		Count int64 `gorm:"column:count"`
	}

	err := r.db.WithContext(ctx).Model(&member.MemberUser{}).
		Select("sex, COUNT(*) as count").
		Where("deleted = ?", false).
		Group("sex").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 性别映射
	sexMap := map[int]string{
		0: "未知",
		1: "男",
		2: "女",
	}

	stats := make([]*resp.MemberSexStatisticsRespVO, 0, len(results))
	for _, r := range results {
		sexName := sexMap[r.Sex]
		if sexName == "" {
			sexName = "未知"
		}
		stats = append(stats, &resp.MemberSexStatisticsRespVO{
			SexName:   sexName,
			UserCount: r.Count,
		})
	}
	return stats, nil
}

// GetMemberTerminalStatisticsList 按照终端获得会员统计列表
func (r *MemberStatisticsRepositoryImpl) GetMemberTerminalStatisticsList(ctx context.Context) ([]*resp.MemberTerminalStatisticsRespVO, error) {
	var results []struct {
		Terminal int   `gorm:"column:register_terminal"`
		Count    int64 `gorm:"column:count"`
	}

	err := r.db.WithContext(ctx).Model(&member.MemberUser{}).
		Select("register_terminal, COUNT(*) as count").
		Where("deleted = ?", false).
		Group("register_terminal").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 终端映射
	terminalMap := map[int]string{
		1: "微信小程序",
		2: "微信公众号",
		3: "H5",
		4: "App",
	}

	stats := make([]*resp.MemberTerminalStatisticsRespVO, 0, len(results))
	for _, r := range results {
		terminalName := terminalMap[r.Terminal]
		if terminalName == "" {
			terminalName = "其他"
		}
		stats = append(stats, &resp.MemberTerminalStatisticsRespVO{
			TerminalName: terminalName,
			UserCount:    r.Count,
		})
	}
	return stats, nil
}

// GetUserCountComparison 获得用户数量对比
func (r *MemberStatisticsRepositoryImpl) GetUserCountComparison(ctx context.Context) (*resp.DataComparisonRespVO[resp.MemberCountRespVO], error) {
	now := time.Now()

	// 今日数据
	todayBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayBegin.Add(24 * time.Hour)
	var todayCount int64
	err := r.db.WithContext(ctx).Model(&member.MemberUser{}).
		Where("deleted = ?", false).
		Where("create_time BETWEEN ? AND ?", todayBegin, todayEnd).
		Count(&todayCount).Error
	if err != nil {
		return nil, err
	}

	// 昨日数据
	yesterdayBegin := todayBegin.AddDate(0, 0, -1)
	yesterdayEnd := todayBegin
	var yesterdayCount int64
	err = r.db.WithContext(ctx).Model(&member.MemberUser{}).
		Where("deleted = ?", false).
		Where("create_time BETWEEN ? AND ?", yesterdayBegin, yesterdayEnd).
		Count(&yesterdayCount).Error
	if err != nil {
		return nil, err
	}

	return &resp.DataComparisonRespVO[resp.MemberCountRespVO]{
		Summary:    &resp.MemberCountRespVO{UserCount: todayCount},
		Comparison: &resp.MemberCountRespVO{UserCount: yesterdayCount},
	}, nil
}

// GetMemberRegisterCountList 获得会员注册数量列表
func (r *MemberStatisticsRepositoryImpl) GetMemberRegisterCountList(ctx context.Context, beginTime, endTime time.Time) ([]*resp.MemberRegisterCountRespVO, error) {
	// 需要按日期分组统计，暂返回空
	return []*resp.MemberRegisterCountRespVO{}, nil
}

// ============ ApiAccessLogStatisticsRepository 实现 ============

// ApiAccessLogStatisticsRepositoryImpl API访问日志统计 Repository 实现
type ApiAccessLogStatisticsRepositoryImpl struct {
	q *query.Query
}

// NewApiAccessLogStatisticsRepository 创建API访问日志统计 Repository
func NewApiAccessLogStatisticsRepository(q *query.Query) service.ApiAccessLogStatisticsRepository {
	return &ApiAccessLogStatisticsRepositoryImpl{q: q}
}

// GetIpCount 获取独立IP数量（UV）
func (r *ApiAccessLogStatisticsRepositoryImpl) GetIpCount(ctx context.Context, userType int, beginTime, endTime interface{}) (int64, error) {
	t := r.q.InfraApiAccessLog

	q := t.WithContext(ctx)
	if userType >= 0 {
		q = q.Where(t.UserType.Eq(userType))
	}
	// beginTime 和 endTime 类型转换
	if bt, ok := beginTime.(time.Time); ok {
		if et, ok := endTime.(time.Time); ok {
			q = q.Where(t.BeginTime.Between(bt, et))
		}
	}

	return q.Distinct(t.UserIP).Count()
}

// ============ PayWalletStatisticsRepository 实现 ============

// PayWalletStatisticsRepositoryImpl 支付钱包统计 Repository 实现
type PayWalletStatisticsRepositoryImpl struct {
	q *query.Query
}

// NewPayWalletStatisticsRepository 创建支付钱包统计 Repository
func NewPayWalletStatisticsRepository(q *query.Query) service.PayWalletStatisticsRepository {
	return &PayWalletStatisticsRepositoryImpl{q: q}
}

// GetRechargePriceSummary 获取充值金额汇总
func (r *PayWalletStatisticsRepositoryImpl) GetRechargePriceSummary(ctx context.Context) (int64, error) {
	// pay_wallet_recharge 表未生成 gorm gen，返回 0
	return 0, nil
}
