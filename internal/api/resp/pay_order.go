package resp

import (
	"time"
)

type PayOrderResp struct {
	ID              int64      `json:"id"`
	AppID           int64      `json:"appId"`
	AppName         string     `json:"appName"` // From PayApp
	ChannelID       int64      `json:"channelId"`
	ChannelCode     string     `json:"channelCode"`
	MerchantOrderId string     `json:"merchantOrderId"`
	Subject         string     `json:"subject"`
	Body            string     `json:"body"`
	NotifyURL       string     `json:"notifyUrl"`
	Price           int64      `json:"price"` // 改为 int64
	ChannelFeeRate  float64    `json:"channelFeeRate"`
	ChannelFeePrice int        `json:"channelFeePrice"`
	Status          int        `json:"status"`
	UserIP          string     `json:"userIp"`
	ExpireTime      time.Time  `json:"expireTime"`
	SuccessTime     *time.Time `json:"successTime"`
	ExtensionID     int64      `json:"extensionId"`
	No              string     `json:"no"`
	RefundPrice     int64      `json:"refundPrice"` // 改为 int64
	ChannelUserID   string     `json:"channelUserId"`
	ChannelOrderNo  string     `json:"channelOrderNo"`
	CreateTime      time.Time  `json:"createTime"`
	UpdateTime      time.Time  `json:"updateTime"`
	Creator         string     `json:"creator"`
	Updater         string     `json:"updater"`
}

type PayOrderDetailsResp struct {
	PayOrderResp
	Extension *PayOrderExtensionResp `json:"extension"`
	App       *PayAppResp            `json:"app"`
}

type PayOrderExtensionResp struct {
	ID                int64     `json:"id"`
	No                string    `json:"no"`
	OrderID           int64     `json:"orderId"`
	ChannelID         int64     `json:"channelId"`
	ChannelCode       string    `json:"channelCode"`
	UserIP            string    `json:"userIp"`
	Status            int       `json:"status"`
	ChannelExtras     string    `json:"channelExtras"`
	ChannelErrorCode  string    `json:"channelErrorCode"`
	ChannelErrorMsg   string    `json:"channelErrorMsg"`
	ChannelNotifyData string    `json:"channelNotifyData"`
	CreateTime        time.Time `json:"createTime"`
}

type PayOrderSubmitResp struct {
	Status         int    `json:"status"`
	DisplayMode    string `json:"displayMode"`
	DisplayContent string `json:"displayContent"`
}
