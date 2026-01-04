package errors

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Is = errors.Is

// Global Error Codes - 与 Java GlobalErrorCodeConstants 对齐
const (
	// 成功
	SuccessCode = 0

	// 客户端错误 (4xx)
	ParamErrCode     = 400 // 参数错误
	UnauthorizedCode = 401 // 未授权/未登录
	ForbiddenCode    = 403 // 禁止访问
	NotFoundCode     = 404 // 资源不存在
	ConflictCode     = 409 // 冲突

	// 服务器错误 (5xx)
	ServerErrCode      = 500 // 系统异常
	NotImplementCode   = 501 // 未实现
	ServiceUnavailCode = 503 // 服务不可用

	// 业务失败标识（来自 core/pay_consts.go）
	FailCode = 1
)

// BizError 业务异常
type BizError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func NewBizError(code int, msg string) *BizError {
	return &BizError{
		Code: code,
		Msg:  msg,
	}
}

// 常用错误快捷方式
var (
	ErrUnknown      = NewBizError(ServerErrCode, "系统异常")
	ErrParam        = NewBizError(ParamErrCode, "参数错误")
	ErrUnauthorized = NewBizError(UnauthorizedCode, "未登录")
	ErrForbidden    = NewBizError(ForbiddenCode, "禁止访问")
	ErrNotFound     = NewBizError(NotFoundCode, "资源不存在")
	ErrConflict     = NewBizError(ConflictCode, "资源冲突")
)

// ParseBindingError 解析 Gin binding 错误，返回有意义的中文错误提示
// 它将 validator.ValidationErrors 转换为用户友好的错误消息
func ParseBindingError(err error) *BizError {
	if err == nil {
		return nil
	}

	// 处理 validator.ValidationErrors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		if len(validationErrors) > 0 {
			// 获取第一个验证错误
			fieldError := validationErrors[0]
			return buildFieldErrorMessage(fieldError)
		}
	}

	// 默认返回参数错误
	return ErrParam
}

// buildFieldErrorMessage 根据字段验证错误构建错误消息
func buildFieldErrorMessage(fieldError validator.FieldError) *BizError {
	fieldName := getChineseFieldName(fieldError.Field())
	tag := fieldError.Tag()

	switch tag {
	case "required":
		return NewBizError(ParamErrCode, fmt.Sprintf("%s 不能为空", fieldName))
	case "min":
		return NewBizError(ParamErrCode, fmt.Sprintf("%s 不能少于 %s 个字符", fieldName, fieldError.Param()))
	case "max":
		return NewBizError(ParamErrCode, fmt.Sprintf("%s 不能超过 %s 个字符", fieldName, fieldError.Param()))
	case "email":
		return NewBizError(ParamErrCode, fmt.Sprintf("%s 邮箱格式错误", fieldName))
	case "numeric":
		return NewBizError(ParamErrCode, fmt.Sprintf("%s 必须为数字", fieldName))
	case "len":
		return NewBizError(ParamErrCode, fmt.Sprintf("%s 长度必须为 %s", fieldName, fieldError.Param()))
	default:
		return NewBizError(ParamErrCode, fmt.Sprintf("%s 验证失败", fieldName))
	}
}

// getChineseFieldName 获取字段的中文名称
// 支持 JSON tag 中的自定义名称，否则返回字段原名
func getChineseFieldName(fieldName string) string {
	// 这里可以根据需要添加字段映射表
	// 例如: "ID" -> "ID", "Status" -> "状态"
	fieldNameMap := map[string]string{
		"ID":       "ID",
		"Status":   "状态",
		"Name":     "名称",
		"Code":     "编码",
		"Username": "用户名",
		"Password": "密码",
		"Email":    "邮箱",
		"Mobile":   "手机号",
	}

	if chinese, ok := fieldNameMap[fieldName]; ok {
		return chinese
	}

	// 默认返回原字段名
	return fieldName
}

// BindingErr 是一个便捷函数，用于在 handler 中处理 binding 错误
// 用法: if err := c.ShouldBindJSON(&req); err != nil {
//         response.WriteBizError(c, BindingErr(err))
//         return
//       }
func BindingErr(err error) error {
	if err == nil {
		return nil
	}

	bizErr := ParseBindingError(err)
	if bizErr != nil {
		return bizErr
	}

	return ErrParam
}
