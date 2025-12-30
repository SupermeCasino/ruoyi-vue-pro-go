package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppBannerHandler struct {
	svc *promotion.PromotionBannerService
}

func NewAppBannerHandler(svc *promotion.PromotionBannerService) *AppBannerHandler {
	return &AppBannerHandler{svc: svc}
}

// GetBannerList 获得首页 Banner 列表
// @Summary 获得首页 Banner 列表
// @Tags 用户 APP - 营销 Banner
// @Produce json
// @Param position query int true "位置"
// @Router /app-api/promotion/banner/list [get]
func (h *AppBannerHandler) GetBannerList(c *gin.Context) {
	var r req.AppBannerListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetAppBannerList(c, r.Position)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
