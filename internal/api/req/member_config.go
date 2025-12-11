package req

type MemberConfigSaveReq struct {
	PointTradeDeductEnable    int `json:"pointTradeDeductEnable"`    // 积分抵扣开关
	PointTradeDeductUnitPrice int `json:"pointTradeDeductUnitPrice"` // 积分抵扣单位价格
	PointTradeDeductMaxPrice  int `json:"pointTradeDeductMaxPrice"`  // 积分抵扣最大值
	PointTradeGivePoint       int `json:"pointTradeGivePoint"`       // 1 元赠送多少分
}
