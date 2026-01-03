package system

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// DictTypePageReq 字典类型分页请求
type DictTypePageReq struct {
	pagination.PageParam
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
	pagination.PageParam
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

// DictTypeSimpleRespVO 字典类型精简信息
type DictTypeSimpleRespVO struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// DictTypeRespVO 字典类型详细信息
type DictTypeRespVO struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Status     int32     `json:"status"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}

// DictDataSimpleRespVO 字典数据精简信息
type DictDataSimpleRespVO struct {
	DictType  string `json:"dictType"`
	Value     string `json:"value"`
	Label     string `json:"label"`
	ColorType string `json:"colorType"`
	CssClass  string `json:"cssClass"`
}

// DictDataRespVO 字典数据详细信息
type DictDataRespVO struct {
	ID         int64     `json:"id"`
	Sort       int32     `json:"sort"`
	Label      string    `json:"label"`
	Value      string    `json:"value"`
	DictType   string    `json:"dictType"`
	Status     int32     `json:"status"`
	ColorType  string    `json:"colorType"`
	CssClass   string    `json:"cssClass"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}
