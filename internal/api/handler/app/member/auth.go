package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"

	"github.com/gin-gonic/gin"
)

type AppAuthHandler struct {
	svc *member.MemberAuthService
}

func NewAppAuthHandler(svc *member.MemberAuthService) *AppAuthHandler {
	return &AppAuthHandler{svc: svc}
}

// Login 手机+密码登录
// @Router /member/auth/login [post]
func (h *AppAuthHandler) Login(c *gin.Context) {
	var r req.AppAuthLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.Login(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// SmsLogin 手机+验证码登录
// @Router /member/auth/sms-login [post]
func (h *AppAuthHandler) SmsLogin(c *gin.Context) {
	var r req.AppAuthSmsLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.SmsLogin(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// SocialLogin 社交登录
// @Router /member/auth/social-login [post]
func (h *AppAuthHandler) SocialLogin(c *gin.Context) {
	var r req.AppAuthSocialLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.SocialLogin(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// SendSmsCode 发送手机验证码
// @Router /member/auth/send-sms-code [post]
func (h *AppAuthHandler) SendSmsCode(c *gin.Context) {
	var r req.AppAuthSmsSendReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.SendSmsCode(c, &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// ValidateSmsCode 校验手机验证码
// @Router /member/auth/validate-sms-code [post]
func (h *AppAuthHandler) ValidateSmsCode(c *gin.Context) {
	var r req.AppAuthSmsValidateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.ValidateSmsCode(c, &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// RefreshToken 刷新令牌
// @Router /member/auth/refresh-token [post]
func (h *AppAuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refreshToken")
	if refreshToken == "" {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.svc.RefreshToken(c, refreshToken)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// Logout 退出登录
// @Router /member/auth/logout [post]
func (h *AppAuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if err := h.svc.Logout(c, token); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}
