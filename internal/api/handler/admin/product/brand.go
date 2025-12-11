package product

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductBrandHandler struct {
	svc *product.ProductBrandService
}

func NewProductBrandHandler(svc *product.ProductBrandService) *ProductBrandHandler {
	return &ProductBrandHandler{svc: svc}
}

// CreateBrand 创建品牌
// @Router /admin-api/product/brand/create [post]
func (h *ProductBrandHandler) CreateBrand(c *gin.Context) {
	var r req.ProductBrandCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateBrand(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateBrand 更新品牌
// @Router /admin-api/product/brand/update [put]
func (h *ProductBrandHandler) UpdateBrand(c *gin.Context) {
	var r req.ProductBrandUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateBrand(c, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteBrand 删除品牌
// @Router /admin-api/product/brand/delete [delete]
func (h *ProductBrandHandler) DeleteBrand(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.DeleteBrand(c, id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// GetBrand 获得品牌
// @Router /admin-api/product/brand/get [get]
func (h *ProductBrandHandler) GetBrand(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetBrand(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetBrandPage 获得品牌分页
// @Router /admin-api/product/brand/page [get]
func (h *ProductBrandHandler) GetBrandPage(c *gin.Context) {
	var r req.ProductBrandPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetBrandPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetBrandList 获得品牌列表
// @Router /admin-api/product/brand/list [get]
// @Router /admin-api/product/brand/list-all-simple [get]
func (h *ProductBrandHandler) GetBrandList(c *gin.Context) {
	var r req.ProductBrandListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetBrandList(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}
