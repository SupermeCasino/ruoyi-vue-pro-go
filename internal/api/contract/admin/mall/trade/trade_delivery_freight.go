package trade

import "time"

// DeliveryFreightTemplatePageReq 运费模板分页 Request
type DeliveryFreightTemplatePageReq struct {
	PageNo   int    `form:"pageNo"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
}

// DeliveryFreightTemplateSaveReq 运费模板保存 Request
type DeliveryFreightTemplateSaveReq struct {
	ID         int64                                  `json:"id"`
	Name       string                                 `json:"name"`
	ChargeMode int                                    `json:"chargeMode"`
	Sort       int                                    `json:"sort"`
	Charges    []DeliveryFreightTemplateChargeSaveReq `json:"charges"`
	Frees      []DeliveryFreightTemplateFreeSaveReq   `json:"frees"`
}

type DeliveryFreightTemplateChargeSaveReq struct {
	AreaIDs    []int   `json:"areaIds"` // 接收数组，转换成逗号分隔字符串
	StartCount float64 `json:"startCount"`
	StartPrice int     `json:"startPrice"`
	ExtraCount float64 `json:"extraCount"`
	ExtraPrice int     `json:"extraPrice"`
}

type DeliveryFreightTemplateFreeSaveReq struct {
	AreaIDs   []int `json:"areaIds"`
	FreePrice int   `json:"freePrice"`
	FreeCount int   `json:"freeCount"`
}

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
	AreaIDs    []int   `json:"areaIds"`
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

type SimpleDeliveryFreightTemplateResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
