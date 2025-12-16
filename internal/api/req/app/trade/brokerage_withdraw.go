package trade

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

type AppBrokerageWithdrawPageReqVO struct {
	core.PageParam
	Type   int `form:"type"`
	Status int `form:"status"`
}

type AppBrokerageWithdrawCreateReqVO struct {
	Type        int    `json:"type"`        // 提现类型
	Price       int    `json:"price"`       // 提现金额
	Name        string `json:"name"`        // 真实姓名 (Bank/Alipay)
	Account     string `json:"account"`     // 账号 (Bank/Alipay)
	BankName    string `json:"bankName"`    // 银行名称 (Bank)
	BankAddress string `json:"bankAddress"` // 开户地址 (Bank)
	QrCodeUrl   string `json:"qrCodeUrl"`   // 收款码 (Wechat)
	Code        string `json:"code"`        // 微信渠道需要
}
