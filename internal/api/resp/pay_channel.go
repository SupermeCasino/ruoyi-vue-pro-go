package resp

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"time"
)

type PayChannelResp struct {
	ID         int64                `json:"id"`
	Code       string               `json:"code"`
	Status     int                  `json:"status"`
	FeeRate    float64              `json:"feeRate"`
	Remark     string               `json:"remark"`
	AppID      int64                `json:"appId"`
	Config     *pay.PayClientConfig `json:"config"` // 支付渠道配置
	CreateTime time.Time            `json:"createTime"`
}
