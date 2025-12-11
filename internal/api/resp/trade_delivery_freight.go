package resp

import "time"

// DeliveryFreightTemplateResp 运费模板 Response
type DeliveryFreightTemplateResp struct {
	ID         int64                               `json:"id"`
	Name       string                              `json:"name"`
	Type       int                                 `json:"type"`
	ChargeMode int                                 `json:"chargeMode"`
	Sort       int                                 `json:"sort"`
	Status     int                                 `json:"status"`
	Remark     string                              `json:"remark"`
	CreateTime time.Time                           `json:"createTime"`
	Charges    []DeliveryFreightTemplateChargeResp `json:"charges"`
	Frees      []DeliveryFreightTemplateFreeResp   `json:"frees"`
}

type DeliveryFreightTemplateChargeResp struct {
	AreaIDs    []int   `json:"areaIds"` // 返回数组
	StartCount float64 `json:"startCount"`
	StartPrice int     `json:"startPrice"`
	ExtraCount float64 `json:"extraCount"`
	ExtraPrice int     `json:"extraPrice"`
}

type DeliveryFreightTemplateFreeResp struct {
	AreaIDs   []int   `json:"areaIds"`
	FreePrice int     `json:"freePrice"`
	FreeCount float64 `json:"freeCount"`
}

// SimpleDeliveryFreightTemplateResp 运费模板精简 Response
type SimpleDeliveryFreightTemplateResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
