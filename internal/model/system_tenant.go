package model

import (
	"time"
)

// SystemTenant 租户表
type SystemTenant struct {
	ID            int64             `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string            `gorm:"column:name;not null" json:"name"`
	ContactUserID int64             `gorm:"column:contact_user_id" json:"contactUserId"`
	ContactName   string            `gorm:"column:contact_name" json:"contactName"`
	ContactMobile string            `gorm:"column:contact_mobile" json:"contactMobile"`
	Status        int32             `gorm:"column:status;not null;default:0" json:"status"`
	Websites      StringListFromCSV `gorm:"column:website" json:"websites"` // 对齐 Java: List<String> + StringListTypeHandler (CSV)
	PackageID     int64             `gorm:"column:package_id" json:"packageId"`
	ExpireDate    time.Time         `gorm:"column:expire_time" json:"expireTime"` // 对齐 Java 契约字段名
	AccountCount  int32             `gorm:"column:account_count" json:"accountCount"`
	BaseDO
}

func (SystemTenant) TableName() string {
	return "system_tenant"
}
