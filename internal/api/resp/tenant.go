package resp

import "time"

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
	Domain        string    `json:"domain"`
	PackageID     int64     `json:"packageId"`
	ExpireDate    int64     `json:"expireTime"` // Timestamp
	AccountCount  int       `json:"accountCount"`
	CreateTime    time.Time `json:"createTime"`
}
