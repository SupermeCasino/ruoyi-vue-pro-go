package promotion

import (
	"testing"
)

// TestDocumentationCompleteness 测试文档完整性
// Feature: promotion-constants-refactor, Property 4: For any constant definition, it should include comprehensive documentation comments
func TestDocumentationCompleteness(t *testing.T) {
	validations, err := ValidateDocumentationCompleteness()
	if err != nil {
		t.Fatalf("验证文档完整性失败: %v", err)
	}

	for _, validation := range validations {
		t.Run(validation.ConstantName, func(t *testing.T) {
			if len(validation.MissingElements) > 0 {
				t.Errorf("常量 %s 文档不完整: %v", validation.ConstantName, validation.MissingElements)
			}
		})
	}
}

// TestAllDocumentationComplete 测试所有常量是否都有完整的文档
func TestAllDocumentationComplete(t *testing.T) {
	isComplete, err := IsAllDocumentationComplete()
	if err != nil {
		t.Fatalf("检查文档完整性失败: %v", err)
	}

	if !isComplete {
		incomplete, err := GetIncompleteDocumentationConstants()
		if err != nil {
			t.Fatalf("获取不完整文档常量失败: %v", err)
		}

		t.Errorf("发现 %d 个常量文档不完整:", len(incomplete))
		for _, validation := range incomplete {
			t.Errorf("  - %s (%s): %v",
				validation.ConstantName, validation.File, validation.MissingElements)
		}
	}
}

// TestSpecificDocumentationRequirements 测试特定的文档需求
func TestSpecificDocumentationRequirements(t *testing.T) {
	validations, err := ValidateDocumentationCompleteness()
	if err != nil {
		t.Fatalf("验证文档完整性失败: %v", err)
	}

	// 检查关键常量是否有Java参考
	keyConstants := []string{
		"CouponStatusUnused", "CouponStatusUsed", "CouponStatusExpired",
		"CouponTakeTypeUser", "CouponTakeTypeAdmin", "CouponTakeTypeRegister",
		"ActivityStatusWait", "ActivityStatusRun", "ActivityStatusEnd", "ActivityStatusClose",
		"ProductScopeAll", "ProductScopeSpu", "ProductScopeCategory",
		"DiscountTypePrice", "DiscountTypePercent",
	}

	for _, constantName := range keyConstants {
		t.Run("JavaReference_"+constantName, func(t *testing.T) {
			found := false
			for _, validation := range validations {
				if validation.ConstantName == constantName {
					found = true
					if !validation.HasJavaReference {
						t.Errorf("关键常量 %s 缺少Java参考文档", constantName)
					}
					break
				}
			}
			if !found {
				t.Errorf("未找到常量 %s", constantName)
			}
		})
	}
}

// TestDocumentationReport 测试文档报告生成
func TestDocumentationReport(t *testing.T) {
	// 这个测试主要用于手动验证，打印文档报告
	if testing.Verbose() {
		err := PrintDocumentationReport()
		if err != nil {
			t.Errorf("生成文档报告失败: %v", err)
		}
	}

	// 验证报告生成不会出错
	validations, err := ValidateDocumentationCompleteness()
	if err != nil {
		t.Errorf("验证文档完整性失败: %v", err)
	}

	if len(validations) == 0 {
		t.Error("验证列表不应该为空")
	}
}

// TestPromotionConstantDocumentation 测试促销常量的特定文档要求
func TestPromotionConstantDocumentation(t *testing.T) {
	validations, err := ValidateDocumentationCompleteness()
	if err != nil {
		t.Fatalf("验证文档完整性失败: %v", err)
	}

	// 验证所有促销常量都有基本文档
	for _, validation := range validations {
		if isPromotionConstantForDoc(validation.ConstantName) {
			t.Run("BasicDoc_"+validation.ConstantName, func(t *testing.T) {
				if !validation.HasDocumentation {
					t.Errorf("促销常量 %s 缺少基本文档", validation.ConstantName)
				}
			})
		}
	}
}
