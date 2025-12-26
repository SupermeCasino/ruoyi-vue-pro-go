package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductFavorite 商品收藏 DO
type ProductFavorite struct {
	ID     int64 `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	UserID int64 `gorm:"index;not null;comment:用户编号" json:"userId"`
	SpuID  int64 `gorm:"index;not null;comment:商品 SPU 编号" json:"spuId"`
	model.TenantBaseDO
}

func (ProductFavorite) TableName() string {
	return "product_favorite"
}
