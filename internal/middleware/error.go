package middleware

import (
	"net/http"

	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/logger"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 处理请求过程中的错误
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 1. 如果是业务异常
			if bizErr, ok := err.(*errors.BizError); ok {
				c.JSON(http.StatusOK, response.Result[any]{
					Code: bizErr.Code,
					Msg:  bizErr.Msg,
					Data: nil,
				})
				return
			}

			// 2. 如果是其他未知错误
			logger.Error("internal server error", zap.Error(err), zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusInternalServerError, response.Error(errors.ServerErrCode, "系统内部异常"))
		}
	}
}
