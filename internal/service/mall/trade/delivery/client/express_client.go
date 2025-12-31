package client

type ExpressClient interface {
	// GetExpressTrackList 获得物流轨迹
	GetExpressTrackList(req *ExpressTrackQueryReqDTO) ([]ExpressTrackRespDTO, error)
}

// ExpressTrackQueryReqDTO 快递查询请求 DTO
type ExpressTrackQueryReqDTO struct {
	ExpressCode string // 快递公司编码
	LogisticsNo string // 快递单号
	Phone       string // 手机号
}

// ExpressTrackRespDTO 快递查询响应 DTO
type ExpressTrackRespDTO struct {
	Time    string `json:"time"`
	Context string `json:"context"`
}
