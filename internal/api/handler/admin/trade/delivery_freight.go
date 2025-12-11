package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/trade"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeliveryFreightTemplateHandler struct {
	svc    *trade.DeliveryFreightTemplateService
	logger *zap.Logger
}

func NewDeliveryFreightTemplateHandler(svc *trade.DeliveryFreightTemplateService, logger *zap.Logger) *DeliveryFreightTemplateHandler {
	return &DeliveryFreightTemplateHandler{
		svc:    svc,
		logger: logger,
	}
}

// CreateDeliveryFreightTemplate 创建运费模板
func (h *DeliveryFreightTemplateHandler) CreateDeliveryFreightTemplate(c *gin.Context) {
	var r req.DeliveryFreightTemplateSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.svc.CreateDeliveryFreightTemplate(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("创建运费模板失败", zap.Error(err))
		core.WriteError(c, 500, "创建失败")
		return
	}

	core.WriteSuccess(c, id)
}

// UpdateDeliveryFreightTemplate 更新运费模板
func (h *DeliveryFreightTemplateHandler) UpdateDeliveryFreightTemplate(c *gin.Context) {
	var r req.DeliveryFreightTemplateSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdateDeliveryFreightTemplate(c.Request.Context(), &r); err != nil {
		h.logger.Error("更新运费模板失败", zap.Error(err))
		core.WriteError(c, 500, "更新失败")
		return
	}

	core.WriteSuccess(c, true)
}

// DeleteDeliveryFreightTemplate 删除运费模板
func (h *DeliveryFreightTemplateHandler) DeleteDeliveryFreightTemplate(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteDeliveryFreightTemplate(c.Request.Context(), id); err != nil {
		h.logger.Error("删除运费模板失败", zap.Error(err))
		core.WriteError(c, 500, "删除失败")
		return
	}

	core.WriteSuccess(c, true)
}

// GetDeliveryFreightTemplate 获取运费模板详情
func (h *DeliveryFreightTemplateHandler) GetDeliveryFreightTemplate(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	template, err := h.svc.GetDeliveryFreightTemplate(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取运费模板失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}

	core.WriteSuccess(c, template)
}

// GetDeliveryFreightTemplatePage 获取运费模板分页
func (h *DeliveryFreightTemplateHandler) GetDeliveryFreightTemplatePage(c *gin.Context) {
	var r req.DeliveryFreightTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.svc.GetDeliveryFreightTemplatePage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("获取运费模板分页失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}

	list := make([]resp.DeliveryFreightTemplateResp, len(page.List))
	for i, item := range page.List {
		list[i] = resp.DeliveryFreightTemplateResp{
			ID:         item.ID,
			Name:       item.Name,
			Type:       item.Type,
			ChargeMode: item.ChargeMode,
			Sort:       item.Sort,
			Status:     item.Status,
			Remark:     item.Remark,
			CreateTime: item.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.DeliveryFreightTemplateResp]{
		List:  list,
		Total: page.Total,
	})
}

// GetSimpleDeliveryFreightTemplateList 获取所有运费模板精简列表
func (h *DeliveryFreightTemplateHandler) GetSimpleDeliveryFreightTemplateList(c *gin.Context) {
	list, err := h.svc.GetSimpleDeliveryFreightTemplateList(c.Request.Context())
	if err != nil {
		h.logger.Error("获取运费模板列表失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}
	core.WriteSuccess(c, list)
}
