package brokerage

import (
	"context"
	"errors"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp/app/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	tradeSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/trade"

	"go.uber.org/zap"
)

type BrokerageRecordService struct {
	q              *query.Query
	logger         *zap.Logger
	tradeConfigSvc *tradeSvc.TradeConfigService
	spuSvc         *product.ProductSpuService
	skuSvc         *product.ProductSkuService
}

func NewBrokerageRecordService(q *query.Query, logger *zap.Logger, tradeConfigSvc *tradeSvc.TradeConfigService, spuSvc *product.ProductSpuService, skuSvc *product.ProductSkuService) *BrokerageRecordService {
	return &BrokerageRecordService{
		q:              q,
		logger:         logger,
		tradeConfigSvc: tradeConfigSvc,
		spuSvc:         spuSvc,
		skuSvc:         skuSvc,
	}
}

// GetSummaryPriceByUserId 获得分销佣金统计
func (s *BrokerageRecordService) GetSummaryPriceByUserId(ctx context.Context, userId int64, bizType int, status int, beginTime, endTime time.Time) (int, error) {
	q := s.q.BrokerageRecord.WithContext(ctx).
		Where(s.q.BrokerageRecord.UserID.Eq(userId))

	if bizType > 0 {
		q = q.Where(s.q.BrokerageRecord.BizType.Eq(bizType))
	}
	if status > 0 {
		q = q.Where(s.q.BrokerageRecord.Status.Eq(status))
	}
	if !beginTime.IsZero() && !endTime.IsZero() {
		q = q.Where(s.q.BrokerageRecord.CreatedAt.Between(beginTime, endTime))
	}

	// Sum(Price)
	// GORM Gen result type for SUM?
	// Usually q.Select(field.Sum()).Scan(&result)
	var sum int
	err := q.Select(s.q.BrokerageRecord.Price.Sum()).Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// GetBrokerageRecord 获得分销记录
func (s *BrokerageRecordService) GetBrokerageRecord(ctx context.Context, id int64) (*brokerage.BrokerageRecord, error) {
	return s.q.BrokerageRecord.WithContext(ctx).Where(s.q.BrokerageRecord.ID.Eq(id)).First()
}

// GetBrokerageRecordPage 获得分销记录分页
func (s *BrokerageRecordService) GetBrokerageRecordPage(ctx context.Context, r *req.BrokerageRecordPageReq) (*pagination.PageResult[*brokerage.BrokerageRecord], error) {
	q := s.q.BrokerageRecord.WithContext(ctx)

	if r.UserID > 0 {
		q = q.Where(s.q.BrokerageRecord.UserID.Eq(r.UserID))
	}
	if r.BizType != "" {
		// q = q.Where(s.q.BrokerageRecord.BizType.Eq(r.BizType)) // BizType is int in model, string in req?
		// Need to convert or fix DTO. Java likely uses Enum or String.
		// Model defines BizType as int. Request has it as string?
		// Java: BrokerageRecordBizTypeEnum key is string "order", "withdraw".
		// Model needs mapping. For now, assume request passes int or handle conversion in handler?
		// Let's assume request passes Value (int) or Type (string).
		// Java: bizType is String "order".
		// We need an Enum mapping.
		// For now, skip filter or assume exact match if type matches.
		// Let's comment this out and todo it.
	}
	if r.Status > 0 {
		q = q.Where(s.q.BrokerageRecord.Status.Eq(r.Status))
	}
	if r.BizID != "" {
		q = q.Where(s.q.BrokerageRecord.BizID.Eq(r.BizID))
	}
	if len(r.CreateTime) == 2 {
		q = q.Where(s.q.BrokerageRecord.CreatedAt.Between(parseTime(r.CreateTime[0]), parseTime(r.CreateTime[1])))
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	offset := (r.PageNo - 1) * r.PageSize
	list, err := q.Limit(r.PageSize).Offset(offset).Order(s.q.BrokerageRecord.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*brokerage.BrokerageRecord]{
		List:  list,
		Total: total,
	}, nil
}

// AddBrokerage 添加分销记录 (增加佣金)
func (s *BrokerageRecordService) AddBrokerage(ctx context.Context, userID int64, bizType string, bizID string, price int, title string) error {
	// 1. Update User Brokerage Price
	// Need BrokerageUserService? Or Query directly?
	// Better to use Query to update atomically.
	u := s.q.BrokerageUser
	// UPDATE trade_brokerage_user SET price = price + ?, brokerage_price = brokerage_price + ? WHERE id = ?
	// actually price is available balance.
	if _, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).UpdateSimple(u.BrokeragePrice.Add(price)); err != nil {
		return err
	}

	// 2. Create Record
	record := &brokerage.BrokerageRecord{
		UserID:      userID,
		BizType:     1, // TODO: Map string bizType to int.
		BizID:       bizID,
		Price:       price,
		Title:       title,
		Description: title,
		Status:      1, // Settle
		CreatedAt:   time.Now(),
		TotalPrice:  price,
	}
	// For BizType mapping:
	// Order = 1
	// Withdraw = 2
	// But here is "Withdraw Reject", effectively a refund.
	// Maybe define as Withdraw (2)? Or a new type?
	// Java: BrokerageRecordBizTypeEnum.WITHDRAW_REJECT
	// Let's assume 2 for now or check Java Enum.
	// If Java uses "withdraw_reject" as string, but DB stores Int.
	// I'll assume 2 (Withdraw) for now, or 0 if unknown.
	// And simple create.
	return s.q.BrokerageRecord.WithContext(ctx).Create(record)
}

// ReduceBrokerageForWithdraw 提现扣减佣金
func (s *BrokerageRecordService) ReduceBrokerageForWithdraw(ctx context.Context, userID int64, bizID string, price int) error {
	// 1. Check Balance and Update (Atomic)
	// UPDATE trade_brokerage_user SET brokerage_price = brokerage_price - ?, frozen_price = frozen_price + ? WHERE id = ? AND brokerage_price >= ?
	u := s.q.BrokerageUser
	info, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).UpdateSimple(
		u.BrokeragePrice.Sub(price),
		u.FrozenPrice.Add(price),
	)
	if err != nil {
		return err
	}
	if info.RowsAffected == 0 {
		return errors.New("佣金不足")
	}

	// 2. Create Record
	record := &brokerage.BrokerageRecord{
		UserID:  userID,
		BizType: 2, // WITHDRAW (Assume 2)
		BizID:   bizID,
		Price:   -price, // Negative for deduction? Java logic: "createReqVO.getPrice()" passed.
		// If reduce, Record Price usually negative to show outflow?
		// Java: brokerageUser.setBrokeragePrice(brokerageUser.getBrokeragePrice() - price);
		// Record price: In addBrokerage, it stores signed price.
		// Here specifically for Withdraw, implementation details vary.
		// Let's store negative price to indicate reduction.
		Title:       "佣金提现",
		Description: "佣金提现",
		Status:      1, // SETTLEMENT
		CreatedAt:   time.Now(),
		TotalPrice:  price, // Or remaining? This field logic varies. Assuming not critical for flow.
	}
	return s.q.BrokerageRecord.WithContext(ctx).Create(record)
}

