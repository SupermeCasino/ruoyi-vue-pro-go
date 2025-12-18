package req

import (
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
