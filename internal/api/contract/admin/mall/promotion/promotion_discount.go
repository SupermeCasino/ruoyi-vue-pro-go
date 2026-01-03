package promotion

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

// DiscountProductRespVO 限时折扣商品 Response VO
type DiscountProductRespVO struct {
	ID              int64 `json:"id"`
	ActivityID      int64 `json:"activityId"`
	SpuID           int64 `json:"spuId"`
	SkuID           int64 `json:"skuId"`
	DiscountType    int   `json:"discountType"`
	DiscountPercent int   `json:"discountPercent"`
	DiscountPrice   int   `json:"discountPrice"`
}

// DiscountActivityRespVO 限时折扣活动 Response VO
type DiscountActivityRespVO struct {
	ID         int64                    `json:"id"`
	Name       string                   `json:"name"`
	Status     int                      `json:"status"`
	StartTime  time.Time                `json:"startTime"`
	EndTime    time.Time                `json:"endTime"`
	Remark     string                   `json:"remark"`
	CreateTime time.Time                `json:"createTime"`
	Products   []*DiscountProductRespVO `json:"products"`
}
