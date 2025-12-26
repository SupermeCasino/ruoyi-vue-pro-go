package model

// CommonStatus 通用状态枚举
// 对应 Java: CommonStatusEnum
const (
	// CommonStatusEnable 开启
	CommonStatusEnable = 0
	// CommonStatusDisable 禁用
	CommonStatusDisable = 1
)

// ProductSpuStatus 商品 SPU 状态枚举
// 对应 Java: ProductSpuStatusEnum
const (
	// ProductSpuStatusRecycle 回收站
	ProductSpuStatusRecycle = -1
	// ProductSpuStatusDisable 下架
	ProductSpuStatusDisable = 0
	// ProductSpuStatusEnable 上架
	ProductSpuStatusEnable = 1
)

// BargainRecordStatus 砍价记录状态枚举
// 对应 Java: BargainRecordStatusEnum
const (
	// BargainRecordStatusInProgress 进行中
	BargainRecordStatusInProgress = 0
	// BargainRecordStatusSuccess 成功
	BargainRecordStatusSuccess = 1
	// BargainRecordStatusFailed 失败
	BargainRecordStatusFailed = 2
)

const (
	// PromotionCombinationRecordStatusInProgress 进行中
	PromotionCombinationRecordStatusInProgress = 0
	// PromotionCombinationRecordStatusSuccess 成功
	PromotionCombinationRecordStatusSuccess = 1
	// PromotionCombinationRecordStatusFailed 失败
	PromotionCombinationRecordStatusFailed = 2
)

// 拼团模块常量
const (
	// PromotionCombinationRecordHeadIDGroup 团长 ID (0 代表团长)
	PromotionCombinationRecordHeadIDGroup = 0
	// AppCombinationRecordSummaryAvatarCount 拼团摘要头像展示数量
	AppCombinationRecordSummaryAvatarCount = 7
)
