package promotion

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidationTool 综合验证工具
type ValidationTool struct {
	framework       *ValidationFramework
	namingValidator *NamingValidator
	rootPath        string
}

// NewValidationTool 创建新的验证工具
func NewValidationTool(rootPath string) *ValidationTool {
	return &ValidationTool{
		framework:       NewValidationFramework(),
		namingValidator: NewNamingValidator(),
		rootPath:        rootPath,
	}
}

// ComprehensiveValidationReport 综合验证报告
type ComprehensiveValidationReport struct {
	// 基础验证结果
	ValidationResult *ValidationResult

	// 命名约定违规
	NamingViolations []NamingViolation

	// 重复常量
	DuplicateConstants []DuplicateConstant

	// 总体状态
	OverallPassed bool
	Summary       string

	// 详细统计
	Statistics ValidationStatistics
}

// ValidationStatistics 验证统计信息
type ValidationStatistics struct {
	TotalFiles          int // 扫描的文件总数
	TotalConstants      int // 常量总数
	MagicNumbersFound   int // 发现的魔法数字数量
	MisalignedConstants int // 与Java不对齐的常量数量
	NamingViolations    int // 命名违规数量
	DuplicateConstants  int // 重复常量数量
	OrganizationErrors  int // 组织结构错误数量
}

// RunComprehensiveValidation 运行综合验证
func (vt *ValidationTool) RunComprehensiveValidation() (*ComprehensiveValidationReport, error) {
	report := &ComprehensiveValidationReport{}

	// 1. 运行基础验证框架
	fmt.Println("正在运行基础验证框架...")
	validationResult, err := vt.framework.RunComprehensiveValidation(vt.rootPath)
	if err != nil {
		return nil, fmt.Errorf("基础验证失败: %w", err)
	}
	report.ValidationResult = validationResult

	// 2. 运行命名约定验证
	fmt.Println("正在验证命名约定...")
	namingViolations, duplicateConstants, err := vt.namingValidator.RunNamingValidation(vt.rootPath)
	if err != nil {
		return nil, fmt.Errorf("命名验证失败: %w", err)
	}
	report.NamingViolations = namingViolations
	report.DuplicateConstants = duplicateConstants

	// 3. 计算统计信息
	report.Statistics = vt.calculateStatistics(validationResult, namingViolations, duplicateConstants)

	// 4. 确定总体状态
	report.OverallPassed = validationResult.Passed &&
		len(namingViolations) == 0 &&
		len(duplicateConstants) == 0

	// 5. 生成总结
	report.Summary = vt.generateSummary(report)

	return report, nil
}

// calculateStatistics 计算验证统计信息
func (vt *ValidationTool) calculateStatistics(
	validationResult *ValidationResult,
	namingViolations []NamingViolation,
	duplicateConstants []DuplicateConstant,
) ValidationStatistics {

	// 统计文件数量
	totalFiles := vt.countPromotionFiles()

	// 统计常量数量
	totalConstants := vt.countTotalConstants()

	return ValidationStatistics{
		TotalFiles:          totalFiles,
		TotalConstants:      totalConstants,
		MagicNumbersFound:   len(validationResult.MagicNumbers),
		MisalignedConstants: len(validationResult.Misaligned),
		NamingViolations:    len(namingViolations),
		DuplicateConstants:  len(duplicateConstants),
		OrganizationErrors:  len(validationResult.Errors),
	}
}

// countPromotionFiles 统计promotion模块文件数量
func (vt *ValidationTool) countPromotionFiles() int {
	promotionPath := filepath.Join(vt.rootPath, "backend-go/internal/model/promotion")
	files, err := os.ReadDir(promotionPath)
	if err != nil {
		return 0
	}

	count := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") && !strings.HasSuffix(file.Name(), "_test.go") {
			count++
		}
	}
	return count
}

// countTotalConstants 统计常量总数
func (vt *ValidationTool) countTotalConstants() int {
	constants, err := vt.framework.getGoConstants()
	if err != nil {
		return 0
	}
	return len(constants)
}

