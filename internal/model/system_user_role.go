package model

import (
	"time"
)

// SystemUserRole 用户和角色关联表
type SystemUserRole struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null" json:"userId"`
	RoleID    int64     `gorm:"column:role_id;not null" json:"roleId"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
	TenantID  int64     `gorm:"column:tenant_id;default:0" json:"tenantId"` // 注意：UserRoleDO 继承 BaseDO，但在 RuoYi 中通常也包含 tenant_id，需检查数据库。假设有。
}

func (SystemUserRole) TableName() string {
	return "system_user_role"
}
