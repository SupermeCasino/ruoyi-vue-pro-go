package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayRefundCreateReq struct {
	AppKey           string `json:"appKey"`
	UserIP           string `json:"userIp"`
	MerchantOrderId  string `json:"merchantOrderId"`
	MerchantRefundId string `json:"merchantRefundId"`
	Reason           string `json:"reason"`
	Price            int    `json:"price"`
}

type PayRefundPageReq struct {
	pagination.PageParam
	AppID            int64    `form:"appId"`
	ChannelCode      string   `form:"channelCode"`
	MerchantOrderId  string   `form:"merchantOrderId"`
	MerchantRefundId string   `form:"merchantRefundId"`
	ChannelOrderNo   string   `form:"channelOrderNo"`
	ChannelRefundNo  string   `form:"channelRefundNo"`
	Status           *int     `form:"status"`
	CreateTime       []string `form:"createTime[]"` // Range search
}

type PayRefundExportReq struct {
	AppID            int64    `form:"appId"`
	ChannelCode      string   `form:"channelCode"`
	MerchantOrderId  string   `form:"merchantOrderId"`
	MerchantRefundId string   `form:"merchantRefundId"`
	ChannelOrderNo   string   `form:"channelOrderNo"`
	ChannelRefundNo  string   `form:"channelRefundNo"`
	Status           *int     `form:"status"`
	CreateTime       []string `form:"createTime[]"` // Range search
}

// PayRefundNotifyReqDTO 支付退款回调通知 Request DTO
type PayRefundNotifyReqDTO struct {
	MerchantOrderId  string `json:"merchantOrderId"`
	MerchantRefundId string `json:"merchantRefundId"`
	PayRefundId      int64  `json:"payRefundId"`
	Status           int    `json:"status"` // PayRefundStatusEnum
}

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
	PayPrice          int64      `json:"payPrice"`    // 改为 int64
	RefundPrice       int64      `json:"refundPrice"` // 改为 int64
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
	App   *PayAppResp  `json:"app"`
	Order *RefundOrder `json:"order"` // 关联的原订单信息
}

// RefundOrder 退款中的订单信息（用于展示原订单详情）
type RefundOrder struct {
	Subject string `json:"subject"`
}
