package member

import (
	"regexp"
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.Login(c, &r, c.ClientIP(), c.Request.UserAgent(), h.getTerminal(c))
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// SmsLogin 手机+验证码登录
// @Router /member/auth/sms-login [post]
func (h *AppAuthHandler) SmsLogin(c *gin.Context) {
	var r req.AppAuthSmsLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 额外的验证码格式校验
	if !regexp.MustCompile(`^\d+$`).MatchString(r.Code) {
		response.WriteBizError(c, errors.NewBizError(40001, "验证码必须为数字"))
		return
	}
	res, err := h.svc.SmsLogin(c, &r, c.ClientIP(), c.Request.UserAgent(), h.getTerminal(c))
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// SocialLogin 社交登录
// @Router /member/auth/social-login [post]
func (h *AppAuthHandler) SocialLogin(c *gin.Context) {
	var r req.AppAuthSocialLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.SocialLogin(c, &r, c.ClientIP(), c.Request.UserAgent(), h.getTerminal(c))
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// SendSmsCode 发送手机验证码
// @Router /member/auth/send-sms-code [post]
func (h *AppAuthHandler) SendSmsCode(c *gin.Context) {
	var r req.AppAuthSmsSendReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.SendSmsCode(c, &r, c.ClientIP()); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// ValidateSmsCode 校验手机验证码
// @Router /member/auth/validate-sms-code [post]
func (h *AppAuthHandler) ValidateSmsCode(c *gin.Context) {
	var r req.AppAuthSmsValidateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.ValidateSmsCode(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// RefreshToken 刷新令牌
// @Router /member/auth/refresh-token [post]
func (h *AppAuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refreshToken")
	if refreshToken == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.RefreshToken(c, refreshToken, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// Logout 退出登录
// @Router /member/auth/logout [post]
func (h *AppAuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if err := h.svc.Logout(c, token, c.ClientIP(), c.Request.UserAgent()); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// SocialAuthRedirect 社交授权跳转
// @Router /member/auth/social-auth-redirect [get]
func (h *AppAuthHandler) SocialAuthRedirect(c *gin.Context) {
	socialType := c.Query("type")
	redirectUri := c.Query("redirectUri")
	if socialType == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	typeInt, _ := strconv.Atoi(socialType)
	url, err := h.svc.GetSocialAuthorizeUrl(c, typeInt, redirectUri)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, url)
}

// WeixinMiniAppLogin 微信小程序登录
// @Router /member/auth/weixin-mini-app-login [post]
func (h *AppAuthHandler) WeixinMiniAppLogin(c *gin.Context) {
	var r req.AppAuthWeixinMiniAppLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.WeixinMiniAppLogin(c, &r, c.ClientIP(), c.Request.UserAgent(), h.getTerminal(c))
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// CreateWeixinMpJsapiSignature 创建微信 JS SDK 初始化所需的签名
// @Router /member/auth/create-weixin-jsapi-signature [post]
func (h *AppAuthHandler) CreateWeixinMpJsapiSignature(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.CreateWeixinMpJsapiSignature(c, url)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

func (h *AppAuthHandler) getTerminal(c *gin.Context) int32 {
	terminal := c.GetHeader("terminal")
	if terminal == "" {
		return 0 // UNKNOWN
	}
	t, _ := strconv.Atoi(terminal)
	return int32(t)
}
