package consts

// PromotionTypeEnum 营销类型枚举
// 对应 Java 端的 PromotionTypeEnum
const (
	PromotionTypeSeckillActivity     = 1 // 秒杀活动
	PromotionTypeBargainActivity     = 2 // 砍价活动
	PromotionTypeCombinationActivity = 3 // 拼团活动
	PromotionTypeDiscountActivity    = 4 // 限时折扣
	PromotionTypeRewardActivity      = 5 // 满减送
	PromotionTypeMemberLevel         = 6 // 会员折扣
	PromotionTypeCoupon              = 7 // 优惠劵
	PromotionTypePoint               = 8 // 积分
)

// PromotionTypeValues 营销类型值数组
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

// Discount Type Constants (对齐 Java PromotionDiscountTypeEnum)
const (
	DiscountTypePrice   = 1 // 满减 (Java: PRICE) - 具体金额
	DiscountTypePercent = 2 // 折扣 (Java: PERCENT) - 百分比
)

// DiscountTypeValues 折扣类型值数组
var DiscountTypeValues = []int{DiscountTypePrice, DiscountTypePercent}

// IsValidDiscountType 验证折扣类型是否有效
func IsValidDiscountType(discountType int) bool {
	for _, v := range DiscountTypeValues {
		if v == discountType {
			return true
		}
	}
	return false
}

// IsDiscountTypePrice 判断是否为满减类型
func IsDiscountTypePrice(discountType int) bool {
	return discountType == DiscountTypePrice
}

// IsDiscountTypePercent 判断是否为折扣类型
func IsDiscountTypePercent(discountType int) bool {
	return discountType == DiscountTypePercent
}

// Condition Type Constants (条件类型常量，对齐 Java PromotionConditionTypeEnum)
const (
	ConditionTypePrice = 10 // 满 N 元
	ConditionTypeCount = 20 // 满 N 件
)

// ConditionTypeValues 条件类型值数组
var ConditionTypeValues = []int{ConditionTypePrice, ConditionTypeCount}

// IsValidConditionType 验证条件类型是否有效
func IsValidConditionType(conditionType int) bool {
	for _, v := range ConditionTypeValues {
		if v == conditionType {
			return true
		}
	}
	return false
}

// IsConditionTypePrice 判断是否为满金额条件
func IsConditionTypePrice(conditionType int) bool {
	return conditionType == ConditionTypePrice
}

// IsConditionTypeCount 判断是否为满数量条件
func IsConditionTypeCount(conditionType int) bool {
	return conditionType == ConditionTypeCount
}

// KeFuMessageContentType 客服消息类型 (对齐 Java: KeFuMessageContentTypeEnum)
const (
	// KeFuMessageContentTypeText 文本消息
	KeFuMessageContentTypeText = 1
	// KeFuMessageContentTypeImage 图片消息
	KeFuMessageContentTypeImage = 2
	// KeFuMessageContentTypeVoice 语音消息
	KeFuMessageContentTypeVoice = 3
	// KeFuMessageContentTypeVideo 视频消息
	KeFuMessageContentTypeVideo = 4
	// KeFuMessageContentTypeSystem 系统消息
	KeFuMessageContentTypeSystem = 5
	// KeFuMessageContentTypeProduct 商品消息
	KeFuMessageContentTypeProduct = 10
	// KeFuMessageContentTypeOrder 订单消息
	KeFuMessageContentTypeOrder = 11
)

// DiyPageEnum 装修页面枚举 (对齐 Java: DiyPageEnum)
const (
	// DiyPageIndex 首页
	DiyPageIndex = 1
	// DiyPageMy 我的
	DiyPageMy = 2
)
