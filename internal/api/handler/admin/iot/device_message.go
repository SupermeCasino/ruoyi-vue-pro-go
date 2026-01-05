package iot

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type DeviceMessageHandler struct {
	svc *iotsvc.DeviceMessageService
}

func NewDeviceMessageHandler(svc *iotsvc.DeviceMessageService) *DeviceMessageHandler {
	return &DeviceMessageHandler{svc: svc}
}

func (h *DeviceMessageHandler) GetPage(c *gin.Context) {
	var req iot2.IotDeviceMessagePageReqVO
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	pageResult, err := h.svc.GetDeviceMessagePage(c.Request.Context(), &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotDeviceMessageRespVO, 0, len(pageResult.List))
	for _, do := range pageResult.List {
		list = append(list, h.toVO(do))
	}

	response.WriteSuccess(c, &pagination.PageResult[*iot2.IotDeviceMessageRespVO]{
		List:  list,
		Total: pageResult.Total,
	})
}

func (h *DeviceMessageHandler) GetPairPage(c *gin.Context) {
	var req iot2.IotDeviceMessagePageReqVO
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 1.1 先按照条件，查询 request 的消息（非 reply）
	reqReply := false
	req.Reply = &reqReply
	requestMessagePageResult, err := h.svc.GetDeviceMessagePage(c.Request.Context(), &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	if len(requestMessagePageResult.List) == 0 {
		response.WriteSuccess(c, pagination.NewEmptyPageResult[*iot2.IotDeviceMessageRespPairVO]())
		return
	}

	// 1.2 接着按照 requestIds，批量查询 reply 消息
	requestIDs := make([]string, 0, len(requestMessagePageResult.List))
	for _, m := range requestMessagePageResult.List {
		if m.RequestID != "" {
			requestIDs = append(requestIDs, m.RequestID)
		}
	}

	var replyMessageList []*model.IotDeviceMessageDO
	if len(requestIDs) > 0 {
		replyMessageList, err = h.svc.GetDeviceMessageListByRequestIdsAndReply(c.Request.Context(), req.DeviceID, requestIDs, true)
		if err != nil {
			response.WriteBizError(c, err)
			return
		}
	}

	replyMessagesMap := make(map[string]*model.IotDeviceMessageDO)
	for _, m := range replyMessageList {
		replyMessagesMap[m.RequestID] = m
	}

	// 2. 组装结果
	pairMessages := make([]*iot2.IotDeviceMessageRespPairVO, 0, len(requestMessagePageResult.List))
	for _, requestMessage := range requestMessagePageResult.List {
		pair := &iot2.IotDeviceMessageRespPairVO{
			Request: h.toVO(requestMessage),
		}
		if replyMessage, ok := replyMessagesMap[requestMessage.RequestID]; ok {
			pair.Reply = h.toVO(replyMessage)
		}
		pairMessages = append(pairMessages, pair)
	}

	response.WriteSuccess(c, &pagination.PageResult[*iot2.IotDeviceMessageRespPairVO]{
		List:  pairMessages,
		Total: requestMessagePageResult.Total,
	})
}

func (h *DeviceMessageHandler) Send(c *gin.Context) {
	var req iot2.IotDeviceMessageSendReqVO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	if err := h.svc.SendDeviceMessage(c.Request.Context(), &req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DeviceMessageHandler) toVO(do *model.IotDeviceMessageDO) *iot2.IotDeviceMessageRespVO {
	if do == nil {
		return nil
	}
	reportTime := time.UnixMilli(do.ReportTime)
	ts := time.UnixMilli(do.TS)

	vo := &iot2.IotDeviceMessageRespVO{
		ID:         do.ID,
		ReportTime: &reportTime,
		TS:         &ts,
		DeviceID:   do.DeviceID,
		ServerID:   do.ServerID,
		Upstream:   do.Upstream,
		Reply:      do.Reply,
		Identifier: do.Identifier,
		RequestID:  do.RequestID,
		Method:     do.Method,
		Code:       do.Code,
		Msg:        do.Msg,
	}

	if do.Params != "" {
		var params interface{}
		if err := json.Unmarshal([]byte(do.Params), &params); err == nil {
			vo.Params = params
		} else {
			vo.Params = do.Params
		}
	}
	if do.Data != "" {
		var data interface{}
		if err := json.Unmarshal([]byte(do.Data), &data); err == nil {
			vo.Data = data
		} else {
			vo.Data = do.Data
		}
	}

	return vo
}
