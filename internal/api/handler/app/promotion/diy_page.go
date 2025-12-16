package promotion

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type AppDiyPageHandler struct {
	svc promotion.DiyPageService
}

func NewAppDiyPageHandler(svc promotion.DiyPageService) *AppDiyPageHandler {
	return &AppDiyPageHandler{svc: svc}
}

func (h *AppDiyPageHandler) GetDiyPage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.GetDiyPage(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}
