package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
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

// PayTransferResp 转账单 Response VO
// PayTransferCreateResp 转账单创建响应
type PayTransferCreateResp struct {
	ID                 int64  `json:"id"`
	Status             int    `json:"status"`
	ChannelPackageInfo string `json:"channelPackageInfo,omitempty"` // 渠道 package 信息 (WeChat)
}

type PayTransferResp struct {
	ID                 int64             `json:"id"`
	No                 string            `json:"no"`
	AppID              int64             `json:"appId"`
	AppName            string            `json:"appName"` // From PayApp
	ChannelID          int64             `json:"channelId"`
	ChannelCode        string            `json:"channelCode"`
	MerchantTransferID string            `json:"merchantTransferId"`
	Subject            string            `json:"subject"`
	Price              int               `json:"price"`
	UserAccount        string            `json:"userAccount"`
	UserName           string            `json:"userName"`
	Status             int               `json:"status"`
	SuccessTime        *time.Time        `json:"successTime"`
	NotifyURL          string            `json:"notifyUrl"`
	UserIP             string            `json:"userIp"`
	ChannelExtras      map[string]string `json:"channelExtras"`
	ChannelTransferNo  string            `json:"channelTransferNo"`
	ChannelErrorCode   string            `json:"channelErrorCode"`
	ChannelErrorMsg    string            `json:"channelErrorMsg"`
	ChannelNotifyData  string            `json:"channelNotifyData"`
	ChannelPackageInfo string            `json:"channelPackageInfo"`
	CreateTime         time.Time         `json:"createTime"`
	UpdateTime         time.Time         `json:"updateTime"`
	Creator            string            `json:"creator"`
	Updater            string            `json:"updater"`
	Deleted            model.BitBool     `json:"deleted"`
	TenantID           int64             `json:"tenantId"`
}
