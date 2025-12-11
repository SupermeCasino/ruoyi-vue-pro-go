package promotion

import (
	"time"
)

// PromotionDiscountActivity 限时折扣活动
type PromotionDiscountActivity struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Status    int       `gorm:"column:status" json:"status"` // 0: Disable, 1: Enable
	StartTime time.Time `gorm:"column:start_time" json:"startTime"`
	EndTime   time.Time `gorm:"column:end_time" json:"endTime"`
	Remark    string    `gorm:"column:remark" json:"remark"`

	Creator    string    `gorm:"column:creator" json:"creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Updater    string    `gorm:"column:updater" json:"updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"column:deleted" json:"deleted"`
	TenantID   int64     `gorm:"column:tenant_id" json:"tenantId"`
}

func (PromotionDiscountActivity) TableName() string {
	return "promotion_discount_activity"
}

// PromotionDiscountProduct 限时折扣商品
type PromotionDiscountProduct struct {
	ID              int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ActivityID      int64 `gorm:"column:activity_id" json:"activityId"`
	SpuID           int64 `gorm:"column:spu_id" json:"spuId"`
	SkuID           int64 `gorm:"column:sku_id" json:"skuId"`
	DiscountType    int   `gorm:"column:discount_type" json:"discountType"` // 1: Price, 2: Percent
	DiscountPercent int   `gorm:"column:discount_percent" json:"discountPercent"`
	DiscountPrice   int   `gorm:"column:discount_price" json:"discountPrice"`

	// Redundant fields for easier querying/display
	ActivityName      string    `gorm:"column:activity_name" json:"activityName"`
	ActivityStatus    int       `gorm:"column:activity_status" json:"activityStatus"`
	ActivityStartTime time.Time `gorm:"column:activity_start_time" json:"activityStartTime"`
	ActivityEndTime   time.Time `gorm:"column:activity_end_time" json:"activityEndTime"`

	Creator    string    `gorm:"column:creator" json:"creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Updater    string    `gorm:"column:updater" json:"updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"column:deleted" json:"deleted"`
	TenantID   int64     `gorm:"column:tenant_id" json:"tenantId"`
}

func (PromotionDiscountProduct) TableName() string {
	return "promotion_discount_product"
}
