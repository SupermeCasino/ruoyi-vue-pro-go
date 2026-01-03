package member

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/common"
)

// MemberAnalyseReqVO 会员分析请求
type MemberAnalyseReqVO struct {
	Times []time.Time `form:"times" binding:"required,len=2" time_format:"2006-01-02 15:04:05"` // 时间范围 [开始时间, 结束时间]
}

// MemberSummaryRespVO 会员统计摘要响应
type MemberSummaryRespVO struct {
	TotalUserCount  int64 `json:"totalUserCount"`  // 总用户数
	ActiveUserCount int64 `json:"activeUserCount"` // 活跃用户数
	RegisterCount   int64 `json:"registerCount"`   // 注册数
	VisitUserCount  int64 `json:"visitUserCount"`  // 访客数
	OrderUserCount  int64 `json:"orderUserCount"`  // 下单用户数
	PayUserCount    int64 `json:"payUserCount"`    // 支付用户数
}

// MemberAnalyseRespVO 会员分析响应
type MemberAnalyseRespVO struct {
	VisitUserCount int64                            `json:"visitUserCount"` // 访客数
	OrderUserCount int64                            `json:"orderUserCount"` // 下单用户数
	PayUserCount   int64                            `json:"payUserCount"`   // 支付用户数
	ATV            int64                            `json:"atv"`            // 客单价
	ComparisonData common.DataComparisonRespVO[any] `json:"comparisonData"` // 对比数据
}

// MemberAreaStatisticsRespVO 会员地区统计响应
type MemberAreaStatisticsRespVO struct {
	AreaName  string `json:"areaName"`  // 地区名称
	UserCount int64  `json:"userCount"` // 用户数
}

// MemberSexStatisticsRespVO 会员性别统计响应
type MemberSexStatisticsRespVO struct {
	SexName   string `json:"sexName"`   // 性别名称
	UserCount int64  `json:"userCount"` // 用户数
}

// MemberTerminalStatisticsRespVO 会员终端统计响应
type MemberTerminalStatisticsRespVO struct {
	TerminalName string `json:"terminalName"` // 终端名称
	UserCount    int64  `json:"userCount"`    // 用户数
}

// MemberCountRespVO 会员数量响应
type MemberCountRespVO struct {
	Date      string `json:"date"`      // 日期
	UserCount int64  `json:"userCount"` // 用户数
}

// MemberRegisterCountRespVO 会员注册数量响应
type MemberRegisterCountRespVO struct {
	Date          string `json:"date"`          // 日期
	RegisterCount int64  `json:"registerCount"` // 注册数
}
