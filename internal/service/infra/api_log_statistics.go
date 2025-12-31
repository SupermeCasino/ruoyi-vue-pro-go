package infra

import (
	"context"
)

// ApiAccessLogStatisticsService API 访问日志统计服务接口
type ApiAccessLogStatisticsService interface {
	GetIpCount(ctx context.Context, userType int, beginTime, endTime interface{}) (int64, error)
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
