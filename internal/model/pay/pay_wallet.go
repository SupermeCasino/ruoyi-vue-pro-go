package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayWallet 会员钱包表
type PayWallet struct {
	ID            int64 `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	UserID        int64 `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	UserType      int   `gorm:"column:user_type;not null;default:0;comment:用户类型" json:"userType"`
	Balance       int   `gorm:"column:balance;not null;default:0;comment:余额" json:"balance"`                // 单位：分
	TotalExpense  int   `gorm:"column:total_expense;not null;default:0;comment:累计支出" json:"totalExpense"`   // 单位：分
	TotalRecharge int   `gorm:"column:total_recharge;not null;default:0;comment:累计充值" json:"totalRecharge"` // 单位：分
	FreezePrice   int   `gorm:"column:freeze_price;not null;default:0;comment:冻结金额" json:"freezePrice"`     // 单位：分
	model.TenantBaseDO
}

func (PayWallet) TableName() string {
	return "pay_wallet"
}
