package iot

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建设备
func (h *DeviceHandler) Create(c *gin.Context) {
	var r iot2.IotDeviceSaveReqVO
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

// Update 更新设备
func (h *DeviceHandler) Update(c *gin.Context) {
	var r iot2.IotDeviceSaveReqVO
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

// Delete 删除设备
func (h *DeviceHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteList 批量删除设备
func (h *DeviceHandler) DeleteList(c *gin.Context) {
	idsStr := c.QueryArray("ids")
	ids := make([]int64, 0, len(idsStr))
	for _, s := range idsStr {
		id, _ := strconv.ParseInt(s, 10, 64)
		ids = append(ids, id)
	}
	if err := h.svc.DeleteList(c, ids); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取设备
func (h *DeviceHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	device, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if device == nil {
		response.WriteSuccess(c, nil)
		return
	}

	var groupIDs []int64
	_ = json.Unmarshal(device.GroupIDs, &groupIDs)
	resp := &iot2.IotDeviceRespVO{
		ID:           device.ID,
		DeviceName:   device.DeviceName,
		Nickname:     device.Nickname,
		SerialNumber: device.SerialNumber,
		PicURL:       device.PicURL,
		GroupIDs:     groupIDs,
		ProductID:    device.ProductID,
		ProductKey:   device.ProductKey,
		DeviceType:   device.DeviceType,
		GatewayID:    device.GatewayID,
		State:        device.State,
		OnlineTime:   device.OnlineTime,
		OfflineTime:  device.OfflineTime,
		ActiveTime:   device.ActiveTime,
		DeviceSecret: device.DeviceSecret,
		Config:       string(device.Config),
		LocationType: device.LocationType,
		Latitude:     device.Latitude,
		Longitude:    device.Longitude,
		CreateTime:   device.CreateTime,
	}
	response.WriteSuccess(c, resp)
}

// GetAuthInfo 获取设备认证信息
func (h *DeviceHandler) GetAuthInfo(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	authInfo, err := h.svc.GetAuthInfo(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, authInfo)
}

// GetListByProductKeyAndNames 根据产品Key和设备名称列表获取设备
func (h *DeviceHandler) GetListByProductKeyAndNames(c *gin.Context) {
	var r iot2.IotDeviceByProductKeyAndNamesReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	list, err := h.svc.GetListByProductKeyAndNames(c, r.ProductKey, r.DeviceNames)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

// Page 获取设备分页
func (h *DeviceHandler) Page(c *gin.Context) {
	var r iot2.IotDevicePageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotDeviceRespVO, 0, len(page.List))
	for _, item := range page.List {
		var groupIDs []int64
		_ = json.Unmarshal(item.GroupIDs, &groupIDs)
		list = append(list, &iot2.IotDeviceRespVO{
			ID:           item.ID,
			DeviceName:   item.DeviceName,
			Nickname:     item.Nickname,
			SerialNumber: item.SerialNumber,
			PicURL:       item.PicURL,
			GroupIDs:     groupIDs,
			ProductID:    item.ProductID,
			ProductKey:   item.ProductKey,
			DeviceType:   item.DeviceType,
			GatewayID:    item.GatewayID,
			State:        item.State,
			OnlineTime:   item.OnlineTime,
			OfflineTime:  item.OfflineTime,
			ActiveTime:   item.ActiveTime,
			DeviceSecret: item.DeviceSecret,
			Config:       string(item.Config),
			LocationType: item.LocationType,
			Latitude:     item.Latitude,
			Longitude:    item.Longitude,
			CreateTime:   item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// UpdateGroup 更新设备分组
func (h *DeviceHandler) UpdateGroup(c *gin.Context) {
	var r iot2.IotDeviceUpdateGroupReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	if err := h.svc.UpdateGroup(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetCount 获取设备数量
func (h *DeviceHandler) GetCount(c *gin.Context) {
	productIDStr := c.Query("productId")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	count, err := h.svc.GetCountByProductID(c, productID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, count)
}
