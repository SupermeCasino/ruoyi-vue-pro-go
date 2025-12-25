package middleware

import (
	"net/http"

	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/logger"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProductErrorHandler 商品模块专用错误处理中间件
// 统一处理商品模块的各种错误，确保错误响应格式与Java版本一致
func ProductErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 1. 处理商品模块的业务异常
			if bizErr, ok := err.(*errors.BizError); ok {
				// 记录商品模块的业务错误日志
				if isProductError(bizErr.Code) {
					logger.Info("商品模块业务异常",
						zap.Int("code", bizErr.Code),
						zap.String("message", bizErr.Msg),
						zap.String("path", c.Request.URL.Path),
						zap.String("method", c.Request.Method),
					)
				}

				// 返回与Java版本一致的错误响应格式
				c.JSON(http.StatusOK, response.Result[any]{
					Code: bizErr.Code,
					Msg:  bizErr.Msg,
					Data: nil,
				})
				return
			}

			// 2. 处理其他未知错误
			logger.Error("商品模块系统异常",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			c.JSON(http.StatusOK, response.Error(errors.ServerErrCode, "系统内部异常"))
		}
	}
}

// isProductError 判断是否为商品模块的错误码
// 商品模块错误码范围：1-008-000-000 到 1-008-999-999
func isProductError(code int) bool {
	return code >= 1008000000 && code <= 1008999999
}

// HandleProductError 商品模块错误处理辅助函数
// 用于在Handler中统一处理商品模块的错误
func HandleProductError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 1. 处理商品模块的业务异常
	if bizErr, ok := err.(*errors.BizError); ok {
		// 记录商品模块的业务错误日志
		if isProductError(bizErr.Code) {
			logger.Info("商品模块业务异常",
				zap.Int("code", bizErr.Code),
				zap.String("message", bizErr.Msg),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
		}

		// 返回与Java版本一致的错误响应格式
		response.WriteError(c, bizErr.Code, bizErr.Msg)
		return
	}

	// 2. 处理其他未知错误
	logger.Error("商品模块系统异常",
		zap.Error(err),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)
	response.WriteError(c, errors.ServerErrCode, "系统内部异常")
}

// ProductErrorMapping 商品模块错误映射
// 将常见的系统错误映射为商品模块的业务错误
func ProductErrorMapping(err error) error {
	if err == nil {
		return nil
	}

	// 根据错误类型映射为商品模块的业务错误
	switch err.Error() {
	case "record not found":
		// 根据上下文判断是SPU还是SKU不存在
		// 这里默认返回SPU不存在，具体的Handler可以覆盖
		return product.ErrSpuNotExists
	case "invalid input":
		return errors.ErrParam
	default:
		return err
	}
}

// ValidateProductParams 商品模块参数验证辅助函数
// 统一处理商品模块的参数验证错误
func ValidateProductParams(c *gin.Context, err error) bool {
	if err != nil {
		logger.Info("商品模块参数验证失败",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		response.WriteError(c, errors.ParamErrCode, "参数错误: "+err.Error())
		return false
	}
	return true
}

// ProductNotFoundError 商品不存在错误处理
func ProductNotFoundError(c *gin.Context, productType string, id int64) {
	var err *errors.BizError
	switch productType {
	case "spu":
		err = product.ErrSpuNotExists
	case "sku":
		err = product.ErrSkuNotExists
	case "category":
		err = product.ErrCategoryNotExists
	case "brand":
		err = product.ErrBrandNotExists
	default:
		err = product.ErrSpuNotExists
	}

	logger.Info("商品不存在",
		zap.String("type", productType),
		zap.Int64("id", id),
		zap.String("path", c.Request.URL.Path),
	)
	response.WriteError(c, err.Code, err.Msg)
}
