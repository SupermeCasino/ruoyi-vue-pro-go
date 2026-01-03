package trade

import (
	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeliveryExpressHandler struct {
	svc    *trade.DeliveryExpressService
	logger *zap.Logger
}

func NewDeliveryExpressHandler(svc *trade.DeliveryExpressService, logger *zap.Logger) *DeliveryExpressHandler {
	return &DeliveryExpressHandler{
		svc:    svc,
		logger: logger,
	}
}

// CreateDeliveryExpress 创建物流公司
func (h *DeliveryExpressHandler) CreateDeliveryExpress(c *gin.Context) {
	var r trade2.DeliveryExpressSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.svc.CreateDeliveryExpress(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("创建物流公司失败", zap.Error(err))
		response.WriteError(c, 500, "创建失败")
		return
	}

	response.WriteSuccess(c, id)
}

// UpdateDeliveryExpress 更新物流公司
func (h *DeliveryExpressHandler) UpdateDeliveryExpress(c *gin.Context) {
	var r trade2.DeliveryExpressSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdateDeliveryExpress(c.Request.Context(), &r); err != nil {
		h.logger.Error("更新物流公司失败", zap.Error(err))
		response.WriteError(c, 500, "更新失败")
		return
	}

	response.WriteSuccess(c, true)
}

// DeleteDeliveryExpress 删除物流公司
func (h *DeliveryExpressHandler) DeleteDeliveryExpress(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteDeliveryExpress(c.Request.Context(), id); err != nil {
		h.logger.Error("删除物流公司失败", zap.Error(err))
		response.WriteError(c, 500, "删除失败")
		return
	}

	response.WriteSuccess(c, true)
}

// GetDeliveryExpress 获取物流公司
func (h *DeliveryExpressHandler) GetDeliveryExpress(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	express, err := h.svc.GetDeliveryExpress(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取物流公司失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	response.WriteSuccess(c, trade2.DeliveryExpressResp{
		ID:         express.ID,
		Code:       express.Code,
		Name:       express.Name,
		Logo:       express.Logo,
		Sort:       express.Sort,
		Status:     express.Status,
		CreateTime: express.CreateTime,
	})
}

// GetDeliveryExpressPage 获取物流公司分页
func (h *DeliveryExpressHandler) GetDeliveryExpressPage(c *gin.Context) {
	var r trade2.DeliveryExpressPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.svc.GetDeliveryExpressPage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("获取物流公司分页失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	list := make([]trade2.DeliveryExpressResp, len(page.List))
	for i, item := range page.List {
		list[i] = trade2.DeliveryExpressResp{
			ID:         item.ID,
			Code:       item.Code,
			Name:       item.Name,
			Logo:       item.Logo,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[trade2.DeliveryExpressResp]{
		List:  list,
		Total: page.Total,
	})
}

// GetSimpleDeliveryExpressList 获取物流公司精简列表
func (h *DeliveryExpressHandler) GetSimpleDeliveryExpressList(c *gin.Context) {
	list, err := h.svc.GetSimpleDeliveryExpressList(c.Request.Context())
	if err != nil {
		h.logger.Error("获取物流公司列表失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	res := make([]trade2.DeliveryExpressResp, len(list))
	for i, item := range list {
		res[i] = trade2.DeliveryExpressResp{
			ID:         item.ID,
			Code:       item.Code,
			Name:       item.Name,
			Logo:       item.Logo,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreateTime,
		}
	}
	response.WriteSuccess(c, res)
}

// ExportDeliveryExpress 导出物流公司列表
// @Router /admin-api/trade/delivery/express/export-excel [get]
func (h *DeliveryExpressHandler) ExportDeliveryExpress(c *gin.Context) {
	var r trade2.DeliveryExpressPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	r.PageNo = 1
	r.PageSize = 10000 // 导出最大数量

	page, err := h.svc.GetDeliveryExpressPage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("导出物流公司失败", zap.Error(err))
		response.WriteError(c, 500, "导出失败")
		return
	}

	list := make([]*trade2.DeliveryExpressExcelVO, len(page.List))
	for i, item := range page.List {
		list[i] = &trade2.DeliveryExpressExcelVO{
			ID:         item.ID,
			Code:       item.Code,
			Name:       item.Name,
			Logo:       item.Logo,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreateTime,
		}
	}

	if err := excel.WriteExcel(c, "快递公司.xls", "数据", list); err != nil {
		h.logger.Error("导出Excel失败", zap.Error(err))
		response.WriteError(c, 500, "导出失败")
	}
}

type DeliveryPickUpStoreHandler struct {
	svc    *trade.DeliveryPickUpStoreService
	logger *zap.Logger
}

func NewDeliveryPickUpStoreHandler(svc *trade.DeliveryPickUpStoreService, logger *zap.Logger) *DeliveryPickUpStoreHandler {
	return &DeliveryPickUpStoreHandler{
		svc:    svc,
		logger: logger,
	}
}

// CreateDeliveryPickUpStore 创建自提门店
func (h *DeliveryPickUpStoreHandler) CreateDeliveryPickUpStore(c *gin.Context) {
	var r trade2.DeliveryPickUpStoreSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.svc.CreateDeliveryPickUpStore(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("创建自提门店失败", zap.Error(err))
		response.WriteError(c, 500, "创建失败")
		return
	}

	response.WriteSuccess(c, id)
}

// UpdateDeliveryPickUpStore 更新自提门店
func (h *DeliveryPickUpStoreHandler) UpdateDeliveryPickUpStore(c *gin.Context) {
	var r trade2.DeliveryPickUpStoreSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdateDeliveryPickUpStore(c.Request.Context(), &r); err != nil {
		h.logger.Error("更新自提门店失败", zap.Error(err))
		response.WriteError(c, 500, "更新失败")
		return
	}

	response.WriteSuccess(c, true)
}

// DeleteDeliveryPickUpStore 删除自提门店
func (h *DeliveryPickUpStoreHandler) DeleteDeliveryPickUpStore(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteDeliveryPickUpStore(c.Request.Context(), id); err != nil {
		h.logger.Error("删除自提门店失败", zap.Error(err))
		response.WriteError(c, 500, "删除失败")
		return
	}

	response.WriteSuccess(c, true)
}

// GetDeliveryPickUpStore 获取自提门店
func (h *DeliveryPickUpStoreHandler) GetDeliveryPickUpStore(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	store, err := h.svc.GetDeliveryPickUpStore(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取自提门店失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	response.WriteSuccess(c, trade2.DeliveryPickUpStoreResp{
		ID:            store.ID,
		Name:          store.Name,
		Introduction:  store.Introduction,
		Phone:         store.Phone,
		AreaID:        store.AreaID,
		DetailAddress: store.DetailAddress,
		Logo:          store.Logo,
		OpeningTime:   store.OpeningTime,
		ClosingTime:   store.ClosingTime,
		Latitude:      store.Latitude,
		Longitude:     store.Longitude,
		Status:        store.Status,
		CreateTime:    store.CreateTime,
	})
}

// GetDeliveryPickUpStorePage 获取自提门店分页
func (h *DeliveryPickUpStoreHandler) GetDeliveryPickUpStorePage(c *gin.Context) {
	var r trade2.DeliveryPickUpStorePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.svc.GetDeliveryPickUpStorePage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("获取自提门店分页失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	list := make([]trade2.DeliveryPickUpStoreResp, len(page.List))
	for i, item := range page.List {
		list[i] = trade2.DeliveryPickUpStoreResp{
			ID:            item.ID,
			Name:          item.Name,
			Introduction:  item.Introduction,
			Phone:         item.Phone,
			AreaID:        item.AreaID,
			DetailAddress: item.DetailAddress,
			Logo:          item.Logo,
			OpeningTime:   item.OpeningTime,
			ClosingTime:   item.ClosingTime,
			Latitude:      item.Latitude,
			Longitude:     item.Longitude,
			Status:        item.Status,
			CreateTime:    item.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[trade2.DeliveryPickUpStoreResp]{
		List:  list,
		Total: page.Total,
	})
}

// GetSimpleDeliveryPickUpStoreList 获取自提门店精简列表
func (h *DeliveryPickUpStoreHandler) GetSimpleDeliveryPickUpStoreList(c *gin.Context) {
	list, err := h.svc.GetSimpleDeliveryPickUpStoreList(c.Request.Context())
	if err != nil {
		h.logger.Error("获取自提门店列表失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	res := make([]trade2.DeliveryPickUpStoreResp, len(list))
	for i, item := range list {
		res[i] = trade2.DeliveryPickUpStoreResp{
			ID:            item.ID,
			Name:          item.Name,
			Introduction:  item.Introduction,
			Phone:         item.Phone,
			AreaID:        item.AreaID,
			DetailAddress: item.DetailAddress,
			Logo:          item.Logo,
			OpeningTime:   item.OpeningTime,
			ClosingTime:   item.ClosingTime,
			Latitude:      item.Latitude,
			Longitude:     item.Longitude,
			Status:        item.Status,
			CreateTime:    item.CreateTime,
		}
	}
	response.WriteSuccess(c, res)
}

// BindDeliveryPickUpStore 绑定自提门店核销员工
func (h *DeliveryPickUpStoreHandler) BindDeliveryPickUpStore(c *gin.Context) {
	var r trade2.DeliveryPickUpBindReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	if err := h.svc.BindDeliveryPickUpStore(c.Request.Context(), &r); err != nil {
		h.logger.Error("绑定自提门店核销员工失败", zap.Error(err))
		response.WriteError(c, 500, "绑定失败")
		return
	}

	response.WriteSuccess(c, true)
}

type DeliveryExpressTemplateHandler struct {
	svc    *trade.DeliveryExpressTemplateService
	logger *zap.Logger
}

func NewDeliveryExpressTemplateHandler(svc *trade.DeliveryExpressTemplateService, logger *zap.Logger) *DeliveryExpressTemplateHandler {
	return &DeliveryExpressTemplateHandler{
		svc:    svc,
		logger: logger,
	}
}

// CreateDeliveryExpressTemplate 创建运费模板
func (h *DeliveryExpressTemplateHandler) CreateDeliveryExpressTemplate(c *gin.Context) {
	var r trade2.DeliveryFreightTemplateSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.svc.CreateDeliveryExpressTemplate(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("创建运费模板失败", zap.Error(err))
		response.WriteError(c, 500, "创建失败")
		return
	}

	response.WriteSuccess(c, id)
}

// UpdateDeliveryExpressTemplate 更新运费模板
func (h *DeliveryExpressTemplateHandler) UpdateDeliveryExpressTemplate(c *gin.Context) {
	var r trade2.DeliveryFreightTemplateSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdateDeliveryExpressTemplate(c.Request.Context(), &r); err != nil {
		h.logger.Error("更新运费模板失败", zap.Error(err))
		response.WriteError(c, 500, "更新失败")
		return
	}

	response.WriteSuccess(c, true)
}

// DeleteDeliveryExpressTemplate 删除运费模板
func (h *DeliveryExpressTemplateHandler) DeleteDeliveryExpressTemplate(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteDeliveryExpressTemplate(c.Request.Context(), id); err != nil {
		h.logger.Error("删除运费模板失败", zap.Error(err))
		response.WriteError(c, 500, "删除失败")
		return
	}

	response.WriteSuccess(c, true)
}

// GetDeliveryExpressTemplate 获取运费模板详情
func (h *DeliveryExpressTemplateHandler) GetDeliveryExpressTemplate(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	template, err := h.svc.GetDeliveryExpressTemplate(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取运费模板失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	response.WriteSuccess(c, template)
}

// GetDeliveryExpressTemplatePage 获取运费模板分页
func (h *DeliveryExpressTemplateHandler) GetDeliveryExpressTemplatePage(c *gin.Context) {
	var r trade2.DeliveryFreightTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.svc.GetDeliveryExpressTemplatePage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("获取运费模板分页失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}

	list := make([]trade2.DeliveryFreightTemplateResp, len(page.List))
	for i, item := range page.List {
		list[i] = trade2.DeliveryFreightTemplateResp{
			ID:         item.ID,
			Name:       item.Name,
			ChargeMode: item.ChargeMode,
			Sort:       item.Sort,
			CreateTime: item.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[trade2.DeliveryFreightTemplateResp]{
		List:  list,
		Total: page.Total,
	})
}

// GetSimpleDeliveryExpressTemplateList 获取所有运费模板精简列表
func (h *DeliveryExpressTemplateHandler) GetSimpleDeliveryExpressTemplateList(c *gin.Context) {
	list, err := h.svc.GetSimpleDeliveryExpressTemplateList(c.Request.Context())
	if err != nil {
		h.logger.Error("获取运费模板列表失败", zap.Error(err))
		response.WriteError(c, 500, "获取失败")
		return
	}
	response.WriteSuccess(c, list)
}
