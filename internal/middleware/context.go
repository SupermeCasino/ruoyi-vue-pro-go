package middleware

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

// InjectContext 中间件：将 gin.Context 注入到 context.Context 中
// 用于在 service 层和 GORM Hook 中访问 gin.Context
func InjectContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将 gin.Context 放入 request context
		ctx := context.WithValue(c.Request.Context(), core.CtxGinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
