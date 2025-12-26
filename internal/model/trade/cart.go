package trade

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// Cart 购物车
type Cart struct {
	ID       int64         `gorm:"primaryKey;autoIncrement;comment:购物项编号" json:"id"`
	UserID   int64         `gorm:"index;not null;comment:用户编号" json:"userId"`
	SpuID    int64         `gorm:"index;not null;comment:商品 SPU 编号" json:"spuId"`
	SkuID    int64         `gorm:"index;not null;comment:商品 SKU 编号" json:"skuId"`
	Count    int           `gorm:"not null;default:1;comment:商品数量" json:"count"`
	Selected model.BitBool `gorm:"not null;comment:是否选中" json:"selected"`
	model.TenantBaseDO
}

func (*Cart) TableName() string {
	return "trade_cart"
}
