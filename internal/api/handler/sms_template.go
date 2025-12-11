package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SmsTemplateHandler struct {
	smsTemplateSvc *service.SmsTemplateService
}

func NewSmsTemplateHandler(smsTemplateSvc *service.SmsTemplateService) *SmsTemplateHandler {
	return &SmsTemplateHandler{
		smsTemplateSvc: smsTemplateSvc,
	}
}

// CreateSmsTemplate 创建短信模板
func (h *SmsTemplateHandler) CreateSmsTemplate(c *gin.Context) {
	var req req.SmsTemplateSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	id, err := h.smsTemplateSvc.CreateSmsTemplate(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateSmsTemplate 更新短信模板
func (h *SmsTemplateHandler) UpdateSmsTemplate(c *gin.Context) {
	var req req.SmsTemplateSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	if err := h.smsTemplateSvc.UpdateSmsTemplate(c, &req); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteSmsTemplate 删除短信模板
func (h *SmsTemplateHandler) DeleteSmsTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	if err := h.smsTemplateSvc.DeleteSmsTemplate(c, id); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

// GetSmsTemplate 获得短信模板
func (h *SmsTemplateHandler) GetSmsTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	res, err := h.smsTemplateSvc.GetSmsTemplate(c, id)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}

// GetSmsTemplatePage 获得短信模板分页
func (h *SmsTemplateHandler) GetSmsTemplatePage(c *gin.Context) {
	var req req.SmsTemplatePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	res, err := h.smsTemplateSvc.GetSmsTemplatePage(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}
