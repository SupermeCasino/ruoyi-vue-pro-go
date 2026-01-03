package system

import "time"

// ApiAccessLogResp API访问日志响应
type ApiAccessLogResp struct {
	ID              int64     `json:"id"`
	TraceID         string    `json:"traceId"`
	UserID          int64     `json:"userId"`
	UserType        int       `json:"userType"`
	ApplicationName string    `json:"applicationName"`
	RequestMethod   string    `json:"requestMethod"`
	RequestURL      string    `json:"requestUrl"`
	RequestParams   string    `json:"requestParams"`
	ResponseBody    string    `json:"responseBody"`
	UserIP          string    `json:"userIp"`
	UserAgent       string    `json:"userAgent"`
	OperateModule   string    `json:"operateModule"`
	OperateName     string    `json:"operateName"`
	OperateType     int       `json:"operateType"`
	BeginTime       time.Time `json:"beginTime"`
	EndTime         time.Time `json:"endTime"`
	Duration        int       `json:"duration"`
	ResultCode      int       `json:"resultCode"`
	ResultMsg       string    `json:"resultMsg"`
	CreateTime      time.Time `json:"createTime"`
}

// ApiErrorLogResp API错误日志响应
type ApiErrorLogResp struct {
	ID                        int64      `json:"id"`
	TraceID                   string     `json:"traceId"`
	UserID                    int64      `json:"userId"`
	UserType                  int        `json:"userType"`
	ApplicationName           string     `json:"applicationName"`
	RequestMethod             string     `json:"requestMethod"`
	RequestURL                string     `json:"requestUrl"`
	RequestParams             string     `json:"requestParams"`
	UserIP                    string     `json:"userIp"`
	UserAgent                 string     `json:"userAgent"`
	ExceptionTime             time.Time  `json:"exceptionTime"`
	ExceptionName             string     `json:"exceptionName"`
	ExceptionMessage          string     `json:"exceptionMessage"`
	ExceptionRootCauseMessage string     `json:"exceptionRootCauseMessage"`
	ExceptionStackTrace       string     `json:"exceptionStackTrace"`
	ExceptionClassName        string     `json:"exceptionClassName"`
	ExceptionFileName         string     `json:"exceptionFileName"`
	ExceptionMethodName       string     `json:"exceptionMethodName"`
	ExceptionLineNumber       int        `json:"exceptionLineNumber"`
	ProcessStatus             int        `json:"processStatus"`
	ProcessTime               *time.Time `json:"processTime"`
	ProcessUserID             int64      `json:"processUserId"`
	CreateTime                time.Time  `json:"createTime"`
}
