package weixin

import (
	"backend-go/internal/service/pay/client"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func init() {
	client.RegisterCreator("wx_pub", NewWxPayClientAsClient)
	client.RegisterCreator("wx_lite", NewWxPayClientAsClient)
	client.RegisterCreator("wx_app", NewWxPayClientAsClient)
	client.RegisterCreator("wx_native", NewWxPayClientAsClient)
	client.RegisterCreator("wx_wap", NewWxPayClientAsClient)
	client.RegisterCreator("wx_bar", NewWxPayClientAsClient)
}

func NewWxPayClientAsClient(channelID int64, config string) (client.PayClient, error) {
	return NewWxPayClient(channelID, "wx_unknown", config)
}

type WxPayClient struct {
	*client.BaseClient
	config     *WxPayClientConfig
	coreClient *core.Client
	privateKey *rsa.PrivateKey
}

func NewWxPayClient(channelID int64, channelCode string, config string) (*WxPayClient, error) {
	return &WxPayClient{
		BaseClient: client.NewBaseClient(channelID, channelCode, config),
	}, nil
}

func (c *WxPayClient) Init() error {
	// 1. 解析配置
	var cfg WxPayClientConfig
	if err := json.Unmarshal([]byte(c.Config), &cfg); err != nil {
		return fmt.Errorf("解析微信支付配置失败: %w", err)
	}
	c.config = &cfg

	// 2. V3 版本初始化
	if cfg.APIVersion == APIVersionV3 {
		return c.initV3Client()
	}

	// V2 版本暂不支持
	return errors.New("暂不支持微信支付 V2 版本")
}

func (c *WxPayClient) initV3Client() error {
	cfg := c.config

	// 加载商户私钥
	privateKey, err := utils.LoadPrivateKey(cfg.PrivateKeyContent)
	if err != nil {
		return fmt.Errorf("加载商户私钥失败: %w", err)
	}
	c.privateKey = privateKey

	// 加载微信支付公钥
	publicKey, err := utils.LoadPublicKey(cfg.PublicKeyContent)
	if err != nil {
		return fmt.Errorf("加载微信支付公钥失败: %w", err)
	}

	// 使用微信支付公钥模式创建客户端
	opts := []core.ClientOption{
		option.WithWechatPayPublicKeyAuthCipher(
			cfg.MchID,
			cfg.CertSerialNo,
			privateKey,
			cfg.PublicKeyID,
			publicKey,
		),
	}

	coreClient, err := core.NewClient(context.Background(), opts...)
	if err != nil {
		return fmt.Errorf("创建微信支付客户端失败: %w", err)
	}
	c.coreClient = coreClient

	fmt.Printf("微信支付客户端初始化成功 [Channel: %d, MchID: %s]\n", c.ChannelID, cfg.MchID)
	return nil
}

// UnifiedOrder 统一下单
func (c *WxPayClient) UnifiedOrder(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
	// 根据渠道类型选择支付方式
	switch c.ChannelCode {
	case "wx_native":
		return c.nativeOrder(ctx, req)
	case "wx_pub", "wx_lite":
		return c.jsapiOrder(ctx, req)
	default:
		return nil, fmt.Errorf("暂不支持的微信支付渠道: %s", c.ChannelCode)
	}
}

// nativeOrder Native 扫码支付
func (c *WxPayClient) nativeOrder(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
	svc := native.NativeApiService{Client: c.coreClient}

	resp, result, err := svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(c.config.AppID),
		Mchid:       core.String(c.config.MchID),
		Description: core.String(req.Subject),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(req.NotifyURL),
		Amount: &native.Amount{
			Total:    core.Int64(int64(req.Price)),
			Currency: core.String("CNY"),
		},
	})

	if err != nil {
		return &client.OrderResp{
			Status:           20, // CLOSED
			OutTradeNo:       req.OutTradeNo,
			ChannelErrorCode: "NATIVE_PREPAY_ERROR",
			ChannelErrorMsg:  err.Error(),
		}, nil
	}

	_ = result

	return &client.OrderResp{
		Status:         0, // WAITING
		OutTradeNo:     req.OutTradeNo,
		DisplayMode:    "qr_code",
		DisplayContent: *resp.CodeUrl,
	}, nil
}

