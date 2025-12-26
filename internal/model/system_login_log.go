package model

// SystemLoginLog 登录日志
type SystemLoginLog struct {
	ID        int64  `gorm:"primaryKey;autoIncrement;comment:日志编号" json:"id"`
	LogType   int    `gorm:"column:log_type;type:int;not null;comment:日志类型" json:"logType"`
	TraceID   string `gorm:"column:trace_id;type:varchar(64);comment:链路追踪编号" json:"traceId"`
	UserID    int64  `gorm:"column:user_id;type:bigint;not null;default:0;comment:用户编号" json:"userId"`
	UserType  int    `gorm:"column:user_type;type:tinyint;not null;comment:用户类型" json:"userType"`
	Username  string `gorm:"column:username;type:varchar(50);not null;default:'';comment:用户账号" json:"username"`
	Result    int    `gorm:"column:result;type:tinyint;not null;comment:登录结果" json:"result"`
	UserIP    string `gorm:"column:user_ip;type:varchar(50);not null;default:'';comment:用户 IP" json:"userIp"`
	UserAgent string `gorm:"column:user_agent;type:varchar(512);not null;default:'';comment:浏览器 UA" json:"userAgent"`
	TenantBaseDO
}

func (SystemLoginLog) TableName() string {
	return "system_login_log"
}
