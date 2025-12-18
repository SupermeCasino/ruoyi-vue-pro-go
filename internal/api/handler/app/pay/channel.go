package pay

import (
	"github.com/gin-gonic/gin"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

type AppPayChannelHandler struct {
	svc *paySvc.PayChannelService
}

func NewAppPayChannelHandler(svc *paySvc.PayChannelService) *AppPayChannelHandler {
	return &AppPayChannelHandler{svc: svc}
}

// GetEnableChannelCodeList 获得指定应用的开启的支付渠道编码列表
func (h *AppPayChannelHandler) GetEnableChannelCodeList(c *gin.Context) {
	appId := utils.ParseInt64(c.Query("appId"))
	if appId == 0 {
		response.WriteError(c, 400, "参数错误")
		return
	}

	channels, err := h.svc.GetEnableChannelList(c, appId)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 提取 code 集合
	codes := make([]string, 0, len(channels))
	for _, channel := range channels {
		codes = append(codes, channel.Code)
	}

	response.WriteSuccess(c, codes)
}
