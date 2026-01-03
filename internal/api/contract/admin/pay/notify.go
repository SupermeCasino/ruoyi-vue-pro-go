package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// Notify Template Requests
type NotifyTemplateCreateReq struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Type     int    `json:"type" binding:"required"`
	Status   int    `json:"status" binding:"required"`
	Remark   string `json:"remark"`
}

type NotifyTemplateUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Type     int    `json:"type" binding:"required"`
	Status   int    `json:"status" binding:"required"`
	Remark   string `json:"remark"`
}

type NotifyTemplatePageReq struct {
	pagination.PageParam
	Name      string `form:"name"`
	Code      string `form:"code"`
	Status    *int   `form:"status"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

type NotifyTemplateSendReq struct {
	UserID         int64                  `json:"userId" binding:"required"`
	UserType       int                    `json:"userType" binding:"required"`
	TemplateCode   string                 `json:"templateCode" binding:"required"`
	TemplateParams map[string]interface{} `json:"templateParams"`
}

// Notify Message Requests
type NotifyMessagePageReq struct {
	pagination.PageParam
	UserID       int64  `form:"userId"`
	UserType     int    `form:"userType"`
	TemplateCode string `form:"templateCode"`
	TemplateType *int   `form:"templateType"`
	ReadStatus   *bool  `form:"readStatus"`
	StartDate    string `form:"startDate"`
	EndDate      string `form:"endDate"`
}

// My Message Page Req
type MyNotifyMessagePageReq struct {
	pagination.PageParam
	ReadStatus *bool `form:"readStatus"`
}

type NotifyMessageUpdateReadReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type NotifyMessageReadAllReq struct {
}
