package promotion

import "github.com/wxlbd/ruoyi-mall-go/internal/model"

// Activity Status Constants (对齐 Java PromotionActivityStatusEnum)
const (
	ActivityStatusWait  = 10 // 未开始 (Java: WAIT)
	ActivityStatusRun   = 20 // 进行中 (Java: RUN)
	ActivityStatusEnd   = 30 // 已结束 (Java: END)
	ActivityStatusClose = 40 // 已关闭 (Java: CLOSE)
)

// ActivityStatusValues 活动状态值数组 (对齐 Java ARRAYS pattern)
var ActivityStatusValues = []int{ActivityStatusWait, ActivityStatusRun, ActivityStatusEnd, ActivityStatusClose}

// IsValidActivityStatus 验证活动状态是否有效
func IsValidActivityStatus(status int) bool {
	for _, v := range ActivityStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// IsActivityStatusWait 判断是否为等待状态
func IsActivityStatusWait(status int) bool {
	return status == ActivityStatusWait
}

// IsActivityStatusRun 判断是否为进行中状态
func IsActivityStatusRun(status int) bool {
	return status == ActivityStatusRun
}

// IsActivityStatusEnd 判断是否为已结束状态
func IsActivityStatusEnd(status int) bool {
	return status == ActivityStatusEnd
}

// IsActivityStatusClose 判断是否为已关闭状态
func IsActivityStatusClose(status int) bool {
	return status == ActivityStatusClose
}

// Promotion Type Constants (使用现有的 PromotionTypeEnum 常量)
// 这些常量已在 promotion_type.go 中定义，这里提供验证函数和值数组

// PromotionTypeValues 营销类型值数组 (使用 promotion_type.go 中的常量)
var PromotionTypeValues = []int{
	PromotionTypeSeckillActivity, PromotionTypeBargainActivity, PromotionTypeCombinationActivity,
	PromotionTypeDiscountActivity, PromotionTypeRewardActivity, PromotionTypeMemberLevel,
	PromotionTypeCoupon, PromotionTypePoint,
}

// IsValidPromotionType 验证营销类型是否有效
func IsValidPromotionType(promotionType int) bool {
	for _, v := range PromotionTypeValues {
		if v == promotionType {
			return true
		}
	}
	return false
}

// Combination Record Status Constants (使用 model 中现有的常量)
// 拼团记录状态常量已在 internal/model/consts.go 中定义：
// - model.PromotionCombinationRecordStatusInProgress = 0 (进行中)
// - model.PromotionCombinationRecordStatusSuccess = 1    (拼团成功)
// - model.PromotionCombinationRecordStatusFailed = 2     (拼团失败)

// CombinationRecordStatusValues 拼团记录状态值数组 (使用现有的 model 常量)
var CombinationRecordStatusValues = []int{
	model.PromotionCombinationRecordStatusInProgress,
	model.PromotionCombinationRecordStatusSuccess,
	model.PromotionCombinationRecordStatusFailed,
}

// IsValidCombinationRecordStatus 验证拼团记录状态是否有效
func IsValidCombinationRecordStatus(status int) bool {
	for _, v := range CombinationRecordStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// Bargain Record Status Constants (使用 model 中现有的常量)
// 砍价记录状态常量已在 internal/model/consts.go 中定义：
// - model.BargainRecordStatusInProgress = 0 (砍价中)
// - model.BargainRecordStatusSuccess = 1    (砍价成功)
// - model.BargainRecordStatusFailed = 2     (砍价失败)

// BargainRecordStatusValues 砍价记录状态值数组 (使用现有的 model 常量)
var BargainRecordStatusValues = []int{
	model.BargainRecordStatusInProgress,
	model.BargainRecordStatusSuccess,
	model.BargainRecordStatusFailed,
}

// IsValidBargainRecordStatus 验证砍价记录状态是否有效
func IsValidBargainRecordStatus(status int) bool {
	for _, v := range BargainRecordStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// Seckill Activity Status Constants (秒杀活动状态，对齐 Java PromotionActivityStatusEnum)
const (
	SeckillActivityStatusWait  = 10 // 未开始 (对齐 Java PromotionActivityStatusEnum.WAIT)
	SeckillActivityStatusRun   = 20 // 进行中 (对齐 Java PromotionActivityStatusEnum.RUN)
	SeckillActivityStatusEnd   = 30 // 已结束 (对齐 Java PromotionActivityStatusEnum.END)
	SeckillActivityStatusClose = 40 // 已关闭 (对齐 Java PromotionActivityStatusEnum.CLOSE)
)

// SeckillActivityStatusValues 秒杀活动状态值数组
var SeckillActivityStatusValues = []int{SeckillActivityStatusWait, SeckillActivityStatusRun, SeckillActivityStatusEnd, SeckillActivityStatusClose}

// IsValidSeckillActivityStatus 验证秒杀活动状态是否有效
func IsValidSeckillActivityStatus(status int) bool {
	for _, v := range SeckillActivityStatusValues {
		if v == status {
			return true
		}
	}
	return false
}
