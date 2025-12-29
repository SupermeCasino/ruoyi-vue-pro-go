package resp

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

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

// DeliveryPickUpStoreResp 自提门店 Response
type DeliveryPickUpStoreResp struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	Introduction  string           `json:"introduction"`
	Phone         string           `json:"phone"`
	AreaID        int              `json:"areaId"`
	DetailAddress string           `json:"detailAddress"`
	Logo          string           `json:"logo"`
	OpeningTime   model.TimeOfDay  `json:"openingTime"`
	ClosingTime   model.TimeOfDay  `json:"closingTime"`
	Latitude      float64          `json:"latitude"`
	Longitude     float64          `json:"longitude"`
	Status        int              `json:"status"`
	Sort          int              `json:"sort"`
	CreateTime    time.Time        `json:"createTime"`
	VerifyUsers   []UserSimpleResp `json:"verifyUsers"` // 核销用户数组
}

// UserSimpleResp 用户精简信息
type UserSimpleResp struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ExpressTrackRespVO 物流轨迹 Response
type ExpressTrackRespVO struct {
	Time    string `json:"time"`
	Content string `json:"content"`
}

// DeliveryExpressExcelVO 物流公司导出 Response
type DeliveryExpressExcelVO struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Logo       string    `json:"logo"`
	Sort       int       `json:"sort"`
	Status     int       `json:"status" dict_type:"common_status"`
	CreateTime time.Time `json:"createTime"`
}
