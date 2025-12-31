package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// AppSeckillConfigHandler App 端秒杀时段 Handler
type AppSeckillConfigHandler struct {
	svc *promotion.SeckillConfigService
}

// NewAppSeckillConfigHandler 创建 Handler
func NewAppSeckillConfigHandler(svc *promotion.SeckillConfigService) *AppSeckillConfigHandler {
	return &AppSeckillConfigHandler{svc: svc}
}

// GetSeckillConfigList 获得启用的秒杀时段列表
// 对齐 Java: AppSeckillConfigController.getSeckillConfigList
func (h *AppSeckillConfigHandler) GetSeckillConfigList(c *gin.Context) {
	list, err := h.svc.GetSeckillConfigListByStatus(c.Request.Context(), consts.CommonStatusEnable)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 构建响应
	result := make([]resp.AppSeckillConfigResp, len(list))
	for i, cfg := range list {
		result[i] = resp.AppSeckillConfigResp{
			ID:            cfg.ID,
			StartTime:     cfg.StartTime,
			EndTime:       cfg.EndTime,
			SliderPicUrls: cfg.SliderPicUrls,
		}
	}

	response.WriteSuccess(c, result)
}
