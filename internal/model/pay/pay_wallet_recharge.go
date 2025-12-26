package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayWalletRecharge 会员钱包充值表
type PayWalletRecharge struct {
	ID               int64      `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	WalletID         int64      `gorm:"column:wallet_id;not null;comment:钱包编号" json:"walletId"`
	TotalPrice       int        `gorm:"column:total_price;not null;default:0;comment:充值实际金额" json:"totalPrice"` // 单位：分
	PayPrice         int        `gorm:"column:pay_price;not null;default:0;comment:实际支付金额" json:"payPrice"`     // 单位：分
	BonusPrice       int        `gorm:"column:bonus_price;not null;default:0;comment:赠送金额" json:"bonusPrice"`   // 单位：分
	PackageID        int64      `gorm:"column:package_id;default:0;comment:充值套餐编号" json:"packageId"`
	PayStatus        bool       `gorm:"column:pay_status;not null;default:0;comment:支付状态" json:"payStatus"`
	PayOrderID       int64      `gorm:"column:pay_order_id;comment:支付订单编号" json:"payOrderId"`
	PayChannelCode   string     `gorm:"column:pay_channel_code;size:32;comment:支付渠道" json:"payChannelCode"`
	PayTime          *time.Time `gorm:"column:pay_time;comment:订单支付时间" json:"payTime"`
	RefundStatus     int        `gorm:"column:refund_status;not null;default:0;comment:退款状态" json:"refundStatus"` // 0: 无, 10: 退款中, 20: 成功, 30: 失败
	PayRefundID      int64      `gorm:"column:pay_refund_id;default:0;comment:支付退款单编号" json:"payRefundId"`
	RefundTotalPrice int        `gorm:"column:refund_total_price;not null;default:0;comment:退款金额" json:"refundTotalPrice"`   // 单位：分
	RefundPayPrice   int        `gorm:"column:refund_pay_price;not null;default:0;comment:退款支付金额" json:"refundPayPrice"`     // 单位：分
	RefundBonusPrice int        `gorm:"column:refund_bonus_price;not null;default:0;comment:退款赠送金额" json:"refundBonusPrice"` // 单位：分
	RefundTime       *time.Time `gorm:"column:refund_time;comment:退款生效时间" json:"refundTime"`
	model.TenantBaseDO
}

func (PayWalletRecharge) TableName() string {
	return "pay_wallet_recharge"
}
