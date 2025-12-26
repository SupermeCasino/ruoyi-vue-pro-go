package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayOrderExtension 支付订单拓展
type PayOrderExtension struct {
	ID                int64  `gorm:"column:id;primaryKey;autoIncrement"`
	No                string `gorm:"column:no"`
	OrderID           int64  `gorm:"column:order_id"`
	ChannelID         int64  `gorm:"column:channel_id"`
	ChannelCode       string `gorm:"column:channel_code"`
	UserIP            string `gorm:"column:user_ip"`
	Status            int    `gorm:"column:status"`
	ChannelExtras     string `gorm:"column:channel_extras"` // JSON String
	ChannelErrorCode  string `gorm:"column:channel_error_code"`
	ChannelErrorMsg   string `gorm:"column:channel_error_msg"`
	ChannelNotifyData string `gorm:"column:channel_notify_data"`
	model.TenantBaseDO
}

func (PayOrderExtension) TableName() string {
	return "pay_order_extension"
}
