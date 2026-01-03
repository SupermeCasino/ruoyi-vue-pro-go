package pay

import (
	"strconv"

	reqPay "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
	servicePay "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type PayTransferHandler struct {
	transferSvc *servicePay.PayTransferService
}

func NewPayTransferHandler(transferSvc *servicePay.PayTransferService) *PayTransferHandler {
	return &PayTransferHandler{
		transferSvc: transferSvc,
	}
}

// GetTransfer 获得转账订单
func (h *PayTransferHandler) GetTransfer(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	transfer, err := h.transferSvc.GetTransfer(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, transfer)
}

// GetTransferPage 获得转账订单分页
func (h *PayTransferHandler) GetTransferPage(c *gin.Context) {
	var req reqPay.PayTransferPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	pageResult, err := h.transferSvc.GetTransferPage(c.Request.Context(), &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, pageResult)
}
