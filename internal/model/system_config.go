package model

import (
	"time"
)

// SystemConfig 参数配置表
type SystemConfig struct {
	ID        int64   `gorm:"primaryKey;autoIncrement;comment:参数主键" json:"id"`
	Category  string  `gorm:"size:50;not null;comment:参数分类" json:"category"`
	Name      string  `gorm:"size:100;not null;comment:参数名称" json:"name"`
	ConfigKey string  `gorm:"size:100;not null;comment:参数键名" json:"configKey"`
	Value     string  `gorm:"size:500;not null;comment:参数键值" json:"value"`
	Type      int32   `gorm:"size:4;not null;default:1;comment:参数类型" json:"type"`
	Visible   BitBool `gorm:"not null;default:true;comment:是否可见" json:"visible"`
	Remark    string  `gorm:"size:500;comment:备注" json:"remark"`

	// Base fields
	Creator   string    `gorm:"column:creator;size:64;comment:创建者" json:"creator"`
	Updater   string    `gorm:"column:updater;size:64;comment:更新者" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"-"`
}

func (SystemConfig) TableName() string {
	return "infra_config"
}
