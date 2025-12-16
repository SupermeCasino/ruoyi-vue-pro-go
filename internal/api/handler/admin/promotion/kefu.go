package promotion

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type KefuHandler struct {
	svc promotion.KefuService
}

func NewKefuHandler(svc promotion.KefuService) *KefuHandler {
	return &KefuHandler{svc: svc}
}

// GetConversationPage 获得客服会话分页
func (h *KefuHandler) GetConversationPage(c *gin.Context) {
	var r req.KefuConversationPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetConversationPage(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
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

// GetMessagePage 获得客服消息分页
func (h *KefuHandler) GetMessagePage(c *gin.Context) {
	var r req.KefuMessagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetMessagePage(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// SendMessage 发送客服消息
func (h *KefuHandler) SendMessage(c *gin.Context) {
	var r req.KefuMessageCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 获取当前登录管理员ID
	adminID := c.GetInt64("uid")                     // 假设中间件注入了 admin uid
	id, err := h.svc.CreateMessage(c, r, adminID, 2) // SenderType 2 = Admin
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}
