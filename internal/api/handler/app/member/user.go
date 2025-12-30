package member

import (
	"regexp"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppMemberUserHandler struct {
	svc *member.MemberUserService
}

func NewAppMemberUserHandler(svc *member.MemberUserService) *AppMemberUserHandler {
	return &AppMemberUserHandler{svc: svc}
}

// GetUserInfo 获得基本信息
// @Router /member/user/get [get]
func (h *AppMemberUserHandler) GetUserInfo(c *gin.Context) {
	userId := c.GetInt64(context.CtxUserIDKey)

	res, err := h.svc.GetUserInfo(c, userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(res))
}

// UpdateUser 修改基本信息
// @Router /member/user/update [put]
func (h *AppMemberUserHandler) UpdateUser(c *gin.Context) {
	var r req.AppMemberUserUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	userId := c.GetInt64(context.CtxUserIDKey)
	if err := h.svc.UpdateUser(c, userId, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateUserMobile 修改用户手机
// @Router /member/user/update-mobile [put]
func (h *AppMemberUserHandler) UpdateUserMobile(c *gin.Context) {
	var r req.AppMemberUserUpdateMobileReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 额外的参数验证
	if !regexp.MustCompile(`^\d{4,6}$`).MatchString(r.Code) {
		response.WriteBizError(c, errors.NewBizError(40001, "手机验证码长度为 4-6 位"))
		return
	}
	if r.OldCode != "" && !regexp.MustCompile(`^\d{4,6}$`).MatchString(r.OldCode) {
		response.WriteBizError(c, errors.NewBizError(40001, "原手机验证码长度为 4-6 位"))
		return
	}

	userId := c.GetInt64(context.CtxUserIDKey)
	if err := h.svc.UpdateUserMobile(c, userId, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateUserPassword 修改用户密码
// @Router /member/user/update-password [put]
func (h *AppMemberUserHandler) UpdateUserPassword(c *gin.Context) {
	var r req.AppMemberUserUpdatePasswordReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 额外的参数验证
	if !regexp.MustCompile(`^\d{4,6}$`).MatchString(r.Code) {
		response.WriteBizError(c, errors.NewBizError(40001, "手机验证码长度为 4-6 位"))
		return
	}

	if err := h.svc.UpdateUserPassword(c, c.GetInt64(context.CtxUserIDKey), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// ResetUserPassword 重置用户密码 (忘记密码)
// @Router /member/user/reset-password [put]
func (h *AppMemberUserHandler) ResetUserPassword(c *gin.Context) {
	var r req.AppMemberUserResetPasswordReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.ResetUserPassword(c, &r); err != nil {
			response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateUserMobileByWeixin 微信小程序更新手机号
// @Router /member/user/update-mobile-by-weixin [put]
func (h *AppMemberUserHandler) UpdateUserMobileByWeixin(c *gin.Context) {
	var r struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	userId := c.GetInt64(context.CtxUserIDKey)
	if err := h.svc.UpdateUserMobileByWeixin(c, userId, r.Code); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}
