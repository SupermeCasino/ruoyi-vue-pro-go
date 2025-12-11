package req

import (
	"backend-go/internal/pkg/core"
)

// MailAccount Requests
type MailAccountCreateReq struct {
	Mail           string `json:"mail" binding:"required,email"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Host           string `json:"host" binding:"required"`
	Port           int    `json:"port" binding:"required"`
	SslEnable      bool   `json:"sslEnable"`
	StarttlsEnable bool   `json:"starttlsEnable"`
}

type MailAccountUpdateReq struct {
	ID             int64  `json:"id" binding:"required"`
	Mail           string `json:"mail" binding:"required,email"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Host           string `json:"host" binding:"required"`
	Port           int    `json:"port" binding:"required"`
	SslEnable      bool   `json:"sslEnable"`
	StarttlsEnable bool   `json:"starttlsEnable"`
}

type MailAccountPageReq struct {
	core.PageParam
	Mail     string `form:"mail"`
	Username string `form:"username"`
}

// MailTemplate Requests
type MailTemplateCreateReq struct {
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	AccountID int64  `json:"accountId" binding:"required"`
	Nickname  string `json:"nickname"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	Remark    string `json:"remark"`
}

type MailTemplateUpdateReq struct {
	ID        int64  `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	AccountID int64  `json:"accountId" binding:"required"`
	Nickname  string `json:"nickname"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	Remark    string `json:"remark"`
}

type MailTemplatePageReq struct {
	core.PageParam
	Name      string `form:"name"`
	Code      string `form:"code"`
	AccountID int64  `form:"accountId"`
	Status    *int   `form:"status"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

type MailTemplateSendReq struct {
	TemplateCode   string                 `json:"templateCode" binding:"required"`
	ToMail         string                 `json:"toMail" binding:"required,email"`
	TemplateParams map[string]interface{} `json:"templateParams"`
}

// MailLog Requests
type MailLogPageReq struct {
	core.PageParam
	UserID     int64  `form:"userId"`
	UserType   int    `form:"userType"`
	ToMail     string `form:"toMail"`
	AccountID  int64  `form:"accountId"`
	TemplateID int64  `form:"templateId"`
	SendStatus *int   `form:"sendStatus"`
	StartDate  string `form:"startDate"`
	EndDate    string `form:"endDate"`
}
