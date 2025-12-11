package resp

import "time"

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
