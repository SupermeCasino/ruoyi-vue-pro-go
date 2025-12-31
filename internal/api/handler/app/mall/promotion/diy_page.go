package promotion

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetDiyPageModel(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, &resp.AppDiyPagePropertyResp{
		ID:       page.ID,
		Name:     page.Name,
		Property: page.Property,
	})
}
