package trade

import (
	"backend-go/internal/model"
	"time"
)

// TradeConfig 交易中心 - 交易配置
// Table: trade_config
type TradeConfig struct {
	ID                          int64         `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	AppID                       int64         `gorm:"column:app_id;not null;comment:支付应用ID" json:"appId"`
	AfterSaleDeadlineDays       int           `gorm:"column:after_sale_deadline_days;not null;comment:售后期限(天)" json:"afterSaleDeadlineDays"`
	PayTimeoutMinutes           int           `gorm:"column:pay_timeout_minutes;not null;comment:支付超时(分钟)" json:"payTimeoutMinutes"`
	AutoReceiveDays             int           `gorm:"column:auto_receive_days;not null;comment:自动收货(天)" json:"autoReceiveDays"`
	AutoCommentDays             int           `gorm:"column:auto_comment_days;not null;comment:自动好评(天)" json:"autoCommentDays"`
	BrokerageWithdrawMinPrice   int           `gorm:"column:brokerage_withdraw_min_price;default:0;comment:提现最低金额" json:"brokerageWithdrawMinPrice"`
	BrokerageWithdrawFeePercent int           `gorm:"column:brokerage_withdraw_fee_percent;default:0;comment:提现手续费百分比" json:"brokerageWithdrawFeePercent"`
	BrokerageEnabled            model.BitBool `gorm:"column:brokerage_enabled;default:0;comment:是否开启分销" json:"brokerageEnabled"`
	BrokerageFrozenDays         int           `gorm:"column:brokerage_frozen_days;default:0;comment:分销佣金冻结时间" json:"brokerageFrozenDays"`
	BrokerageFirstPercent       int           `gorm:"column:brokerage_first_percent;default:0;comment:一级分销比例" json:"brokerageFirstPercent"`
	BrokerageSecondPercent      int           `gorm:"column:brokerage_second_percent;default:0;comment:二级分销比例" json:"brokerageSecondPercent"`
	BrokeragePosterUrls         string        `gorm:"column:brokerage_poster_urls;default:'';comment:分销海报图" json:"brokeragePosterUrls"`
	Creator                     string        `gorm:"column:creator;size:64;default:'';comment:创建者"`
	CreateTime                  time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	Updater                     string        `gorm:"column:updater;size:64;default:'';comment:更新者"`
	UpdateTime                  time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	Deleted                     model.BitBool `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
	TenantID                    int64         `gorm:"column:tenant_id;not null;default:0;comment:租户编号"`
}

func (TradeConfig) TableName() string {
	return "trade_config"
}
