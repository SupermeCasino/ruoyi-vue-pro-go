package trade

// AppCartAddReq 添加购物车请求
type AppCartAddReq struct {
	SkuID int64 `json:"skuId" binding:"required"`
	Count int   `json:"count" binding:"required,min=1"`
}

// AppCartUpdateCountReq 更新购物车数量请求
type AppCartUpdateCountReq struct {
	ID    int64 `json:"id" binding:"required"`
	Count int   `json:"count" binding:"required,min=1"`
}

// AppCartUpdateSelectedReq 更新购物车选中状态请求
type AppCartUpdateSelectedReq struct {
	IDs      []int64 `json:"ids" binding:"required"`
	Selected *bool   `json:"selected" binding:"required"`
}

// AppCartResetReq 重置购物车请求
type AppCartResetReq struct {
	ID    int64 `json:"id" binding:"required"`
	SkuID int64 `json:"skuId" binding:"required"`
	Count int   `json:"count" binding:"required,min=1"`
}

// AppCartDeleteReq 删除购物车请求
type AppCartDeleteReq struct {
	IDs []int64 `form:"ids" binding:"required"`
}

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

// AppCartSpuInfo 购物车中的 SPU 信息
type AppCartSpuInfo struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	PicURL     string `json:"picUrl"`
	CategoryID int64  `json:"categoryId"`
	Stock      int    `json:"stock"`
	Status     int    `json:"status"`
}

// AppCartSkuInfo 购物车中的 SKU 信息
type AppCartSkuInfo struct {
	ID         int64                    `json:"id"`
	PicURL     string                   `json:"picUrl"`
	Price      int                      `json:"price"`
	Stock      int                      `json:"stock"`
	Properties []AppCartSkuPropertyInfo `json:"properties"`
}

// AppCartSkuPropertyInfo 购物车中的 SKU 属性信息
type AppCartSkuPropertyInfo struct {
	PropertyID   int64  `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	ValueID      int64  `json:"valueId"`
	ValueName    string `json:"valueName"`
}
