package handler

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type NotifyHandler struct {
	svc *service.NotifyService
}

func NewNotifyHandler(svc *service.NotifyService) *NotifyHandler {
	return &NotifyHandler{svc: svc}
}

// ================= Template Handlers =================

func (h *NotifyHandler) CreateNotifyTemplate(c *gin.Context) {
	var r req.NotifyTemplateCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateNotifyTemplate(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *NotifyHandler) UpdateNotifyTemplate(c *gin.Context) {
	var r req.NotifyTemplateUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateNotifyTemplate(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *NotifyHandler) DeleteNotifyTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteNotifyTemplate(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *NotifyHandler) GetNotifyTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	t, err := h.svc.GetNotifyTemplate(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, t)
}

func (h *NotifyHandler) GetNotifyTemplatePage(c *gin.Context) {
	var r req.NotifyTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetNotifyTemplatePage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

func (h *NotifyHandler) SendNotify(c *gin.Context) {
	var r req.NotifyTemplateSendReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.SendNotify(c, r.UserID, r.UserType, r.TemplateCode, r.TemplateParams)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// ================= Message Handlers =================

func (h *NotifyHandler) GetNotifyMessagePage(c *gin.Context) {
	var r req.NotifyMessagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// Admin view
	page, err := h.svc.GetNotifyMessagePage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

func (h *NotifyHandler) GetMyNotifyMessagePage(c *gin.Context) {
	var r req.MyNotifyMessagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	loginUser := context.GetLoginUser(c)
	if loginUser == nil {
		response.WriteBizError(c, errors.NewBizError(401, "未登录"))
		return
	}

	page, err := h.svc.GetMyNotifyMessagePage(c, loginUser.UserID, loginUser.UserType, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

func (h *NotifyHandler) UpdateNotifyMessageRead(c *gin.Context) {
	var r req.NotifyMessageUpdateReadReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	loginUser := context.GetLoginUser(c)
	if loginUser == nil {
		response.WriteBizError(c, errors.NewBizError(401, "未登录"))
		return
	}
	if err := h.svc.UpdateNotifyMessageRead(c, loginUser.UserID, loginUser.UserType, r.IDs); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *NotifyHandler) UpdateAllNotifyMessageRead(c *gin.Context) {
	loginUser := context.GetLoginUser(c)
	if loginUser == nil {
		response.WriteBizError(c, errors.NewBizError(401, "未登录"))
		return
	}
	if err := h.svc.UpdateAllNotifyMessageRead(c, loginUser.UserID, loginUser.UserType); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *NotifyHandler) GetUnreadNotifyMessageCount(c *gin.Context) {
	loginUser := context.GetLoginUser(c)
	if loginUser == nil {
		response.WriteBizError(c, errors.NewBizError(401, "未登录"))
		return
	}
	count, err := h.svc.GetUnreadNotifyMessageCount(c, loginUser.UserID, loginUser.UserType)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, count)
}

// GetNotifyMessage 获取单条站内信 (对齐 Java: NotifyMessageController.getNotifyMessage)
func (h *NotifyHandler) GetNotifyMessage(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	msg, err := h.svc.GetNotifyMessage(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, msg)
}

// GetUnreadNotifyMessageList 获取未读站内信列表 (对齐 Java: NotifyMessageController.getUnreadNotifyMessageList)
func (h *NotifyHandler) GetUnreadNotifyMessageList(c *gin.Context) {
	loginUser := context.GetLoginUser(c)
	if loginUser == nil {
		response.WriteBizError(c, errors.NewBizError(401, "未登录"))
		return
	}
	sizeStr := c.DefaultQuery("size", "10")
	size, _ := strconv.Atoi(sizeStr)
	if size <= 0 {
		size = 10
	}
	list, err := h.svc.GetUnreadNotifyMessageList(c, loginUser.UserID, loginUser.UserType, size)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}
