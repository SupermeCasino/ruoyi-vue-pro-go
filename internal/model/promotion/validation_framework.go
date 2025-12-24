package promotion

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

// ValidationFramework 常量验证框架
type ValidationFramework struct {
	fileSet *token.FileSet
}

// NewValidationFramework 创建新的验证框架实例
func NewValidationFramework() *ValidationFramework {
	return &ValidationFramework{
		fileSet: token.NewFileSet(),
	}
}

// NumericLiteral 表示源码中的数字字面量
type NumericLiteral struct {
	Value    int    // 数字值
	Position string // 位置信息
	Context  string // 上下文代码
	File     string // 文件路径
}

// JavaEnumMapping 表示Java枚举映射
type JavaEnumMapping struct {
	JavaEnum    string // Java枚举名称
	JavaValue   int    // Java枚举值
	GoConstant  string // Go常量名称
	GoValue     int    // Go常量值
	Description string // 描述
}

// ConstantDefinition 表示常量定义
type ConstantDefinition struct {
	Name    string // 常量名称
	Value   int    // 常量值
	File    string // 定义文件
	Package string // 包名
	Comment string // 注释
	JavaRef string // Java引用
}

// ValidationResult 验证结果
type ValidationResult struct {
	Passed       bool     // 是否通过验证
	Errors       []string // 错误信息
	Warnings     []string // 警告信息
	MagicNumbers []NumericLiteral
	Misaligned   []JavaEnumMapping
	Duplicates   []ConstantDefinition
}

// ScanSourceForNumericLiterals 扫描源码中的数字字面量
func (vf *ValidationFramework) ScanSourceForNumericLiterals(rootPath string) ([]NumericLiteral, error) {
	var literals []NumericLiteral

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理Go源文件，排除测试文件和生成的文件
		if !strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") ||
			strings.Contains(path, "wire_gen.go") ||
			strings.Contains(path, ".gen.go") {
			return nil
		}

		fileLiterals, err := vf.scanFileForNumericLiterals(path)
		if err != nil {
			return fmt.Errorf("扫描文件 %s 失败: %w", path, err)
		}

		literals = append(literals, fileLiterals...)
		return nil
	})

	return literals, err
}

// scanFileForNumericLiterals 扫描单个文件中的数字字面量
func (vf *ValidationFramework) scanFileForNumericLiterals(filePath string) ([]NumericLiteral, error) {
	var literals []NumericLiteral

	// 解析Go源文件
	src, err := parser.ParseFile(vf.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析文件失败: %w", err)
	}

	// 遍历AST查找数字字面量
	ast.Inspect(src, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.BasicLit:
			if node.Kind == token.INT {
				value, err := strconv.Atoi(node.Value)
				if err == nil && vf.isSuspiciousMagicNumber(value, filePath) {
					pos := vf.fileSet.Position(node.Pos())
					literals = append(literals, NumericLiteral{
						Value:    value,
						Position: fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column),
						Context:  vf.getContext(src, node),
						File:     filePath,
					})
				}
			}
		}
		return true
	})

	return literals, nil
}

// isSuspiciousMagicNumber 判断是否为可疑的魔法数字
func (vf *ValidationFramework) isSuspiciousMagicNumber(value int, filePath string) bool {
	// 如果在常量文件中，不算魔法数字
	if strings.Contains(filePath, "constants.go") {
		return false
	}

	// 检查是否为promotion模块相关的枚举值 - 这些应该被识别为魔法数字
	promotionEnumValues := map[int]bool{
		// 优惠券状态
		1: true, 2: true, 3: true,
		// 活动状态
		10: true, 20: true, 30: true, 40: true,
		// 营销类型
		6: true, 7: true, 8: true,
	}

	// 只有promotion相关的枚举值才被认为是可疑的魔法数字
	return promotionEnumValues[value]
}

// getContext 获取代码上下文
func (vf *ValidationFramework) getContext(file *ast.File, node ast.Node) string {
	// 简化实现，返回节点周围的代码片段
	pos := vf.fileSet.Position(node.Pos())
	return fmt.Sprintf("第%d行附近", pos.Line)
}

