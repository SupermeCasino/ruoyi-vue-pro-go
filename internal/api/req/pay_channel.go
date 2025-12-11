package req

import "backend-go/internal/pkg/core"

type PayChannelCreateReq struct {
	Code    string  `json:"code" binding:"required"`
	Status  int     `json:"status" binding:"required"`
	FeeRate float64 `json:"feeRate" binding:"required"`
	Remark  string  `json:"remark"`
	AppID   int64   `json:"appId" binding:"required"`
	Config  string  `json:"config" binding:"required"` // JSON String
}

type PayChannelUpdateReq struct {
	ID int64 `json:"id" binding:"required"`
	PayChannelCreateReq
}

type PayChannelPageReq struct {
	core.PageParam
	Code   string `json:"code" form:"code"`
	Status *int   `json:"status" form:"status"`
	AppID  int64  `json:"appId" form:"appId"`
}
