package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ConfigSaveReq struct {
	ID       int64  `json:"id"`
	Category string `json:"category" binding:"required,max=50"`
	Name     string `json:"name" binding:"required,max=100"`
	Key      string `json:"key" binding:"required,max=100"`
	Value    string `json:"value" binding:"required,max=500"`
	Visible  *bool  `json:"visible" binding:"required"` // Java: @NotNull
	Remark   string `json:"remark"`
}

type ConfigPageReq struct {
	pagination.PageParam
	Name       string   `form:"name"`
	Key        string   `form:"key"`
	Type       *int32   `form:"type"`
	CreateTime []string `form:"createTime[]"`
}
