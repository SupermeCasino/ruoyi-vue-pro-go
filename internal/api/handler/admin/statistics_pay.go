package admin

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// PayStatisticsHandler 支付统计处理器
type PayStatisticsHandler struct {
	payWalletStatisticsService pay.PayWalletStatisticsService
}

// NewPayStatisticsHandler 创建支付统计处理器
func NewPayStatisticsHandler(payWalletStatisticsService pay.PayWalletStatisticsService) *PayStatisticsHandler {
	return &PayStatisticsHandler{
		payWalletStatisticsService: payWalletStatisticsService,
	}
}

// GetWalletRechargePrice 获取充值金额
// GET /statistics/pay/summary
func (h *PayStatisticsHandler) GetWalletRechargePrice(c *gin.Context) {
	rechargePrice, err := h.payWalletStatisticsService.GetRechargePriceSummary(c)
	if err != nil {
		response.WriteBizError(c, errors.ErrUnknown)
		return
	}

	result := &resp.PaySummaryRespVO{
		RechargePrice: rechargePrice,
	}

	response.WriteSuccess(c, result)
}
