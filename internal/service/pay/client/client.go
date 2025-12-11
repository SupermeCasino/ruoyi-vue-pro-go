package client

import "context"

// PayClient 支付客户端接口
type PayClient interface {
	// GetID 获得渠道编号
	GetID() int64

	// Init 初始化
	Init() error

	// UnifiedOrder 调用支付渠道，统一下单
	UnifiedOrder(ctx context.Context, req *UnifiedOrderReq) (*OrderResp, error)

	// UnifiedRefund 调用支付渠道，进行退款
	UnifiedRefund(ctx context.Context, req *UnifiedRefundReq) (*RefundResp, error)

	// GetOrder 获得支付订单信息
	GetOrder(ctx context.Context, outTradeNo string) (*OrderResp, error)

	// GetRefund 获得退款订单信息
	GetRefund(ctx context.Context, outTradeNo, outRefundNo string) (*RefundResp, error)

	// ParseOrderNotify 解析 order 回调数据
	ParseOrderNotify(req *NotifyData) (*OrderResp, error)

	// ParseRefundNotify 解析 refund 回调数据
	ParseRefundNotify(req *NotifyData) (*RefundResp, error)

	// UnifiedTransfer 调用支付渠道，进行转账
	UnifiedTransfer(ctx context.Context, req *UnifiedTransferReq) (*TransferResp, error)
}

type NotifyData struct {
	Params  map[string]string
	Body    string
	Headers map[string]string
}
