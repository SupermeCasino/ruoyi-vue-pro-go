package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type MemberLevelCreateReq struct {
	Name            string `json:"name" binding:"required"`
	Level           int    `json:"level" binding:"required"`
	Experience      int    `json:"experience" binding:"required"`
	DiscountPercent int    `json:"discountPercent" binding:"required"`
	Icon            string `json:"icon"`
	BackgroundURL   string `json:"backgroundUrl"`
	Status          int    `json:"status" binding:"required"`
}

type MemberLevelUpdateReq struct {
	ID              int64  `json:"id" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Level           int    `json:"level" binding:"required"`
	Experience      int    `json:"experience" binding:"required"`
	DiscountPercent int    `json:"discountPercent" binding:"required"`
	Icon            string `json:"icon"`
	BackgroundURL   string `json:"backgroundUrl"`
	Status          int    `json:"status" binding:"required"`
}

type MemberLevelPageReq struct {
	pagination.PageParam
	Name   string `form:"name"`
	Status *int   `form:"status"`
}
