package pay

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/pay/client"
	"backend-go/pkg/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PayOrderService struct {
	q          *query.Query
	appSvc     *PayAppService
	channelSvc *PayChannelService
	clientFac  *client.PayClientFactory
	notifySvc  *PayNotifyService
}

func NewPayOrderService(q *query.Query, appSvc *PayAppService, channelSvc *PayChannelService, clientFac *client.PayClientFactory, notifySvc *PayNotifyService) *PayOrderService {
	return &PayOrderService{
		q:          q,
		appSvc:     appSvc,
		channelSvc: channelSvc,
		clientFac:  clientFac,
		notifySvc:  notifySvc,
	}
}

// GetOrder 获得支付订单
func (s *PayOrderService) GetOrder(ctx context.Context, id int64) (*pay.PayOrder, error) {
	return s.q.PayOrder.WithContext(ctx).Where(s.q.PayOrder.ID.Eq(id)).First()
}

// GetOrderPage 获得支付订单分页
func (s *PayOrderService) GetOrderPage(ctx context.Context, req *req.PayOrderPageReq) (*core.PageResult[*pay.PayOrder], error) {
	q := s.q.PayOrder.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayOrder.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayOrder.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayOrder.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.Subject != "" {
		q = q.Where(s.q.PayOrder.Subject.Like("%" + req.Subject + "%"))
	}
	if req.No != "" {
		q = q.Where(s.q.PayOrder.No.Eq(req.No))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayOrder.Status.Eq(*req.Status))
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}
	list, err := q.Limit(req.GetLimit()).Offset(req.GetOffset()).Order(s.q.PayOrder.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*pay.PayOrder]{
		List:  list,
		Total: total,
	}, nil
}

// CreateOrder 创建支付单
func (s *PayOrderService) CreateOrder(ctx context.Context, reqDTO *req.PayOrderCreateReq) (int64, error) {
	app, err := s.appSvc.GetApp(ctx, reqDTO.AppID)
	if err != nil {
		return 0, err
	}
	if app == nil || app.Status != 0 {
		return 0, errors.New("App disabled or not found")
	}

	existOrder, _ := s.q.PayOrder.WithContext(ctx).
		Where(s.q.PayOrder.AppID.Eq(app.ID), s.q.PayOrder.MerchantOrderId.Eq(reqDTO.MerchantOrderId)).
		First()
	if existOrder != nil {
		return existOrder.ID, nil
	}

	// 创建支付交易单 (对齐 Java: 使用 app.OrderNotifyURL)
	order := &pay.PayOrder{
		AppID:           app.ID,
		MerchantOrderId: reqDTO.MerchantOrderId,
		Subject:         reqDTO.Subject,
		Body:            reqDTO.Body,
		NotifyURL:       app.OrderNotifyURL, // 对齐 Java: 使用 app 的回调地址
		Price:           reqDTO.Price,
		ExpireTime:      time.Now().Add(2 * time.Hour),
		Status:          PayOrderStatusWaiting,
		RefundPrice:     0,
		UserIP:          reqDTO.UserIP,
	}

	if err := s.q.PayOrder.WithContext(ctx).Create(order); err != nil {
		return 0, err
	}
	return order.ID, nil
}

// ... GetOrderCountByAppId ...

// SubmitOrder 提交支付订单
func (s *PayOrderService) SubmitOrder(ctx context.Context, reqVO *req.PayOrderSubmitReq, userIP string) (*resp.PayOrderSubmitResp, error) {
	order, err := s.validateOrderCanSubmit(ctx, reqVO.ID)
	if err != nil {
		return nil, err
	}

	channel, err := s.validateChannelCanSubmit(ctx, order.AppID, reqVO.ChannelCode)
	if err != nil {
		return nil, err
	}

	// Generate No
	no := s.generateNo()

	// Create Extension
	ext := &pay.PayOrderExtension{
		OrderID:     order.ID,
		No:          no,
		ChannelID:   channel.ID,
		ChannelCode: channel.Code,
		UserIP:      userIP,
		Status:      PayOrderStatusWaiting,
	}
	if err := s.q.PayOrderExtension.WithContext(ctx).Create(ext); err != nil {
		return nil, err
	}

	// Get Pay Client
	payClient := s.clientFac.GetPayClient(channel.ID)
	if payClient == nil {
		// Lazy create if not exists
		var err error
		payClient, err = s.clientFac.CreateOrUpdatePayClient(channel.ID, channel.Code, channel.Config.ToJSON())
		if err != nil {
			return nil, err
		}
	}

	// Call UnifiedOrder (对齐 Java: 使用渠道特定的回调 URL)
	unifiedReq := &client.UnifiedOrderReq{
		UserIP:     userIP,
		OutTradeNo: no,
		Subject:    order.Subject,
		Body:       order.Body,
		NotifyURL:  s.genChannelOrderNotifyUrl(channel), // 对齐 Java: 渠道回调 URL
		// ReturnURL:   reqVO.ReturnUrl,
		Price:       order.Price,
		ExpireTime:  order.ExpireTime,
		DisplayMode: reqVO.DisplayMode,
	}
	unifiedResp, err := payClient.UnifiedOrder(ctx, unifiedReq)
	if err != nil {
		return nil, err
	}

	// Return response
	return &resp.PayOrderSubmitResp{
		Status:         unifiedResp.Status,
		DisplayMode:    unifiedResp.DisplayMode,
		DisplayContent: unifiedResp.DisplayContent,
	}, nil
}

