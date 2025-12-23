package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

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
	CreateTime          time.Time         `json:"createTime"`
	UpdateTime          time.Time         `json:"updateTime"`
	Creator            string            `json:"creator"`
	Updater            string            `json:"updater"`
	Deleted            model.BitBool     `json:"deleted"`
	TenantID           int64             `json:"tenantId"`
}
