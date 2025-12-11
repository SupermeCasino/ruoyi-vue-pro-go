package req

import "backend-go/internal/pkg/core"

// MemberSignInConfigCreateReq 签到规则创建请求
type MemberSignInConfigCreateReq struct {
	Day        int `json:"day" validate:"required,gt=0"`
	Point      int `json:"point" validate:"gte=0"`
	Experience int `json:"experience" validate:"gte=0"`
	Status     int `json:"status" validate:"required,oneof=0 1"`
}

// MemberSignInConfigUpdateReq 签到规则更新请求
type MemberSignInConfigUpdateReq struct {
	ID         int64 `json:"id" validate:"required,gt=0"`
	Day        int   `json:"day" validate:"required,gt=0"`
	Point      int   `json:"point" validate:"gte=0"`
	Experience int   `json:"experience" validate:"gte=0"`
	Status     int   `json:"status" validate:"required,oneof=0 1"`
}

// MemberSignInConfigPageReq 签到规则分页请求
// No specific page req in Java Controller for Config, just List?
// Java: MemberSignInConfigController.java: getSignInConfigList(@RequestParam("status") Integer status)
// Actually standard simple list or page. Java uses simple list for configs usually?
// Check Java Controller again...
/*
   @GetMapping("/list")
   public CommonResult<List<MemberSignInConfigRespVO>> getSignInConfigList(@RequestParam(value = "status", required = false) Integer status) {
*/
// It's a list.

// MemberSignInRecordPageReq 签到记录分页请求
type MemberSignInRecordPageReq struct {
	core.PageParam
	Nickname string `form:"nickname"`
	Day      *int   `form:"day"`
	UserID   int64  `form:"userId"` // Optional filter
}

// AppMemberSignInRecordPageReq App签到记录分页请求
type AppMemberSignInRecordPageReq struct {
	core.PageParam
}
