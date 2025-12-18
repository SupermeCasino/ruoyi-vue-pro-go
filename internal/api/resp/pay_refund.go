package resp

import (
	"time"
)

type PayRefundResp struct {
	ID                int64      `json:"id"`
	No                string     `json:"no"`
	AppID             int64      `json:"appId"`
	ChannelID         int64      `json:"channelId"`
	ChannelCode       string     `json:"channelCode"`
	OrderID           int64      `json:"orderId"`
	OrderNo           string     `json:"orderNo"`
	UserID            int64      `json:"userId"`
	UserType          int        `json:"userType"`
	MerchantOrderId   string     `json:"merchantOrderId"`
	MerchantRefundId  string     `json:"merchantRefundId"`
	NotifyURL         string     `json:"notifyUrl"`
	Status            int        `json:"status"`
	PayPrice          int64      `json:"payPrice"`        // 改为 int64
	RefundPrice       int64      `json:"refundPrice"`     // 改为 int64
	Reason            string     `json:"reason"`
	UserIP            string     `json:"userIp"`
	ChannelOrderNo    string     `json:"channelOrderNo"`
	ChannelRefundNo   string     `json:"channelRefundNo"`
	SuccessTime       *time.Time `json:"successTime"`
	ChannelErrorCode  string     `json:"channelErrorCode"`
	ChannelErrorMsg   string     `json:"channelErrorMsg"`
	ChannelNotifyData string     `json:"channelNotifyData"`
	CreateTime        time.Time  `json:"createTime"`
	UpdateTime        time.Time  `json:"updateTime"`
	AppName           string     `json:"appName"` // Enrichment
}

// PayRefundDetailsResp 详情 Response
type PayRefundDetailsResp struct {
	PayRefundResp
	App   *PayAppResp   `json:"app"`
	Order *RefundOrder  `json:"order"` // 关联的原订单信息
}

// RefundOrder 退款中的订单信息（用于展示原订单详情）
type RefundOrder struct {
	Subject string `json:"subject"`
}
