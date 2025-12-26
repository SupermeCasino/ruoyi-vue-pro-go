package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductPropertyValue 商品属性值
type ProductPropertyValue struct {
	ID         int64  `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	PropertyID int64  `gorm:"column:property_id;not null;comment:属性项的编号" json:"propertyId"`
	Name       string `gorm:"size:255;not null;comment:名称" json:"name"`
	Remark     string `gorm:"size:500;default:'';comment:备注" json:"remark"`
	model.TenantBaseDO
}

func (ProductPropertyValue) TableName() string {
	return "product_property_value"
}
