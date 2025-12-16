package req

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

// ProductFavoritePageReq (Admin)
type ProductFavoritePageReq struct {
	core.PageParam
	UserId int64 `form:"userId"` // 用户编号
	SpuId  int64 `form:"spuId"`  // 商品 SPU 编号
}

// AppFavoritePageReq (App)
type AppFavoritePageReq struct {
	core.PageParam
}

// AppFavoriteReq (App) - For Create, Delete, Exists
type AppFavoriteReq struct {
	SpuId int64 `json:"spuId" form:"spuId" binding:"required"` // 商品 SPU 编号
}

type AppFavoriteCreateReq struct {
	SpuId int64 `json:"spuId" binding:"required"` // 商品 SPU 编号
}
