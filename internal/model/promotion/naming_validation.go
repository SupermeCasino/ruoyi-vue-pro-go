package promotion

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// NamingValidator 命名约定验证器
type NamingValidator struct {
	fileSet *token.FileSet
}

// NewNamingValidator 创建新的命名验证器
func NewNamingValidator() *NamingValidator {
	return &NamingValidator{
		fileSet: token.NewFileSet(),
	}
}

// NamingViolation 命名违规
type NamingViolation struct {
	Type        string // 违规类型
	Name        string // 名称
	File        string // 文件路径
	Line        int    // 行号
	Expected    string // 期望的命名
	Actual      string // 实际的命名
	Description string // 描述
}

// DuplicateConstant 重复常量
type DuplicateConstant struct {
	Name        string   // 常量名称
	Value       int      // 常量值
	Files       []string // 定义文件列表
	Occurrences int      // 出现次数
}

// NamingConventionRules 命名约定规则
type NamingConventionRules struct {
	// 常量命名规则
	ConstantPatterns map[string]*regexp.Regexp
	// 函数命名规则
	FunctionPatterns map[string]*regexp.Regexp
	// 变量命名规则
	VariablePatterns map[string]*regexp.Regexp
}

// GetDefaultNamingRules 获取默认命名规则
func GetDefaultNamingRules() *NamingConventionRules {
	return &NamingConventionRules{
		ConstantPatterns: map[string]*regexp.Regexp{
			// 优惠券相关常量
			"coupon_status":        regexp.MustCompile(`^CouponStatus[A-Z][a-zA-Z]*$`),
			"coupon_take_type":     regexp.MustCompile(`^CouponTakeType[A-Z][a-zA-Z]*$`),
			"coupon_validity_type": regexp.MustCompile(`^CouponValidityType[A-Z][a-zA-Z]*$`),

			// 活动相关常量
			"activity_status":         regexp.MustCompile(`^ActivityStatus[A-Z][a-zA-Z]*$`),
			"seckill_activity_status": regexp.MustCompile(`^SeckillActivityStatus[A-Z][a-zA-Z]*$`),

			// 商品范围常量
			"product_scope": regexp.MustCompile(`^ProductScope[A-Z][a-zA-Z]*$`),

			// 折扣类型常量
			"discount_type": regexp.MustCompile(`^DiscountType[A-Z][a-zA-Z]*$`),

			// Banner相关常量
			"banner_position": regexp.MustCompile(`^BannerPosition[A-Z][a-zA-Z]*$`),
			"banner_priority": regexp.MustCompile(`^BannerPriority[A-Z][a-zA-Z]*$`),
			"banner_type":     regexp.MustCompile(`^BannerType[A-Z][a-zA-Z]*$`),

			// 营销类型常量
			"promotion_type": regexp.MustCompile(`^PromotionType[A-Z][a-zA-Z]*$`),

			// 条件类型常量
			"condition_type": regexp.MustCompile(`^ConditionType[A-Z][a-zA-Z]*$`),

			// 发送者类型常量
			"sender_type": regexp.MustCompile(`^SenderType[A-Z][a-zA-Z]*$`),
		},
		FunctionPatterns: map[string]*regexp.Regexp{
			// 验证函数
			"validation_function": regexp.MustCompile(`^IsValid[A-Z][a-zA-Z]*$`),
			// 判断函数
			"predicate_function": regexp.MustCompile(`^Is[A-Z][a-zA-Z]*$`),
		},
		VariablePatterns: map[string]*regexp.Regexp{
			// 值数组变量
			"values_array": regexp.MustCompile(`^[A-Z][a-zA-Z]*Values$`),
		},
	}
}

