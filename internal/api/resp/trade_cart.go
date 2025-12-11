package resp

// AppCartListResp 购物车列表响应
type AppCartListResp struct {
	ValidList   []AppCartItem `json:"validList"`
	InvalidList []AppCartItem `json:"invalidList"`
}

// AppCartItem 购物车项
type AppCartItem struct {
	ID       int64           `json:"id"`
	Count    int             `json:"count"`
	Selected bool            `json:"selected"`
	Spu      *AppCartSpuInfo `json:"spu"`
	Sku      *AppCartSkuInfo `json:"sku"`
}

// AppCartSpuInfo 购物车 SPU 信息
type AppCartSpuInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	PicURL string `json:"picUrl"`
}

// AppCartSkuInfo 购物车 SKU 信息
type AppCartSkuInfo struct {
	ID         int64                    `json:"id"`
	Name       string                   `json:"name"`
	PicURL     string                   `json:"picUrl"`
	Price      int                      `json:"price"`
	Stock      int                      `json:"stock"`
	Properties []AppCartSkuPropertyInfo `json:"properties"`
}

// AppCartSkuPropertyInfo SKU 属性信息
type AppCartSkuPropertyInfo struct {
	PropertyID   int64  `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	ValueID      int64  `json:"valueId"`
	ValueName    string `json:"valueName"`
}
