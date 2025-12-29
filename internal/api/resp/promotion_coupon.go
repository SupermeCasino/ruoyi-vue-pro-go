package resp

import "time"

// CouponTemplateResp 管理后台 - 优惠券模板响应 (对齐 Java: CouponTemplateRespVO)
type CouponTemplateResp struct {
	// RespVO 字段
	ID         int64     `json:"id"`
	Status     int       `json:"status"`
	TakeCount  int       `json:"takeCount"` // 领取数量
	UseCount   int       `json:"useCount"`  // 使用数量
	CreateTime time.Time `json:"createTime"`

	// BaseVO 字段
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	TotalCount         int        `json:"totalCount"`
	TakeLimitCount     int        `json:"takeLimitCount"`
	TakeType           int        `json:"takeType"`
	UsePrice           int        `json:"usePrice"`
	ProductScope       int        `json:"productScope"`
	ProductScopeValues []int64    `json:"productScopeValues"`
	ValidityType       int        `json:"validityType"`
	ValidStartTime     *time.Time `json:"validStartTime"`
	ValidEndTime       *time.Time `json:"validEndTime"`
	FixedStartTerm     *int       `json:"fixedStartTerm"`
	FixedEndTerm       *int       `json:"fixedEndTerm"`
	DiscountType       int        `json:"discountType"`
	DiscountPercent    *int       `json:"discountPercent"`
	DiscountPrice      *int       `json:"discountPrice"`
	DiscountLimitPrice *int       `json:"discountLimitPrice"`
}

// CouponPageResp 管理后台 - 优惠券分页响应 (对齐 Java: CouponPageItemRespVO)
type CouponPageResp struct {
	// RespVO 字段
	ID         int64     `json:"id"`
	CreateTime time.Time `json:"createTime"`

	// BaseVO 字段 - 基本信息
	TemplateID int64  `json:"templateId"`
	Name       string `json:"name"`
	Status     int    `json:"status"`

	// BaseVO 字段 - 领取情况
	UserID   int64 `json:"userId"`
	TakeType int   `json:"takeType"`

	// BaseVO 字段 - 使用规则
	UsePrice           int       `json:"usePrice"`
	ValidStartTime     time.Time `json:"validStartTime"`
	ValidEndTime       time.Time `json:"validEndTime"`
	ProductScope       int       `json:"productScope"`
	ProductScopeValues []int64   `json:"productScopeValues"`

	// BaseVO 字段 - 使用效果
	DiscountType       int  `json:"discountType"`
	DiscountPercent    *int `json:"discountPercent"`
	DiscountPrice      *int `json:"discountPrice"`
	DiscountLimitPrice *int `json:"discountLimitPrice"`

	// BaseVO 字段 - 使用情况
	UseOrderID *int64     `json:"useOrderId"`
	UseTime    *time.Time `json:"useTime"`

	// PageItemRespVO 字段 - 关联字段
	Nickname string `json:"nickname"` // 用户昵称（关联查询 member_user）
}
