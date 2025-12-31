package repo

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/dto"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

// ============ TradeStatisticsRepository 实现 ============

// TradeStatisticsRepositoryImpl 交易统计 Repository 实现（基于 gorm gen）
type TradeStatisticsRepositoryImpl struct {
	q *query.Query
}

// NewTradeStatisticsRepository 创建交易统计 Repository
func NewTradeStatisticsRepository(q *query.Query) *TradeStatisticsRepositoryImpl {
	return &TradeStatisticsRepositoryImpl{q: q}
}

// GetByDateRange 获取指定日期范围的统计汇总
func (r *TradeStatisticsRepositoryImpl) GetByDateRange(ctx context.Context, beginTime, endTime time.Time) (*dto.TradeStatisticsDTO, error) {
	ts := r.q.TradeStatistics

	var result struct {
		OrderCreateCount         int64 `gorm:"column:order_create_count"`
		OrderPayCount            int64 `gorm:"column:order_pay_count"`
		OrderPayPrice            int64 `gorm:"column:order_pay_price"`
		AfterSaleCount           int64 `gorm:"column:after_sale_count"`
		AfterSaleRefundPrice     int64 `gorm:"column:after_sale_refund_price"`
		BrokerageSettlementPrice int64 `gorm:"column:brokerage_settlement_price"`
		WalletPayPrice           int64 `gorm:"column:wallet_pay_price"`
		RechargePayCount         int64 `gorm:"column:recharge_pay_count"`
		RechargePayPrice         int64 `gorm:"column:recharge_pay_price"`
		RechargeRefundCount      int64 `gorm:"column:recharge_refund_count"`
		RechargeRefundPrice      int64 `gorm:"column:recharge_refund_price"`
	}

	err := ts.WithContext(ctx).
		Select(
			ts.OrderCreateCount.Sum().As("order_create_count"),
			ts.OrderPayCount.Sum().As("order_pay_count"),
			ts.OrderPayPrice.Sum().As("order_pay_price"),
			ts.AfterSaleCount.Sum().As("after_sale_count"),
			ts.AfterSaleRefundPrice.Sum().As("after_sale_refund_price"),
			ts.BrokerageSettlementPrice.Sum().As("brokerage_settlement_price"),
			ts.WalletPayPrice.Sum().As("wallet_pay_price"),
			ts.RechargePayCount.Sum().As("recharge_pay_count"),
			ts.RechargePayPrice.Sum().As("recharge_pay_price"),
			ts.RechargeRefundCount.Sum().As("recharge_refund_count"),
			ts.RechargeRefundPrice.Sum().As("recharge_refund_price"),
		).
		Where(ts.Time.Between(beginTime, endTime)).
		Where(ts.Time.Between(beginTime, endTime)).
		// Where(ts.Deleted.Is(false)). // Removed explicit check
		Scan(&result)
	if err != nil {
		return nil, err
	}

	return &dto.TradeStatisticsDTO{
		StatisticsTime:           beginTime,
		OrderCreateCount:         int(result.OrderCreateCount),
		OrderPayCount:            int(result.OrderPayCount),
		OrderPayPrice:            int(result.OrderPayPrice),
		AfterSaleCount:           int(result.AfterSaleCount),
		AfterSaleRefundPrice:     int(result.AfterSaleRefundPrice),
		BrokerageSettlementPrice: int(result.BrokerageSettlementPrice),
		WalletPayPrice:           int(result.WalletPayPrice),
		RechargePayCount:         int(result.RechargePayCount),
		RechargePayPrice:         int(result.RechargePayPrice),
		RechargeRefundCount:      int(result.RechargeRefundCount),
		RechargeRefundPrice:      int(result.RechargeRefundPrice),
	}, nil
}

// GetByMonthRange 获取指定月份范围的统计汇总
func (r *TradeStatisticsRepositoryImpl) GetByMonthRange(ctx context.Context, beginTime, endTime time.Time) (*dto.TradeStatisticsDTO, error) {
	return r.GetByDateRange(ctx, beginTime, endTime)
}

