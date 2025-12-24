package promotion

import (
	"testing"
)

// TestJavaAlignmentValidation 测试所有常量与Java版本的对齐情况
// Feature: promotion-constants-refactor, Property 3: For any defined constant, its value should exactly match the corresponding Java enum value
func TestJavaAlignmentValidation(t *testing.T) {
	validations := ValidateJavaAlignment()

	for _, validation := range validations {
		t.Run(validation.ConstantName, func(t *testing.T) {
			if !validation.IsAligned {
				t.Errorf("常量 %s 与Java版本不对齐: Go值=%d, Java值=%d (Java枚举: %s)",
					validation.ConstantName, validation.GoValue, validation.JavaValue, validation.JavaEnum)
			}
		})
	}
}

// TestAllConstantsAligned 测试所有常量是否都与Java版本对齐
func TestAllConstantsAligned(t *testing.T) {
	if !IsAllConstantsAligned() {
		misaligned := GetMisalignedConstants()
		t.Errorf("发现 %d 个常量与Java版本不对齐:", len(misaligned))
		for _, validation := range misaligned {
			t.Errorf("  - %s: Go值=%d, Java值=%d",
				validation.ConstantName, validation.GoValue, validation.JavaValue)
		}
	}
}

// TestSpecificJavaAlignmentRequirements 测试特定的Java对齐需求
func TestSpecificJavaAlignmentRequirements(t *testing.T) {
	testCases := []struct {
		name        string
		goValue     int
		javaValue   int
		requirement string
	}{
		// Requirement 1.3: 优惠券状态常量对齐Java CouponStatusEnum
		{"CouponStatusUnused", CouponStatusUnused, 1, "1.3"},
		{"CouponStatusUsed", CouponStatusUsed, 2, "1.3"},
		{"CouponStatusExpired", CouponStatusExpired, 3, "1.3"},

		// Requirement 2.3: 优惠券领取类型常量对齐Java CouponTakeTypeEnum
		{"CouponTakeTypeUser", CouponTakeTypeUser, 1, "2.3"},
		{"CouponTakeTypeAdmin", CouponTakeTypeAdmin, 2, "2.3"},
		{"CouponTakeTypeRegister", CouponTakeTypeRegister, 3, "2.3"},

		// Requirement 3.3: 商品范围常量对齐Java PromotionProductScopeEnum
		{"ProductScopeAll", ProductScopeAll, 1, "3.3"},
		{"ProductScopeSpu", ProductScopeSpu, 2, "3.3"},
		{"ProductScopeCategory", ProductScopeCategory, 3, "3.3"},

		// Requirement 4.3: 折扣类型常量对齐Java PromotionDiscountTypeEnum
		{"DiscountTypePrice", DiscountTypePrice, 1, "4.3"},
		{"DiscountTypePercent", DiscountTypePercent, 2, "4.3"},

		// Requirement 5.3: 有效期类型常量对齐Java CouponTemplateValidityTypeEnum
		{"CouponValidityTypeDate", CouponValidityTypeDate, 1, "5.3"},
		{"CouponValidityTypeTerm", CouponValidityTypeTerm, 2, "5.3"},

		// Requirement 6.3: 活动状态常量对齐Java PromotionActivityStatusEnum
		{"ActivityStatusWait", ActivityStatusWait, 10, "6.3"},
		{"ActivityStatusRun", ActivityStatusRun, 20, "6.3"},
		{"ActivityStatusEnd", ActivityStatusEnd, 30, "6.3"},
		{"ActivityStatusClose", ActivityStatusClose, 40, "6.3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.goValue != tc.javaValue {
				t.Errorf("Requirement %s: 常量 %s 与Java版本不对齐: Go值=%d, 期望Java值=%d",
					tc.requirement, tc.name, tc.goValue, tc.javaValue)
			}
		})
	}
}

// TestJavaAlignmentReport 测试Java对齐报告生成
func TestJavaAlignmentReport(t *testing.T) {
	// 这个测试主要用于手动验证，打印对齐报告
	if testing.Verbose() {
		PrintJavaAlignmentReport()
	}

	// 验证报告生成不会出错
	validations := ValidateJavaAlignment()
	if len(validations) == 0 {
		t.Error("验证列表不应该为空")
	}
}
