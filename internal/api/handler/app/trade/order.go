package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/trade"

	"github.com/gin-gonic/gin"
)

type AppTradeOrderHandler struct {
	svc      *trade.TradeOrderUpdateService
	querySvc *trade.TradeOrderQueryService
}

func NewAppTradeOrderHandler(svc *trade.TradeOrderUpdateService, querySvc *trade.TradeOrderQueryService) *AppTradeOrderHandler {
	return &AppTradeOrderHandler{
		svc:      svc,
		querySvc: querySvc,
	}
}

// SettlementOrder 获得订单结算信息
func (h *AppTradeOrderHandler) SettlementOrder(c *gin.Context) {
	var r req.AppTradeOrderSettlementReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}

	// GET request with query params mapped to JSON struct?
	// Usually Settlement is POST if complex, or GET with query.
	// Java: GET /settlement, Valid AppTradeOrderSettlementReqVO. Complex objects in GET are tricky.
	// But let's assume it accepts query params or we change to POST if needed.
	// If query, use ShouldBindQuery. But items list is hard in query.
	// Let's check Java Controller. @GetMapping("/settlement").
	// Frontend likely sends items[0].skuId=...

	// For simplicity, if complex, maybe POST is better, but stick to GET if we can bind.
	// Gin's BindQuery can handle array/map.
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}

	res, err := h.svc.SettlementOrder(c, core.GetUserId(c), &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// CreateOrder 创建订单
func (h *AppTradeOrderHandler) CreateOrder(c *gin.Context) {
	var r req.AppTradeOrderCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.CreateOrder(c, core.GetUserId(c), &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	// Return ID and PayOrderID (mocked)
	core.WriteSuccess(c, map[string]interface{}{
		"id":         res.ID,
		"payOrderId": res.PayOrderID,
	})
}

// GetOrderDetail 获得订单详情
func (h *AppTradeOrderHandler) GetOrderDetail(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	// 1. Get Order
	order, err := h.querySvc.GetOrder(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	if order == nil || order.UserID != core.GetUserId(c) {
		core.WriteBizError(c, core.ErrNotFound)
		return
	}

	// 2. Get Items
	items, err := h.querySvc.GetOrderItemListByOrderId(c, order.ID)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	// 3. Assemble DTO
	itemResps := make([]resp.AppTradeOrderItemResp, len(items))
	for i, item := range items {
		itemResps[i] = resp.AppTradeOrderItemResp{
			ID:            item.ID,
			OrderID:       item.OrderID,
			SpuID:         item.SpuID,
			SpuName:       item.SpuName,
			SkuID:         item.SkuID,
			PicURL:        item.PicURL,
			Count:         item.Count,
			CommentStatus: item.CommentStatus,
			Price:         item.Price,
			PayPrice:      item.PayPrice,
			// Properties: item.Properties (Need deserialize if stored as JSON/String)
		}
	}

	var payOrderID int64
	if order.PayOrderID != nil {
		payOrderID = *order.PayOrderID
	}

	res := resp.AppTradeOrderDetailResp{
		ID:                    order.ID,
		No:                    order.No,
		Type:                  order.Type,
		CreateTime:            order.CreatedAt,
		UserRemark:            order.UserRemark,
		Status:                order.Status,
		ProductCount:          order.ProductCount,
		FinishTime:            order.FinishTime,
		CancelTime:            order.CancelTime,
		CommentStatus:         order.CommentStatus,
		PayStatus:             order.PayStatus,
		PayOrderID:            payOrderID,
		PayTime:               order.PayTime,
		PayChannelCode:        order.PayChannelCode,
		TotalPrice:            order.TotalPrice,
		DiscountPrice:         order.DiscountPrice,
		DeliveryPrice:         order.DeliveryPrice,
		AdjustPrice:           order.AdjustPrice,
		PayPrice:              order.PayPrice,
		DeliveryType:          order.DeliveryType,
		LogisticsID:           order.LogisticsID,
		LogisticsNo:           order.LogisticsNo,
		DeliveryTime:          order.DeliveryTime,
		ReceiveTime:           order.ReceiveTime,
		ReceiverName:          order.ReceiverName,
		ReceiverMobile:        order.ReceiverMobile,
		ReceiverAreaID:        order.ReceiverAreaID,
		ReceiverDetailAddress: order.ReceiverDetailAddress,
		RefundStatus:          order.RefundStatus,
		RefundPrice:           order.RefundPrice,
		CouponID:              order.CouponID,
		CouponPrice:           order.CouponPrice,
		Items:                 itemResps,
	}

	core.WriteSuccess(c, res)
}

// GetOrderPage 获得订单分页
func (h *AppTradeOrderHandler) GetOrderPage(c *gin.Context) {
	var r req.AppTradeOrderPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	// 1. Get Page
	pageResult, err := h.querySvc.GetOrderPage(c, core.GetUserId(c), &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// 2. Get Items for all orders
	if len(pageResult.List) > 0 {
		orderIds := make([]int64, len(pageResult.List))
		for i, o := range pageResult.List {
			orderIds[i] = o.ID
		}
		items, err := h.querySvc.GetOrderItemListByOrderIds(c, orderIds)
		if err != nil {
			// Log error but continue? Or fail.
			core.WriteError(c, 500, err.Error())
			return
		}

		// Map items by OrderID
		itemMap := make(map[int64][]resp.AppTradeOrderItemResp)
		for _, item := range items {
			itemResp := resp.AppTradeOrderItemResp{
				ID:            item.ID,
				OrderID:       item.OrderID,
				SpuID:         item.SpuID,
				SpuName:       item.SpuName,
				SkuID:         item.SkuID,
				PicURL:        item.PicURL,
				Count:         item.Count,
				CommentStatus: item.CommentStatus,
				Price:         item.Price,
				PayPrice:      item.PayPrice,
			}
			itemMap[item.OrderID] = append(itemMap[item.OrderID], itemResp)
		}

		// 3. Assemble DTO List
		list := make([]resp.AppTradeOrderPageItemResp, len(pageResult.List))
		for i, o := range pageResult.List {
			var payOrderID int64
			if o.PayOrderID != nil {
				payOrderID = *o.PayOrderID
			}
			list[i] = resp.AppTradeOrderPageItemResp{
				ID:            o.ID,
				No:            o.No,
				Type:          o.Type,
				Status:        o.Status,
				ProductCount:  o.ProductCount,
				CommentStatus: o.CommentStatus,
				CreateTime:    o.CreatedAt,
				PayOrderID:    payOrderID,
				PayPrice:      o.PayPrice,
				DeliveryType:  o.DeliveryType,
				Items:         itemMap[o.ID],
			}
		}

		core.WriteSuccess(c, core.PageResult[resp.AppTradeOrderPageItemResp]{
			List:  list,
			Total: pageResult.Total,
		})
	} else {
		core.WriteSuccess(c, core.PageResult[resp.AppTradeOrderPageItemResp]{
			List:  []resp.AppTradeOrderPageItemResp{},
			Total: 0,
		})
	}
}

// CancelOrder 取消订单
func (h *AppTradeOrderHandler) CancelOrder(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "id is required")
		return
	}
	err := h.svc.CancelOrder(c, core.GetUserId(c), id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteOrder 删除订单
func (h *AppTradeOrderHandler) DeleteOrder(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "id is required")
		return
	}
	err := h.svc.DeleteOrder(c, core.GetUserId(c), id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetOrderCount 获得交易订单数量
func (h *AppTradeOrderHandler) GetOrderCount(c *gin.Context) {
	// Simple map return as per Java
	uid := core.GetUserId(c)
	res := make(map[string]int64)

	// Helper to get count
	getCount := func(status *int, commentStatus *bool) int64 {
		count, _ := h.querySvc.GetOrderCount(c, uid, status, commentStatus)
		return count
	}

	// Constants (Should use Enum in real, here hardcoded for migration speed matching Java Enum)
	// UNPAID(0), UNDELIVERED(10), DELIVERED(20), COMPLETED(30), CANCELLED(40)
	unpaid := 0
	undelivered := 10
	delivered := 20
	completed := 30

	res["allCount"] = getCount(nil, nil)
	res["unpaidCount"] = getCount(&unpaid, nil)
	res["undeliveredCount"] = getCount(&undelivered, nil)
	res["deliveredCount"] = getCount(&delivered, nil)

	commentStatus := false
	res["uncommentedCount"] = getCount(&completed, &commentStatus)

	// TODO: AfterSale count (requires AfterSaleService)
	res["afterSaleCount"] = 0

	core.WriteSuccess(c, res)
}

// CreateOrderItemComment 创建订单项评价
func (h *AppTradeOrderHandler) CreateOrderItemComment(c *gin.Context) {
	var r req.AppTradeOrderItemCommentCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.svc.CreateOrderItemCommentByMember(c, core.GetUserId(c), &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}

// GetOrderExpressTrackList 获得物流轨迹
func (h *AppTradeOrderHandler) GetOrderExpressTrackList(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "id is required")
		return
	}
	res, err := h.querySvc.GetExpressTrackList(c, id, core.GetUserId(c))
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}

// ReceiveOrder 确认收货
func (h *AppTradeOrderHandler) ReceiveOrder(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "id is required")
		return
	}
	err := h.svc.ReceiveOrder(c, core.GetUserId(c), id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}
