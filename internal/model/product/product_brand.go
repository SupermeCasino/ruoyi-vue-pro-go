package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductBrand 商品品牌
type ProductBrand struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name        string `gorm:"size:255;not null;comment:品牌名称" json:"name"`
	PicURL      string `gorm:"column:pic_url;size:255;default:'';comment:品牌图片" json:"picUrl"`
	Sort        int    `gorm:"default:0;comment:品牌排序" json:"sort"`
	Description string `gorm:"size:1024;default:'';comment:品牌描述" json:"description"`
	Status      int    `gorm:"default:0;comment:状态" json:"status"` // 0: 开启, 1: 关闭
	model.TenantBaseDO
}

func (ProductBrand) TableName() string {
	return "product_brand"
}
