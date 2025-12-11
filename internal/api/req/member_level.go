package req

import "backend-go/internal/pkg/core"

type MemberLevelCreateReq struct {
	Name            string `json:"name" binding:"required"`
	Level           int    `json:"level" binding:"required"`
	Experience      int    `json:"experience" binding:"required"`
	DiscountPercent int    `json:"discountPercent" binding:"required"`
	Icon            string `json:"icon"`
	BackgroundURL   string `json:"backgroundUrl"`
	Status          int    `json:"status" binding:"required"`
	Remark          string `json:"remark"`
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
	Remark          string `json:"remark"`
}

type MemberLevelPageReq struct {
	core.PageParam
	Name   string `form:"name"`
	Status *int   `form:"status"`
}
