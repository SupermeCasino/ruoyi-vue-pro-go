package iot

import (
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建产品
func (h *ProductHandler) Create(c *gin.Context) {
	var r iot2.IotProductSaveReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	id, err := h.svc.Create(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// Update 更新产品
func (h *ProductHandler) Update(c *gin.Context) {
	var r iot2.IotProductSaveReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	if err := h.svc.Update(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateStatus 更新产品状态
func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Query("id")
	statusStr := c.Query("status")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	status, _ := strconv.ParseInt(statusStr, 10, 8)

	if err := h.svc.UpdateStatus(c, id, int8(status)); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Delete 删除产品
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取产品
func (h *ProductHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	product, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if product == nil {
		response.WriteSuccess(c, nil)
		return
	}
	response.WriteSuccess(c, product)
}

// GetByKey 获取产品 (by Key)
func (h *ProductHandler) GetByKey(c *gin.Context) {
	productKey := c.Query("productKey")
	product, err := h.svc.GetByKey(c, productKey)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, product)
}

// SimpleList 获取产品精简列表
func (h *ProductHandler) SimpleList(c *gin.Context) {
	list, err := h.svc.GetSimpleList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

// Page 获取产品分页
func (h *ProductHandler) Page(c *gin.Context) {
	var r iot2.IotProductPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotProductRespVO, 0, len(page.List))
	for _, item := range page.List {
		list = append(list, &iot2.IotProductRespVO{
			ID:           item.ID,
			Name:         item.Name,
			ProductKey:   item.ProductKey,
			CategoryID:   item.CategoryID,
			Icon:         item.Icon,
			PicURL:       item.PicURL,
			Description:  item.Description,
			Status:       item.Status,
			DeviceType:   item.DeviceType,
			NetType:      item.NetType,
			LocationType: item.LocationType,
			CodecType:    item.CodecType,
			CreateTime:   item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}
