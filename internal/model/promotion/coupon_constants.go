package promotion

// Coupon Status Constants (对齐 Java CouponStatusEnum)
const (
	CouponStatusUnused  = 1 // 未使用 (Java: UNUSED)
	CouponStatusUsed    = 2 // 已使用 (Java: USED)
	CouponStatusExpired = 3 // 已过期 (Java: EXPIRE)
)

// CouponStatusValues 优惠券状态值数组 (对齐 Java ARRAYS pattern)
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
	CouponTakeTypeUser     = 1 // 直接领取 (Java: USER) - 用户可在首页、每日领劵直接领取
	CouponTakeTypeAdmin    = 2 // 指定发放 (Java: ADMIN) - 后台指定会员赠送优惠劵
	CouponTakeTypeRegister = 3 // 新人券 (Java: REGISTER) - 注册时自动领取
)

// CouponTakeTypeValues 优惠券领取类型值数组 (对齐 Java ARRAYS pattern)
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

// IsCouponTakeTypeUser 判断是否为用户领取类型 (对齐 Java isUser 方法)
func IsCouponTakeTypeUser(takeType int) bool {
	return takeType == CouponTakeTypeUser
}

// Template Validity Type Constants (对齐 Java CouponTemplateValidityTypeEnum)
const (
	CouponValidityTypeDate = 1 // 固定日期 (Java: DATE)
	CouponValidityTypeTerm = 2 // 领取之后 (Java: TERM)
)

// CouponValidityTypeValues 优惠券有效期类型值数组 (对齐 Java ARRAYS pattern)
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

// IsCouponValidityTypeDate 判断是否为固定日期类型 (对齐 Java DATE.getType().equals())
func IsCouponValidityTypeDate(validityType int) bool {
	return validityType == CouponValidityTypeDate
}

// IsCouponValidityTypeTerm 判断是否为领取之后类型 (对齐 Java TERM.getType().equals())
func IsCouponValidityTypeTerm(validityType int) bool {
	return validityType == CouponValidityTypeTerm
}
