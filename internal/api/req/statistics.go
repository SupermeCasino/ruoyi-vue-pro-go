package req

import "time"

// TradeStatisticsReqVO 交易统计请求
type TradeStatisticsReqVO struct {
	Times []time.Time `form:"times" binding:"required,len=2"` // 时间范围 [开始时间, 结束时间]
}

// ProductStatisticsReqVO 商品统计请求
type ProductStatisticsReqVO struct {
	Times []time.Time `form:"times" binding:"required,len=2"` // 时间范围 [开始时间, 结束时间]
}

// MemberAnalyseReqVO 会员分析请求
type MemberAnalyseReqVO struct {
	Times []time.Time `form:"times" binding:"required,len=2"` // 时间范围 [开始时间, 结束时间]
}

// TradeOrderTrendReqVO 订单趋势请求
type TradeOrderTrendReqVO struct {
	Times []time.Time `form:"times" binding:"required,len=2"` // 时间范围 [开始时间, 结束时间]
}
