package alipay

import (
	"backend-go/internal/service/pay/client"
	"context"
	"errors"
	"fmt"
)

func init() {
	client.RegisterCreator("alipay_pc", NewAlipayPayClientAsClient)
	client.RegisterCreator("alipay_wap", NewAlipayPayClientAsClient)
	client.RegisterCreator("alipay_app", NewAlipayPayClientAsClient)
	client.RegisterCreator("alipay_qr", NewAlipayPayClientAsClient)
	client.RegisterCreator("alipay_bar", NewAlipayPayClientAsClient)
}

func NewAlipayPayClientAsClient(channelID int64, config string) (client.PayClient, error) {
	// We pass empty channelCode because it's determined by the map key in factory if we want precise control.
	// Or we can parse it from somewhere. For now, let's just make it work.
	return NewAlipayPayClient(channelID, "alipay_unknown", config)
}

type AlipayPayClient struct {
	*client.BaseClient
}

func NewAlipayPayClient(channelID int64, channelCode string, config string) (*AlipayPayClient, error) {
	return &AlipayPayClient{
		BaseClient: client.NewBaseClient(channelID, channelCode, config),
	}, nil
}

func (c *AlipayPayClient) Init() error {
	fmt.Printf("Initializing Alipay Client for Channel %d\n", c.ChannelID)
	// Parse config JSON to Alipay Config struct
	return nil
}

func (c *AlipayPayClient) UnifiedOrder(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
	fmt.Printf("Alipay UnifiedOrder: %s\n", req.OutTradeNo)
	// Implement real Alipay SDK call here
	return &client.OrderResp{
		Status:         0, // WAITING
		OutTradeNo:     req.OutTradeNo,
		DisplayMode:    req.DisplayMode,
		DisplayContent: "https://openapi.alipay.com/gateway.do?mock=" + req.OutTradeNo,
	}, nil
}

func (c *AlipayPayClient) UnifiedRefund(ctx context.Context, req *client.UnifiedRefundReq) (*client.RefundResp, error) {
	return nil, errors.New("not implemented")
}

func (c *AlipayPayClient) GetOrder(ctx context.Context, outTradeNo string) (*client.OrderResp, error) {
	return nil, errors.New("not implemented")
}

func (c *AlipayPayClient) GetRefund(ctx context.Context, outTradeNo, outRefundNo string) (*client.RefundResp, error) {
	return nil, errors.New("not implemented")
}

func (c *AlipayPayClient) ParseOrderNotify(req *client.NotifyData) (*client.OrderResp, error) {
	return nil, errors.New("not implemented")
}

func (c *AlipayPayClient) ParseRefundNotify(req *client.NotifyData) (*client.RefundResp, error) {
	return nil, errors.New("not implemented")
}

func (c *AlipayPayClient) UnifiedTransfer(ctx context.Context, req *client.UnifiedTransferReq) (*client.TransferResp, error) {
	fmt.Printf("Alipay UnifiedTransfer: %s -> %s\n", req.OutTradeNo, req.ChannelUserID)
	return &client.TransferResp{
		ChannelTransferNo: "MOCK_ALIPAY_TRANSFER_" + req.OutTradeNo,
		Status:            10, // SUCCESS
	}, nil
}
