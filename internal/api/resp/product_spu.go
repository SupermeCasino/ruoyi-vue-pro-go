package resp

import "time"

// ProductSpuResp 商品 SPU Response
type ProductSpuResp struct {
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
	Skus               []*ProductSkuResp `json:"skus,omitempty"` // 详情时返回
}

// ProductSkuResp 商品 SKU Response
type ProductSkuResp struct {
	ID                   int64                    `json:"id"`
	Name                 string                   `json:"name"` // ✅ 对齐 Java: 商品 SKU 名字
	SpuID                int64                    `json:"spuId"`
	Properties           []ProductSkuPropertyResp `json:"properties"`
	Price                int                      `json:"price"`
	MarketPrice          int                      `json:"marketPrice"`
	CostPrice            int                      `json:"costPrice"`
	BarCode              string                   `json:"barCode"`
	PicURL               string                   `json:"picUrl"`
	Stock                int                      `json:"stock"`
	Weight               float64                  `json:"weight"`
	Volume               float64                  `json:"volume"`
	FirstBrokeragePrice  int                      `json:"firstBrokeragePrice"`
	SecondBrokeragePrice int                      `json:"secondBrokeragePrice"`
	SalesCount           int                      `json:"salesCount"`
	VipPrice             int                      `json:"vipPrice"`
}

type ProductSkuPropertyResp struct {
	PropertyID   int64  `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	ValueID      int64  `json:"valueId"`
	ValueName    string `json:"valueName"`
}

// ProductSpuSimpleResp 精简商品 Response
type ProductSpuSimpleResp struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PicURL      string `json:"picUrl"`
	Price       int    `json:"price"`
	MarketPrice int    `json:"marketPrice"`
	CostPrice   int    `json:"costPrice"`
	Stock       int    `json:"stock"`
}

// ProductSpuExportVO 商品导出 VO
type ProductSpuExportVO struct {
	ID         int64  `json:"id" label:"商品编号"`
	Name       string `json:"name" label:"商品名称"`
	CategoryID int64  `json:"categoryId" label:"分类编号"`
	Price      int    `json:"price" label:"价格(分)"`
	Stock      int    `json:"stock" label:"库存"`
	Status     int    `json:"status" label:"状态"`
	SalesCount int    `json:"salesCount" label:"销量"`
}
