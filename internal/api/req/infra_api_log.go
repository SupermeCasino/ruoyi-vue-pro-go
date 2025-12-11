package req

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
