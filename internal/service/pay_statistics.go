package service

import (
	"context"
)

// PayWalletStatisticsService 支付钱包统计服务接口
type PayWalletStatisticsService interface {
	GetRechargePriceSummary(ctx context.Context) (int64, error)
}

// ApiAccessLogStatisticsService API 访问日志统计服务接口
type ApiAccessLogStatisticsService interface {
	GetIpCount(ctx context.Context, userType int, beginTime, endTime interface{}) (int64, error)
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

// ApiAccessLogStatisticsRepository API 访问日志统计数据访问接口
type ApiAccessLogStatisticsRepository interface {
	GetIpCount(ctx context.Context, userType int, beginTime, endTime interface{}) (int64, error)
}

// ApiAccessLogStatisticsServiceImpl API 访问日志统计服务实现
type ApiAccessLogStatisticsServiceImpl struct {
	apiAccessLogStatisticsRepo ApiAccessLogStatisticsRepository
}

// NewApiAccessLogStatisticsService 创建 API 访问日志统计服务
func NewApiAccessLogStatisticsService(repo ApiAccessLogStatisticsRepository) ApiAccessLogStatisticsService {
	return &ApiAccessLogStatisticsServiceImpl{
		apiAccessLogStatisticsRepo: repo,
	}
}

// GetIpCount 获得 IP 访问数
func (s *ApiAccessLogStatisticsServiceImpl) GetIpCount(ctx context.Context, userType int, beginTime, endTime interface{}) (int64, error) {
	return s.apiAccessLogStatisticsRepo.GetIpCount(ctx, userType, beginTime, endTime)
}
