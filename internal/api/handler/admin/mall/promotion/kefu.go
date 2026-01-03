package promotion

import (
	"strconv"

	"github.com/gin-gonic/gin"
	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	promotionModel "github.com/wxlbd/ruoyi-mall-go/internal/consts"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type KefuHandler struct {
	svc promotion.KefuService
}

func NewKefuHandler(svc promotion.KefuService) *KefuHandler {
	return &KefuHandler{svc: svc}
}

// GetConversation 获得客服会话
func (h *KefuHandler) GetConversation(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetConversation(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetConversationList 获得客服会话列表 (对齐 Java /list 端点)
func (h *KefuHandler) GetConversationList(c *gin.Context) {
	// Java 使用无分页的列表，Go 端暂时也返回全量列表
	res, err := h.svc.GetConversationList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// UpdateConversationPinned 置顶/取消置顶客服会话
func (h *KefuHandler) UpdateConversationPinned(c *gin.Context) {
	var r promotion2.KeFuConversationUpdatePinnedReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateConversationPinned(c, r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteConversation 删除客服会话
func (h *KefuHandler) DeleteConversation(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteConversation(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetMessageList 获得客服消息列表 (对齐 Java /list 端点)
func (h *KefuHandler) GetMessageList(c *gin.Context) {
	var r promotion2.KefuMessageListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetMessageList(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// UpdateMessageReadStatus 更新客服消息已读状态
func (h *KefuHandler) UpdateMessageReadStatus(c *gin.Context) {
	conversationID, _ := strconv.ParseInt(c.Query("conversationId"), 10, 64)
	if conversationID == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 获取当前登录管理员ID
	adminID := c.GetInt64("uid")
	if err := h.svc.UpdateMessageReadStatus(c, conversationID, adminID, promotionModel.SenderTypeAdmin); err != nil { // 使用客服发送者类型常量替代魔法数字 2
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// SendMessage 发送客服消息
func (h *KefuHandler) SendMessage(c *gin.Context) {
	var r promotion2.KefuMessageCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 获取当前登录管理员ID
	adminID := c.GetInt64("uid")                                                  // 假设中间件注入了 admin uid
	id, err := h.svc.CreateMessage(c, r, adminID, promotionModel.SenderTypeAdmin) // 使用客服发送者类型常量替代魔法数字 2
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}
