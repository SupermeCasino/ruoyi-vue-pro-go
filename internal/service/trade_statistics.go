package service

import (
	"context"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/statistics"
	"time"
)

// ============ TradeStatisticsService 接口定义 ============

// TradeStatisticsService 交易统计服务接口
type TradeStatisticsService interface {
	GetTradeSummaryByDays(ctx context.Context, days int) (*resp.TradeSummaryItemVO, error)
	GetTradeSummaryByMonths(ctx context.Context, months int) (*resp.TradeSummaryItemVO, error)
	GetTradeStatisticsAnalyse(ctx context.Context, beginTime, endTime time.Time) (*resp.DataComparisonRespVO[resp.TradeTrendSummaryRespVO], error)
	GetTradeStatisticsList(ctx context.Context, beginTime, endTime time.Time) ([]*resp.TradeTrendSummaryRespVO, error)
}

// TradeOrderStatisticsService 交易订单统计服务接口
type TradeOrderStatisticsService interface {
	GetCountByStatusAndDeliveryType(ctx context.Context, status int, deliveryType int) (int64, error)
	GetPayUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error)
	GetOrderPayPrice(ctx context.Context, beginTime, endTime time.Time) (int64, error)
	GetOrderUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error)
	GetOrderComparison(ctx context.Context) (*resp.DataComparisonRespVO[resp.TradeOrderSummaryRespVO], error)
	GetOrderCountTrendComparison(ctx context.Context) ([]*resp.DataComparisonRespVO[resp.TradeOrderTrendRespVO], error)
}

// AfterSaleStatisticsService 售后统计服务接口
type AfterSaleStatisticsService interface {
	GetCountByStatus(ctx context.Context, status int) (int64, error)
}

// BrokerageStatisticsService 佣金统计服务接口
type BrokerageStatisticsService interface {
	GetWithdrawCountByStatus(ctx context.Context, status int) (int64, error)
}

// ============ TradeStatisticsService 实现 ============

// TradeStatisticsModel 交易统计模型（用于 Repository 和 Service 间传递）- 严格对齐 Java TradeStatisticsDO
type TradeStatisticsModel struct {
	StatisticsTime           time.Time // 统计日期
	OrderCreateCount         int       // 创建订单数
	OrderPayCount            int       // 支付订单商品数
	OrderPayPrice            int       // 总支付金额(分)
	AfterSaleCount           int       // 退款订单数
	AfterSaleRefundPrice     int       // 总退款金额(分)
	BrokerageSettlementPrice int       // 佣金金额已结算(分)
	WalletPayPrice           int       // 总支付金额余额(分)
	RechargePayCount         int       // 充值订单数
	RechargePayPrice         int       // 充值金额(分)
	RechargeRefundCount      int       // 充值退款订单数
	RechargeRefundPrice      int       // 充值退款金额(分)
}

// TradeStatisticsRepository 交易统计数据访问接口
type TradeStatisticsRepository interface {
	GetByDateRange(ctx context.Context, beginTime, endTime time.Time) (*TradeStatisticsModel, error)
	GetByMonthRange(ctx context.Context, beginTime, endTime time.Time) (*TradeStatisticsModel, error)
	GetListByDateRange(ctx context.Context, beginTime, endTime time.Time) ([]*TradeStatisticsModel, error)
	Insert(ctx context.Context, stats *TradeStatisticsModel) error
}

// TradeStatisticsServiceImpl 交易统计服务实现
type TradeStatisticsServiceImpl struct {
	tradeStatisticsRepo         TradeStatisticsRepository
	tradeOrderStatisticsService TradeOrderStatisticsService
	afterSaleStatisticsService  AfterSaleStatisticsService
	brokerageStatisticsService  BrokerageStatisticsService
}

// NewTradeStatisticsService 创建交易统计服务
func NewTradeStatisticsService(
	tradeStatisticsRepo TradeStatisticsRepository,
	tradeOrderStatisticsService TradeOrderStatisticsService,
	afterSaleStatisticsService AfterSaleStatisticsService,
	brokerageStatisticsService BrokerageStatisticsService,
) TradeStatisticsService {
	return &TradeStatisticsServiceImpl{
		tradeStatisticsRepo:         tradeStatisticsRepo,
		tradeOrderStatisticsService: tradeOrderStatisticsService,
		afterSaleStatisticsService:  afterSaleStatisticsService,
		brokerageStatisticsService:  brokerageStatisticsService,
	}
}

