package pay

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	payWalletSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type PayOrderHandler struct {
	svc       *paySvc.PayOrderService
	appSvc    *paySvc.PayAppService
	walletSvc *payWalletSvc.PayWalletService
}

func NewPayOrderHandler(svc *paySvc.PayOrderService, appSvc *paySvc.PayAppService, walletSvc *payWalletSvc.PayWalletService) *PayOrderHandler {
	return &PayOrderHandler{svc: svc, appSvc: appSvc, walletSvc: walletSvc}
}

// GetOrder 获得支付订单
func (h *PayOrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	// Handle sync param
	syncStr := c.Query("sync")
	if syncStr == "true" {
		order, err := h.svc.GetOrder(c, id)
		if err == nil && order.Status == paySvc.PayOrderStatusWaiting {
			h.svc.SyncOrderQuietly(c, id)
		}
	}

	order, err := h.svc.GetOrder(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(convertOrderResp(order)))
}

// GetOrderDetail 获得支付订单详情
func (h *PayOrderHandler) GetOrderDetail(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	order, err := h.svc.GetOrder(c, id)
	if err != nil {
		// If order not found, return Success(nil) as per Java logic?
		// Or Error? Java: return success(null) if order is null.
		// Go gorm returns generic error usually.
		c.Error(err)
		return
	}

	app, _ := h.appSvc.GetApp(c, order.AppID)
	extension, _ := h.svc.GetOrderExtension(c, order.ExtensionID)

	detail := &resp.PayOrderDetailsResp{
		PayOrderResp: *convertOrderResp(order),
		Extension:    convertExtensionResp(extension),
		App:          convertAppResp(app), // Need to export/access this converter or rewrite
	}
	c.JSON(200, response.Success(detail))
}

// GetOrderPage 获得支付订单分页
func (h *PayOrderHandler) GetOrderPage(c *gin.Context) {
	var r req.PayOrderPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	pageResult, err := h.svc.GetOrderPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// Fetch apps for mapping
	appIds := make([]int64, 0, len(pageResult.List))
	for _, order := range pageResult.List {
		appIds = append(appIds, order.AppID)
	}
	appMap, _ := h.appSvc.GetAppMap(c, appIds)

	list := make([]resp.PayOrderResp, 0, len(pageResult.List))
	for _, order := range pageResult.List {
		r := *convertOrderResp(order)
		if app, ok := appMap[order.AppID]; ok {
			r.AppName = app.Name
		}
		list = append(list, r)
	}

	c.JSON(200, response.Success(pagination.PageResult[resp.PayOrderResp]{
		List:  list,
		Total: pageResult.Total,
	}))
}

// SubmitPayOrder 提交支付订单
func (h *PayOrderHandler) SubmitPayOrder(c *gin.Context) {
	var r req.PayOrderSubmitReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	// 1. Wallet payment case
	if r.ChannelCode == "wallet" { // Assuming "wallet" is the code
		if r.ChannelExtras == nil {
			r.ChannelExtras = make(map[string]string)
		}
		userID := context.GetLoginUserID(c)
		user := context.GetLoginUser(c)
		userType := 0
		if user != nil {
			userType = user.UserType
		}
		wallet, err := h.walletSvc.GetOrCreateWallet(c, userID, userType)
		if err != nil {
			c.Error(err)
			return
		}
		r.ChannelExtras["wallet_id"] = strconv.FormatInt(wallet.ID, 10)
	}

	// 2. Submit Order
	respVO, err := h.svc.SubmitOrder(c, &r, c.ClientIP())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(respVO))
}

// Helpers

func convertOrderResp(order *pay.PayOrder) *resp.PayOrderResp {
	if order == nil {
		return nil
	}
	return &resp.PayOrderResp{
		ID:              order.ID,
		AppID:           order.AppID,
		ChannelID:       order.ChannelID,
		ChannelCode:     order.ChannelCode,
		UserID:          order.UserID,
		UserType:        order.UserType,
		MerchantOrderId: order.MerchantOrderId,
		Subject:         order.Subject,
		Body:            order.Body,
		NotifyURL:       order.NotifyURL,
		Price:           order.Price,
		ChannelFeeRate:  order.ChannelFeeRate,
		ChannelFeePrice: order.ChannelFeePrice,
		Status:          order.Status,
		UserIP:          order.UserIP,
		ExpireTime:      order.ExpireTime,
		SuccessTime:     order.SuccessTime,
		ExtensionID:     order.ExtensionID,
		No:              order.No,
		RefundPrice:     order.RefundPrice,
		ChannelUserID:   order.ChannelUserID,
		ChannelOrderNo:  order.ChannelOrderNo,
		CreateTime:      order.CreatedAt,
		UpdateTime:      order.UpdatedAt,
		Creator:         order.Creator,
		Updater:         order.Updater,
		Deleted:         order.Deleted,
	}
}

func convertExtensionResp(ext *pay.PayOrderExtension) *resp.PayOrderExtensionResp {
	if ext == nil {
		return nil
	}
	return &resp.PayOrderExtensionResp{
		ID:                ext.ID,
		No:                ext.No,
		OrderID:           ext.OrderID,
		ChannelID:         ext.ChannelID,
		ChannelCode:       ext.ChannelCode,
		UserIP:            ext.UserIP,
		Status:            ext.Status,
		ChannelExtras:     ext.ChannelExtras,
		ChannelErrorCode:  ext.ChannelErrorCode,
		ChannelErrorMsg:   ext.ChannelErrorMsg,
		ChannelNotifyData: ext.ChannelNotifyData,
		CreateTime:        ext.CreatedAt,
	}
}

func convertAppResp(app *pay.PayApp) *resp.PayAppResp {
	// Duplicated from PayAppHandler to avoid import cycle or dependency issues
	if app == nil {
		return nil
	}
	return &resp.PayAppResp{
		ID:                app.ID,
		AppKey:            app.AppKey,
		Name:              app.Name,
		Status:            app.Status,
		Remark:            app.Remark,
		OrderNotifyURL:    app.OrderNotifyURL,
		RefundNotifyURL:   app.RefundNotifyURL,
		TransferNotifyURL: app.TransferNotifyURL,
		CreateTime:        app.CreatedAt,
	}
}
