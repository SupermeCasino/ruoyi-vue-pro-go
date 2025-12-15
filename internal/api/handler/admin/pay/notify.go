package pay

import (
	"io"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	paySvc "backend-go/internal/service/pay"
	"backend-go/internal/service/pay/client"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type PayNotifyHandler struct {
	svc         *paySvc.PayNotifyService
	appSvc      *paySvc.PayAppService
	channelSvc  *paySvc.PayChannelService
	orderSvc    *paySvc.PayOrderService
	refundSvc   *paySvc.PayRefundService
	transferSvc *paySvc.PayTransferService
	logger      *zap.Logger
}

func NewPayNotifyHandler(
	svc *paySvc.PayNotifyService,
	appSvc *paySvc.PayAppService,
	channelSvc *paySvc.PayChannelService,
	orderSvc *paySvc.PayOrderService,
	refundSvc *paySvc.PayRefundService,
	transferSvc *paySvc.PayTransferService,
	logger *zap.Logger,
) *PayNotifyHandler {
	return &PayNotifyHandler{
		svc:         svc,
		appSvc:      appSvc,
		channelSvc:  channelSvc,
		orderSvc:    orderSvc,
		refundSvc:   refundSvc,
		transferSvc: transferSvc,
		logger:      logger,
	}
}

// NotifyOrder 支付渠道的统一【支付】回调
// POST /pay/notify/order/:channelId
// 对齐 Java: PayNotifyController.notifyOrder
func (h *PayNotifyHandler) NotifyOrder(c *gin.Context) {
	channelId := core.ParseInt64(c.Param("channelId"))
	h.logger.Info("[NotifyOrder] 收到支付回调", zap.Int64("channelId", channelId))

	// 1. 获取 PayClient
	payClient := h.channelSvc.GetPayClient(channelId)
	if payClient == nil {
		h.logger.Error("[NotifyOrder] 渠道编号找不到对应的支付客户端", zap.Int64("channelId", channelId))
		c.String(400, "渠道不存在")
		return
	}

	// 2. 解析回调数据
	body, _ := io.ReadAll(c.Request.Body)
	notifyData := &client.NotifyData{
		Params:  h.queryToMap(c),
		Body:    string(body),
		Headers: h.headerToMap(c),
	}

	orderResp, err := payClient.ParseOrderNotify(notifyData)
	if err != nil {
		h.logger.Error("[NotifyOrder] 解析回调数据失败", zap.Error(err))
		c.String(400, "解析失败: "+err.Error())
		return
	}

	// 3. 处理回调
	if err := h.orderSvc.NotifyOrder(c.Request.Context(), channelId, orderResp); err != nil {
		h.logger.Error("[NotifyOrder] 处理回调失败", zap.Error(err))
		c.String(500, "处理失败: "+err.Error())
		return
	}

	h.logger.Info("[NotifyOrder] 支付回调处理成功", zap.Int64("channelId", channelId))
	c.String(200, "success")
}

// NotifyRefund 支付渠道的统一【退款】回调
// POST /pay/notify/refund/:channelId
// 对齐 Java: PayNotifyController.notifyRefund
func (h *PayNotifyHandler) NotifyRefund(c *gin.Context) {
	channelId := core.ParseInt64(c.Param("channelId"))
	h.logger.Info("[NotifyRefund] 收到退款回调", zap.Int64("channelId", channelId))

	// 1. 获取 PayClient
	payClient := h.channelSvc.GetPayClient(channelId)
	if payClient == nil {
		h.logger.Error("[NotifyRefund] 渠道编号找不到对应的支付客户端", zap.Int64("channelId", channelId))
		c.String(400, "渠道不存在")
		return
	}

	// 2. 解析回调数据
	body, _ := io.ReadAll(c.Request.Body)
	notifyData := &client.NotifyData{
		Params:  h.queryToMap(c),
		Body:    string(body),
		Headers: h.headerToMap(c),
	}

	refundResp, err := payClient.ParseRefundNotify(notifyData)
	if err != nil {
		h.logger.Error("[NotifyRefund] 解析回调数据失败", zap.Error(err))
		c.String(400, "解析失败: "+err.Error())
		return
	}

	// 3. 处理回调
	if err := h.refundSvc.NotifyRefund(c.Request.Context(), channelId, refundResp); err != nil {
		h.logger.Error("[NotifyRefund] 处理回调失败", zap.Error(err))
		c.String(500, "处理失败: "+err.Error())
		return
	}

	h.logger.Info("[NotifyRefund] 退款回调处理成功", zap.Int64("channelId", channelId))
	c.String(200, "success")
}

// NotifyTransfer 支付渠道的统一【转账】回调
// POST /pay/notify/transfer/:channelId
// 对齐 Java: PayNotifyController.notifyTransfer
func (h *PayNotifyHandler) NotifyTransfer(c *gin.Context) {
	channelId := core.ParseInt64(c.Param("channelId"))
	h.logger.Info("[NotifyTransfer] 收到转账回调", zap.Int64("channelId", channelId))

	// 1. 获取 PayClient
	payClient := h.channelSvc.GetPayClient(channelId)
	if payClient == nil {
		h.logger.Error("[NotifyTransfer] 渠道编号找不到对应的支付客户端", zap.Int64("channelId", channelId))
		c.String(400, "渠道不存在")
		return
	}

	// 2. 解析回调数据
	body, _ := io.ReadAll(c.Request.Body)
	notifyData := &client.NotifyData{
		Params:  h.queryToMap(c),
		Body:    string(body),
		Headers: h.headerToMap(c),
	}

	transferResp, err := payClient.ParseTransferNotify(notifyData)
	if err != nil {
		h.logger.Error("[NotifyTransfer] 解析回调数据失败", zap.Error(err))
		c.String(400, "解析失败: "+err.Error())
		return
	}

	// 3. 处理回调
	// 注意：需要注入 transferSvc
	if h.transferSvc != nil {
		if err := h.transferSvc.NotifyTransfer(c.Request.Context(), channelId, transferResp); err != nil {
			h.logger.Error("[NotifyTransfer] 处理回调失败", zap.Error(err))
			c.String(500, "处理失败: "+err.Error())
			return
		}
	}

	h.logger.Info("[NotifyTransfer] 转账回调处理成功", zap.Int64("channelId", channelId))
	c.String(200, "success")
}

// GetNotifyTaskDetail 获得回调通知详情 (Task + Logs)
func (h *PayNotifyHandler) GetNotifyTaskDetail(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	task, err := h.svc.GetNotifyTask(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	if task == nil {
		c.JSON(200, core.Success(&resp.PayNotifyTaskDetailResp{}))
		return
	}

	logs, _ := h.svc.GetNotifyLogList(c, id)
	app, _ := h.appSvc.GetApp(c, task.AppID)

	r := &resp.PayNotifyTaskDetailResp{}
	copier.Copy(&r.PayNotifyTaskResp, task)
	if app != nil {
		r.AppName = app.Name
	}

	logResps := make([]*resp.PayNotifyLogResp, 0, len(logs))
	for _, log := range logs {
		lr := &resp.PayNotifyLogResp{}
		copier.Copy(lr, log)
		logResps = append(logResps, lr)
	}
	r.Logs = logResps

	c.JSON(200, core.Success(r))
}

// GetNotifyTaskPage 获得回调通知分页
func (h *PayNotifyHandler) GetNotifyTaskPage(c *gin.Context) {
	var r req.PayNotifyTaskPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	pageResult, err := h.svc.GetNotifyTaskPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	var appIds []int64
	for _, item := range pageResult.List {
		appIds = append(appIds, item.AppID)
	}
	appMap, _ := h.appSvc.GetAppMap(c, appIds)

	list := make([]*resp.PayNotifyTaskResp, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		tr := &resp.PayNotifyTaskResp{}
		copier.Copy(tr, item)
		if app, ok := appMap[item.AppID]; ok {
			tr.AppName = app.Name
		}
		list = append(list, tr)
	}

	c.JSON(200, core.Success(core.PageResult[*resp.PayNotifyTaskResp]{
		List:  list,
		Total: pageResult.Total,
	}))
}

// Helpers: 将 query 和 header 转换为 map
func (h *PayNotifyHandler) queryToMap(c *gin.Context) map[string]string {
	result := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

func (h *PayNotifyHandler) headerToMap(c *gin.Context) map[string]string {
	result := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

// Helpers
func convertNotifyTaskResp(task *pay.PayNotifyTask, app *pay.PayApp) *resp.PayNotifyTaskResp {
	r := &resp.PayNotifyTaskResp{}
	copier.Copy(r, task)
	if app != nil {
		r.AppName = app.Name
	}
	return r
}
