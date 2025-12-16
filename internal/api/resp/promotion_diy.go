package resp

import "time"

// DiyTemplateResp 装修模板响应
type DiyTemplateResp struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	CoverImage   string    `json:"coverImage"`
	PreviewImage string    `json:"previewImage"`
	Property     string    `json:"property"`
	Sort         int       `json:"sort"`
	Remark       string    `json:"remark"`
	CreateTime   time.Time `json:"createTime"`
}

// DiyPageResp 装修页面响应
type DiyPageResp struct {
	ID         int64     `json:"id"`
	TemplateID int64     `json:"templateId"`
	Name       string    `json:"name"`
	Remark     string    `json:"remark"`
	Status     int       `json:"status"`
	Property   string    `json:"property"`
	CreateTime time.Time `json:"createTime"`
}

type AppDiyTemplatePropertyResp struct {
	DiyTemplateResp
	Home string `json:"home"`
	User string `json:"user"`
}
