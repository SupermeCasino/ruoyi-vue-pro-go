package system

import (
	"encoding/json"
	"strconv"

	system2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
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
	var r system2.MailAccountSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
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
	var r system2.MailAccountSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	if err := h.svc.UpdateMailAccount(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) DeleteMailAccount(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
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

func (h *MailHandler) DeleteMailAccountList(c *gin.Context) {
	var ids []int64
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteMailAccountList(c, ids); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) GetMailAccount(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	account, err := h.svc.GetMailAccount(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, system2.MailAccountRespVO{
		ID:             account.ID,
		Mail:           account.Mail,
		Username:       account.Username,
		Password:       account.Password,
		Host:           account.Host,
		Port:           account.Port,
		SslEnable:      bool(account.SslEnable),
		StarttlsEnable: bool(account.StarttlsEnable),
		CreateTime:     account.CreateTime,
	})
}

func (h *MailHandler) GetMailAccountPage(c *gin.Context) {
	var r system2.MailAccountPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetMailAccountPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*system2.MailAccountRespVO, 0, len(page.List))
	for _, account := range page.List {
		list = append(list, &system2.MailAccountRespVO{
			ID:             account.ID,
			Mail:           account.Mail,
			Username:       account.Username,
			Password:       account.Password,
			Host:           account.Host,
			Port:           account.Port,
			SslEnable:      bool(account.SslEnable),
			StarttlsEnable: bool(account.StarttlsEnable),
			CreateTime:     account.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

func (h *MailHandler) GetSimpleMailAccountList(c *gin.Context) {
	list, err := h.svc.GetSimpleMailAccountList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	respList := make([]system2.MailAccountSimpleRespVO, 0, len(list))
	for _, item := range list {
		respList = append(respList, system2.MailAccountSimpleRespVO{
			ID:   item.ID,
			Mail: item.Mail,
		})
	}
	response.WriteSuccess(c, respList)
}

// ================= Mail Template Request Handlers =================

func (h *MailHandler) CreateMailTemplate(c *gin.Context) {
	var r system2.MailTemplateSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
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
	var r system2.MailTemplateSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	if err := h.svc.UpdateMailTemplate(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) DeleteMailTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
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

func (h *MailHandler) DeleteMailTemplateList(c *gin.Context) {
	var ids []int64
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteMailTemplateList(c, ids); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *MailHandler) GetMailTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	template, err := h.svc.GetMailTemplate(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, system2.MailTemplateRespVO{
		ID:         template.ID,
		Name:       template.Name,
		Code:       template.Code,
		AccountID:  template.AccountID,
		Nickname:   template.Nickname,
		Title:      template.Title,
		Content:    template.Content,
		Params:     template.Params,
		Status:     template.Status,
		Remark:     template.Remark,
		CreateTime: template.CreateTime,
	})
}

func (h *MailHandler) GetSimpleMailTemplateList(c *gin.Context) {
	list, err := h.svc.GetSimpleMailTemplateList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	respList := make([]system2.MailTemplateSimpleRespVO, 0, len(list))
	for _, item := range list {
		respList = append(respList, system2.MailTemplateSimpleRespVO{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	response.WriteSuccess(c, respList)
}

func (h *MailHandler) GetMailTemplatePage(c *gin.Context) {
	var r system2.MailTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetMailTemplatePage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*system2.MailTemplateRespVO, 0, len(page.List))
	for _, template := range page.List {
		list = append(list, &system2.MailTemplateRespVO{
			ID:         template.ID,
			Name:       template.Name,
			Code:       template.Code,
			AccountID:  template.AccountID,
			Nickname:   template.Nickname,
			Title:      template.Title,
			Content:    template.Content,
			Params:     template.Params,
			Status:     template.Status,
			Remark:     template.Remark,
			CreateTime: template.CreateTime,
		})
	}
	response.WritePage(c, page.Total, list)
}

func (h *MailHandler) SendMail(c *gin.Context) {
	var r system2.MailTemplateSendReq
	user := context.GetLoginUser(c)
	userId := int64(0)
	userType := 1 // AdminType
	if user != nil {
		userId = user.UserID
		userType = user.UserType
	}

	id, err := h.svc.SendSingleMail(c, r.ToMails, r.CcMails, r.BccMails, userId, userType, r.TemplateCode, r.TemplateParams)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// ================= Mail Log Request Handlers =================

func (h *MailHandler) GetMailLogPage(c *gin.Context) {
	var r system2.MailLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetMailLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*system2.MailLogRespVO, 0, len(page.List))
	for _, log := range page.List {
		resp := &system2.MailLogRespVO{
			ID:               log.ID,
			UserID:           log.UserID,
			UserType:         log.UserType,
			ToMails:          log.ToMails,
			CcMails:          log.CcMails,
			BccMails:         log.BccMails,
			AccountID:        log.AccountID,
			FromMail:         log.FromMail,
			TemplateID:       log.TemplateID,
			TemplateCode:     log.TemplateCode,
			TemplateNickname: log.TemplateNickname,
			TemplateTitle:    log.TemplateTitle,
			TemplateContent:  log.TemplateContent,
			SendStatus:       log.SendStatus,
			SendTime:         log.SendTime,
			SendMessageID:    log.SendMessageID,
			SendException:    log.SendException,
			CreateTime:       log.CreateTime,
		}
		if log.TemplateParams != "" {
			_ = json.Unmarshal([]byte(log.TemplateParams), &resp.TemplateParams)
		}
		list = append(list, resp)
	}
	response.WritePage(c, page.Total, list)
}

func (h *MailHandler) GetMailLog(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	log, err := h.svc.GetMailLog(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	resp := system2.MailLogRespVO{
		ID:               log.ID,
		UserID:           log.UserID,
		UserType:         log.UserType,
		ToMails:          log.ToMails,
		CcMails:          log.CcMails,
		BccMails:         log.BccMails,
		AccountID:        log.AccountID,
		FromMail:         log.FromMail,
		TemplateID:       log.TemplateID,
		TemplateCode:     log.TemplateCode,
		TemplateNickname: log.TemplateNickname,
		TemplateTitle:    log.TemplateTitle,
		TemplateContent:  log.TemplateContent,
		SendStatus:       log.SendStatus,
		SendTime:         log.SendTime,
		SendMessageID:    log.SendMessageID,
		SendException:    log.SendException,
		CreateTime:       log.CreateTime,
	}
	if log.TemplateParams != "" {
		_ = json.Unmarshal([]byte(log.TemplateParams), &resp.TemplateParams)
	}
	response.WriteSuccess(c, resp)
}