// generateSummary 生成验证总结
func (vt *ValidationTool) generateSummary(report *ComprehensiveValidationReport) string {
	var summary strings.Builder

	summary.WriteString("=== 促销模块常量验证报告 ===\n\n")

	// 总体状态
	if report.OverallPassed {
		summary.WriteString("✅ 验证通过：所有检查项目均符合要求\n\n")
	} else {
		summary.WriteString("❌ 验证失败：发现需要修复的问题\n\n")
	}

	// 统计信息
	stats := report.Statistics
	summary.WriteString("📊 统计信息：\n")
	summary.WriteString(fmt.Sprintf("  - 扫描文件数量: %d\n", stats.TotalFiles))
	summary.WriteString(fmt.Sprintf("  - 常量总数: %d\n", stats.TotalConstants))
	summary.WriteString(fmt.Sprintf("  - 发现魔法数字: %d\n", stats.MagicNumbersFound))
	summary.WriteString(fmt.Sprintf("  - Java对齐问题: %d\n", stats.MisalignedConstants))
	summary.WriteString(fmt.Sprintf("  - 命名约定违规: %d\n", stats.NamingViolations))
	summary.WriteString(fmt.Sprintf("  - 重复常量: %d\n", stats.DuplicateConstants))
	summary.WriteString("\n")

	// 详细问题报告
	if !report.OverallPassed {
		summary.WriteString("🔍 详细问题：\n")

		// 魔法数字
		if len(report.ValidationResult.MagicNumbers) > 0 {
			summary.WriteString("\n📍 发现的魔法数字：\n")
			for _, magic := range report.ValidationResult.MagicNumbers {
				summary.WriteString(fmt.Sprintf("  - 值 %d 在 %s (%s)\n",
					magic.Value, magic.Position, magic.Context))
			}
		}

		// Java对齐问题
		if len(report.ValidationResult.Misaligned) > 0 {
			summary.WriteString("\n🔄 Java对齐问题：\n")
			for _, misaligned := range report.ValidationResult.Misaligned {
				summary.WriteString(fmt.Sprintf("  - %s: Java值=%d, Go值=%d\n",
					misaligned.GoConstant, misaligned.JavaValue, misaligned.GoValue))
			}
		}

		// 命名约定违规
		if len(report.NamingViolations) > 0 {
			summary.WriteString("\n📝 命名约定违规：\n")
			for _, violation := range report.NamingViolations {
				summary.WriteString(fmt.Sprintf("  - %s (第%d行): %s\n",
					violation.Name, violation.Line, violation.Description))
			}
		}

		// 重复常量
		if len(report.DuplicateConstants) > 0 {
			summary.WriteString("\n🔁 重复常量：\n")
			for _, duplicate := range report.DuplicateConstants {
				summary.WriteString(fmt.Sprintf("  - %s (值=%d): 出现%d次\n",
					duplicate.Name, duplicate.Value, duplicate.Occurrences))
			}
		}

		// 组织结构错误
		if len(report.ValidationResult.Errors) > 0 {
			summary.WriteString("\n📁 组织结构问题：\n")
			for _, err := range report.ValidationResult.Errors {
				summary.WriteString(fmt.Sprintf("  - %s\n", err))
			}
		}
	}

	// 建议
	summary.WriteString("\n💡 建议：\n")
	if report.OverallPassed {
		summary.WriteString("  - 常量定义规范，继续保持良好的编码习惯\n")
		summary.WriteString("  - 定期运行验证工具确保代码质量\n")
	} else {
		if stats.MagicNumbersFound > 0 {
			summary.WriteString("  - 将发现的魔法数字替换为有意义的常量\n")
		}
		if stats.MisalignedConstants > 0 {
			summary.WriteString("  - 检查并修正与Java不对齐的常量值\n")
		}
		if stats.NamingViolations > 0 {
			summary.WriteString("  - 按照Go命名约定重命名违规的标识符\n")
		}
		if stats.DuplicateConstants > 0 {
			summary.WriteString("  - 消除重复的常量定义，统一使用单一定义\n")
		}
	}

	return summary.String()
}

