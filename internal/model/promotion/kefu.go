package promotion

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PromotionKefuConversation 客服会话表
type PromotionKefuConversation struct {
	ID int64 `gorm:"primaryKey;autoIncrement;comment:编号"`

	// 会话所属用户
	UserID int64 `gorm:"comment:用户编号;index"`

	// 最后聊天相关
	LastMessageTime        time.Time `gorm:"comment:最后聊天时间"`
	LastMessageContent     string    `gorm:"comment:最后聊天内容"`
	LastMessageContentType int       `gorm:"comment:最后发送的消息类型"` // 枚举 KeFuMessageContentTypeEnum: 1-文本 2-图片 3-商品 4-订单

	//======================= 会话操作相关 =======================

	// 管理端置顶
	AdminPinned model.BitBool `gorm:"comment:管理端置顶"` // false-未置顶 true-已置顶

	// 用户是否可见 (false - 可见，默认值; true - 不可见，用户删除时设置为 true)
	UserDeleted model.BitBool `gorm:"comment:用户是否可见"` // false-可见 true-不可见

	// 管理员是否可见 (false - 可见，默认值; true - 不可见，管理员删除时设置为 true)
	AdminDeleted model.BitBool `gorm:"comment:管理员是否可见"` // false-可见 true-不可见

	// 管理员未读消息数 (用户发送消息时增加，管理员查看后扣减)
	AdminUnreadMessageCount int `gorm:"comment:管理员未读消息数"`

	Creator    string        `gorm:"column:creator"`
	CreateTime time.Time     `gorm:"column:create_time;autoCreateTime"`
	Updater    string        `gorm:"column:updater"`
	UpdateTime time.Time     `gorm:"column:update_time;autoUpdateTime"`
	Deleted    model.BitBool `gorm:"column:deleted"`
	TenantID   int64         `gorm:"column:tenant_id"`
}

func (PromotionKefuConversation) TableName() string {
	return "promotion_kefu_conversation"
}

// PromotionKefuMessage 客服消息表
type PromotionKefuMessage struct {
	ID             int64         `gorm:"primaryKey;autoIncrement;comment:编号"`
	ConversationID int64         `gorm:"comment:会话编号;index"`
	SenderID       int64         `gorm:"comment:发送人编号"`
	SenderType     int           `gorm:"comment:发送人类型"` // 1-用户 2-客服
	ReceiverID     int64         `gorm:"comment:接收人编号"`
	ReceiverType   int           `gorm:"comment:接收人类型"` // 1-用户 2-客服
	ContentType    int           `gorm:"comment:消息类型"`  // 1-文本 2-图片 3-商品 4-订单
	Content        string        `gorm:"comment:消息内容"`  // JSON 结构，根据 ContentType 解析
	ReadStatus     model.BitBool `gorm:"comment:是否已读"`  // false-未读 true-已读

	Creator    string        `gorm:"column:creator"`
	CreateTime time.Time     `gorm:"column:create_time;autoCreateTime"`
	Updater    string        `gorm:"column:updater"`
	UpdateTime time.Time     `gorm:"column:update_time;autoUpdateTime"`
	Deleted    model.BitBool `gorm:"column:deleted"`
	TenantID   int64         `gorm:"column:tenant_id"`
}

func (PromotionKefuMessage) TableName() string {
	return "promotion_kefu_message"
}
