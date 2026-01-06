package iot

import (
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建设备分组
func (h *DeviceGroupHandler) Create(c *gin.Context) {
	var r iot2.IotDeviceGroupSaveReqVO
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

// Update 更新设备分组
func (h *DeviceGroupHandler) Update(c *gin.Context) {
	var r iot2.IotDeviceGroupSaveReqVO
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

// Delete 删除设备分组
func (h *DeviceGroupHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取设备分组
func (h *DeviceGroupHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	group, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if group == nil {
		response.WriteSuccess(c, nil)
		return
	}
	resp := &iot2.IotDeviceGroupRespVO{
		ID:          group.ID,
		Name:        group.Name,
		Status:      group.Status,
		Description: group.Description,
		CreateTime:  group.CreateTime,
	}
	response.WriteSuccess(c, resp)
}

// Page 获取设备分组分页
func (h *DeviceGroupHandler) Page(c *gin.Context) {
	var r iot2.IotDeviceGroupPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotDeviceGroupRespVO, 0, len(page.List))
	for _, item := range page.List {
		list = append(list, &iot2.IotDeviceGroupRespVO{
			ID:          item.ID,
			Name:        item.Name,
			Status:      item.Status,
			Description: item.Description,
			CreateTime:  item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// SimpleList 获取设备分组简单列表
func (h *DeviceGroupHandler) SimpleList(c *gin.Context) {
	// 查询状态为启用(0)的设备分组
	list, err := h.svc.GetListByStatus(c, 0)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 只返回 id 和 name 字段
	result := make([]*iot2.IotDeviceGroupRespVO, 0, len(list))
	for _, item := range list {
		result = append(result, &iot2.IotDeviceGroupRespVO{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	response.WriteSuccess(c, result)
}
