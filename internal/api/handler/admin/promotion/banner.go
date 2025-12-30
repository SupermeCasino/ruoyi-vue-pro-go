package promotion

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
	svc *promotion.PromotionBannerService
}

func NewBannerHandler(svc *promotion.PromotionBannerService) *BannerHandler {
	return &BannerHandler{svc: svc}
}

// CreateBanner 创建 Banner
// @Summary 创建 Banner
// @Tags 管理后台 - 营销 Banner
// @Produce json
// @Param req body req.PromotionBannerCreateReq true "请求参数"
// @Success 200 {object} core.Response
// @Router /admin-api/promotion/banner/create [post]
func (h *BannerHandler) CreateBanner(c *gin.Context) {
	var r req.PromotionBannerCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateBanner(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateBanner 更新 Banner
// @Summary 更新 Banner
// @Tags 管理后台 - 营销 Banner
// @Produce json
// @Param req body req.PromotionBannerUpdateReq true "请求参数"
// @Success 200 {object} core.Response
// @Router /admin-api/promotion/banner/update [put]
func (h *BannerHandler) UpdateBanner(c *gin.Context) {
	var r req.PromotionBannerUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateBanner(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteBanner 删除 Banner
// @Summary 删除 Banner
// @Tags 管理后台 - 营销 Banner
// @Produce json
// @Param id query int true "编号"
// @Success 200 {object} core.Response
// @Router /admin-api/promotion/banner/delete [delete]
func (h *BannerHandler) DeleteBanner(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteBanner(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetBanner 获得 Banner
// @Summary 获得 Banner
// @Tags 管理后台 - 营销 Banner
// @Produce json
// @Param id query int true "编号"
// @Router /admin-api/promotion/banner/get [get]
func (h *BannerHandler) GetBanner(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetBanner(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetBannerPage 获得 Banner 分页
// @Summary 获得 Banner 分页
// @Tags 管理后台 - 营销 Banner
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "页数"
// @Param title query string false "标题"
// @Param status query int false "状态"
// @Router /admin-api/promotion/banner/page [get]
func (h *BannerHandler) GetBannerPage(c *gin.Context) {
	var r req.PromotionBannerPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetBannerPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
