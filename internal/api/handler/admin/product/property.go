package product

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductPropertyHandler struct {
	svc      *product.ProductPropertyService
	valueSvc *product.ProductPropertyValueService
}

func NewProductPropertyHandler(svc *product.ProductPropertyService, valueSvc *product.ProductPropertyValueService) *ProductPropertyHandler {
	return &ProductPropertyHandler{
		svc:      svc,
		valueSvc: valueSvc,
	}
}

func (h *ProductPropertyHandler) CreateProperty(c *gin.Context) {
	var r req.ProductPropertyCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateProperty(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *ProductPropertyHandler) UpdateProperty(c *gin.Context) {
	var r req.ProductPropertyUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateProperty(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *ProductPropertyHandler) DeleteProperty(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteProperty(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetProperty 获得属性项
// @Router /admin-api/product/property/get [get]
func (h *ProductPropertyHandler) GetProperty(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetProperty(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetPropertyPage 获得属性项分页
// @Router /admin-api/product/property/page [get]
func (h *ProductPropertyHandler) GetPropertyPage(c *gin.Context) {
	var r req.ProductPropertyPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetPropertyPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetPropertySimpleList 获得属性项精简列表
// @Router /admin-api/product/property/simple-list [get]
func (h *ProductPropertyHandler) GetPropertySimpleList(c *gin.Context) {
	var r req.ProductPropertyListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetPropertyList(c, &r) // Reusing GetPropertyList for simple-list
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// --- Value Handlers ---

// CreatePropertyValue 创建属性值
// @Router /admin-api/product/property/value/create [post]
func (h *ProductPropertyHandler) CreatePropertyValue(c *gin.Context) {
	var r req.ProductPropertyValueCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.valueSvc.CreatePropertyValue(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdatePropertyValue 更新属性值
// @Router /admin-api/product/property/value/update [put]
func (h *ProductPropertyHandler) UpdatePropertyValue(c *gin.Context) {
	var r req.ProductPropertyValueUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.valueSvc.UpdatePropertyValue(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeletePropertyValue 删除属性值
// @Router /admin-api/product/property/value/delete [delete]
func (h *ProductPropertyHandler) DeletePropertyValue(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.valueSvc.DeletePropertyValue(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetPropertyValue 获得属性值
// @Router /admin-api/product/property/value/get [get]
func (h *ProductPropertyHandler) GetPropertyValue(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.valueSvc.GetPropertyValue(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetPropertyValuePage 获得属性值分页
// @Router /admin-api/product/property/value/page [get]
func (h *ProductPropertyHandler) GetPropertyValuePage(c *gin.Context) {
	var r req.ProductPropertyValuePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.valueSvc.GetPropertyValuePage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetPropertyValueSimpleList 获得属性值精简列表
// @Router /admin-api/product/property/value/simple-list [get]
func (h *ProductPropertyHandler) GetPropertyValueSimpleList(c *gin.Context) {
	propertyIDStr := c.Query("propertyId")
	propertyID, err := strconv.ParseInt(propertyIDStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.valueSvc.GetPropertyValueListByPropertyIds(c, []int64{propertyID})
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	// 转换为简化响应 (只返回 id, name)
	type simpleVO struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	simpleList := make([]simpleVO, len(res))
	for i, v := range res {
		simpleList[i] = simpleVO{ID: v.ID, Name: v.Name}
	}
	response.WriteSuccess(c, simpleList)
}
