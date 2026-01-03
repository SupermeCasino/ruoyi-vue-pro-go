package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayNotifyTaskPageReq struct {
	pagination.PageParam
	AppID           int64    `form:"appId"`
	Type            *int     `form:"type"`
	DataID          int64    `form:"dataId"`
	MerchantOrderId string   `form:"merchantOrderId"`
	Status          *int     `form:"status"`
	CreateTime      []string `form:"createTime[]"` // Range search
}
type PayNotifyTaskResp struct {
	ID                 int64      `json:"id"`
	AppID              int64      `json:"appId"`
	Type               int        `json:"type"`
	DataID             int64      `json:"dataId"`
	MerchantOrderId    string     `json:"merchantOrderId"`
	MerchantRefundId   string     `json:"merchantRefundId"`
	MerchantTransferId string     `json:"merchantTransferId"`
	Status             int        `json:"status"`
	NextNotifyTime     *time.Time `json:"nextNotifyTime"`
	LastExecuteTime    *time.Time `json:"lastExecuteTime"`
	NotifyTimes        int        `json:"notifyTimes"`
	MaxNotifyTimes     int        `json:"maxNotifyTimes"`
	NotifyURL          string     `json:"notifyUrl"`
	CreateTime         time.Time  `json:"createTime"`
	UpdateTime         time.Time  `json:"updateTime"`
	AppName            string     `json:"appName"` // Enrichment
}

type PayNotifyTaskDetailResp struct {
	PayNotifyTaskResp
	Logs []*PayNotifyLogResp `json:"logs"`
}

type PayNotifyLogResp struct {
	ID          int64     `json:"id"`
	TaskID      int64     `json:"taskId"`
	NotifyTimes int       `json:"notifyTimes"`
	Response    string    `json:"response"`
	Status      int       `json:"status"`
	CreateTime  time.Time `json:"createTime"`
}
