package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PromotionPointActivity 积分商城活动
// 对应 Java: PointActivityDO
type PromotionPointActivity struct {
	ID         int64  `gorm:"primaryKey;autoIncrement;comment:活动编号"`
	SpuID      int64  `gorm:"column:spu_id;type:bigint;not null;comment:商品SPU编号"`
	Status     int    `gorm:"column:status;type:int;not null;comment:活动状态"` // 0: Disable, 1: Enable
	Remark     string `gorm:"column:remark;size:255;default:'';comment:备注"`
	Sort       int    `gorm:"column:sort;type:int;not null;default:0;comment:排序"`
	Stock      int    `gorm:"column:stock;type:int;not null;default:0;comment:活动库存"` // 剩余库存
	TotalStock int    `gorm:"column:total_stock;type:int;not null;default:0;comment:活动总库存"`
	model.TenantBaseDO
}

func (PromotionPointActivity) TableName() string {
	return "promotion_point_activity"
}
