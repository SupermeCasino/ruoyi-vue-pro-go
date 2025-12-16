package trade

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"

	"github.com/gin-gonic/gin"
)

type AppTradeConfigHandler struct {
	svc *trade.TradeConfigService
}

func NewAppTradeConfigHandler(svc *trade.TradeConfigService) *AppTradeConfigHandler {
	return &AppTradeConfigHandler{svc: svc}
}

// GetTradeConfig @Summary 获得交易配置
// @Router /app-api/trade/config/get [GET]
func (h *AppTradeConfigHandler) GetTradeConfig(c *gin.Context) {
	res, err := h.svc.GetTradeConfig(c)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}
