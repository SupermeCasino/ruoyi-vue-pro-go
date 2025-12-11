package product

import (
	"time"
)

// ProductStatistics 商品统计 - 严格对齐 Java ProductStatisticsDO
// 表名: product_statistics
type ProductStatistics struct {
	ID                   int64     `gorm:"column:id;primaryKey;autoIncrement;comment:编号" json:"id"`
	Time                 time.Time `gorm:"column:time;not null;comment:统计日期" json:"time"`
	SpuID                int64     `gorm:"column:spu_id;not null;comment:商品SPU编号" json:"spuId"`
	BrowseCount          int       `gorm:"column:browse_count;default:0;comment:浏览量" json:"browseCount"`
	BrowseUserCount      int       `gorm:"column:browse_user_count;default:0;comment:访客量" json:"browseUserCount"`
	FavoriteCount        int       `gorm:"column:favorite_count;default:0;comment:收藏数量" json:"favoriteCount"`
	CartCount            int       `gorm:"column:cart_count;default:0;comment:加购数量" json:"cartCount"`
	OrderCount           int       `gorm:"column:order_count;default:0;comment:下单件数" json:"orderCount"`
	OrderPayCount        int       `gorm:"column:order_pay_count;default:0;comment:支付件数" json:"orderPayCount"`
	OrderPayPrice        int       `gorm:"column:order_pay_price;default:0;comment:支付金额(分)" json:"orderPayPrice"`
	AfterSaleCount       int       `gorm:"column:after_sale_count;default:0;comment:退款件数" json:"afterSaleCount"`
	AfterSaleRefundPrice int       `gorm:"column:after_sale_refund_price;default:0;comment:退款金额(分)" json:"afterSaleRefundPrice"`
	BrowseConvertPercent int       `gorm:"column:browse_convert_percent;default:0;comment:访客支付转化率(百分比)" json:"browseConvertPercent"`
	Creator              string    `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater              string    `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt            time.Time `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt            time.Time `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted              bool      `gorm:"column:deleted;default:0;comment:是否删除" json:"deleted"`
}

func (ProductStatistics) TableName() string {
	return "product_statistics"
}
