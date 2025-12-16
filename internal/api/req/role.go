package req

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
)

// --- Role ---

type RolePageReq struct {
	core.PageParam
	Name         string     `form:"name"`
	Code         string     `form:"code"`
	Status       *int       `form:"status"`
	CreateTimeGe *time.Time `form:"createTime[0]"`
	CreateTimeLe *time.Time `form:"createTime[1]"`
}

type RoleSaveReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   int32  `json:"sort"`
	Status *int   `json:"status" binding:"required"`
	Remark string `json:"remark"`
}

type RoleUpdateStatusReq struct {
	ID     int64 `json:"id" binding:"required"`
	Status *int  `json:"status" binding:"required"`
}

// --- Permission ---

type PermissionAssignRoleMenuReq struct {
	RoleID  int64   `json:"roleId" binding:"required"`
	MenuIDs []int64 `json:"menuIds" binding:"required"`
}

type PermissionAssignRoleDataScopeReq struct {
	RoleID           int64   `json:"roleId" binding:"required"`
	DataScope        int     `json:"dataScope" binding:"required"`
	DataScopeDeptIDs []int64 `json:"dataScopeDeptIds"`
}

type PermissionAssignUserRoleReq struct {
	UserID  int64   `json:"userId" binding:"required"`
	RoleIDs []int64 `json:"roleIds" binding:"required"`
}
