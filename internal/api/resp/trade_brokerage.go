package resp

import "time"

// BrokerageUserResp 分销用户 Response
type BrokerageUserResp struct {
	ID               int64      `json:"id"`
	BindUserID       int64      `json:"bindUserId"`
	BindUserTime     *time.Time `json:"bindUserTime"`
	BrokerageEnabled bool       `json:"brokerageEnabled"`
	BrokerageTime    *time.Time `json:"brokerageTime"`
	Price            int        `json:"price"` // BrokeragePrice
	FrozenPrice      int        `json:"frozenPrice"`
	CreateTime       time.Time  `json:"createTime"`

	// User Info
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`

	// Brokerage Info (Aggregated)
	BrokerageUserCount  int `json:"brokerageUserCount"`
	BrokerageOrderCount int `json:"brokerageOrderCount"`
	BrokerageOrderPrice int `json:"brokerageOrderPrice"`

	// Withdraw Info
	WithdrawPrice int `json:"withdrawPrice"`
	WithdrawCount int `json:"withdrawCount"`
}

// BrokerageRecordResp 分销记录 Response
type BrokerageRecordResp struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"userId"`
	BizID           string     `json:"bizId"`
	BizType         int        `json:"bizType"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Price           int        `json:"price"`
	TotalPrice      int        `json:"totalPrice"`
	Status          int        `json:"status"`
	FrozenDays      int        `json:"frozenDays"`
	UnfreezeTime    *time.Time `json:"unfreezeTime"`
	SourceUserLevel int        `json:"sourceUserLevel"`
	SourceUserID    int64      `json:"sourceUserId"`
	CreateTime      time.Time  `json:"createTime"`

	// User Info
	BrokerageUserResp
}

// BrokerageWithdrawResp 分销提现 Response
type BrokerageWithdrawResp struct {
	ID                  int64      `json:"id"`
	UserID              int64      `json:"userId"`
	Price               int        `json:"price"`
	FeePrice            int        `json:"feePrice"`
	TotalPrice          int        `json:"totalPrice"`
	Type                int        `json:"type"`
	UserName            string     `json:"userName"`
	UserAccount         string     `json:"userAccount"`
	QRCodeUrl           string     `json:"qrCodeUrl"`
	BankName            string     `json:"bankName"`
	BankAddress         string     `json:"bankAddress"`
	Status              int        `json:"status"`
	AuditReason         string     `json:"auditReason"`
	AuditTime           *time.Time `json:"auditTime"`
	Remark              string     `json:"remark"`
	PayTransferID       int64      `json:"payTransferId"`
	TransferChannelCode string     `json:"transferChannelCode"`
	TransferTime        *time.Time `json:"transferTime"`
	TransferErrorMsg    string     `json:"transferErrorMsg"`
	CreateTime          time.Time  `json:"createTime"`

	// User Info
	BrokerageUserResp
}
