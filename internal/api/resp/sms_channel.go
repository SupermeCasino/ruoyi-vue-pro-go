package resp

import "time"

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
