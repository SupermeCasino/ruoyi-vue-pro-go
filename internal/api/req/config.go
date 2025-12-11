package req

import "backend-go/internal/pkg/core"

type ConfigSaveReq struct {
	ID       int64  `json:"id"`
	Category string `json:"category" binding:"required,max=50"`
	Name     string `json:"name" binding:"required,max=100"`
	Key      string `json:"key" binding:"required,max=100"`
	Value    string `json:"value" binding:"required,max=500"`
	Type     int32  `json:"type" binding:"required"`
	Visible  bool   `json:"visible"`
	Remark   string `json:"remark"`
}

type ConfigPageReq struct {
	core.PageParam
	Name       string   `form:"name"`
	Key        string   `form:"key"`
	Type       *int32   `form:"type"`
	CreateTime []string `form:"createTime[]"`
}
