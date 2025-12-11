package resp

import "time"

type RoleRespVO struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	Sort             int32     `json:"sort"`
	Status           int32     `json:"status"`
	Type             int32     `json:"type"`
	Remark           string    `json:"remark"`
	DataScope        int32     `json:"dataScope"`
	DataScopeDeptIDs []int64   `json:"dataScopeDeptIds"`
	CreateTime       time.Time `json:"createTime"`
}

type RoleSimpleRespVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
