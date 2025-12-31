package trade

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

// PromotionPriceCalculator 促销活动价格计算器接口
// 用于解决循环依赖问题，将促销活动价格计算逻辑抽象为接口
type PromotionPriceCalculator interface {
	// CalculateCombinationPrice 计算拼团活动价格
	CalculateCombinationPrice(ctx context.Context, userID int64, activityID int64, headID int64, skuID int64, count int) (int, error)

	// CalculateBargainPrice 计算砍价活动价格
	CalculateBargainPrice(ctx context.Context, userID int64, recordID int64, skuID int64, count int) (int, error)

	// CalculatePointPrice 计算积分活动价格
	CalculatePointPrice(ctx context.Context, activityID int64, spuID int64, skuID int64, count int) (int, error)
}

// DefaultPromotionPriceCalculator 默认的促销活动价格计算器实现
type DefaultPromotionPriceCalculator struct {
	// 使用延迟初始化来避免循环依赖
	q *query.Query
}

// NewDefaultPromotionPriceCalculator 创建默认促销活动价格计算器
func NewDefaultPromotionPriceCalculator(q *query.Query) *DefaultPromotionPriceCalculator {
	return &DefaultPromotionPriceCalculator{q: q}
}

// CalculateCombinationPrice 计算拼团活动价格
func (c *DefaultPromotionPriceCalculator) CalculateCombinationPrice(ctx context.Context, userID int64, activityID int64, headID int64, skuID int64, count int) (int, error) {
	// 直接查询数据库获取拼团商品价格，避免循环依赖
	combinationProduct, err := c.q.PromotionCombinationProduct.WithContext(ctx).
		Where(c.q.PromotionCombinationProduct.ActivityID.Eq(activityID)).
		Where(c.q.PromotionCombinationProduct.SkuID.Eq(skuID)).
		First()
	if err != nil {
		return 0, err
	}

	return combinationProduct.CombinationPrice * count, nil
}

// CalculateBargainPrice 计算砍价活动价格
func (c *DefaultPromotionPriceCalculator) CalculateBargainPrice(ctx context.Context, userID int64, recordID int64, skuID int64, count int) (int, error) {
	// 直接查询数据库获取砍价记录，避免循环依赖
	bargainRecord, err := c.q.PromotionBargainRecord.WithContext(ctx).
		Where(c.q.PromotionBargainRecord.ID.Eq(recordID)).
		First()
	if err != nil {
		return 0, err
	}

	// 验证砍价记录
	if bargainRecord.UserID != userID {
		return 0, errors.New("砍价记录不属于当前用户")
	}
	if bargainRecord.SkuID != skuID {
		return 0, errors.New("砍价记录商品不匹配")
	}

	return bargainRecord.BargainPrice * count, nil
}

// CalculatePointPrice 计算积分活动价格
func (c *DefaultPromotionPriceCalculator) CalculatePointPrice(ctx context.Context, activityID int64, spuID int64, skuID int64, count int) (int, error) {
	// 直接查询数据库获取积分商品价格，避免循环依赖
	pointProduct, err := c.q.PromotionPointProduct.WithContext(ctx).
		Where(c.q.PromotionPointProduct.ActivityID.Eq(activityID)).
		Where(c.q.PromotionPointProduct.SkuID.Eq(skuID)).
		First()
	if err != nil {
		return 0, err
	}

	// 积分商城：如果Price为0则纯积分兑换，否则积分+金额
	return pointProduct.Price * count, nil
}
