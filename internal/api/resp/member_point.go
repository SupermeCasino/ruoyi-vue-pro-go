package resp

import "time"

type MemberPointRecordResp struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Nickname    string    `json:"nickname"`
	BizID       string    `json:"bizId"`
	BizType     int       `json:"bizType"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	TotalPoint  int       `json:"totalPoint"`
	CreatedAt   time.Time `json:"createTime"`
}

type AppMemberPointRecordResp struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	CreatedAt   time.Time `json:"createTime"`
}
