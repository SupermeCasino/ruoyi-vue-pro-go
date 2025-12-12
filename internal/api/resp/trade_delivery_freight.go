package resp

import "time"

// DeliveryFreightTemplateResp 运费模板 Response
type DeliveryFreightTemplateResp struct {
	ID         int64                               `json:"id"`
	Name       string                              `json:"name"`
	ChargeMode int                                 `json:"chargeMode"`
	Sort       int                                 `json:"sort"`
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
	AreaIDs   []int `json:"areaIds"`
	FreePrice int   `json:"freePrice"`
	FreeCount int   `json:"freeCount"`
}

// SimpleDeliveryFreightTemplateResp 运费模板精简 Response
type SimpleDeliveryFreightTemplateResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
