package req

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
