package consts

// Banner Position Constants (对齐 Java BannerPositionEnum)
const (
	BannerPositionHome        = 1 // 首页
	BannerPositionSeckill     = 2 // 秒杀活动页
	BannerPositionCombination = 3 // 砍价活动页
	BannerPositionDiscount    = 4 // 限时折扣页
	BannerPositionReward      = 5 // 满减送页
)

// BannerPositionValues Banner位置值数组
var BannerPositionValues = []int{
	BannerPositionHome, BannerPositionSeckill, BannerPositionCombination,
	BannerPositionDiscount, BannerPositionReward,
}

// IsValidBannerPosition 验证Banner位置是否有效
func IsValidBannerPosition(position int) bool {
	for _, v := range BannerPositionValues {
		if v == position {
			return true
		}
	}
	return false
}

// IsBannerPositionHome 判断是否为首页位置
func IsBannerPositionHome(position int) bool {
	return position == BannerPositionHome
}

// IsBannerPositionSeckill 判断是否为秒杀活动页位置
func IsBannerPositionSeckill(position int) bool {
	return position == BannerPositionSeckill
}

// IsBannerPositionCombination 判断是否为砍价活动页位置
func IsBannerPositionCombination(position int) bool {
	return position == BannerPositionCombination
}

// IsBannerPositionDiscount 判断是否为限时折扣页位置
func IsBannerPositionDiscount(position int) bool {
	return position == BannerPositionDiscount
}

// IsBannerPositionReward 判断是否为满减送页位置
func IsBannerPositionReward(position int) bool {
	return position == BannerPositionReward
}

// BannerStatusValues Banner状态值数组
var BannerStatusValues = []int{CommonStatusEnable, CommonStatusDisable}

// IsValidBannerStatus 验证Banner状态是否有效
func IsValidBannerStatus(status int) bool {
	for _, v := range BannerStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// IsBannerStatusEnable 判断是否为启用状态
func IsBannerStatusEnable(status int) bool {
	return status == CommonStatusEnable
}

// IsBannerStatusDisable 判断是否为禁用状态
func IsBannerStatusDisable(status int) bool {
	return status == CommonStatusDisable
}

// Banner Display Priority Constants
const (
	BannerPriorityLow    = 1  // 低优先级
	BannerPriorityNormal = 5  // 普通优先级
	BannerPriorityHigh   = 10 // 高优先级
)

// BannerPriorityValues Banner优先级值数组
var BannerPriorityValues = []int{BannerPriorityLow, BannerPriorityNormal, BannerPriorityHigh}

// IsValidBannerPriority 验证Banner优先级是否有效
func IsValidBannerPriority(priority int) bool {
	for _, v := range BannerPriorityValues {
		if v == priority {
			return true
		}
	}
	return false
}

// Banner Type Constants
const (
	BannerTypeImage = 1 // 图片Banner
	BannerTypeVideo = 2 // 视频Banner
)

// BannerTypeValues Banner类型值数组
var BannerTypeValues = []int{BannerTypeImage, BannerTypeVideo}

// IsValidBannerType 验证Banner类型是否有效
func IsValidBannerType(bannerType int) bool {
	for _, v := range BannerTypeValues {
		if v == bannerType {
			return true
		}
	}
	return false
}

// IsBannerTypeImage 判断是否为图片Banner
func IsBannerTypeImage(bannerType int) bool {
	return bannerType == BannerTypeImage
}

// IsBannerTypeVideo 判断是否为视频Banner
func IsBannerTypeVideo(bannerType int) bool {
	return bannerType == BannerTypeVideo
}
