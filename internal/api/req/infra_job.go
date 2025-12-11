package req

import "time"

// JobSaveReq 定时任务创建/修改请求
type JobSaveReq struct {
	ID             *int64 `json:"id"`
	Name           string `json:"name" binding:"required"`
	HandlerName    string `json:"handlerName" binding:"required"`
	HandlerParam   string `json:"handlerParam"`
	CronExpression string `json:"cronExpression" binding:"required"`
	RetryCount     int    `json:"retryCount"`
	RetryInterval  int    `json:"retryInterval"`
	MonitorTimeout *int   `json:"monitorTimeout"`
}

// JobPageReq 定时任务分页请求
type JobPageReq struct {
	PageNo      int    `form:"pageNo" json:"pageNo"`
	PageSize    int    `form:"pageSize" json:"pageSize"`
	Name        string `form:"name" json:"name"`
	HandlerName string `form:"handlerName" json:"handlerName"`
	Status      *int   `form:"status" json:"status"`
}

// JobLogPageReq 定时任务日志分页请求
type JobLogPageReq struct {
	PageNo      int         `form:"pageNo" json:"pageNo"`
	PageSize    int         `form:"pageSize" json:"pageSize"`
	JobID       *int64      `form:"jobId" json:"jobId"`
	HandlerName string      `form:"handlerName" json:"handlerName"`
	BeginTime   []time.Time `form:"beginTime[]" json:"beginTime"`
	Status      *int        `form:"status" json:"status"`
}
