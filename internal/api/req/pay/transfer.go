package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// PayTransferPageReq 转账单分页请求
type PayTransferCreateReq struct {
	AppID              int64             `json:"appId" binding:"required"`
	ChannelCode        string            `json:"channelCode" binding:"required"`
	MerchantTransferID string            `json:"merchantTransferId" binding:"required"`
	Type               int               `json:"type" binding:"required"`
	Price              int               `json:"price" binding:"required"`
	Subject            string            `json:"subject" binding:"required"`
	UserName           string            `json:"userName" binding:"required"`
	UserAccount        string            `json:"userAccount"`   // 用户账号
	AlipayLogonID      string            `json:"alipayLogonId"` // 支付宝登录号
	OpenID             string            `json:"openid"`        // 微信 OpenID
	NotifyURL          string            `json:"notifyUrl"`
	UserIP             string            `json:"userIp"`
	ChannelExtras      map[string]string `json:"channelExtras"`
}

type PayTransferPageReq struct {
	pagination.PageParam
	No                string   `form:"no" json:"no"`                               // 转账单号
	AppID             int64    `form:"appId" json:"appId"`                         // 应用编号
	ChannelCode       string   `form:"channelCode" json:"channelCode"`             // 渠道编码
	MerchantOrderId   string   `form:"merchantOrderId" json:"merchantOrderId"`     // 商户转账单编号
	Status            *int     `form:"status" json:"status"`                       // 转账状态
	UserName          string   `form:"userName" json:"userName"`                   // 收款人姓名
	UserAccount       string   `form:"userAccount" json:"userAccount"`             // 收款人账号
	ChannelTransferNo string   `form:"channelTransferNo" json:"channelTransferNo"` // 渠道转账单号
	CreateTime        []string `form:"createTime[]" json:"createTime"`             // 创建时间
}

// PayTransferReq 转账单详情请求
type PayTransferReq struct {
	ID int64 `form:"id" json:"id"`
}
