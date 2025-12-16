package req

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

// SmsChannelSaveReq 短信渠道创建/修改 Request
type SmsChannelSaveReq struct {
	ID          int64  `json:"id"`
	Signature   string `json:"signature" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Status      int32  `json:"status" binding:"required"`
	Remark      string `json:"remark"`
	ApiKey      string `json:"apiKey" binding:"required"`
	ApiSecret   string `json:"apiSecret"`
	CallbackUrl string `json:"callbackUrl" binding:"omitempty,url"`
}

// SmsChannelPageReq 短信渠道分页 Request
type SmsChannelPageReq struct {
	core.PageParam
	Signature string `form:"signature"`
	Status    *int32 `form:"status"`
}