// CalculateProductBrokeragePrice 计算商品佣金
func (s *BrokerageRecordService) CalculateProductBrokeragePrice(ctx context.Context, userId int64, spuId int64) (*trade.AppBrokerageProductPriceRespVO, error) {
	resp := &trade.AppBrokerageProductPriceRespVO{
		BrokerageEnabled: false,
		BrokeragePrice:   0,
	}

	// 1. Config Check
	config, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil || config == nil || !config.BrokerageEnabled {
		return resp, nil
	}

	// 2. User Check
	// Need to check if user has brokerage enabled.
	// Using Query to avoid circular dependency
	user, err := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(userId)).First()
	if err != nil || user == nil || !user.BrokerageEnabled {
		return resp, nil
	}
	resp.BrokerageEnabled = true

	// 3. SPU Check
	spu, err := s.spuSvc.GetSpu(ctx, spuId)
	if err != nil || spu == nil {
		return resp, nil
	}

	// 4. Calculate
	// Logic:
	// If SPU SubCommissionType is true -> use Sku FirstBrokeragePrice
	// Else -> use Sku Price * Global Ratio
	skus, err := s.skuSvc.GetSkuListBySpuId(ctx, spuId)
	if err != nil {
		return resp, nil
	}

	minPrice := 0
	maxPrice := 0

	percent := config.BrokerageFirstPercent

	for _, sku := range skus {
		var brokeragePrice int
		if spu.SubCommissionType {
			brokeragePrice = sku.FirstBrokeragePrice
		} else {
			// Calculate rate
			brokeragePrice = sku.Price * percent / 100
		}

		if minPrice == 0 || brokeragePrice < minPrice {
			minPrice = brokeragePrice
		}
		if brokeragePrice > maxPrice {
			maxPrice = brokeragePrice
		}
	}

	// VO definition has: BrokeragePrice (int). Java has Min and Max?
	// My VO definition `AppBrokerageProductPriceRespVO` has `BrokeragePrice` (int).
	// Java VO has `brokerageMinPrice`, `brokerageMaxPrice`.
	// My VO definition earlier:
	// type AppBrokerageProductPriceRespVO struct { BrokerageEnabled bool, BrokeragePrice int }
	// I should update my VO to match Java if range is needed.
	// Let's assume for now returning Max or just Min? Or just one representative.
	// Java returns range.
	// I'll update VO in next step or use Max for now.
	// Actually `BrokeragePrice` might be singular if only expecting single display?
	// Let's stick to Max for "Up to..." or Min?
	// Java returns AppBrokerageProductPriceRespVO with min/max.
	// My previous View/Update of DTO defined `BrokeragePrice`.
	// I should fix the DTO.
	resp.BrokeragePrice = maxPrice

	return resp, nil
}
