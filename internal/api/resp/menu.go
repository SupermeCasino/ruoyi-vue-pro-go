package resp

import "time"

// MenuResp 菜单信息响应
type MenuResp struct {
	ID            int64     `json:"id"`
	ParentID      int64     `json:"parentId"`
	Name          string    `json:"name"`
	Type          int32     `json:"type"`
	Sort          int32     `json:"sort"`
	Path          string    `json:"path"`
	Icon          string    `json:"icon"`
	Component     string    `json:"component"`
	ComponentName string    `json:"componentName"`
	Permission    string    `json:"permission"`
	Status        int32     `json:"status"`
	Visible       bool      `json:"visible"`
	KeepAlive     bool      `json:"keepAlive"`
	AlwaysShow    bool      `json:"alwaysShow"`
	CreateTime    time.Time `json:"createTime"`
}

// MenuSimpleResp 菜单精简响应
type MenuSimpleResp struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parentId"`
	Name     string `json:"name"`
	Type     int32  `json:"type"`
}
