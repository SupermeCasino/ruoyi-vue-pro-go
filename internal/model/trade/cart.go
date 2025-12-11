package trade

import "time"

// Cart 购物车
type Cart struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:购物项编号" json:"id"`
	UserID    int64     `gorm:"index;not null;comment:用户编号" json:"userId"`
	SpuID     int64     `gorm:"index;not null;comment:商品 SPU 编号" json:"spuId"`
	SkuID     int64     `gorm:"index;not null;comment:商品 SKU 编号" json:"skuId"`
	Count     int       `gorm:"not null;default:1;comment:商品数量" json:"count"`
	Selected  bool      `gorm:"not null;default:true;comment:是否选中" json:"selected"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt int64     `gorm:"softDelete:milli;default:0;comment:删除时间" json:"deletedAt"`
	Deleted   bool      `gorm:"default:false;comment:是否删除" json:"deleted"`
}

func (*Cart) TableName() string {
	return "trade_cart"
}
