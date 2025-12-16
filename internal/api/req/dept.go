package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// --- Dept ---

type DeptListReq struct {
	Name   string `form:"name"`
	Status *int   `form:"status"`
}

type DeptSaveReq struct {
	ID           int64  `json:"id"`
	Name         string `json:"name" binding:"required"`
	ParentID     int64  `json:"parentId"` // Root is 0
	Sort         int32  `json:"sort"`
	LeaderUserID int64  `json:"leaderUserId"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Status       int    `json:"status" binding:"required"`
}

// --- Post ---

type PostPageReq struct {
	pagination.PageParam
	Code   string `form:"code"`
	Name   string `form:"name"`
	Status *int   `form:"status"`
}

type PostSaveReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   int32  `json:"sort"`
	Status int    `json:"status" binding:"required"`
	Remark string `json:"remark"`
}
