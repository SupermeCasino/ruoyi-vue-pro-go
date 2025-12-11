package req

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
