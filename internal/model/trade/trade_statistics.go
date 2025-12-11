package trade

import (
	"time"
)

// TradeStatistics 交易统计 - 严格对齐 Java TradeStatisticsDO
// 表名: trade_statistics
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
	Creator                  string    `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater                  string    `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt                time.Time `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt                time.Time `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted                  bool      `gorm:"column:deleted;default:0;comment:是否删除" json:"deleted"`
}

func (TradeStatistics) TableName() string {
	return "trade_statistics"
}
