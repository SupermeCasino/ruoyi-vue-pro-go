package trade

import (
	"context"

	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
)

// OrderHandler 订单处理器接口
type OrderHandler interface {
	// Handle 处理订单操作
	Handle(ctx context.Context, handleReq *OrderHandleRequest) (*OrderHandleResponse, error)

	// GetHandlerType 获取处理器类型
	GetHandlerType() string

	// CanHandle 判断是否能处理指定操作
	CanHandle(operation string) bool

	// PreCheck 前置检查
	PreCheck(ctx context.Context, handleReq *OrderHandleRequest) error

	// PostProcess 后置处理
	PostProcess(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error

	// Lifecycle Hooks
	BeforeOrderCreate(ctx context.Context, handleReq *OrderHandleRequest) error
	AfterOrderCreate(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error
	AfterPayOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error
	AfterCancelOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error
	AfterReceiveOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error
	AfterDeliveryOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error
}

// OrderHandleRequest 订单处理请求
type OrderHandleRequest struct {
	Operation string `json:"operation"` // 操作类型

	// 通用字段
	UserID  int64  `json:"userId"`  // 用户ID
	OrderID int64  `json:"orderId"` // 订单ID
	AdminID int64  `json:"adminId"` // 管理员ID
	Remark  string `json:"remark"`  // 备注

	// 创建订单相关
	CreateReq         *trade2.AppTradeOrderCreateReq `json:"createReq"`         // 创建订单请求
	PriceCalculateReq *TradePriceCalculateReqBO      `json:"priceCalculateReq"` // 价格计算请求
	CartIDs           []int64                        `json:"cartIds"`           // 购物车ID数组
	StockItems        []StockDeductItem              `json:"stockItems"`        // 库存扣减项
	OrderItems        []*trade.TradeOrderItem        `json:"orderItems"`        // 订单项数组

	// 支付订单相关
	PayOrderID int64 `json:"payOrderId"` // 支付订单ID

	// 发货相关
	LogisticsID int64  `json:"logisticsId"` // 物流公司ID
	TrackingNo  string `json:"trackingNo"`  // 物流单号

	// 取消订单相关
	CancelType   int    `json:"cancelType"`   // 取消类型
	CancelReason string `json:"cancelReason"` // 取消原因

	// 退款相关
	RefundPrice  int64  `json:"refundPrice"`  // 退款金额
	RefundReason string `json:"refundReason"` // 退款原因
}

// OrderHandleResponse 订单处理响应
type OrderHandleResponse struct {
	Order      *trade.TradeOrder      `json:"order"`      // 订单信息
	PayOrderID int64                  `json:"payOrderId"` // 支付订单ID
	Success    bool                   `json:"success"`    // 处理是否成功
	Message    string                 `json:"message"`    // 处理消息
	Data       map[string]interface{} `json:"data"`       // 额外数据
}

// StockDeductItem 库存扣减项
type StockDeductItem struct {
	SkuID int64 `json:"skuId"` // 商品SKU ID
	Count int   `json:"count"` // 扣减数量
}

// BaseOrderHandler 订单处理器基类
type BaseOrderHandler struct {
	handlerType string
	operations  []string
}

// NewBaseOrderHandler 创建订单处理器基类
func NewBaseOrderHandler(handlerType string, operations []string) *BaseOrderHandler {
	return &BaseOrderHandler{
		handlerType: handlerType,
		operations:  operations,
	}
}

// GetHandlerType 获取处理器类型
func (h *BaseOrderHandler) GetHandlerType() string {
	return h.handlerType
}

// CanHandle 判断是否能处理指定操作
func (h *BaseOrderHandler) CanHandle(operation string) bool {
	for _, op := range h.operations {
		if op == operation {
			return true
		}
	}
	return false
}

// PreCheck 前置检查（默认实现）
func (h *BaseOrderHandler) PreCheck(ctx context.Context, handleReq *OrderHandleRequest) error {
	// 默认不做检查
	return nil
}

// PostProcess 后置处理（默认实现）
func (h *BaseOrderHandler) PostProcess(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error {
	// 默认不做处理
	return nil
}

// BeforeOrderCreate 订单创建前置处理
func (h *BaseOrderHandler) BeforeOrderCreate(ctx context.Context, handleReq *OrderHandleRequest) error {
	return nil
}

// AfterOrderCreate 订单创建后置处理
func (h *BaseOrderHandler) AfterOrderCreate(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error {
	return nil
}

// AfterPayOrder 订单支付后置处理
func (h *BaseOrderHandler) AfterPayOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error {
	return nil
}

// AfterCancelOrder 订单取消后置处理
func (h *BaseOrderHandler) AfterCancelOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error {
	return nil
}

// AfterReceiveOrder 订单收货后置处理
func (h *BaseOrderHandler) AfterReceiveOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error {
	return nil
}

// AfterDeliveryOrder 订单发货后置处理
func (h *BaseOrderHandler) AfterDeliveryOrder(ctx context.Context, handleReq *OrderHandleRequest, resp *OrderHandleResponse) error {
	return nil
}
