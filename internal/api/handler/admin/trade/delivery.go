package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/pkg/excel"
	"backend-go/internal/service/trade"

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
	var r req.DeliveryExpressSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.svc.CreateDeliveryExpress(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("创建物流公司失败", zap.Error(err))
		core.WriteError(c, 500, "创建失败")
		return
	}

	core.WriteSuccess(c, id)
}

// UpdateDeliveryExpress 更新物流公司
func (h *DeliveryExpressHandler) UpdateDeliveryExpress(c *gin.Context) {
	var r req.DeliveryExpressSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdateDeliveryExpress(c.Request.Context(), &r); err != nil {
		h.logger.Error("更新物流公司失败", zap.Error(err))
		core.WriteError(c, 500, "更新失败")
		return
	}

	core.WriteSuccess(c, true)
}

// DeleteDeliveryExpress 删除物流公司
func (h *DeliveryExpressHandler) DeleteDeliveryExpress(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteDeliveryExpress(c.Request.Context(), id); err != nil {
		h.logger.Error("删除物流公司失败", zap.Error(err))
		core.WriteError(c, 500, "删除失败")
		return
	}

	core.WriteSuccess(c, true)
}

// GetDeliveryExpress 获取物流公司
func (h *DeliveryExpressHandler) GetDeliveryExpress(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	express, err := h.svc.GetDeliveryExpress(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取物流公司失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}

	core.WriteSuccess(c, resp.DeliveryExpressResp{
		ID:         express.ID,
		Code:       express.Code,
		Name:       express.Name,
		Logo:       express.Logo,
		Sort:       express.Sort,
		Status:     express.Status,
		CreateTime: express.CreatedAt,
	})
}

// GetDeliveryExpressPage 获取物流公司分页
func (h *DeliveryExpressHandler) GetDeliveryExpressPage(c *gin.Context) {
	var r req.DeliveryExpressPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.svc.GetDeliveryExpressPage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("获取物流公司分页失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}

	list := make([]resp.DeliveryExpressResp, len(page.List))
	for i, item := range page.List {
		list[i] = resp.DeliveryExpressResp{
			ID:         item.ID,
			Code:       item.Code,
			Name:       item.Name,
			Logo:       item.Logo,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.DeliveryExpressResp]{
		List:  list,
		Total: page.Total,
	})
}

// GetSimpleDeliveryExpressList 获取物流公司精简列表
func (h *DeliveryExpressHandler) GetSimpleDeliveryExpressList(c *gin.Context) {
	list, err := h.svc.GetSimpleDeliveryExpressList(c.Request.Context())
	if err != nil {
		h.logger.Error("获取物流公司列表失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}

	res := make([]resp.DeliveryExpressResp, len(list))
	for i, item := range list {
		res[i] = resp.DeliveryExpressResp{
			ID:         item.ID,
			Code:       item.Code,
			Name:       item.Name,
			Logo:       item.Logo,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreatedAt,
		}
	}
	core.WriteSuccess(c, res)
}

// ExportDeliveryExpress 导出物流公司列表
// @Router /admin-api/trade/delivery/express/export-excel [get]
func (h *DeliveryExpressHandler) ExportDeliveryExpress(c *gin.Context) {
	var r req.DeliveryExpressPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	r.PageNo = 1
	r.PageSize = 10000 // 导出最大数量

	page, err := h.svc.GetDeliveryExpressPage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("导出物流公司失败", zap.Error(err))
		core.WriteError(c, 500, "导出失败")
		return
	}

	list := make([]*resp.DeliveryExpressExcelVO, len(page.List))
	for i, item := range page.List {
		list[i] = &resp.DeliveryExpressExcelVO{
			ID:         item.ID,
			Code:       item.Code,
			Name:       item.Name,
			Logo:       item.Logo,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreatedAt,
		}
	}

	if err := excel.WriteExcel(c, "快递公司.xls", "数据", list); err != nil {
		h.logger.Error("导出Excel失败", zap.Error(err))
		core.WriteError(c, 500, "导出失败")
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
	var r req.DeliveryPickUpStoreSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.svc.CreateDeliveryPickUpStore(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("创建自提门店失败", zap.Error(err))
		core.WriteError(c, 500, "创建失败")
		return
	}

	core.WriteSuccess(c, id)
}

// UpdateDeliveryPickUpStore 更新自提门店
func (h *DeliveryPickUpStoreHandler) UpdateDeliveryPickUpStore(c *gin.Context) {
	var r req.DeliveryPickUpStoreSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdateDeliveryPickUpStore(c.Request.Context(), &r); err != nil {
		h.logger.Error("更新自提门店失败", zap.Error(err))
		core.WriteError(c, 500, "更新失败")
		return
	}

	core.WriteSuccess(c, true)
}

// DeleteDeliveryPickUpStore 删除自提门店
func (h *DeliveryPickUpStoreHandler) DeleteDeliveryPickUpStore(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteDeliveryPickUpStore(c.Request.Context(), id); err != nil {
		h.logger.Error("删除自提门店失败", zap.Error(err))
		core.WriteError(c, 500, "删除失败")
		return
	}

	core.WriteSuccess(c, true)
}

// GetDeliveryPickUpStore 获取自提门店
func (h *DeliveryPickUpStoreHandler) GetDeliveryPickUpStore(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	store, err := h.svc.GetDeliveryPickUpStore(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取自提门店失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}

	core.WriteSuccess(c, resp.DeliveryPickUpStoreResp{
		ID:            store.ID,
		Name:          store.Name,
		Introduction:  store.Introduction,
		Phone:         store.Phone,
		AreaID:        store.AreaID,
		DetailAddress: store.DetailAddress,
		Logo:          store.Logo,
		Latitude:      store.Latitude,
		Longitude:     store.Longitude,
		Status:        store.Status,
		Sort:          store.Sort,
		CreateTime:    store.CreatedAt,
	})
}

// GetDeliveryPickUpStorePage 获取自提门店分页
func (h *DeliveryPickUpStoreHandler) GetDeliveryPickUpStorePage(c *gin.Context) {
	var r req.DeliveryPickUpStorePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.svc.GetDeliveryPickUpStorePage(c.Request.Context(), &r)
	if err != nil {
		h.logger.Error("获取自提门店分页失败", zap.Error(err))
		core.WriteError(c, 500, "获取失败")
		return
	}

	list := make([]resp.DeliveryPickUpStoreResp, len(page.List))
	for i, item := range page.List {
		list[i] = resp.DeliveryPickUpStoreResp{
			ID:            item.ID,
			Name:          item.Name,
			Introduction:  item.Introduction,
			Phone:         item.Phone,
			AreaID:        item.AreaID,
			DetailAddress: item.DetailAddress,
			Logo:          item.Logo,
			Latitude:      item.Latitude,
			Longitude:     item.Longitude,
			Status:        item.Status,
			Sort:          item.Sort,
			CreateTime:    item.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.DeliveryPickUpStoreResp]{
		List:  list,
		Total: page.Total,
	})
}
