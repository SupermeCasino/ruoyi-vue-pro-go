package iot

import (
	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type StatisticsHandler struct {
	statisticsService *iotsvc.StatisticsService
}

func NewStatisticsHandler(statisticsService *iotsvc.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHandler) GetSummary(c *gin.Context) {
	summary, err := h.statisticsService.GetStatisticsSummary(c.Request.Context())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, summary)
}

func (h *StatisticsHandler) GetDeviceMessageSummaryByDate(c *gin.Context) {
	var req iot2.IotStatisticsDeviceMessageReqVO
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	list, err := h.statisticsService.GetDeviceMessageSummaryByDate(c.Request.Context(), &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}
