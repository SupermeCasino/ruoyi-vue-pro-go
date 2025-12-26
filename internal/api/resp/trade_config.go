package resp

// TradeConfigResp Admin 端 - 交易配置响应 (对齐 Java: TradeConfigRespVO)
type TradeConfigResp struct {
	ID                          int64    `json:"id"`                          // 自增主键
	AppID                       int64    `json:"appId"`                       // 支付应用 ID
	AfterSaleDeadlineDays       int      `json:"afterSaleDeadlineDays"`       // 售后期限(天)
	PayTimeoutMinutes           int      `json:"payTimeoutMinutes"`           // 支付超时(分钟)
	AutoReceiveDays             int      `json:"autoReceiveDays"`             // 自动收货(天)
	AutoCommentDays             int      `json:"autoCommentDays"`             // 自动好评(天)
	AfterSaleRefundReasons      []string `json:"afterSaleRefundReasons"`      // 售后的退款理由
	AfterSaleReturnReasons      []string `json:"afterSaleReturnReasons"`      // 售后的退货理由
	DeliveryExpressFreeEnabled  bool     `json:"deliveryExpressFreeEnabled"`  // 是否启用全场包邮
	DeliveryExpressFreePrice    int      `json:"deliveryExpressFreePrice"`    // 全场包邮的最小金额
	DeliveryPickUpEnabled       bool     `json:"deliveryPickUpEnabled"`       // 是否开启自提
	BrokerageWithdrawMinPrice   int      `json:"brokerageWithdrawMinPrice"`   // 提现最低金额
	BrokerageWithdrawFeePercent int      `json:"brokerageWithdrawFeePercent"` // 提现手续费百分比
	BrokerageEnabled            bool     `json:"brokerageEnabled"`            // 是否开启分销
	BrokerageFrozenDays         int      `json:"brokerageFrozenDays"`         // 分销佣金冻结时间（天）
	BrokerageFirstPercent       int      `json:"brokerageFirstPercent"`       // 一级分销比例
	BrokerageSecondPercent      int      `json:"brokerageSecondPercent"`      // 二级分销比例
	BrokerageEnabledCondition   int      `json:"brokerageEnabledCondition"`   // 分销资格启用条件 1:人人分销 2:仅指定用户
	BrokerageBindMode           int      `json:"brokerageBindMode"`           // 分销关系绑定模式 1:首次绑定 2:注册绑定 3:覆盖绑定
	BrokeragePosterUrls         []string `json:"brokeragePosterUrls"`         // 分销海报图
	BrokerageWithdrawTypes      []int    `json:"brokerageWithdrawTypes"`      // 提现方式列表
	TencentLbsKey               string   `json:"tencentLbsKey"`               // 腾讯地图 Key
}

// AppTradeConfigResp App 端 - 交易配置响应 (对齐 Java: AppTradeConfigRespVO)
type AppTradeConfigResp struct {
	TencentLbsKey             string   `json:"tencentLbsKey"`             // 腾讯地图 Key
	DeliveryPickUpEnabled     bool     `json:"deliveryPickUpEnabled"`     // 是否启用自提
	AfterSaleRefundReasons    []string `json:"afterSaleRefundReasons"`    // 售后退款原因列表
	AfterSaleReturnReasons    []string `json:"afterSaleReturnReasons"`    // 售后退货原因列表
	BrokeragePosterUrls       []string `json:"brokeragePosterUrls"`       // 分销海报图
	BrokerageFrozenDays       int      `json:"brokerageFrozenDays"`       // 分销佣金冻结时间（天）
	BrokerageWithdrawMinPrice int      `json:"brokerageWithdrawMinPrice"` // 提现最低金额
	BrokerageWithdrawTypes    []int    `json:"brokerageWithdrawTypes"`    // 提现方式列表
}
