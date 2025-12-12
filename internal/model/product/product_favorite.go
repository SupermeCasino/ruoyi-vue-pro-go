package product

import "time"

// ProductFavorite 商品收藏 DO
type ProductFavorite struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	UserID    int64     `gorm:"index;not null;comment:用户编号" json:"userId"`
	SpuID     int64     `gorm:"index;not null;comment:商品 SPU 编号" json:"spuId"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	Deleted   bool      `gorm:"default:false;comment:是否删除" json:"deleted"` // 逻辑删除标志
}

func (ProductFavorite) TableName() string {
	return "product_favorite"
}