// ValidateNamingConventions 验证命名约定
func (nv *NamingValidator) ValidateNamingConventions(rootPath string) ([]NamingViolation, error) {
	var violations []NamingViolation
	rules := GetDefaultNamingRules()

	// 扫描promotion包中的常量文件
	constantFiles := []string{
		"backend-go/internal/model/promotion/coupon_constants.go",
		"backend-go/internal/model/promotion/activity_constants.go",
		"backend-go/internal/model/promotion/common_constants.go",
		"backend-go/internal/model/promotion/banner_constants.go",
		"backend-go/internal/model/promotion/promotion_type.go",
	}

	for _, filePath := range constantFiles {
		fileViolations, err := nv.validateFileNaming(filePath, rules)
		if err != nil {
			continue // 文件可能不存在，跳过
		}
		violations = append(violations, fileViolations...)
	}

	return violations, nil
}

// validateFileNaming 验证单个文件的命名
func (nv *NamingValidator) validateFileNaming(filePath string, rules *NamingConventionRules) ([]NamingViolation, error) {
	var violations []NamingViolation

	src, err := parser.ParseFile(nv.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// 验证常量命名
	for _, decl := range src.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						violation := nv.validateConstantName(name.Name, filePath, nv.fileSet.Position(name.Pos()).Line, rules)
						if violation != nil {
							violations = append(violations, *violation)
						}
					}
				}
			}
		}

		// 验证函数命名
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			violation := nv.validateFunctionName(funcDecl.Name.Name, filePath, nv.fileSet.Position(funcDecl.Pos()).Line, rules)
			if violation != nil {
				violations = append(violations, *violation)
			}
		}
	}

	// 验证变量命名
	for _, decl := range src.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						violation := nv.validateVariableName(name.Name, filePath, nv.fileSet.Position(name.Pos()).Line, rules)
						if violation != nil {
							violations = append(violations, *violation)
						}
					}
				}
			}
		}
	}

	return violations, nil
}

// validateConstantName 验证常量名称
func (nv *NamingValidator) validateConstantName(name, file string, line int, rules *NamingConventionRules) *NamingViolation {
	// 确定常量类型
	constType := nv.getConstantType(name)
	if constType == "" {
		return nil // 未知类型，跳过验证
	}

	// 获取对应的命名规则
	pattern, exists := rules.ConstantPatterns[constType]
	if !exists {
		return nil // 没有对应规则，跳过验证
	}

	// 检查是否符合命名规则
	if !pattern.MatchString(name) {
		expected := nv.generateExpectedName(name, constType)
		return &NamingViolation{
			Type:        "constant_naming",
			Name:        name,
			File:        file,
			Line:        line,
			Expected:    expected,
			Actual:      name,
			Description: fmt.Sprintf("常量 %s 不符合 %s 类型的命名约定", name, constType),
		}
	}

	return nil
}

// validateFunctionName 验证函数名称
func (nv *NamingValidator) validateFunctionName(name, file string, line int, rules *NamingConventionRules) *NamingViolation {
	// 确定函数类型
	funcType := nv.getFunctionType(name)
	if funcType == "" {
		return nil // 未知类型，跳过验证
	}

	// 获取对应的命名规则
	pattern, exists := rules.FunctionPatterns[funcType]
	if !exists {
		return nil // 没有对应规则，跳过验证
	}

	// 检查是否符合命名规则
	if !pattern.MatchString(name) {
		expected := nv.generateExpectedFunctionName(name, funcType)
		return &NamingViolation{
			Type:        "function_naming",
			Name:        name,
			File:        file,
			Line:        line,
			Expected:    expected,
			Actual:      name,
			Description: fmt.Sprintf("函数 %s 不符合 %s 类型的命名约定", name, funcType),
		}
	}

	return nil
}

// validateVariableName 验证变量名称
func (nv *NamingValidator) validateVariableName(name, file string, line int, rules *NamingConventionRules) *NamingViolation {
	// 确定变量类型
	varType := nv.getVariableType(name)
	if varType == "" {
		return nil // 未知类型，跳过验证
	}

	// 获取对应的命名规则
	pattern, exists := rules.VariablePatterns[varType]
	if !exists {
		return nil // 没有对应规则，跳过验证
	}

	// 检查是否符合命名规则
	if !pattern.MatchString(name) {
		expected := nv.generateExpectedVariableName(name, varType)
		return &NamingViolation{
			Type:        "variable_naming",
			Name:        name,
			File:        file,
			Line:        line,
			Expected:    expected,
			Actual:      name,
			Description: fmt.Sprintf("变量 %s 不符合 %s 类型的命名约定", name, varType),
		}
	}

	return nil
}

