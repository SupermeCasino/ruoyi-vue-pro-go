package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductBrowseHistory 商品浏览记录
type ProductBrowseHistory struct {
	ID          int64         `gorm:"primaryKey;autoIncrement;comment:记录编号" json:"id"`
	UserID      int64         `gorm:"index;not null;comment:用户编号" json:"userId"`
	SpuID       int64         `gorm:"index;not null;comment:商品 SPU 编号" json:"spuId"`
	UserDeleted model.BitBool `gorm:"column:user_deleted;comment:用户是否删除" json:"userDeleted"`
	model.TenantBaseDO
}

func (*ProductBrowseHistory) TableName() string {
	return "product_browse_history"
}
