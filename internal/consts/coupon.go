package consts

// Coupon Status Constants (对齐 Java CouponStatusEnum)
const (
	CouponStatusUnused  = 1 // 未使用
	CouponStatusUsed    = 2 // 已使用
	CouponStatusExpired = 3 // 已过期
)

// CouponStatusValues 优惠券状态值数组
var CouponStatusValues = []int{CouponStatusUnused, CouponStatusUsed, CouponStatusExpired}

// IsValidCouponStatus 验证优惠券状态是否有效
func IsValidCouponStatus(status int) bool {
	for _, v := range CouponStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// Coupon Take Type Constants (对齐 Java CouponTakeTypeEnum)
const (
	CouponTakeTypeUser     = 1 // 直接领取 - 用户可在首页、每日领劵直接领取
	CouponTakeTypeAdmin    = 2 // 指定发放 - 后台指定会员赠送优惠劵
	CouponTakeTypeRegister = 3 // 新人券 - 注册时自动领取
)

// CouponTakeTypeValues 优惠券领取类型值数组
var CouponTakeTypeValues = []int{CouponTakeTypeUser, CouponTakeTypeAdmin, CouponTakeTypeRegister}

// IsValidCouponTakeType 验证优惠券领取类型是否有效
func IsValidCouponTakeType(takeType int) bool {
	for _, v := range CouponTakeTypeValues {
		if v == takeType {
			return true
		}
	}
	return false
}

// IsCouponTakeTypeUser 判断是否为用户领取类型
func IsCouponTakeTypeUser(takeType int) bool {
	return takeType == CouponTakeTypeUser
}

// Template Validity Type Constants (对齐 Java CouponTemplateValidityTypeEnum)
const (
	CouponValidityTypeDate = 1 // 固定日期
	CouponValidityTypeTerm = 2 // 领取之后
)

// CouponValidityTypeValues 优惠券有效期类型值数组
var CouponValidityTypeValues = []int{CouponValidityTypeDate, CouponValidityTypeTerm}

// IsValidCouponValidityType 验证优惠券有效期类型是否有效
func IsValidCouponValidityType(validityType int) bool {
	for _, v := range CouponValidityTypeValues {
		if v == validityType {
			return true
		}
	}
	return false
}

// IsCouponValidityTypeDate 判断是否为固定日期类型
func IsCouponValidityTypeDate(validityType int) bool {
	return validityType == CouponValidityTypeDate
}

// IsCouponValidityTypeTerm 判断是否为领取之后类型
func IsCouponValidityTypeTerm(validityType int) bool {
	return validityType == CouponValidityTypeTerm
}

// 优惠券模板相关限制
const (
	CouponTemplateTakeLimitCountMax = -1 // 不限制领取次数
)
