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

// ============= 订单类型常量 (Order Type Constants) =============

const (
	// TradeOrderTypeNormal 普通订单
	TradeOrderTypeNormal = 0
	// TradeOrderTypeSeckill 秒杀订单
	TradeOrderTypeSeckill = 1
	// TradeOrderTypeBargain 砍价订单
	TradeOrderTypeBargain = 2
	// TradeOrderTypeCombination 拼团订单
	TradeOrderTypeCombination = 3
	// TradeOrderTypePoint 积分订单
	TradeOrderTypePoint = 4
)

// ============= 价格计算器优先级常量 (Price Calculator Priority) =============
// 数字越小优先级越高，按照Java版本的TradePriceCalculator优先顺序定义

const (
	// OrderSeckillActivity 秒杀活动计算器优先级
	OrderSeckillActivity = 8
	// OrderBargainActivity 砍价活动计算器优先级
	OrderBargainActivity = 8
	// OrderCombinationActivity 拼团活动计算器优先级
	OrderCombinationActivity = 8
	// OrderPointActivity 积分商城活动计算器优先级
	OrderPointActivity = 8
	// OrderDiscountActivity 限时折扣活动计算器优先级
	OrderDiscountActivity = 10
	// OrderRewardActivity 满减送活动计算器优先级
	OrderRewardActivity = 20
	// OrderCoupon 优惠券计算器优先级
	OrderCoupon = 30
	// OrderPointUse 积分抵扣计算器优先级
	OrderPointUse = 40
	// OrderDelivery 运费计算器优先级
	OrderDelivery = 50
	// OrderPointGive 积分赠送计算器优先级
	OrderPointGive = 999
)

// ============= 促销类型常量 (Promotion Type Constants) =============

const (
	// PromotionTypeNone 无促销
	PromotionTypeNone = 0
	// PromotionTypeDiscountActivity 限时折扣活动
	PromotionTypeDiscountActivity = 10
	// PromotionTypeRewardActivity 满减送活动
	PromotionTypeRewardActivity = 20
	// PromotionTypeMemberLevel 会员等级折扣
	PromotionTypeMemberLevel = 30
	// PromotionTypeCoupon 优惠券
	PromotionTypeCoupon = 40
	// PromotionTypeCombination 拼团活动
	PromotionTypeCombination = 50
	// PromotionTypeBargain 砍价活动
	PromotionTypeBargain = 60
	// PromotionTypeSeckill 秒杀活动
	PromotionTypeSeckill = 70
	// PromotionTypePoint 积分抵扣
	PromotionTypePoint = 80
)

// ============= 折扣类型常量 (Discount Type Constants) =============

const (
	// DiscountTypePrice 减价
	DiscountTypePrice = 1
	// DiscountTypePercent 打折
	DiscountTypePercent = 2
)

// ============= 计算器名称常量 (Calculator Name Constants) =============

const (
	// CalculatorNameSeckill 秒杀活动价格计算器
	CalculatorNameSeckill = "秒杀活动价格计算器"
	// CalculatorNameBargain 砍价活动价格计算器
	CalculatorNameBargain = "砍价活动价格计算器"
	// CalculatorNameCombination 拼团活动价格计算器
	CalculatorNameCombination = "拼团活动价格计算器"
	// CalculatorNamePoint 积分商城价格计算器
	CalculatorNamePoint = "积分商城价格计算器"
	// CalculatorNameDiscount 限时折扣活动价格计算器
	CalculatorNameDiscount = "限时折扣活动价格计算器"
	// CalculatorNameReward 满减送活动价格计算器
	CalculatorNameReward = "满减送活动价格计算器"
	// CalculatorNameCoupon 优惠券价格计算器
	CalculatorNameCoupon = "优惠券价格计算器"
	// CalculatorNamePointUse 积分抵扣价格计算器
	CalculatorNamePointUse = "积分抵扣价格计算器"
	// CalculatorNameDelivery 运费计算器
	CalculatorNameDelivery = "运费计算器"
	// CalculatorNamePointGive 积分赠送计算器
	CalculatorNamePointGive = "积分赠送计算器"
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
	OrderOperateTypeCancel = 41
	// OrderOperateTypeRefund 退款
	OrderOperateTypeRefund = 60
)

const (
	// OrderCancelTypeMember 会员取消
	OrderCancelTypeMember = 10
	// OrderCancelTypeTimeout 支付超时取消
	OrderCancelTypeTimeout = 20
	// OrderCancelTypeAdmin 管理员取消
	OrderCancelTypeAdmin = 30
	// OrderCancelTypeSystem 系统取消
	OrderCancelTypeSystem = 40
	// OrderCancelTypeAfterSaleClose 售后全退关闭
	OrderCancelTypeAfterSaleClose = 50
	// OrderCancelTypePaymentFallback 支付异常回滚
	OrderCancelTypePaymentFallback = 60
	// OrderCancelTypeCombinationClose 拼团关闭取消
	OrderCancelTypeCombinationClose = 70
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
