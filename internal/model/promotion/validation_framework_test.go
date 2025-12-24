package promotion

import (
	"path/filepath"
	"testing"
)

func TestValidationFramework_ScanSourceForNumericLiterals(t *testing.T) {
	vf := NewValidationFramework()

	// 测试扫描当前目录
	rootPath := "."
	literals, err := vf.ScanSourceForNumericLiterals(rootPath)

	if err != nil {
		t.Fatalf("扫描数字字面量失败: %v", err)
	}

	t.Logf("发现 %d 个可疑的数字字面量", len(literals))

	// 输出前几个结果用于调试
	for i, literal := range literals {
		if i >= 5 { // 只输出前5个
			break
		}
		t.Logf("数字字面量 %d: 值=%d, 位置=%s, 上下文=%s",
			i+1, literal.Value, literal.Position, literal.Context)
	}
}

func TestValidationFramework_ValidateJavaAlignment(t *testing.T) {
	vf := NewValidationFramework()

	misaligned, err := vf.ValidateJavaAlignment()
	if err != nil {
		t.Fatalf("验证Java对齐失败: %v", err)
	}

	t.Logf("发现 %d 个与Java不对齐的常量", len(misaligned))

	// 输出不对齐的常量
	for i, mapping := range misaligned {
		if i >= 10 { // 只输出前10个
			break
		}
		t.Logf("不对齐常量 %d: Go常量=%s, Java值=%d, Go值=%d",
			i+1, mapping.GoConstant, mapping.JavaValue, mapping.GoValue)
	}
}

func TestValidationFramework_ValidateConstantOrganization(t *testing.T) {
	vf := NewValidationFramework()

	errors, err := vf.ValidateConstantOrganization()
	if err != nil {
		t.Fatalf("验证常量组织失败: %v", err)
	}

	t.Logf("发现 %d 个组织结构错误", len(errors))

	// 输出组织错误
	for i, errMsg := range errors {
		if i >= 5 { // 只输出前5个
			break
		}
		t.Logf("组织错误 %d: %s", i+1, errMsg)
	}
}

func TestValidationFramework_RunComprehensiveValidation(t *testing.T) {
	vf := NewValidationFramework()

	// 使用当前目录作为根路径
	rootPath := "."
	result, err := vf.RunComprehensiveValidation(rootPath)

	if err != nil {
		t.Fatalf("综合验证失败: %v", err)
	}

	t.Logf("验证结果: 通过=%v", result.Passed)
	t.Logf("发现魔法数字: %d", len(result.MagicNumbers))
	t.Logf("不对齐常量: %d", len(result.Misaligned))
	t.Logf("错误信息: %d", len(result.Errors))
	t.Logf("警告信息: %d", len(result.Warnings))

	// 输出详细信息
	if len(result.Errors) > 0 {
		t.Log("错误详情:")
		for i, err := range result.Errors {
			if i >= 3 { // 只输出前3个
				break
			}
			t.Logf("  错误 %d: %s", i+1, err)
		}
	}

	if len(result.Warnings) > 0 {
		t.Log("警告详情:")
		for i, warning := range result.Warnings {
			if i >= 3 { // 只输出前3个
				break
			}
			t.Logf("  警告 %d: %s", i+1, warning)
		}
	}
}

func TestValidationFramework_GetGoConstants(t *testing.T) {
	vf := NewValidationFramework()

	constants, err := vf.getGoConstants()
	if err != nil {
		t.Fatalf("获取Go常量失败: %v", err)
	}

	t.Logf("发现 %d 个Go常量", len(constants))

	// 输出一些常量示例
	count := 0
	for name, def := range constants {
		if count >= 10 { // 只输出前10个
			break
		}
		t.Logf("常量 %d: %s = %d (文件: %s)",
			count+1, name, def.Value, filepath.Base(def.File))
		count++
	}
}

func TestValidationFramework_IsSuspiciousMagicNumber(t *testing.T) {
	vf := NewValidationFramework()

	testCases := []struct {
		value    int
		filePath string
		expected bool
		desc     string
	}{
		{1, "service/coupon.go", true, "优惠券状态值1应该被识别为魔法数字"},
		{2, "service/coupon.go", true, "优惠券状态值2应该被识别为魔法数字"},
		{10, "handler/activity.go", true, "活动状态值10应该被识别为魔法数字"},
		{0, "service/common.go", false, "常见数字0不应该被识别为魔法数字"},
		{1, "constants.go", false, "常量文件中的数字不应该被识别为魔法数字"},
		{200, "handler/http.go", false, "HTTP状态码200不应该被识别为魔法数字"},
		{999, "service/test.go", false, "非枚举值999不应该被识别为魔法数字"},
	}

	for _, tc := range testCases {
		result := vf.isSuspiciousMagicNumber(tc.value, tc.filePath)
		if result != tc.expected {
			t.Errorf("测试失败: %s - 期望 %v, 实际 %v", tc.desc, tc.expected, result)
		} else {
			t.Logf("测试通过: %s", tc.desc)
		}
	}
}
