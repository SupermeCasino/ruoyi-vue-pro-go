package resp

import "time"

type NoticeRespVO struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Type       int32     `json:"type"`
	Content    string    `json:"content"`
	Status     int32     `json:"status"`
	CreateTime time.Time `json:"createTime"`
}