func (s *PayOrderService) validateOrderCanSubmit(ctx context.Context, id int64) (*pay.PayOrder, error) {
	order, err := s.q.PayOrder.WithContext(ctx).Where(s.q.PayOrder.ID.Eq(id)).First()
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	if order.Status == PayOrderStatusSuccess {
		return nil, errors.New("Order already paid")
	}
	if order.Status != PayOrderStatusWaiting {
		return nil, errors.New("Order status not waiting")
	}
	if order.ExpireTime.Before(time.Now()) {
		return nil, errors.New("Order expired")
	}
	return order, nil
}

func (s *PayOrderService) validateChannelCanSubmit(ctx context.Context, appId int64, channelCode string) (*pay.PayChannel, error) {
	// app validation is implicit or done separately
	return s.channelSvc.GetChannelByAppIdAndCode(ctx, appId, channelCode)
}

// genChannelOrderNotifyUrl 根据支付渠道生成回调地址
// 对齐 Java: payProperties.getOrderNotifyUrl() + "/" + channel.getId()
func (s *PayOrderService) genChannelOrderNotifyUrl(channel *pay.PayChannel) string {
	return fmt.Sprintf("%s/%d", config.C.Pay.OrderNotifyURL, channel.ID)
}

func (s *PayOrderService) generateNo() string {
	// Simple timestamp + random for now.
	// Java uses Redis. We can use core.RDB.Incr if we want strictly strict.
	// For MVP: P + yyyyMMddHHmmss + 6 digit random
	return "P" + time.Now().Format("20060102150405") + core.GenerateRandomString(6) // Need helper?
	// Let's use simplified version
	return "P" + time.Now().Format("20060102150405") + "000000"
}

// GetOrderExtension 获得支付订单拓展
func (s *PayOrderService) GetOrderExtension(ctx context.Context, id int64) (*pay.PayOrderExtension, error) {
	return s.q.PayOrderExtension.WithContext(ctx).Where(s.q.PayOrderExtension.ID.Eq(id)).First()
}

// SyncOrderQuietly 同步订单的支付状态 (Quietly)
// 对齐 Java: PayOrderServiceImpl.syncOrderQuietly
func (s *PayOrderService) SyncOrderQuietly(ctx context.Context, id int64) {
	// 1. 查询待支付订单拓展
	extensions, err := s.q.PayOrderExtension.WithContext(ctx).
		Where(s.q.PayOrderExtension.OrderID.Eq(id), s.q.PayOrderExtension.Status.Eq(PayOrderStatusWaiting)).
		Find()
	if err != nil {
		return
	}

	// 2. 遍历执行同步
	for _, ext := range extensions {
		s.syncOrder(ctx, ext)
	}
}

// syncOrder 同步单个支付拓展单
// 对齐 Java: PayOrderServiceImpl.syncOrder(PayOrderExtensionDO)
func (s *PayOrderService) syncOrder(ctx context.Context, orderExtension *pay.PayOrderExtension) bool {
	// 1.1 查询支付订单信息
	payClient := s.clientFac.GetPayClient(orderExtension.ChannelID)
	if payClient == nil {
		return false
	}

	respDTO, err := payClient.GetOrder(ctx, orderExtension.No)
	if err != nil {
		return false
	}

	// 如果查询到订单不存在,PayClient 返回的状态为关闭。但此时不能关闭订单。
	// 存在以下场景:拉起渠道支付后,短时间内用户未及时完成支付,但是该订单同步定时任务恰巧自动触发了,
	// 主动查询结果为订单不存在。当用户支付成功之后,该订单状态在渠道的回调中无法从已关闭改为已支付,造成重大影响。
	// 考虑此定时任务是异常场景的兜底操作,因此这里不做变更,优先以回调为准。
	if respDTO.Status == PayOrderStatusClosed {
		return false
	}

	// 1.2 回调支付结果
	s.NotifyOrder(ctx, orderExtension.ChannelID, respDTO)

	// 2. 如果是已支付,则返回 true
	return respDTO.Status == PayOrderStatusSuccess
}

