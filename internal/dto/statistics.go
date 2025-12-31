package dto

import "time"

// TradeStatisticsDTO 交易统计数据传输对象
type TradeStatisticsDTO struct {
	StatisticsTime           time.Time // 统计日期
	OrderCreateCount         int       // 创建订单数
	OrderPayCount            int       // 支付订单商品数
	OrderPayPrice            int       // 总支付金额(分)
	AfterSaleCount           int       // 退款订单数
	AfterSaleRefundPrice     int       // 总退款金额(分)
	BrokerageSettlementPrice int       // 佣金金额已结算(分)
	WalletPayPrice           int       // 总支付金额余额(分)
	RechargePayCount         int       // 充值订单数
	RechargePayPrice         int       // 充值金额(分)
	RechargeRefundCount      int       // 充值退款订单数
	RechargeRefundPrice      int       // 充值退款金额(分)
}
