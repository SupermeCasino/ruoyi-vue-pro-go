package member

import "github.com/wxlbd/ruoyi-mall-go/internal/model"

// MemberTag 会员标签
// Table: member_tag
type MemberTag struct {
	ID   int64  `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name string `gorm:"column:name;type:varchar(30);not null;default:'';comment:标签名称" json:"name"`
	model.TenantBaseDO
}

func (MemberTag) TableName() string {
	return "member_tag"
}
