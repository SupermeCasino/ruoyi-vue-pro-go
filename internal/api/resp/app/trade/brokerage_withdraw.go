package trade

import "time"

type AppBrokerageWithdrawRespVO struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"userId"`
	Price       int        `json:"price"`
	FeePrice    int        `json:"feePrice"`
	TotalPrice  int        `json:"totalPrice"`
	Type        int        `json:"type"`
	Name        string     `json:"name"`
	Account     string     `json:"account"`
	BankName    string     `json:"bankName"`
	Status      int        `json:"status"`
	AuditReason string     `json:"auditReason"`
	AuditTime   *time.Time `json:"auditTime"`
	Remark      string     `json:"remark"`
	CreatedAt   time.Time  `json:"createTime"`
	TypeName    string     `json:"typeName"`
	StatusName  string     `json:"statusName"`

	// Wechat specific
	TransferChannelPackageInfo string `json:"transferChannelPackageInfo,omitempty"`
	TransferChannelMchId       string `json:"transferChannelMchId,omitempty"`
}
