package trade

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

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
	ID            *int64          `json:"id"`
	Name          string          `json:"name" binding:"required"`
	Introduction  string          `json:"introduction"`
	Phone         string          `json:"phone" binding:"required"`
	AreaID        int             `json:"areaId" binding:"required"`
	DetailAddress string          `json:"detailAddress" binding:"required"`
	Logo          string          `json:"logo" binding:"required"`
	OpeningTime   model.TimeOfDay `json:"openingTime" binding:"required"`
	ClosingTime   model.TimeOfDay `json:"closingTime" binding:"required"`
	Latitude      float64         `json:"latitude" binding:"required"`
	Longitude     float64         `json:"longitude" binding:"required"`
	Status        int             `json:"status" binding:"required"`
	Sort          int             `json:"sort"`
}

type ExpressTrackRespVO struct {
	Time    string `json:"time"`
	Content string `json:"content"`
}

// ========== Response DTOs ==========

// DeliveryExpressResp 物流公司 Response
type DeliveryExpressResp struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Logo       string    `json:"logo"`
	Sort       int       `json:"sort"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"createTime"`
}

// DeliveryExpressExcelVO 物流公司 Excel 导出 VO
type DeliveryExpressExcelVO struct {
	ID         int64     `json:"id" excel:"编号"`
	Code       string    `json:"code" excel:"物流代码"`
	Name       string    `json:"name" excel:"物流名称"`
	Logo       string    `json:"logo" excel:"物流图标"`
	Sort       int       `json:"sort" excel:"排序"`
	Status     int       `json:"status" excel:"状态"`
	CreateTime time.Time `json:"createTime" excel:"创建时间"`
}

// DeliveryPickUpStoreResp 自提门店 Response
type DeliveryPickUpStoreResp struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	Introduction  string          `json:"introduction"`
	Phone         string          `json:"phone"`
	AreaID        int             `json:"areaId"`
	DetailAddress string          `json:"detailAddress"`
	Logo          string          `json:"logo"`
	OpeningTime   model.TimeOfDay `json:"openingTime"`
	ClosingTime   model.TimeOfDay `json:"closingTime"`
	Latitude      float64         `json:"latitude"`
	Longitude     float64         `json:"longitude"`
	Status        int             `json:"status"`
	CreateTime    time.Time       `json:"createTime"`
}
