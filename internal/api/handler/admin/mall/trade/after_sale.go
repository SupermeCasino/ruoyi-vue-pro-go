package trade

import (
	"strconv"

	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type TradeAfterSaleHandler struct {
	svc *trade.TradeAfterSaleService
}

func NewTradeAfterSaleHandler(svc *trade.TradeAfterSaleService) *TradeAfterSaleHandler {
	return &TradeAfterSaleHandler{svc: svc}
}

// GetAfterSalePage 获得售后分页
func (h *TradeAfterSaleHandler) GetAfterSalePage(c *gin.Context) {
	var r trade2.TradeAfterSalePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.svc.GetAfterSalePage(c, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, res)
}

// AgreeAfterSale 同意售后
func (h *TradeAfterSaleHandler) AgreeAfterSale(c *gin.Context) {
	var r trade2.TradeAfterSaleAgreeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.AgreeAfterSale(c, context.GetUserId(c), r.ID); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// DisagreeAfterSale 拒绝售后
func (h *TradeAfterSaleHandler) DisagreeAfterSale(c *gin.Context) {
	var r trade2.TradeAfterSaleDisagreeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.DisagreeAfterSale(c, context.GetUserId(c), &r); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// RefundAfterSale 退款
func (h *TradeAfterSaleHandler) RefundAfterSale(c *gin.Context) {
	var r trade2.TradeAfterSaleRefundReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.RefundAfterSale(c, context.GetUserId(c), c.ClientIP(), r.ID); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// GetAfterSaleDetail 获得售后订单详情
func (h *TradeAfterSaleHandler) GetAfterSaleDetail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.WriteError(c, 400, "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(c, 400, "invalid id")
		return
	}
	res, err := h.svc.GetAfterSaleDetail(c, id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, res)
}

// ReceiveAfterSale 确认收货
func (h *TradeAfterSaleHandler) ReceiveAfterSale(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.WriteError(c, 400, "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(c, 400, "invalid id")
		return
	}
	if err := h.svc.ReceiveAfterSale(c, context.GetUserId(c), id); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// RefuseAfterSale 拒绝收货
func (h *TradeAfterSaleHandler) RefuseAfterSale(c *gin.Context) {
	var r trade2.TradeAfterSaleRefuseReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.RefuseAfterSale(c, context.GetUserId(c), &r); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateAfterSaleRefunded 更新售后单为已退款 (Callback)
// @Router /admin-api/trade/after-sale/update-refunded [post]
func (h *TradeAfterSaleHandler) UpdateAfterSaleRefunded(c *gin.Context) {
	var r pay.PayRefundNotifyReqDTO // Reuse Pay DTO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateRefunded(c.Request.Context(), &r); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}
