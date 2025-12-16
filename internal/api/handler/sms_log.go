package handler

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type SmsLogHandler struct {
	smsLogSvc *service.SmsLogService
}

func NewSmsLogHandler(smsLogSvc *service.SmsLogService) *SmsLogHandler {
	return &SmsLogHandler{
		smsLogSvc: smsLogSvc,
	}
}

// GetSmsLogPage 获得短信日志分页
func (h *SmsLogHandler) GetSmsLogPage(c *gin.Context) {
	var req req.SmsLogPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	res, err := h.smsLogSvc.GetSmsLogPage(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}
