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
	// OrderTypeNormal 普通订单
	OrderTypeNormal = 1
)

const (
	// OrderOperateTypeCreate 创建订单
	OrderOperateTypeCreate = 10
	// OrderOperateTypePay 支付
	OrderOperateTypePay = 20
	// OrderOperateTypeDelivery 发货
	OrderOperateTypeDelivery = 30
	// OrderOperateTypeReceive 确认收货
	OrderOperateTypeReceive = 40
	// OrderOperateTypePickUp 自提核销
	OrderOperateTypePickUp = 50
	// OrderOperateTypeCancel 取消订单
	OrderOperateTypeCancel = 40
	// OrderOperateTypeRefund 退款
	OrderOperateTypeRefund = 60
)

const (
	// OrderCancelTypeUser 用户取消
	OrderCancelTypeUser = 1
	// OrderCancelTypeAdmin 管理员取消
	OrderCancelTypeAdmin = 2
	// OrderCancelTypeSystem 系统取消
	OrderCancelTypeSystem = 3
)

const (
	// OrderRefundStatusNone 无退款
	OrderRefundStatusNone = 0
	// OrderRefundStatusApply 申请退款
	OrderRefundStatusApply = 10
	// OrderRefundStatusAuditing 审核中
	OrderRefundStatusAuditing = 20
	// OrderRefundStatusRefunded 已退款
	OrderRefundStatusRefunded = 30
)

const (
	// PickUpVerifyCodeLength 核销码长度
	PickUpVerifyCodeLength = 8
)

const (
	// DeliveryStatusEnabled 启用状态
	DeliveryStatusEnabled = 1
	// DeliveryStatusDisabled 禁用状态
	DeliveryStatusDisabled = 0
)

// ============= Brokerage Constants (保留原有常量) =============

const (
	// AfterSaleStatusNone 无售后申请
	AfterSaleStatusNone = 0
	// AfterSaleStatusApply 申请售后
	AfterSaleStatusApply = 10
	// AfterSaleStatusSellerAgree 卖家同意，等待买家退货
	AfterSaleStatusSellerAgree = 20
	// AfterSaleStatusSellerDisagree 卖家拒绝
	AfterSaleStatusSellerDisagree = 30
	// AfterSaleStatusBuyerDelivery 买家已发货，等待卖家收货
	AfterSaleStatusBuyerDelivery = 40
	// AfterSaleStatusWaitRefund 卖家已收货，等待平台退款
	AfterSaleStatusWaitRefund = 50
	// AfterSaleStatusSellerRefuse 卖家拒绝收货
	AfterSaleStatusSellerRefuse = 61
	// AfterSaleStatusComplete 完成
	AfterSaleStatusComplete = 100
	// AfterSaleStatusBuyerCancel 买家取消
	AfterSaleStatusBuyerCancel = 60
)

const (
	// AfterSaleWayRefund 仅退款
	AfterSaleWayRefund = 10
	// AfterSaleWayReturnAndRefund 退货退款
	AfterSaleWayReturnAndRefund = 20
)

const (
	// AfterSaleTypeInSale 售中
	AfterSaleTypeInSale = 10
	// AfterSaleTypeAfterSale 售后
	AfterSaleTypeAfterSale = 20
)

const (
	// AfterSaleOperateTypeMemberCreate 会员申请售后
	AfterSaleOperateTypeMemberCreate = 10
	// AfterSaleOperateTypeAdminAgreeApply 管理员同意售后
	AfterSaleOperateTypeAdminAgreeApply = 20
	// AfterSaleOperateTypeAdminDisagreeApply 管理员拒绝售后
	AfterSaleOperateTypeAdminDisagreeApply = 21
	// AfterSaleOperateTypeMemberDelivery 会员退货
	AfterSaleOperateTypeMemberDelivery = 30
	// AfterSaleOperateTypeAdminAgreeReceive 管理员确认收货
	AfterSaleOperateTypeAdminAgreeReceive = 40
	// AfterSaleOperateTypeAdminDisagreeReceive 管理员拒绝收货
	AfterSaleOperateTypeAdminDisagreeReceive = 41
	// AfterSaleOperateTypeAdminRefund 管理员确认退款
	AfterSaleOperateTypeAdminRefund = 50
	// AfterSaleOperateTypeSystemRefundSuccess 系统确认退款成功
	AfterSaleOperateTypeSystemRefundSuccess = 51
	// AfterSaleOperateTypeSystemRefundFail 系统确认退款失败
	AfterSaleOperateTypeSystemRefundFail = 52
	// AfterSaleOperateTypeMemberCancel 会员取消售后
	AfterSaleOperateTypeMemberCancel = 60
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
