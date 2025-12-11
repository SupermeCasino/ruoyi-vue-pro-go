package req

import (
	"backend-go/internal/pkg/core"
)

// KefuMessageCreateReq 发送消息请求
type KefuMessageCreateReq struct {
	ConversationID int64  `json:"conversationId"`                 // 会话编号 (可选，发送给客服时如果不传，则自动查找或创建)
	ContentType    int    `json:"contentType" binding:"required"` // 1-文本 2-图片 3-商品 4-订单
	Content        string `json:"content" binding:"required"`     // 消息内容
}

// KefuMessagePageReq 消息分页请求
type KefuMessagePageReq struct {
	core.PageParam
	ConversationID int64 `form:"conversationId" binding:"required"` // 会话编号
}

// KefuConversationPageReq 会话列表请求
type KefuConversationPageReq struct {
	core.PageParam
	// 可扩展查询条件，如用户昵称等
}
