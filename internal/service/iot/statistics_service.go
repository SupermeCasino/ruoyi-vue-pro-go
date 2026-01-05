package iot

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
)

type StatisticsService struct {
	productCategoryRepo ProductCategoryRepository
	productRepo         ProductRepository
	deviceRepo          DeviceRepository
	deviceMessageRepo   DeviceMessageRepository
}

func NewStatisticsService(
	productCategoryRepo ProductCategoryRepository,
	productRepo ProductRepository,
	deviceRepo DeviceRepository,
	deviceMessageRepo DeviceMessageRepository,
) *StatisticsService {
	return &StatisticsService{
		productCategoryRepo: productCategoryRepo,
		productRepo:         productRepo,
		deviceRepo:          deviceRepo,
		deviceMessageRepo:   deviceMessageRepo,
	}
}

func (s *StatisticsService) GetStatisticsSummary(ctx context.Context) (*iot.IotStatisticsSummaryRespVO, error) {
	resp := &iot.IotStatisticsSummaryRespVO{}

	// Global counts
	resp.ProductCategoryCount, _ = s.productCategoryRepo.Count(ctx, nil)
	resp.ProductCount, _ = s.productRepo.Count(ctx, nil)
	resp.DeviceCount, _ = s.deviceRepo.Count(ctx, nil)
	resp.DeviceMessageCount, _ = s.deviceMessageRepo.Count(ctx, nil)

	// Today's increments
	todayStart := time.Now().Truncate(24 * time.Hour)
	resp.ProductCategoryTodayCount, _ = s.productCategoryRepo.Count(ctx, &todayStart)
	resp.ProductTodayCount, _ = s.productRepo.Count(ctx, &todayStart)
	resp.DeviceTodayCount, _ = s.deviceRepo.Count(ctx, &todayStart)
	resp.DeviceMessageTodayCount, _ = s.deviceMessageRepo.Count(ctx, &todayStart)

	// Device status distribution
	stateCounts, _ := s.deviceRepo.GetStateCountMap(ctx)
	resp.DeviceOnlineCount = stateCounts[consts.IotDeviceStateOnline]
	resp.DeviceOfflineCount = stateCounts[consts.IotDeviceStateOffline]
	resp.DeviceInactiveCount = stateCounts[consts.IotDeviceStateInactive]

	// Category device counts
	resp.ProductCategoryDeviceCounts, _ = s.productCategoryRepo.GetProductCategoryDeviceCountMap(ctx)

	return resp, nil
}

func (s *StatisticsService) GetDeviceMessageSummaryByDate(ctx context.Context, r *iot.IotStatisticsDeviceMessageReqVO) ([]*iot.IotStatisticsDeviceMessageSummaryByDateRespVO, error) {
	var startTime, endTime *time.Time
	if len(r.Times) == 2 {
		startTime = r.Times[0]
		endTime = r.Times[1]
	}
	return s.deviceMessageRepo.GetSummaryByDate(ctx, r.Interval, startTime, endTime)
}