// NotifyOrder 通知并更新订单的支付结果
// 对齐 Java: PayOrderService.notifyOrder(Long channelId, PayOrderRespDTO notify)
func (s *PayOrderService) NotifyOrder(ctx context.Context, channelID int64, notify *client.OrderResp) error {
	// 校验支付渠道是否有效
	channel, err := s.channelSvc.GetChannel(ctx, channelID)
	if err != nil {
		return err
	}

	// 情况一:支付成功的回调
	if notify.Status == PayOrderStatusSuccess {
		return s.notifyOrderSuccess(ctx, channel, notify)
	}

	// 情况二:支付失败的回调
	if notify.Status == PayOrderStatusClosed {
		return s.notifyOrderClosed(ctx, channel, notify)
	}

	// 情况三:WAITING 无需处理
	// 情况四:REFUND 通过退款回调处理
	return nil
}

// notifyOrderSuccess 处理支付成功的回调
func (s *PayOrderService) notifyOrderSuccess(ctx context.Context, channel *pay.PayChannel, notify *client.OrderResp) error {
	// 1. 更新 PayOrderExtension 支付成功
	orderExtension, err := s.updateOrderExtensionSuccess(ctx, notify)
	if err != nil {
		return err
	}

	// 2. 更新 PayOrder 支付成功
	paid, err := s.updateOrderSuccess(ctx, channel, orderExtension, notify)
	if err != nil {
		return err
	}
	if paid {
		// 如果之前已经成功回调,则直接返回,不用重复记录支付通知记录
		return nil
	}

	// 3. 插入支付通知记录
	s.notifySvc.CreatePayNotifyTask(ctx, PayNotifyTypeOrder, orderExtension.OrderID)

	return nil
}

// updateOrderExtensionSuccess 更新 PayOrderExtension 支付成功
func (s *PayOrderService) updateOrderExtensionSuccess(ctx context.Context, notify *client.OrderResp) (*pay.PayOrderExtension, error) {
	// 1. 查询 PayOrderExtension
	orderExtension, err := s.q.PayOrderExtension.WithContext(ctx).
		Where(s.q.PayOrderExtension.No.Eq(notify.OutTradeNo)).
		First()
	if err != nil {
		return nil, fmt.Errorf("支付订单拓展不存在")
	}

	// 如果已经是成功,直接返回,不用重复更新
	if orderExtension.Status == PayOrderStatusSuccess {
		return orderExtension, nil
	}

	// 校验状态,必须是待支付
	if orderExtension.Status != PayOrderStatusWaiting {
		return nil, fmt.Errorf("支付订单拓展状态不是待支付")
	}

	// 2. 更新 PayOrderExtension (使用乐观锁)
	notifyDataJSON, _ := json.Marshal(notify)
	result, err := s.q.PayOrderExtension.WithContext(ctx).
		Where(s.q.PayOrderExtension.ID.Eq(orderExtension.ID), s.q.PayOrderExtension.Status.Eq(PayOrderStatusWaiting)).
		Updates(map[string]interface{}{
			"status":              PayOrderStatusSuccess,
			"channel_notify_data": string(notifyDataJSON),
		})

	if err != nil || result.RowsAffected == 0 {
		return nil, fmt.Errorf("支付订单拓展状态不是待支付")
	}

	orderExtension.Status = PayOrderStatusSuccess
	return orderExtension, nil
}

// updateOrderSuccess 更新 PayOrder 支付成功
// 返回值: 是否之前已经成功回调
func (s *PayOrderService) updateOrderSuccess(ctx context.Context, channel *pay.PayChannel, orderExtension *pay.PayOrderExtension, notify *client.OrderResp) (bool, error) {
	// 1. 判断 PayOrder 是否处于待支付
	order, err := s.q.PayOrder.WithContext(ctx).
		Where(s.q.PayOrder.ID.Eq(orderExtension.OrderID)).
		First()
	if err != nil {
		return false, fmt.Errorf("支付订单不存在")
	}

	// 如果已经是成功,直接返回,不用重复更新
	if order.Status == PayOrderStatusSuccess && order.ExtensionID == orderExtension.ID {
		return true, nil
	}

	// 校验状态,必须是待支付
	if order.Status != PayOrderStatusWaiting {
		return false, fmt.Errorf("支付订单状态不是待支付")
	}

	// 2. 更新 PayOrder (使用乐观锁)
	channelFeePrice := int(float64(order.Price) * channel.FeeRate / 100)
	now := time.Now()

	result, err := s.q.PayOrder.WithContext(ctx).
		Where(s.q.PayOrder.ID.Eq(order.ID), s.q.PayOrder.Status.Eq(PayOrderStatusWaiting)).
		Updates(map[string]interface{}{
			"status":            PayOrderStatusSuccess,
			"channel_id":        channel.ID,
			"channel_code":      channel.Code,
			"success_time":      &now,
			"extension_id":      orderExtension.ID,
			"no":                orderExtension.No,
			"channel_order_no":  notify.ChannelOrderNo,
			"channel_user_id":   notify.ChannelUserID,
			"channel_fee_rate":  channel.FeeRate,
			"channel_fee_price": channelFeePrice,
		})

	if err != nil || result.RowsAffected == 0 {
		return false, fmt.Errorf("支付订单状态不是待支付")
	}

	return false, nil
}

