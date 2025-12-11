package trade

import "backend-go/internal/pkg/core"

type AppBrokerageRecordPageReqVO struct {
	core.PageParam
	BizType    string   `form:"bizType"`      // 业务类型
	Status     int      `form:"status"`       // 状态
	CreateTime []string `form:"createTime[]"` // 创建时间
}
