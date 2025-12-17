package resp

import (
	"time"

	"gorm.io/datatypes"
)

// DiyTemplateResp 装修模板响应 (Admin 端使用)
type DiyTemplateResp struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	PreviewPicUrls []string       `json:"previewPicUrls"`
	Property       datatypes.JSON `json:"property"`
	Remark         string         `json:"remark"`
	Used           bool           `json:"used"`
	UsedTime       *time.Time     `json:"usedTime"`
	CreateTime     time.Time      `json:"createTime"`
}

// DiyPageResp 装修页面响应
type DiyPageResp struct {
	ID             int64          `json:"id"`
	TemplateID     int64          `json:"templateId"`
	Name           string         `json:"name"`
	Remark         string         `json:"remark"`
	PreviewPicUrls []string       `json:"previewPicUrls"`
	Property       datatypes.JSON `json:"property"`
	CreateTime     time.Time      `json:"createTime"`
}

// AppDiyTemplatePropertyResp 用户App - 装修模板属性响应 (严格对齐 Java: AppDiyTemplatePropertyRespVO)
// Java 使用 @JsonRawValue 注解让 property/home/user 作为原始 JSON 输出
type AppDiyTemplatePropertyResp struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Property datatypes.JSON `json:"property"` // @JsonRawValue 效果
	Home     datatypes.JSON `json:"home"`     // @JsonRawValue 效果
	User     datatypes.JSON `json:"user"`     // @JsonRawValue 效果
}
