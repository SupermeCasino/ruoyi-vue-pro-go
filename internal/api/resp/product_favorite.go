package resp

import "time"

// ProductFavoriteResp (Admin) - 继承 ProductSpuRespVO 全部字段
type ProductFavoriteResp struct {
	// 收藏字段
	UserID int64 `json:"userId"`
	SpuID  int64 `json:"spuId"`
	// SPU 字段
	ID                 int64             `json:"id"`
	Name               string            `json:"name"`
	Keyword            string            `json:"keyword"`
	Introduction       string            `json:"introduction"`
	Description        string            `json:"description"`
	CategoryID         int64             `json:"categoryId"`
	BrandID            int64             `json:"brandId"`
	PicURL             string            `json:"picUrl"`
	SliderPicURLs      []string          `json:"sliderPicUrls"`
	Sort               int               `json:"sort"`
	Status             int               `json:"status"`
	SpecType           bool              `json:"specType"`
	Price              int               `json:"price"`
	MarketPrice        int               `json:"marketPrice"`
	CostPrice          int               `json:"costPrice"`
	Stock              int               `json:"stock"`
	DeliveryTypes      []int             `json:"deliveryTypes"`
	DeliveryTemplateID int64             `json:"deliveryTemplateId"`
	GiveIntegral       int               `json:"giveIntegral"`
	SubCommissionType  bool              `json:"subCommissionType"`
	SalesCount         int               `json:"salesCount"`
	VirtualSalesCount  int               `json:"virtualSalesCount"`
	BrowseCount        int               `json:"browseCount"`
	CreateTime          time.Time         `json:"createTime"`
	Skus               []*ProductSkuResp `json:"skus,omitempty"`
}

// AppFavoriteResp (App)
type AppFavoriteResp struct {
	ID        int64     `json:"id"`
	SpuID     int64     `json:"spuId"`
	SpuName   string    `json:"spuName"`
	PicURL    string    `json:"picUrl"`
	Price     int64     `json:"price"`
	CreateTime time.Time `json:"createdAt"`
}
