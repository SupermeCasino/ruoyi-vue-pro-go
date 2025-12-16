package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// ProductFavoritePageReq (Admin)
type ProductFavoritePageReq struct {
	pagination.PageParam
	UserId int64 `form:"userId"` // 用户编号
	SpuId  int64 `form:"spuId"`  // 商品 SPU 编号
}

// AppFavoritePageReq (App)
type AppFavoritePageReq struct {
	pagination.PageParam
}

// AppFavoriteReq (App) - For Create, Delete, Exists
type AppFavoriteReq struct {
	SpuId int64 `json:"spuId" form:"spuId" binding:"required"` // 商品 SPU 编号
}

type AppFavoriteCreateReq struct {
	SpuId int64 `json:"spuId" binding:"required"` // 商品 SPU 编号
}
