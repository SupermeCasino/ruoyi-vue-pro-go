package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login 登录接口
// @Router /system/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req req.AuthLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		c.Error(err) // 交给 ErrorHandler 中间件处理
		return
	}

	c.JSON(200, core.Success(resp))
}

// GetPermissionInfo 获取权限信息
// @Router /system/auth/get-permission-info [get]
func (h *AuthHandler) GetPermissionInfo(c *gin.Context) {
	resp, err := h.svc.GetPermissionInfo(c)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(resp))
}

// Logout 登出
// @Router /system/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从 Header 获取 token
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	err := h.svc.Logout(c.Request.Context(), token)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// RefreshToken 刷新令牌
// @Router /system/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refreshToken")
	if refreshToken == "" {
		c.JSON(200, core.ErrParam)
		return
	}

	resp, err := h.svc.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(resp))
}

// SmsLogin 短信登录
// @Router /system/auth/sms-login [post]
func (h *AuthHandler) SmsLogin(c *gin.Context) {
	var r req.AuthSmsLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	resp, err := h.svc.SmsLogin(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(resp))
}

// SendSmsCode 发送短信验证码
// @Router /system/auth/send-sms-code [post]
func (h *AuthHandler) SendSmsCode(c *gin.Context) {
	var r req.AuthSmsSendReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	err := h.svc.SendSmsCode(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// Register 注册
// @Router /system/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var r req.AuthRegisterReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(resp))
}

// ResetPassword 重置密码
// @Router /system/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var r req.AuthResetPasswordReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	err := h.svc.ResetPassword(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// SocialAuthRedirect 社交授权跳转
// @Router /system/auth/social-auth-redirect [get]
func (h *AuthHandler) SocialAuthRedirect(c *gin.Context) {
	socialType := c.Query("type")
	redirectUri := c.Query("redirectUri")

	if socialType == "" {
		c.JSON(200, core.ErrParam)
		return
	}

	// 转换类型
	typeInt, err := strconv.Atoi(socialType)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	url, err := h.svc.SocialAuthRedirect(c.Request.Context(), typeInt, redirectUri)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(url))
}

// SocialLogin 社交登录
// @Router /system/auth/social-login [post]
func (h *AuthHandler) SocialLogin(c *gin.Context) {
	var r req.AuthSocialLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	resp, err := h.svc.SocialLogin(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(resp))
}
