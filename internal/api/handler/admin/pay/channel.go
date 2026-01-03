package pay

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
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
	var r pay.PayChannelCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateChannel(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateChannel 更新支付渠道
func (h *PayChannelHandler) UpdateChannel(c *gin.Context) {
	var r pay.PayChannelUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.UpdateChannel(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteChannel 删除支付渠道
func (h *PayChannelHandler) DeleteChannel(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err = h.svc.DeleteChannel(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetChannel 获得支付渠道
func (h *PayChannelHandler) GetChannel(c *gin.Context) {
	idStr := c.Query("id")
	appIdStr := c.Query("appId")
	code := c.Query("code")

	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteBizError(c, errors.ErrParam)
			return
		}
		channel, err := h.svc.GetChannel(c, id)
		if err != nil {
			response.WriteBizError(c, err)
			return
		}
		// 对齐Java: 查询不到返回null，Go版本返回nil
		response.WriteSuccess(c, convertChannelResp(channel))
		return
	}

	if appIdStr != "" && code != "" {
		appId, err := strconv.ParseInt(appIdStr, 10, 64)
		if err != nil {
			response.WriteBizError(c, errors.ErrParam)
			return
		}
		channel, err := h.svc.GetChannelByAppIdAndCode(c, appId, code)
		if err != nil {
			response.WriteBizError(c, err)
			return
		}
		// 对齐Java: 查询不到返回null，Go版本返回nil
		response.WriteSuccess(c, convertChannelResp(channel))
		return
	}

	response.WriteBizError(c, errors.ErrParam)
}

// GetEnableChannelCodeList 获得指定应用的开启的支付渠道编码列表
func (h *PayChannelHandler) GetEnableChannelCodeList(c *gin.Context) {
	appIdStr := c.Query("appId")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	channels, err := h.svc.GetEnableChannelList(c, appId)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	codes := make([]string, 0, len(channels))
	for _, ch := range channels {
		codes = append(codes, ch.Code)
	}
	response.WriteSuccess(c, codes)
}

func convertChannelResp(channel *payModel.PayChannel) *pay.PayChannelResp {
	if channel == nil {
		return nil
	}
	return &pay.PayChannelResp{
		ID:         channel.ID,
		Code:       channel.Code,
		Status:     channel.Status,
		FeeRate:    channel.FeeRate,
		Remark:     channel.Remark,
		AppID:      channel.AppID,
		Config:     channel.Config.ToJSON(),
		CreateTime: channel.CreateTime,
	}
}
