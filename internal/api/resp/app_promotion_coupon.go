package resp

import "time"

// AppCouponResp 用户 App - 优惠劵响应 (对齐 Java: AppCouponRespVO)
type AppCouponResp struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Status             int        `json:"status"`             // 参见 CouponStatusEnum: 1-未使用 2-已使用 3-已过期
	UsePrice           int        `json:"usePrice"`           // 满多少可用 (单位: 分), 0-不限制
	ProductScope       int        `json:"productScope"`       // 商品范围: 1-全部 2-分类 3-商品
	ProductScopeValues []int64    `json:"productScopeValues"` // 商品范围编号数组
	ValidStartTime     *time.Time `json:"validStartTime"`     // 生效开始时间
	ValidEndTime       *time.Time `json:"validEndTime"`       // 生效结束时间
	DiscountType       int        `json:"discountType"`       // 优惠类型: 1-满减 2-折扣
	DiscountPercent    int        `json:"discountPercent"`    // 折扣百分比 (80 表示 80%)
	DiscountPrice      int        `json:"discountPrice"`      // 优惠金额 (单位: 分)
	DiscountLimitPrice int        `json:"discountLimitPrice"` // 折扣上限 (单位: 分)
}

// AppCouponTemplateResp 用户 App - 优惠劵模板响应 (对齐 Java: AppCouponTemplateRespVO)
type AppCouponTemplateResp struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`        // 优惠券说明
	TotalCount         int         `json:"totalCount"`         // 发行总量, -1 表示不限制
	TakeLimitCount     int         `json:"takeLimitCount"`     // 每人限领个数, -1 表示不限制
	UsePrice           int         `json:"usePrice"`           // 满多少可用 (单位: 分), 0-不限制
	ProductScope       int         `json:"productScope"`       // 商品范围: 1-全部 2-分类 3-商品
	ProductScopeValues []int64     `json:"productScopeValues"` // 商品范围编号数组
	ValidityType       int         `json:"validityType"`       // 生效日期类型: 1-固定日期 2-领取后N天
	ValidStartTime     interface{} `json:"validStartTime"`     // 固定日期 - 生效开始时间 (毫秒级时间戳或null)
	ValidEndTime       interface{} `json:"validEndTime"`       // 固定日期 - 生效结束时间 (毫秒级时间戳或null)
	FixedStartTerm     interface{} `json:"fixedStartTerm"`     // 领取日期 - 开始天数 (int或null)
	FixedEndTerm       interface{} `json:"fixedEndTerm"`       // 领取日期 - 结束天数 (int或null)
	DiscountType       int         `json:"discountType"`       // 优惠类型: 1-满减 2-折扣
	DiscountPercent    int         `json:"discountPercent"`    // 折扣百分比 (80 表示 80%)
	DiscountPrice      int         `json:"discountPrice"`      // 优惠金额 (单位: 分)
	DiscountLimitPrice int         `json:"discountLimitPrice"` // 折扣上限 (单位: 分)
	TakeCount          int         `json:"takeCount"`          // 已领取数量
	CanTake            bool        `json:"canTake"`            // 当前用户是否可领取
}
