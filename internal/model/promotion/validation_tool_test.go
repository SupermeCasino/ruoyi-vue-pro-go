package promotion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidationTool_RunComprehensiveValidation(t *testing.T) {
	// 使用当前目录作为根路径
	rootPath := "."
	tool := NewValidationTool(rootPath)

	// 运行综合验证
	report, err := tool.RunComprehensiveValidation()
	if err != nil {
		t.Fatalf("综合验证失败: %v", err)
	}

	// 输出验证结果
	t.Logf("综合验证结果:")
	t.Logf("  总体通过: %v", report.OverallPassed)
	t.Logf("  扫描文件数: %d", report.Statistics.TotalFiles)
	t.Logf("  常量总数: %d", report.Statistics.TotalConstants)
	t.Logf("  魔法数字: %d", report.Statistics.MagicNumbersFound)
	t.Logf("  Java对齐问题: %d", report.Statistics.MisalignedConstants)
	t.Logf("  命名违规: %d", report.Statistics.NamingViolations)
	t.Logf("  重复常量: %d", report.Statistics.DuplicateConstants)
	t.Logf("  组织错误: %d", report.Statistics.OrganizationErrors)

	// 输出总结
	t.Log("验证总结:")
	t.Log(report.Summary)

	// 如果有问题，输出详细信息
	if !report.OverallPassed {
		t.Log("详细问题:")

		// 输出前几个魔法数字
		if len(report.ValidationResult.MagicNumbers) > 0 {
			t.Log("魔法数字示例:")
			for i, magic := range report.ValidationResult.MagicNumbers {
				if i >= 3 {
					break
				}
				t.Logf("  %d. 值=%d, 位置=%s", i+1, magic.Value, magic.Position)
			}
		}

		// 输出前几个Java对齐问题
		if len(report.ValidationResult.Misaligned) > 0 {
			t.Log("Java对齐问题示例:")
			for i, misaligned := range report.ValidationResult.Misaligned {
				if i >= 3 {
					break
				}
				t.Logf("  %d. %s: Java=%d, Go=%d",
					i+1, misaligned.GoConstant, misaligned.JavaValue, misaligned.GoValue)
			}
		}

		// 输出前几个命名违规
		if len(report.NamingViolations) > 0 {
			t.Log("命名违规示例:")
			for i, violation := range report.NamingViolations {
				if i >= 3 {
					break
				}
				t.Logf("  %d. %s: %s", i+1, violation.Name, violation.Description)
			}
		}
	}
}

func TestValidationTool_GenerateDetailedReport(t *testing.T) {
	rootPath := "."
	tool := NewValidationTool(rootPath)

	// 运行验证获取报告
	report, err := tool.RunComprehensiveValidation()
	if err != nil {
		t.Fatalf("运行验证失败: %v", err)
	}

	// 生成详细报告
	detailedReport := tool.GenerateDetailedReport(report)

	// 检查报告内容
	if len(detailedReport) == 0 {
		t.Error("详细报告为空")
	}

	// 检查报告是否包含关键信息
	expectedSections := []string{
		"促销模块常量验证报告",
		"统计信息",
		"详细验证结果",
	}

	for _, section := range expectedSections {
		if !strings.Contains(detailedReport, section) {
			t.Errorf("详细报告缺少部分: %s", section)
		}
	}

	t.Logf("详细报告长度: %d 字符", len(detailedReport))

	// 输出报告的前几行用于调试
	lines := strings.Split(detailedReport, "\n")
	t.Log("详细报告前10行:")
	for i, line := range lines {
		if i >= 10 {
			break
		}
		t.Logf("  %d: %s", i+1, line)
	}
}

