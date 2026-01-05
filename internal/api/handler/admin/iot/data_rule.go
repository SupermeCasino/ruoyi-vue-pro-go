package iot

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建数据规则
func (h *DataRuleHandler) Create(c *gin.Context) {
	var r iot2.IotDataRuleSaveReqVO
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

// Update 更新数据规则
func (h *DataRuleHandler) Update(c *gin.Context) {
	var r iot2.IotDataRuleSaveReqVO
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

// Delete 删除数据规则
func (h *DataRuleHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取数据规则
func (h *DataRuleHandler) Get(c *gin.Context) {
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

	var sourceConfigs []iot2.IotDataRuleSourceConfig
	var sinkIDs []int64
	_ = json.Unmarshal(rule.SourceConfigs, &sourceConfigs)
	_ = json.Unmarshal(rule.SinkIDs, &sinkIDs)

	resp := &iot2.IotDataRuleRespVO{
		ID:            rule.ID,
		Name:          rule.Name,
		Description:   rule.Description,
		Status:        rule.Status,
		SourceConfigs: sourceConfigs,
		SinkIDs:       sinkIDs,
		CreateTime:    rule.CreateTime,
	}
	response.WriteSuccess(c, resp)
}

// Page 获取数据规则分页
func (h *DataRuleHandler) Page(c *gin.Context) {
	var r iot2.IotDataRulePageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotDataRuleRespVO, 0, len(page.List))
	for _, item := range page.List {
		var sourceConfigs []iot2.IotDataRuleSourceConfig
		var sinkIDs []int64
		_ = json.Unmarshal(item.SourceConfigs, &sourceConfigs)
		_ = json.Unmarshal(item.SinkIDs, &sinkIDs)

		list = append(list, &iot2.IotDataRuleRespVO{
			ID:            item.ID,
			Name:          item.Name,
			Description:   item.Description,
			Status:        item.Status,
			SourceConfigs: sourceConfigs,
			SinkIDs:       sinkIDs,
			CreateTime:    item.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}