// GenerateDetailedReport 生成详细报告
func (vt *ValidationTool) GenerateDetailedReport(report *ComprehensiveValidationReport) string {
	var detailed strings.Builder

	detailed.WriteString(report.Summary)
	detailed.WriteString("\n" + strings.Repeat("=", 60) + "\n")
	detailed.WriteString("详细验证结果\n")
	detailed.WriteString(strings.Repeat("=", 60) + "\n\n")

	// 1. 魔法数字详细信息
	if len(report.ValidationResult.MagicNumbers) > 0 {
		detailed.WriteString("1. 魔法数字详细信息：\n")
		detailed.WriteString(strings.Repeat("-", 40) + "\n")
		for i, magic := range report.ValidationResult.MagicNumbers {
			detailed.WriteString(fmt.Sprintf("%d. 值: %d\n", i+1, magic.Value))
			detailed.WriteString(fmt.Sprintf("   位置: %s\n", magic.Position))
			detailed.WriteString(fmt.Sprintf("   上下文: %s\n", magic.Context))
			detailed.WriteString(fmt.Sprintf("   文件: %s\n\n", magic.File))
		}
	}

	// 2. Java对齐详细信息
	if len(report.ValidationResult.Misaligned) > 0 {
		detailed.WriteString("2. Java对齐详细信息：\n")
		detailed.WriteString(strings.Repeat("-", 40) + "\n")
		for i, misaligned := range report.ValidationResult.Misaligned {
			detailed.WriteString(fmt.Sprintf("%d. Go常量: %s\n", i+1, misaligned.GoConstant))
			detailed.WriteString(fmt.Sprintf("   Java枚举: %s\n", misaligned.JavaEnum))
			detailed.WriteString(fmt.Sprintf("   Java值: %d\n", misaligned.JavaValue))
			detailed.WriteString(fmt.Sprintf("   Go值: %d\n", misaligned.GoValue))
			detailed.WriteString(fmt.Sprintf("   描述: %s\n\n", misaligned.Description))
		}
	}

	// 3. 命名约定违规详细信息
	if len(report.NamingViolations) > 0 {
		detailed.WriteString("3. 命名约定违规详细信息：\n")
		detailed.WriteString(strings.Repeat("-", 40) + "\n")
		for i, violation := range report.NamingViolations {
			detailed.WriteString(fmt.Sprintf("%d. 名称: %s\n", i+1, violation.Name))
			detailed.WriteString(fmt.Sprintf("   类型: %s\n", violation.Type))
			detailed.WriteString(fmt.Sprintf("   文件: %s (第%d行)\n", violation.File, violation.Line))
			detailed.WriteString(fmt.Sprintf("   期望: %s\n", violation.Expected))
			detailed.WriteString(fmt.Sprintf("   实际: %s\n", violation.Actual))
			detailed.WriteString(fmt.Sprintf("   描述: %s\n\n", violation.Description))
		}
	}

	// 4. 重复常量详细信息
	if len(report.DuplicateConstants) > 0 {
		detailed.WriteString("4. 重复常量详细信息：\n")
		detailed.WriteString(strings.Repeat("-", 40) + "\n")
		for i, duplicate := range report.DuplicateConstants {
			detailed.WriteString(fmt.Sprintf("%d. 常量名: %s\n", i+1, duplicate.Name))
			detailed.WriteString(fmt.Sprintf("   值: %d\n", duplicate.Value))
			detailed.WriteString(fmt.Sprintf("   出现次数: %d\n", duplicate.Occurrences))
			detailed.WriteString(fmt.Sprintf("   文件: %s\n\n", strings.Join(duplicate.Files, ", ")))
		}
	}

	return detailed.String()
}

// SaveReportToFile 将报告保存到文件
func (vt *ValidationTool) SaveReportToFile(report *ComprehensiveValidationReport, filename string) error {
	content := vt.GenerateDetailedReport(report)

	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// ValidateSpecificFile 验证特定文件
func (vt *ValidationTool) ValidateSpecificFile(filePath string) (*ComprehensiveValidationReport, error) {
	// 创建临时验证工具，只验证指定文件
	tempTool := &ValidationTool{
		framework:       NewValidationFramework(),
		namingValidator: NewNamingValidator(),
		rootPath:        filepath.Dir(filePath),
	}

	// 只验证指定文件
	report := &ComprehensiveValidationReport{}

	// 扫描魔法数字
	magicNumbers, err := tempTool.framework.scanFileForNumericLiterals(filePath)
	if err != nil {
		return nil, fmt.Errorf("扫描文件魔法数字失败: %w", err)
	}

	// 验证命名约定
	rules := GetDefaultNamingRules()
	namingViolations, err := tempTool.namingValidator.validateFileNaming(filePath, rules)
	if err != nil {
		return nil, fmt.Errorf("验证文件命名失败: %w", err)
	}

	// 构建报告
	report.ValidationResult = &ValidationResult{
		MagicNumbers: magicNumbers,
		Passed:       len(magicNumbers) == 0,
	}
	report.NamingViolations = namingViolations
	report.OverallPassed = len(magicNumbers) == 0 && len(namingViolations) == 0
	report.Summary = tempTool.generateSummary(report)

	return report, nil
}
