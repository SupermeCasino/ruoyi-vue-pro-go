package resp

import "time"

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
