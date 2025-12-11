package resp

import "time"

// TenantSimpleResp 租户精简信息响应
type TenantSimpleResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TenantRespVO 租户详细信息响应（完整版，后续 CRUD 使用）
type TenantRespVO struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	ContactName   string     `json:"contactName,omitempty"`
	ContactMobile string     `json:"contactMobile,omitempty"`
	Status        int32      `json:"status"`
	Websites      []string   `json:"websites,omitempty"` // 转换为数组（注意：Java 中可能存储为 JSON）
	PackageID     int64      `json:"packageId,omitempty"`
	ExpireTime    *time.Time `json:"expireTime,omitempty"`
	AccountCount  int32      `json:"accountCount,omitempty"`
	CreateTime    *time.Time `json:"createTime,omitempty"`
}
