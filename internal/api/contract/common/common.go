package common

// DataComparisonRespVO 数据对比响应 VO
type DataComparisonRespVO[T any] struct {
	Summary    *T `json:"summary"`    // 今日/本月数据
	Comparison *T `json:"comparison"` // 昨日/上月数据
}
