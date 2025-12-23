package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductCategory 商品分类
type ProductCategory struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;comment:分类编号" json:"id"`
	ParentID    int64  `gorm:"column:parent_id;not null;default:0;comment:父分类编号" json:"parentId"`
	Name        string `gorm:"size:255;not null;comment:分类名称" json:"name"`
	PicURL      string `gorm:"column:pic_url;size:255;default:'';comment:移动端分类图" json:"picUrl"`
	Sort        int32  `gorm:"default:0;comment:分类排序" json:"sort"`
	Status      int32  `gorm:"default:0;comment:开启状态" json:"status"` // 参见 CommonStatusEnum
	Description string `gorm:"size:512;default:'';comment:分类描述" json:"description"`

	model.TenantBaseDO
}

func (ProductCategory) TableName() string {
	return "product_category"
}
