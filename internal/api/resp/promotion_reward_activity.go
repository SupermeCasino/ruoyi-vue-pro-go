package resp

import "time"

type PromotionRewardActivityResp struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	Status             int       `json:"status"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	ProductScope       int       `json:"productScope"`
	ProductScopeValues []int64   `json:"productScopeValues"`
	ConditionType      int       `json:"conditionType"`
	Rules              []Rule    `json:"rules"`
	Remark             string    `json:"remark"`
	CreateTime         time.Time `json:"createTime"`
}

type Rule struct {
	Limit                    int           `json:"limit"`
	ReducePrice              int           `json:"reducePrice"`
	FreeDelivery             bool          `json:"freeDelivery"`
	Point                    int           `json:"point"`
	GiveCouponTemplateCounts map[int64]int `json:"giveCouponTemplateCounts"`
}
