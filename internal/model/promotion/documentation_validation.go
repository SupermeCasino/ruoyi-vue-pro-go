package promotion

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DocumentationValidation 文档验证结果
type DocumentationValidation struct {
	ConstantName     string
	File             string
	HasDocumentation bool
	Documentation    []string
	HasJavaReference bool
	MissingElements  []string
}

// ValidateDocumentationCompleteness 验证所有常量的文档完整性
func ValidateDocumentationCompleteness() ([]DocumentationValidation, error) {
	var validations []DocumentationValidation

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	promotionDir := filepath.Join(wd, ".")

	// 遍历所有Go文件
	err = filepath.Walk(promotionDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理常量文件，跳过测试文件
		if !strings.HasSuffix(path, "_constants.go") && !strings.HasSuffix(path, "promotion_type.go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 解析Go文件
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// 提取常量及其文档
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.GenDecl:
				if x.Tok == token.CONST {
					for _, spec := range x.Specs {
						if valueSpec, ok := spec.(*ast.ValueSpec); ok {
							for _, name := range valueSpec.Names {
								// 收集文档注释
								var comments []string
								if x.Doc != nil {
									for _, comment := range x.Doc.List {
										comments = append(comments, comment.Text)
									}
								}

								// 检查行内注释
								if valueSpec.Comment != nil {
									for _, comment := range valueSpec.Comment.List {
										comments = append(comments, comment.Text)
									}
								}

								validation := DocumentationValidation{
									ConstantName:     name.Name,
									File:             filepath.Base(path),
									Documentation:    comments,
									HasDocumentation: len(comments) > 0,
								}

								// 检查是否有Java参考
								for _, comment := range comments {
									if strings.Contains(comment, "Java:") || strings.Contains(comment, "对齐") {
										validation.HasJavaReference = true
										break
									}
								}

								// 检查缺失的文档元素
								validation.MissingElements = checkMissingDocumentationElements(validation)

								validations = append(validations, validation)
							}
						}
					}
				}
			}
			return true
		})

		return nil
	})

	return validations, err
}

// checkMissingDocumentationElements 检查缺失的文档元素
func checkMissingDocumentationElements(validation DocumentationValidation) []string {
	var missing []string

	if !validation.HasDocumentation {
		missing = append(missing, "缺少注释")
		return missing
	}

	// 将所有注释合并为一个字符串进行检查
	allComments := strings.Join(validation.Documentation, " ")

	// 检查是否有中文描述
	hasChineseDescription := false
	for _, comment := range validation.Documentation {
		// 检查是否包含中文字符
		for _, r := range comment {
			if r >= 0x4e00 && r <= 0x9fff {
				hasChineseDescription = true
				break
			}
		}
		if hasChineseDescription {
			break
		}
	}

	if !hasChineseDescription {
		missing = append(missing, "缺少中文描述")
	}

	// 检查是否有Java参考（对于促销模块常量）
	if isPromotionConstantForDoc(validation.ConstantName) && !validation.HasJavaReference {
		missing = append(missing, "缺少Java参考")
	}

	// 检查特定常量的特殊要求
	if strings.HasPrefix(validation.ConstantName, "Coupon") ||
		strings.HasPrefix(validation.ConstantName, "Activity") ||
		strings.HasPrefix(validation.ConstantName, "Product") ||
		strings.HasPrefix(validation.ConstantName, "Discount") ||
		strings.HasPrefix(validation.ConstantName, "Banner") {

		// 这些常量应该有详细的用途说明
		if len(allComments) < 20 { // 简单的长度检查
			missing = append(missing, "文档过于简短")
		}
	}

	return missing
}

// isPromotionConstantForDoc 检查是否为促销相关常量（文档验证专用）
func isPromotionConstantForDoc(name string) bool {
	promotionPrefixes := []string{
		"CouponStatus", "CouponTakeType", "CouponValidityType",
		"ActivityStatus", "PromotionType", "SeckillActivityStatus",
		"CombinationRecordStatus", "BargainRecordStatus",
		"ProductScope", "DiscountType", "ConditionType",
		"BannerPosition", "BannerStatus", "BannerPriority", "BannerType",
		"DefaultPageSize", "MaxPageSize", "MinPrice", "MaxPrice",
		"MinDiscountPercent", "MaxDiscountPercent", "CouponTemplateTakeLimitCountMax",
	}

	for _, prefix := range promotionPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

// PrintDocumentationReport 打印文档完整性报告
func PrintDocumentationReport() error {
	validations, err := ValidateDocumentationCompleteness()
	if err != nil {
		return err
	}

	fmt.Println("=== 文档完整性验证报告 ===")
	fmt.Println()

	completeCount := 0
	incompleteCount := 0

	for _, validation := range validations {
		if len(validation.MissingElements) == 0 {
			completeCount++
			fmt.Printf("✅ %s (%s) - 文档完整\n", validation.ConstantName, validation.File)
		} else {
			incompleteCount++
			fmt.Printf("❌ %s (%s) - 文档不完整: %s\n",
				validation.ConstantName, validation.File, strings.Join(validation.MissingElements, ", "))
		}
	}

	fmt.Println()
	fmt.Printf("总计: %d 个常量\n", len(validations))
	fmt.Printf("文档完整: %d 个\n", completeCount)
	fmt.Printf("文档不完整: %d 个\n", incompleteCount)

	if incompleteCount == 0 {
		fmt.Println("🎉 所有常量都有完整的文档！")
	} else {
		fmt.Printf("⚠️  发现 %d 个常量文档不完整，需要改进\n", incompleteCount)
	}

	return nil
}

// GetIncompleteDocumentationConstants 获取所有文档不完整的常量
func GetIncompleteDocumentationConstants() ([]DocumentationValidation, error) {
	validations, err := ValidateDocumentationCompleteness()
	if err != nil {
		return nil, err
	}

	var incomplete []DocumentationValidation
	for _, validation := range validations {
		if len(validation.MissingElements) > 0 {
			incomplete = append(incomplete, validation)
		}
	}

	return incomplete, nil
}

// IsAllDocumentationComplete 检查是否所有常量都有完整的文档
func IsAllDocumentationComplete() (bool, error) {
	incomplete, err := GetIncompleteDocumentationConstants()
	if err != nil {
		return false, err
	}
	return len(incomplete) == 0, nil
}
