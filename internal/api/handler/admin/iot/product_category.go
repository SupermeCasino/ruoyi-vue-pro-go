package iot

import (
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建产品分类
func (h *ProductCategoryHandler) Create(c *gin.Context) {
	var r iot2.IotProductCategorySaveReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	id, err := h.svc.CreateProductCategory(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// Update 更新产品分类
func (h *ProductCategoryHandler) Update(c *gin.Context) {
	var r iot2.IotProductCategorySaveReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	if err := h.svc.UpdateProductCategory(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Delete 删除产品分类
func (h *ProductCategoryHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteProductCategory(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取产品分类
func (h *ProductCategoryHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	category, err := h.svc.GetProductCategory(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if category == nil {
		response.WriteSuccess(c, nil)
		return
	}
	response.WriteSuccess(c, &iot2.IotProductCategoryRespVO{
		ID:          category.ID,
		Name:        category.Name,
		Sort:        category.Sort,
		Status:      category.Status,
		Description: category.Description,
		CreateTime:  category.CreateTime,
	})
}

// Page 获取产品分类分页
func (h *ProductCategoryHandler) Page(c *gin.Context) {
	var r iot2.IotProductCategoryPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetProductCategoryPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotProductCategoryRespVO, 0, len(page.List))
	for _, item := range page.List {
		list = append(list, &iot2.IotProductCategoryRespVO{
			ID:          item.ID,
			Name:        item.Name,
			Sort:        item.Sort,
			Status:      item.Status,
			Description: item.Description,
			CreateTime:  item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// SimpleList 获取产品分类精简列表
func (h *ProductCategoryHandler) SimpleList(c *gin.Context) {
	list, err := h.svc.GetProductCategoryListByStatus(c, 0) // 0: 启用
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	resp := make([]*iot2.IotProductCategoryRespVO, 0, len(list))
	for _, item := range list {
		resp = append(resp, &iot2.IotProductCategoryRespVO{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	response.WriteSuccess(c, resp)
}
