package pay

import (
	"github.com/gin-gonic/gin"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

type AppPayTransferHandler struct {
	svc *paySvc.PayTransferService
}

func NewAppPayTransferHandler(svc *paySvc.PayTransferService) *AppPayTransferHandler {
	return &AppPayTransferHandler{svc: svc}
}

// SyncTransfer 同步转账单
func (h *AppPayTransferHandler) SyncTransfer(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 400, "参数错误")
		return
	}

	err := h.svc.SyncTransferById(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, true)
}
