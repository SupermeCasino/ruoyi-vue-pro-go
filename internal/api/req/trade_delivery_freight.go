package req

// DeliveryFreightTemplatePageReq 运费模板分页 Request
type DeliveryFreightTemplatePageReq struct {
	PageNo   int    `form:"pageNo"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
	Status   *int   `form:"status"` // 状态
}

// DeliveryFreightTemplateSaveReq 运费模板保存 Request
type DeliveryFreightTemplateSaveReq struct {
	ID         int64                                  `json:"id"`
	Name       string                                 `json:"name"`
	Type       int                                    `json:"type"`
	ChargeMode int                                    `json:"chargeMode"`
	Sort       int                                    `json:"sort"`
	Status     int                                    `json:"status"`
	Remark     string                                 `json:"remark"`
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
	AreaIDs   []int   `json:"areaIds"`
	FreePrice int     `json:"freePrice"`
	FreeCount float64 `json:"freeCount"`
}
