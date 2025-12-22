package member

import "github.com/wxlbd/ruoyi-mall-go/pkg/errors"

// 会员认证模块错误码常量
// 参考 Java: yudao-module-member-api/src/main/java/cn/iocoder/yudao/module/member/enums/ErrorCodeConstants.java

var (
	// ========== AUTH 模块 ==========
	ErrAuthLoginBadCredentials = errors.NewBizError(1004003002, "账号或密码不正确")
	ErrAuthLoginUserDisabled   = errors.NewBizError(1004003001, "用户已被禁用")
	ErrAuthSocialUserNotFound  = errors.NewBizError(1002002000, "社交账号不存在")
	ErrAuthUserNotTokenValid   = errors.NewBizError(401, "Token无效或已过期")
	ErrAuthUserNotFound        = errors.NewBizError(1004003005, "用户不存在")
	ErrMobileFormatInvalid     = errors.NewBizError(1004003006, "手机号格式不正确")
	ErrSmsCodeFormatInvalid    = errors.NewBizError(1004003007, "验证码必须为数字")
)
