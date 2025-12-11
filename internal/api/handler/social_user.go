package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SocialUserHandler struct {
	socialUserService *service.SocialUserService
	logger            *zap.Logger
}

func NewSocialUserHandler(socialUserService *service.SocialUserService, logger *zap.Logger) *SocialUserHandler {
	return &SocialUserHandler{
		socialUserService: socialUserService,
		logger:            logger,
	}
}

// BindSocialUser 绑定社交用户
func (h *SocialUserHandler) BindSocialUser(c *gin.Context) {
	var req req.SocialUserBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	// 从上下文获取当前用户信息 (假设已通过认证中间件)
	// 从上下文获取当前用户信息
	userID := core.GetLoginUserID(c)
	userType := 2 // 2=System/Admin

	if err := h.socialUserService.BindSocialUser(c.Request.Context(), userID, userType, &req); err != nil {
		h.logger.Error("绑定社交用户失败", zap.Error(err))
		core.WriteError(c, 500, err.Error())
		return
	}

	core.WriteSuccess(c, true)
}

// UnbindSocialUser 解绑社交用户
func (h *SocialUserHandler) UnbindSocialUser(c *gin.Context) {
	var req req.SocialUserUnbindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	// 从上下文获取当前用户信息
	// 从上下文获取当前用户信息
	userID := core.GetLoginUserID(c)
	userType := 2 // 2=System/Admin

	if err := h.socialUserService.UnbindSocialUser(c.Request.Context(), userID, userType, req.Type, req.Openid); err != nil {
		h.logger.Error("解绑社交用户失败", zap.Error(err))
		core.WriteError(c, 500, err.Error())
		return
	}

	core.WriteSuccess(c, true)
}

// GetSocialUserList 获取用户绑定的社交账号列表
func (h *SocialUserHandler) GetSocialUserList(c *gin.Context) {
	// 从上下文获取当前用户信息
	// 从上下文获取当前用户信息
	userID := core.GetLoginUserID(c)
	userType := 2 // 2=System/Admin

	list, err := h.socialUserService.GetSocialUserList(c.Request.Context(), userID, userType)
	if err != nil {
		h.logger.Error("获取社交用户列表失败", zap.Error(err))
		core.WriteError(c, 500, err.Error())
		return
	}

	result := make([]resp.SocialUserResp, len(list))
	for i, user := range list {
		result[i] = resp.SocialUserResp{
			ID:          user.ID,
			Type:        user.Type,
			Openid:      user.Openid,
			Token:       user.Token,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			RawUserInfo: user.RawUserInfo,
			Code:        user.Code,
			State:       user.State,
			CreateTime:  user.CreatedAt,
		}
	}
	core.WriteSuccess(c, result)
}

// GetSocialUser 获取社交用户
func (h *SocialUserHandler) GetSocialUser(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))

	user, err := h.socialUserService.GetSocialUser(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取社交用户失败", zap.Error(err))
		core.WriteError(c, 500, err.Error())
		return
	}

	result := resp.SocialUserResp{
		ID:          user.ID,
		Type:        user.Type,
		Openid:      user.Openid,
		Token:       user.Token,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		RawUserInfo: user.RawUserInfo,
		Code:        user.Code,
		State:       user.State,
		CreateTime:  user.CreatedAt,
	}
	core.WriteSuccess(c, result)
}

// GetSocialUserPage 获取社交用户分页
func (h *SocialUserHandler) GetSocialUserPage(c *gin.Context) {
	var req req.SocialUserPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.socialUserService.GetSocialUserPage(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("获取社交用户分页失败", zap.Error(err))
		core.WriteError(c, 500, err.Error())
		return
	}

	list := make([]resp.SocialUserResp, len(page.List))
	for i, user := range page.List {
		list[i] = resp.SocialUserResp{
			ID:          user.ID,
			Type:        user.Type,
			Openid:      user.Openid,
			Token:       user.Token,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			RawUserInfo: user.RawUserInfo,
			Code:        user.Code,
			State:       user.State,
			CreateTime:  user.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.SocialUserResp]{
		List:  list,
		Total: page.Total,
	})
}
