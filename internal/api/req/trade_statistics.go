package req

import "time"

// TradeStatisticsAnalysisReq 交易状况分析 Request
type TradeStatisticsAnalysisReq struct {
	Times []time.Time `form:"times[]" time_format:"2006-01-02 15:04:05"` // 时间范围 [start, end]
}
