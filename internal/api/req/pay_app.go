package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayAppCreateReq struct {
	AppKey            string `json:"appKey" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Status            int    `json:"status" binding:"required"`
	Remark            string `json:"remark"`
	OrderNotifyURL    string `json:"orderNotifyUrl" binding:"required,url"`
	RefundNotifyURL   string `json:"refundNotifyUrl" binding:"required,url"`
	TransferNotifyURL string `json:"transferNotifyUrl" binding:"url"`
}

type PayAppUpdateReq struct {
	ID int64 `json:"id" binding:"required"`
	PayAppCreateReq
}

type PayAppUpdateStatusReq struct {
	ID     int64 `json:"id" binding:"required"`
	Status int   `json:"status" binding:"required"`
}

type PayAppPageReq struct {
	pagination.PageParam
	Name   string `json:"name" form:"name"`
	Status *int   `json:"status" form:"status"`
	Remark string `json:"remark" form:"remark"`
}
