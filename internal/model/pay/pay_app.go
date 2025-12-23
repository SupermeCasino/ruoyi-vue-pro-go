package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayApp 支付应用 DO
type PayApp struct {
	ID                int64  `gorm:"primaryKey;autoIncrement;comment:应用编号" json:"id"`
	AppKey            string `gorm:"size:64;not null;comment:应用标识" json:"appKey"`
	Name              string `gorm:"size:64;not null;comment:应用名" json:"name"`
	Status            int    `gorm:"default:0;not null;comment:状态" json:"status"` // 参见 CommonStatusEnum
	Remark            string `gorm:"size:255;default:'';comment:备注" json:"remark"`
	OrderNotifyURL    string `gorm:"column:order_notify_url;size:1024;not null;comment:支付结果的回调地址" json:"orderNotifyUrl"`
	RefundNotifyURL   string `gorm:"column:refund_notify_url;size:1024;not null;comment:退款结果的回调地址" json:"refundNotifyUrl"`
	TransferNotifyURL string `gorm:"column:transfer_notify_url;size:1024;default:'';comment:转账结果的回调地址" json:"transferNotifyUrl"`
	model.TenantBaseDO
}

func (PayApp) TableName() string {
	return "pay_app"
}
