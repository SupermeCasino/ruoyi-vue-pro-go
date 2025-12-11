package model

import (
	"time"
)

// SystemNotice 通知公告表
type SystemNotice struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"column:title;not null" json:"title"`
	Type      int32     `gorm:"column:type;not null" json:"type"`
	Content   string    `gorm:"column:content;not null" json:"content"`
	Status    int32     `gorm:"column:status;not null;default:0" json:"status"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
	TenantID  int64     `gorm:"column:tenant_id;default:0" json:"tenantId"`
}

func (SystemNotice) TableName() string {
	return "system_notice"
}
