package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type NoticeSaveReq struct {
	ID      int64  `json:"id"`
	Title   string `json:"title" binding:"required"`
	Type    *int32 `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  *int32 `json:"status" binding:"required"`
}

type NoticePageReq struct {
	pagination.PageParam
	Title  string `form:"title"`
	Status *int32 `form:"status"`
}
