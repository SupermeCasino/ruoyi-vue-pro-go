package req

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

// ProductBrowseHistoryPageReq (Admin)
type ProductBrowseHistoryPageReq struct {
	core.PageParam
	UserId int64 `form:"userId"` // 用户编号
	SpuId  int64 `form:"spuId"`  // 商品 SPU 编号
}

// AppProductBrowseHistoryPageReq (App)
type AppProductBrowseHistoryPageReq struct {
	core.PageParam
}

// AppProductBrowseHistoryDeleteReq (App)
type AppProductBrowseHistoryDeleteReq struct {
	SpuIds []int64 `json:"spuIds" binding:"required"` // 商品 SPU 编号数组
}
