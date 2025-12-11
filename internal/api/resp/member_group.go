package resp

import "time"

// MemberGroupResp 用户分组 Response
type MemberGroupResp struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Remark    string    `json:"remark"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createTime"`
}

// MemberGroupSimpleResp 用户分组精简信息 Response
type MemberGroupSimpleResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
