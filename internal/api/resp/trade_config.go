package resp

type TradeConfigResp struct {
	ID                          int64    `json:"id"`
	AppID                       int64    `json:"appId"`                       // 支付应用 ID
	AfterSaleDeadlineDays       int      `json:"afterSaleDeadlineDays"`       // 售后期限(天)
	PayTimeoutMinutes           int      `json:"payTimeoutMinutes"`           // 支付超时(分钟)
	AutoReceiveDays             int      `json:"autoReceiveDays"`             // 自动收货(天)
	AutoCommentDays             int      `json:"autoCommentDays"`             // 自动好评(天)
	BrokerageWithdrawMinPrice   int      `json:"brokerageWithdrawMinPrice"`   // 提现最低金额
	BrokerageWithdrawFeePercent int      `json:"brokerageWithdrawFeePercent"` // 提现手续费百分比
	BrokerageEnabled            bool     `json:"brokerageEnabled"`            // 是否开启分销
	BrokerageFrozenDays         int      `json:"brokerageFrozenDays"`         // 分销佣金冻结时间（天）
	BrokerageFirstPercent       int      `json:"brokerageFirstPercent"`       // 一级分销比例
	BrokerageSecondPercent      int      `json:"brokerageSecondPercent"`      // 二级分销比例
	BrokerageEnabledCondition   int      `json:"brokerageEnabledCondition"`   // 分销资格启用条件 1:人人分销 2:仅指定用户
	BrokerageBindMode           int      `json:"brokerageBindMode"`           // 分销关系绑定模式 1:首次绑定 2:注册绑定 3:覆盖绑定
	BrokeragePosterUrls         []string `json:"brokeragePosterUrls"`         // 分销海报图
}
