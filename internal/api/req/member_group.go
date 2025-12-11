package req

import (
	"backend-go/internal/pkg/core"
)

// MemberGroupCreateReq 创建用户分组 Request
type MemberGroupCreateReq struct {
	Name   string `json:"name" binding:"required"`   // 名称
	Remark string `json:"remark"`                    // 备注
	Status int    `json:"status" binding:"required"` // 状态
}

// MemberGroupUpdateReq 更新用户分组 Request
type MemberGroupUpdateReq struct {
	ID     int64  `json:"id" binding:"required"`     // 编号
	Name   string `json:"name" binding:"required"`   // 名称
	Remark string `json:"remark"`                    // 备注
	Status int    `json:"status" binding:"required"` // 状态
}

// MemberGroupPageReq 用户分组分页 Request
type MemberGroupPageReq struct {
	core.PageParam
	Name       string   `form:"name"`       // 名称
	Status     *int     `form:"status"`     // 状态
	CreateTime []string `form:"createTime"` // 创建时间
}
