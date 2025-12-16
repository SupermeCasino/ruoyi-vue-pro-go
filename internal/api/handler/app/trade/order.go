package trade

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AppTradeOrderHandler struct {
	svc          *trade.TradeOrderUpdateService
	querySvc     *trade.TradeOrderQueryService
	afterSaleSvc *trade.TradeAfterSaleService
	priceSvc     *trade.TradePriceService
}

func NewAppTradeOrderHandler(
	svc *trade.TradeOrderUpdateService,
	querySvc *trade.TradeOrderQueryService,
	afterSaleSvc *trade.TradeAfterSaleService,
	priceSvc *trade.TradePriceService,
) *AppTradeOrderHandler {
	return &AppTradeOrderHandler{
		svc:          svc,
		querySvc:     querySvc,
		afterSaleSvc: afterSaleSvc,
		priceSvc:     priceSvc,
	}
}

// SettlementOrder 获得订单结算信息
func (h *AppTradeOrderHandler) SettlementOrder(c *gin.Context) {
	var r req.AppTradeOrderSettlementReq
	// 支持 POST (JSON) 和 GET (Query) 以对齐多样化的 App/Web 使用场景
	if c.Request.Method == "POST" {
		if err := c.ShouldBindJSON(&r); err != nil {
			response.WriteBizError(c, errors.ErrParam)
			return
		}
	} else {
		if err := c.ShouldBindQuery(&r); err != nil {
			response.WriteBizError(c, errors.ErrParam)
			return
		}
	}

	res, err := h.svc.SettlementOrder(c, context.GetUserId(c), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// CreateOrder 创建订单
func (h *AppTradeOrderHandler) CreateOrder(c *gin.Context) {
	var r req.AppTradeOrderCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.CreateOrder(c, context.GetUserId(c), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 严格对齐 Java: 返回 AppTradeOrderCreateResp
	vo := resp.AppTradeOrderCreateResp{
		ID:         res.ID,
		PayOrderID: 0, // 默认 0
	}
	if res.PayOrderID != nil {
		vo.PayOrderID = *res.PayOrderID
	}

	response.WriteSuccess(c, vo)
}

// UpdateOrderPaid 更新订单为已支付
// 这是一个回调接口,通常由 Pay 模块通过 HTTP 调用
func (h *AppTradeOrderHandler) UpdateOrderPaid(c *gin.Context) {
	var r req.PayOrderNotifyReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	id := utils.ParseInt64(r.MerchantOrderId)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	err := h.svc.UpdateOrderPaid(c, id, r.PayOrderID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetOrderDetail 获得订单详情
func (h *AppTradeOrderHandler) GetOrderDetail(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 1. 获得订单
	order, err := h.querySvc.GetOrder(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if order == nil || order.UserID != context.GetUserId(c) {
		response.WriteBizError(c, errors.ErrNotFound)
		return
	}

	// 2. 获得订单项
	items, err := h.querySvc.GetOrderItemListByOrderId(c, order.ID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 3. 拼接数据
	itemResps := make([]resp.AppTradeOrderItemResp, len(items))
	for i, item := range items {
		properties := make([]resp.ProductSkuPropertyResp, len(item.Properties))
		for j, p := range item.Properties {
			properties[j] = resp.ProductSkuPropertyResp{
				PropertyID:   p.PropertyID,
				PropertyName: p.PropertyName,
				ValueID:      p.ValueID,
				ValueName:    p.ValueName,
			}
		}

		itemResps[i] = resp.AppTradeOrderItemResp{
			ID:              item.ID,
			OrderID:         item.OrderID,
			SpuID:           item.SpuID,
			SpuName:         item.SpuName,
			SkuID:           item.SkuID,
			PicURL:          item.PicURL,
			Count:           item.Count,
			CommentStatus:   bool(item.CommentStatus),
			Price:           item.Price,
			PayPrice:        item.PayPrice,
			AfterSaleID:     item.AfterSaleID,
			AfterSaleStatus: item.AfterSaleStatus,
			Properties:      properties,
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
		CommentStatus:         bool(order.CommentStatus),
		PayStatus:             bool(order.PayStatus),
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

	response.WriteSuccess(c, res)
}

// GetOrderPage 获得订单分页
func (h *AppTradeOrderHandler) GetOrderPage(c *gin.Context) {
	var r req.AppTradeOrderPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	// 1. 获得分页列表
	pageResult, err := h.querySvc.GetOrderPage(c, context.GetUserId(c), &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	// 2. 获得订单项
	if len(pageResult.List) > 0 {
		orderIds := make([]int64, len(pageResult.List))
		for i, o := range pageResult.List {
			orderIds[i] = o.ID
		}
		items, err := h.querySvc.GetOrderItemListByOrderIds(c, orderIds)
		if err != nil {
			// 记录日志，但此处继续返回，或返回错误
			response.WriteError(c, 500, err.Error())
			return
		}

		// 根据 OrderID 映射订单项
		itemMap := make(map[int64][]resp.AppTradeOrderItemResp)
		for _, item := range items {
			properties := make([]resp.ProductSkuPropertyResp, len(item.Properties))
			for j, p := range item.Properties {
				properties[j] = resp.ProductSkuPropertyResp{
					PropertyID:   p.PropertyID,
					PropertyName: p.PropertyName,
					ValueID:      p.ValueID,
					ValueName:    p.ValueName,
				}
			}

			itemResp := resp.AppTradeOrderItemResp{
				ID:              item.ID,
				OrderID:         item.OrderID,
				SpuID:           item.SpuID,
				SpuName:         item.SpuName,
				SkuID:           item.SkuID,
				PicURL:          item.PicURL,
				Count:           item.Count,
				CommentStatus:   bool(item.CommentStatus),
				Price:           item.Price,
				PayPrice:        item.PayPrice,
				AfterSaleID:     item.AfterSaleID,
				AfterSaleStatus: item.AfterSaleStatus,
				Properties:      properties,
			}
			itemMap[item.OrderID] = append(itemMap[item.OrderID], itemResp)
		}

		// 3. 拼接返回 VO
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
				CommentStatus: bool(o.CommentStatus),
				CreateTime:    o.CreatedAt,
				PayOrderID:    payOrderID,
				PayPrice:      o.PayPrice,
				DeliveryType:  o.DeliveryType,
				Items:         itemMap[o.ID],
			}
		}

		response.WriteSuccess(c, pagination.PageResult[resp.AppTradeOrderPageItemResp]{
			List:  list,
			Total: pageResult.Total,
		})
	} else {
		response.WriteSuccess(c, pagination.PageResult[resp.AppTradeOrderPageItemResp]{
			List:  []resp.AppTradeOrderPageItemResp{},
			Total: 0,
		})
	}
}

// CancelOrder 取消订单
func (h *AppTradeOrderHandler) CancelOrder(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "id is required")
		return
	}
	err := h.svc.CancelOrder(c, context.GetUserId(c), id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteOrder 删除订单
func (h *AppTradeOrderHandler) DeleteOrder(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "id is required")
		return
	}
	err := h.svc.DeleteOrder(c, context.GetUserId(c), id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// GetOrderCount 获得交易订单数量
func (h *AppTradeOrderHandler) GetOrderCount(c *gin.Context) {
	// 简单的 Map 返回，对齐 Java
	uid := context.GetUserId(c)
	res := make(map[string]int64)

	// 获取数量辅助函数
	getCount := func(status *int, commentStatus *bool) int64 {
		count, _ := h.querySvc.GetOrderCount(c, uid, status, commentStatus)
		return count
	}

	// 常量
	// UNPAID(0), UNDELIVERED(10), DELIVERED(20), COMPLETED(30), CANCELLED(40)
	unpaid := tradeModel.TradeOrderStatusUnpaid
	undelivered := tradeModel.TradeOrderStatusUndelivered
	delivered := tradeModel.TradeOrderStatusDelivered
	completed := tradeModel.TradeOrderStatusCompleted
	commentStatus := false

	res["allCount"] = getCount(nil, nil)
	res["unpaidCount"] = getCount(&unpaid, nil)
	res["undeliveredCount"] = getCount(&undelivered, nil)
	res["deliveredCount"] = getCount(&delivered, nil)
	res["uncommentedCount"] = getCount(&completed, &commentStatus)

	// 售后数量 (需要 TradeAfterSaleService)
	afterSaleCount, _ := h.afterSaleSvc.GetUserAfterSaleCount(c, uid)
	res["afterSaleCount"] = afterSaleCount

	response.WriteSuccess(c, res)
}

// CreateOrderItemComment 创建订单项评价
func (h *AppTradeOrderHandler) CreateOrderItemComment(c *gin.Context) {
	var r req.AppTradeOrderItemCommentCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.svc.CreateOrderItemCommentByMember(c, context.GetUserId(c), &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, res)
}

// GetOrderExpressTrackList 获得物流轨迹
func (h *AppTradeOrderHandler) GetOrderExpressTrackList(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "id is required")
		return
	}
	res, err := h.querySvc.GetExpressTrackList(c, id, context.GetUserId(c))
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, res)
}

// ReceiveOrder 确认收货
func (h *AppTradeOrderHandler) ReceiveOrder(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "id is required")
		return
	}
	err := h.svc.ReceiveOrder(c, context.GetUserId(c), id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// SettlementProduct 获得商品结算信息
func (h *AppTradeOrderHandler) SettlementProduct(c *gin.Context) {
	// spuIds=1,2,3 or spuIds=1&spuIds=2
	var req struct {
		SpuIDs []int64 `form:"spuIds"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam) // FIXED: syntax error previously here? No, duplicate WriteSuccess logic in partial
		return
	}
	if len(req.SpuIDs) == 0 {
		// Try parsing from string "spuIds" if user sends "1,2,3"
		str := c.Query("spuIds")
		if str != "" {
			req.SpuIDs = utils.SplitToInt64(str)
		}
	}

	res, err := h.priceSvc.CalculateProductPrice(c, context.GetUserId(c), req.SpuIDs)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetOrderItem 获得交易订单项
func (h *AppTradeOrderHandler) GetOrderItem(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	item, err := h.querySvc.GetOrderItem(c, context.GetUserId(c), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	properties := make([]resp.ProductSkuPropertyResp, len(item.Properties))
	for j, p := range item.Properties {
		properties[j] = resp.ProductSkuPropertyResp{
			PropertyID:   p.PropertyID,
			PropertyName: p.PropertyName,
			ValueID:      p.ValueID,
			ValueName:    p.ValueName,
		}
	}

	res := resp.AppTradeOrderItemResp{
		ID:              item.ID,
		OrderID:         item.OrderID,
		SpuID:           item.SpuID,
		SpuName:         item.SpuName,
		SkuID:           item.SkuID,
		PicURL:          item.PicURL,
		Count:           item.Count,
		CommentStatus:   bool(item.CommentStatus),
		Price:           item.Price,
		PayPrice:        item.PayPrice,
		AfterSaleID:     item.AfterSaleID,
		AfterSaleStatus: item.AfterSaleStatus,
		Properties:      properties,
	}
	response.WriteSuccess(c, res)
}
