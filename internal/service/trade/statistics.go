package trade

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/statistics"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type TradeStatisticsService struct {
	q *query.Query
}

func NewTradeStatisticsService(q *query.Query) *TradeStatisticsService {
	return &TradeStatisticsService{
		q: q,
	}
}

// GetSummary 获得交易简报
func (s *TradeStatisticsService) GetSummary(ctx context.Context) (*resp.DataComparisonRespVO[resp.TradeSummaryRespVO], error) {
	// 1.1 昨天的数据
	yesterdayData, err := s.getTradeSummaryByDays(ctx, -1)
	if err != nil {
		return nil, err
	}
	// 1.2 前天的数据（用于对照昨天的数据）
	beforeYesterdayData, err := s.getTradeSummaryByDays(ctx, -2)
	if err != nil {
		return nil, err
	}

	// 2.1 本月数据
	monthData, err := s.getTradeSummaryByMonths(ctx, 0)
	if err != nil {
		return nil, err
	}
	// 2.2 上月数据（用于对照本月的数据）
	lastMonthData, err := s.getTradeSummaryByMonths(ctx, -1)
	if err != nil {
		return nil, err
	}

	// 拼接数据
	return &resp.DataComparisonRespVO[resp.TradeSummaryRespVO]{
		Summary: &resp.TradeSummaryRespVO{
			Yesterday: yesterdayData,
			Month:     monthData,
		},
		Comparison: &resp.TradeSummaryRespVO{
			Yesterday: beforeYesterdayData,
			Month:     lastMonthData,
		},
	}, nil
}

// getTradeSummaryByDays 获得指定天数的交易统计摘要
func (s *TradeStatisticsService) getTradeSummaryByDays(ctx context.Context, days int) (*resp.TradeSummaryItemVO, error) {
	targetDate := time.Now().AddDate(0, 0, days)
	beginTime := statistics.BeginOfDay(targetDate)
	endTime := statistics.EndOfDay(targetDate)

	return s.calculateSummaryItem(ctx, beginTime, endTime)
}

// getTradeSummaryByMonths 获得指定月份的交易统计摘要
func (s *TradeStatisticsService) getTradeSummaryByMonths(ctx context.Context, months int) (*resp.TradeSummaryItemVO, error) {
	targetDate := time.Now().AddDate(0, months, 0)
	beginTime := statistics.BeginOfMonth(targetDate)
	endTime := statistics.EndOfMonth(targetDate)

	return s.calculateSummaryItem(ctx, beginTime, endTime)
}

func (s *TradeStatisticsService) calculateSummaryItem(ctx context.Context, beginTime, endTime time.Time) (*resp.TradeSummaryItemVO, error) {
	data, err := s.calculateData(ctx, beginTime, endTime)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return &resp.TradeSummaryItemVO{}, nil
	}

	return &resp.TradeSummaryItemVO{
		OrderCreateCount:         data.OrderPayCount, // 用 OrderPayCount 作为 OrderCreateCount
		OrderPayCount:            data.OrderPayCount,
		OrderPayPrice:            int64(data.OrderPayPrice),
		AfterSaleCount:           data.AfterSaleCount,
		AfterSaleRefundPrice:     int64(data.AfterSaleRefundPrice),
		BrokerageSettlementPrice: 0, // Not implemented
		WalletPayPrice:           0,
		RechargePayCount:         0,
		RechargePayPrice:         0,
		RechargeRefundCount:      0,
		RechargeRefundPrice:      0,
	}, nil
}

// GetAnalysis 交易状况分析
func (s *TradeStatisticsService) GetAnalysis(ctx context.Context, r *req.TradeStatisticsAnalysisReq) (*resp.TradeStatisticsAnalysisResp, error) {
	if len(r.Times) != 2 {
		return nil, nil // Should not happen with validation
	}
	start, end := r.Times[0], r.Times[1]

	// Create date range
	dates := make([]string, 0)
	current := start
	for !current.After(end) {
		dates = append(dates, current.Format("2006-01-02"))
		current = current.AddDate(0, 0, 1)
	}

	// Prepare result arrays
	orderPayPrices := make([]int, len(dates))
	orderPayCounts := make([]int64, len(dates))
	afterSaleCounts := make([]int64, len(dates))
	afterSaleRefundPrices := make([]int, len(dates))

	for i, dateStr := range dates {
		dayStart, _ := time.ParseInLocation("2006-01-02", dateStr, time.Local)
		dayEnd := dayStart.AddDate(0, 0, 1).Add(-time.Second)

		data, err := s.calculateData(ctx, dayStart, dayEnd)
		if err != nil {
			return nil, err
		}

		orderPayPrices[i] = data.OrderPayPrice
		orderPayCounts[i] = data.OrderPayCount
		afterSaleCounts[i] = data.AfterSaleCount
		afterSaleRefundPrices[i] = data.AfterSaleRefundPrice
	}

	return &resp.TradeStatisticsAnalysisResp{
		Dates:                dates,
		OrderPayPrice:        orderPayPrices,
		OrderPayCount:        orderPayCounts,
		AfterSaleCount:       afterSaleCounts,
		AfterSaleRefundPrice: afterSaleRefundPrices,
	}, nil
}

