package calculators

import (
	"context"

	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/consts"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"go.uber.org/zap"
)

// PointGivePriceCalculator 积分赠送计算器
type PointGivePriceCalculator struct {
	*trade.BasePriceCalculator
	memberConfigSvc *memberSvc.MemberConfigService
	helper          *trade.PriceCalculatorHelper
}

// NewPointGivePriceCalculator 创建积分赠送计算器
func NewPointGivePriceCalculator(
	memberConfigSvc *memberSvc.MemberConfigService,
	helper *trade.PriceCalculatorHelper,
	logger *zap.Logger,
) *PointGivePriceCalculator {
	return &PointGivePriceCalculator{
		BasePriceCalculator: trade.NewBasePriceCalculator(
			tradeModel.CalculatorNamePointGive,
			tradeModel.OrderPointGive,
			helper,
			logger,
		),
		memberConfigSvc: memberConfigSvc,
		helper:          helper, // Initialize the helper field
	}
}

// Calculate 执行积分赠送计算
func (c *PointGivePriceCalculator) Calculate(ctx context.Context, req *trade.TradePriceCalculateReqBO, resp *trade.TradePriceCalculateRespBO) error {
	c.LogCalculation(ctx, req, "开始执行积分赠送计算")

	// 获取会员配置
	config, err := c.memberConfigSvc.GetConfig(ctx)
	if err != nil || config == nil || config.PointTradeGivePoint <= 0 {
		c.LogCalculation(ctx, req, "积分赠送功能未启用")
		return nil
	}

	// 计算总支付金额（用于积分赠送计算）
	totalPayPrice := c.helper.CalculateTotalPayPrice(resp.Items, true)
	if totalPayPrice <= 0 {
		return nil
	}

	// 计算赠送积分
	// 通常按支付金额的一定比例赠送积分
	givePoints := totalPayPrice / config.PointTradeGivePoint
	if givePoints <= 0 {
		return nil
	}

	c.LogCalculation(ctx, req, "积分赠送计算结果",
		zap.Int("totalPayPrice", totalPayPrice),
		zap.Int("givePointRate", config.PointTradeGivePoint),
		zap.Int("givePoints", givePoints),
	)

	// 分摊赠送积分到各项
	divideGivePoints := c.helper.DividePrice(resp.Items, givePoints)

	totalGivePoints := 0
	for i := range resp.Items {
		if !resp.Items[i].Selected {
			continue
		}

		resp.Items[i].GivePoint = divideGivePoints[i]
		totalGivePoints += divideGivePoints[i]

		c.LogCalculation(ctx, req, "分摊赠送积分",
			zap.Int64("skuId", resp.Items[i].SkuID),
			zap.Int("givePoint", resp.Items[i].GivePoint),
		)
	}

	// 更新响应中的赠送积分总数
	resp.GivePoint = totalGivePoints

	c.LogCalculation(ctx, req, "积分赠送计算完成",
		zap.Int("totalGivePoints", totalGivePoints),
	)

	return nil
}

// IsApplicable 判断是否适用于当前订单类型
func (c *PointGivePriceCalculator) IsApplicable(orderType int) bool {
	// 积分赠送适用于所有订单类型
	return true
}
