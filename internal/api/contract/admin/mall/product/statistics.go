package product

import "time"

// ProductStatisticsReqVO 商品统计请求
type ProductStatisticsReqVO struct {
	Times []time.Time `form:"times" binding:"required,len=2" time_format:"2006-01-02 15:04:05"` // 时间范围 [开始时间, 结束时间]
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

// DataComparisonRespVO 数据对比响应 (泛型)
type DataComparisonRespVO[T any] struct {
	Summary      *T  `json:"summary"`      // 当前数据
	Comparison   *T  `json:"comparison"`   // 对比数据
	IncreaseRate int `json:"increaseRate"` // 增长率
}
