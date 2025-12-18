package pay

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	payModel "github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type PayChannelHandler struct {
	svc *paySvc.PayChannelService
}

func NewPayChannelHandler(svc *paySvc.PayChannelService) *PayChannelHandler {
	return &PayChannelHandler{svc: svc}
}

// CreateChannel 创建支付渠道
func (h *PayChannelHandler) CreateChannel(c *gin.Context) {
	var r req.PayChannelCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateChannel(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(id))
}

// UpdateChannel 更新支付渠道
func (h *PayChannelHandler) UpdateChannel(c *gin.Context) {
	var r req.PayChannelUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	err := h.svc.UpdateChannel(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// DeleteChannel 删除支付渠道
func (h *PayChannelHandler) DeleteChannel(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	err = h.svc.DeleteChannel(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// GetChannel 获得支付渠道
func (h *PayChannelHandler) GetChannel(c *gin.Context) {
	idStr := c.Query("id")
	appIdStr := c.Query("appId")
	code := c.Query("code")

	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(200, errors.ErrParam)
			return
		}
		channel, err := h.svc.GetChannel(c, id)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(200, response.Success(convertChannelResp(channel)))
		return
	}

	if appIdStr != "" && code != "" {
		appId, err := strconv.ParseInt(appIdStr, 10, 64)
		if err != nil {
			c.JSON(200, errors.ErrParam)
			return
		}
		channel, err := h.svc.GetChannelByAppIdAndCode(c, appId, code)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(200, response.Success(convertChannelResp(channel)))
		return
	}

	c.JSON(200, errors.ErrParam)
}

// GetEnableChannelCodeList 获得指定应用的开启的支付渠道编码列表
func (h *PayChannelHandler) GetEnableChannelCodeList(c *gin.Context) {
	appIdStr := c.Query("appId")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	channels, err := h.svc.GetEnableChannelList(c, appId)
	if err != nil {
		c.Error(err)
		return
	}

	codes := make([]string, 0, len(channels))
	for _, ch := range channels {
		codes = append(codes, ch.Code)
	}
	c.JSON(200, response.Success(codes))
}

func convertChannelResp(channel *payModel.PayChannel) *resp.PayChannelResp {
	if channel == nil {
		return nil
	}
	return &resp.PayChannelResp{
		ID:         channel.ID,
		Code:       channel.Code,
		Status:     channel.Status,
		FeeRate:    channel.FeeRate,
		Remark:     channel.Remark,
		AppID:      channel.AppID,
		Config:     channel.Config.ToJSON(),
		CreateTime: channel.CreatedAt,
	}
}
