package handler

import (
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	svc *service.TenantService
}

func NewTenantHandler(svc *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

// GetTenantSimpleList 获取租户精简列表
// @Router /system/tenant/simple-list [get]
func (h *TenantHandler) GetTenantSimpleList(c *gin.Context) {
	list, err := h.svc.GetTenantSimpleList(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}

// GetTenantByWebsite 根据域名获取租户
// @Router /system/tenant/get-by-website [get]
func (h *TenantHandler) GetTenantByWebsite(c *gin.Context) {
	website := c.Query("website")
	tenant, err := h.svc.GetTenantByWebsite(c.Request.Context(), website)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(tenant))
}

// GetTenantIdByName 根据租户名获取租户ID
// @Router /system/tenant/get-id-by-name [get]
func (h *TenantHandler) GetTenantIdByName(c *gin.Context) {
	name := c.Query("name")
	tenantId, err := h.svc.GetTenantIdByName(c.Request.Context(), name)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(tenantId))
}
