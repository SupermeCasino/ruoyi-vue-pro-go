package middleware

import (
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"backend-go/internal/pkg/core"
	"backend-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 全局异常捕获中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取调用栈
				stack := string(debug.Stack())

				// 记录日志
				logger.Error("panic recover",
					zap.Any("error", err),
					zap.String("stack", stack),
				)

				// 检查连接断开
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				// 返回 500
				c.JSON(http.StatusInternalServerError, core.Error(core.ServerErrCode, "系统异常，请联系管理员"))
				c.Abort()
			}
		}()
		c.Next()
	}
}

// BizErrorHandle 业务错误处理中间件 (可选，如果想把 controller 的 return error 统一处理)
// 但 Gin 的 handler 签名没有 error 返回值。这里我们通常采用 c.Error() 机制或者封装 HandlerFunc
// 简单起见，我们推荐 Controller 显式调用 core.Error 或 core.Success
