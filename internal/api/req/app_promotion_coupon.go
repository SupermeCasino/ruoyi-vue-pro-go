package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/types"
)

// AppCouponTakeReq 领取优惠券 Request
type AppCouponTakeReq struct {
	TemplateID types.FlexInt64 `json:"templateId" binding:"required"`
}

// AppCouponPageReq 我的优惠券分页 Request
type AppCouponPageReq struct {
	PageNo   int  `form:"pageNo,default=1"`
	PageSize int  `form:"pageSize,default=10"`
	Status   *int `form:"status"` // 1: 未使用, 2: 已使用, 3: 已过期
}

// AppCouponMatchReq 匹配优惠券 Request
type AppCouponMatchReq struct {
	Price       int64   `json:"price"`
	SkuIDs      []int64 `json:"skuIds"`
	CategoryIDs []int64 `json:"categoryIds"`
	SpuIDs      []int64 `json:"spuIds"`
}

// AppCouponTemplatePageReq 优惠券模板分页 Request (对齐 Java: AppCouponTemplatePageReqVO)
type AppCouponTemplatePageReq struct {
	PageNo       int   `form:"pageNo,default=1"`
	PageSize     int   `form:"pageSize,default=10"`
	SpuID        int64 `form:"spuId"`        // 商品 SPU 编号
	ProductScope *int  `form:"productScope"` // 使用类型
}

// AppCouponTemplateListReq 优惠券模板列表 Request (对齐 Java: AppCouponTemplateController.getCouponTemplateList)
type AppCouponTemplateListReq struct {
	SpuID        *int64 `form:"spuId"`        // 商品 SPU 编号
	ProductScope *int   `form:"productScope"` // 使用类型
	Count        int    `form:"count,default=10"`
}
