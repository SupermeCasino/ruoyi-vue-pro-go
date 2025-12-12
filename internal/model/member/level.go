package member

import (
	"backend-go/internal/model"

	"time"
)

// MemberLevel 会员等级
// Table: member_level
type MemberLevel struct {
	ID              int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name            string        `gorm:"column:name;type:varchar(64);not null;default:'';comment:等级名称" json:"name"`
	Level           int           `gorm:"column:level;type:int;not null;default:0;comment:等级" json:"level"`
	Experience      int           `gorm:"column:experience;type:int;not null;default:0;comment:升级经验" json:"experience"`
	DiscountPercent int           `gorm:"column:discount_percent;type:int;not null;default:100;comment:享受折扣" json:"discountPercent"`
	Icon            string        `gorm:"column:icon;type:varchar(255);default:'';comment:等级图标" json:"icon"`
	BackgroundURL   string        `gorm:"column:background_url;type:varchar(255);default:'';comment:等级背景图" json:"backgroundUrl"`
	Status          int           `gorm:"column:status;type:int;not null;default:0;comment:状态" json:"status"` // 0: 开启, 1: 关闭
	Result          string        `gorm:"-" json:"result"`                                                    // Ignore
	Remark          string        `gorm:"column:remark;type:varchar(255);default:'';comment:备注"`              // BaseDO usually has remark? Java DO didn't show it but BaseDO might.
	Creator         string        `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater         string        `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt       time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt       time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	Deleted         model.BitBool `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
	TenantID         int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (MemberLevel) TableName() string {
	return "member_level"
}
