package resp

import "time"

type PayChannelResp struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Status     int       `json:"status"`
	FeeRate    float64   `json:"feeRate"`
	Remark     string    `json:"remark"`
	AppID      int64     `json:"appId"`
	Config     string    `json:"config"` // JSON String
	CreateTime time.Time `json:"createTime"`
}
