package req

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

type MemberPointRecordPageReq struct {
	core.PageParam
	Nickname string `form:"nickname"` // 用户昵称
	BizType  string `form:"bizType"`  // 业务类型
	Title    string `form:"title"`    // 积分标题
}

type AppMemberPointRecordPageReq struct {
	core.PageParam
	AddStatus *bool `form:"addStatus"` // 是否增加积分, nil-全部, true-增加, false-减少
}
