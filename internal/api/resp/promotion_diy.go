package resp

import "time"

// DiyTemplateResp 装修模板响应
type DiyTemplateResp struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	PreviewPicUrls []string   `json:"previewPicUrls"`
	Property       string     `json:"property"`
	Remark         string     `json:"remark"`
	Used           bool       `json:"used"`
	UsedTime       *time.Time `json:"usedTime"`
	CreateTime     time.Time  `json:"createTime"`
}

// DiyPageResp 装修页面响应
type DiyPageResp struct {
	ID             int64     `json:"id"`
	TemplateID     int64     `json:"templateId"`
	Name           string    `json:"name"`
	Remark         string    `json:"remark"`
	PreviewPicUrls []string  `json:"previewPicUrls"`
	Property       string    `json:"property"`
	CreateTime     time.Time `json:"createTime"`
}

type AppDiyTemplatePropertyResp struct {
	DiyTemplateResp
	Home string `json:"home"`
	User string `json:"user"`
}
