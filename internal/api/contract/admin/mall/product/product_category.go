package product

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductCategoryCreateReq 创建商品分类请求
type ProductCategoryCreateReq struct {
	ParentID model.FlexInt64  `json:"parentId" binding:"min=0"` // 允许为 0 (根节点)
	Name     string           `json:"name" binding:"required,max=255"`
	PicURL   string           `json:"picUrl"`
	Sort     model.FlexInt32  `json:"sort" binding:"min=0"`
	Status   *model.FlexInt32 `json:"status" binding:"required,oneof=0 1"`
}

// ProductCategoryUpdateReq 更新商品分类请求
type ProductCategoryUpdateReq struct {
	ID       model.FlexInt64  `json:"id" binding:"required"`
	ParentID model.FlexInt64  `json:"parentId" binding:"min=0"`
	Name     string           `json:"name" binding:"required,max=255"`
	PicURL   string           `json:"picUrl"`
	Sort     model.FlexInt32  `json:"sort" binding:"min=0"`
	Status   *model.FlexInt32 `json:"status" binding:"required,oneof=0 1"`
}

// ProductCategoryListReq 列表查询请求
type ProductCategoryListReq struct {
	Name      string  `form:"name"`
	ParentID  *int64  `form:"parentId"`
	ParentIDs []int64 `form:"parentIds"` // 父分类编号数组
	Status    *int32  `form:"status"`
}

// ProductCategoryResp 商品分类响应
type ProductCategoryResp struct {
	ID         int64     `json:"id"`
	ParentID   int64     `json:"parentId"`
	Name       string    `json:"name"`
	PicURL     string    `json:"picUrl"`
	Sort       int32     `json:"sort"`
	Status     int32     `json:"status"`
	CreateTime time.Time `json:"createTime"`
}
