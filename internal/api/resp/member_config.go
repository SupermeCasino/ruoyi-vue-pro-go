package resp

type MemberConfigResp struct {
	ID                        int64 `json:"id"`
	PointTradeDeductEnable    int   `json:"pointTradeDeductEnable"`
	PointTradeDeductUnitPrice int   `json:"pointTradeDeductUnitPrice"`
	PointTradeDeductMaxPrice  int   `json:"pointTradeDeductMaxPrice"`
	PointTradeGivePoint       int   `json:"pointTradeGivePoint"`
}
