package system

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type SmsChannelHandler struct {
	smsChannelSvc *system.SmsChannelService
}

func NewSmsChannelHandler(smsChannelSvc *system.SmsChannelService) *SmsChannelHandler {
	return &SmsChannelHandler{
		smsChannelSvc: smsChannelSvc,
	}
}

// CreateSmsChannel 创建短信渠道
func (h *SmsChannelHandler) CreateSmsChannel(c *gin.Context) {
	var req req.SmsChannelSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.smsChannelSvc.CreateSmsChannel(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateSmsChannel 更新短信渠道
func (h *SmsChannelHandler) UpdateSmsChannel(c *gin.Context) {
	var req req.SmsChannelSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.smsChannelSvc.UpdateSmsChannel(c, &req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteSmsChannel 删除短信渠道
func (h *SmsChannelHandler) DeleteSmsChannel(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.smsChannelSvc.DeleteSmsChannel(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetSmsChannel 获得短信渠道
func (h *SmsChannelHandler) GetSmsChannel(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.smsChannelSvc.GetSmsChannel(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetSmsChannelPage 获得短信渠道分页
func (h *SmsChannelHandler) GetSmsChannelPage(c *gin.Context) {
	var req req.SmsChannelPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.smsChannelSvc.GetSmsChannelPage(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetSimpleSmsChannelList 获得短信渠道精简列表
func (h *SmsChannelHandler) GetSimpleSmsChannelList(c *gin.Context) {
	res, err := h.smsChannelSvc.GetSimpleSmsChannelList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
