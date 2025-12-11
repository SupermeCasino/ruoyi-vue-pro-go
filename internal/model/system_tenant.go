package model

import (
	"time"
)

// SystemTenant 租户表
type SystemTenant struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string     `gorm:"column:name;not null" json:"name"`
	ContactUserId int64      `gorm:"column:contact_user_id" json:"contactUserId"`
	ContactName   string     `gorm:"column:contact_name;default:''" json:"contactName"`
	ContactMobile string     `gorm:"column:contact_mobile;default:''" json:"contactMobile"`
	Status        int32      `gorm:"column:status;not null;default:0" json:"status"`
	Websites      string     `gorm:"column:websites;default:''" json:"websites"` // JSON 数组存储
	PackageID     int64      `gorm:"column:package_id;default:0" json:"packageId"`
	ExpireTime    *time.Time `gorm:"column:expire_time" json:"expireTime"`
	AccountCount  int32      `gorm:"column:account_count;default:-1" json:"accountCount"`
	Creator       string     `gorm:"column:creator;default:''" json:"creator"`
	Updater       string     `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt     time.Time  `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt     time.Time  `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted       BitBool    `gorm:"column:deleted;softDelete:flag" json:"-"`
}

func (SystemTenant) TableName() string {
	return "system_tenant"
}
