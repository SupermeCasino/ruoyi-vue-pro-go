package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayOrder 支付订单
type PayOrder struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement"`
	AppID           int64      `gorm:"column:app_id"`
	ChannelID       int64      `gorm:"column:channel_id"`
	ChannelCode     string     `gorm:"column:channel_code"`
	MerchantOrderId string     `gorm:"column:merchant_order_id"`
	Subject         string     `gorm:"column:subject"`
	Body            string     `gorm:"column:body"`
	NotifyURL       string     `gorm:"column:notify_url"`
	Price           int        `gorm:"column:price"` // Unit: fen
	ChannelFeeRate  float64    `gorm:"column:channel_fee_rate"`
	ChannelFeePrice int        `gorm:"column:channel_fee_price"` // Unit: fen
	Status          int        `gorm:"column:status"`
	UserIP          string     `gorm:"column:user_ip"`
	ExpireTime      time.Time  `gorm:"column:expire_time"`
	SuccessTime     *time.Time `gorm:"column:success_time"`
	ExtensionID     int64      `gorm:"column:extension_id"`
	No              string     `gorm:"column:no"`
	RefundPrice     int        `gorm:"column:refund_price"` // Unit: fen
	ChannelUserID   string     `gorm:"column:channel_user_id"`
	ChannelOrderNo  string     `gorm:"column:channel_order_no"`
	model.TenantBaseDO
}

func (PayOrder) TableName() string {
	return "pay_order"
}
