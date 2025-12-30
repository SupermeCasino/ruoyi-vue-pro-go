package promotion

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppRewardActivityHandler struct {
	svc *promotion.RewardActivityService
}

func NewAppRewardActivityHandler(svc *promotion.RewardActivityService) *AppRewardActivityHandler {
	return &AppRewardActivityHandler{svc: svc}
}

// GetRewardActivity 获得满减送活动
// @Summary 获得满减送活动
// @Router /app-api/promotion/reward-activity/get [get]
// Java: AppRewardActivityController#getRewardActivity
func (h *AppRewardActivityHandler) GetRewardActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	res, err := h.svc.GetRewardActivityForApp(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, res)
}