// GetListByDateRange 获取指定日期范围的统计列表
func (r *TradeStatisticsRepositoryImpl) GetListByDateRange(ctx context.Context, beginTime, endTime time.Time) ([]*dto.TradeStatisticsDTO, error) {
	ts := r.q.TradeStatistics

	list, err := ts.WithContext(ctx).
		Where(ts.Time.Between(beginTime, endTime)).
		Where(ts.Time.Between(beginTime, endTime)).
		// Where(ts.Deleted.Is(false)). // Removed explicit check
		Order(ts.Time).
		Find()
	if err != nil {
		return nil, err
	}

	result := make([]*dto.TradeStatisticsDTO, 0, len(list))
	for _, item := range list {
		result = append(result, &dto.TradeStatisticsDTO{
			StatisticsTime:           item.Time,
			OrderCreateCount:         item.OrderCreateCount,
			OrderPayCount:            item.OrderPayCount,
			OrderPayPrice:            item.OrderPayPrice,
			AfterSaleCount:           item.AfterSaleCount,
			AfterSaleRefundPrice:     item.AfterSaleRefundPrice,
			BrokerageSettlementPrice: item.BrokerageSettlementPrice,
			WalletPayPrice:           item.WalletPayPrice,
			RechargePayCount:         item.RechargePayCount,
			RechargePayPrice:         item.RechargePayPrice,
			RechargeRefundCount:      item.RechargeRefundCount,
			RechargeRefundPrice:      item.RechargeRefundPrice,
		})
	}
	return result, nil
}

// Insert 插入统计记录
func (r *TradeStatisticsRepositoryImpl) Insert(ctx context.Context, stats *dto.TradeStatisticsDTO) error {
	ts := r.q.TradeStatistics

	record := &trade.TradeStatistics{
		Time:                     stats.StatisticsTime,
		OrderCreateCount:         stats.OrderCreateCount,
		OrderPayCount:            stats.OrderPayCount,
		OrderPayPrice:            stats.OrderPayPrice,
		AfterSaleCount:           stats.AfterSaleCount,
		AfterSaleRefundPrice:     stats.AfterSaleRefundPrice,
		BrokerageSettlementPrice: stats.BrokerageSettlementPrice,
		WalletPayPrice:           stats.WalletPayPrice,
		RechargePayCount:         stats.RechargePayCount,
		RechargePayPrice:         stats.RechargePayPrice,
		RechargeRefundCount:      stats.RechargeRefundCount,
		RechargeRefundPrice:      stats.RechargeRefundPrice,
	}

	return ts.WithContext(ctx).Create(record)
}

// ============ TradeOrderStatisticsRepository 实现 ============

// TradeOrderStatisticsRepositoryImpl 交易订单统计 Repository 实现
type TradeOrderStatisticsRepositoryImpl struct {
	q *query.Query
}

// NewTradeOrderStatisticsRepository 创建交易订单统计 Repository
func NewTradeOrderStatisticsRepository(q *query.Query) *TradeOrderStatisticsRepositoryImpl {
	return &TradeOrderStatisticsRepositoryImpl{q: q}
}

// GetCountByCreateTime 获取指定时间范围内的订单数量
func (r *TradeOrderStatisticsRepositoryImpl) GetCountByCreateTime(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	t := r.q.TradeOrder
	return t.WithContext(ctx).
		// Where(t.Deleted.Is(false)).
		Where(t.CreateTime.Between(beginTime, endTime)).
		Count()
}

// GetPayPriceSummary 获取指定时间范围内的支付金额汇总
func (r *TradeOrderStatisticsRepositoryImpl) GetPayPriceSummary(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	t := r.q.TradeOrder

	var result struct {
		TotalPayPrice int64 `gorm:"column:total_pay_price"`
	}
	err := t.WithContext(ctx).
		Select(t.PayPrice.Sum().As("total_pay_price")).
		// Where(t.Deleted.Is(false)).
		Where(t.PayTime.Between(beginTime, endTime)).
		Scan(&result)
	if err != nil {
		return 0, err
	}
	return result.TotalPayPrice, nil
}

// GetUserCountByCreateTime 获取指定时间范围内的下单用户数
func (r *TradeOrderStatisticsRepositoryImpl) GetUserCountByCreateTime(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	t := r.q.TradeOrder
	return t.WithContext(ctx).
		// Where(t.Deleted.Is(false)).
		Where(t.CreateTime.Between(beginTime, endTime)).
		Distinct(t.UserID).
		Count()
}

