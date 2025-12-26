package model

// LoginLogTypeEnum 登录日志类型
const (
	LoginLogTypeUsername = 100 // 使用账号登录
	LoginLogTypeSocial   = 101 // 使用社交登录
	LoginLogTypeMobile   = 103 // 使用手机登录
	LoginLogTypeSms      = 104 // 使用短信登录

	LogoutLogTypeSelf   = 200 // 主动登出
	LogoutLogTypeDelete = 202 // 强制退出
)

// LoginResultEnum 登录结果
const (
	LoginResultSuccess          = 0  // 成功
	LoginResultBadCredentials   = 10 // 账号或密码不正确
	LoginResultUserDisabled     = 20 // 用户被禁用
	LoginResultCaptchaNotFound  = 30 // 验证码不存在
	LoginResultCaptchaCodeError = 31 // 验证码不正确
)
