package req

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
)

type PayNotifyTaskPageReq struct {
	core.PageParam
	AppID           int64    `form:"appId"`
	Type            *int     `form:"type"`
	DataID          int64    `form:"dataId"`
	MerchantOrderId string   `form:"merchantOrderId"`
	Status          *int     `form:"status"`
	CreateTime      []string `form:"createTime[]"` // Range search
}
