package promotion

import (
	"testing"
)

func TestNamingValidator_ValidateNamingConventions(t *testing.T) {
	nv := NewNamingValidator()

	// 测试验证命名约定
	rootPath := "."
	violations, err := nv.ValidateNamingConventions(rootPath)

	if err != nil {
		t.Fatalf("验证命名约定失败: %v", err)
	}

	t.Logf("发现 %d 个命名违规", len(violations))

	// 输出违规详情
	for i, violation := range violations {
		if i >= 5 { // 只输出前5个
			break
		}
		t.Logf("违规 %d: 类型=%s, 名称=%s, 文件=%s, 行=%d",
			i+1, violation.Type, violation.Name, violation.File, violation.Line)
		t.Logf("  期望: %s, 实际: %s", violation.Expected, violation.Actual)
		t.Logf("  描述: %s", violation.Description)
	}
}

func TestNamingValidator_DetectDuplicateConstants(t *testing.T) {
	nv := NewNamingValidator()

	// 测试检测重复常量
	rootPath := "."
	duplicates, err := nv.DetectDuplicateConstants(rootPath)

	if err != nil {
		t.Fatalf("检测重复常量失败: %v", err)
	}

	t.Logf("发现 %d 个重复常量", len(duplicates))

	// 输出重复常量详情
	for i, duplicate := range duplicates {
		if i >= 5 { // 只输出前5个
			break
		}
		t.Logf("重复常量 %d: 名称=%s, 值=%d, 出现次数=%d",
			i+1, duplicate.Name, duplicate.Value, duplicate.Occurrences)
		t.Logf("  文件: %v", duplicate.Files)
	}
}

func TestNamingValidator_GetConstantType(t *testing.T) {
	nv := NewNamingValidator()

	testCases := []struct {
		name     string
		expected string
		desc     string
	}{
		{"CouponStatusUnused", "coupon_status", "优惠券状态常量"},
		{"CouponTakeTypeUser", "coupon_take_type", "优惠券领取类型常量"},
		{"ActivityStatusWait", "activity_status", "活动状态常量"},
		{"ProductScopeAll", "product_scope", "商品范围常量"},
		{"BannerPositionHome", "banner_position", "Banner位置常量"},
		{"PromotionTypeCoupon", "promotion_type", "营销类型常量"},
		{"UnknownConstant", "", "未知常量类型"},
	}

	for _, tc := range testCases {
		result := nv.getConstantType(tc.name)
		if result != tc.expected {
			t.Errorf("测试失败: %s - 期望 '%s', 实际 '%s'", tc.desc, tc.expected, result)
		} else {
			t.Logf("测试通过: %s -> %s", tc.name, result)
		}
	}
}

func TestNamingValidator_GetFunctionType(t *testing.T) {
	nv := NewNamingValidator()

	testCases := []struct {
		name     string
		expected string
		desc     string
	}{
		{"IsValidCouponStatus", "validation_function", "验证函数"},
		{"IsCouponStatusUnused", "predicate_function", "判断函数"},
		{"IsActivityStatusWait", "predicate_function", "判断函数"},
		{"CreateOrder", "", "普通函数"},
		{"GetUserInfo", "", "普通函数"},
	}

	for _, tc := range testCases {
		result := nv.getFunctionType(tc.name)
		if result != tc.expected {
			t.Errorf("测试失败: %s - 期望 '%s', 实际 '%s'", tc.desc, tc.expected, result)
		} else {
			t.Logf("测试通过: %s -> %s", tc.name, result)
		}
	}
}

func TestNamingValidator_GetVariableType(t *testing.T) {
	nv := NewNamingValidator()

	testCases := []struct {
		name     string
		expected string
		desc     string
	}{
		{"CouponStatusValues", "values_array", "值数组变量"},
		{"ActivityStatusValues", "values_array", "值数组变量"},
		{"UserInfo", "", "普通变量"},
		{"OrderData", "", "普通变量"},
	}

	for _, tc := range testCases {
		result := nv.getVariableType(tc.name)
		if result != tc.expected {
			t.Errorf("测试失败: %s - 期望 '%s', 实际 '%s'", tc.desc, tc.expected, result)
		} else {
			t.Logf("测试通过: %s -> %s", tc.name, result)
		}
	}
}

