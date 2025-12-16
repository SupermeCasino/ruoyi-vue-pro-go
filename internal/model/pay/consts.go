package pay

// PayTransfer 转账状态枚举
// 对齐 Java: PayTransferStatusEnum
const (
	PayTransferStatusWaiting    = 0  // 等待转账
	PayTransferStatusProcessing = 5  // 转账进行中
	PayTransferStatusSuccess    = 10 // 转账成功
	PayTransferStatusClosed     = 20 // 转账关闭

	// PayTransferType 转账类型
	PayTransferTypeAlipayBalance = 1 // 支付宝 - 余额
	PayTransferTypeWxBalance     = 2 // 微信 - 余额
	PayTransferTypeBankCard      = 3 // 银行卡
	PayTransferTypeWallet        = 4 // 钱包余额
)

// IsPayTransferStatusSuccess 判断是否转账成功
func IsPayTransferStatusSuccess(status int) bool {
	return status == PayTransferStatusSuccess
}

// IsPayTransferStatusClosed 判断是否转账关闭
func IsPayTransferStatusClosed(status int) bool {
	return status == PayTransferStatusClosed
}

// IsPayTransferStatusWaiting 判断是否等待转账
func IsPayTransferStatusWaiting(status int) bool {
	return status == PayTransferStatusWaiting
}

// IsPayTransferStatusProcessing 判断是否转账进行中
func IsPayTransferStatusProcessing(status int) bool {
	return status == PayTransferStatusProcessing
}

// IsPayTransferStatusWaitingOrProcessing 判断是否处于待转账或转账中状态
func IsPayTransferStatusWaitingOrProcessing(status int) bool {
	return IsPayTransferStatusWaiting(status) || IsPayTransferStatusProcessing(status)
}

// IsPayTransferStatusSuccessOrClosed 判断是否处于成功或关闭状态
func IsPayTransferStatusSuccessOrClosed(status int) bool {
	return IsPayTransferStatusSuccess(status) || IsPayTransferStatusClosed(status)
}

// PayWalletBizType 钱包业务类型 (对齐 Java: PayWalletBizTypeEnum)
const (
	PayWalletBizTypeRecharge       = 1 // 充值
	PayWalletBizTypeRechargeRefund = 2 // 充值退款
	PayWalletBizTypePayment        = 3 // 支付
	PayWalletBizTypePaymentRefund  = 4 // 支付退款
	PayWalletBizTypeUpdateBalance  = 5 // 更新余额 (Admin)
)
