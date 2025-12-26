package calculators

import (
	"context"

	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade"

	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"go.uber.org/zap"
)

// DiscountActivityPriceCalculator 限时折扣活动价格计算器
type DiscountActivityPriceCalculator struct {
	*trade.BasePriceCalculator
	discountActivitySvc promotion.DiscountActivityService
	memberUserSvc       *memberSvc.MemberUserService
	memberLevelSvc      *memberSvc.MemberLevelService
}

// NewDiscountActivityPriceCalculator 创建限时折扣活动价格计算器
func NewDiscountActivityPriceCalculator(
	discountActivitySvc promotion.DiscountActivityService,
	memberUserSvc *memberSvc.MemberUserService,
	memberLevelSvc *memberSvc.MemberLevelService,
	Helper *trade.PriceCalculatorHelper,
	logger *zap.Logger,
) *DiscountActivityPriceCalculator {
	return &DiscountActivityPriceCalculator{
		BasePriceCalculator: trade.NewBasePriceCalculator(
			tradeModel.CalculatorNameDiscount,
			tradeModel.OrderDiscountActivity,
			Helper,
			logger,
		),
		discountActivitySvc: discountActivitySvc,
		memberUserSvc:       memberUserSvc,
		memberLevelSvc:      memberLevelSvc,
	}
}

// Calculate 执行限时折扣活动价格计算
func (c *DiscountActivityPriceCalculator) Calculate(ctx context.Context, req *trade.TradePriceCalculateReqBO, resp *trade.TradePriceCalculateRespBO) error {
	// 只处理普通订单
	if resp.Type != tradeModel.TradeOrderTypeNormal {
		return nil
	}

	c.LogCalculation(ctx, req, "开始执行限时折扣活动价格计算")

	// 获取所有SKU ID
	var skuIDs []int64
	for _, item := range resp.Items {
		if item.Selected {
			skuIDs = append(skuIDs, item.SkuID)
		}
	}

	if len(skuIDs) == 0 {
		return nil
	}

	// 获取限时折扣活动
	discountActivityMap, err := c.discountActivitySvc.GetMatchDiscountProductMap(ctx, skuIDs)
	if err != nil {
		c.LogError(ctx, req, err, "获取限时折扣活动失败")
		return err
	}

	// 获取会员等级信息
	var memberLevel *memberModel.MemberLevel
	if req.UserID > 0 {
		user, _ := c.memberUserSvc.GetUser(ctx, req.UserID)
		if user != nil && user.LevelID > 0 {
			memberLevel, _ = c.memberLevelSvc.GetLevel(ctx, user.LevelID)
		}
	}

	// 计算每个商品项的折扣
	for i := range resp.Items {
		if !resp.Items[i].Selected {
			continue
		}

		item := &resp.Items[i]
		itemPayPrice := item.Price * item.Count

		// A. 计算限时折扣
		activitySaving := 0
		if discount, ok := discountActivityMap[item.SkuID]; ok {
			if discount.DiscountType == tradeModel.DiscountTypePrice { // 减价
				activitySaving = discount.DiscountPrice * item.Count
			} else if discount.DiscountType == tradeModel.DiscountTypePercent { // 打折
				targetPrice := c.Helper.CalculateRatePrice(itemPayPrice, discount.DiscountPercent)
				activitySaving = itemPayPrice - targetPrice
			}
		}

		// B. 计算 VIP 折扣
		vipSaving := 0
		if memberLevel != nil && memberLevel.Status == 0 && memberLevel.DiscountPercent > 0 {
			targetPrice := c.Helper.CalculateRatePrice(itemPayPrice, memberLevel.DiscountPercent)
			vipSaving = itemPayPrice - targetPrice
		}

		// C. 取其优（对齐 Java TradeDiscountActivityPriceCalculator）
		if activitySaving > 0 || vipSaving > 0 {
			if activitySaving >= vipSaving {
				item.DiscountPrice += activitySaving
				item.VipPrice = 0

				// 记录限时折扣活动明细
				if discount, ok := discountActivityMap[item.SkuID]; ok {
					c.Helper.AddPromotion(resp, &trade.TradePriceCalculatePromotionBO{
						ID:            discount.ActivityID,
						Name:          "限时折扣",
						Type:          tradeModel.PromotionTypeDiscountActivity,
						TotalPrice:    itemPayPrice,
						DiscountPrice: activitySaving,
						Match:         true,
						Items: []trade.TradePriceCalculatePromotionItemBO{
							{
								SkuID:         item.SkuID,
								TotalPrice:    itemPayPrice,
								DiscountPrice: activitySaving,
								PayPrice:      itemPayPrice - activitySaving,
							},
						},
					})
				}

				c.LogCalculation(ctx, req, "应用限时折扣活动",
					zap.Int64("skuId", item.SkuID),
					zap.Int("activitySaving", activitySaving),
					zap.Int("vipSaving", vipSaving),
					zap.String("choice", "activity"),
				)
			} else {
				item.VipPrice += vipSaving

				// 记录会员等级折扣明细
				c.Helper.AddPromotion(resp, &trade.TradePriceCalculatePromotionBO{
					ID:            memberLevel.ID,
					Name:          memberLevel.Name,
					Type:          tradeModel.PromotionTypeMemberLevel,
					TotalPrice:    itemPayPrice,
					DiscountPrice: vipSaving,
					Match:         true,
					Items: []trade.TradePriceCalculatePromotionItemBO{
						{
							SkuID:         item.SkuID,
							TotalPrice:    itemPayPrice,
							DiscountPrice: vipSaving,
							PayPrice:      itemPayPrice - vipSaving,
						},
					},
				})

				c.LogCalculation(ctx, req, "应用VIP折扣",
					zap.Int64("skuId", item.SkuID),
					zap.Int("activitySaving", activitySaving),
					zap.Int("vipSaving", vipSaving),
					zap.String("choice", "vip"),
				)
			}

			// 重新计算支付金额
			c.Helper.RecountPayPrice(item)
		}
	}

	return nil
}

// IsApplicable 判断是否适用于当前订单类型
func (c *DiscountActivityPriceCalculator) IsApplicable(orderType int) bool {
	return orderType == tradeModel.TradeOrderTypeNormal
}
