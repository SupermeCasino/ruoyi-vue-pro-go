package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayWalletTransaction 会员钱包流水表
type PayWalletTransaction struct {
	ID       int64  `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	WalletID int64  `gorm:"column:wallet_id;not null;comment:钱包编号" json:"walletId"`
	BizType  int    `gorm:"column:biz_type;not null;comment:关联业务类型" json:"bizType"` // 1: 充值, 2: 支付...
	BizID    string `gorm:"column:biz_id;size:64;not null;comment:关联业务编号" json:"bizId"`
	No       string `gorm:"column:no;size:64;not null;comment:流水号" json:"no"`
	Title    string `gorm:"size:128;not null;comment:流水标题" json:"title"`
	Price    int    `gorm:"column:price;not null;default:0;comment:交易金额" json:"price"`      // 单位：分。正数：收入，负数：支出
	Balance  int    `gorm:"column:balance;not null;default:0;comment:交易后余额" json:"balance"` // 单位：分
	model.TenantBaseDO
}

func (PayWalletTransaction) TableName() string {
	return "pay_wallet_transaction"
}
