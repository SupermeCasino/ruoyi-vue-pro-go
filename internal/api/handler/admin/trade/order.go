package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/member"
	"backend-go/internal/service/trade"

	"github.com/gin-gonic/gin"
)

type TradeOrderHandler struct {
	svc                        *trade.TradeOrderUpdateService
	querySvc                   *trade.TradeOrderQueryService
	memberSvc                  *member.MemberUserService
	deliveryFreightTemplateSvc *trade.DeliveryExpressTemplateService
}

func NewTradeOrderHandler(svc *trade.TradeOrderUpdateService, querySvc *trade.TradeOrderQueryService, memberSvc *member.MemberUserService, deliveryFreightTemplateSvc *trade.DeliveryExpressTemplateService) *TradeOrderHandler {
	return &TradeOrderHandler{
		svc:                        svc,
		querySvc:                   querySvc,
		memberSvc:                  memberSvc,
		deliveryFreightTemplateSvc: deliveryFreightTemplateSvc,
	}
}

// GetOrderPage 获得订单分页
func (h *TradeOrderHandler) GetOrderPage(c *gin.Context) {
	var r req.TradeOrderPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	// Call Service
	pageResult, err := h.querySvc.GetOrderPageForAdmin(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// Fetch Items
	var resultList []resp.TradeOrderPageItemResp
	if len(pageResult.List) > 0 {
		orderIds := make([]int64, len(pageResult.List))
		userIds := make([]int64, 0, len(pageResult.List))
		brokerageUserIds := make([]int64, 0, len(pageResult.List))

		for i, o := range pageResult.List {
			orderIds[i] = o.ID
			userIds = append(userIds, o.UserID)
			if o.BrokerageUserID != nil {
				brokerageUserIds = append(brokerageUserIds, *o.BrokerageUserID)
			}
		}

		// Query Items
		items, err := h.querySvc.GetOrderItemListByOrderIds(c, orderIds)
		if err != nil {
			core.WriteError(c, 500, err.Error())
			return
		}
		itemMap := make(map[int64][]resp.TradeOrderItemBase)
		for _, item := range items {
			itemResp := resp.TradeOrderItemBase{
				ID:       item.ID,
				UserID:   item.UserID,
				OrderID:  item.OrderID,
				SpuID:    item.SpuID,
				SpuName:  item.SpuName,
				SkuID:    item.SkuID,
				PicURL:   item.PicURL,
				Count:    item.Count,
				Price:    item.Price,
				PayPrice: item.PayPrice,
			}
			itemMap[item.OrderID] = append(itemMap[item.OrderID], itemResp)
		}

		// Query Users
		userMap, err := h.memberSvc.GetUserRespMap(c, userIds)
		if err != nil {
			// Log error but continue? Or fail? Java fails if user query fails usually.
			core.WriteError(c, 500, err.Error())
			return
		}

		// Query Brokerage Users
		var brokerageUserMap map[int64]*resp.MemberUserResp
		if len(brokerageUserIds) > 0 {
			brokerageUserMap, err = h.memberSvc.GetUserRespMap(c, brokerageUserIds)
			if err != nil {
				core.WriteError(c, 500, err.Error())
				return
			}
		}

		resultList = make([]resp.TradeOrderPageItemResp, len(pageResult.List))
		for i, o := range pageResult.List {
			var payOrderID int64
			if o.PayOrderID != nil {
				payOrderID = *o.PayOrderID
			}

			var brokerageUser *resp.MemberUserResp
			if o.BrokerageUserID != nil {
				brokerageUser = brokerageUserMap[*o.BrokerageUserID]
			}

			resultList[i] = resp.TradeOrderPageItemResp{
				TradeOrderBase: resp.TradeOrderBase{
					ID:                    o.ID,
					No:                    o.No,
					CreateTime:            o.CreatedAt,
					Type:                  o.Type,
					Terminal:              o.Terminal,
					UserID:                o.UserID,
					UserIP:                o.UserIP,
					UserRemark:            o.UserRemark,
					Status:                o.Status,
					ProductCount:          o.ProductCount,
					FinishTime:            o.FinishTime,
					CancelTime:            o.CancelTime,
					CancelType:            o.CancelType,
					Remark:                o.Remark,
					PayOrderID:            payOrderID,
					PayStatus:             o.PayStatus,
					PayTime:               o.PayTime,
					PayChannelCode:        o.PayChannelCode,
					TotalPrice:            o.TotalPrice,
					DiscountPrice:         o.DiscountPrice,
					DeliveryPrice:         o.DeliveryPrice,
					AdjustPrice:           o.AdjustPrice,
					PayPrice:              o.PayPrice,
					DeliveryType:          o.DeliveryType,
					LogisticsID:           o.LogisticsID,
					LogisticsNo:           o.LogisticsNo,
					DeliveryTime:          o.DeliveryTime,
					ReceiveTime:           o.ReceiveTime,
					ReceiverName:          o.ReceiverName,
					ReceiverMobile:        o.ReceiverMobile,
					ReceiverAreaID:        o.ReceiverAreaID,
					ReceiverDetailAddress: o.ReceiverDetailAddress,
					RefundPrice:           o.RefundPrice,
					CouponID:              o.CouponID,
					CouponPrice:           o.CouponPrice,
				},
				Items:            itemMap[o.ID],
				User:             userMap[o.UserID],
				BrokerageUser:    brokerageUser,
				ReceiverAreaName: "", // TODO: Implement Area Service to get name
			}
		}
	} else {
		resultList = []resp.TradeOrderPageItemResp{}
	}

	core.WriteSuccess(c, core.PageResult[resp.TradeOrderPageItemResp]{
		List:  resultList,
		Total: pageResult.Total,
	})
}

// GetOrderDetail 获得订单详情
func (h *TradeOrderHandler) GetOrderDetail(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "id is required")
		return
	}
	// 1. Get Order
	order, err := h.querySvc.GetOrder(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if order == nil {
		core.WriteError(c, 404, "订单不存在")
		return
	}

	// 2. Get Items
	items, err := h.querySvc.GetOrderItemListByOrderId(c, order.ID)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// 3. Get Logs
	logs, err := h.querySvc.GetOrderLogListByOrderId(c, order.ID)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// 4. Assemble DTO
	itemResps := make([]resp.TradeOrderItemBase, len(items))
	for i, item := range items {
		itemResps[i] = resp.TradeOrderItemBase{
			ID:       item.ID,
			UserID:   item.UserID,
			OrderID:  item.OrderID,
			SpuID:    item.SpuID,
			SpuName:  item.SpuName,
			SkuID:    item.SkuID,
			PicURL:   item.PicURL,
			Count:    item.Count,
			Price:    item.Price,
			PayPrice: item.PayPrice,
			// ... other fields
		}
	}

	logResps := make([]resp.TradeOrderLogResp, len(logs))
	for i, l := range logs {
		logResps[i] = resp.TradeOrderLogResp{
			Content:    l.Content,
			CreateTime: l.CreatedAt,
			UserType:   l.UserType,
		}
	}

	var payOrderID int64
	if order.PayOrderID != nil {
		payOrderID = *order.PayOrderID
	}

	// Prepare User Info if needed, skip for now

	res := resp.TradeOrderDetailResp{
		TradeOrderBase: resp.TradeOrderBase{
			ID:                    order.ID,
			No:                    order.No,
			Type:                  order.Type,
			Terminal:              order.Terminal,
			UserID:                order.UserID,
			UserIP:                order.UserIP,
			UserRemark:            order.UserRemark,
			Status:                order.Status,
			ProductCount:          order.ProductCount,
			FinishTime:            order.FinishTime,
			CancelTime:            order.CancelTime,
			CancelType:            order.CancelType,
			Remark:                order.Remark,
			PayOrderID:            payOrderID,
			PayStatus:             order.PayStatus,
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
			RefundPrice:           order.RefundPrice,
			CouponID:              order.CouponID,
			CouponPrice:           order.CouponPrice,
		},
		Items: itemResps,
		Logs:  logResps,
		// User: Fetch User if possible, or leave nil for now.
	}

	core.WriteSuccess(c, res)
}

