package system

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type TenantCreateReq struct {
	Name          string   `json:"name" binding:"required"`
	ContactName   string   `json:"contactName" binding:"required"`
	ContactMobile string   `json:"contactMobile"`
	Status        *int     `json:"status" binding:"required"`    // 改为指针以支持 0 值 (通过 validator 检查不为 nil 即可)
	PackageID     *int64   `json:"packageId" binding:"required"` // 改为指针支持 0 (系统套餐)
	AccountCount  int      `json:"accountCount" binding:"required"`
	ExpireTime    int64    `json:"expireTime,string" binding:"required"` // 增加 ,string 支持字符串类型的数字
	Websites      []string `json:"websites"`
	Username      string   `json:"username" binding:"required"`
	Password      string   `json:"password" binding:"required"`
}

type TenantUpdateReq struct {
	ID            int64    `json:"id" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	ContactName   string   `json:"contactName" binding:"required"`
	ContactMobile string   `json:"contactMobile" binding:"required"`
	Status        *int     `json:"status" binding:"required"`
	PackageID     *int64   `json:"packageId" binding:"required"`
	AccountCount  int      `json:"accountCount" binding:"required"`
	ExpireTime    int64    `json:"expireTime,string" binding:"required"`
	Websites      []string `json:"websites"`
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

// TenantSimpleResp 租户精简信息响应
type TenantSimpleResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TenantRespVO 租户详细信息响应（完整版，后续 CRUD 使用）
type TenantRespVO struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	ContactUserID int64     `json:"contactUserId"`
	ContactName   string    `json:"contactName"`
	ContactMobile string    `json:"contactMobile"`
	Status        int       `json:"status"`
	Websites      []string  `json:"websites"` // 对齐 Java: List<String>
	PackageID     int64     `json:"packageId"`
	ExpireTime    int64     `json:"expireTime"` // Timestamp (ms)
	AccountCount  int       `json:"accountCount"`
	CreateTime    time.Time `json:"createTime"`
}
