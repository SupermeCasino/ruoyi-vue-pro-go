package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
)

func main() {
	// 命令行参数
	var (
		rootPath   = flag.String("root", ".", "项目根目录路径")
		outputFile = flag.String("output", "", "输出报告文件路径")
		fileMode   = flag.String("file", "", "验证特定文件")
		verbose    = flag.Bool("verbose", false, "详细输出")
	)
	flag.Parse()

	fmt.Println("🔍 促销模块常量验证工具")
	fmt.Println("========================")

	// 创建验证工具
	tool := promotion.NewValidationTool(*rootPath)

	var report *promotion.ComprehensiveValidationReport
	var err error

	// 根据模式运行验证
	if *fileMode != "" {
		fmt.Printf("正在验证文件: %s\n", *fileMode)
		report, err = tool.ValidateSpecificFile(*fileMode)
	} else {
		fmt.Printf("正在验证项目: %s\n", *rootPath)
		report, err = tool.RunComprehensiveValidation()
	}

	if err != nil {
		log.Fatalf("验证失败: %v", err)
	}

	// 输出结果
	if *verbose {
		fmt.Println(tool.GenerateDetailedReport(report))
	} else {
		fmt.Println(report.Summary)
	}

	// 保存报告到文件
	if *outputFile != "" {
		if err := tool.SaveReportToFile(report, *outputFile); err != nil {
			log.Printf("保存报告失败: %v", err)
		} else {
			fmt.Printf("\n📄 报告已保存到: %s\n", *outputFile)
		}
	}

	// 设置退出码
	if !report.OverallPassed {
		fmt.Println("\n❌ 验证未通过，请修复发现的问题")
		os.Exit(1)
	} else {
		fmt.Println("\n✅ 验证通过！")
		os.Exit(0)
	}
}
