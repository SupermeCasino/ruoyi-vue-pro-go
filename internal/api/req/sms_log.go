package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// SmsLogPageReq 短信日志分页 Request
type SmsLogPageReq struct {
	pagination.PageParam
	ChannelId     *int64   `form:"channelId"`
	TemplateId    *int64   `form:"templateId"`
	Mobile        string   `form:"mobile"`
	SendStatus    *int32   `form:"sendStatus"`
	SendTime      []string `form:"sendTime[]"`
	ReceiveStatus *int32   `form:"receiveStatus"`
	ReceiveTime   []string `form:"receiveTime[]"`
}
