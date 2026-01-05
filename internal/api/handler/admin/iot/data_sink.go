package iot

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建数据目的
func (h *DataSinkHandler) Create(c *gin.Context) {
	var r iot2.IotDataSinkSaveReqVO
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

// Update 更新数据目的
func (h *DataSinkHandler) Update(c *gin.Context) {
	var r iot2.IotDataSinkSaveReqVO
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

// Delete 删除数据目的
func (h *DataSinkHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取数据目的
func (h *DataSinkHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	sink, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if sink == nil {
		response.WriteSuccess(c, nil)
		return
	}

	var config interface{}
	_ = json.Unmarshal(sink.Config, &config)

	resp := &iot2.IotDataSinkRespVO{
		ID:          sink.ID,
		Name:        sink.Name,
		Description: sink.Description,
		Status:      sink.Status,
		Type:        sink.Type,
		Config:      config,
		CreateTime:  sink.CreateTime,
	}
	response.WriteSuccess(c, resp)
}

// Page 获取数据目的分页
func (h *DataSinkHandler) Page(c *gin.Context) {
	var r iot2.IotDataSinkPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotDataSinkRespVO, 0, len(page.List))
	for _, item := range page.List {
		var config interface{}
		_ = json.Unmarshal(item.Config, &config)

		list = append(list, &iot2.IotDataSinkRespVO{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			Type:        item.Type,
			Config:      config,
			CreateTime:  item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// SimpleList 获取数据目的精简列表
func (h *DataSinkHandler) SimpleList(c *gin.Context) {
	list, err := h.svc.GetListByStatus(c, 0) // 0: 启用
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	resp := make([]*iot2.IotDataSinkRespVO, 0, len(list))
	for _, item := range list {
		resp = append(resp, &iot2.IotDataSinkRespVO{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	response.WriteSuccess(c, resp)
}
