package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// MemberSignInConfig 签到规则
type MemberSignInConfig struct {
	ID         int64 `gorm:"primaryKey;autoIncrement;comment:规则自增主键" json:"id"`
	Day        int   `gorm:"comment:签到第 x 天" json:"day"`
	Point      int   `gorm:"comment:奖励积分" json:"point"`
	Experience int   `gorm:"comment:奖励经验" json:"experience"`
	Status     int   `gorm:"default:0;comment:状态" json:"status"` // 参见 CommonStatusEnum
	model.TenantBaseDO
}

func (MemberSignInConfig) TableName() string {
	return "member_sign_in_config"
}

// MemberSignInRecord 签到记录
type MemberSignInRecord struct {
	ID         int64 `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	UserID     int64 `gorm:"column:user_id;comment:签到用户" json:"userId"`
	Day        int   `gorm:"comment:第几天签到" json:"day"`
	Point      int   `gorm:"comment:签到的积分" json:"point"`
	Experience int   `gorm:"comment:签到的经验" json:"experience"`
	model.TenantBaseDO
}

func (MemberSignInRecord) TableName() string {
	return "member_sign_in_record"
}
