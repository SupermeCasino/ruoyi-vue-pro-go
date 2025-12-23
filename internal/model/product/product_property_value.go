package product

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductPropertyValue 商品属性值
type ProductPropertyValue struct {
	ID         int64         `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	PropertyID int64         `gorm:"column:property_id;not null;comment:属性项的编号" json:"propertyId"`
	Name       string        `gorm:"size:255;not null;comment:名称" json:"name"`
	Remark     string        `gorm:"size:500;default:'';comment:备注" json:"remark"`
	Creator    string        `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string        `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreateTime  time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdateTime  time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID   int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (ProductPropertyValue) TableName() string {
	return "product_property_value"
}
