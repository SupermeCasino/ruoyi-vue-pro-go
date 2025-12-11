package resp

import "time"

type DiscountActivityRespVO struct {
	ID         int64                    `json:"id"`
	Name       string                   `json:"name"`
	Status     int                      `json:"status"`
	StartTime  time.Time                `json:"startTime"`
	EndTime    time.Time                `json:"endTime"`
	Remark     string                   `json:"remark"`
	CreateTime time.Time                `json:"createTime"`
	Products   []*DiscountProductRespVO `json:"products"`
}

type DiscountProductRespVO struct {
	ID              int64 `json:"id"`
	ActivityID      int64 `json:"activityId"`
	SpuID           int64 `json:"spuId"`
	SkuID           int64 `json:"skuId"`
	DiscountType    int   `json:"discountType"`
	DiscountPercent int   `json:"discountPercent"`
	DiscountPrice   int   `json:"discountPrice"`
}
