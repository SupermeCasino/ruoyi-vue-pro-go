package calculators

import (
	"context"

	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"go.uber.org/zap"
)

// RewardActivityPriceCalculator 满减送活动价格计算器
type RewardActivityPriceCalculator struct {
	*trade.BasePriceCalculator
	rewardActivitySvc *promotion.RewardActivityService
}

// NewRewardActivityPriceCalculator 创建满减送活动价格计算器
func NewRewardActivityPriceCalculator(
	rewardActivitySvc *promotion.RewardActivityService,
	helper *trade.PriceCalculatorHelper,
	logger *zap.Logger,
) *RewardActivityPriceCalculator {
	return &RewardActivityPriceCalculator{
		BasePriceCalculator: trade.NewBasePriceCalculator(
			tradeModel.CalculatorNameReward,
			tradeModel.OrderRewardActivity,
			helper,
			logger,
		),
		rewardActivitySvc: rewardActivitySvc,
	}
}

// Calculate 执行满减送活动价格计算
func (c *RewardActivityPriceCalculator) Calculate(ctx context.Context, req *trade.TradePriceCalculateReqBO, resp *trade.TradePriceCalculateRespBO) error {
	// 只处理普通订单
	if resp.Type != tradeModel.TradeOrderTypeNormal {
		return nil
	}

	c.LogCalculation(ctx, req, "开始执行满减送活动价格计算")

	// 构建活动匹配项
	matchItems := make([]promotion.ActivityMatchItem, 0)
	for _, item := range resp.Items {
		if !item.Selected { // 过滤未选中项，对齐Java版本逻辑
			continue
		}
		matchItems = append(matchItems, promotion.ActivityMatchItem{
			SkuID:      item.SkuID,
			SpuID:      item.SpuID,
			CategoryID: item.CategoryID,
			Price:      item.Price,
			Count:      item.Count,
		})
	}

	if len(matchItems) == 0 {
		return nil
	}

	// 计算满减送活动
	activityDiscount, rewardResults, err := c.rewardActivitySvc.CalculateRewardActivity(ctx, matchItems)
	if err != nil {
		c.LogError(ctx, req, err, "计算满减送活动失败")
		return err
	}

	if activityDiscount <= 0 {
		return nil
	}

	c.LogCalculation(ctx, req, "满减送活动计算结果",
		zap.Int("activityDiscount", activityDiscount),
		zap.Int("resultCount", len(rewardResults)),
	)

	// 分摊满减折扣到各项
	divideActivityDiscounts := c.Helper.DividePrice(resp.Items, activityDiscount)
	for i := range resp.Items {
		if !resp.Items[i].Selected {
			continue
		}

		resp.Items[i].DiscountPrice += divideActivityDiscounts[i]

		// 重新计算支付金额
		c.Helper.RecountPayPrice(&resp.Items[i])

		c.LogCalculation(ctx, req, "分摊满减送折扣",
			zap.Int64("skuId", resp.Items[i].SkuID),
			zap.Int("dividedDiscount", divideActivityDiscounts[i]),
			zap.Int("totalDiscountPrice", resp.Items[i].DiscountPrice),
		)
	}

	// 添加促销活动明细到响应
	for _, res := range rewardResults {
		p := &trade.TradePriceCalculatePromotionBO{
			ID:            res.ActivityID,
			Name:          res.ActivityName,
			Type:          tradeModel.OrderRewardActivity,
			TotalPrice:    res.TotalPrice,
			DiscountPrice: res.TotalDiscount,
			Match:         true,
		}

		// 添加商品项明细
		for _, skuID := range res.SkuIDs {
			p.Items = append(p.Items, trade.TradePriceCalculatePromotionItemBO{
				SkuID: skuID,
			})
		}

		c.Helper.AddPromotion(resp, p)
	}

	return nil
}

// IsApplicable 判断是否适用于当前订单类型
func (c *RewardActivityPriceCalculator) IsApplicable(orderType int) bool {
	return orderType == tradeModel.TradeOrderTypeNormal
}
