package system

import (
	"time"

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

// SmsLogRespVO 短信日志 Response
type SmsLogRespVO struct {
	ID              int64                  `json:"id"`
	ChannelId       int64                  `json:"channelId"`
	ChannelCode     string                 `json:"channelCode"`
	TemplateId      int64                  `json:"templateId"`
	TemplateCode    string                 `json:"templateCode"`
	TemplateType    int32                  `json:"templateType"`
	TemplateContent string                 `json:"templateContent"`
	TemplateParams  map[string]interface{} `json:"templateParams"`
	ApiTemplateId   string                 `json:"apiTemplateId"`
	Mobile          string                 `json:"mobile"`
	UserId          int64                  `json:"userId"`
	UserType        int32                  `json:"userType"`
	SendStatus      int32                  `json:"sendStatus"`
	SendTime        *time.Time             `json:"sendTime"`
	ApiSendCode     string                 `json:"apiSendCode"`
	ApiSendMsg      string                 `json:"apiSendMsg"`
	ApiRequestId    string                 `json:"apiRequestId"`
	ApiSerialNo     string                 `json:"apiSerialNo"`
	ReceiveStatus   int32                  `json:"receiveStatus"`
	ReceiveTime     *time.Time             `json:"receiveTime"`
	ApiReceiveCode  string                 `json:"apiReceiveCode"`
	ApiReceiveMsg   string                 `json:"apiReceiveMsg"`
	CreateTime      time.Time              `json:"createTime"`
}
