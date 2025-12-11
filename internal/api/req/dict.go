package req

import "backend-go/internal/pkg/core"

// DictTypePageReq 字典类型分页请求
type DictTypePageReq struct {
	core.PageParam
	Name   string `form:"name"`
	Type   string `form:"type"`
	Status *int   `form:"status"` // 指针允许空值
}

// DictTypeSaveReq 字典类型创建/修改请求
type DictTypeSaveReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Status int    `json:"status" binding:"required"`
	Remark string `json:"remark"`
}

// DictDataPageReq 字典数据分页请求
type DictDataPageReq struct {
	core.PageParam
	Label    string `form:"label"`
	DictType string `form:"dictType"`
	Status   *int   `form:"status"`
}

// DictDataSaveReq 字典数据创建/修改请求
type DictDataSaveReq struct {
	ID        int64  `json:"id"`
	Sort      int32  `json:"sort"`
	Label     string `json:"label" binding:"required"`
	Value     string `json:"value" binding:"required"`
	DictType  string `json:"dictType" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	ColorType string `json:"colorType"`
	CssClass  string `json:"cssClass"`
	Remark    string `json:"remark"`
}
