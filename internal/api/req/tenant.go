package req

import (
	"time"

	"backend-go/internal/pkg/core"
)

type TenantCreateReq struct {
	Name          string `json:"name" binding:"required"`
	ContactName   string `json:"contactName" binding:"required"`
	ContactMobile string `json:"contactMobile" binding:"required"`
	Status        int    `json:"status" binding:"required"`
	PackageID     int64  `json:"packageId" binding:"required"`
	AccountCount  int    `json:"accountCount" binding:"required"`
	ExpireDate    int64  `json:"expireTime" binding:"required"` // Timestamp
	Domain        string `json:"domain"`
	Username      string `json:"username" binding:"required"` // Admin username
	Password      string `json:"password" binding:"required"` // Admin password
}

type TenantUpdateReq struct {
	ID            int64  `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	ContactName   string `json:"contactName" binding:"required"`
	ContactMobile string `json:"contactMobile" binding:"required"`
	Status        int    `json:"status" binding:"required"`
	PackageID     int64  `json:"packageId" binding:"required"`
	AccountCount  int    `json:"accountCount" binding:"required"`
	ExpireDate    int64  `json:"expireTime" binding:"required"`
	Domain        string `json:"domain"`
}

type TenantPageReq struct {
	core.PageParam
	Name          string     `form:"name"`
	ContactName   string     `form:"contactName"`
	ContactMobile string     `form:"contactMobile"`
	Status        *int       `form:"status"`
	CreateTimeGe  *time.Time `form:"createTime[0]"`
	CreateTimeLe  *time.Time `form:"createTime[1]"`
}

type TenantExportReq struct {
	Name          string     `form:"name"`
	ContactName   string     `form:"contactName"`
	ContactMobile string     `form:"contactMobile"`
	Status        *int       `form:"status"`
	CreateTimeGe  *time.Time `form:"createTime[0]"`
	CreateTimeLe  *time.Time `form:"createTime[1]"`
}
