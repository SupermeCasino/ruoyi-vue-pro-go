package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type AppTenantHandler struct {
	svc *system.TenantService
}

func NewAppTenantHandler(svc *system.TenantService) *AppTenantHandler {
	return &AppTenantHandler{svc: svc}
}

// GetTenantByWebsite 使用域名获得租户信息
// @Router /app-api/system/tenant/get-by-website [get]
func (h *AppTenantHandler) GetTenantByWebsite(c *gin.Context) {
	website := c.Query("website")
	if website == "" {
		response.WriteError(c, 400, "域名不能为空")
		return
	}

	res, err := h.svc.GetTenantByWebsite(c, website)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 对齐 Java 逻辑：查不到或禁用返回 success(null)
	// Service 层 GetTenantByWebsite 已处理 RecordNotFound 返回 nil, nil
	response.WriteSuccess(c, res)
}
