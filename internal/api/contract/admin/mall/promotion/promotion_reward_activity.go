package promotion

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// PromotionRewardActivityCreateReq 创建 Request
type PromotionRewardActivityCreateReq struct {
	Name               string    `json:"name" binding:"required"`
	StartTime          time.Time `json:"startTime" binding:"required"`
	EndTime            time.Time `json:"endTime" binding:"required"`
	Remark             string    `json:"remark"`
	ConditionType      int       `json:"conditionType" binding:"required"` // 10: 满N元, 20: 满N件
	ProductScope       int       `json:"productScope" binding:"required"`  // 1: All, 2: Spu, 3: Category
	ProductScopeValues []int64   `json:"productScopeValues"`               // Array
	Rules              []Rule    `json:"rules" binding:"required,dive"`
}

type Rule struct {
	Limit                    int           `json:"limit"`         // 门槛 (分 or 件)
	DiscountPrice            int           `json:"discountPrice"` // 优惠价格，单位：分
	FreeDelivery             bool          `json:"freeDelivery"`
	Point                    int           `json:"point"`
	GiveCouponTemplateCounts map[int64]int `json:"giveCouponTemplateCounts"`
}

// PromotionRewardActivityUpdateReq 更新 Request
type PromotionRewardActivityUpdateReq struct {
	ID                 int64     `json:"id" binding:"required"`
	Name               string    `json:"name" binding:"required"`
	StartTime          time.Time `json:"startTime" binding:"required"`
	EndTime            time.Time `json:"endTime" binding:"required"`
	Remark             string    `json:"remark"`
	ConditionType      int       `json:"conditionType" binding:"required"`
	ProductScope       int       `json:"productScope" binding:"required"`
	ProductScopeValues []int64   `json:"productScopeValues"`
	Rules              []Rule    `json:"rules" binding:"required,dive"`
}

// PromotionRewardActivityPageReq 分页 Request
type PromotionRewardActivityPageReq struct {
	pagination.PageParam
	Name       string       `form:"name"`
	Status     *int         `form:"status"` // 0: Open, 1: Close
	CreateTime []*time.Time `form:"createTime"`
}

// PromotionRewardActivityResp 满减送活动 Response VO
type PromotionRewardActivityResp struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	Status             int       `json:"status"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	Remark             string    `json:"remark"`
	ConditionType      int       `json:"conditionType"`
	ProductScope       int       `json:"productScope"`
	ProductScopeValues []int64   `json:"productScopeValues"`
	Rules              []Rule    `json:"rules"`
	CreateTime         time.Time `json:"createTime"`
}
