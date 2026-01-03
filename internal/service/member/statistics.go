package member

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/common"
)

// MemberStatisticsService 会员统计服务接口
type MemberStatisticsService interface {
	GetMemberSummary(ctx context.Context) (*member.MemberSummaryRespVO, error)
	GetMemberAnalyseComparisonData(ctx context.Context, beginTime, endTime time.Time) (*common.DataComparisonRespVO[interface{}], error)
	GetMemberAreaStatisticsList(ctx context.Context) ([]*member.MemberAreaStatisticsRespVO, error)
	GetMemberSexStatisticsList(ctx context.Context) ([]*member.MemberSexStatisticsRespVO, error)
	GetMemberTerminalStatisticsList(ctx context.Context) ([]*member.MemberTerminalStatisticsRespVO, error)
	GetUserCountComparison(ctx context.Context) (*common.DataComparisonRespVO[member.MemberCountRespVO], error)
	GetMemberRegisterCountList(ctx context.Context, beginTime, endTime time.Time) ([]*member.MemberRegisterCountRespVO, error)
}

// MemberStatisticsRepository 会员统计数据访问接口
type MemberStatisticsRepository interface {
	GetMemberSummary(ctx context.Context) (*member.MemberSummaryRespVO, error)
	GetMemberAreaStatisticsList(ctx context.Context) ([]*member.MemberAreaStatisticsRespVO, error)
	GetMemberSexStatisticsList(ctx context.Context) ([]*member.MemberSexStatisticsRespVO, error)
	GetMemberTerminalStatisticsList(ctx context.Context) ([]*member.MemberTerminalStatisticsRespVO, error)
	GetUserCountComparison(ctx context.Context) (*common.DataComparisonRespVO[member.MemberCountRespVO], error)
	GetMemberRegisterCountList(ctx context.Context, beginTime, endTime time.Time) ([]*member.MemberRegisterCountRespVO, error)
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
func (s *MemberStatisticsServiceImpl) GetMemberSummary(ctx context.Context) (*member.MemberSummaryRespVO, error) {
	return s.memberStatisticsRepo.GetMemberSummary(ctx)
}

// GetMemberAnalyseComparisonData 获得会员分析对比数据
func (s *MemberStatisticsServiceImpl) GetMemberAnalyseComparisonData(ctx context.Context, beginTime, endTime time.Time) (*common.DataComparisonRespVO[interface{}], error) {
	return &common.DataComparisonRespVO[interface{}]{}, nil
}

// GetMemberAreaStatisticsList 按照省份获得会员统计列表
func (s *MemberStatisticsServiceImpl) GetMemberAreaStatisticsList(ctx context.Context) ([]*member.MemberAreaStatisticsRespVO, error) {
	return s.memberStatisticsRepo.GetMemberAreaStatisticsList(ctx)
}

// GetMemberSexStatisticsList 按照性别获得会员统计列表
func (s *MemberStatisticsServiceImpl) GetMemberSexStatisticsList(ctx context.Context) ([]*member.MemberSexStatisticsRespVO, error) {
	return s.memberStatisticsRepo.GetMemberSexStatisticsList(ctx)
}

// GetMemberTerminalStatisticsList 按照终端获得会员统计列表
func (s *MemberStatisticsServiceImpl) GetMemberTerminalStatisticsList(ctx context.Context) ([]*member.MemberTerminalStatisticsRespVO, error) {
	return s.memberStatisticsRepo.GetMemberTerminalStatisticsList(ctx)
}

// GetUserCountComparison 获得用户数量对比
func (s *MemberStatisticsServiceImpl) GetUserCountComparison(ctx context.Context) (*common.DataComparisonRespVO[member.MemberCountRespVO], error) {
	return s.memberStatisticsRepo.GetUserCountComparison(ctx)
}

// GetMemberRegisterCountList 获得会员注册数量列表
func (s *MemberStatisticsServiceImpl) GetMemberRegisterCountList(ctx context.Context, beginTime, endTime time.Time) ([]*member.MemberRegisterCountRespVO, error) {
	return s.memberStatisticsRepo.GetMemberRegisterCountList(ctx, beginTime, endTime)
}
