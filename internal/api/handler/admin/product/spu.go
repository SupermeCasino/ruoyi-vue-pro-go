package product

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/pkg/excel"
	"backend-go/internal/service/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductSpuHandler struct {
	svc *product.ProductSpuService
}

func NewProductSpuHandler(svc *product.ProductSpuService) *ProductSpuHandler {
	return &ProductSpuHandler{svc: svc}
}

// CreateSpu 创建 SPU
// @Router /admin-api/product/spu/create [post]
func (h *ProductSpuHandler) CreateSpu(c *gin.Context) {
	var r req.ProductSpuSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateSpu(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateSpu 更新 SPU
// @Router /admin-api/product/spu/update [put]
func (h *ProductSpuHandler) UpdateSpu(c *gin.Context) {
	var r req.ProductSpuSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateSpu(c, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// UpdateSpuStatus 更新 SPU 状态
// @Router /admin-api/product/spu/update-status [put]
func (h *ProductSpuHandler) UpdateSpuStatus(c *gin.Context) {
	var r req.ProductSpuUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateSpuStatus(c, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteSpu 删除 SPU
// @Router /admin-api/product/spu/delete [delete]
func (h *ProductSpuHandler) DeleteSpu(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.DeleteSpu(c, id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// GetSpuDetail 获得 SPU 详情
// @Router /admin-api/product/spu/get-detail [get]
func (h *ProductSpuHandler) GetSpuDetail(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetSpuDetail(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetSpuPage 获得 SPU 分页
// @Router /admin-api/product/spu/page [get]
func (h *ProductSpuHandler) GetSpuPage(c *gin.Context) {
	var r req.ProductSpuPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetSpuPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetTabsCount 获得 SPU Tab 统计
// @Router /admin-api/product/spu/get-count [get]
func (h *ProductSpuHandler) GetTabsCount(c *gin.Context) {
	res, err := h.svc.GetTabsCount(c)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetSpuSimpleList 获得 SPU 精简列表
// @Router /admin-api/product/spu/list-all-simple [get]
func (h *ProductSpuHandler) GetSpuSimpleList(c *gin.Context) {
	res, err := h.svc.GetSpuSimpleList(c)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetSpuList 根据 ID 列表获得 SPU 详情列表
// @Router /admin-api/product/spu/list [get]
func (h *ProductSpuHandler) GetSpuList(c *gin.Context) {
	var r req.ProductSpuListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetSpuList(c, r.SpuIDs)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// ExportSpuList 导出商品列表
// @Router /admin-api/product/spu/export [get]
func (h *ProductSpuHandler) ExportSpuList(c *gin.Context) {
	var r req.ProductSpuPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	// 导出不分页
	r.PageNo = 1
	r.PageSize = 10000

	res, err := h.svc.GetSpuPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// 转换为导出 VO
	exportList := make([]*resp.ProductSpuExportVO, len(res.List))
	for i, spu := range res.List {
		exportList[i] = &resp.ProductSpuExportVO{
			ID:         spu.ID,
			Name:       spu.Name,
			CategoryID: spu.CategoryID,
			Price:      spu.Price,
			Stock:      spu.Stock,
			Status:     spu.Status,
			SalesCount: spu.SalesCount,
		}
	}

	if err := excel.WriteExcel(c, "商品列表.xlsx", "数据", exportList); err != nil {
		c.Error(err)
		return
	}
}
