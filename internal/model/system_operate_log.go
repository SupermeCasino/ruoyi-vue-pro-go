package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemOperateLog 操作日志
type SystemOperateLog struct {
	ID            int64          `gorm:"primaryKey;autoIncrement;comment:日志编号" json:"id"`
	TraceID       string         `gorm:"column:trace_id;type:varchar(64);comment:链路追踪编号" json:"traceId"`
	UserID        int64          `gorm:"column:user_id;type:bigint;not null;default:0;comment:用户编号" json:"userId"`
	UserType      int            `gorm:"column:user_type;type:tinyint;not null;comment:用户类型" json:"userType"`
	Type          string         `gorm:"column:type;type:varchar(50);not null;default:'';comment:操作模块类型" json:"type"`
	SubType       string         `gorm:"column:sub_type;type:varchar(50);not null;default:'';comment:操作名" json:"subType"`
	BizID         int64          `gorm:"column:biz_id;type:bigint;not null;default:0;comment:操作模块业务编号" json:"bizId"`
	Action        string         `gorm:"column:action;type:varchar(2000);not null;default:'';comment:操作内容" json:"action"`
	Extra         string         `gorm:"column:extra;type:varchar(2000);not null;default:'';comment:拓展字段" json:"extra"`
	RequestMethod string         `gorm:"column:request_method;type:varchar(16);not null;default:'';comment:请求方法名" json:"requestMethod"`
	RequestURL    string         `gorm:"column:request_url;type:varchar(255);not null;default:'';comment:请求地址" json:"requestUrl"`
	UserIP        string         `gorm:"column:user_ip;type:varchar(50);not null;default:'';comment:用户 IP" json:"userIp"`
	UserAgent     string         `gorm:"column:user_agent;type:varchar(512);not null;default:'';comment:浏览器 UA" json:"userAgent"`
	Creator       string         `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater       string         `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt     time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt     time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted;index;comment:是否删除" json:"-"`
}

func (SystemOperateLog) TableName() string {
	return "system_operate_log"
}
