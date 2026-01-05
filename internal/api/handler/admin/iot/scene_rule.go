package iot

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建场景规则
func (h *SceneRuleHandler) Create(c *gin.Context) {
	var r iot2.IotSceneRuleSaveReqVO
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

// Update 更新场景规则
func (h *SceneRuleHandler) Update(c *gin.Context) {
	var r iot2.IotSceneRuleSaveReqVO
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

// UpdateStatus 更新场景规则状态
func (h *SceneRuleHandler) UpdateStatus(c *gin.Context) {
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

// Delete 删除场景规则
func (h *SceneRuleHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取场景规则
func (h *SceneRuleHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	rule, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if rule == nil {
		response.WriteSuccess(c, nil)
		return
	}

	var triggers []iot2.IotSceneRuleTrigger
	var actions []iot2.IotSceneRuleAction
	_ = json.Unmarshal(rule.Triggers, &triggers)
	_ = json.Unmarshal(rule.Actions, &actions)

	resp := &iot2.IotSceneRuleRespVO{
		ID:          rule.ID,
		Name:        rule.Name,
		Description: rule.Description,
		Status:      rule.Status,
		Triggers:    triggers,
		Actions:     actions,
		CreateTime:  rule.CreateTime,
	}
	response.WriteSuccess(c, resp)
}

// Page 获取场景规则分页
func (h *SceneRuleHandler) Page(c *gin.Context) {
	var r iot2.IotSceneRulePageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotSceneRuleRespVO, 0, len(page.List))
	for _, item := range page.List {
		var triggers []iot2.IotSceneRuleTrigger
		var actions []iot2.IotSceneRuleAction
		_ = json.Unmarshal(item.Triggers, &triggers)
		_ = json.Unmarshal(item.Actions, &actions)

		list = append(list, &iot2.IotSceneRuleRespVO{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			Triggers:    triggers,
			Actions:     actions,
			CreateTime:  item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

// SimpleList 获取场景规则精简列表
func (h *SceneRuleHandler) SimpleList(c *gin.Context) {
	list, err := h.svc.GetListByStatus(c, 0) // 0: 启用
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	resp := make([]*iot2.IotSceneRuleRespVO, 0, len(list))
	for _, item := range list {
		resp = append(resp, &iot2.IotSceneRuleRespVO{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	response.WriteSuccess(c, resp)
}
