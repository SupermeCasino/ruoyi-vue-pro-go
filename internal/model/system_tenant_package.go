package model

import (
	"time"
)

// SystemTenantPackage 租户套餐表
type SystemTenantPackage struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	PackageID int64     `gorm:"column:package_id;default:0" json:"packageId"` // 某些旧逻辑可能用，通常用 ID
	Status    int32     `gorm:"column:status;not null;default:0" json:"status"`
	MenuIDs   string    `gorm:"column:menu_ids;type:text" json:"menuIds"` // JSON 数组存储
	Remark    string    `gorm:"column:remark;default:''" json:"remark"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
}

func (SystemTenantPackage) TableName() string {
	return "system_tenant_package"
}
