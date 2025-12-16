package promotion

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type DiyTemplateHandler struct {
	svc promotion.DiyTemplateService
}

func NewDiyTemplateHandler(svc promotion.DiyTemplateService) *DiyTemplateHandler {
	return &DiyTemplateHandler{svc: svc}
}

func (h *DiyTemplateHandler) CreateDiyTemplate(c *gin.Context) {
	var r req.DiyTemplateCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateDiyTemplate(c, r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

func (h *DiyTemplateHandler) UpdateDiyTemplate(c *gin.Context) {
	var r req.DiyTemplateUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.UpdateDiyTemplate(c, r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *DiyTemplateHandler) DeleteDiyTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.DeleteDiyTemplate(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

func (h *DiyTemplateHandler) GetDiyTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.GetDiyTemplate(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

func (h *DiyTemplateHandler) GetDiyTemplatePage(c *gin.Context) {
	var r req.DiyTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.GetDiyTemplatePage(c, r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// GetDiyTemplateProperty 获得装修模板属性
func (h *DiyTemplateHandler) GetDiyTemplateProperty(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.GetDiyTemplateProperty(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// UseDiyTemplate 使用装修模板
// Java: DiyTemplateController#useDiyTemplate
func (h *DiyTemplateHandler) UseDiyTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.UseDiyTemplate(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// UpdateDiyTemplateProperty 更新装修模板属性
// Java: DiyTemplateController#updateDiyTemplateProperty
func (h *DiyTemplateHandler) UpdateDiyTemplateProperty(c *gin.Context) {
	var r req.DiyTemplatePropertyUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.UpdateDiyTemplateProperty(c, r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}
