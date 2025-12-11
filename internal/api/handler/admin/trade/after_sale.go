package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/trade"
	"strconv"

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
	var r req.TradeAfterSalePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.svc.GetAfterSalePage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}

// AgreeAfterSale 同意售后
func (h *TradeAfterSaleHandler) AgreeAfterSale(c *gin.Context) {
	var r req.TradeAfterSaleAgreeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.AgreeAfterSale(c, r.ID); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// DisagreeAfterSale 拒绝售后
func (h *TradeAfterSaleHandler) DisagreeAfterSale(c *gin.Context) {
	var r req.TradeAfterSaleDisagreeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.DisagreeAfterSale(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// RefundAfterSale 退款
func (h *TradeAfterSaleHandler) RefundAfterSale(c *gin.Context) {
	var r req.TradeAfterSaleRefundReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.RefundAfterSale(c, r.ID); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetAfterSaleDetail 获得售后订单详情
func (h *TradeAfterSaleHandler) GetAfterSaleDetail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		core.WriteError(c, 400, "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "invalid id")
		return
	}
	res, err := h.svc.GetAfterSaleDetail(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}

// ReceiveAfterSale 确认收货
func (h *TradeAfterSaleHandler) ReceiveAfterSale(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		core.WriteError(c, 400, "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		core.WriteError(c, 400, "invalid id")
		return
	}
	if err := h.svc.ReceiveAfterSale(c, id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// UpdateAfterSaleRefunded 更新售后单为已退款 (Callback)
// @Router /admin-api/trade/after-sale/update-refunded [post]
func (h *TradeAfterSaleHandler) UpdateAfterSaleRefunded(c *gin.Context) {
	var r req.PayRefundNotifyReqDTO // Reuse Pay DTO
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateRefunded(c.Request.Context(), &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}
