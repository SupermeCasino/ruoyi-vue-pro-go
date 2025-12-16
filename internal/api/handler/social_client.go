package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
	"go.uber.org/zap"
)

type SocialClientHandler struct {
	socialClientService *service.SocialClientService
	logger              *zap.Logger
}

func NewSocialClientHandler(socialClientService *service.SocialClientService, logger *zap.Logger) *SocialClientHandler {
	return &SocialClientHandler{
		socialClientService: socialClientService,
		logger:              logger,
	}
}

// CreateSocialClient 创建社交客户端
func (h *SocialClientHandler) CreateSocialClient(c *gin.Context) {
	var req req.SocialClientSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	id, err := h.socialClientService.CreateSocialClient(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("创建社交客户端失败", zap.Error(err))
		response.WriteError(c, 500, err.Error())
		return
	}

	response.WriteSuccess(c, id)
}

// UpdateSocialClient 更新社交客户端
func (h *SocialClientHandler) UpdateSocialClient(c *gin.Context) {
	var req req.SocialClientSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	if err := h.socialClientService.UpdateSocialClient(c.Request.Context(), &req); err != nil {
		h.logger.Error("更新社交客户端失败", zap.Error(err))
		response.WriteError(c, 500, err.Error())
		return
	}

	response.WriteSuccess(c, true)
}

// DeleteSocialClient 删除社交客户端
func (h *SocialClientHandler) DeleteSocialClient(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))

	if err := h.socialClientService.DeleteSocialClient(c.Request.Context(), id); err != nil {
		h.logger.Error("删除社交客户端失败", zap.Error(err))
		response.WriteError(c, 500, err.Error())
		return
	}

	response.WriteSuccess(c, true)
}

// GetSocialClient 获取社交客户端
func (h *SocialClientHandler) GetSocialClient(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))

	client, err := h.socialClientService.GetSocialClient(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取社交客户端失败", zap.Error(err))
		response.WriteError(c, 500, err.Error())
		return
	}

	result := resp.SocialClientResp{
		ID:           client.ID,
		Name:         client.Name,
		SocialType:   client.SocialType,
		UserType:     client.UserType,
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
		AgentId:      client.AgentId,
		Status:       client.Status,
		CreateTime:   client.CreatedAt,
	}
	response.WriteSuccess(c, result)
}

// GetSocialClientPage 获取社交客户端分页
func (h *SocialClientHandler) GetSocialClientPage(c *gin.Context) {
	var req req.SocialClientPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	page, err := h.socialClientService.GetSocialClientPage(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("获取社交客户端分页失败", zap.Error(err))
		response.WriteError(c, 500, err.Error())
		return
	}

	list := make([]resp.SocialClientResp, len(page.List))
	for i, client := range page.List {
		list[i] = resp.SocialClientResp{
			ID:           client.ID,
			Name:         client.Name,
			SocialType:   client.SocialType,
			UserType:     client.UserType,
			ClientId:     client.ClientId,
			ClientSecret: client.ClientSecret,
			AgentId:      client.AgentId,
			Status:       client.Status,
			CreateTime:   client.CreatedAt,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.SocialClientResp]{
		List:  list,
		Total: page.Total,
	})
}