func TestNamingValidator_ToPascalCase(t *testing.T) {
	nv := NewNamingValidator()

	testCases := []struct {
		input    string
		expected string
		desc     string
	}{
		{"unused", "Unused", "小写转PascalCase"},
		{"user", "User", "小写转PascalCase"},
		{"WAIT", "WAIT", "大写保持不变"},
		{"", "", "空字符串"},
		{"a", "A", "单字符"},
	}

	for _, tc := range testCases {
		result := nv.toPascalCase(tc.input)
		if result != tc.expected {
			t.Errorf("测试失败: %s - 期望 '%s', 实际 '%s'", tc.desc, tc.expected, result)
		} else {
			t.Logf("测试通过: %s -> %s", tc.input, result)
		}
	}
}

func TestNamingValidator_IsReasonableDuplication(t *testing.T) {
	nv := NewNamingValidator()

	testCases := []struct {
		constants []string
		value     int
		expected  bool
		desc      string
	}{
		{[]string{"CouponStatusUnused", "ActivityStatusWait"}, 1, true, "不同模块的相同值是合理的"},
		{[]string{"CouponStatusUnused", "CouponStatusUsed"}, 1, false, "同一模块的不同常量不合理"},
		{[]string{"CommonStatusEnable", "BannerStatusEnable"}, 0, true, "通用状态值的重复是合理的"},
		{[]string{"ActivityStatusWait", "SeckillActivityStatusWait"}, 10, true, "活动状态值的重复是合理的"},
	}

	for _, tc := range testCases {
		result := nv.isReasonableDuplication(tc.constants, tc.value)
		if result != tc.expected {
			t.Errorf("测试失败: %s - 期望 %v, 实际 %v", tc.desc, tc.expected, result)
		} else {
			t.Logf("测试通过: %s", tc.desc)
		}
	}
}

func TestNamingValidator_GetConstantModule(t *testing.T) {
	nv := NewNamingValidator()

	testCases := []struct {
		constName string
		expected  string
		desc      string
	}{
		{"CouponStatusUnused", "coupon", "优惠券模块"},
		{"ActivityStatusWait", "activity", "活动模块"},
		{"BannerPositionHome", "banner", "Banner模块"},
		{"ProductScopeAll", "product", "商品模块"},
		{"DiscountTypePrice", "discount", "折扣模块"},
		{"PromotionTypeCoupon", "promotion", "营销模块"},
		{"CommonStatusEnable", "common", "通用模块"},
		{"UnknownConstant", "common", "未知常量归为通用模块"},
	}

	for _, tc := range testCases {
		result := nv.getConstantModule(tc.constName)
		if result != tc.expected {
			t.Errorf("测试失败: %s - 期望 '%s', 实际 '%s'", tc.desc, tc.expected, result)
		} else {
			t.Logf("测试通过: %s -> %s", tc.constName, result)
		}
	}
}

func TestNamingValidator_RunNamingValidation(t *testing.T) {
	nv := NewNamingValidator()

	// 运行完整的命名验证
	rootPath := "."
	violations, duplicates, err := nv.RunNamingValidation(rootPath)

	if err != nil {
		t.Fatalf("运行命名验证失败: %v", err)
	}

	t.Logf("命名验证结果:")
	t.Logf("  命名违规: %d", len(violations))
	t.Logf("  重复常量: %d", len(duplicates))

	// 输出总结
	if len(violations) == 0 && len(duplicates) == 0 {
		t.Log("✅ 命名验证通过")
	} else {
		t.Log("❌ 发现命名问题")

		if len(violations) > 0 {
			t.Log("命名违规详情:")
			for i, violation := range violations {
				if i >= 3 { // 只输出前3个
					break
				}
				t.Logf("  %d. %s: %s", i+1, violation.Name, violation.Description)
			}
		}

		if len(duplicates) > 0 {
			t.Log("重复常量详情:")
			for i, duplicate := range duplicates {
				if i >= 3 { // 只输出前3个
					break
				}
				t.Logf("  %d. %s (值=%d): 出现%d次",
					i+1, duplicate.Name, duplicate.Value, duplicate.Occurrences)
			}
		}
	}
}
