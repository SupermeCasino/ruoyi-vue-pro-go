package pay

import (
	"context"
)

// PayWalletStatisticsService 支付钱包统计服务接口
type PayWalletStatisticsService interface {
	GetRechargePriceSummary(ctx context.Context) (int64, error)
}

// PayWalletStatisticsRepository 支付钱包统计数据访问接口
type PayWalletStatisticsRepository interface {
	GetRechargePriceSummary(ctx context.Context) (int64, error)
}

// PayWalletStatisticsServiceImpl 支付钱包统计服务实现
type PayWalletStatisticsServiceImpl struct {
	payWalletStatisticsRepo PayWalletStatisticsRepository
}

// NewPayWalletStatisticsService 创建支付钱包统计服务
func NewPayWalletStatisticsService(repo PayWalletStatisticsRepository) PayWalletStatisticsService {
	return &PayWalletStatisticsServiceImpl{
		payWalletStatisticsRepo: repo,
	}
}

// GetRechargePriceSummary 获得充值金额总和
func (s *PayWalletStatisticsServiceImpl) GetRechargePriceSummary(ctx context.Context) (int64, error) {
	return s.payWalletStatisticsRepo.GetRechargePriceSummary(ctx)
}
