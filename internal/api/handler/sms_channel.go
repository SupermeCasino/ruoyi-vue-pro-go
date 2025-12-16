package handler

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type SmsChannelHandler struct {
	smsChannelSvc *service.SmsChannelService
}

func NewSmsChannelHandler(smsChannelSvc *service.SmsChannelService) *SmsChannelHandler {
	return &SmsChannelHandler{
		smsChannelSvc: smsChannelSvc,
	}
}

// CreateSmsChannel 创建短信渠道
func (h *SmsChannelHandler) CreateSmsChannel(c *gin.Context) {
	var req req.SmsChannelSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	id, err := h.smsChannelSvc.CreateSmsChannel(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(id))
}

// UpdateSmsChannel 更新短信渠道
func (h *SmsChannelHandler) UpdateSmsChannel(c *gin.Context) {
	var req req.SmsChannelSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	if err := h.smsChannelSvc.UpdateSmsChannel(c, &req); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

// DeleteSmsChannel 删除短信渠道
func (h *SmsChannelHandler) DeleteSmsChannel(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	if err := h.smsChannelSvc.DeleteSmsChannel(c, id); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

// GetSmsChannel 获得短信渠道
func (h *SmsChannelHandler) GetSmsChannel(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	res, err := h.smsChannelSvc.GetSmsChannel(c, id)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

// GetSmsChannelPage 获得短信渠道分页
func (h *SmsChannelHandler) GetSmsChannelPage(c *gin.Context) {
	var req req.SmsChannelPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	res, err := h.smsChannelSvc.GetSmsChannelPage(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

// GetSimpleSmsChannelList 获得短信渠道精简列表
func (h *SmsChannelHandler) GetSimpleSmsChannelList(c *gin.Context) {
	res, err := h.smsChannelSvc.GetSimpleSmsChannelList(c)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}
