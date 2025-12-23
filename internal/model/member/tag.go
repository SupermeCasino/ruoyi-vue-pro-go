package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"

	"time"
)

// MemberTag 会员标签
// Table: member_tag
type MemberTag struct {
	ID        int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name      string        `gorm:"column:name;type:varchar(30);not null;default:'';comment:标签名称" json:"name"`
	Creator   string        `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater   string        `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	Deleted   model.BitBool `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
	TenantID  int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (MemberTag) TableName() string {
	return "member_tag"
}
