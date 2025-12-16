package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// DiyTemplateCreateReq 装修模板创建请求
type DiyTemplateCreateReq struct {
	Name         string `json:"name" binding:"required"`
	CoverImage   string `json:"coverImage"`
	PreviewImage string `json:"previewImage"`
	Status       int    `json:"status" binding:"required"` // 0-开启 1-关闭
	Property     string `json:"property"`
	Sort         int    `json:"sort"`
	Remark       string `json:"remark"`
}

// DiyTemplateUpdateReq 装修模板更新请求
type DiyTemplateUpdateReq struct {
	ID           int64  `json:"id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	CoverImage   string `json:"coverImage"`
	PreviewImage string `json:"previewImage"`
	Status       int    `json:"status" binding:"required"`
	Property     string `json:"property"`
	Sort         int    `json:"sort"`
	Remark       string `json:"remark"`
}

// DiyTemplatePageReq 装修模板分页请求
type DiyTemplatePageReq struct {
	pagination.PageParam
	Name string `form:"name"`
}

// DiyPageCreateReq 装修页面创建请求
type DiyPageCreateReq struct {
	TemplateID int64  `json:"templateId" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Remark     string `json:"remark"`
	Status     int    `json:"status" binding:"required"` // 0-开启 1-关闭
	Property   string `json:"property"`
}

// DiyPageUpdateReq 装修页面更新请求
type DiyPageUpdateReq struct {
	ID         int64  `json:"id" binding:"required"`
	TemplateID int64  `json:"templateId" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Remark     string `json:"remark"`
	Status     int    `json:"status" binding:"required"`
	Property   string `json:"property"`
}

// DiyPagePageReq 装修页面分页请求
type DiyPagePageReq struct {
	pagination.PageParam
	Name   string `form:"name"`
	Status *int   `form:"status"`
}
