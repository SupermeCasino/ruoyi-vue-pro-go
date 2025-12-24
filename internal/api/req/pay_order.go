package req

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayOrderPageReq struct {
	pagination.PageParam
	AppID           int64    `form:"appId"`
	ChannelCode     string   `form:"channelCode"`
	MerchantOrderId string   `form:"merchantOrderId"`
	Subject         string   `form:"subject"`
	No              string   `form:"no"`
	Status          *int     `form:"status"`
	CreateTime      []string `form:"createTime[]"` // time range
}

type PayOrderExportReq struct {
	AppID           int64    `form:"appId"`
	ChannelCode     string   `form:"channelCode"`
	MerchantOrderId string   `form:"merchantOrderId"`
	Subject         string   `form:"subject"`
	No              string   `form:"no"`
	Status          *int     `form:"status"`
	CreateTime      []string `form:"createTime[]"`
}

type PayOrderSubmitReq struct {
	ID            int64             `json:"id" binding:"required"`
	ChannelCode   string            `json:"channelCode" binding:"required"`
	ChannelExtras map[string]string `json:"channelExtras"`
	DisplayMode   string            `json:"displayMode"` // PayOrderDisplayModeEnum
	ReturnUrl     string            `json:"returnUrl"`
}

type PayOrderCreateReq struct {
	AppKey          string    `json:"appKey" binding:"required"`
	UserIP          string    `json:"userIp" binding:"required"`
	MerchantOrderId string    `json:"merchantOrderId" binding:"required"`
	Subject         string    `json:"subject" binding:"required"`
	Body            string    `json:"body"`
	Price           int       `json:"price" binding:"required,min=0"`
	ExpireTime      time.Time `json:"expireTime" binding:"required"`
}

type PayOrderNotifyReq struct {
	MerchantOrderId string `json:"merchantOrderId" binding:"required"`
	PayOrderID      int64  `json:"payOrderId" binding:"required"`
}
