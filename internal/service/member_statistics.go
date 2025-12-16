package service

import (
	"context"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"time"
)

// MemberStatisticsService 会员统计服务接口
type MemberStatisticsService interface {
	GetMemberSummary(ctx context.Context) (*resp.MemberSummaryRespVO, error)
	GetMemberAnalyseComparisonData(ctx context.Context, beginTime, endTime time.Time) (*resp.DataComparisonRespVO[interface{}], error)
	GetMemberAreaStatisticsList(ctx context.Context) ([]*resp.MemberAreaStatisticsRespVO, error)
	GetMemberSexStatisticsList(ctx context.Context) ([]*resp.MemberSexStatisticsRespVO, error)
	GetMemberTerminalStatisticsList(ctx context.Context) ([]*resp.MemberTerminalStatisticsRespVO, error)
	GetUserCountComparison(ctx context.Context) (*resp.DataComparisonRespVO[resp.MemberCountRespVO], error)
	GetMemberRegisterCountList(ctx context.Context, beginTime, endTime time.Time) ([]*resp.MemberRegisterCountRespVO, error)
}

// MemberStatisticsRepository 会员统计数据访问接口
type MemberStatisticsRepository interface {
	GetMemberSummary(ctx context.Context) (*resp.MemberSummaryRespVO, error)
	GetMemberAreaStatisticsList(ctx context.Context) ([]*resp.MemberAreaStatisticsRespVO, error)
	GetMemberSexStatisticsList(ctx context.Context) ([]*resp.MemberSexStatisticsRespVO, error)
	GetMemberTerminalStatisticsList(ctx context.Context) ([]*resp.MemberTerminalStatisticsRespVO, error)
	GetUserCountComparison(ctx context.Context) (*resp.DataComparisonRespVO[resp.MemberCountRespVO], error)
	GetMemberRegisterCountList(ctx context.Context, beginTime, endTime time.Time) ([]*resp.MemberRegisterCountRespVO, error)
}

// MemberStatisticsServiceImpl 会员统计服务实现
type MemberStatisticsServiceImpl struct {
	memberStatisticsRepo MemberStatisticsRepository
}

// NewMemberStatisticsService 创建会员统计服务
func NewMemberStatisticsService(repo MemberStatisticsRepository) MemberStatisticsService {
	return &MemberStatisticsServiceImpl{
		memberStatisticsRepo: repo,
	}
}

// GetMemberSummary 获得会员统计摘要
func (s *MemberStatisticsServiceImpl) GetMemberSummary(ctx context.Context) (*resp.MemberSummaryRespVO, error) {
	return s.memberStatisticsRepo.GetMemberSummary(ctx)
}

// GetMemberAnalyseComparisonData 获得会员分析对比数据
func (s *MemberStatisticsServiceImpl) GetMemberAnalyseComparisonData(ctx context.Context, beginTime, endTime time.Time) (*resp.DataComparisonRespVO[interface{}], error) {
	return &resp.DataComparisonRespVO[interface{}]{}, nil
}

// GetMemberAreaStatisticsList 按照省份获得会员统计列表
func (s *MemberStatisticsServiceImpl) GetMemberAreaStatisticsList(ctx context.Context) ([]*resp.MemberAreaStatisticsRespVO, error) {
	return s.memberStatisticsRepo.GetMemberAreaStatisticsList(ctx)
}

// GetMemberSexStatisticsList 按照性别获得会员统计列表
func (s *MemberStatisticsServiceImpl) GetMemberSexStatisticsList(ctx context.Context) ([]*resp.MemberSexStatisticsRespVO, error) {
	return s.memberStatisticsRepo.GetMemberSexStatisticsList(ctx)
}

// GetMemberTerminalStatisticsList 按照终端获得会员统计列表
func (s *MemberStatisticsServiceImpl) GetMemberTerminalStatisticsList(ctx context.Context) ([]*resp.MemberTerminalStatisticsRespVO, error) {
	return s.memberStatisticsRepo.GetMemberTerminalStatisticsList(ctx)
}

// GetUserCountComparison 获得用户数量对比
func (s *MemberStatisticsServiceImpl) GetUserCountComparison(ctx context.Context) (*resp.DataComparisonRespVO[resp.MemberCountRespVO], error) {
	return s.memberStatisticsRepo.GetUserCountComparison(ctx)
}

// GetMemberRegisterCountList 获得会员注册数量列表
func (s *MemberStatisticsServiceImpl) GetMemberRegisterCountList(ctx context.Context, beginTime, endTime time.Time) ([]*resp.MemberRegisterCountRespVO, error) {
	return s.memberStatisticsRepo.GetMemberRegisterCountList(ctx, beginTime, endTime)
}
