package trade

import (
	"context"
	"fmt"
	"time"

	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"go.uber.org/zap"
)

// OrderHandlerRouter 订单处理器路由
type OrderHandlerRouter struct {
	factory *OrderHandlerFactory
	logger  *zap.Logger
}

// NewOrderHandlerRouter 创建订单处理器路由
func NewOrderHandlerRouter(factory *OrderHandlerFactory, logger *zap.Logger) *OrderHandlerRouter {
	return &OrderHandlerRouter{
		factory: factory,
		logger:  logger,
	}
}

// Route 路由处理订单操作
func (r *OrderHandlerRouter) Route(ctx context.Context, req *OrderHandleRequest) (*OrderHandleResponse, error) {
	startTime := time.Now()

	// 记录开始处理日志
	r.logger.Info("开始处理订单操作",
		zap.String("operation", req.Operation),
		zap.Int64("userId", req.UserID),
		zap.Int64("orderId", req.OrderID),
	)

	// 获取处理器
	handler, err := r.factory.GetHandlerByOperation(req.Operation)
	if err != nil {
		r.logger.Error("获取订单处理器失败",
			zap.String("operation", req.Operation),
			zap.Error(err),
		)
		return nil, fmt.Errorf("不支持的订单操作: %s", req.Operation)
	}

	// 执行处理器
	resp, err := handler.Handle(ctx, req)
	if err != nil {
		r.logger.Error("订单处理器执行失败",
			zap.String("operation", req.Operation),
			zap.String("handlerType", handler.GetHandlerType()),
			zap.Error(err),
			zap.Duration("duration", time.Since(startTime)),
		)
		return nil, err
	}

	// 记录完成处理日志
	r.logger.Info("订单操作处理完成",
		zap.String("operation", req.Operation),
		zap.String("handlerType", handler.GetHandlerType()),
		zap.Bool("success", resp.Success),
		zap.Duration("duration", time.Since(startTime)),
	)

	return resp, nil
}

// RouteCreateOrder 路由创建订单操作
func (r *OrderHandlerRouter) RouteCreateOrder(ctx context.Context, userID int64, createReq *OrderHandleRequest) (*OrderHandleResponse, error) {
	createReq.Operation = GetOrderOperationName(tradeModel.OrderOperateTypeCreate)
	createReq.UserID = userID

	return r.Route(ctx, createReq)
}

// RoutePayOrder 路由支付订单操作
func (r *OrderHandlerRouter) RoutePayOrder(ctx context.Context, orderID int64, payOrderID int64) (*OrderHandleResponse, error) {
	req := &OrderHandleRequest{
		Operation:  GetOrderOperationName(tradeModel.OrderOperateTypePay),
		OrderID:    orderID,
		PayOrderID: payOrderID,
	}

	return r.Route(ctx, req)
}

// RouteDeliveryOrder 路由发货订单操作
func (r *OrderHandlerRouter) RouteDeliveryOrder(ctx context.Context, orderID int64, logisticsID int64, trackingNo string) (*OrderHandleResponse, error) {
	req := &OrderHandleRequest{
		Operation:   GetOrderOperationName(tradeModel.OrderOperateTypeDelivery),
		OrderID:     orderID,
		LogisticsID: logisticsID,
		TrackingNo:  trackingNo,
	}

	return r.Route(ctx, req)
}

// RouteReceiveOrder 路由确认收货操作
func (r *OrderHandlerRouter) RouteReceiveOrder(ctx context.Context, userID int64, orderID int64) (*OrderHandleResponse, error) {
	req := &OrderHandleRequest{
		Operation: GetOrderOperationName(tradeModel.OrderOperateTypeReceive),
		UserID:    userID,
		OrderID:   orderID,
	}

	return r.Route(ctx, req)
}

// RouteCancelOrder 路由取消订单操作
func (r *OrderHandlerRouter) RouteCancelOrder(ctx context.Context, userID int64, orderID int64, cancelType int, cancelReason string) (*OrderHandleResponse, error) {
	req := &OrderHandleRequest{
		Operation:    GetOrderOperationName(tradeModel.OrderOperateTypeCancel),
		UserID:       userID,
		OrderID:      orderID,
		CancelType:   cancelType,
		CancelReason: cancelReason,
	}

	return r.Route(ctx, req)
}

// RouteRefundOrder 路由退款订单操作
func (r *OrderHandlerRouter) RouteRefundOrder(ctx context.Context, orderID int64, refundPrice int64, refundReason string) (*OrderHandleResponse, error) {
	req := &OrderHandleRequest{
		Operation:    GetOrderOperationName(tradeModel.OrderOperateTypeRefund),
		OrderID:      orderID,
		RefundPrice:  refundPrice,
		RefundReason: refundReason,
	}

	return r.Route(ctx, req)
}

// GetOrderOperationName 获取订单操作名称
func GetOrderOperationName(operationType int) string {
	switch operationType {
	case tradeModel.OrderOperateTypeCreate:
		return "create"
	case tradeModel.OrderOperateTypePay:
		return "pay"
	case tradeModel.OrderOperateTypeDelivery:
		return "delivery"
	case tradeModel.OrderOperateTypeReceive:
		return "receive"
	case tradeModel.OrderOperateTypePickUp:
		return "pickup"
	case tradeModel.OrderOperateTypeCancel:
		return "cancel"
	case tradeModel.OrderOperateTypeRefund:
		return "refund"
	default:
		return "unknown"
	}
}