// GetTradeSummaryByDays 获得指定天数的交易统计摘要
func (s *TradeStatisticsServiceImpl) GetTradeSummaryByDays(ctx context.Context, days int) (*resp.TradeSummaryItemVO, error) {
	targetDate := time.Now().AddDate(0, 0, days)
	beginTime := statistics.BeginOfDay(targetDate)
	endTime := statistics.EndOfDay(targetDate)

	stats, err := s.tradeStatisticsRepo.GetByDateRange(ctx, beginTime, endTime)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return &resp.TradeSummaryItemVO{}, nil
	}

	return &resp.TradeSummaryItemVO{
		OrderCreateCount:         int64(stats.OrderCreateCount),
		OrderPayCount:            int64(stats.OrderPayCount),
		OrderPayPrice:            int64(stats.OrderPayPrice),
		AfterSaleCount:           int64(stats.AfterSaleCount),
		AfterSaleRefundPrice:     int64(stats.AfterSaleRefundPrice),
		BrokerageSettlementPrice: int64(stats.BrokerageSettlementPrice),
		WalletPayPrice:           int64(stats.WalletPayPrice),
		RechargePayCount:         int64(stats.RechargePayCount),
		RechargePayPrice:         int64(stats.RechargePayPrice),
		RechargeRefundCount:      int64(stats.RechargeRefundCount),
		RechargeRefundPrice:      int64(stats.RechargeRefundPrice),
	}, nil
}

// GetTradeSummaryByMonths 获得指定月份的交易统计摘要
func (s *TradeStatisticsServiceImpl) GetTradeSummaryByMonths(ctx context.Context, months int) (*resp.TradeSummaryItemVO, error) {
	monthDate := time.Now().AddDate(0, months, 0)
	beginTime := statistics.BeginOfMonth(monthDate)
	endTime := statistics.EndOfMonth(monthDate)

	stats, err := s.tradeStatisticsRepo.GetByMonthRange(ctx, beginTime, endTime)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return &resp.TradeSummaryItemVO{}, nil
	}

	return &resp.TradeSummaryItemVO{
		OrderCreateCount:         int64(stats.OrderCreateCount),
		OrderPayCount:            int64(stats.OrderPayCount),
		OrderPayPrice:            int64(stats.OrderPayPrice),
		AfterSaleCount:           int64(stats.AfterSaleCount),
		AfterSaleRefundPrice:     int64(stats.AfterSaleRefundPrice),
		BrokerageSettlementPrice: int64(stats.BrokerageSettlementPrice),
		WalletPayPrice:           int64(stats.WalletPayPrice),
		RechargePayCount:         int64(stats.RechargePayCount),
		RechargePayPrice:         int64(stats.RechargePayPrice),
		RechargeRefundCount:      int64(stats.RechargeRefundCount),
		RechargeRefundPrice:      int64(stats.RechargeRefundPrice),
	}, nil
}

// GetTradeStatisticsAnalyse 获得交易统计分析
func (s *TradeStatisticsServiceImpl) GetTradeStatisticsAnalyse(ctx context.Context, beginTime, endTime time.Time) (*resp.DataComparisonRespVO[resp.TradeTrendSummaryRespVO], error) {
	currentStats, err := s.tradeStatisticsRepo.GetByDateRange(ctx, beginTime, endTime)
	if err != nil {
		return nil, err
	}

	duration := endTime.Sub(beginTime)
	referenceBeginTime := beginTime.Add(-duration)
	referenceEndTime := beginTime

	referenceStats, err := s.tradeStatisticsRepo.GetByDateRange(ctx, referenceBeginTime, referenceEndTime)
	if err != nil {
		return nil, err
	}

	result := &resp.DataComparisonRespVO[resp.TradeTrendSummaryRespVO]{
		Summary: &resp.TradeTrendSummaryRespVO{
			StatisticsTime:           currentStats.StatisticsTime,
			OrderCreateCount:         int64(currentStats.OrderCreateCount),
			OrderPayCount:            int64(currentStats.OrderPayCount),
			OrderPayPrice:            int64(currentStats.OrderPayPrice),
			AfterSaleCount:           int64(currentStats.AfterSaleCount),
			AfterSaleRefundPrice:     int64(currentStats.AfterSaleRefundPrice),
			BrokerageSettlementPrice: int64(currentStats.BrokerageSettlementPrice),
			WalletPayPrice:           int64(currentStats.WalletPayPrice),
			RechargePayCount:         int64(currentStats.RechargePayCount),
			RechargePayPrice:         int64(currentStats.RechargePayPrice),
			RechargeRefundCount:      int64(currentStats.RechargeRefundCount),
			RechargeRefundPrice:      int64(currentStats.RechargeRefundPrice),
		},
		Comparison: &resp.TradeTrendSummaryRespVO{
			StatisticsTime:           referenceStats.StatisticsTime,
			OrderCreateCount:         int64(referenceStats.OrderCreateCount),
			OrderPayCount:            int64(referenceStats.OrderPayCount),
			OrderPayPrice:            int64(referenceStats.OrderPayPrice),
			AfterSaleCount:           int64(referenceStats.AfterSaleCount),
			AfterSaleRefundPrice:     int64(referenceStats.AfterSaleRefundPrice),
			BrokerageSettlementPrice: int64(referenceStats.BrokerageSettlementPrice),
			WalletPayPrice:           int64(referenceStats.WalletPayPrice),
			RechargePayCount:         int64(referenceStats.RechargePayCount),
			RechargePayPrice:         int64(referenceStats.RechargePayPrice),
			RechargeRefundCount:      int64(referenceStats.RechargeRefundCount),
			RechargeRefundPrice:      int64(referenceStats.RechargeRefundPrice),
		},
	}

	return result, nil
}

