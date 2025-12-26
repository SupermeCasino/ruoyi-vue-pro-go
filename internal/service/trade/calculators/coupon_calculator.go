package calculators

import (
	"context"

	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"go.uber.org/zap"
)

// CouponPriceCalculator 优惠券价格计算器
type CouponPriceCalculator struct {
	*trade.BasePriceCalculator
	couponSvc *promotion.CouponUserService
}

// NewCouponPriceCalculator 创建优惠券价格计算器
func NewCouponPriceCalculator(
	couponSvc *promotion.CouponUserService,
	helper *trade.PriceCalculatorHelper,
	logger *zap.Logger,
) *CouponPriceCalculator {
	return &CouponPriceCalculator{
		BasePriceCalculator: trade.NewBasePriceCalculator(
			tradeModel.CalculatorNameCoupon,
			tradeModel.OrderCoupon,
			helper,
			logger,
		),
		couponSvc: couponSvc,
	}
}

// Calculate 执行优惠券价格计算
func (c *CouponPriceCalculator) Calculate(ctx context.Context, req *trade.TradePriceCalculateReqBO, resp *trade.TradePriceCalculateRespBO) error {
	// 只处理普通订单且指定了优惠券
	if resp.Type != tradeModel.TradeOrderTypeNormal || req.CouponID == nil || *req.CouponID <= 0 {
		return nil
	}

	c.LogCalculation(ctx, req, "开始执行优惠券价格计算",
		zap.Int64("couponId", *req.CouponID),
	)

	// 收集SPU ID和分类ID
	var spuIDs []int64
	var categoryIDs []int64
	spuMapForCoupon := make(map[int64]bool)
	catMapForCoupon := make(map[int64]bool)

	for _, item := range resp.Items {
		if !item.Selected {
			continue
		}

		if !spuMapForCoupon[item.SpuID] {
			spuIDs = append(spuIDs, item.SpuID)
			spuMapForCoupon[item.SpuID] = true
		}
		if !catMapForCoupon[item.CategoryID] {
			categoryIDs = append(categoryIDs, item.CategoryID)
			catMapForCoupon[item.CategoryID] = true
		}
	}

	if len(spuIDs) == 0 {
		return nil
	}

	// 计算当前支付金额
	currentPayPrice := c.Helper.CalculateTotalPayPrice(resp.Items, true)

	// 计算优惠券折扣
	couponPrice, err := c.couponSvc.CalculateCoupon(ctx, req.UserID, *req.CouponID, int64(currentPayPrice), spuIDs, categoryIDs)
	if err != nil {
		c.LogError(ctx, req, err, "计算优惠券折扣失败")
		return err
	}

	couponPriceInt := int(couponPrice)
	if couponPriceInt <= 0 {
		return nil
	}

	c.LogCalculation(ctx, req, "优惠券计算结果",
		zap.Int("currentPayPrice", currentPayPrice),
		zap.Int("couponPrice", couponPriceInt),
	)

	// 分摊优惠券折扣到各项
	divideCouponPrices := c.Helper.DividePrice(resp.Items, couponPriceInt)
	for i := range resp.Items {
		if !resp.Items[i].Selected {
			continue
		}

		resp.Items[i].CouponPrice += divideCouponPrices[i]

		// 重新计算支付金额
		c.Helper.RecountPayPrice(&resp.Items[i])

		c.LogCalculation(ctx, req, "分摊优惠券折扣",
			zap.Int64("skuId", resp.Items[i].SkuID),
			zap.Int("dividedCouponPrice", divideCouponPrices[i]),
			zap.Int("totalCouponPrice", resp.Items[i].CouponPrice),
		)
	}

	// 更新响应中的优惠券信息
	resp.CouponID = *req.CouponID

	return nil
}

// IsApplicable 判断是否适用于当前订单类型
func (c *CouponPriceCalculator) IsApplicable(orderType int) bool {
	return orderType == tradeModel.TradeOrderTypeNormal
}
