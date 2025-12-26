package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// MemberGroup 会员分组
type MemberGroup struct {
	ID     int64  `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name   string `gorm:"column:name;comment:名称" json:"name"`     // 名称
	Remark string `gorm:"column:remark;comment:备注" json:"remark"` // 备注
	Status int    `gorm:"column:status;comment:状态" json:"status"` // 状态
	model.TenantBaseDO
}

// TableName 表名
func (MemberGroup) TableName() string {
	return "member_group"
}
