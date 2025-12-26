package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// ProductSku 商品 SKU
type ProductSku struct {
	ID                   int64                `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	SpuID                int64                `gorm:"column:spu_id;not null;comment:SPU编号" json:"spuId"`
	Properties           []ProductSkuProperty `gorm:"type:json;serializer:json;comment:属性数组" json:"properties"`
	Price                int                  `gorm:"default:0;comment:商品价格" json:"price"`      // 单位：分
	MarketPrice          int                  `gorm:"default:0;comment:市场价" json:"marketPrice"` // 单位：分
	CostPrice            int                  `gorm:"default:0;comment:成本价" json:"costPrice"`   // 单位：分
	BarCode              string               `gorm:"size:64;default:'';comment:商品条码" json:"barCode"`
	PicURL               string               `gorm:"column:pic_url;size:255;default:'';comment:图片地址" json:"picUrl"`
	Stock                int                  `gorm:"default:0;comment:库存" json:"stock"`
	Weight               float64              `gorm:"default:0;comment:商品重量" json:"weight"` // 单位：kg
	Volume               float64              `gorm:"default:0;comment:商品体积" json:"volume"` // 单位：m^3
	FirstBrokeragePrice  int                  `gorm:"default:0;comment:一级分销的佣金" json:"firstBrokeragePrice"`
	SecondBrokeragePrice int                  `gorm:"default:0;comment:二级分销的佣金" json:"secondBrokeragePrice"`
	SalesCount           int                  `gorm:"default:0;comment:商品销量" json:"salesCount"`
	model.TenantBaseDO
}

type ProductSkuProperty struct {
	PropertyID   int64  `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	ValueID      int64  `json:"valueId"`
	ValueName    string `json:"valueName"`
}

func (ProductSku) TableName() string {
	return "product_sku"
}
