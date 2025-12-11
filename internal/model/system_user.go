package model

import (
	"time"
)

// SystemUser 管理后台用户
type SystemUser struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string     `gorm:"column:username;not null" json:"username"`
	Password  string     `gorm:"column:password;not null" json:"-"`
	Nickname  string     `gorm:"column:nickname;not null" json:"nickname"`
	Remark    string     `gorm:"column:remark" json:"remark"`
	DeptID    int64      `gorm:"column:dept_id" json:"deptId"`
	PostIDs   string     `gorm:"column:post_ids" json:"postIds"` // 简单处理为 JSON 字符串，后续可加 Serializer
	Email     string     `gorm:"column:email" json:"email"`
	Mobile    string     `gorm:"column:mobile" json:"mobile"`
	Sex       int32      `gorm:"column:sex" json:"sex"`
	Avatar    string     `gorm:"column:avatar" json:"avatar"`
	Status    int32      `gorm:"column:status;not null" json:"status"`
	LoginIP   string     `gorm:"column:login_ip" json:"loginIp"`
	LoginDate *time.Time `gorm:"column:login_date" json:"loginDate"`

	// TenantBaseDO 字段
	TenantID  int64     `gorm:"column:tenant_id" json:"tenantId"`
	Creator   string    `gorm:"column:creator" json:"creator"`
	Updater   string    `gorm:"column:updater" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"` // Map create_time to CreatedAt
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"` // Map update_time to UpdatedAt
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`             // Use soft_delete plugin
}

// TableName 指定表名
func (SystemUser) TableName() string {
	return "system_users"
}
