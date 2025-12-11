package req

// AuthLoginReq 登录请求
type AuthLoginReq struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	TenantName string `json:"tenantName"` // 租户名, 某些版本前端可能传 tenantName
	// CaptchaVerificationReqVO fields (Skipping strict validation for now)
	CaptchaVerification string `json:"captchaVerification"`
}

// AuthSmsLoginReq 短信登录请求
type AuthSmsLoginReq struct {
	Mobile string `json:"mobile" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

// AuthSmsSendReq 发送短信验证码请求
type AuthSmsSendReq struct {
	Mobile string `json:"mobile" binding:"required"`
	Scene  int    `json:"scene" binding:"required"` // 场景：1-登录 2-注册 3-重置密码
}

// AuthRegisterReq 注册请求
type AuthRegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResetPasswordReq 重置密码请求
type AuthResetPasswordReq struct {
	Mobile   string `json:"mobile" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthSocialLoginReq 社交登录请求
type AuthSocialLoginReq struct {
	Type        int    `json:"type" binding:"required"`
	Code        string `json:"code" binding:"required"`
	State       string `json:"state"`
	RedirectUri string `json:"redirectUri"`
}
