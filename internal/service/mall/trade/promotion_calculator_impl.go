package trade

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
)

// PromotionPriceCalculatorImpl 促销活动价格计算器的完整实现
// 通过延迟初始化的方式解决循环依赖问题
type PromotionPriceCalculatorImpl struct {
	q                    *query.Query
	combinationRecordSvc promotion.CombinationRecordService
	bargainRecordSvc     *promotion.BargainRecordService
	pointActivitySvc     *promotion.PointActivityService
}

// NewPromotionPriceCalculatorImpl 创建促销活动价格计算器实现
func NewPromotionPriceCalculatorImpl(
	q *query.Query,
	combinationRecordSvc promotion.CombinationRecordService,
	bargainRecordSvc *promotion.BargainRecordService,
	pointActivitySvc *promotion.PointActivityService,
) PromotionPriceCalculator {
	return &PromotionPriceCalculatorImpl{
		q:                    q,
		combinationRecordSvc: combinationRecordSvc,
		bargainRecordSvc:     bargainRecordSvc,
		pointActivitySvc:     pointActivitySvc,
	}
}

// CalculateCombinationPrice 计算拼团活动价格
// 对齐Java: TradeCombinationActivityPriceCalculator
func (c *PromotionPriceCalculatorImpl) CalculateCombinationPrice(ctx context.Context, userID int64, activityID int64, headID int64, skuID int64, count int) (int, error) {
	// 1. 校验是否可以参与拼团
	_, combinationProd, err := c.combinationRecordSvc.ValidateCombinationRecord(ctx, userID, activityID, headID, skuID, count)
	if err != nil {
		return 0, err
	}

	// 2. 返回拼团价格
	return combinationProd.CombinationPrice * count, nil
}

// CalculateBargainPrice 计算砍价活动价格
// 对齐Java: TradeBargainActivityPriceCalculator
func (c *PromotionPriceCalculatorImpl) CalculateBargainPrice(ctx context.Context, userID int64, recordID int64, skuID int64, count int) (int, error) {
	// 1. 获取砍价记录
	bargainRecord, err := c.bargainRecordSvc.GetBargainRecord(ctx, recordID)
	if err != nil {
		return 0, err
	}

	// 2. 验证砍价记录是否有效
	if bargainRecord.UserID != userID {
		return 0, errors.New("砍价记录不属于当前用户")
	}
	if bargainRecord.SkuID != skuID {
		return 0, errors.New("砍价记录商品不匹配")
	}

	// 3. 返回砍价价格
	return bargainRecord.BargainPrice * count, nil
}

// CalculatePointPrice 计算积分活动价格
// 对齐Java: TradePointActivityPriceCalculator
func (c *PromotionPriceCalculatorImpl) CalculatePointPrice(ctx context.Context, activityID int64, spuID int64, skuID int64, count int) (int, error) {
	// 1. 校验是否可以参与积分商城活动
	pointProd, err := c.pointActivitySvc.ValidateJoinPointActivity(ctx, activityID, spuID, skuID, count)
	if err != nil {
		return 0, err
	}

	// 2. 计算积分商城价格
	// 情况一：纯积分兑换，价格为0
	if pointProd.Price == 0 {
		return 0, nil
	}

	// 情况二：积分+金额，返回金额部分
	return pointProd.Price * count, nil
}
