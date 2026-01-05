package iot

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建告警配置
func (h *AlertConfigHandler) Create(c *gin.Context) {
	var r iot2.IotAlertConfigSaveReqVO
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

// Update 更新告警配置
func (h *AlertConfigHandler) Update(c *gin.Context) {
	var r iot2.IotAlertConfigSaveReqVO
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

// Delete 删除告警配置
func (h *AlertConfigHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取告警配置
func (h *AlertConfigHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	config, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if config == nil {
		response.WriteSuccess(c, nil)
		return
	}

	var sceneRuleIDs []int64
	var receiveUserIDs []int64
	var receiveTypes []int
	_ = json.Unmarshal(config.SceneRuleIDs, &sceneRuleIDs)
	_ = json.Unmarshal(config.ReceiveUserIDs, &receiveUserIDs)
	_ = json.Unmarshal(config.ReceiveTypes, &receiveTypes)

	// 获取用户姓名
	var receiveUserNames []string
	if len(receiveUserIDs) > 0 {
		// TODO:
		// 这里简单调用或者依赖注入 AdminUserApi
		// 为了简化，目前先留空或者从全局获取
	}

	resp := &iot2.IotAlertConfigRespVO{
		ID:               config.ID,
		Name:             config.Name,
		Description:      config.Description,
		Level:            config.Level,
		Status:           config.Status,
		SceneRuleIDs:     sceneRuleIDs,
		ReceiveUserIDs:   receiveUserIDs,
		ReceiveUserNames: receiveUserNames,
		ReceiveTypes:     receiveTypes,
		CreateTime:       config.CreateTime,
	}
	response.WriteSuccess(c, resp)
}

// Page 获取告警配置分页
func (h *AlertConfigHandler) Page(c *gin.Context) {
	var r iot2.IotAlertConfigPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotAlertConfigRespVO, 0, len(page.List))
	for _, item := range page.List {
		var sceneRuleIDs []int64
		var receiveUserIDs []int64
		var receiveTypes []int
		_ = json.Unmarshal(item.SceneRuleIDs, &sceneRuleIDs)
		_ = json.Unmarshal(item.ReceiveUserIDs, &receiveUserIDs)
		_ = json.Unmarshal(item.ReceiveTypes, &receiveTypes)

		list = append(list, &iot2.IotAlertConfigRespVO{
			ID:             item.ID,
			Name:           item.Name,
			Description:    item.Description,
			Level:          item.Level,
			Status:         item.Status,
			SceneRuleIDs:   sceneRuleIDs,
			ReceiveUserIDs: receiveUserIDs,
			ReceiveTypes:   receiveTypes,
			CreateTime:     item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// SimpleList 获取告警配置简单列表
func (h *AlertConfigHandler) SimpleList(c *gin.Context) {
	list, err := h.svc.GetListByStatus(c, consts.CommonStatusEnable) // 1: 启用
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	resp := make([]*iot2.IotAlertConfigRespVO, 0, len(list))
	for _, item := range list {
		resp = append(resp, &iot2.IotAlertConfigRespVO{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	response.WriteSuccess(c, resp)
}
