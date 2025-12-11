package resp

// TradeStatisticsSummaryResp 交易统计数据摘要 Response
type TradeStatisticsSummaryResp struct {
	Yesterday *TradeStatisticsData `json:"yesterday"`
	Today     *TradeStatisticsData `json:"today"`
	Month     *TradeStatisticsData `json:"month"`
}

type TradeStatisticsData struct {
	OrderPayPrice        int   `json:"orderPayPrice"`        // 订单支付金额（分）
	OrderPayCount        int64 `json:"orderPayCount"`        // 订单支付数量
	AfterSaleCount       int64 `json:"afterSaleCount"`       // 退单数量
	AfterSaleRefundPrice int   `json:"afterSaleRefundPrice"` // 退款金额（分）
}

// TradeStatisticsAnalysisResp 交易状况分析 Response
type TradeStatisticsAnalysisResp struct {
	Dates                []string `json:"dates"`
	OrderPayPrice        []int    `json:"orderPayPrice"`
	OrderPayCount        []int64  `json:"orderPayCount"`
	AfterSaleCount       []int64  `json:"afterSaleCount"`
	AfterSaleRefundPrice []int    `json:"afterSaleRefundPrice"`
}
