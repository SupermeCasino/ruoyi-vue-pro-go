package system

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// MailAccount Requests
type MailAccountSaveReq struct {
	ID             int64  `json:"id"`
	Mail           string `json:"mail" binding:"required,email"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Host           string `json:"host" binding:"required"`
	Port           int    `json:"port" binding:"required"`
	SslEnable      bool   `json:"sslEnable"`
	StarttlsEnable bool   `json:"starttlsEnable"`
}

// 移除 MailAccountUpdateReq，统一使用 MailAccountSaveReq

type MailAccountPageReq struct {
	pagination.PageParam
	Mail     string `form:"mail"`
	Username string `form:"username"`
}

// MailTemplate Requests
type MailTemplateSaveReq struct {
	ID        int64  `json:"id"`
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	AccountID int64  `json:"accountId" binding:"required"`
	Nickname  string `json:"nickname"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	Remark    string `json:"remark"`
}

// 移除 MailTemplateUpdateReq，统一使用 MailTemplateSaveReq

type MailTemplatePageReq struct {
	pagination.PageParam
	Name       string   `form:"name"`
	Code       string   `form:"code"`
	AccountID  int64    `form:"accountId"`
	Status     *int     `form:"status"`
	CreateTime []string `form:"createTime"`
}

type MailTemplateSendReq struct {
	TemplateCode   string                 `json:"templateCode" binding:"required"`
	ToMails        []string               `json:"toMails" binding:"required"`
	CcMails        []string               `json:"ccMails"`
	BccMails       []string               `json:"bccMails"`
	TemplateParams map[string]interface{} `json:"templateParams"`
}

// MailLog Requests
type MailLogPageReq struct {
	pagination.PageParam
	UserID     int64    `form:"userId"`
	UserType   int      `form:"userType"`
	ToMail     string   `form:"toMail"`
	AccountID  int64    `form:"accountId"`
	TemplateID int64    `form:"templateId"`
	SendStatus *int     `form:"sendStatus"`
	SendTime   []string `form:"sendTime"`
}

// MailAccount Responses
type MailAccountRespVO struct {
	ID             int64     `json:"id"`
	Mail           string    `json:"mail"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Host           string    `json:"host"`
	Port           int       `json:"port"`
	SslEnable      bool      `json:"sslEnable"`
	StarttlsEnable bool      `json:"starttlsEnable"`
	CreateTime     time.Time `json:"createTime"`
}

type MailAccountSimpleRespVO struct {
	ID   int64  `json:"id"`
	Mail string `json:"mail"`
}

// MailTemplate Responses
type MailTemplateRespVO struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	AccountID  int64     `json:"accountId"`
	Nickname   string    `json:"nickname"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Params     []string  `json:"params"`
	Status     int       `json:"status"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}

type MailTemplateSimpleRespVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// MailLog Responses
type MailLogRespVO struct {
	ID               int64                  `json:"id"`
	UserID           int64                  `json:"userId"`
	UserType         int                    `json:"userType"`
	ToMails          []string               `json:"toMails"`
	CcMails          []string               `json:"ccMails"`
	BccMails         []string               `json:"bccMails"`
	AccountID        int64                  `json:"accountId"`
	FromMail         string                 `json:"fromMail"`
	TemplateID       int64                  `json:"templateId"`
	TemplateCode     string                 `json:"templateCode"`
	TemplateNickname string                 `json:"templateNickname"`
	TemplateTitle    string                 `json:"templateTitle"`
	TemplateContent  string                 `json:"templateContent"`
	TemplateParams   map[string]interface{} `json:"templateParams"`
	SendStatus       int                    `json:"sendStatus"`
	SendTime         *time.Time             `json:"sendTime"`
	SendMessageID    string                 `json:"sendMessageId"`
	SendException    string                 `json:"sendException"`
	CreateTime       time.Time              `json:"createTime"`
}
