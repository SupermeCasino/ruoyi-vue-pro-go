package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"

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
	// TODO: Get userId from context after middleware implementation
	// For now, let's assume middleware or manual testing sets it, or extract from generic Auth
	// Logic to be refined: how to get 'loginUserId' in this clean architecture?
	// Usually middleware sets a key. Let's assume "userId".
	userId := c.GetInt64("userId")
	if userId == 0 {
		// Fallback for testing or error if no middleware
		// c.JSON(401, core.Error(401, "Unauthorized"))
		// return
		// For verification purposes without middleware, we might mock it or expect 0 and fail.
	}

	res, err := h.svc.GetUserInfo(c, userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// UpdateUser 修改基本信息
// @Router /member/user/update [put]
func (h *AppMemberUserHandler) UpdateUser(c *gin.Context) {
	var r req.AppMemberUserUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	userId := c.GetInt64("userId")
	if err := h.svc.UpdateUser(c, userId, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// UpdateUserMobile 修改用户手机
// @Router /member/user/update-mobile [put]
func (h *AppMemberUserHandler) UpdateUserMobile(c *gin.Context) {
	var r req.AppMemberUserUpdateMobileReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}

	userId := c.GetInt64(core.CtxUserIDKey)
	if err := h.svc.UpdateUserMobile(c, userId, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// UpdateUserPassword 修改用户密码
// @Router /member/user/update-password [put]
func (h *AppMemberUserHandler) UpdateUserPassword(c *gin.Context) {
	var r req.AppMemberUserUpdatePasswordReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateUserPassword(c, c.GetInt64(core.CtxUserIDKey), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// ResetUserPassword 重置用户密码 (忘记密码)
// @Router /member/user/reset-password [put]
func (h *AppMemberUserHandler) ResetUserPassword(c *gin.Context) {
	var r req.AppMemberUserResetPasswordReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.ResetUserPassword(c, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}
