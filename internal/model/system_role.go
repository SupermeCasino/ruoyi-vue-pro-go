package model

import (
	"time"
)

// SystemRole 角色表
type SystemRole struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name             string    `gorm:"column:name;not null" json:"name"`
	Code             string    `gorm:"column:code;not null" json:"code"`
	Sort             int32     `gorm:"column:sort" json:"sort"`
	DataScope        int32     `gorm:"column:data_scope;default:1" json:"dataScope"`
	DataScopeDeptIds []int64   `gorm:"column:data_scope_dept_ids;serializer:json" json:"dataScopeDeptIds"` // 使用 serializer:json 自动序列化
	Status           int32     `gorm:"column:status;not null" json:"status"`
	Type             int32     `gorm:"column:type;not null;default:1" json:"type"` // 角色类型(1:内置角色 2:自定义角色)
	Remark           string    `gorm:"column:remark" json:"remark"`
	Creator          string    `gorm:"column:creator;default:''" json:"creator"`
	Updater          string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt        time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt        time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted          BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
	TenantID         int64     `gorm:"column:tenant_id;default:0" json:"tenantId"`
}

func (SystemRole) TableName() string {
	return "system_role"
}
