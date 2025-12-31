package system

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SocialUserHandler struct {
	socialUserService *system.SocialUserService
	logger            *zap.Logger
}

func NewSocialUserHandler(socialUserService *system.SocialUserService, logger *zap.Logger) *SocialUserHandler {
	return &SocialUserHandler{
		socialUserService: socialUserService,
		logger:            logger,
	}
}

// BindSocialUser 绑定社交用户
func (h *SocialUserHandler) BindSocialUser(c *gin.Context) {
	var req req.SocialUserBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 从上下文获取当前用户信息 (假设已通过认证中间件)
	// 从上下文获取当前用户信息
	userID := context.GetLoginUserID(c)
	userType := consts.UserTypeAdmin // 管理员用户类型

	if _, err := h.socialUserService.BindSocialUser(c.Request.Context(), userID, userType, &req); err != nil {
		h.logger.Error("绑定社交用户失败", zap.Error(err))
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, true)
}

// UnbindSocialUser 解绑社交用户
func (h *SocialUserHandler) UnbindSocialUser(c *gin.Context) {
	var req req.SocialUserUnbindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 从上下文获取当前用户信息
	// 从上下文获取当前用户信息
	userID := context.GetLoginUserID(c)
	userType := consts.UserTypeAdmin // 管理员用户类型

	if err := h.socialUserService.UnbindSocialUser(c.Request.Context(), userID, userType, req.Type, req.Openid); err != nil {
		h.logger.Error("解绑社交用户失败", zap.Error(err))
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, true)
}

// GetSocialUserList 获取用户绑定的社交账号列表
func (h *SocialUserHandler) GetSocialUserList(c *gin.Context) {
	// 从上下文获取当前用户信息
	// 从上下文获取当前用户信息
	userID := context.GetLoginUserID(c)
	userType := consts.UserTypeAdmin // 管理员用户类型

	list, err := h.socialUserService.GetSocialUserList(c.Request.Context(), userID, userType)
	if err != nil {
		h.logger.Error("获取社交用户列表失败", zap.Error(err))
		response.WriteBizError(c, err)
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
			CreateTime:  user.CreateTime,
		}
	}
	response.WriteSuccess(c, result)
}

// GetSocialUser 获取社交用户
func (h *SocialUserHandler) GetSocialUser(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))

	user, err := h.socialUserService.GetSocialUser(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取社交用户失败", zap.Error(err))
		response.WriteBizError(c, err)
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
		CreateTime:  user.CreateTime,
	}
	response.WriteSuccess(c, result)
}

// GetSocialUserPage 获取社交用户分页
func (h *SocialUserHandler) GetSocialUserPage(c *gin.Context) {
	var req req.SocialUserPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	page, err := h.socialUserService.GetSocialUserPage(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("获取社交用户分页失败", zap.Error(err))
		response.WriteBizError(c, err)
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
			CreateTime:  user.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.SocialUserResp]{
		List:  list,
		Total: page.Total,
	})
}