// ValidateJavaAlignment 验证与Java枚举的对齐
func (vf *ValidationFramework) ValidateJavaAlignment() ([]JavaEnumMapping, error) {
	var misaligned []JavaEnumMapping

	// 定义Java枚举映射表
	javaEnumMappings := []JavaEnumMapping{
		// CouponStatusEnum
		{"CouponStatusEnum.UNUSED", 1, "CouponStatusUnused", 1, "优惠券未使用状态"},
		{"CouponStatusEnum.USED", 2, "CouponStatusUsed", 2, "优惠券已使用状态"},
		{"CouponStatusEnum.EXPIRE", 3, "CouponStatusExpired", 3, "优惠券已过期状态"},

		// CouponTakeTypeEnum
		{"CouponTakeTypeEnum.USER", 1, "CouponTakeTypeUser", 1, "用户直接领取"},
		{"CouponTakeTypeEnum.ADMIN", 2, "CouponTakeTypeAdmin", 2, "管理员指定发放"},
		{"CouponTakeTypeEnum.REGISTER", 3, "CouponTakeTypeRegister", 3, "注册时自动领取"},

		// CouponTemplateValidityTypeEnum
		{"CouponTemplateValidityTypeEnum.DATE", 1, "CouponValidityTypeDate", 1, "固定日期有效期"},
		{"CouponTemplateValidityTypeEnum.TERM", 2, "CouponValidityTypeTerm", 2, "领取后有效期"},

		// PromotionProductScopeEnum
		{"PromotionProductScopeEnum.ALL", 1, "ProductScopeAll", 1, "全部商品"},
		{"PromotionProductScopeEnum.SPU", 2, "ProductScopeSpu", 2, "指定商品"},
		{"PromotionProductScopeEnum.CATEGORY", 3, "ProductScopeCategory", 3, "指定品类"},

		// PromotionDiscountTypeEnum
		{"PromotionDiscountTypeEnum.PRICE", 1, "DiscountTypePrice", 1, "满减折扣"},
		{"PromotionDiscountTypeEnum.PERCENT", 2, "DiscountTypePercent", 2, "百分比折扣"},

		// PromotionActivityStatusEnum
		{"PromotionActivityStatusEnum.WAIT", 10, "ActivityStatusWait", 10, "活动未开始"},
		{"PromotionActivityStatusEnum.RUN", 20, "ActivityStatusRun", 20, "活动进行中"},
		{"PromotionActivityStatusEnum.END", 30, "ActivityStatusEnd", 30, "活动已结束"},
		{"PromotionActivityStatusEnum.CLOSE", 40, "ActivityStatusClose", 40, "活动已关闭"},

		// PromotionTypeEnum
		{"PromotionTypeEnum.SECKILL_ACTIVITY", 1, "PromotionTypeSeckillActivity", 1, "秒杀活动"},
		{"PromotionTypeEnum.BARGAIN_ACTIVITY", 2, "PromotionTypeBargainActivity", 2, "砍价活动"},
		{"PromotionTypeEnum.COMBINATION_ACTIVITY", 3, "PromotionTypeCombinationActivity", 3, "拼团活动"},
		{"PromotionTypeEnum.DISCOUNT_ACTIVITY", 4, "PromotionTypeDiscountActivity", 4, "限时折扣"},
		{"PromotionTypeEnum.REWARD_ACTIVITY", 5, "PromotionTypeRewardActivity", 5, "满减送"},
		{"PromotionTypeEnum.MEMBER_LEVEL", 6, "PromotionTypeMemberLevel", 6, "会员折扣"},
		{"PromotionTypeEnum.COUPON", 7, "PromotionTypeCoupon", 7, "优惠券"},
		{"PromotionTypeEnum.POINT", 8, "PromotionTypePoint", 8, "积分"},

		// BannerPositionEnum
		{"BannerPositionEnum.HOME_POSITION", 1, "BannerPositionHome", 1, "首页Banner"},
		{"BannerPositionEnum.SECKILL_POSITION", 2, "BannerPositionSeckill", 2, "秒杀活动页Banner"},
		{"BannerPositionEnum.COMBINATION_POSITION", 3, "BannerPositionCombination", 3, "砍价活动页Banner"},
		{"BannerPositionEnum.DISCOUNT_POSITION", 4, "BannerPositionDiscount", 4, "限时折扣页Banner"},
		{"BannerPositionEnum.REWARD_POSITION", 5, "BannerPositionReward", 5, "满减送页Banner"},
	}

	// 获取当前Go常量定义
	goConstants, err := vf.getGoConstants()
	if err != nil {
		return nil, fmt.Errorf("获取Go常量失败: %w", err)
	}

	// 检查对齐情况
	for _, mapping := range javaEnumMappings {
		if goConst, exists := goConstants[mapping.GoConstant]; exists {
			if goConst.Value != mapping.JavaValue {
				mapping.GoValue = goConst.Value
				misaligned = append(misaligned, mapping)
			}
		} else {
			// Go常量不存在
			mapping.GoValue = -1
			misaligned = append(misaligned, mapping)
		}
	}

	return misaligned, nil
}

