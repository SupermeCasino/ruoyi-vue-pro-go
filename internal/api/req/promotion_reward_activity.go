package req

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"time"
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
	Sort               int       `json:"sort"`
}

type Rule struct {
	Limit       int `json:"limit"`       // 门槛 (分 or 件)
	ReducePrice int `json:"reducePrice"` // 减多少分
	// TODO: Add Gift, Point, Coupon logic if needed
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
	Sort               int       `json:"sort"`
}

// PromotionRewardActivityPageReq 分页 Request
type PromotionRewardActivityPageReq struct {
	core.PageParam
	Name   string `form:"name"`
	Status *int   `form:"status"` // 0: Open, 1: Close
}
