package resp

import "time"

// AuthLoginResp 登录响应
type AuthLoginResp struct {
	UserId       int64     `json:"userId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresTime  time.Time `json:"expiresTime"`
}

type AuthPermissionInfoResp struct {
	User        UserVO   `json:"user"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Menus       []MenuVO `json:"menus"`
}

type UserVO struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	DeptID   int64  `json:"deptId"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type MenuVO struct {
	ID            int64    `json:"id"`
	ParentID      int64    `json:"parentId"`
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Component     string   `json:"component"`
	ComponentName string   `json:"componentName"`
	Icon          string   `json:"icon"`
	Visible       bool     `json:"visible"`
	KeepAlive     bool     `json:"keepAlive"`
	AlwaysShow    bool     `json:"alwaysShow"`
	Children      []MenuVO `json:"children,omitempty"`
}
