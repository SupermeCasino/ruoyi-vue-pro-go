package trade

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
)

// 消息模板常量 (对齐 Java: MessageTemplateConstants)
const (
	// SMSOrderDelivery 订单发货短信模板编号
	SMSOrderDelivery = "order_delivery"
	// WXAOrderDelivery 小程序订阅消息模版标题
	WXAOrderDelivery = "订单发货通知"
)

// TradeOrderMessageWhenDeliveryOrderReq 订单发货时通知创建请求 (对齐 Java: TradeOrderMessageWhenDeliveryOrderReqBO)
type TradeOrderMessageWhenDeliveryOrderReq struct {
	OrderID int64  // 订单编号
	UserID  int64  // 用户编号
	Message string // 消息
}

// TradeMessageService Trade 消息 service (对齐 Java: TradeMessageService)
type TradeMessageService struct {
	notifySvc *service.NotifyService
}

// NewTradeMessageService 创建 TradeMessageService
func NewTradeMessageService(notifySvc *service.NotifyService) *TradeMessageService {
	return &TradeMessageService{
		notifySvc: notifySvc,
	}
}

// SendMessageWhenDeliveryOrder 订单发货时发送通知 (对齐 Java: TradeMessageServiceImpl.sendMessageWhenDeliveryOrder)
func (s *TradeMessageService) SendMessageWhenDeliveryOrder(ctx context.Context, req *TradeOrderMessageWhenDeliveryOrderReq) {

	// 1、构造消息参数
	msgMap := map[string]interface{}{
		"orderId":         req.OrderID,
		"deliveryMessage": req.Message,
	}

	// 2、发送站内信 (对齐 Java: notifyMessageSendApi.sendSingleMessageToMember)
	// userType = 1 表示会员用户 (model.UserTypeMember)
	_, _ = s.notifySvc.SendNotify(ctx, req.UserID, consts.UserTypeMember, SMSOrderDelivery, msgMap)
}
