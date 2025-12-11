package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/trade"

	"github.com/gin-gonic/gin"
)

type AppTradeAfterSaleHandler struct {
	svc *trade.TradeAfterSaleService
}

func NewAppTradeAfterSaleHandler(svc *trade.TradeAfterSaleService) *AppTradeAfterSaleHandler {
	return &AppTradeAfterSaleHandler{svc: svc}
}

// CreateAfterSale 申请售后
func (h *AppTradeAfterSaleHandler) CreateAfterSale(c *gin.Context) {
	var r req.AppAfterSaleCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateAfterSale(c, core.GetUserId(c), &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, id)
}

// GetAfterSalePage 获得售后分页
func (h *AppTradeAfterSaleHandler) GetAfterSalePage(c *gin.Context) {
	var r req.TradeAfterSalePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	uid := core.GetUserId(c)
	r.UserID = &uid

	res, err := h.svc.GetAfterSalePage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}

// GetAfterSale 获得售后详情
func (h *AppTradeAfterSaleHandler) GetAfterSale(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	res, err := h.svc.GetAfterSale(c, core.GetUserId(c), id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, res)
}

// CancelAfterSale 取消售后
func (h *AppTradeAfterSaleHandler) CancelAfterSale(c *gin.Context) {
	var r req.AppAfterSaleCancelReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.CancelAfterSale(c, core.GetUserId(c), r.ID); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// DeliveryAfterSale 退回货物 (Submit Logistics)
func (h *AppTradeAfterSaleHandler) DeliveryAfterSale(c *gin.Context) {
	var r req.AppAfterSaleDeliveryReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.DeliveryAfterSale(c, core.GetUserId(c), &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}
