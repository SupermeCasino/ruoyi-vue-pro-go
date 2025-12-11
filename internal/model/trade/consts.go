package trade

const (
	// TradeOrderStatusUnpaid 待支付
	TradeOrderStatusUnpaid = 0
	// TradeOrderStatusUndelivered 待发货
	TradeOrderStatusUndelivered = 10
	// TradeOrderStatusDelivered 待收货
	TradeOrderStatusDelivered = 20
	// TradeOrderStatusCompleted 完成
	TradeOrderStatusCompleted = 30
	// TradeOrderStatusCanceled 取消
	TradeOrderStatusCanceled = 40
)

const (
	// DeliveryTypeExpress 快递发货
	DeliveryTypeExpress = 1
	// DeliveryTypePickUp 到店自提
	DeliveryTypePickUp = 2
)

const (
	// AfterSaleStatusApply 申请售后
	AfterSaleStatusApply = 10
)

const (
	// BrokerageWithdrawStatusAuditing 审核中
	BrokerageWithdrawStatusAuditing = 1
)
