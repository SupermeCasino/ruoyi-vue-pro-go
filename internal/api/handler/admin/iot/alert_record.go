package iot

import (
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Page 获取告警记录分页
func (h *AlertRecordHandler) Page(c *gin.Context) {
	var r iot2.IotAlertRecordPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotAlertRecordRespVO, 0, len(page.List))
	for _, item := range page.List {
		list = append(list, &iot2.IotAlertRecordRespVO{
			ID:            item.ID,
			ConfigID:      item.ConfigID,
			ConfigName:    item.ConfigName,
			ConfigLevel:   item.ConfigLevel,
			SceneRuleID:   item.SceneRuleID,
			ProductID:     item.ProductID,
			DeviceID:      item.DeviceID,
			DeviceMessage: item.DeviceMessage,
			ProcessStatus: item.ProcessStatus,
			ProcessRemark: item.ProcessRemark,
			CreateTime:    item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// Process 处理告警记录
func (h *AlertRecordHandler) Process(c *gin.Context) {
	var r iot2.IotAlertRecordProcessReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	if err := h.svc.Process(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取告警记录
func (h *AlertRecordHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	record, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if record == nil {
		response.WriteSuccess(c, nil)
		return
	}
	resp := &iot2.IotAlertRecordRespVO{
		ID:            record.ID,
		ConfigID:      record.ConfigID,
		ConfigName:    record.ConfigName,
		ConfigLevel:   record.ConfigLevel,
		SceneRuleID:   record.SceneRuleID,
		ProductID:     record.ProductID,
		DeviceID:      record.DeviceID,
		DeviceMessage: record.DeviceMessage,
		ProcessStatus: record.ProcessStatus,
		ProcessRemark: record.ProcessRemark,
		CreateTime:    record.CreateTime,
	}
	response.WriteSuccess(c, resp)
}
