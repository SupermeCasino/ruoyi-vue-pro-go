package handler

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configSvc *service.ConfigService
}

func NewConfigHandler(configSvc *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configSvc: configSvc,
	}
}

// GetConfigPage 获得参数配置分页
func (h *ConfigHandler) GetConfigPage(c *gin.Context) {
	var req req.ConfigPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	res, err := h.configSvc.GetConfigPage(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

// GetConfig 获得参数配置
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	res, err := h.configSvc.GetConfig(c, id)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

// GetConfigKey 根据参数键名查询参数值
func (h *ConfigHandler) GetConfigKey(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(400, response.Error(400, "key is required"))
		return
	}
	config, err := h.configSvc.GetConfigByKey(c, key)
	if err != nil || config == nil {
		c.JSON(500, response.Error(500, "config not found"))
		return
	}
	if !config.Visible {
		c.JSON(500, response.Error(500, "不可见的配置，不允许返回给前端"))
		return
	}
	c.JSON(200, response.Success(config.Value))
}

// CreateConfig 创建参数配置
func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	var req req.ConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	id, err := h.configSvc.CreateConfig(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(id))
}

// UpdateConfig 更新参数配置
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var req req.ConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	if err := h.configSvc.UpdateConfig(c, &req); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

// DeleteConfig 删除参数配置
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	if err := h.configSvc.DeleteConfig(c, id); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}
