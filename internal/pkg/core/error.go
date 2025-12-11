package core

import "fmt"

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
