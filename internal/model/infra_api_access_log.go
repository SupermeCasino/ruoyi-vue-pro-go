package model

import (
	"time"

	"gorm.io/gorm"
)

// InfraApiAccessLog API访问日志
type InfraApiAccessLog struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;comment:日志编号" json:"id"`
	TraceID         string         `gorm:"column:trace_id;type:varchar(64);comment:链路追踪编号" json:"traceId"`
	UserID          int64          `gorm:"column:user_id;type:bigint;default:0;comment:用户编号" json:"userId"`
	UserType        int            `gorm:"column:user_type;type:tinyint;default:0;comment:用户类型" json:"userType"`
	ApplicationName string         `gorm:"column:application_name;type:varchar(50);not null;comment:应用名" json:"applicationName"`
	RequestMethod   string         `gorm:"column:request_method;type:varchar(16);not null;comment:请求方法名" json:"requestMethod"`
	RequestURL      string         `gorm:"column:request_url;type:varchar(255);not null;comment:请求地址" json:"requestUrl"`
	RequestParams   string         `gorm:"column:request_params;type:text;comment:请求参数" json:"requestParams"`
	ResponseBody    string         `gorm:"column:response_body;type:text;comment:响应结果" json:"responseBody"`
	UserIP          string         `gorm:"column:user_ip;type:varchar(50);comment:用户IP" json:"userIp"`
	UserAgent       string         `gorm:"column:user_agent;type:varchar(512);comment:浏览器UA" json:"userAgent"`
	OperateModule   string         `gorm:"column:operate_module;type:varchar(50);comment:操作模块" json:"operateModule"`
	OperateName     string         `gorm:"column:operate_name;type:varchar(50);comment:操作名" json:"operateName"`
	OperateType     int            `gorm:"column:operate_type;type:tinyint;default:0;comment:操作分类" json:"operateType"`
	BeginTime       time.Time      `gorm:"column:begin_time;comment:开始请求时间" json:"beginTime"`
	EndTime         time.Time      `gorm:"column:end_time;comment:结束请求时间" json:"endTime"`
	Duration        int            `gorm:"column:duration;type:int;default:0;comment:执行时长" json:"duration"`
	ResultCode      int            `gorm:"column:result_code;type:int;default:0;comment:结果码" json:"resultCode"`
	ResultMsg       string         `gorm:"column:result_msg;type:varchar(512);comment:结果提示" json:"resultMsg"`
	Creator         string         `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater         string         `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt       time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt       time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted;index;comment:是否删除" json:"-"`
}

func (InfraApiAccessLog) TableName() string {
	return "infra_api_access_log"
}
