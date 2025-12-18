package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type AppAfterSaleCreateReq struct {
	OrderItemID      int64    `json:"orderItemId" binding:"required"`
	RefundPrice      int      `json:"refundPrice" binding:"required"`
	Count            int      `json:"count" binding:"required"`
	Way              int      `json:"way" binding:"required"`  // 10: Refund only, 20: Return & Refund
	Type             int      `json:"type" binding:"required"` // 10: Want refund, 20: Not received
	ApplyReason      string   `json:"applyReason" binding:"required"`
	ApplyDescription string   `json:"applyDescription"`
	ApplyPicURLs     []string `json:"applyPicUrls"`
}

type AppAfterSaleCancelReq struct {
	ID int64 `json:"id" binding:"required"`
}

type AppAfterSalePageReq struct {
	pagination.PageParam
	Status *int `form:"status"`
}

type AppAfterSaleDeliveryReq struct {
	ID          int64  `json:"id" binding:"required"`
	LogisticsId int64  `json:"logisticsId" binding:"required"`
	LogisticsNo string `json:"logisticsNo" binding:"required"`
}

type TradeAfterSalePageReq struct {
	pagination.PageParam
	No          string   `form:"no"`
	UserID      *int64   `form:"userId"` // Admin can filter by userId, App sets it manually
	Status      *int     `form:"status"`
	OrderItemID *int64   `form:"orderItemId"`
	OrderNo     string   `form:"orderNo"`
	SpuName     string   `form:"spuName"`
	CreateTime  []string `form:"createTime[]"`
}

type TradeAfterSaleAgreeReq struct {
	ID int64 `json:"id" binding:"required"`
}

type TradeAfterSaleDisagreeReq struct {
	ID          int64  `json:"id" binding:"required"`
	AuditReason string `json:"auditReason" binding:"required"`
}

type TradeAfterSaleRefundReq struct {
	ID int64 `json:"id" binding:"required"`
}

type TradeAfterSaleRefuseReq struct {
	ID         int64  `json:"id" binding:"required"`
	RefuseMemo string `json:"refuseMemo" binding:"required"`
}
