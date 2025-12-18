package trade

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

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
		response.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateAfterSale(c, context.GetUserId(c), &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, id)
}

// GetAfterSalePage 获得售后分页
func (h *AppTradeAfterSaleHandler) GetAfterSalePage(c *gin.Context) {
	var r req.AppAfterSalePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.svc.GetUserAfterSalePage(c, context.GetUserId(c), &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, res)
}

// GetAfterSale 获得售后详情
func (h *AppTradeAfterSaleHandler) GetAfterSale(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	res, err := h.svc.GetAfterSale(c, context.GetUserId(c), id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, res)
}

// CancelAfterSale 取消售后
func (h *AppTradeAfterSaleHandler) CancelAfterSale(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "id is required")
		return
	}
	if err := h.svc.CancelAfterSale(c, context.GetUserId(c), id); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}

// DeliveryAfterSale 退回货物 (Submit Logistics)
func (h *AppTradeAfterSaleHandler) DeliveryAfterSale(c *gin.Context) {
	var r req.AppAfterSaleDeliveryReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.DeliveryAfterSale(c, context.GetUserId(c), &r); err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, true)
}
