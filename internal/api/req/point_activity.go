package req

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// PointActivityCreateReq 创建积分商城活动 Request
type PointActivityCreateReq struct {
	SpuID    int64                 `json:"spuId"`
	Status   int                   `json:"status"`
	Remark   string                `json:"remark"`
	Sort     int                   `json:"sort"`
	Products []PointProductSaveReq `json:"products"`
}

// PointActivityUpdateReq 更新积分商城活动 Request
type PointActivityUpdateReq struct {
	ID       int64                 `json:"id"`
	SpuID    int64                 `json:"spuId"`
	Status   int                   `json:"status"`
	Remark   string                `json:"remark"`
	Sort     int                   `json:"sort"`
	Products []PointProductSaveReq `json:"products"`
}

// PointProductSaveReq 保存积分商城商品 Request
type PointProductSaveReq struct {
	SkuID int64 `json:"skuId"`
	Count int   `json:"count"`
	Point int   `json:"point"`
	Price int   `json:"price"` // 单位：分
	Stock int   `json:"stock"`
}

// PointActivityPageReq 积分商城活动分页 Request
type PointActivityPageReq struct {
	pagination.PageParam
	Status     *int         `form:"status"`
	CreateTime []*time.Time `form:"createTime"`
}
