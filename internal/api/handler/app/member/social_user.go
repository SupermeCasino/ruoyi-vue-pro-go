package member

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppSocialUserHandler struct {
	svc *service.SocialUserService
}

func NewAppSocialUserHandler(svc *service.SocialUserService) *AppSocialUserHandler {
	return &AppSocialUserHandler{svc: svc}
}

// Bind 绑定社交用户
// @Router /member/social-user/bind [post]
func (h *AppSocialUserHandler) Bind(c *gin.Context) {
	var r req.AppSocialUserBindReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, err)
		return
	}

	userID := c.GetInt64("userId")
	openid, err := h.svc.BindSocialUser(c, userID, 1, &req.SocialUserBindReq{
		Type:  r.Type,
		Code:  r.Code,
		State: r.State,
	})
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, openid)
}

// Unbind 解绑社交用户
// @Router /member/social-user/unbind [delete]
func (h *AppSocialUserHandler) Unbind(c *gin.Context) {
	var r req.AppSocialUserUnbindReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, err)
		return
	}

	userID := c.GetInt64("userId")
	err := h.svc.UnbindSocialUser(c, userID, 1, r.Type, r.OpenID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, true)
}

// Get 获取社交用户
// @Router /member/social-user/get [get]
func (h *AppSocialUserHandler) Get(c *gin.Context) {
	socialTypeStr := c.Query("type")
	socialType, _ := strconv.Atoi(socialTypeStr)
	userID := c.GetInt64("userId")

	socialUsers, err := h.svc.GetSocialUserList(c, userID, 1)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	var target *resp.AppSocialUserResp
	for _, su := range socialUsers {
		if su.Type == socialType {
			target = &resp.AppSocialUserResp{
				Openid:   su.Openid,
				Nickname: su.Nickname,
				Avatar:   su.Avatar,
			}
			break
		}
	}

	response.WriteSuccess(c, target)
}

// GetWxaQrcode 获得微信小程序码
// @Router /member/social-user/wxa-qrcode [post]
func (h *AppSocialUserHandler) GetWxaQrcode(c *gin.Context) {
	path := c.Query("path")
	width, _ := strconv.Atoi(c.DefaultQuery("width", "430"))

	qrCode, err := h.svc.GetWxaQrcode(c, 1, path, width)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, qrCode)
}

// GetSubscribeTemplateList 获得微信小程序订阅模板列表
// @Router /member/social-user/get-subscribe-template-list [get]
func (h *AppSocialUserHandler) GetSubscribeTemplateList(c *gin.Context) {
	templates, err := h.svc.GetSubscribeTemplateList(c, 1)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, templates)
}
