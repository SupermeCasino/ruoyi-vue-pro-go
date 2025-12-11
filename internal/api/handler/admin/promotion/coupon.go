package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	svc *promotion.CouponService
}

func NewCouponHandler(svc *promotion.CouponService) *CouponHandler {
	return &CouponHandler{svc: svc}
}

// CreateCouponTemplate 创建模板
func (h *CouponHandler) CreateCouponTemplate(c *gin.Context) {
	var r req.CouponTemplateCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateCouponTemplate(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateCouponTemplate 更新模板
func (h *CouponHandler) UpdateCouponTemplate(c *gin.Context) {
	var r req.CouponTemplateUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateCouponTemplate(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetCouponTemplatePage 模板分页
func (h *CouponHandler) GetCouponTemplatePage(c *gin.Context) {
	var r req.CouponTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	list, err := h.svc.GetCouponTemplatePage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, list)
}

// GetCouponPage 发放记录
func (h *CouponHandler) GetCouponPage(c *gin.Context) {
	var r req.CouponPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	list, err := h.svc.GetCouponPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, list)
}
