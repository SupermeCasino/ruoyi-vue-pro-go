package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayChannelCreateReq struct {
	Code    string               `json:"code" binding:"required"`
	Status  int                  `json:"status" binding:"required"`
	FeeRate float64              `json:"feeRate" binding:"required"`
	Remark  string               `json:"remark"`
	AppID   int64                `json:"appId" binding:"required"`
	Config  *pay.PayClientConfig `json:"config" binding:"required"` // 支付渠道配置
}

type PayChannelUpdateReq struct {
	ID int64 `json:"id" binding:"required"`
	PayChannelCreateReq
}

type PayChannelPageReq struct {
	pagination.PageParam
	Code   string `json:"code" form:"code"`
	Status *int   `json:"status" form:"status"`
	AppID  int64  `json:"appId" form:"appId"`
}

type PayChannelResp struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Status     int       `json:"status"`
	FeeRate    float64   `json:"feeRate"`
	Remark     string    `json:"remark"`
	AppID      int64     `json:"appId"`
	Config     string    `json:"config"` // 支付渠道配置
	CreateTime time.Time `json:"createTime"`
}
