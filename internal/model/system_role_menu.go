package model

import (
	"time"
)

// SystemRoleMenu 角色和菜单关联表
type SystemRoleMenu struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID    int64     `gorm:"column:role_id;not null" json:"roleId"`
	MenuID    int64     `gorm:"column:menu_id;not null" json:"menuId"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
	TenantID  int64     `gorm:"column:tenant_id;default:0" json:"tenantId"`
}

func (SystemRoleMenu) TableName() string {
	return "system_role_menu"
}
