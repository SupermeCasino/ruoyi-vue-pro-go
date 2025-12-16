package trade

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type TradeConfigHandler struct {
	svc *trade.TradeConfigService
}

func NewTradeConfigHandler(svc *trade.TradeConfigService) *TradeConfigHandler {
	return &TradeConfigHandler{svc: svc}
}

// GetTradeConfig @Summary 获得交易配置
// @Router /admin-api/trade/config/get [GET]
func (h *TradeConfigHandler) GetTradeConfig(c *gin.Context) {
	res, err := h.svc.GetTradeConfig(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// SaveTradeConfig @Summary 保存交易配置
// @Router /admin-api/trade/config/save [PUT]
func (h *TradeConfigHandler) SaveTradeConfig(c *gin.Context) {
	var r req.TradeConfigSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	if err := h.svc.SaveTradeConfig(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}
