package resp

import "time"

// --- Dept ---

type DeptRespVO struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ParentID     int64     `json:"parentId"`
	Sort         int32     `json:"sort"`
	LeaderUserID int64     `json:"leaderUserId"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Status       int32     `json:"status"`
	CreateTime   time.Time `json:"createTime"`
}

type DeptSimpleRespVO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID int64  `json:"parentId"`
}

// --- Post ---

type PostRespVO struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	Sort       int32     `json:"sort"`
	Status     int32     `json:"status"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}

type PostSimpleRespVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
