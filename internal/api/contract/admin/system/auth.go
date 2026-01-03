package system

import "time"

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

// AuthLoginResp 登录响应
type AuthLoginResp struct {
	UserId       int64     `json:"userId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresTime  time.Time `json:"expiresTime"`
}

type AuthPermissionInfoResp struct {
	User        UserVO   `json:"user"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Menus       []MenuVO `json:"menus"`
}

type UserVO struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	DeptID   int64  `json:"deptId"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type MenuVO struct {
	ID            int64    `json:"id"`
	ParentID      int64    `json:"parentId"`
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Component     string   `json:"component"`
	ComponentName string   `json:"componentName"`
	Icon          string   `json:"icon"`
	Visible       bool     `json:"visible"`
	KeepAlive     bool     `json:"keepAlive"`
	AlwaysShow    bool     `json:"alwaysShow"`
	Children      []MenuVO `json:"children,omitempty"`
}