// getGoConstants 获取Go常量定义
func (vf *ValidationFramework) getGoConstants() (map[string]ConstantDefinition, error) {
	constants := make(map[string]ConstantDefinition)

	// 扫描promotion包中的常量文件
	constantFiles := []string{
		"internal/model/promotion/coupon_constants.go",
		"internal/model/promotion/activity_constants.go",
		"internal/model/promotion/common_constants.go",
		"internal/model/promotion/banner_constants.go",
		"internal/model/promotion/promotion_type.go",
	}

	for _, filePath := range constantFiles {
		fileConstants, err := vf.parseConstantsFromFile(filePath)
		if err != nil {
			continue // 文件可能不存在，跳过
		}

		for name, def := range fileConstants {
			constants[name] = def
		}
	}

	return constants, nil
}

// parseConstantsFromFile 从文件中解析常量定义
func (vf *ValidationFramework) parseConstantsFromFile(filePath string) (map[string]ConstantDefinition, error) {
	constants := make(map[string]ConstantDefinition)

	src, err := parser.ParseFile(vf.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// 遍历AST查找常量定义
	for _, decl := range src.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						if i < len(valueSpec.Values) {
							if basicLit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
								if value, err := strconv.Atoi(basicLit.Value); err == nil {
									constants[name.Name] = ConstantDefinition{
										Name:    name.Name,
										Value:   value,
										File:    filePath,
										Package: src.Name.Name,
										Comment: vf.getConstantComment(genDecl, valueSpec),
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return constants, nil
}

// getConstantComment 获取常量注释
func (vf *ValidationFramework) getConstantComment(genDecl *ast.GenDecl, valueSpec *ast.ValueSpec) string {
	if genDecl.Doc != nil {
		return genDecl.Doc.Text()
	}
	if valueSpec.Doc != nil {
		return valueSpec.Doc.Text()
	}
	if valueSpec.Comment != nil {
		return valueSpec.Comment.Text()
	}
	return ""
}

// ValidateConstantOrganization 验证常量组织结构
func (vf *ValidationFramework) ValidateConstantOrganization() ([]string, error) {
	var errors []string

	// 检查常量是否在正确的文件中定义
	organizationRules := map[string][]string{
		"coupon_constants.go": {
			"CouponStatus", "CouponTakeType", "CouponValidityType",
		},
		"activity_constants.go": {
			"ActivityStatus", "SeckillActivityStatus", "CombinationRecordStatus", "BargainRecordStatus",
		},
		"common_constants.go": {
			"ProductScope", "DiscountType", "ConditionType", "SenderType",
		},
		"banner_constants.go": {
			"BannerPosition", "BannerPriority", "BannerType",
		},
		"promotion_type.go": {
			"PromotionType",
		},
	}

	// 获取所有常量定义
	allConstants, err := vf.getGoConstants()
	if err != nil {
		return nil, fmt.Errorf("获取常量定义失败: %w", err)
	}

	// 检查每个常量是否在正确的文件中
	for constName, constDef := range allConstants {
		expectedFile := vf.getExpectedFileForConstant(constName, organizationRules)
		actualFile := filepath.Base(constDef.File)

		if expectedFile != "" && expectedFile != actualFile {
			errors = append(errors, fmt.Sprintf(
				"常量 %s 应该定义在 %s 中，但实际在 %s 中",
				constName, expectedFile, actualFile,
			))
		}
	}

	return errors, nil
}

// getExpectedFileForConstant 获取常量应该定义的文件
func (vf *ValidationFramework) getExpectedFileForConstant(constName string, rules map[string][]string) string {
	for fileName, prefixes := range rules {
		for _, prefix := range prefixes {
			if strings.HasPrefix(constName, prefix) {
				return fileName
			}
		}
	}
	return ""
}

// RunComprehensiveValidation 运行综合验证
func (vf *ValidationFramework) RunComprehensiveValidation(rootPath string) (*ValidationResult, error) {
	result := &ValidationResult{
		Passed:   true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 1. 扫描魔法数字
	magicNumbers, err := vf.ScanSourceForNumericLiterals(rootPath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("扫描魔法数字失败: %v", err))
		result.Passed = false
	} else {
		result.MagicNumbers = magicNumbers
		if len(magicNumbers) > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("发现 %d 个可疑的魔法数字", len(magicNumbers)))
		}
	}

	// 2. 验证Java对齐
	misaligned, err := vf.ValidateJavaAlignment()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("验证Java对齐失败: %v", err))
		result.Passed = false
	} else {
		result.Misaligned = misaligned
		if len(misaligned) > 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("发现 %d 个与Java不对齐的常量", len(misaligned)))
			result.Passed = false
		}
	}

	// 3. 验证常量组织
	orgErrors, err := vf.ValidateConstantOrganization()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("验证常量组织失败: %v", err))
		result.Passed = false
	} else {
		result.Errors = append(result.Errors, orgErrors...)
		if len(orgErrors) > 0 {
			result.Passed = false
		}
	}

	return result, nil
}
