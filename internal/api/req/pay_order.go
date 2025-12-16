package req

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
)

type PayOrderPageReq struct {
	core.PageParam
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
	AppID           int64  `json:"appId" binding:"required"`
	UserIP          string `json:"userIp" binding:"required"`
	MerchantOrderId string `json:"merchantOrderId" binding:"required"`
	Subject         string `json:"subject" binding:"required"`
	Body            string `json:"body"`
	NotifyUrl       string `json:"notifyUrl" binding:"required"`
	Price           int    `json:"price" binding:"required,min=0"`
	ExpireTime      string `json:"expireTime" binding:"required"` // time string
}
