package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"time"
)

// PayWallet 会员钱包表
type PayWallet struct {
	ID            int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	UserID        int64         `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	UserType      int           `gorm:"column:user_type;not null;default:0;comment:用户类型" json:"userType"`
	Balance       int           `gorm:"column:balance;not null;default:0;comment:余额" json:"balance"`                // 单位：分
	TotalExpense  int           `gorm:"column:total_expense;not null;default:0;comment:累计支出" json:"totalExpense"`   // 单位：分
	TotalRecharge int           `gorm:"column:total_recharge;not null;default:0;comment:累计充值" json:"totalRecharge"` // 单位：分
	FreezePrice   int           `gorm:"column:freeze_price;not null;default:0;comment:冻结金额" json:"freezePrice"`     // 单位：分
	Creator       string        `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater       string        `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt     time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt     time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted       model.BitBool `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID      int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (PayWallet) TableName() string {
	return "pay_wallet"
}
