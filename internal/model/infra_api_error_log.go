package model

import (
	"time"
)

// InfraApiErrorLog API错误日志
type InfraApiErrorLog struct {
	ID                        int64      `gorm:"primaryKey;autoIncrement;comment:日志编号" json:"id"`
	TraceID                   string     `gorm:"column:trace_id;type:varchar(64);comment:链路追踪编号" json:"traceId"`
	UserID                    int64      `gorm:"column:user_id;type:bigint;default:0;comment:用户编号" json:"userId"`
	UserType                  int        `gorm:"column:user_type;type:tinyint;default:0;comment:用户类型" json:"userType"`
	ApplicationName           string     `gorm:"column:application_name;type:varchar(50);not null;comment:应用名" json:"applicationName"`
	RequestMethod             string     `gorm:"column:request_method;type:varchar(16);not null;comment:请求方法名" json:"requestMethod"`
	RequestURL                string     `gorm:"column:request_url;type:varchar(255);not null;comment:请求地址" json:"requestUrl"`
	RequestParams             string     `gorm:"column:request_params;type:text;comment:请求参数" json:"requestParams"`
	UserIP                    string     `gorm:"column:user_ip;type:varchar(50);comment:用户IP" json:"userIp"`
	UserAgent                 string     `gorm:"column:user_agent;type:varchar(512);comment:浏览器UA" json:"userAgent"`
	ExceptionTime             time.Time  `gorm:"column:exception_time;comment:异常发生时间" json:"exceptionTime"`
	ExceptionName             string     `gorm:"column:exception_name;type:varchar(128);comment:异常名" json:"exceptionName"`
	ExceptionMessage          string     `gorm:"column:exception_message;type:text;comment:异常导致的消息" json:"exceptionMessage"`
	ExceptionRootCauseMessage string     `gorm:"column:exception_root_cause_message;type:text;comment:异常导致的根消息" json:"exceptionRootCauseMessage"`
	ExceptionStackTrace       string     `gorm:"column:exception_stack_trace;type:text;comment:异常的栈轨迹" json:"exceptionStackTrace"`
	ExceptionClassName        string     `gorm:"column:exception_class_name;type:varchar(512);comment:异常发生的类全名" json:"exceptionClassName"`
	ExceptionFileName         string     `gorm:"column:exception_file_name;type:varchar(512);comment:异常发生的类文件" json:"exceptionFileName"`
	ExceptionMethodName       string     `gorm:"column:exception_method_name;type:varchar(512);comment:异常发生的方法名" json:"exceptionMethodName"`
	ExceptionLineNumber       int        `gorm:"column:exception_line_number;type:int;comment:异常发生的方法所在行" json:"exceptionLineNumber"`
	ProcessStatus             int        `gorm:"column:process_status;type:tinyint;default:0;comment:处理状态" json:"processStatus"`
	ProcessTime               *time.Time `gorm:"column:process_time;comment:处理时间" json:"processTime"`
	ProcessUserID             int64      `gorm:"column:process_user_id;type:bigint;default:0;comment:处理用户编号" json:"processUserId"`
	TenantBaseDO
}

func (InfraApiErrorLog) TableName() string {
	return "infra_api_error_log"
}
