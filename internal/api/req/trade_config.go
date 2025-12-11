package req

// TradeConfigSaveReq 交易配置 - 保存 Request
type TradeConfigSaveReq struct {
	AfterSaleDeadlineDays       *int     `json:"afterSaleDeadlineDays" binding:"required"` // 售后期限(天)
	PayTimeoutMinutes           *int     `json:"payTimeoutMinutes" binding:"required"`     // 支付超时(分钟)
	AutoReceiveDays             *int     `json:"autoReceiveDays" binding:"required"`       // 自动收货(天)
	AutoCommentDays             *int     `json:"autoCommentDays" binding:"required"`       // 自动好评(天)
	BrokerageWithdrawMinPrice   *int     `json:"brokerageWithdrawMinPrice"`                // 提现最低金额
	BrokerageWithdrawFeePercent *int     `json:"brokerageWithdrawFeePercent"`              // 提现手续费百分比
	BrokerageEnabled            *bool    `json:"brokerageEnabled"`                         // 是否开启分销
	BrokerageFrozenDays         *int     `json:"brokerageFrozenDays"`                      // 分销佣金冻结时间（天）
	BrokerageFirstPercent       *int     `json:"brokerageFirstPercent"`                    // 一级分销比例
	BrokerageSecondPercent      *int     `json:"brokerageSecondPercent"`                   // 二级分销比例
	BrokeragePosterUrls         []string `json:"brokeragePosterUrls"`                      // 分销海报图
}
