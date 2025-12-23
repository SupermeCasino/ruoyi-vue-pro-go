package resp

import "time"

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
	CreatedAt  time.Time  `json:"createTime"`
	// 扩展字段
	Point      int32    `json:"point"`
	TagNames   []string `json:"tagNames"`
	LevelName  string   `json:"levelName"`
	GroupName  string   `json:"groupName"`
	Experience int32    `json:"experience"`
}
