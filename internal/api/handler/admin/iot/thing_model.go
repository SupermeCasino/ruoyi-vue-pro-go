package iot

import (
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建物模型
func (h *ThingModelHandler) Create(c *gin.Context) {
	var r iot2.IotThingModelSaveReqVO
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

// Update 更新物模型
func (h *ThingModelHandler) Update(c *gin.Context) {
	var r iot2.IotThingModelSaveReqVO
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

// Delete 删除物模型
func (h *ThingModelHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取物模型
func (h *ThingModelHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tm, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if tm == nil {
		response.WriteSuccess(c, nil)
		return
	}

	resp := &iot2.IotThingModelRespVO{
		ID:          tm.ID,
		Identifier:  tm.Identifier,
		Name:        tm.Name,
		Description: tm.Description,
		ProductID:   tm.ProductID,
		ProductKey:  tm.ProductKey,
		Type:        tm.Type,
		CreateTime:  tm.CreateTime,
	}

	// 使用 Data() 获取值并取地址
	if prop := tm.Property.Data(); prop.Identifier != "" {
		resp.Property = &prop
	}
	if event := tm.Event.Data(); event.Identifier != "" {
		resp.Event = &event
	}
	if service := tm.Service.Data(); service.Identifier != "" {
		resp.Service = &service
	}

	response.WriteSuccess(c, resp)
}

// List 获取物模型列表
func (h *ThingModelHandler) List(c *gin.Context) {
	var r iot2.IotThingModelListReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	list, err := h.svc.GetList(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	respList := make([]*iot2.IotThingModelRespVO, 0, len(list))
	for _, item := range list {
		respList = append(respList, convertToRespVO(item))
	}
	response.WriteSuccess(c, respList)
}

// Page 获取物模型分页
func (h *ThingModelHandler) Page(c *gin.Context) {
	var r iot2.IotThingModelPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotThingModelRespVO, 0, len(page.List))
	for _, item := range page.List {
		list = append(list, convertToRespVO(item))
	}
	response.WritePage(c, page.Total, list)
}

// GetTSL 获取物模型 TSL
func (h *ThingModelHandler) GetTSL(c *gin.Context) {
	productIdStr := c.Query("productId")
	productId, _ := strconv.ParseInt(productIdStr, 10, 64)
	tsl, err := h.svc.GetTSL(c, productId)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, tsl)
}


// convertToRespVO 转换数据库实体为响应 VO
func convertToRespVO(tm *model.IotThingModelDO) *iot2.IotThingModelRespVO {
resp := &iot2.IotThingModelRespVO{
ID:          tm.ID,
Identifier:  tm.Identifier,
Name:        tm.Name,
Description: tm.Description,
ProductID:   tm.ProductID,
ProductKey:  tm.ProductKey,
Type:        tm.Type,
CreateTime:  tm.CreateTime,
}
if prop := tm.Property.Data(); prop.Identifier != "" {
resp.Property = &prop
}
if event := tm.Event.Data(); event.Identifier != "" {
resp.Event = &event
}
if service := tm.Service.Data(); service.Identifier != "" {
resp.Service = &service
}
return resp
}