// GetCountByStatusAndDeliveryType 获取指定状态和配送类型的订单数量
func (r *TradeOrderStatisticsRepositoryImpl) GetCountByStatusAndDeliveryType(ctx context.Context, status, deliveryType int) (int64, error) {
	t := r.q.TradeOrder
	q := t.WithContext(ctx) // .Where(t.Deleted.Is(false))
	if status >= 0 {
		q = q.Where(t.Status.Eq(status))
	}
	if deliveryType >= 0 {
		q = q.Where(t.DeliveryType.Eq(deliveryType))
	}
	return q.Count()
}

// GetPayUserCount 获取指定时间范围内的支付用户数
func (r *TradeOrderStatisticsRepositoryImpl) GetPayUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	t := r.q.TradeOrder
	return t.WithContext(ctx).
		// Where(t.Deleted.Is(false)).
		Where(t.PayTime.Between(beginTime, endTime)).
		Distinct(t.UserID).
		Count()
}

// GetOrderPayPrice 获取指定时间范围内的订单支付金额
func (r *TradeOrderStatisticsRepositoryImpl) GetOrderPayPrice(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	return r.GetPayPriceSummary(ctx, beginTime, endTime)
}

// GetOrderUserCount 获取指定时间范围内的下单用户数
func (r *TradeOrderStatisticsRepositoryImpl) GetOrderUserCount(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	return r.GetUserCountByCreateTime(ctx, beginTime, endTime)
}

// ============ AfterSaleStatisticsRepository 实现 ============

// AfterSaleStatisticsRepositoryImpl 售后统计 Repository 实现
type AfterSaleStatisticsRepositoryImpl struct {
	q *query.Query
}

// NewAfterSaleStatisticsRepository 创建售后统计 Repository
func NewAfterSaleStatisticsRepository(q *query.Query) *AfterSaleStatisticsRepositoryImpl {
	return &AfterSaleStatisticsRepositoryImpl{q: q}
}

// GetRefundPriceSummary 获取指定时间范围内的退款金额汇总
func (r *AfterSaleStatisticsRepositoryImpl) GetRefundPriceSummary(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	a := r.q.AfterSale

	var result struct {
		TotalRefundPrice int64 `gorm:"column:total_refund_price"`
	}
	err := a.WithContext(ctx).
		Select(a.RefundPrice.Sum().As("total_refund_price")).
		// Where(a.Deleted.Is(false)).
		Where(a.RefundTime.Between(beginTime, endTime)).
		Scan(&result)
	if err != nil {
		return 0, err
	}
	return result.TotalRefundPrice, nil
}

// GetCountByStatus 获取指定状态的售后数量
func (r *AfterSaleStatisticsRepositoryImpl) GetCountByStatus(ctx context.Context, status int) (int64, error) {
	a := r.q.AfterSale
	return a.WithContext(ctx).
		// Where(a.Deleted.Is(false)).
		Where(a.Status.Eq(status)).
		Count()
}

// ============ BrokerageStatisticsRepository 实现 ============

// BrokerageStatisticsRepositoryImpl 佣金统计 Repository 实现
type BrokerageStatisticsRepositoryImpl struct {
	q *query.Query
}

// NewBrokerageStatisticsRepository 创建佣金统计 Repository
func NewBrokerageStatisticsRepository(q *query.Query) *BrokerageStatisticsRepositoryImpl {
	return &BrokerageStatisticsRepositoryImpl{q: q}
}

// GetSettlementPriceSummary 获取指定时间范围内已结算的佣金汇总
func (r *BrokerageStatisticsRepositoryImpl) GetSettlementPriceSummary(ctx context.Context, status int, beginTime, endTime time.Time) (int64, error) {
	// brokerage_record 表需要使用 TradeStatistics 表中的预计算数据
	// 或直接从订单表中计算佣金
	ts := r.q.TradeStatistics

	var result struct {
		TotalSettlementPrice int64 `gorm:"column:total"`
	}
	err := ts.WithContext(ctx).
		Select(ts.BrokerageSettlementPrice.Sum().As("total")).
		Where(ts.Time.Between(beginTime, endTime)).
		// Where(ts.Deleted.Is(false)).
		Scan(&result)
	if err != nil {
		return 0, err
	}
	return result.TotalSettlementPrice, nil
}

// GetWithdrawCountByStatus 获取指定状态的提现申请数量
func (r *BrokerageStatisticsRepositoryImpl) GetWithdrawCountByStatus(ctx context.Context, status int) (int64, error) {
	// brokerage_withdraw 表暂无 gorm gen，返回 0
	// 需要后续生成 brokerage_withdraw 表的 gorm gen
	return 0, nil
}
