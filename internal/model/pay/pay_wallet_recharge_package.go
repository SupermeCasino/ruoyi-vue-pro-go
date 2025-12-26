package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayWalletRechargePackage 充值套餐表
type PayWalletRechargePackage struct {
	ID         int64  `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name       string `gorm:"size:64;not null;comment:套餐名称" json:"name"`
	PayPrice   int    `gorm:"column:pay_price;not null;default:0;comment:支付金额" json:"payPrice"`     // 单位：分
	BonusPrice int    `gorm:"column:bonus_price;not null;default:0;comment:赠送金额" json:"bonusPrice"` // 单位：分
	Status     int    `gorm:"column:status;not null;default:0;comment:状态" json:"status"`            // 0: 开启, 1: 关闭
	model.TenantBaseDO
}

func (PayWalletRechargePackage) TableName() string {
	return "pay_wallet_recharge_package"
}
