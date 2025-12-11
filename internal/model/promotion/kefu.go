package promotion

import (
	"time"
)

// PromotionKefuConversation 客服会话表
type PromotionKefuConversation struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement;comment:编号"`
	UserID             int64     `gorm:"comment:用户编号;index"`
	AdminID            int64     `gorm:"comment:客服编号;index"`
	UserDeleted        bool      `gorm:"comment:用户是否删除"` // false-未删除 true-已删除
	AdminDeleted       bool      `gorm:"comment:客服是否删除"` // false-未删除 true-已删除
	AdminUnreadCount   int       `gorm:"comment:客服未读消息数"`
	UserUnreadCount    int       `gorm:"comment:用户未读消息数"`
	LastMessageTime    time.Time `gorm:"comment:最后一条消息时间"`
	LastMessageContent string    `gorm:"comment:最后一条消息内容"`
	LastMessageType    int       `gorm:"comment:最后一条消息类型"` // 1-文本 2-图片 3-商品 4-订单
	Status             int       `gorm:"comment:会话状态"`     // 0-接待中 1-结束

	Creator    string    `gorm:"column:creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	Updater    string    `gorm:"column:updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
	Deleted    int       `gorm:"column:deleted"`
	TenantID   int64     `gorm:"column:tenant_id"`
}

func (PromotionKefuConversation) TableName() string {
	return "promotion_kefu_conversation"
}

// PromotionKefuMessage 客服消息表
type PromotionKefuMessage struct {
	ID             int64  `gorm:"primaryKey;autoIncrement;comment:编号"`
	ConversationID int64  `gorm:"comment:会话编号;index"`
	SenderID       int64  `gorm:"comment:发送人编号"`
	SenderType     int    `gorm:"comment:发送人类型"` // 1-用户 2-客服
	ReceiverID     int64  `gorm:"comment:接收人编号"`
	ReceiverType   int    `gorm:"comment:接收人类型"` // 1-用户 2-客服
	ContentType    int    `gorm:"comment:消息类型"`  // 1-文本 2-图片 3-商品 4-订单
	Content        string `gorm:"comment:消息内容"`  // JSON 结构，根据 ContentType 解析
	ReadStatus     bool   `gorm:"comment:是否已读"`  // false-未读 true-已读

	Creator    string    `gorm:"column:creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	Updater    string    `gorm:"column:updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
	Deleted    int       `gorm:"column:deleted"`
	TenantID   int64     `gorm:"column:tenant_id"`
}

func (PromotionKefuMessage) TableName() string {
	return "promotion_kefu_message"
}
