package system

import (
	"encoding/json"
	"time"
)

// FileConfigRespVO 文件配置 Response
type FileConfigRespVO struct {
	ID         int64            `json:"id"`
	Name       string           `json:"name"`
	Storage    int32            `json:"storage"`
	Master     bool             `json:"master"`
	Config     *json.RawMessage `json:"config"`
	Remark     string           `json:"remark"`
	CreateTime time.Time        `json:"createTime"`
}

// FileRespVO 文件 Response
type FileRespVO struct {
	ID         int64     `json:"id"`
	ConfigId   int64     `json:"configId"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Url        string    `json:"url"`
	Type       string    `json:"type"`
	Size       int       `json:"size"`
	CreateTime time.Time `json:"createTime"`
}
