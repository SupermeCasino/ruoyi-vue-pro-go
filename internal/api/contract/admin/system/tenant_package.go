package system

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type TenantPackageSaveReq struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name" binding:"required,max=30"`
	Status  *int    `json:"status" binding:"required"`
	Remark  string  `json:"remark"`
	MenuIds []int64 `json:"menuIds" binding:"required"`
}

type TenantPackagePageReq struct {
	pagination.PageParam
	Name         string     `form:"name"`
	Status       *int       `form:"status"`
	Remark       string     `form:"remark"`
	CreateTimeGe *time.Time `form:"createTime[0]"`
	CreateTimeLe *time.Time `form:"createTime[1]"`
}

type TenantPackageResp struct {
	ID         int64     `json:"id"`         // 套餐编号
	Name       string    `json:"name"`       // 套餐名
	Status     int       `json:"status"`     // 状态
	Remark     string    `json:"remark"`     // 备注
	MenuIDs    []int64   `json:"menuIds"`    // 关联菜单ID
	CreateTime time.Time `json:"createTime"` // 创建时间
}

type TenantPackageSimpleResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
