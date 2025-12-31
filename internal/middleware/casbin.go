package middleware

import (
	"fmt"
	"net/http"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware Casbin 权限中间件
type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
	permSvc  *system.PermissionService
}

func NewCasbinMiddleware(enforcer *casbin.Enforcer, permSvc *system.PermissionService) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer: enforcer,
		permSvc:  permSvc,
	}
}

// RequirePermission 检查权限
// permission: 权限字符串，如 system:user:query
func (m *CasbinMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := context.GetLoginUser(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "未登录"))
			return
		}

		// 1. 超级管理员直接放行
		isSuper, err := m.permSvc.IsSuperAdmin(c.Request.Context(), user.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(500, "判断超级管理员失败"))
			return
		}
		if isSuper {
			c.Next()
			return
		}

		// 2. Casbin 鉴权
		// Subject: user:{userId}
		// Object: permission
		// Action: access
		// 注意：Adapter 加载的 g 策略是 g, user:{userId}, role:{roleId}
		// Adapter 加载的 p 策略是 p, role:{roleId}, permission, access
		// Casbin 会自动推导 user -> role -> permission
		sub := fmt.Sprintf("user:%d", user.UserID)
		ok, err := m.enforcer.Enforce(sub, permission, "access")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(500, "权限校验错误"))
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error(403, "权限不足"))
			return
		}

		c.Next()
	}
}
