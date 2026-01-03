package promotion

// AppPointActivityRespVO 积分商城活动 Response (App)
type AppPointActivityRespVO struct {
	ID          int64  `json:"id"`
	SpuID       int64  `json:"spuId"`
	Status      int    `json:"status"`
	Stock       int    `json:"stock"`      // 剩余库存
	TotalStock  int    `json:"totalStock"` // 总库存
	SpuName     string `json:"spuName"`
	PicUrl      string `json:"picUrl"`
	MarketPrice int    `json:"marketPrice"`
	Point       int    `json:"point"` // 最低积分
	Price       int    `json:"price"` // 最低金额
}

// AppPointActivityDetailRespVO 积分商城活动详情 Response (App)
type AppPointActivityDetailRespVO struct {
	ID         int64  `json:"id"`
	SpuID      int64  `json:"spuId"`
	Status     int    `json:"status"`
	Stock      int    `json:"stock"`
	TotalStock int    `json:"totalStock"`
	Remark     string `json:"remark"`
	Point      int    `json:"point"` // 最低积分
	Price      int    `json:"price"` // 最低金额

	Products []AppPointProductRespVO `json:"products"`
}

// AppPointProductRespVO 积分商城商品 Response (App)
type AppPointProductRespVO struct {
	ID    int64 `json:"id"`
	SkuID int64 `json:"skuId"`
	Count int   `json:"count"`
	Point int   `json:"point"`
	Price int   `json:"price"`
	Stock int   `json:"stock"`
}
