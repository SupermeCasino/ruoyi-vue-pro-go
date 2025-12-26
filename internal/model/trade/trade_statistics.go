package trade

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

type TradeStatistics struct {
	ID                       int64     `gorm:"column:id;primaryKey;autoIncrement;comment:编号" json:"id"`
	Time                     time.Time `gorm:"column:time;not null;comment:统计日期" json:"time"`
	OrderCreateCount         int       `gorm:"column:order_create_count;default:0;comment:创建订单数" json:"orderCreateCount"`
	OrderPayCount            int       `gorm:"column:order_pay_count;default:0;comment:支付订单商品数" json:"orderPayCount"`
	OrderPayPrice            int       `gorm:"column:order_pay_price;default:0;comment:总支付金额(分)" json:"orderPayPrice"`
	AfterSaleCount           int       `gorm:"column:after_sale_count;default:0;comment:退款订单数" json:"afterSaleCount"`
	AfterSaleRefundPrice     int       `gorm:"column:after_sale_refund_price;default:0;comment:总退款金额(分)" json:"afterSaleRefundPrice"`
	BrokerageSettlementPrice int       `gorm:"column:brokerage_settlement_price;default:0;comment:佣金金额已结算(分)" json:"brokerageSettlementPrice"`
	WalletPayPrice           int       `gorm:"column:wallet_pay_price;default:0;comment:总支付金额余额(分)" json:"walletPayPrice"`
	RechargePayCount         int       `gorm:"column:recharge_pay_count;default:0;comment:充值订单数" json:"rechargePayCount"`
	RechargePayPrice         int       `gorm:"column:recharge_pay_price;default:0;comment:充值金额(分)" json:"rechargePayPrice"`
	RechargeRefundCount      int       `gorm:"column:recharge_refund_count;default:0;comment:充值退款订单数" json:"rechargeRefundCount"`
	RechargeRefundPrice      int       `gorm:"column:recharge_refund_price;default:0;comment:充值退款金额(分)" json:"rechargeRefundPrice"`
	model.TenantBaseDO
}

func (TradeStatistics) TableName() string {
	return "trade_statistics"
}