// jsapiOrder JSAPI 公众号/小程序支付
func (c *WxPayClient) jsapiOrder(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
	svc := jsapi.JsapiApiService{Client: c.coreClient}

	// OpenID 从 ChannelExtras 获取
	openID := ""
	if req.ChannelExtras != nil {
		openID = req.ChannelExtras["openid"]
	}
	if openID == "" {
		return nil, errors.New("JSAPI 支付需要 openid")
	}

	resp, result, err := svc.Prepay(ctx, jsapi.PrepayRequest{
		Appid:       core.String(c.config.AppID),
		Mchid:       core.String(c.config.MchID),
		Description: core.String(req.Subject),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(req.NotifyURL),
		Amount: &jsapi.Amount{
			Total:    core.Int64(int64(req.Price)),
			Currency: core.String("CNY"),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(openID),
		},
	})

	if err != nil {
		return &client.OrderResp{
			Status:           20, // CLOSED
			OutTradeNo:       req.OutTradeNo,
			ChannelErrorCode: "JSAPI_PREPAY_ERROR",
			ChannelErrorMsg:  err.Error(),
		}, nil
	}

	_ = result

	return &client.OrderResp{
		Status:         0, // WAITING
		OutTradeNo:     req.OutTradeNo,
		DisplayMode:    "app",
		DisplayContent: *resp.PrepayId,
	}, nil
}

// UnifiedRefund 统一退款
func (c *WxPayClient) UnifiedRefund(ctx context.Context, req *client.UnifiedRefundReq) (*client.RefundResp, error) {
	// TODO: 实现退款逻辑
	return nil, errors.New("退款功能暂未实现")
}

// GetOrder 查询订单
func (c *WxPayClient) GetOrder(ctx context.Context, outTradeNo string) (*client.OrderResp, error) {
	svc := jsapi.JsapiApiService{Client: c.coreClient}

	resp, _, err := svc.QueryOrderByOutTradeNo(ctx, jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(c.config.MchID),
	})

	if err != nil {
		return nil, err
	}

	status := 0 // WAITING
	var successTime time.Time
	if *resp.TradeState == "SUCCESS" {
		status = 10 // SUCCESS
		successTime, _ = time.Parse(time.RFC3339, *resp.SuccessTime)
	} else if *resp.TradeState == "CLOSED" || *resp.TradeState == "PAYERROR" {
		status = 20 // CLOSED
	}

	return &client.OrderResp{
		Status:         status,
		OutTradeNo:     outTradeNo,
		ChannelOrderNo: *resp.TransactionId,
		ChannelUserID:  *resp.Payer.Openid,
		SuccessTime:    successTime,
	}, nil
}

// GetRefund 查询退款
func (c *WxPayClient) GetRefund(ctx context.Context, outTradeNo, outRefundNo string) (*client.RefundResp, error) {
	return nil, errors.New("退款查询功能暂未实现")
}

// ParseOrderNotify 解析支付回调
func (c *WxPayClient) ParseOrderNotify(req *client.NotifyData) (*client.OrderResp, error) {
	// TODO: 使用 SDK 解密并验签
	return nil, errors.New("回调解析功能暂未实现")
}

// ParseRefundNotify 解析退款回调
func (c *WxPayClient) ParseRefundNotify(req *client.NotifyData) (*client.RefundResp, error) {
	return nil, errors.New("退款回调解析功能暂未实现")
}

// UnifiedTransfer 统一转账
func (c *WxPayClient) UnifiedTransfer(ctx context.Context, req *client.UnifiedTransferReq) (*client.TransferResp, error) {
	svc := transferbatch.TransferBatchApiService{Client: c.coreClient}

	// 创建转账批次
	batchReq := transferbatch.InitiateBatchTransferRequest{
		Appid:       core.String(c.config.AppID),
		OutBatchNo:  core.String(req.OutTradeNo),
		BatchName:   core.String(req.Subject),
		BatchRemark: core.String(req.Subject),
		TotalAmount: core.Int64(int64(req.Price)),
		TotalNum:    core.Int64(1),
		TransferDetailList: []transferbatch.TransferDetailInput{
			{
				OutDetailNo:    core.String(req.OutTradeNo + "_1"),
				TransferAmount: core.Int64(int64(req.Price)),
				TransferRemark: core.String(req.Subject),
				Openid:         core.String(req.ChannelUserID),
				UserName:       nil, // 如需要真实姓名需加密
			},
		},
	}

	resp, _, err := svc.InitiateBatchTransfer(ctx, batchReq)
	if err != nil {
		return &client.TransferResp{
			Status:           20, // CLOSED
			OutTradeNo:       req.OutTradeNo,
			ChannelErrorCode: "TRANSFER_ERROR",
			ChannelErrorMsg:  err.Error(),
		}, nil
	}

	return &client.TransferResp{
		Status:            10, // SUCCESS (or PROCESSING)
		OutTradeNo:        req.OutTradeNo,
		ChannelTransferNo: *resp.BatchId,
	}, nil
}
