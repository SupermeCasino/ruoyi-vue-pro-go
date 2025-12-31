package system

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	svc *system.MailService
}

func NewMailHandler(svc *system.MailService) *MailHandler {
	return &MailHandler{svc: svc}
}

// ================= Mail Account Request Handlers =================

func (h *MailHandler) CreateMailAccount(c *gin.Context) {
	var r req.MailAccountCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateMailAccount(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *MailHandler) UpdateMailAccount(c *gin.Context) {
	var r req.MailAccountUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateMailAccount(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) DeleteMailAccount(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteMailAccount(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) GetMailAccount(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	account, err := h.svc.GetMailAccount(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, account)
}

func (h *MailHandler) GetMailAccountPage(c *gin.Context) {
	var r req.MailAccountPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetMailAccountPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

func (h *MailHandler) GetSimpleMailAccountList(c *gin.Context) {
	list, err := h.svc.GetSimpleMailAccountList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

// ================= Mail Template Request Handlers =================

func (h *MailHandler) CreateMailTemplate(c *gin.Context) {
	var r req.MailTemplateCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateMailTemplate(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *MailHandler) UpdateMailTemplate(c *gin.Context) {
	var r req.MailTemplateUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateMailTemplate(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) DeleteMailTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteMailTemplate(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) GetMailTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	template, err := h.svc.GetMailTemplate(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, template)
}

func (h *MailHandler) GetMailTemplatePage(c *gin.Context) {
	var r req.MailTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetMailTemplatePage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

func (h *MailHandler) SendMail(c *gin.Context) {
	var r req.MailTemplateSendReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// TODO: Get Current User ID from context
	// For now, assuming user ID 0 or passed (but req doesn't have it).
	// In admin api, usually we might send to self or test?
	// The API `/system/mail/template/send-mail` usually takes `toMail` and `templateCode`?
	// RuoYi: "send-mail" testing API often requires `mail` and `templateCode`.
	// My `MailTemplateSendReq` has `toMail`.

	// Assuming logic SendMail(ctx, userId, userType, ...)
	// context.Get("userId")
	userId := int64(0)
	// For test purporse mostly.

	id, err := h.svc.SendMail(c, userId, 1, r.ToMail, r.TemplateCode, r.TemplateParams)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// ================= Mail Log Request Handlers =================

func (h *MailHandler) GetMailLogPage(c *gin.Context) {
	var r req.MailLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetMailLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}
