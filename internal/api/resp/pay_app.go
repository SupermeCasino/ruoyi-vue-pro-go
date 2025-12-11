package resp

import "time"

type PayAppResp struct {
	ID                int64     `json:"id"`
	AppKey            string    `json:"appKey"`
	Name              string    `json:"name"`
	Status            int       `json:"status"`
	Remark            string    `json:"remark"`
	OrderNotifyURL    string    `json:"orderNotifyUrl"`
	RefundNotifyURL   string    `json:"refundNotifyUrl"`
	TransferNotifyURL string    `json:"transferNotifyUrl"`
	CreateTime        time.Time `json:"createTime"`
}

type PayAppPageItemResp struct {
	PayAppResp
	ChannelCodes []string `json:"channelCodes"`
}
