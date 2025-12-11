package resp

import "time"

type MemberTagResp struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Remark    string    `json:"remark"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createTime"`
}
