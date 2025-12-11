package req

// AppCouponTakeReq 领取优惠券 Request
type AppCouponTakeReq struct {
	TemplateID int64 `json:"templateId" binding:"required"`
}

// AppCouponPageReq 我的优惠券分页 Request
type AppCouponPageReq struct {
	PageNo   int  `form:"pageNo,default=1"`
	PageSize int  `form:"pageSize,default=10"`
	Status   *int `form:"status"` // 1: Unused, 2: Used, 3: Expired
}

// AppCouponMatchReq 匹配优惠券 Request
type AppCouponMatchReq struct {
	Price       int64   `json:"price"`
	SkuIDs      []int64 `json:"skuIds"`
	CategoryIDs []int64 `json:"categoryIds"` // Optional, if frontend knows. But usually backend resolves.
	SpuIDs      []int64 `json:"spuIds"`      // Optional
}
