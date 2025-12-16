package weixin

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client"

	"io"
	"net/http"
	"strings"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
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
	publicKey  *rsa.PublicKey
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
	c.publicKey = publicKey

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
	case "wx_h5":
		return c.h5Order(ctx, req)
	case "wx_app":
		return c.appOrder(ctx, req)
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

// h5Order H5 支付
func (c *WxPayClient) h5Order(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
	svc := h5.H5ApiService{Client: c.coreClient}

	// 构造 H5 场景信息
	sceneInfo := &h5.SceneInfo{
		PayerClientIp: core.String(req.UserIP),
		H5Info: &h5.H5Info{
			Type: core.String("Wap"),
		},
	}
	// 如果 channelExtras 中有 info，可以解析覆盖默认值
	if req.ChannelExtras != nil {
		if appName, ok := req.ChannelExtras["app_name"]; ok {
			sceneInfo.H5Info.AppName = core.String(appName)
		}
		if bundleId, ok := req.ChannelExtras["bundle_id"]; ok {
			sceneInfo.H5Info.BundleId = core.String(bundleId)
		}
	}

	resp, result, err := svc.Prepay(ctx, h5.PrepayRequest{
		Appid:       core.String(c.config.AppID),
		Mchid:       core.String(c.config.MchID),
		Description: core.String(req.Subject),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(req.NotifyURL),
		Amount: &h5.Amount{
			Total:    core.Int64(int64(req.Price)),
			Currency: core.String("CNY"),
		},
		SceneInfo: sceneInfo,
	})

	if err != nil {
		return &client.OrderResp{
			Status:           20, // CLOSED
			OutTradeNo:       req.OutTradeNo,
			ChannelErrorCode: "H5_PREPAY_ERROR",
			ChannelErrorMsg:  err.Error(),
		}, nil
	}
	_ = result

	return &client.OrderResp{
		Status:         0, // WAITING
		OutTradeNo:     req.OutTradeNo,
		DisplayMode:    "url",
		DisplayContent: *resp.H5Url,
	}, nil
}

