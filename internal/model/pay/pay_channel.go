package pay

import "github.com/wxlbd/ruoyi-mall-go/internal/model"

// PayChannel 支付渠道 DO
type PayChannel struct {
	ID      int64            `gorm:"primaryKey;autoIncrement;comment:渠道编号" json:"id"`
	Code    string           `gorm:"size:32;not null;comment:渠道编码" json:"code"`
	Status  int              `gorm:"default:0;not null;comment:状态" json:"status"` // 参见 CommonStatusEnum
	FeeRate float64          `gorm:"default:0;comment:渠道费率" json:"feeRate"`       // 单位：百分比
	Remark  string           `gorm:"size:255;default:'';comment:备注" json:"remark"`
	AppID   int64            `gorm:"column:app_id;not null;comment:应用编号" json:"appId"`
	Config  *PayClientConfig `gorm:"type:json;comment:支付渠道配置" json:"config"` // 支付渠道配置 (微信/支付宝等)
	model.TenantBaseDO
}

func (PayChannel) TableName() string {
	return "pay_channel"
}
