package model

// SystemNotice 通知公告表
type SystemNotice struct {
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title   string `gorm:"column:title;not null" json:"title"`
	Type    int32  `gorm:"column:type;not null" json:"type"`
	Content string `gorm:"column:content;not null" json:"content"`
	Status  int32  `gorm:"column:status;not null;default:0" json:"status"`
	TenantBaseDO
}

func (SystemNotice) TableName() string {
	return "system_notice"
}
