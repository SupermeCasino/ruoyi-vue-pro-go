package promotion

// PromotionTypeEnum 营销类型枚举
// 对应 Java 端的 PromotionTypeEnum
const (
	PromotionTypeSeckillActivity     = 1 // 秒杀活动 (Java: SECKILL_ACTIVITY = 1)
	PromotionTypeBargainActivity     = 2 // 砍价活动 (Java: BARGAIN_ACTIVITY = 2)
	PromotionTypeCombinationActivity = 3 // 拼团活动 (Java: COMBINATION_ACTIVITY = 3)
	PromotionTypeDiscountActivity    = 4 // 限时折扣 (Java: DISCOUNT_ACTIVITY = 4)
	PromotionTypeRewardActivity      = 5 // 满减送 (Java: REWARD_ACTIVITY = 5)
	PromotionTypeMemberLevel         = 6 // 会员折扣 (Java: MEMBER_LEVEL = 6)
	PromotionTypeCoupon              = 7 // 优惠劵 (Java: COUPON = 7)
	PromotionTypePoint               = 8 // 积分 (Java: POINT = 8)
)
