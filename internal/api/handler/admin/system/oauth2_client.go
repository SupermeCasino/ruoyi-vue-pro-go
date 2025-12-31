package system

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type OAuth2ClientHandler struct {
	svc *system.OAuth2ClientService
}

func NewOAuth2ClientHandler(svc *system.OAuth2ClientService) *OAuth2ClientHandler {
	return &OAuth2ClientHandler{svc: svc}
}

func (h *OAuth2ClientHandler) CreateOAuth2Client(c *gin.Context) {
	var r req.OAuth2ClientCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateOAuth2Client(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *OAuth2ClientHandler) UpdateOAuth2Client(c *gin.Context) {
	var r req.OAuth2ClientUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateOAuth2Client(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *OAuth2ClientHandler) DeleteOAuth2Client(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteOAuth2Client(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *OAuth2ClientHandler) GetOAuth2Client(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	client, err := h.svc.GetOAuth2Client(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, client)
}

func (h *OAuth2ClientHandler) GetOAuth2ClientPage(c *gin.Context) {
	var r req.OAuth2ClientPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetOAuth2ClientPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}