// notifyOrderClosed 处理支付关闭的回调
func (s *PayOrderService) notifyOrderClosed(ctx context.Context, channel *pay.PayChannel, notify *client.OrderResp) error {
	// 查询 PayOrderExtension
	orderExtension, err := s.q.PayOrderExtension.WithContext(ctx).
		Where(s.q.PayOrderExtension.No.Eq(notify.OutTradeNo)).
		First()
	if err != nil {
		return fmt.Errorf("支付订单拓展不存在")
	}

	// 如果已经是关闭,直接返回
	if orderExtension.Status == PayOrderStatusClosed {
		return nil
	}

	// 一般出现先是支付成功,然后支付关闭,都是全部退款导致关闭的场景
	// 这个情况,我们不更新支付拓展单,只通过退款流程,更新支付单
	if orderExtension.Status == PayOrderStatusSuccess {
		return nil
	}

	// 校验状态,必须是待支付
	if orderExtension.Status != PayOrderStatusWaiting {
		return fmt.Errorf("支付订单拓展状态不是待支付")
	}

	// 更新 PayOrderExtension
	notifyDataJSON, _ := json.Marshal(notify)
	result, err := s.q.PayOrderExtension.WithContext(ctx).
		Where(s.q.PayOrderExtension.ID.Eq(orderExtension.ID), s.q.PayOrderExtension.Status.Eq(PayOrderStatusWaiting)).
		Updates(map[string]interface{}{
			"status":              PayOrderStatusClosed,
			"channel_notify_data": string(notifyDataJSON),
			"channel_error_code":  notify.ChannelErrorCode,
			"channel_error_msg":   notify.ChannelErrorMsg,
		})

	if err != nil || result.RowsAffected == 0 {
		return fmt.Errorf("支付订单拓展状态不是待支付")
	}

	return nil
}

// UpdateOrderRefundPrice 更新订单退款金额
// 对齐 Java: PayOrderService.updateOrderRefundPrice(Long id, Integer incrRefundPrice)
func (s *PayOrderService) UpdateOrderRefundPrice(ctx context.Context, id int64, incrRefundPrice int) error {
	order, err := s.q.PayOrder.WithContext(ctx).Where(s.q.PayOrder.ID.Eq(id)).First()
	if err != nil {
		return fmt.Errorf("支付订单不存在")
	}

	// 校验状态：必须是已支付或已退款
	if order.Status != PayOrderStatusSuccess && order.Status != PayOrderStatusRefund {
		return fmt.Errorf("支付订单状态不是已支付或已退款")
	}

	// 校验退款金额不能超过支付金额
	if order.RefundPrice+incrRefundPrice > order.Price {
		return fmt.Errorf("退款金额超过支付金额")
	}

	// 更新订单 (使用乐观锁)
	result, err := s.q.PayOrder.WithContext(ctx).
		Where(s.q.PayOrder.ID.Eq(id), s.q.PayOrder.Status.Eq(order.Status)).
		Updates(map[string]interface{}{
			"refund_price": order.RefundPrice + incrRefundPrice,
			"status":       PayOrderStatusRefund,
		})

	if err != nil || result.RowsAffected == 0 {
		return fmt.Errorf("支付订单状态不是已支付或已退款")
	}

	return nil
}

// GetOrderList 获得支付订单列表 (Export)
func (s *PayOrderService) GetOrderList(ctx context.Context, req *req.PayOrderExportReq) ([]*pay.PayOrder, error) {
	q := s.q.PayOrder.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayOrder.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayOrder.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayOrder.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.Subject != "" {
		q = q.Where(s.q.PayOrder.Subject.Like("%" + req.Subject + "%"))
	}
	if req.No != "" {
		q = q.Where(s.q.PayOrder.No.Eq(req.No))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayOrder.Status.Eq(*req.Status))
	}
	return q.Order(s.q.PayOrder.ID.Desc()).Find()
}
