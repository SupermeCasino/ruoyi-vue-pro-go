package member

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// AppMemberUserInfoResp 用户个人信息响应
type AppMemberUserInfoResp struct {
	ID               int64                   `json:"id"`
	Nickname         string                  `json:"nickname"`
	Avatar           string                  `json:"avatar"`
	Mobile           string                  `json:"mobile"`
	Sex              int32                   `json:"sex"`
	Point            int32                   `json:"point"`
	Experience       int32                   `json:"experience"`
	Level            *AppMemberUserLevelResp `json:"level"`
	BrokerageEnabled bool                    `json:"brokerageEnabled"`
}

type AppMemberUserLevelResp struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Icon  string `json:"icon"`
}

// ========== Admin API Response VO ==========

// MemberUserResp Admin 会员用户响应
type MemberUserResp struct {
	ID         int64      `json:"id"`
	Mobile     string     `json:"mobile"`
	Status     int32      `json:"status"`
	Nickname   string     `json:"nickname"`
	Avatar     string     `json:"avatar"`
	Name       string     `json:"name"`
	Sex        int32      `json:"sex"`
	AreaID     int64      `json:"areaId"`
	AreaName   string     `json:"areaName"`
	Birthday   *time.Time `json:"birthday"`
	Mark       string     `json:"mark"`
	TagIDs     []int64    `json:"tagIds"`
	LevelID    int64      `json:"levelId"`
	GroupID    int64      `json:"groupId"`
	RegisterIP string     `json:"registerIp"`
	LoginIP    string     `json:"loginIp"`
	LoginDate  *time.Time `json:"loginDate"`
	CreateTime time.Time  `json:"createTime"`
	// 扩展字段
	Point      int32    `json:"point"`
	TagNames   []string `json:"tagNames"`
	LevelName  string   `json:"levelName"`
	GroupName  string   `json:"groupName"`
	Experience int32    `json:"experience"`
}

type AppMemberUserUpdateReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Sex      int    `json:"sex" binding:"required"`
}

type AppMemberUserUpdateMobileReq struct {
	Mobile  string `json:"mobile" binding:"required"`
	Code    string `json:"code" binding:"required"`
	OldCode string `json:"oldCode"`
}

type AppMemberUserResetPasswordReq struct {
	Mobile   string `json:"mobile" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AppMemberUserUpdatePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

type MemberUserUpdateReq struct {
	ID       int64      `json:"id" binding:"required"`
	Mobile   string     `json:"mobile" binding:"required"`
	Status   int        `json:"status" binding:"required"`
	Nickname string     `json:"nickname" binding:"required"`
	Avatar   string     `json:"avatar"`
	Name     string     `json:"name"`
	Sex      int        `json:"sex"`
	AreaID   int        `json:"areaId"`
	Birthday *time.Time `json:"birthday"`
	Mark     string     `json:"mark"`
	TagIDs   []int64    `json:"tagIds"`
	LevelID  *int64     `json:"levelId"`
	GroupID  *int64     `json:"groupId"`
}

type MemberUserPageReq struct {
	pagination.PageParam
	Mobile     string       `form:"mobile"`
	Nickname   string       `form:"nickname"`
	TagIDs     []int64      `form:"tagIds"`
	LevelID    *int64       `form:"levelId"`
	GroupID    *int64       `form:"groupId"`
	LoginDate  []*time.Time `form:"loginDate"`  // 最近登录时间范围
	CreateTime []*time.Time `form:"createTime"` // 创建时间范围
}

type MemberUserUpdateLevelReq struct {
	ID      int64  `json:"id" binding:"required"`
	LevelID int64  `json:"levelId" binding:"required"`
	Reason  string `json:"reason"`
}

type MemberUserUpdatePointReq struct {
	ID    int64 `json:"id" binding:"required"`
	Point int32 `json:"point" binding:"required"`
}
