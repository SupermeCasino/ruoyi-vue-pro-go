package resp

import "time"

// KefuConversationResp 客服会话 Response
type KefuConversationResp struct {
	ID                      int64     `json:"id"`
	UserID                  int64     `json:"userId"`
	UserAvatar              string    `json:"userAvatar"`   // 用户头像 (需关联查询)
	UserNickname            string    `json:"userNickname"` // 用户昵称 (需关联查询)
	LastMessageTime         time.Time `json:"lastMessageTime"`
	LastMessageContent      string    `json:"lastMessageContent"`
	LastMessageContentType  int       `json:"lastMessageContentType"`  // 对齐 Java: lastMessageContentType
	AdminPinned             bool      `json:"adminPinned"`             // 管理端置顶
	UserDeleted             bool      `json:"userDeleted"`             // 用户是否删除
	AdminDeleted            bool      `json:"adminDeleted"`            // 管理员是否删除
	AdminUnreadMessageCount int       `json:"adminUnreadMessageCount"` // 对齐 Java: adminUnreadMessageCount
	CreateTime              time.Time `json:"createTime"`
}

// KefuMessageResp 客服消息 Response
type KefuMessageResp struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversationId"`
	SenderID       int64     `json:"senderId"`
	SenderType     int       `json:"senderType"`
	SenderAvatar   string    `json:"senderAvatar,omitempty"` // Java 动态添加
	ReceiverID     int64     `json:"receiverId"`
	ReceiverType   int       `json:"receiverType"`
	ContentType    int       `json:"contentType"`
	Content        string    `json:"content"`
	ReadStatus     bool      `json:"readStatus"`
	CreateTime     time.Time `json:"createTime"`
}
