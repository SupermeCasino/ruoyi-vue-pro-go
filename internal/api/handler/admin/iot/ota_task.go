package iot

import (
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建OTA任务
func (h *OtaTaskHandler) Create(c *gin.Context) {
	var r iot2.IotOtaTaskCreateReqVO
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

// Get 获取OTA任务
func (h *OtaTaskHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	task, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if task == nil {
		response.WriteSuccess(c, nil)
		return
	}
	resp := &iot2.IotOtaTaskRespVO{
		ID:                 task.ID,
		Name:               task.Name,
		Description:        task.Description,
		FirmwareID:         task.FirmwareID,
		Status:             task.Status,
		DeviceScope:        task.DeviceScope,
		DeviceTotalCount:   task.DeviceTotalCount,
		DeviceSuccessCount: task.DeviceSuccessCount,
		CreateTime:         task.CreateTime,
	}
	response.WriteSuccess(c, resp)
}

// Cancel 取消OTA任务
func (h *OtaTaskHandler) Cancel(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Cancel(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Page 获取OTA任务分页
func (h *OtaTaskHandler) Page(c *gin.Context) {
	var r iot2.IotOtaTaskPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotOtaTaskRespVO, 0, len(page.List))
	for _, item := range page.List {
		list = append(list, &iot2.IotOtaTaskRespVO{
			ID:                 item.ID,
			Name:               item.Name,
			Description:        item.Description,
			FirmwareID:         item.FirmwareID,
			Status:             item.Status,
			DeviceScope:        item.DeviceScope,
			DeviceTotalCount:   item.DeviceTotalCount,
			DeviceSuccessCount: item.DeviceSuccessCount,
			CreateTime:         item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// RecordPage 获取OTA任务记录分页
func (h *OtaTaskHandler) RecordPage(c *gin.Context) {
	var r iot2.IotOtaTaskRecordPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetRecordPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotOtaTaskRecordRespVO, 0, len(page.List))
	for _, item := range page.List {
		list = append(list, &iot2.IotOtaTaskRecordRespVO{
			ID:             item.ID,
			FirmwareID:     item.FirmwareID,
			TaskID:         item.TaskID,
			DeviceID:       item.DeviceID,
			FromFirmwareID: item.FromFirmwareID,
			Status:         item.Status,
			Progress:       item.Progress,
			Description:    item.Description,
			CreateTime:     item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}
