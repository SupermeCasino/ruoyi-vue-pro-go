package pay

// PayOrderStatusEnum 支付订单状态
const (
	PayOrderStatusWaiting = 0  // 待支付
	PayOrderStatusSuccess = 10 // 支付成功
	PayOrderStatusClosed  = 20 // 支付关闭
	PayOrderStatusRefund  = 30 // 已退款
)

// PayNotifyTypeEnum 支付通知类型
const (
	PayNotifyTypeOrder  = 1 // 支付单
	PayNotifyTypeRefund = 2 // 退款单
)

// PayNotifyStatusEnum 支付通知状态
const (
	PayNotifyStatusWaiting        = 0  // 等待通知
	PayNotifyStatusSuccess        = 10 // 通知成功
	PayNotifyStatusFailure        = 20 // 通知失败 (多次尝试，彻底失败)
	PayNotifyStatusRequestSuccess = 21 // 请求成功，但是结果失败
	PayNotifyStatusRequestFailure = 22 // 请求失败
)
