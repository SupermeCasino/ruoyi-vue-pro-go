package product

import (
	"time"
)

// ProductBrowseHistory 商品浏览记录
type ProductBrowseHistory struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;comment:记录编号" json:"id"`
	UserID      int64     `gorm:"index;not null;comment:用户编号" json:"userId"`
	SpuID       int64     `gorm:"index;not null;comment:商品 SPU 编号" json:"spuId"`
	UserDeleted bool      `gorm:"column:user_deleted;default:false;comment:用户是否删除" json:"userDeleted"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   int64     `gorm:"softDelete:milli;default:0;comment:删除时间" json:"deletedAt"`
	Deleted     bool      `gorm:"default:false;comment:是否删除" json:"deleted"`
}

func (*ProductBrowseHistory) TableName() string {
	return "product_browse_history"
}
