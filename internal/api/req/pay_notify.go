package req

import (
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
