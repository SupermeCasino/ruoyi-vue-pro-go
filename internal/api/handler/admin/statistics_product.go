package admin

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/pkg/excel"
	"backend-go/internal/service"
	productService "backend-go/internal/service/product"

	"github.com/gin-gonic/gin"
)

// ProductStatisticsHandler 商品统计处理器
type ProductStatisticsHandler struct {
	productStatisticsService service.ProductStatisticsService
	productSpuService        *productService.ProductSpuService
}

// NewProductStatisticsHandler 创建商品统计处理器
func NewProductStatisticsHandler(
	productStatisticsService service.ProductStatisticsService,
	productSpuService *productService.ProductSpuService,
) *ProductStatisticsHandler {
	return &ProductStatisticsHandler{
		productStatisticsService: productStatisticsService,
		productSpuService:        productSpuService,
	}
}

// GetProductStatisticsAnalyse 获得商品统计分析
// GET /statistics/product/analyse
func (h *ProductStatisticsHandler) GetProductStatisticsAnalyse(c *gin.Context) {
	var reqVO req.ProductStatisticsReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		core.WriteError(c, core.ParamErrCode, err.Error())
		return
	}

	result, err := h.productStatisticsService.GetProductStatisticsAnalyse(c, &reqVO)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	core.WriteSuccess(c, result)
}

// GetProductStatisticsList 获得商品统计明细
// GET /statistics/product/list
func (h *ProductStatisticsHandler) GetProductStatisticsList(c *gin.Context) {
	var reqVO req.ProductStatisticsReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		core.WriteError(c, core.ParamErrCode, err.Error())
		return
	}

	result, err := h.productStatisticsService.GetProductStatisticsList(c, &reqVO)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	// 拼接商品信息
	spuIds := make([]int64, len(result))
	for i, item := range result {
		spuIds[i] = item.SpuID
	}
	if len(spuIds) > 0 {
		spuList, err := h.productSpuService.GetSpuList(c, spuIds)
		if err != nil {
			core.WriteError(c, core.ServerErrCode, err.Error())
			return
		}
		spuMap := make(map[int64]*resp.ProductSpuResp)
		for _, spu := range spuList {
			spuMap[spu.ID] = spu
		}
		for _, item := range result {
			if spu, ok := spuMap[item.SpuID]; ok {
				item.Name = spu.Name
				item.PicUrl = spu.PicURL
			}
		}
	}

	core.WriteSuccess(c, result)
}

// GetProductStatisticsRankPage 获得商品统计排行榜分页
// GET /statistics/product/rank-page
func (h *ProductStatisticsHandler) GetProductStatisticsRankPage(c *gin.Context) {
	var reqVO req.ProductStatisticsReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		core.WriteError(c, core.ParamErrCode, err.Error())
		return
	}

	var pageParam core.PageParam
	if err := c.ShouldBindQuery(&pageParam); err != nil {
		core.WriteError(c, core.ParamErrCode, err.Error())
		return
	}

	// 1. 获取统计数据
	pageResult, err := h.productStatisticsService.GetProductStatisticsRankPage(c, &reqVO, &pageParam)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	// 2. 拼接商品信息
	// List interface{} -> []*resp.ProductStatisticsRespVO
	// Note: The service currently returns interface{}. We need to assert it.
	// In Go, since we did manual pagination in service returning []interface{}, we need to cast back.
	// Actually, modifying service to return *core.PageResult[*resp.ProductStatisticsRespVO] would be better,
	// but reusing interface{} PageResult is common in this codebase for loose coupling.
	// Let's iterate and collect IDs.

	// Wait, internal/core/PageResult is generic in newer Go versions or just struct?
	// Looking at service code: return &core.PageResult[interface{}]...
	// So it is generic.
	// But `interface{}` is tricky.

	// Using lo.Map or manual loop
	spuIds := []int64{}

	// Cast interface{} back to *resp.ProductStatisticsRespVO
	// The service implementation put *resp.ProductStatisticsRespVO into the list.
	realList := make([]*resp.ProductStatisticsRespVO, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		if val, ok := item.(*resp.ProductStatisticsRespVO); ok {
			realList = append(realList, val)
			spuIds = append(spuIds, val.SpuID)
		}
	}

	if len(spuIds) > 0 {
		spuList, err := h.productSpuService.GetSpuList(c, spuIds)
		if err != nil {
			core.WriteError(c, core.ServerErrCode, err.Error())
			return
		}
		spuMap := make(map[int64]*resp.ProductSpuResp)
		for _, spu := range spuList {
			spuMap[spu.ID] = spu
		}
		for _, item := range realList {
			if spu, ok := spuMap[item.SpuID]; ok {
				item.Name = spu.Name
				item.PicUrl = spu.PicURL
			}
		}
	}

	// Update the list in pageResult with enriched data (pointers modified in place, but need to assign back to be safe if we created new slice)
	// Actually we modified the pointers in realList, which point to the same objects as pageResult.List[i].
	// So pageResult.List is already updated? Yes.

	core.WriteSuccess(c, pageResult)
}

// ExportProductStatisticsExcel 导出商品统计 Excel
// GET /statistics/product/export-excel
func (h *ProductStatisticsHandler) ExportProductStatisticsExcel(c *gin.Context) {
	var reqVO req.ProductStatisticsReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		core.WriteError(c, core.ParamErrCode, err.Error())
		return
	}

	// 1. 查询数据
	list, err := h.productStatisticsService.GetProductStatisticsList(c, &reqVO)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	// 2. 拼接商品信息
	spuIds := make([]int64, len(list))
	for i, item := range list {
		spuIds[i] = item.SpuID
	}
	if len(spuIds) > 0 {
		spuList, err := h.productSpuService.GetSpuList(c, spuIds)
		if err != nil {
			core.WriteError(c, core.ServerErrCode, err.Error())
			return
		}
		spuMap := make(map[int64]*resp.ProductSpuResp)
		for _, spu := range spuList {
			spuMap[spu.ID] = spu
		}
		for _, item := range list {
			if spu, ok := spuMap[item.SpuID]; ok {
				item.Name = spu.Name
				item.PicUrl = spu.PicURL
			}
		}
	}

	// 3. 导出 Excel
	if err = excel.WriteExcel(c, "商品分析.xlsx", "数据", list); err != nil {
		core.WriteError(c, core.ServerErrCode, "导出 Excel 失败: "+err.Error())
		return
	}
}
