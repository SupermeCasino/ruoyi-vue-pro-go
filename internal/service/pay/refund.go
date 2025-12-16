package pay

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayRefundService struct {
	q          *query.Query
	channelSvc *PayChannelService
	orderSvc   *PayOrderService
	notifySvc  *PayNotifyService
}

func NewPayRefundService(q *query.Query, channelSvc *PayChannelService, orderSvc *PayOrderService, notifySvc *PayNotifyService) *PayRefundService {
	return &PayRefundService{
		q:          q,
		channelSvc: channelSvc,
		orderSvc:   orderSvc,
		notifySvc:  notifySvc,
	}
}

// GetRefund 获得退款订单
func (s *PayRefundService) GetRefund(ctx context.Context, id int64) (*pay.PayRefund, error) {
	return s.q.PayRefund.WithContext(ctx).Where(s.q.PayRefund.ID.Eq(id)).First()
}

// GetRefundPage 获得退款订单分页
func (s *PayRefundService) GetRefundPage(ctx context.Context, req *req.PayRefundPageReq) (*pagination.PageResult[*pay.PayRefund], error) {
	q := s.q.PayRefund.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayRefund.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayRefund.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayRefund.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.MerchantRefundId != "" {
		q = q.Where(s.q.PayRefund.MerchantRefundId.Eq(req.MerchantRefundId))
	}
	if req.ChannelOrderNo != "" {
		q = q.Where(s.q.PayRefund.ChannelOrderNo.Eq(req.ChannelOrderNo))
	}
	if req.ChannelRefundNo != "" {
		q = q.Where(s.q.PayRefund.ChannelRefundNo.Eq(req.ChannelRefundNo))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayRefund.Status.Eq(*req.Status))
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}
	list, err := q.Limit(req.GetLimit()).Offset(req.GetOffset()).Order(s.q.PayRefund.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[*pay.PayRefund]{
		List:  list,
		Total: total,
	}, nil
}

// GetRefundList 获得退款订单列表 (Export)
func (s *PayRefundService) GetRefundList(ctx context.Context, req *req.PayRefundExportReq) ([]*pay.PayRefund, error) {
	q := s.q.PayRefund.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayRefund.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayRefund.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayRefund.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.MerchantRefundId != "" {
		q = q.Where(s.q.PayRefund.MerchantRefundId.Eq(req.MerchantRefundId))
	}
	if req.ChannelOrderNo != "" {
		q = q.Where(s.q.PayRefund.ChannelOrderNo.Eq(req.ChannelOrderNo))
	}
	if req.ChannelRefundNo != "" {
		q = q.Where(s.q.PayRefund.ChannelRefundNo.Eq(req.ChannelRefundNo))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayRefund.Status.Eq(*req.Status))
	}
	return q.Order(s.q.PayRefund.ID.Desc()).Find()
}

// NotifyRefund 处理退款回调通知
// 对齐 Java: PayRefundService.notifyRefund(Long channelId, PayRefundRespDTO notify)
func (s *PayRefundService) NotifyRefund(ctx context.Context, channelID int64, notify *client.RefundResp) error {
	// 校验支付渠道是否有效
	channel, err := s.channelSvc.ValidPayChannel(ctx, channelID)
	if err != nil {
		return err
	}

	// 情况一：退款成功
	if notify.Status == PayRefundStatusSuccess {
		return s.notifyRefundSuccess(ctx, channel, notify)
	}

	// 情况二：退款失败
	if notify.Status == PayRefundStatusFailure {
		return s.notifyRefundFailure(ctx, channel, notify)
	}

	return nil
}

// notifyRefundSuccess 处理退款成功
func (s *PayRefundService) notifyRefundSuccess(ctx context.Context, channel *pay.PayChannel, notify *client.RefundResp) error {
	// 1.1 查询 PayRefund
	refund, err := s.q.PayRefund.WithContext(ctx).
		Where(s.q.PayRefund.AppID.Eq(channel.AppID), s.q.PayRefund.No.Eq(notify.OutRefundNo)).
		First()
	if err != nil {
		return fmt.Errorf("退款订单不存在")
	}

	// 如果已经是成功，直接返回
	if refund.Status == PayRefundStatusSuccess {
		return nil
	}

	// 校验状态，必须是等待状态
	if refund.Status != PayRefundStatusWaiting {
		return fmt.Errorf("退款订单状态不是待退款")
	}

	// 1.2 更新 PayRefund (使用乐观锁)
	notifyDataJSON, _ := json.Marshal(notify)
	result, err := s.q.PayRefund.WithContext(ctx).
		Where(s.q.PayRefund.ID.Eq(refund.ID), s.q.PayRefund.Status.Eq(PayRefundStatusWaiting)).
		Updates(map[string]interface{}{
			"status":              PayRefundStatusSuccess,
			"success_time":        notify.SuccessTime,
			"channel_refund_no":   notify.ChannelRefundNo,
			"channel_notify_data": string(notifyDataJSON),
		})

	if err != nil || result.RowsAffected == 0 {
		return fmt.Errorf("退款订单状态不是待退款")
	}

	// 2. 更新订单退款金额
	if err := s.orderSvc.UpdateOrderRefundPrice(ctx, refund.OrderID, refund.RefundPrice); err != nil {
		return err
	}

	// 3. 插入退款通知记录
	s.notifySvc.CreatePayNotifyTask(ctx, PayNotifyTypeRefund, refund.ID)

	return nil
}

// notifyRefundFailure 处理退款失败
func (s *PayRefundService) notifyRefundFailure(ctx context.Context, channel *pay.PayChannel, notify *client.RefundResp) error {
	// 1.1 查询 PayRefund
	refund, err := s.q.PayRefund.WithContext(ctx).
		Where(s.q.PayRefund.AppID.Eq(channel.AppID), s.q.PayRefund.No.Eq(notify.OutRefundNo)).
		First()
	if err != nil {
		return fmt.Errorf("退款订单不存在")
	}

	// 如果已经是失败，直接返回
	if refund.Status == PayRefundStatusFailure {
		return nil
	}

	// 校验状态，必须是等待状态
	if refund.Status != PayRefundStatusWaiting {
		return fmt.Errorf("退款订单状态不是待退款")
	}

	// 1.2 更新 PayRefund (使用乐观锁)
	notifyDataJSON, _ := json.Marshal(notify)
	result, err := s.q.PayRefund.WithContext(ctx).
		Where(s.q.PayRefund.ID.Eq(refund.ID), s.q.PayRefund.Status.Eq(PayRefundStatusWaiting)).
		Updates(map[string]interface{}{
			"status":              PayRefundStatusFailure,
			"channel_refund_no":   notify.ChannelRefundNo,
			"channel_notify_data": string(notifyDataJSON),
			"channel_error_code":  notify.ChannelErrorCode,
			"channel_error_msg":   notify.ChannelErrorMsg,
		})

	if err != nil || result.RowsAffected == 0 {
		return fmt.Errorf("退款订单状态不是待退款")
	}

	// 2. 插入退款通知记录
	s.notifySvc.CreatePayNotifyTask(ctx, PayNotifyTypeRefund, refund.ID)

	return nil
}
