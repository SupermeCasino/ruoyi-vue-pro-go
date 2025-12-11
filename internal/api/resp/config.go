package resp

import "time"

type ConfigRespVO struct {
	ID         int64     `json:"id"`
	Category   string    `json:"category"`
	Name       string    `json:"name"`
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	Type       int32     `json:"type"`
	Visible    bool      `json:"visible"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}
