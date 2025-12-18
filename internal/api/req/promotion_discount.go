package req

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type DiscountActivityPageReq struct {
	pagination.PageParam
	Name       string       `form:"name"`
	Status     *int         `form:"status"`
	CreateTime []*time.Time `form:"createTime"`
}

type DiscountActivityCreateReq struct {
	Name      string               `json:"name" binding:"required"`
	StartTime time.Time            `json:"startTime" binding:"required"`
	EndTime   time.Time            `json:"endTime" binding:"required"`
	Remark    string               `json:"remark"`
	Products  []DiscountProductReq `json:"products" binding:"required,dive"`
}

type DiscountActivityUpdateReq struct {
	ID        int64                `json:"id" binding:"required"`
	Name      string               `json:"name" binding:"required"`
	StartTime time.Time            `json:"startTime" binding:"required"`
	EndTime   time.Time            `json:"endTime" binding:"required"`
	Remark    string               `json:"remark"`
	Products  []DiscountProductReq `json:"products" binding:"required,dive"`
}

type DiscountProductReq struct {
	SpuID           int64 `json:"spuId" binding:"required"`
	SkuID           int64 `json:"skuId" binding:"required"`
	DiscountType    int   `json:"discountType" binding:"required"` // 1: Price, 2: Percent
	DiscountPercent int   `json:"discountPercent"`
	DiscountPrice   int   `json:"discountPrice"`
}