// appOrder APP 支付
func (c *WxPayClient) appOrder(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
	svc := app.AppApiService{Client: c.coreClient}

	resp, result, err := svc.Prepay(ctx, app.PrepayRequest{
		Appid:       core.String(c.config.AppID),
		Mchid:       core.String(c.config.MchID),
		Description: core.String(req.Subject),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(req.NotifyURL),
		Amount: &app.Amount{
			Total:    core.Int64(int64(req.Price)),
			Currency: core.String("CNY"),
		},
	})

	if err != nil {
		return &client.OrderResp{
			Status:           20, // CLOSED
			OutTradeNo:       req.OutTradeNo,
			ChannelErrorCode: "APP_PREPAY_ERROR",
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
	svc := refunddomestic.RefundsApiService{Client: c.coreClient}

	resp, _, err := svc.Create(ctx, refunddomestic.CreateRequest{
		OutTradeNo:  core.String(req.OutTradeNo),
		OutRefundNo: core.String(req.OutRefundNo),
		Reason:      core.String(req.Reason),
		NotifyUrl:   core.String(req.NotifyURL),
		Amount: &refunddomestic.AmountReq{
			Currency: core.String("CNY"),
			Refund:   core.Int64(int64(req.RefundPrice)),
			Total:    core.Int64(int64(req.PayPrice)),
		},
	})

	if err != nil {
		return &client.RefundResp{
			Status:           20, // FALLBACK/FAILURE (Need better mapping)
			OutTradeNo:       req.OutTradeNo,
			OutRefundNo:      req.OutRefundNo,
			ChannelErrorCode: "REFUND_ERROR",
			ChannelErrorMsg:  err.Error(),
		}, nil
	}

	status := 0 // WAITING
	if *resp.Status == refunddomestic.STATUS_SUCCESS {
		status = 10 // SUCCESS
	} else if *resp.Status == refunddomestic.STATUS_CLOSED || *resp.Status == refunddomestic.STATUS_ABNORMAL {
		status = 20 // FAILURE
	}

	return &client.RefundResp{
		Status:          status,
		OutTradeNo:      req.OutTradeNo,
		OutRefundNo:     req.OutRefundNo,
		ChannelRefundNo: *resp.RefundId,
		SuccessTime:     time.Now(), // TODO: Parse SuccessTime if available
	}, nil
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
	// 1. 构造 http.Request
	httpReq := &http.Request{
		Header: http.Header{},
		Body:   io.NopCloser(strings.NewReader(req.Body)),
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 2. 初始化 NotifyHandler
	verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(c.config.PublicKeyID, *c.publicKey)
	handler, err := notify.NewRSANotifyHandler(c.config.APIV3Key, verifier)
	if err != nil {
		return nil, fmt.Errorf("创建回调处理器失败: %v", err)
	}
	// 3. 解析并验证签名
	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), httpReq, transaction)
	if err != nil {
		return nil, fmt.Errorf("解析支付回调失败: %w", err)
	}

	_ = notifyReq

	// 4. 转换结果
	status := 0
	if *transaction.TradeState == "SUCCESS" {
		status = 10
	} else if *transaction.TradeState == "CLOSED" || *transaction.TradeState == "PAYERROR" {
		status = 20
	}

	var successTime time.Time
	if transaction.SuccessTime != nil {
		successTime, _ = time.Parse(time.RFC3339, *transaction.SuccessTime)
	}

	return &client.OrderResp{
		Status:         status,
		OutTradeNo:     *transaction.OutTradeNo,
		ChannelOrderNo: *transaction.TransactionId,
		ChannelUserID:  *transaction.Payer.Openid,
		SuccessTime:    successTime,
		RawData:        req.Body,
	}, nil
}

// ParseRefundNotify 解析退款回调
func (c *WxPayClient) ParseRefundNotify(req *client.NotifyData) (*client.RefundResp, error) {
	// 1. 构造 http.Request
	httpReq := &http.Request{
		Header: http.Header{},
		Body:   io.NopCloser(strings.NewReader(req.Body)),
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 2. 初始化 NotifyHandler
	verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(c.config.PublicKeyID, *c.publicKey)
	handler := notify.NewNotifyHandler(c.config.APIV3Key, verifier)

	// 3. 解析并验证签名
	refundNotify := new(refunddomestic.Refund)
	_, err := handler.ParseNotifyRequest(context.Background(), httpReq, refundNotify)
	if err != nil {
		return nil, fmt.Errorf("解析退款回调失败: %w", err)
	}

	// 4. 转换结果
	status := 0
	if *refundNotify.Status == refunddomestic.STATUS_SUCCESS {
		status = 10
	} else {
		status = 20
	}

	var successTime time.Time
	if refundNotify.SuccessTime != nil {
		successTime = *refundNotify.SuccessTime
	}

	return &client.RefundResp{
		Status:          status,
		OutTradeNo:      *refundNotify.OutTradeNo,
		OutRefundNo:     *refundNotify.OutRefundNo,
		ChannelRefundNo: *refundNotify.RefundId,
		SuccessTime:     successTime,
		RawData:         req.Body,
	}, nil
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

// ParseTransferNotify 解析转账回调
// 对齐 Java: AbstractWxPayClient.parseTransferNotifyV3
// 注意: 仅支持 V3 版本，V2 不支持转账回调
func (c *WxPayClient) ParseTransferNotify(req *client.NotifyData) (*client.TransferResp, error) {
	// 1. 构造 http.Request
	httpReq := &http.Request{
		Header: http.Header{},
		Body:   io.NopCloser(strings.NewReader(req.Body)),
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 2. 初始化 NotifyHandler
	verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(c.config.PublicKeyID, *c.publicKey)
	handler, err := notify.NewRSANotifyHandler(c.config.APIV3Key, verifier)
	if err != nil {
		return nil, fmt.Errorf("创建回调处理器失败: %v", err)
	}

	// 3. 解析并验证签名 (使用 map 接收通用回调数据)
	content := make(map[string]interface{})
	_, err = handler.ParseNotifyRequest(context.Background(), httpReq, &content)
	if err != nil {
		return nil, fmt.Errorf("解析转账回调失败: %w", err)
	}

	// 4. 提取字段
	state, _ := content["state"].(string)
	outBizNo, _ := content["out_bill_no"].(string)
	transferBillNo, _ := content["transfer_bill_no"].(string)
	updateTimeStr, _ := content["update_time"].(string)
	failReason, _ := content["fail_reason"].(string)

	// 5. 解析时间
	var successTime time.Time
	if updateTimeStr != "" {
		successTime, _ = time.Parse(time.RFC3339, updateTimeStr)
	}

	// 6. 根据状态转换 (对齐 Java 的状态判断逻辑)
	var transferStatus int
	// ACCEPTED, PROCESSING, WAIT_USER_CONFIRM, TRANSFERING -> 处理中 (5)
	if state == "ACCEPTED" || state == "PROCESSING" || state == "WAIT_USER_CONFIRM" || state == "TRANSFERING" {
		transferStatus = 5 // PayTransferStatusProcessing
	} else if state == "SUCCESS" {
		transferStatus = 10 // PayTransferStatusSuccess
	} else {
		// 其他状态视为关闭 (20)
		transferStatus = 20 // PayTransferStatusClosed
	}

	return &client.TransferResp{
		Status:            transferStatus,
		OutTradeNo:        outBizNo,
		ChannelTransferNo: transferBillNo,
		SuccessTime:       successTime,
		ChannelErrorMsg:   failReason,
		RawData:           req.Body,
	}, nil
}
