package brokerage

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"

	"time"
)

// BrokerageRecord 佣金记录
type BrokerageRecord struct {
	ID              int64      `gorm:"primaryKey;autoIncrement;comment:编号"`
	UserID          int64      `gorm:"column:user_id;not null;comment:用户编号"`
	BizID           string     `gorm:"column:biz_id;size:64;not null;comment:业务编号"`
	BizType         int        `gorm:"column:biz_type;not null;comment:业务类型"`
	Title           string     `gorm:"column:title;size:64;not null;comment:标题"`
	Description     string     `gorm:"column:description;size:255;not null;comment:说明"`
	Price           int        `gorm:"column:price;not null;comment:金额"`
	TotalPrice      int        `gorm:"column:total_price;not null;comment:当前总佣金"`
	Status          int        `gorm:"column:status;not null;comment:状态"`
	FrozenDays      int        `gorm:"column:frozen_days;default:0;comment:冻结时间（天）"`
	UnfreezeTime    *time.Time `gorm:"column:unfreeze_time;comment:解冻时间"`
	SourceUserLevel int        `gorm:"column:source_user_level;default:0;comment:来源用户等级"`
	SourceUserID    int64      `gorm:"column:source_user_id;default:0;comment:来源用户编号"`
	model.TenantBaseDO
}

func (BrokerageRecord) TableName() string {
	return "trade_brokerage_record"
}
