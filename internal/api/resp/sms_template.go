package resp

import (
	"time"

	"gorm.io/datatypes"
)

// SmsTemplateRespVO 短信模板信息 Response
type SmsTemplateRespVO struct {
	ID            int64                       `json:"id"`
	Type          int32                       `json:"type"`
	Status        int32                       `json:"status"`
	Code          string                      `json:"code"`
	Name          string                      `json:"name"`
	Content       string                      `json:"content"`
	Params        datatypes.JSONSlice[string] `json:"params"`
	Remark        string                      `json:"remark"`
	ApiTemplateId string                      `json:"apiTemplateId"`
	ChannelId     int64                       `json:"channelId"`
	ChannelCode   string                      `json:"channelCode"`
	CreateTime    time.Time                   `json:"createTime"`
}
