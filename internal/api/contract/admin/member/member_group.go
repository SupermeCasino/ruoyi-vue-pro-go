package member

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// MemberGroupResp 用户分组 Response
type MemberGroupResp struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"createTime"`
}

// MemberGroupSimpleResp 用户分组精简信息 Response
type MemberGroupSimpleResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MemberGroupCreateReq struct {
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}

type MemberGroupUpdateReq struct {
	ID     int64  `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}

type MemberGroupPageReq struct {
	pagination.PageParam
	Name       string   `form:"name"`
	Status     *int     `form:"status"`
	CreateTime []string `form:"createTime[]"`
}
