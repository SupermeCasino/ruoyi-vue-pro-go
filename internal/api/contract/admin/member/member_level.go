package member

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type MemberLevelResp struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Level           int       `json:"level"`
	Experience      int       `json:"experience"`
	DiscountPercent int       `json:"discountPercent"`
	Icon            string    `json:"icon"`
	BackgroundURL   string    `json:"backgroundUrl"`
	CreateTime      time.Time `json:"createTime"`
}

type MemberLevelCreateReq struct {
	Name            string `json:"name" binding:"required"`
	Level           int    `json:"level" binding:"required"`
	Experience      int    `json:"experience" binding:"required"`
	DiscountPercent int    `json:"discountPercent" binding:"required"`
	Icon            string `json:"icon"`
	BackgroundURL   string `json:"backgroundUrl"`
	Remark          string `json:"remark"`
	Status          int    `json:"status"`
}

type MemberLevelUpdateReq struct {
	ID              int64  `json:"id" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Level           int    `json:"level" binding:"required"`
	Experience      int    `json:"experience" binding:"required"`
	DiscountPercent int    `json:"discountPercent" binding:"required"`
	Icon            string `json:"icon"`
	BackgroundURL   string `json:"backgroundUrl"`
	Remark          string `json:"remark"`
	Status          int    `json:"status"`
}

type MemberLevelPageReq struct {
	pagination.PageParam
	Name   string `form:"name"`
	Status *int   `form:"status"`
}
