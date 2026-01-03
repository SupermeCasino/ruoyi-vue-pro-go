package infra

import "time"

// ApiAccessLogPageReq API访问日志分页请求
type ApiAccessLogPageReq struct {
	PageNo          int         `form:"pageNo" json:"pageNo"`
	PageSize        int         `form:"pageSize" json:"pageSize"`
	UserID          *int64      `form:"userId" json:"userId"`
	UserType        *int        `form:"userType" json:"userType"`
	ApplicationName string      `form:"applicationName" json:"applicationName"`
	RequestURL      string      `form:"requestUrl" json:"requestUrl"`
	BeginTime       []time.Time `form:"beginTime[]" json:"beginTime"`
	Duration        *int        `form:"duration" json:"duration"`
	ResultCode      *int        `form:"resultCode" json:"resultCode"`
}

// ApiErrorLogPageReq API错误日志分页请求
type ApiErrorLogPageReq struct {
	PageNo          int         `form:"pageNo" json:"pageNo"`
	PageSize        int         `form:"pageSize" json:"pageSize"`
	UserID          *int64      `form:"userId" json:"userId"`
	UserType        *int        `form:"userType" json:"userType"`
	ApplicationName string      `form:"applicationName" json:"applicationName"`
	RequestURL      string      `form:"requestUrl" json:"requestUrl"`
	ExceptionTime   []time.Time `form:"exceptionTime[]" json:"exceptionTime"`
	ProcessStatus   *int        `form:"processStatus" json:"processStatus"`
}

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
	ProcessUserID             int64      `json:"processUserID"`
	CreateTime                time.Time  `json:"createTime"`
}
