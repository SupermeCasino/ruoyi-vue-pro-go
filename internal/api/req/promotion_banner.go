package req

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

// PromotionBannerCreateReq 创建 Request
type PromotionBannerCreateReq struct {
	Title    string `json:"title" binding:"required"`
	PicURL   string `json:"picUrl" binding:"required"`
	Url      string `json:"url" binding:"required"`
	Status   int    `json:"status"` // Default 0
	Sort     int    `json:"sort"`
	Position int    `json:"position" binding:"required"`
	Memo     string `json:"memo"`
}

// PromotionBannerUpdateReq 更新 Request
type PromotionBannerUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	PicURL   string `json:"picUrl" binding:"required"`
	Url      string `json:"url" binding:"required"`
	Status   int    `json:"status"`
	Sort     int    `json:"sort"`
	Position int    `json:"position" binding:"required"`
	Memo     string `json:"memo"`
}

// PromotionBannerPageReq 分页 Request
type PromotionBannerPageReq struct {
	core.PageParam
	Title  string `form:"title"`
	Status *int   `form:"status"`
}

// AppBannerListReq App 列表 Request
type AppBannerListReq struct {
	Position int `form:"position" binding:"required"`
}
