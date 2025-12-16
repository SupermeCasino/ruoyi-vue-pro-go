package resp

import "time"

// AppProductSpuResp 用户 APP - 商品 SPU Response VO (List)
type AppProductSpuResp struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	PicURL      string    `json:"picUrl"`
	Price       int       `json:"price"`       // 最低价
	MarketPrice int       `json:"marketPrice"` // 市场价
	SalesCount  int       `json:"salesCount"`  // 销量 + 虚拟销量
	VIPPrice    int       `json:"vipPrice"`    // VIP 价格
	Stock       int       `json:"stock"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"createTime"`
}
