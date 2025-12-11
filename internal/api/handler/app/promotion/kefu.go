package promotion

import (
	"github.com/gin-gonic/gin"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/promotion"
)

type AppKefuHandler struct {
	svc promotion.KefuService
}

func NewAppKefuHandler(svc promotion.KefuService) *AppKefuHandler {
	return &AppKefuHandler{svc: svc}
}

// GetMessagePage 获得消息分页
func (h *AppKefuHandler) GetMessagePage(c *gin.Context) {
	var r req.KefuMessagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.GetMessagePage(c, r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// SendMessage 发送消息
func (h *AppKefuHandler) SendMessage(c *gin.Context) {
	var r req.KefuMessageCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	// 获取当前登录用户ID
	userID := c.GetInt64("userId")                  // 假设中间件注入了 userId
	id, err := h.svc.CreateMessage(c, r, userID, 1) // SenderType 1 = User
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}
