package resp

import "time"

// AppProductSpuResp 用户 APP - 商品 SPU Response VO (List)
type AppProductSpuResp struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Introduction  string    `json:"introduction"` // 商品简介
	CategoryID    int64     `json:"categoryId"`   // 分类编号
	PicURL        string    `json:"picUrl"`
	SliderPicURLs []string  `json:"sliderPicUrls"` // 轮播图
	SpecType      bool      `json:"specType"`      // 规格类型
	Price         int       `json:"price"`         // 最低价
	MarketPrice   int       `json:"marketPrice"`   // 市场价
	Stock         int       `json:"stock"`
	SalesCount    int       `json:"salesCount"`    // 销量 + 虚拟销量
	DeliveryTypes []int     `json:"deliveryTypes"` // 配送方式
	VIPPrice      int       `json:"vipPrice"`      // VIP 价格
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"createTime"`
}
