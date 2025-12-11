package model

import (
	"time"
)

// SystemUserPost 用户和岗位关联表
type SystemUserPost struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null" json:"userId"`
	PostID    int64     `gorm:"column:post_id;not null" json:"postId"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
	TenantID  int64     `gorm:"column:tenant_id;default:0" json:"tenantId"`
}

func (SystemUserPost) TableName() string {
	return "system_user_post"
}