// getConstantType 获取常量类型
func (nv *NamingValidator) getConstantType(name string) string {
	switch {
	case strings.HasPrefix(name, "CouponStatus"):
		return "coupon_status"
	case strings.HasPrefix(name, "CouponTakeType"):
		return "coupon_take_type"
	case strings.HasPrefix(name, "CouponValidityType"):
		return "coupon_validity_type"
	case strings.HasPrefix(name, "ActivityStatus"):
		return "activity_status"
	case strings.HasPrefix(name, "SeckillActivityStatus"):
		return "seckill_activity_status"
	case strings.HasPrefix(name, "ProductScope"):
		return "product_scope"
	case strings.HasPrefix(name, "DiscountType"):
		return "discount_type"
	case strings.HasPrefix(name, "BannerPosition"):
		return "banner_position"
	case strings.HasPrefix(name, "BannerPriority"):
		return "banner_priority"
	case strings.HasPrefix(name, "BannerType"):
		return "banner_type"
	case strings.HasPrefix(name, "PromotionType"):
		return "promotion_type"
	case strings.HasPrefix(name, "ConditionType"):
		return "condition_type"
	case strings.HasPrefix(name, "SenderType"):
		return "sender_type"
	default:
		return ""
	}
}

// getFunctionType 获取函数类型
func (nv *NamingValidator) getFunctionType(name string) string {
	switch {
	case strings.HasPrefix(name, "IsValid"):
		return "validation_function"
	case strings.HasPrefix(name, "Is") && !strings.HasPrefix(name, "IsValid"):
		return "predicate_function"
	default:
		return ""
	}
}

// getVariableType 获取变量类型
func (nv *NamingValidator) getVariableType(name string) string {
	if strings.HasSuffix(name, "Values") {
		return "values_array"
	}
	return ""
}

// generateExpectedName 生成期望的常量名称
func (nv *NamingValidator) generateExpectedName(name, constType string) string {
	// 简化实现，返回基于类型的建议名称
	switch constType {
	case "coupon_status":
		return "CouponStatus" + nv.toPascalCase(strings.TrimPrefix(name, "CouponStatus"))
	case "coupon_take_type":
		return "CouponTakeType" + nv.toPascalCase(strings.TrimPrefix(name, "CouponTakeType"))
	default:
		return name
	}
}

// generateExpectedFunctionName 生成期望的函数名称
func (nv *NamingValidator) generateExpectedFunctionName(name, funcType string) string {
	switch funcType {
	case "validation_function":
		if !strings.HasPrefix(name, "IsValid") {
			return "IsValid" + nv.toPascalCase(name)
		}
	case "predicate_function":
		if !strings.HasPrefix(name, "Is") {
			return "Is" + nv.toPascalCase(name)
		}
	}
	return name
}

// generateExpectedVariableName 生成期望的变量名称
func (nv *NamingValidator) generateExpectedVariableName(name, varType string) string {
	if varType == "values_array" && !strings.HasSuffix(name, "Values") {
		return name + "Values"
	}
	return name
}

// toPascalCase 转换为PascalCase
func (nv *NamingValidator) toPascalCase(s string) string {
	if s == "" {
		return s
	}

	// 将首字母大写
	runes := []rune(s)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes)
}

