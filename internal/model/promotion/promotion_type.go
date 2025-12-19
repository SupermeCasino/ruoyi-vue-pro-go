package promotion

// PromotionTypeEnum 营销类型枚举
// 对应 Java 端的 PromotionTypeEnum
const (
	PromotionTypeCombinationActivity = 1 // 拼团活动 (Java: COMBINATION_ACTIVITY = 1)
	PromotionTypeSeckillActivity     = 2 // 秒杀活动 (Java: SECKILL_ACTIVITY = 2)
	PromotionTypeBargainActivity     = 3 // 砍价活动 (Java: BARGAIN_ACTIVITY = 3)
	PromotionTypeDiscountActivity    = 4 // 限时折扣
	PromotionTypeRewardActivity      = 5 // 满减送
	PromotionTypeMemberLevel         = 6 // 会员折扣
	PromotionTypeCoupon              = 7 // 优惠劵
	PromotionTypePoint               = 8 // 积分
)
