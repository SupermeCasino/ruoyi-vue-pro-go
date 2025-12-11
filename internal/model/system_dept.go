package model

import (
	"time"
)

// SystemDept 部门表
type SystemDept struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:name;not null;default:''" json:"name"`
	ParentID     int64     `gorm:"column:parent_id;not null;default:0" json:"parentId"`
	Sort         int32     `gorm:"column:sort;not null;default:0" json:"sort"`
	LeaderUserID int64     `gorm:"column:leader_user_id;default:0" json:"leaderUserId"`
	Phone        string    `gorm:"column:phone;default:''" json:"phone"`
	Email        string    `gorm:"column:email;default:''" json:"email"`
	Status       int32     `gorm:"column:status;not null;default:0" json:"status"` // 0:开启, 1:禁用
	Creator      string    `gorm:"column:creator;default:''" json:"creator"`
	Updater      string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt    time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt    time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted      BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
	TenantID     int64     `gorm:"column:tenant_id;default:0" json:"tenantId"`
}

func (SystemDept) TableName() string {
	return "system_dept"
}

// SystemPost 岗位表
type SystemPost struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;not null;default:''" json:"name"`
	Code      string    `gorm:"column:code;not null;default:''" json:"code"`
	Sort      int32     `gorm:"column:sort;not null;default:0" json:"sort"`
	Status    int32     `gorm:"column:status;not null;default:0" json:"status"` // 0:开启, 1:禁用
	Remark    string    `gorm:"column:remark;default:''" json:"remark"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
}

func (SystemPost) TableName() string {
	return "system_post"
}
