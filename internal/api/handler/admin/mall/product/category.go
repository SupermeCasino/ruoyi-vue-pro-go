package product

import (
	"strconv"

	product2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductCategoryHandler struct {
	svc *product.ProductCategoryService
}

func NewProductCategoryHandler(svc *product.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{svc: svc}
}

// CreateCategory 创建商品分类
func (h *ProductCategoryHandler) CreateCategory(c *gin.Context) {
	var r product2.ProductCategoryCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateCategory(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateCategory 更新商品分类
func (h *ProductCategoryHandler) UpdateCategory(c *gin.Context) {
	var r product2.ProductCategoryUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateCategory(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteCategory 删除商品分类
func (h *ProductCategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteCategory(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetCategory 获得商品分类
func (h *ProductCategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetCategory(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetCategoryList 获得商品分类列表
func (h *ProductCategoryHandler) GetCategoryList(c *gin.Context) {
	var r product2.ProductCategoryListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetCategoryList(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
