package trade

import (
	"os"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
// 对齐 Java: AppTradeConfigController.getTradeConfig
func (h *AppTradeConfigHandler) GetTradeConfig(c *gin.Context) {
	res, err := h.svc.GetAppTradeConfig(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	// 设置腾讯地图 Key（对齐 Java: @Value("${yudao.tencent-lbs-key}")）
	res.TencentLbsKey = os.Getenv("YUDAO_TENCENT_LBS_KEY")
	response.WriteSuccess(c, res)
}
