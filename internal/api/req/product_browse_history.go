package req

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// ProductBrowseHistoryPageReq (Admin)
type ProductBrowseHistoryPageReq struct {
	pagination.PageParam
	UserId      int64       `form:"userId"`      // 用户编号
	UserDeleted *bool       `form:"userDeleted"` // 用户是否删除
	SpuId       int64       `form:"spuId"`       // 商品 SPU 编号
	CreateTime  []time.Time `form:"createTime"`  // 创建时间
}

// AppProductBrowseHistoryPageReq (App)
type AppProductBrowseHistoryPageReq struct {
	pagination.PageParam
}

// AppProductBrowseHistoryDeleteReq (App)
type AppProductBrowseHistoryDeleteReq struct {
	SpuIds []int64 `json:"spuIds" binding:"required"` // 商品 SPU 编号数组
}