// GetTradeStatisticsList 获得交易统计列表
func (s *TradeStatisticsServiceImpl) GetTradeStatisticsList(ctx context.Context, beginTime, endTime time.Time) ([]*resp.TradeTrendSummaryRespVO, error) {
	statsList, err := s.tradeStatisticsRepo.GetListByDateRange(ctx, beginTime, endTime)
	if err != nil {
		return nil, err
	}

	result := make([]*resp.TradeTrendSummaryRespVO, 0, len(statsList))
	for _, stats := range statsList {
		result = append(result, &resp.TradeTrendSummaryRespVO{
			StatisticsTime:           stats.StatisticsTime,
			OrderCreateCount:         int64(stats.OrderCreateCount),
			OrderPayCount:            int64(stats.OrderPayCount),
			OrderPayPrice:            int64(stats.OrderPayPrice),
			AfterSaleCount:           int64(stats.AfterSaleCount),
			AfterSaleRefundPrice:     int64(stats.AfterSaleRefundPrice),
			BrokerageSettlementPrice: int64(stats.BrokerageSettlementPrice),
			WalletPayPrice:           int64(stats.WalletPayPrice),
			RechargePayCount:         int64(stats.RechargePayCount),
			RechargePayPrice:         int64(stats.RechargePayPrice),
			RechargeRefundCount:      int64(stats.RechargeRefundCount),
			RechargeRefundPrice:      int64(stats.RechargeRefundPrice),
		})
	}

	return result, nil
}

// ============ TradeOrderStatisticsService 实现 ============

// TradeOrderStatisticsRepository 交易订单统计数据访问接口
type TradeOrderStatisticsRepository interface {
	GetCountByStatusAndDeliveryType(ctx context.Context, status int, deliveryType int) (int64, error)
	GetPayUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error)
	GetOrderPayPrice(ctx context.Context, beginTime, endTime time.Time) (int64, error)
	GetOrderUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error)
}

// TradeOrderStatisticsServiceImplV2 交易订单统计服务实现
type TradeOrderStatisticsServiceImplV2 struct {
	tradeOrderStatisticsRepo TradeOrderStatisticsRepository
}

// NewTradeOrderStatisticsServiceV2 创建交易订单统计服务
func NewTradeOrderStatisticsServiceV2(repo TradeOrderStatisticsRepository) TradeOrderStatisticsService {
	return &TradeOrderStatisticsServiceImplV2{
		tradeOrderStatisticsRepo: repo,
	}
}

// GetCountByStatusAndDeliveryType 获得指定状态和配送方式的订单数
func (s *TradeOrderStatisticsServiceImplV2) GetCountByStatusAndDeliveryType(ctx context.Context, status int, deliveryType int) (int64, error) {
	return s.tradeOrderStatisticsRepo.GetCountByStatusAndDeliveryType(ctx, status, deliveryType)
}

// GetPayUserCount 获得支付用户数
func (s *TradeOrderStatisticsServiceImplV2) GetPayUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	return s.tradeOrderStatisticsRepo.GetPayUserCount(ctx, beginTime, endTime)
}

// GetOrderPayPrice 获得订单支付金额
func (s *TradeOrderStatisticsServiceImplV2) GetOrderPayPrice(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	return s.tradeOrderStatisticsRepo.GetOrderPayPrice(ctx, beginTime, endTime)
}

