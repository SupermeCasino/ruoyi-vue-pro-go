package resp

import "time"

// DataComparisonRespVO 数据对比响应
type DataComparisonRespVO[T any] struct {
	Summary      *T   `json:"summary"`      // 当前数据
	Comparison   *T   `json:"comparison"`   // 对比数据
	IncreaseRate int  `json:"increaseRate"` // 增长率
	IsIncrease   bool `json:"isIncrease"`   // 是否增长
}

// TradeSummaryRespVO 交易统计摘要响应
type TradeSummaryRespVO struct {
	Yesterday *TradeSummaryItemVO `json:"yesterday"`
	Month     *TradeSummaryItemVO `json:"month"`
}

// TradeSummaryItemVO 交易统计单项响应 - 严格对齐 Java TradeStatisticsDO
type TradeSummaryItemVO struct {
	OrderCreateCount         int64 `json:"orderCreateCount"`         // 创建订单数
	OrderPayCount            int64 `json:"orderPayCount"`            // 支付订单商品数
	OrderPayPrice            int64 `json:"orderPayPrice"`            // 总支付金额(分)
	AfterSaleCount           int64 `json:"afterSaleCount"`           // 退款订单数
	AfterSaleRefundPrice     int64 `json:"afterSaleRefundPrice"`     // 总退款金额(分)
	BrokerageSettlementPrice int64 `json:"brokerageSettlementPrice"` // 佣金金额已结算(分)
	WalletPayPrice           int64 `json:"walletPayPrice"`           // 总支付金额余额(分)
	RechargePayCount         int64 `json:"rechargePayCount"`         // 充值订单数
	RechargePayPrice         int64 `json:"rechargePayPrice"`         // 充值金额(分)
	RechargeRefundCount      int64 `json:"rechargeRefundCount"`      // 充值退款订单数
	RechargeRefundPrice      int64 `json:"rechargeRefundPrice"`      // 充值退款金额(分)
}

// TradeTrendSummaryRespVO 交易趋势摘要响应 - 严格对齐 Java TradeStatisticsDO
type TradeTrendSummaryRespVO struct {
	StatisticsTime           time.Time `json:"statisticsTime" label:"日期"`             // 统计时间
	OrderCreateCount         int64     `json:"orderCreateCount" label:"创建订单数"`        // 创建订单数
	OrderPayCount            int64     `json:"orderPayCount" label:"支付订单商品数"`         // 支付订单商品数
	OrderPayPrice            int64     `json:"orderPayPrice" label:"总支付金额"`           // 总支付金额(分)
	AfterSaleCount           int64     `json:"afterSaleCount" label:"退款订单数"`          // 退款订单数
	AfterSaleRefundPrice     int64     `json:"afterSaleRefundPrice" label:"总退款金额"`    // 总退款金额(分)
	BrokerageSettlementPrice int64     `json:"brokerageSettlementPrice" label:"佣金金额"` // 佣金金额已结算(分)
	WalletPayPrice           int64     `json:"walletPayPrice" label:"余额支付"`           // 总支付金额余额(分)
	RechargePayCount         int64     `json:"rechargePayCount" label:"充值订单数"`        // 充值订单数
	RechargePayPrice         int64     `json:"rechargePayPrice" label:"充值金额"`         // 充值金额(分)
	RechargeRefundCount      int64     `json:"rechargeRefundCount" label:"充值退款订单数"`   // 充值退款订单数
	RechargeRefundPrice      int64     `json:"rechargeRefundPrice" label:"充值退款金额"`    // 充值退款金额(分)
}

// TradeOrderCountRespVO 交易订单数量响应
type TradeOrderCountRespVO struct {
	UndeliveredCount      int64 `json:"undeliveredCount"`      // 待发货数
	PickUpCount           int64 `json:"pickUpCount"`           // 待自提数
	AfterSaleApplyCount   int64 `json:"afterSaleApplyCount"`   // 售后申请数
	AuditingWithdrawCount int64 `json:"auditingWithdrawCount"` // 待审核提现数
}

// TradeOrderSummaryRespVO 交易订单摘要响应
type TradeOrderSummaryRespVO struct {
	OrderCount int64 `json:"orderCount"` // 订单数
	PayPrice   int64 `json:"payPrice"`   // 支付金额
}

// TradeOrderTrendRespVO 交易订单趋势响应
type TradeOrderTrendRespVO struct {
	StatisticsTime time.Time `json:"statisticsTime"` // 统计时间
	OrderCount     int64     `json:"orderCount"`     // 订单数
	PayPrice       int64     `json:"payPrice"`       // 支付金额
}

// ProductStatisticsRespVO 商品统计响应
type ProductStatisticsRespVO struct {
	ID             int64     `json:"id"`
	StatisticsTime time.Time `json:"statisticsTime"`
	SpuID          int64     `json:"spuId"`
	Name           string    `json:"name"`          // 商品名称
	PicUrl         string    `json:"picUrl"`        // 商品图片
	BuyCount       int64     `json:"buyCount"`      // 购买数
	BuyPrice       int64     `json:"buyPrice"`      // 购买金额
	BrowseCount    int64     `json:"browseCount"`   // 浏览数
	FavoriteCount  int64     `json:"favoriteCount"` // 收藏数
	CommentCount   int64     `json:"commentCount"`  // 评价数
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
	VisitUserCount int64                             `json:"visitUserCount"` // 访客数
	OrderUserCount int64                             `json:"orderUserCount"` // 下单用户数
	PayUserCount   int64                             `json:"payUserCount"`   // 支付用户数
	ATV            int64                             `json:"atv"`            // 客单价
	ComparisonData DataComparisonRespVO[interface{}] `json:"comparisonData"` // 对比数据
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

// PaySummaryRespVO 支付摘要响应
type PaySummaryRespVO struct {
	RechargePrice int64 `json:"rechargePrice"` // 充值金额
}
