package middleware

import (
	"backend-go/internal/pkg/core"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidatorMiddleware 参数验证中间件
// 与 Java 的 @Valid 注解对齐
func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 验证器已在 Gin 中内置，通过 c.ShouldBindJSON 等方法自动验证
		// 此中间件用于统一错误处理
		c.Next()
	}
}

// ValidateStruct 验证结构体
func ValidateStruct(data interface{}) error {
	validate := validator.New()
	return validate.Struct(data)
}

// HandleValidationError 处理验证错误
func HandleValidationError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// 构建错误消息
		errMsg := "参数验证失败: "
		for _, fieldError := range validationErrors {
			errMsg += fieldError.Field() + " " + fieldError.Tag() + "; "
		}
		c.JSON(200, core.Error(core.ParamErrCode, errMsg))
		return
	}

	c.JSON(200, core.Error(core.ParamErrCode, "参数错误"))
}
