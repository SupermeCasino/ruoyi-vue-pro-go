package system

import "time"

// LoginLogPageReq 登录日志分页请求
type LoginLogPageReq struct {
	PageNo     int         `form:"pageNo" json:"pageNo"`
	PageSize   int         `form:"pageSize" json:"pageSize"`
	UserIP     string      `form:"userIp" json:"userIp"`
	Username   string      `form:"username" json:"username"`
	Status     *bool       `form:"status" json:"status"`
	CreateTime []time.Time `form:"createTime[]" json:"createTime"`
}

// OperateLogPageReq 操作日志分页请求
type OperateLogPageReq struct {
	PageNo     int         `form:"pageNo" json:"pageNo"`
	PageSize   int         `form:"pageSize" json:"pageSize"`
	UserID     *int64      `form:"userId" json:"userId"`
	BizID      *int64      `form:"bizId" json:"bizId"`
	Type       string      `form:"type" json:"type"`
	SubType    string      `form:"subType" json:"subType"`
	Action     string      `form:"action" json:"action"`
	CreateTime []time.Time `form:"createTime[]" json:"createTime"`
}

// LoginLogResp 登录日志响应
type LoginLogResp struct {
	ID         int64     `json:"id"`
	LogType    int       `json:"logType"`
	UserID     int64     `json:"userId"`
	UserType   int       `json:"userType"`
	TraceID    string    `json:"traceId"`
	Username   string    `json:"username"`
	Result     int       `json:"result"`
	UserIP     string    `json:"userIp"`
	UserAgent  string    `json:"userAgent"`
	CreateTime time.Time `json:"createTime"`
}

// OperateLogResp 操作日志响应
type OperateLogResp struct {
	ID            int64     `json:"id"`
	TraceID       string    `json:"traceId"`
	UserID        int64     `json:"userId"`
	UserName      string    `json:"userName"`
	Type          string    `json:"type"`
	SubType       string    `json:"subType"`
	BizID         int64     `json:"bizId"`
	Action        string    `json:"action"`
	Extra         string    `json:"extra"`
	RequestMethod string    `json:"requestMethod"`
	RequestURL    string    `json:"requestUrl"`
	UserIP        string    `json:"userIp"`
	UserAgent     string    `json:"userAgent"`
	CreateTime    time.Time `json:"createTime"`
}
