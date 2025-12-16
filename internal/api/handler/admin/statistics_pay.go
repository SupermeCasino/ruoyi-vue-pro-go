package admin

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"

	"github.com/gin-gonic/gin"
)

// PayStatisticsHandler 支付统计处理器
type PayStatisticsHandler struct {
	payWalletStatisticsService service.PayWalletStatisticsService
}

// NewPayStatisticsHandler 创建支付统计处理器
func NewPayStatisticsHandler(payWalletStatisticsService service.PayWalletStatisticsService) *PayStatisticsHandler {
	return &PayStatisticsHandler{
		payWalletStatisticsService: payWalletStatisticsService,
	}
}

// GetWalletRechargePrice 获取充值金额
// GET /statistics/pay/summary
func (h *PayStatisticsHandler) GetWalletRechargePrice(c *gin.Context) {
	rechargePrice, err := h.payWalletStatisticsService.GetRechargePriceSummary(c)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	result := &resp.PaySummaryRespVO{
		RechargePrice: rechargePrice,
	}

	core.WriteSuccess(c, result)
}