// GetOrderExpressTrackList 获得交易订单的物流轨迹
func (h *TradeOrderHandler) GetOrderExpressTrackList(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "id is required")
		return
	}
	tracks, err := h.querySvc.GetExpressTrackListById(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, tracks)
}

// DeliveryOrder 订单发货
func (h *TradeOrderHandler) DeliveryOrder(c *gin.Context) {
	var r req.TradeOrderDeliveryReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.DeliveryOrder(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// UpdateOrderRemark 订单备注
func (h *TradeOrderHandler) UpdateOrderRemark(c *gin.Context) {
	var r req.TradeOrderRemarkReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateOrderRemark(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// UpdateOrderPrice 订单调价
func (h *TradeOrderHandler) UpdateOrderPrice(c *gin.Context) {
	var r req.TradeOrderUpdatePriceReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateOrderPrice(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// UpdateOrderAddress 修改订单收货地址
func (h *TradeOrderHandler) UpdateOrderAddress(c *gin.Context) {
	var r req.TradeOrderUpdateAddressReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateOrderAddress(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// PickUpOrderById 订单核销 (By ID)
func (h *TradeOrderHandler) PickUpOrderById(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "id is required")
		return
	}
	if err := h.svc.PickUpOrderByAdmin(c, core.GetUserId(c), id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// PickUpOrderByVerifyCode 订单核销 (By Code)
func (h *TradeOrderHandler) PickUpOrderByVerifyCode(c *gin.Context) {
	code := c.Query("pickUpVerifyCode")
	if code == "" {
		core.WriteError(c, 400, "pickUpVerifyCode is required")
		return
	}
	if err := h.svc.PickUpOrderByVerifyCode(c, core.GetUserId(c), code); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetByPickUpVerifyCode 查询核销码对应的订单
func (h *TradeOrderHandler) GetByPickUpVerifyCode(c *gin.Context) {
	code := c.Query("pickUpVerifyCode")
	if code == "" {
		core.WriteError(c, 400, "pickUpVerifyCode is required")
		return
	}
	res, err := h.svc.GetByPickUpVerifyCode(c, code)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}

// GetOrderSummary 获得交易订单统计
func (h *TradeOrderHandler) GetOrderSummary(c *gin.Context) {
	var r req.TradeOrderPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.querySvc.GetOrderSummary(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}
