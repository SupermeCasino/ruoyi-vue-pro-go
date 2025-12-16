package trade

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

// AppBrokerageUserBindReqVO 绑定推广员 Request
type AppBrokerageUserBindReqVO struct {
	BindUserID int64 `json:"bindUserId" binding:"required"` // 推广员编号
}

// AppBrokerageUserChildSummaryPageReqVO 下级分销统计分页 Request
type AppBrokerageUserChildSummaryPageReqVO struct {
	core.PageParam
	Nickname string `json:"nickname"` // 下级昵称
	Level    int    `json:"level"`    // 分销层级
	Sorting  string `json:"sorting"`  // 排序字段: brokerageTime, userCount, brokeragePrice
}

// AppBrokerageUserRankPageReqVO 分销用户排行分页 Request
type AppBrokerageUserRankPageReqVO struct {
	core.PageParam
	Times []string `json:"times"` // 时间范围 [start, end]
}
