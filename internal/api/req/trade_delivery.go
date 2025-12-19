package req

import "time"

// DeliveryExpressPageReq 物流公司分页 Request
type DeliveryExpressPageReq struct {
	PageNo   int    `form:"pageNo"`
	PageSize int    `form:"pageSize"`
	Code     string `form:"code"`
	Name     string `form:"name"`
	Status   *int   `form:"status"`
}

// DeliveryExpressSaveReq 物流公司保存 Request
type DeliveryExpressSaveReq struct {
	ID     *int64 `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

// DeliveryPickUpStorePageReq 自提门店分页 Request
type DeliveryPickUpStorePageReq struct {
	PageNo   int    `form:"pageNo"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
	Phone    string `form:"phone"`
	Status   *int   `form:"status"`
}

// DeliveryPickUpStoreSaveReq 自提门店保存 Request
type DeliveryPickUpStoreSaveReq struct {
	ID            *int64     `json:"id"`
	Name          string     `json:"name" binding:"required"`
	Introduction  string     `json:"introduction"`
	Phone         string     `json:"phone" binding:"required"`
	AreaID        int        `json:"areaId" binding:"required"`
	DetailAddress string     `json:"detailAddress" binding:"required"`
	Logo          string     `json:"logo" binding:"required"`
	OpeningTime   *time.Time `json:"openingTime" binding:"required"`
	ClosingTime   *time.Time `json:"closingTime" binding:"required"`
	Latitude      float64    `json:"latitude" binding:"required"`
	Longitude     float64    `json:"longitude" binding:"required"`
	Status        int        `json:"status" binding:"required"`
	Sort          int        `json:"sort"`
}
