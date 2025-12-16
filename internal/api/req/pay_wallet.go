package req

import (
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
	WalletID int64  `form:"walletId"`
	BizType  int    `form:"bizType"`
	BizID    string `form:"bizId"`
	No       string `form:"no"`
	Title    string `form:"title"`
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
