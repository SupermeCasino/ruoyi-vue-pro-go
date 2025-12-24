package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// Product Scope Constants (对齐 Java PromotionProductScopeEnum)
const (
	ProductScopeAll      = 1 // 全部商品 (Java: ALL)
	ProductScopeSpu      = 2 // 指定商品 (Java: SPU)
	ProductScopeCategory = 3 // 指定品类 (Java: CATEGORY)
)

// ProductScopeValues 商品范围值数组 (对齐 Java ARRAYS pattern)
var ProductScopeValues = []int{ProductScopeAll, ProductScopeSpu, ProductScopeCategory}

// IsValidProductScope 验证商品范围是否有效
func IsValidProductScope(scope int) bool {
	for _, v := range ProductScopeValues {
		if v == scope {
			return true
		}
	}
	return false
}

// IsProductScopeAll 判断是否为全部商品范围 (对齐 Java isAll 方法)
func IsProductScopeAll(scope int) bool {
	return scope == ProductScopeAll
}

// IsProductScopeSpu 判断是否为指定商品范围 (对齐 Java isSpu 方法)
func IsProductScopeSpu(scope int) bool {
	return scope == ProductScopeSpu
}

// IsProductScopeCategory 判断是否为指定品类范围 (对齐 Java isCategory 方法)
func IsProductScopeCategory(scope int) bool {
	return scope == ProductScopeCategory
}

// Discount Type Constants (对齐 Java PromotionDiscountTypeEnum)
const (
	DiscountTypePrice   = 1 // 满减 (Java: PRICE) - 具体金额
	DiscountTypePercent = 2 // 折扣 (Java: PERCENT) - 百分比
)

// DiscountTypeValues 折扣类型值数组 (对齐 Java ARRAYS pattern)
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

// Common Status Constants (使用 model.CommonStatus* 常量)
// 通用状态常量已在 internal/model/consts.go 中定义：
// - model.CommonStatusEnable = 0  (启用)
// - model.CommonStatusDisable = 1 (禁用)
// 这里提供便捷的验证函数

// CommonStatusValues 通用状态值数组 (使用现有的 model 常量)
var CommonStatusValues = []int{model.CommonStatusEnable, model.CommonStatusDisable}

// IsValidCommonStatus 验证通用状态是否有效
func IsValidCommonStatus(status int) bool {
	for _, v := range CommonStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// IsCommonStatusEnable 判断是否为启用状态 (对齐 Java isEnable 方法)
func IsCommonStatusEnable(status int) bool {
	return status == model.CommonStatusEnable
}

// IsCommonStatusDisable 判断是否为禁用状态 (对齐 Java isDisable 方法)
func IsCommonStatusDisable(status int) bool {
	return status == model.CommonStatusDisable
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

// Sender Type Constants (发送者类型常量，对齐客服系统)
const (
	SenderTypeMember = 1 // 用户发送 (对齐注释：1-用户)
	SenderTypeAdmin  = 2 // 客服发送 (对齐注释：2-客服)
)

// SenderTypeValues 发送者类型值数组
var SenderTypeValues = []int{SenderTypeMember, SenderTypeAdmin}

// IsValidSenderType 验证发送者类型是否有效
func IsValidSenderType(senderType int) bool {
	for _, v := range SenderTypeValues {
		if v == senderType {
			return true
		}
	}
	return false
}

// IsSenderTypeMember 判断是否为用户发送
func IsSenderTypeMember(senderType int) bool {
	return senderType == SenderTypeMember
}

// IsSenderTypeAdmin 判断是否为客服发送
func IsSenderTypeAdmin(senderType int) bool {
	return senderType == SenderTypeAdmin
}

// General Numeric Limits and Defaults (通用数值限制和默认值，Go特有扩展)
const (
	// 优惠券模板相关限制
	CouponTemplateTakeLimitCountMax = -1 // 不限制领取次数 (对齐 Java CouponTemplateDO.TIME_LIMIT_COUNT_MAX)

	// 分页默认值 (Java: 无对应枚举，Go扩展常量)
	DefaultPageSize = 10  // 默认分页大小 (Java: 无对应枚举，Go扩展常量)
	MaxPageSize     = 100 // 最大分页大小 (Java: 无对应枚举，Go扩展常量)

	// 价格相关 (Java: 无对应枚举，Go扩展常量)
	MinPrice = 1         // 最小价格（分） (Java: 无对应枚举，Go扩展常量)
	MaxPrice = 999999999 // 最大价格（分） (Java: 无对应枚举，Go扩展常量)

	// 折扣相关 (Java: 无对应枚举，Go扩展常量)
	MinDiscountPercent = 1    // 最小折扣百分比 (Java: 无对应枚举，Go扩展常量)
	MaxDiscountPercent = 9999 // 最大折扣百分比 (Java: 无对应枚举，Go扩展常量)

	// HTTP状态码常量 (Go特有扩展)
	HTTPStatusOK                  = 200 // 成功
	HTTPStatusBadRequest          = 400 // 请求参数错误
	HTTPStatusInternalServerError = 500 // 服务器内部错误
)

// HTTPStatusValues HTTP状态码值数组
var HTTPStatusValues = []int{HTTPStatusOK, HTTPStatusBadRequest, HTTPStatusInternalServerError}

// IsValidHTTPStatus 验证HTTP状态码是否有效
func IsValidHTTPStatus(status int) bool {
	for _, v := range HTTPStatusValues {
		if v == status {
			return true
		}
	}
	return false
}
