package req

import (
	"backend-go/internal/pkg/core"
)

// SmsTemplateSaveReq 短信模板创建/修改 Request
type SmsTemplateSaveReq struct {
	ID            int64  `json:"id"`
	Type          int32  `json:"type" binding:"required"`
	Status        int32  `json:"status" binding:"required"`
	Code          string `json:"code" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Content       string `json:"content" binding:"required"`
	Remark        string `json:"remark"`
	ApiTemplateId string `json:"apiTemplateId" binding:"required"`
	ChannelId     int64  `json:"channelId" binding:"required"`
}

// SmsTemplatePageReq 短信模板分页 Request
type SmsTemplatePageReq struct {
	core.PageParam
	Type          *int32   `form:"type"`
	Status        *int32   `form:"status"`
	Code          string   `form:"code"`
	Content       string   `form:"content"`
	ApiTemplateId string   `form:"apiTemplateId"`
	ChannelId     *int64   `form:"channelId"`
	CreateTime    []string `form:"createTime[]"`
}
