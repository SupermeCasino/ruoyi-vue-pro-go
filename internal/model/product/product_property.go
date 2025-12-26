package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductProperty 商品属性项
type ProductProperty struct {
	ID     int64  `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name   string `gorm:"size:255;not null;comment:名称" json:"name"`
	Remark string `gorm:"size:500;default:'';comment:备注" json:"remark"`
	model.TenantBaseDO
}

func (ProductProperty) TableName() string {
	return "product_property"
}
