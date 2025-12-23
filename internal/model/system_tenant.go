package model

import (
	"time"
)

// SystemTenant 租户表
type SystemTenant struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	ContactName   string    `gorm:"column:contact_name" json:"contactName"`
	ContactMobile string    `gorm:"column:contact_mobile" json:"contactMobile"`
	Status        int32     `gorm:"column:status;not null;default:0" json:"status"`
	Domain        string    `gorm:"column:domain" json:"domain"`     // Maybe Websites in some versions, but requested Domain in DTO
	Websites      string    `gorm:"column:websites" json:"websites"` // To match service usage
	PackageID     int64     `gorm:"column:package_id" json:"packageId"`
	ExpireDate    time.Time `gorm:"column:expire_time" json:"expireTime"` // Mapped to expire_time
	AccountCount  int32     `gorm:"column:account_count" json:"accountCount"`
	Creator       string    `gorm:"column:creator;default:''" json:"creator"`
	Updater       string    `gorm:"column:updater;default:''" json:"updater"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime     time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted       BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
}

func (SystemTenant) TableName() string {
	return "system_tenant"
}