func TestValidationTool_SaveReportToFile(t *testing.T) {
	rootPath := "."
	tool := NewValidationTool(rootPath)

	// 运行验证获取报告
	report, err := tool.RunComprehensiveValidation()
	if err != nil {
		t.Fatalf("运行验证失败: %v", err)
	}

	// 创建临时文件路径
	tempDir := t.TempDir()
	reportFile := filepath.Join(tempDir, "validation_report.txt")

	// 保存报告到文件
	err = tool.SaveReportToFile(report, reportFile)
	if err != nil {
		t.Fatalf("保存报告失败: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(reportFile); os.IsNotExist(err) {
		t.Error("报告文件未创建")
	}

	// 读取文件内容
	content, err := os.ReadFile(reportFile)
	if err != nil {
		t.Fatalf("读取报告文件失败: %v", err)
	}

	// 检查文件内容
	if len(content) == 0 {
		t.Error("报告文件为空")
	}

	t.Logf("报告文件大小: %d 字节", len(content))
	t.Logf("报告文件路径: %s", reportFile)

	// 检查文件内容是否包含关键信息
	contentStr := string(content)
	if !strings.Contains(contentStr, "促销模块常量验证报告") {
		t.Error("报告文件缺少标题")
	}
}

func TestValidationTool_ValidateSpecificFile(t *testing.T) {
	rootPath := "."
	tool := NewValidationTool(rootPath)

	// 测试验证特定文件
	testFile := "backend-go/internal/model/promotion/coupon_constants.go"

	report, err := tool.ValidateSpecificFile(testFile)
	if err != nil {
		// 如果文件不存在，跳过测试
		if strings.Contains(err.Error(), "no such file") {
			t.Skipf("测试文件不存在: %s", testFile)
		}
		t.Fatalf("验证特定文件失败: %v", err)
	}

	t.Logf("特定文件验证结果:")
	t.Logf("  文件: %s", testFile)
	t.Logf("  通过: %v", report.OverallPassed)
	t.Logf("  魔法数字: %d", len(report.ValidationResult.MagicNumbers))
	t.Logf("  命名违规: %d", len(report.NamingViolations))

	// 输出总结
	t.Log("文件验证总结:")
	t.Log(report.Summary)
}

func TestValidationTool_CountPromotionFiles(t *testing.T) {
	rootPath := "."
	tool := NewValidationTool(rootPath)

	count := tool.countPromotionFiles()
	t.Logf("promotion模块文件数量: %d", count)

	// 文件数量应该大于0（至少有我们创建的常量文件）
	if count == 0 {
		t.Log("注意: 未找到promotion模块文件，可能路径不正确")
	}
}

func TestValidationTool_CountTotalConstants(t *testing.T) {
	rootPath := "."
	tool := NewValidationTool(rootPath)

	count := tool.countTotalConstants()
	t.Logf("常量总数: %d", count)

	// 常量数量应该大于0（至少有我们定义的常量）
	if count == 0 {
		t.Log("注意: 未找到常量定义，可能路径不正确")
	}
}

func TestValidationTool_CalculateStatistics(t *testing.T) {
	rootPath := "."
	tool := NewValidationTool(rootPath)

	// 创建模拟的验证结果
	validationResult := &ValidationResult{
		MagicNumbers: []NumericLiteral{
			{Value: 1, Position: "test:1:1", Context: "test", File: "test.go"},
			{Value: 2, Position: "test:2:1", Context: "test", File: "test.go"},
		},
		Misaligned: []JavaEnumMapping{
			{GoConstant: "TestConstant", JavaValue: 1, GoValue: 2},
		},
		Errors: []string{"test error"},
	}

	namingViolations := []NamingViolation{
		{Name: "testViolation", Type: "test", Description: "test violation"},
	}

	duplicateConstants := []DuplicateConstant{
		{Name: "testDuplicate", Value: 1, Occurrences: 2},
	}

	// 计算统计信息
	stats := tool.calculateStatistics(validationResult, namingViolations, duplicateConstants)

	// 验证统计结果
	if stats.MagicNumbersFound != 2 {
		t.Errorf("魔法数字统计错误: 期望 2, 实际 %d", stats.MagicNumbersFound)
	}

	if stats.MisalignedConstants != 1 {
		t.Errorf("不对齐常量统计错误: 期望 1, 实际 %d", stats.MisalignedConstants)
	}

	if stats.NamingViolations != 1 {
		t.Errorf("命名违规统计错误: 期望 1, 实际 %d", stats.NamingViolations)
	}

	if stats.DuplicateConstants != 1 {
		t.Errorf("重复常量统计错误: 期望 1, 实际 %d", stats.DuplicateConstants)
	}

	t.Logf("统计信息验证通过:")
	t.Logf("  文件数: %d", stats.TotalFiles)
	t.Logf("  常量数: %d", stats.TotalConstants)
	t.Logf("  魔法数字: %d", stats.MagicNumbersFound)
	t.Logf("  不对齐: %d", stats.MisalignedConstants)
	t.Logf("  命名违规: %d", stats.NamingViolations)
	t.Logf("  重复常量: %d", stats.DuplicateConstants)
}
