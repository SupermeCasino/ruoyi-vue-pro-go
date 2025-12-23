package req

import "time"

// AppMemberUserUpdateReq 修改基本信息请求
type AppMemberUserUpdateReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Sex      int32  `json:"sex"`
}

// AppMemberUserUpdateMobileReq 修改手机请求
type AppMemberUserUpdateMobileReq struct {
	Mobile  string `json:"mobile" binding:"required,len=11"`
	Code    string `json:"code" binding:"required"`
	OldCode string `json:"oldCode"` // 旧手机验证码，可选
}

// AppMemberUserUpdatePasswordReq 修改密码请求
type AppMemberUserUpdatePasswordReq struct {
	Password string `json:"password" binding:"required,min=4,max=16"`
	Code     string `json:"code" binding:"required"` // 验证码
	Mobile   string `json:"mobile"`                  // Mobile usually needed for verification context if not from token
}

// AppMemberUserResetPasswordReq 重置密码请求
type AppMemberUserResetPasswordReq struct {
	Mobile   string `json:"mobile" binding:"required,len=11"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=16"`
}

// ========== Admin API Request VO ==========

// MemberUserUpdateReq Admin 更新会员用户请求
type MemberUserUpdateReq struct {
	ID       int64      `json:"id" binding:"required"`
	Mobile   string     `json:"mobile"`
	Status   int32      `json:"status"`
	Nickname string     `json:"nickname"`
	Avatar   string     `json:"avatar"`
	Name     string     `json:"name"`
	Sex      int32      `json:"sex"`
	AreaID   int64      `json:"areaId"`
	Birthday *time.Time `json:"birthday"`
	Mark     string     `json:"mark"`
	TagIDs   []int64    `json:"tagIds"`
	LevelID  *int64     `json:"levelId"`
	GroupID  *int64     `json:"groupId"`
}

// MemberUserUpdateLevelReq 更新会员等级请求
type MemberUserUpdateLevelReq struct {
	ID      int64  `json:"id" binding:"required"`
	LevelID *int64 `json:"levelId"`
	Reason  string `json:"reason"`
}

// MemberUserUpdatePointReq 更新会员积分请求
type MemberUserUpdatePointReq struct {
	ID    int64 `json:"id" binding:"required"`
	Point int   `json:"point" binding:"required"`
}

// MemberUserPageReq 会员用户分页请求
type MemberUserPageReq struct {
	PageNo     int          `form:"pageNo" binding:"required"`
	PageSize   int          `form:"pageSize" binding:"required"`
	Mobile     string       `form:"mobile"`
	Nickname   string       `form:"nickname"`
	TagIDs     []int64      `form:"tagIds"`
	LevelID    *int64       `form:"levelId"`
	GroupID    *int64       `form:"groupId"`
	LoginDate  []*time.Time `form:"loginDate"`  // 最近登录时间范围
	CreateTime []*time.Time `form:"createTime"` // 创建时间范围
}
