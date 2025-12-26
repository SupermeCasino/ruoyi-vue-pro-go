package resp

// AppRewardActivityResp App 端满减送活动响应 VO
// 对应 Java: AppRewardActivityRespVO
type AppRewardActivityResp struct {
	ID                 int64                   `json:"id"`                 // 活动编号
	Status             int                     `json:"status"`             // 活动状态
	Name               string                  `json:"name"`               // 活动标题
	StartTime          int64                   `json:"startTime"`          // 开始时间（毫秒时间戳）
	EndTime            int64                   `json:"endTime"`            // 结束时间（毫秒时间戳）
	ConditionType      int                     `json:"conditionType"`      // 条件类型
	ProductScope       int                     `json:"productScope"`       // 商品范围
	ProductScopeValues []int64                 `json:"productScopeValues"` // 商品 SPU 编号的数组
	Rules              []AppRewardActivityRule `json:"rules"`              // 优惠规则的数组
}

// AppRewardActivityRule App 端满减送活动规则
// 对应 Java: AppRewardActivityRespVO.Rule
type AppRewardActivityRule struct {
	Limit                    int           `json:"limit"`                    // 优惠门槛
	DiscountPrice            int           `json:"discountPrice"`            // 优惠价格，单位：分
	FreeDelivery             bool          `json:"freeDelivery"`             // 是否包邮
	Point                    int           `json:"point"`                    // 赠送的积分
	GiveCouponTemplateCounts map[int64]int `json:"giveCouponTemplateCounts"` // 赠送的优惠劵
	Description              string        `json:"description"`              // 规则描述
}
