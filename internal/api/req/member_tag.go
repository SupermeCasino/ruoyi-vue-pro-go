package req

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

type MemberTagCreateReq struct {
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark"`
}

type MemberTagUpdateReq struct {
	ID     int64  `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark"`
}

type MemberTagPageReq struct {
	core.PageParam
	Name *string `form:"name"`
}
