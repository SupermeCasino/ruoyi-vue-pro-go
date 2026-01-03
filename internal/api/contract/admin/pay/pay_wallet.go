package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// PayWalletRechargePackageCreateReq
type PayWalletRechargePackageCreateReq struct {
	Name       string `json:"name" binding:"required"`
	PayPrice   int    `json:"payPrice" binding:"required"`
	BonusPrice int    `json:"bonusPrice" binding:"required"`
	Status     int    `json:"status" binding:"required"`
}

// PayWalletRechargePackageUpdateReq
type PayWalletRechargePackageUpdateReq struct {
	ID         int64  `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	PayPrice   int    `json:"payPrice" binding:"required"`
	BonusPrice int    `json:"bonusPrice" binding:"required"`
	Status     int    `json:"status" binding:"required"`
}

// PayWalletRechargePackagePageReq
type PayWalletRechargePackagePageReq struct {
	pagination.PageParam
	Name   string `form:"name"`
	Status *int   `form:"status"`
}

// PayWalletRechargeCreateReq
type PayWalletRechargeCreateReq struct {
	UserID     int64 `json:"userId"`   // 可选，后台充值可能需要
	UserType   int   `json:"userType"` // 可选
	PayPrice   int   `json:"payPrice" binding:"required"`
	BonusPrice int   `json:"bonusPrice"`
	PackageID  int64 `json:"packageId"`
}

// PayWalletRechargePageReq
type PayWalletRechargePageReq struct {
	pagination.PageParam
	PayStatus *bool `form:"payStatus"`
}

// PayWalletTransactionPageReq
type PayWalletTransactionPageReq struct {
	pagination.PageParam
	WalletID   int64    `form:"walletId"`
	BizType    int      `form:"bizType"`
	BizID      string   `form:"bizId"`
	No         string   `form:"no"`
	Title      string   `form:"title"`
	Type       int      `form:"type"`         // 1: 收入, 2: 支出
	CreateTime []string `form:"createTime[]"` // [startTime, endTime]
}

// PayWalletPageReq
type PayWalletPageReq struct {
	pagination.PageParam
	UserID   int64 `form:"userId"`
	UserType int   `form:"userType"`
}

type PayWalletUpdateBalanceReq struct {
	UserID  int64 `json:"userId" binding:"required"`
	Balance int   `json:"balance" binding:"required"` // Unit: fen, can be negative
}

// PayWalletRechargePackageResp
type PayWalletRechargePackageResp struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	PayPrice   int       `json:"payPrice"`
	BonusPrice int       `json:"bonusPrice"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"createTime"`
}

// PayWalletTransactionResp
type PayWalletTransactionResp struct {
	ID         int64     `json:"id"`
	WalletID   int64     `json:"walletId"`
	BizType    int       `json:"bizType"`
	BizID      string    `json:"bizId"`
	No         string    `json:"no"`
	Title      string    `json:"title"`
	Price      int       `json:"price"`
	Balance    int       `json:"balance"`
	CreateTime time.Time `json:"createTime"`
}

// PayWalletResp
type PayWalletResp struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"userId"`
	UserType      int       `json:"userType"`
	Balance       int       `json:"balance"`
	TotalExpense  int       `json:"totalExpense"`
	TotalRecharge int       `json:"totalRecharge"`
	FreezePrice   int       `json:"freezePrice"`
	CreateTime    time.Time `json:"createTime"`
}

// PayWalletRechargeResp
type PayWalletRechargeResp struct {
	ID               int64      `json:"id"`
	WalletID         int64      `json:"walletId"`
	TotalPrice       int        `json:"totalPrice"`
	PayPrice         int        `json:"payPrice"`
	BonusPrice       int        `json:"bonusPrice"`
	PackageID        int64      `json:"packageId"`
	PayStatus        bool       `json:"payStatus"`
	PayOrderID       int64      `json:"payOrderId"`
	PayChannelCode   string     `json:"payChannelCode"`
	PayTime          *time.Time `json:"payTime"`
	RefundStatus     int        `json:"refundStatus"`
	PayRefundID      int64      `json:"payRefundId"`
	RefundTotalPrice int        `json:"refundTotalPrice"`
	RefundPayPrice   int        `json:"refundPayPrice"`
	RefundBonusPrice int        `json:"refundBonusPrice"`
	RefundTime       *time.Time `json:"refundTime"`
	CreateTime       time.Time  `json:"createTime"`
}
