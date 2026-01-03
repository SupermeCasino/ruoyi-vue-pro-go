package system

import "time"

// JobResp 定时任务响应
type JobResp struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Status         int       `json:"status"`
	HandlerName    string    `json:"handlerName"`
	HandlerParam   string    `json:"handlerParam"`
	CronExpression string    `json:"cronExpression"`
	RetryCount     int       `json:"retryCount"`
	RetryInterval  int       `json:"retryInterval"`
	MonitorTimeout *int      `json:"monitorTimeout"`
	CreateTime     time.Time `json:"createTime"`
}

// JobLogResp 定时任务日志响应
type JobLogResp struct {
	ID           int64      `json:"id"`
	JobID        int64      `json:"jobId"`
	HandlerName  string     `json:"handlerName"`
	HandlerParam string     `json:"handlerParam"`
	ExecuteIndex int        `json:"executeIndex"`
	BeginTime    time.Time  `json:"beginTime"`
	EndTime      *time.Time `json:"endTime"`
	Duration     *int       `json:"duration"`
	Status       int        `json:"status"`
	Result       string     `json:"result"`
	CreateTime   time.Time  `json:"createTime"`
}
