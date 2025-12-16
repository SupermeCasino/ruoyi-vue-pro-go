package promotion

import (
	"github.com/gin-gonic/gin"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
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

// SendMessage 发送消息
func (h *AppKefuHandler) SendMessage(c *gin.Context) {
	var r req.KefuMessageCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 获取当前登录用户ID
	userID := c.GetInt64("userId")                  // 假设中间件注入了 userId
	id, err := h.svc.CreateMessage(c, r, userID, 1) // SenderType 1 = User
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}
