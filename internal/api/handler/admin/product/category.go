package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductCategoryHandler struct {
	svc *product.ProductCategoryService
}

func NewProductCategoryHandler(svc *product.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{svc: svc}
}

// CreateCategory 创建商品分类
// @Router /admin-api/product/category/create [post]
func (h *ProductCategoryHandler) CreateCategory(c *gin.Context) {
	var r req.ProductCategoryCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateCategory(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateCategory 更新商品分类
// @Router /admin-api/product/category/update [put]
func (h *ProductCategoryHandler) UpdateCategory(c *gin.Context) {
	var r req.ProductCategoryUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateCategory(c, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteCategory 删除商品分类
// @Router /admin-api/product/category/delete [delete]
func (h *ProductCategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.DeleteCategory(c, id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// GetCategory 获得商品分类
// @Router /admin-api/product/category/get [get]
func (h *ProductCategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetCategory(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetCategoryList 获得商品分类列表
// @Router /admin-api/product/category/list [get]
func (h *ProductCategoryHandler) GetCategoryList(c *gin.Context) {
	var r req.ProductCategoryListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetCategoryList(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}
