package req

import "github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"

// BrokerageUserCreateReq 创建分销用户 Request
type BrokerageUserCreateReq struct {
	UserID     int64 `json:"userId" binding:"required"`
	BindUserID int64 `json:"bindUserId"`
}

// BrokerageUserUpdateBindUserReq 修改推广员 Request
type BrokerageUserUpdateBindUserReq struct {
	ID         int64 `json:"id" binding:"required"`
	BindUserID int64 `json:"bindUserId" binding:"required"`
}

// BrokerageUserClearBindUserReq 清除推广员 Request
type BrokerageUserClearBindUserReq struct {
	ID int64 `json:"id" binding:"required"`
}

// BrokerageUserUpdateBrokerageEnabledReq 修改推广资格 Request
type BrokerageUserUpdateBrokerageEnabledReq struct {
	ID      int64 `json:"id" binding:"required"`
	Enabled bool  `json:"enabled"`
}

// BrokerageUserPageReq 分销用户分页 Request
type BrokerageUserPageReq struct {
	core.PageParam
	BindUserID       int64    `form:"bindUserId"`
	BrokerageEnabled *bool    `form:"brokerageEnabled"`
	CreateTime       []string `form:"createTime[]"`   // Range
	Level            int      `form:"level"`          // 推广层级过滤
	BindUserTime     []string `form:"bindUserTime[]"` // Range
}

// BrokerageRecordPageReq 分销记录分页 Request
type BrokerageRecordPageReq struct {
	core.PageParam
	UserID     int64    `form:"userId"`
	BizType    string   `form:"bizType"` // 业务类型: order, withdraw
	Status     int      `form:"status"`
	CreateTime []string `form:"createTime[]"` // Range
	BizID      string   `form:"bizId"`
}

// BrokerageWithdrawPageReq 分销提现分页 Request
type BrokerageWithdrawPageReq struct {
	core.PageParam
	UserID      int64    `form:"userId"`
	Type        int      `form:"type"`
	Status      int      `form:"status"`
	UserName    string   `form:"userName"`
	UserAccount string   `form:"userAccount"`
	BankName    string   `form:"bankName"`
	CreateTime  []string `form:"createTime[]"` // Range
}

// BrokerageWithdrawRejectReq 分销提现驳回 Request
type BrokerageWithdrawRejectReq struct {
	ID          int64  `json:"id" binding:"required"`
	AuditReason string `json:"auditReason" binding:"required"`
}
