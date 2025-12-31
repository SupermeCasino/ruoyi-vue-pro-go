package handler

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configSvc *system.ConfigService
}

func NewConfigHandler(configSvc *system.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configSvc: configSvc,
	}
}

// GetConfigPage 获得参数配置分页
func (h *ConfigHandler) GetConfigPage(c *gin.Context) {
	var req req.ConfigPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.configSvc.GetConfigPage(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetConfig 获得参数配置
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.configSvc.GetConfig(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetConfigKey 根据参数键名查询参数值
func (h *ConfigHandler) GetConfigKey(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	config, err := h.configSvc.GetConfigByKey(c, key)
	if err != nil || config == nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if !config.Visible {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	response.WriteSuccess(c, config.Value)
}

// CreateConfig 创建参数配置
func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	var req req.ConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.configSvc.CreateConfig(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateConfig 更新参数配置
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var req req.ConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.configSvc.UpdateConfig(c, &req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteConfig 删除参数配置
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.configSvc.DeleteConfig(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}
