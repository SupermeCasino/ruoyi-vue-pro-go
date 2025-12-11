package resp

import "time"

// KefuConversationResp 客服会话 Response
type KefuConversationResp struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"userId"`
	UserAvatar         string    `json:"userAvatar"`   // 用户头像 (需关联查询)
	UserNickname       string    `json:"userNickname"` // 用户昵称 (需关联查询)
	AdminID            int64     `json:"adminId"`
	AdminUnreadCount   int       `json:"adminUnreadCount"`
	UserUnreadCount    int       `json:"userUnreadCount"`
	LastMessageTime    time.Time `json:"lastMessageTime"`
	LastMessageContent string    `json:"lastMessageContent"`
	LastMessageType    int       `json:"lastMessageType"`
	Status             int       `json:"status"`
	CreateTime         time.Time `json:"createTime"`
}

// KefuMessageResp 客服消息 Response
type KefuMessageResp struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversationId"`
	SenderID       int64     `json:"senderId"`
	SenderType     int       `json:"senderType"`
	ReceiverID     int64     `json:"receiverId"`
	ReceiverType   int       `json:"receiverType"`
	ContentType    int       `json:"contentType"`
	Content        string    `json:"content"`
	ReadStatus     bool      `json:"readStatus"`
	CreateTime     time.Time `json:"createTime"`
}
