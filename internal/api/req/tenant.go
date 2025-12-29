package req

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type TenantCreateReq struct {
	Name          string `json:"name" binding:"required"`
	ContactName   string `json:"contactName" binding:"required"`
	ContactMobile string `json:"contactMobile" binding:"required"`
	Status        int    `json:"status" binding:"required"`
	PackageID     int64  `json:"packageId" binding:"required"`
	AccountCount  int    `json:"accountCount" binding:"required"`
	ExpireTime    int64  `json:"expireTime" binding:"required"` // Timestamp (ms)
	Website       string `json:"website"`
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
	ExpireTime    int64  `json:"expireTime" binding:"required"`
	Website       string `json:"website"`
}

type TenantPageReq struct {
	pagination.PageParam
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
