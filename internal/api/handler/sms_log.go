package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

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
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	res, err := h.smsLogSvc.GetSmsLogPage(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}