// GetOrderUserCount 获得下单用户数
func (s *TradeOrderStatisticsServiceImplV2) GetOrderUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	return s.tradeOrderStatisticsRepo.GetOrderUserCount(ctx, beginTime, endTime)
}

// GetOrderComparison 获得订单数量对比
func (s *TradeOrderStatisticsServiceImplV2) GetOrderComparison(ctx context.Context) (*resp.DataComparisonRespVO[resp.TradeOrderSummaryRespVO], error) {
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayBegin := statistics.BeginOfDay(yesterday)
	yesterdayEnd := statistics.EndOfDay(yesterday)

	yesterdayOrderCount, err := s.tradeOrderStatisticsRepo.GetOrderUserCount(ctx, yesterdayBegin, yesterdayEnd)
	if err != nil {
		return nil, err
	}

	yesterdayPayPrice, err := s.tradeOrderStatisticsRepo.GetOrderPayPrice(ctx, yesterdayBegin, yesterdayEnd)
	if err != nil {
		return nil, err
	}

	beforeYesterday := time.Now().AddDate(0, 0, -2)
	beforeYesterdayBegin := statistics.BeginOfDay(beforeYesterday)
	beforeYesterdayEnd := statistics.EndOfDay(beforeYesterday)

	beforeYesterdayOrderCount, err := s.tradeOrderStatisticsRepo.GetOrderUserCount(ctx, beforeYesterdayBegin, beforeYesterdayEnd)
	if err != nil {
		return nil, err
	}

	beforeYesterdayPayPrice, err := s.tradeOrderStatisticsRepo.GetOrderPayPrice(ctx, beforeYesterdayBegin, beforeYesterdayEnd)
	if err != nil {
		return nil, err
	}

	result := &resp.DataComparisonRespVO[resp.TradeOrderSummaryRespVO]{
		Summary: &resp.TradeOrderSummaryRespVO{
			OrderCount: yesterdayOrderCount,
			PayPrice:   yesterdayPayPrice,
		},
		Comparison: &resp.TradeOrderSummaryRespVO{
			OrderCount: beforeYesterdayOrderCount,
			PayPrice:   beforeYesterdayPayPrice,
		},
	}

	return result, nil
}

// GetOrderCountTrendComparison 获得订单量趋势对比
func (s *TradeOrderStatisticsServiceImplV2) GetOrderCountTrendComparison(ctx context.Context) ([]*resp.DataComparisonRespVO[resp.TradeOrderTrendRespVO], error) {
	return []*resp.DataComparisonRespVO[resp.TradeOrderTrendRespVO]{}, nil
}

// ============ AfterSaleStatisticsService 实现 ============

// AfterSaleStatisticsRepository 售后统计数据访问接口
type AfterSaleStatisticsRepository interface {
	GetCountByStatus(ctx context.Context, status int) (int64, error)
}

// AfterSaleStatisticsServiceImpl 售后统计服务实现
type AfterSaleStatisticsServiceImpl struct {
	afterSaleStatisticsRepo AfterSaleStatisticsRepository
}

// NewAfterSaleStatisticsService 创建售后统计服务
func NewAfterSaleStatisticsService(repo AfterSaleStatisticsRepository) AfterSaleStatisticsService {
	return &AfterSaleStatisticsServiceImpl{
		afterSaleStatisticsRepo: repo,
	}
}

// GetCountByStatus 获得指定状态的售后数
func (s *AfterSaleStatisticsServiceImpl) GetCountByStatus(ctx context.Context, status int) (int64, error) {
	return s.afterSaleStatisticsRepo.GetCountByStatus(ctx, status)
}

// ============ BrokerageStatisticsService 实现 ============

// BrokerageStatisticsRepository 佣金统计数据访问接口
type BrokerageStatisticsRepository interface {
	GetWithdrawCountByStatus(ctx context.Context, status int) (int64, error)
}

// BrokerageStatisticsServiceImpl 佣金统计服务实现
type BrokerageStatisticsServiceImpl struct {
	brokerageStatisticsRepo BrokerageStatisticsRepository
}

// NewBrokerageStatisticsService 创建佣金统计服务
func NewBrokerageStatisticsService(repo BrokerageStatisticsRepository) BrokerageStatisticsService {
	return &BrokerageStatisticsServiceImpl{
		brokerageStatisticsRepo: repo,
	}
}

// GetWithdrawCountByStatus 获得指定状态的提现数
func (s *BrokerageStatisticsServiceImpl) GetWithdrawCountByStatus(ctx context.Context, status int) (int64, error) {
	return s.brokerageStatisticsRepo.GetWithdrawCountByStatus(ctx, status)
}
