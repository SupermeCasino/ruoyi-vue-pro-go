package resp

import "time"

type UserRespVO struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	Remark     string    `json:"remark"`
	DeptID     int64     `json:"deptId"`
	PostIDs    []int64   `json:"postIds"`
	RoleIDs    []int64   `json:"roleIds"`
	Email      string    `json:"email"`
	Mobile     string    `json:"mobile"`
	Sex        int32     `json:"sex"`
	Avatar     string    `json:"avatar"`
	Status     int32     `json:"status"`
	LoginIP    string    `json:"loginIp"`
	LoginDate  time.Time `json:"loginDate"`
	CreateTime time.Time `json:"createTime"`
}

type UserProfileRespVO struct {
	*UserRespVO
	Roles []*RoleRespVO `json:"roles"`
	Posts []*PostRespVO `json:"posts"`
}

// UserImportRespVO generic import response if needed
type UserImportRespVO struct {
	ReqUsername string `json:"username"`
	FailDesc    string `json:"failDesc"`
}
