package consts

// Activity Status Constants (对齐 Java PromotionActivityStatusEnum)
const (
	ActivityStatusWait  = 10 // 未开始
	ActivityStatusRun   = 20 // 进行中
	ActivityStatusEnd   = 30 // 已结束
	ActivityStatusClose = 40 // 已关闭
)

// ActivityStatusValues 活动状态值数组
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

// Seckill Activity Status Constants (秒杀活动状态)
const (
	SeckillActivityStatusWait  = 10 // 未开始
	SeckillActivityStatusRun   = 20 // 进行中
	SeckillActivityStatusEnd   = 30 // 已结束
	SeckillActivityStatusClose = 40 // 已关闭
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

// BargainRecordStatus 砍价记录状态枚举
// 对应 Java: BargainRecordStatusEnum
const (
	// BargainRecordStatusInProgress 砍价中
	BargainRecordStatusInProgress = 1
	// BargainRecordStatusSuccess 砍价成功
	BargainRecordStatusSuccess = 2
	// BargainRecordStatusFailed 砍价失败
	BargainRecordStatusFailed = 3
)

// BargainRecordStatusValues 砍价记录状态值数组
var BargainRecordStatusValues = []int{
	BargainRecordStatusInProgress,
	BargainRecordStatusSuccess,
	BargainRecordStatusFailed,
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

// PromotionCombinationRecordStatus 拼团记录状态枚举
const (
	// PromotionCombinationRecordStatusInProgress 进行中
	PromotionCombinationRecordStatusInProgress = 0
	// PromotionCombinationRecordStatusSuccess 成功
	PromotionCombinationRecordStatusSuccess = 1
	// PromotionCombinationRecordStatusFailed 失败
	PromotionCombinationRecordStatusFailed = 2
)

// CombinationRecordStatusValues 拼团记录状态值数组
var CombinationRecordStatusValues = []int{
	PromotionCombinationRecordStatusInProgress,
	PromotionCombinationRecordStatusSuccess,
	PromotionCombinationRecordStatusFailed,
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

// 拼团模块常量
const (
	// PromotionCombinationRecordHeadIDGroup 团长 ID (0 代表团长)
	PromotionCombinationRecordHeadIDGroup = 0
	// AppCombinationRecordSummaryAvatarCount 拼团摘要头像展示数量
	AppCombinationRecordSummaryAvatarCount = 7
)
