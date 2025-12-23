package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type MemberTagCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type MemberTagUpdateReq struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type MemberTagPageReq struct {
	pagination.PageParam
	Name *string `form:"name"`
}
