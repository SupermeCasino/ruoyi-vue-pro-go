package member

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type MemberTagResp struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"createTime"`
}

type MemberTagCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type MemberTagUpdateReq struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type MemberTagPageReq struct {
	pagination.PageParam
	Name       string   `form:"name"`
	CreateTime []string `form:"createTime[]"`
}
