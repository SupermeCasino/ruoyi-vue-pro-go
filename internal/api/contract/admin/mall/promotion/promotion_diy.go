package promotion

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

// DiyTemplateCreateReq 装修模板创建请求
type DiyTemplateCreateReq struct {
	Name           string         `json:"name" binding:"required"`
	PreviewPicUrls []string       `json:"previewPicUrls"`
	Property       datatypes.JSON `json:"property"`
	Remark         string         `json:"remark"`
}

// DiyTemplateUpdateReq 装修模板更新请求
type DiyTemplateUpdateReq struct {
	ID             int64          `json:"id" binding:"required"`
	Name           string         `json:"name" binding:"required"`
	PreviewPicUrls []string       `json:"previewPicUrls"`
	Property       datatypes.JSON `json:"property"`
	Remark         string         `json:"remark"`
}

// DiyTemplatePageReq 装修模板分页请求
type DiyTemplatePageReq struct {
	pagination.PageParam
	Name       string   `form:"name"`
	CreateTime []string `form:"createTime[]"` // 时间范围
}

// DiyPageCreateReq 装修页面创建请求
type DiyPageCreateReq struct {
	TemplateID     int64          `json:"templateId" binding:"required"`
	Name           string         `json:"name" binding:"required"`
	Remark         string         `json:"remark"`
	PreviewPicUrls []string       `json:"previewPicUrls"`
	Property       datatypes.JSON `json:"property"`
}

// DiyPageUpdateReq 装修页面更新请求
type DiyPageUpdateReq struct {
	ID             int64          `json:"id" binding:"required"`
	TemplateID     int64          `json:"templateId" binding:"required"`
	Name           string         `json:"name" binding:"required"`
	Remark         string         `json:"remark"`
	PreviewPicUrls []string       `json:"previewPicUrls"`
	Property       datatypes.JSON `json:"property"`
}

// DiyPagePageReq 装修页面分页请求
type DiyPagePageReq struct {
	pagination.PageParam
	Name       string   `form:"name"`
	CreateTime []string `form:"createTime[]"` // 时间范围
}

// DiyPagePropertyUpdateReq 装修页面属性更新请求
type DiyPagePropertyUpdateReq struct {
	ID       int64          `json:"id" binding:"required"`
	Property datatypes.JSON `json:"property" binding:"required"`
}

// DiyPageBase DIY 页面基础 VO
type DiyPageBase struct {
	TemplateID     int64    `json:"templateId"`
	Name           string   `json:"name"`
	Remark         string   `json:"remark"`
	PreviewPicUrls []string `json:"previewPicUrls"`
}

// DiyPageResp DIY 页面 Response VO
type DiyPageResp struct {
	DiyPageBase
	ID         int64     `json:"id"`
	CreateTime time.Time `json:"createTime"`
}

// DiyPagePropertyResp DIY 页面属性 Response VO
type DiyPagePropertyResp struct {
	DiyPageBase
	ID       int64  `json:"id"`
	Property string `json:"property"`
}

// DiyTemplateBase 装修模板基础 VO
type DiyTemplateBase struct {
	Name           string   `json:"name"`
	PreviewPicUrls []string `json:"previewPicUrls"`
	Remark         string   `json:"remark"`
}

// DiyTemplateResp 装修模板 Response VO
type DiyTemplateResp struct {
	DiyTemplateBase
	ID         int64      `json:"id"`
	Used       bool       `json:"used"`
	UsedTime   *time.Time `json:"usedTime"`
	Property   string     `json:"property"`
	CreateTime time.Time  `json:"createTime"`
}

// DiyTemplatePropertyResp 装修模板属性 Response VO
type DiyTemplatePropertyResp struct {
	DiyTemplateBase
	ID       int64                 `json:"id"`
	Property string                `json:"property"`
	Pages    []DiyPagePropertyResp `json:"pages"`
}

// DiyTemplatePropertyUpdateReq 装修模板属性更新请求
type DiyTemplatePropertyUpdateReq struct {
	ID       int64          `json:"id" binding:"required"`
	Property datatypes.JSON `json:"property" binding:"required"`
}
