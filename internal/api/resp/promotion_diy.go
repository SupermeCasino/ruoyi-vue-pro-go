package resp

import (
	"time"

	"gorm.io/datatypes"
)

// DiyTemplateBase 装修模板 Base
type DiyTemplateBase struct {
	Name           string   `json:"name"`
	Remark         string   `json:"remark"`
	PreviewPicUrls []string `json:"previewPicUrls"`
}

// DiyTemplateResp 装修模板响应 (Admin List/Get)
type DiyTemplateResp struct {
	DiyTemplateBase
	ID         int64      `json:"id"`
	Used       bool       `json:"used"`
	UsedTime   *time.Time `json:"usedTime"`
	CreateTime time.Time  `json:"createTime"`
}

// DiyTemplatePropertyResp 装修模板属性响应 (Admin GetProperty)
type DiyTemplatePropertyResp struct {
	DiyTemplateBase
	ID       int64                 `json:"id"`
	Property string                `json:"property"` // JSON String
	Pages    []DiyPagePropertyResp `json:"pages"`
}

// DiyPageBase 装修页面 Base
type DiyPageBase struct {
	TemplateID     int64    `json:"templateId"`
	Name           string   `json:"name"`
	Remark         string   `json:"remark"`
	PreviewPicUrls []string `json:"previewPicUrls"`
}

// DiyPageResp 装修页面响应 (Admin List/Get)
type DiyPageResp struct {
	DiyPageBase
	ID         int64     `json:"id"`
	CreateTime time.Time `json:"createTime"`
}

// DiyPagePropertyResp 装修页面属性响应 (Admin GetProperty)
type DiyPagePropertyResp struct {
	DiyPageBase
	ID       int64  `json:"id"`
	Property string `json:"property"` // JSON String
}

// AppDiyTemplatePropertyResp 用户App - 装修模板属性响应
type AppDiyTemplatePropertyResp struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Property datatypes.JSON `json:"property"`
	Home     datatypes.JSON `json:"home"`
	User     datatypes.JSON `json:"user"`
}

// AppDiyPagePropertyResp 用户App - 装修页面属性响应 (新增)
type AppDiyPagePropertyResp struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Property datatypes.JSON `json:"property"`
}
