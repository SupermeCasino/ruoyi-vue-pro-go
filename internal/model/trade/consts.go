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
	// AfterSaleStatusProcessing 处理中
	AfterSaleStatusProcessing = 20
)

const (
	// BrokerageWithdrawStatusAuditing 审核中
	BrokerageWithdrawStatusAuditing = 0
	// BrokerageWithdrawStatusAuditSuccess 审核通过
	BrokerageWithdrawStatusAuditSuccess = 10
	// BrokerageWithdrawStatusAuditFail 审核不通过
	BrokerageWithdrawStatusAuditFail = 20
	// BrokerageWithdrawStatusWithdrawSuccess 提现成功
	BrokerageWithdrawStatusWithdrawSuccess = 11
	// BrokerageWithdrawStatusWithdrawFail 提现失败
	BrokerageWithdrawStatusWithdrawFail = 21
)

const (
	// BrokerageWithdrawTypeWallet 钱包
	BrokerageWithdrawTypeWallet = 1
	// BrokerageWithdrawTypeBank 银行卡
	BrokerageWithdrawTypeBank = 2
	// BrokerageWithdrawTypeWechat 微信 API
	BrokerageWithdrawTypeWechat = 3
	// BrokerageWithdrawTypeAlipay 支付宝 API
	BrokerageWithdrawTypeAlipay = 4
)

const (
	// BrokerageRecordBizTypeOrder 分销订单
	BrokerageRecordBizTypeOrder = 1
	// BrokerageRecordBizTypeWithdraw 佣金提现
	BrokerageRecordBizTypeWithdraw = 2
	// BrokerageRecordBizTypeWithdrawReject 提现驳回
	BrokerageRecordBizTypeWithdrawReject = 3
)

const (
	// BrokerageRecordStatusWait 待结算
	BrokerageRecordStatusWait = 0
	// BrokerageRecordStatusSettlement 已结算
	BrokerageRecordStatusSettlement = 1
	// BrokerageRecordStatusCancel 已取消
	BrokerageRecordStatusCancel = 2
)

const (
	// BrokerageUserLevelOne 一级
	BrokerageUserLevelOne = 1
	// BrokerageUserLevelTwo 二级
	BrokerageUserLevelTwo = 2
)

const (
	// PayTransferStatusWaiting 转账中
	PayTransferStatusWaiting = 0
	// PayTransferStatusSuccess 转账成功
	PayTransferStatusSuccess = 20
	// PayTransferStatusClosed 转账关闭
	PayTransferStatusClosed = 30
)