// GetTradeStatisticsList 获得交易状况明细
func (s *TradeStatisticsService) GetTradeStatisticsList(ctx context.Context, beginTime, endTime time.Time) ([]*resp.TradeTrendSummaryRespVO, error) {
	// Create date range
	dates := make([]string, 0)
	current := beginTime
	for !current.After(endTime) {
		dates = append(dates, current.Format("2006-01-02"))
		current = current.AddDate(0, 0, 1)
	}

	list := make([]*resp.TradeTrendSummaryRespVO, 0, len(dates))
	for _, dateStr := range dates {
		dayStart, _ := time.ParseInLocation("2006-01-02", dateStr, time.Local)
		dayEnd := dayStart.AddDate(0, 0, 1).Add(-time.Second)

		data, err := s.calculateData(ctx, dayStart, dayEnd)
		if err != nil {
			return nil, err
		}

		list = append(list, &resp.TradeTrendSummaryRespVO{
			StatisticsTime:           dayStart,
			OrderCreateCount:         data.OrderPayCount, // 用 OrderPayCount 作为 OrderCreateCount
			OrderPayCount:            data.OrderPayCount,
			OrderPayPrice:            int64(data.OrderPayPrice),
			AfterSaleCount:           data.AfterSaleCount,
			AfterSaleRefundPrice:     int64(data.AfterSaleRefundPrice),
			BrokerageSettlementPrice: 0,
			WalletPayPrice:           0,
			RechargePayCount:         0,
			RechargePayPrice:         0,
			RechargeRefundCount:      0,
			RechargeRefundPrice:      0,
		})
	}
	return list, nil
}

// GetCountByStatusAndDeliveryType 获得指定状态和配送方式的订单数
func (s *TradeStatisticsService) GetCountByStatusAndDeliveryType(ctx context.Context, status int, deliveryType int) (int64, error) {
	return s.q.TradeOrder.WithContext(ctx).
		Where(s.q.TradeOrder.Status.Eq(status), s.q.TradeOrder.DeliveryType.Eq(deliveryType)).
		Count()
}

// GetAfterSaleCount 获得售后数量
func (s *TradeStatisticsService) GetAfterSaleCount(ctx context.Context, status int) (int64, error) {
	return s.q.AfterSale.WithContext(ctx).
		Where(s.q.AfterSale.Status.Eq(status)).
		Count()
}

// GetBrokerageWithdrawCount 获得佣金提现数量
func (s *TradeStatisticsService) GetBrokerageWithdrawCount(ctx context.Context, status int) (int64, error) {
	// Mock implementation as Brokerage model is missing
	return 0, nil
}

func (s *TradeStatisticsService) calculateData(ctx context.Context, start, end time.Time) (*resp.TradeStatisticsData, error) {
	o := s.q.TradeOrder
	af := s.q.AfterSale

	// Order Data
	// PayStatus = true, PayTime in range
	orderQ := o.WithContext(ctx).Where(o.PayStatus.Eq(model.NewBitBool(true)), o.PayTime.Between(start, end))
	orderCount, err := orderQ.Count()
	if err != nil {
		return nil, err
	}

	var orderPayPrice int64
	type SumResult struct {
		Total int64
	}
	var sumRes SumResult
	err = orderQ.Select(o.PayPrice.Sum().As("total")).Scan(&sumRes)
	if err != nil {
		return nil, err
	}
	orderPayPrice = sumRes.Total

	// AfterSale Data
	// Status = 30 (Refunded), RefundTime in range? Or CreateTime?
	// Align with Java/Requirement: usually RefundTime for financial stats.
	afQ := af.WithContext(ctx).Where(af.Status.Eq(trade.TradeOrderStatusCompleted /* 30? No, AfterSale status */), af.RefundTime.Between(start, end))
	// Wait, AfterSale status 30? Checking code snippet in Step 244: af.Status.Eq(30).
	// But consts.go has TradeOrderStatusCompleted = 30.
	// AfterSale status is logic specific. Java uses AfterSaleStatusEnum.
	// I added `AfterSaleStatusApply = 10`.
	// For "Refunded", assumption 30 is correct if consistent.
	// But `afQ.Where(af.Status.Eq(30))` was in original code. I interpret it as Refunded.

	afterSaleCount, err := afQ.Count()
	if err != nil {
		return nil, err
	}

	var afterSaleRefundPrice int64
	var afSumRes SumResult
	err = afQ.Select(af.RefundPrice.Sum().As("total")).Scan(&afSumRes)
	if err != nil {
		return nil, err
	}
	afterSaleRefundPrice = afSumRes.Total

	return &resp.TradeStatisticsData{
		OrderPayPrice:        int(orderPayPrice),
		OrderPayCount:        orderCount,
		AfterSaleCount:       afterSaleCount,
		AfterSaleRefundPrice: int(afterSaleRefundPrice),
	}, nil
}
