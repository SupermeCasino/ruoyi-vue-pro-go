package system

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// SmsChannelSaveReq 短信渠道创建/修改 Request
type SmsChannelSaveReq struct {
	ID          int64  `json:"id"`
	Signature   string `json:"signature" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Status      *int32 `json:"status" binding:"required"`
	Remark      string `json:"remark"`
	ApiKey      string `json:"apiKey" binding:"required"`
	ApiSecret   string `json:"apiSecret"`
	CallbackUrl string `json:"callbackUrl" binding:"omitempty,url"`
}

// SmsChannelPageReq 短信渠道分页 Request
type SmsChannelPageReq struct {
	pagination.PageParam
	Signature string `form:"signature"`
	Status    *int32 `form:"status"`
}

// SmsChannelRespVO 短信渠道信息 Response
type SmsChannelRespVO struct {
	ID          int64     `json:"id"`
	Signature   string    `json:"signature"`
	Code        string    `json:"code"`
	Status      int32     `json:"status"`
	Remark      string    `json:"remark"`
	ApiKey      string    `json:"apiKey"`
	ApiSecret   string    `json:"apiSecret"`
	CallbackUrl string    `json:"callbackUrl"`
	CreateTime  time.Time `json:"createTime"`
}

// SmsChannelSimpleRespVO 短信渠道精简信息 Response
type SmsChannelSimpleRespVO struct {
	ID        int64  `json:"id"`
	Signature string `json:"signature"`
	Code      string `json:"code"`
}