// DetectDuplicateConstants 检测重复常量
func (nv *NamingValidator) DetectDuplicateConstants(rootPath string) ([]DuplicateConstant, error) {
	var duplicates []DuplicateConstant

	// 收集所有常量定义
	constantMap := make(map[string]map[int][]string) // name -> value -> files

	// 扫描promotion包中的常量文件
	constantFiles := []string{
		"backend-go/internal/model/promotion/coupon_constants.go",
		"backend-go/internal/model/promotion/activity_constants.go",
		"backend-go/internal/model/promotion/common_constants.go",
		"backend-go/internal/model/promotion/banner_constants.go",
		"backend-go/internal/model/promotion/promotion_type.go",
	}

	for _, filePath := range constantFiles {
		constants, err := nv.parseConstantsFromFile(filePath)
		if err != nil {
			continue // 文件可能不存在，跳过
		}

		for name, value := range constants {
			if constantMap[name] == nil {
				constantMap[name] = make(map[int][]string)
			}
			constantMap[name][value] = append(constantMap[name][value], filePath)
		}
	}

	// 检测重复定义
	for name, valueMap := range constantMap {
		totalOccurrences := 0
		var allFiles []string
		var duplicateValue int

		for value, files := range valueMap {
			totalOccurrences += len(files)
			allFiles = append(allFiles, files...)
			if len(files) > 1 {
				duplicateValue = value
			}
		}

		// 如果同一个常量名在多个文件中定义，或者同一个值被多个常量使用
		if totalOccurrences > 1 || len(valueMap) > 1 {
			duplicates = append(duplicates, DuplicateConstant{
				Name:        name,
				Value:       duplicateValue,
				Files:       allFiles,
				Occurrences: totalOccurrences,
			})
		}
	}

	// 检测相同值的不同常量（可能的重复定义）
	valueToConstants := make(map[int][]string)
	for name, valueMap := range constantMap {
		for value := range valueMap {
			valueToConstants[value] = append(valueToConstants[value], name)
		}
	}

	for value, constants := range valueToConstants {
		if len(constants) > 1 {
			// 检查是否为合理的重复（例如，不同模块的相同状态值）
			if !nv.isReasonableDuplication(constants, value) {
				duplicates = append(duplicates, DuplicateConstant{
					Name:        strings.Join(constants, ", "),
					Value:       value,
					Files:       []string{"多个文件"},
					Occurrences: len(constants),
				})
			}
		}
	}

	return duplicates, nil
}

// parseConstantsFromFile 从文件中解析常量
func (nv *NamingValidator) parseConstantsFromFile(filePath string) (map[string]int, error) {
	constants := make(map[string]int)

	src, err := parser.ParseFile(nv.fileSet, filePath, nil, parser.ParseComments)
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
									constants[name.Name] = value
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

// isReasonableDuplication 判断是否为合理的重复
func (nv *NamingValidator) isReasonableDuplication(constants []string, value int) bool {
	// 检查是否为不同模块的相同概念
	modules := make(map[string]bool)
	for _, constName := range constants {
		module := nv.getConstantModule(constName)
		if module != "" {
			modules[module] = true
		}
	}

	// 如果来自同一模块，不合理
	if len(modules) <= 1 {
		return false
	}

	// 一些值在不同模块中重复是合理的
	reasonableValues := map[int]bool{
		0: true, 1: true, 2: true, // 常见的状态值
		10: true, 20: true, 30: true, 40: true, // 活动状态值
	}

	return reasonableValues[value]
}

// getConstantModule 获取常量所属模块
func (nv *NamingValidator) getConstantModule(constName string) string {
	switch {
	case strings.HasPrefix(constName, "Coupon"):
		return "coupon"
	case strings.HasPrefix(constName, "Activity"):
		return "activity"
	case strings.HasPrefix(constName, "Banner"):
		return "banner"
	case strings.HasPrefix(constName, "Product"):
		return "product"
	case strings.HasPrefix(constName, "Discount"):
		return "discount"
	case strings.HasPrefix(constName, "Promotion"):
		return "promotion"
	default:
		return "common"
	}
}

// RunNamingValidation 运行命名验证
func (nv *NamingValidator) RunNamingValidation(rootPath string) ([]NamingViolation, []DuplicateConstant, error) {
	// 验证命名约定
	violations, err := nv.ValidateNamingConventions(rootPath)
	if err != nil {
		return nil, nil, fmt.Errorf("验证命名约定失败: %w", err)
	}

	// 检测重复常量
	duplicates, err := nv.DetectDuplicateConstants(rootPath)
	if err != nil {
		return violations, nil, fmt.Errorf("检测重复常量失败: %w", err)
	}

	return violations, duplicates, nil
}
