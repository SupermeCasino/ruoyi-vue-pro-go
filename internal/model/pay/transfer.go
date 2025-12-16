package pay

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PayTransfer 转账单
type PayTransfer struct {
	// core.BaseDO (Replaced with explicit fields)
	Creator            string            `gorm:"column:creator;default:''" json:"creator"`
	Updater            string            `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt          time.Time         `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Deleted            model.BitBool     `gorm:"column:deleted;softDelete:flag" json:"-"`
	TenantID           int64             `gorm:"column:tenant_id;default:0" json:"tenantId"`
	ID                 int64             `gorm:"primaryKey;autoIncrement;comment:编号"`
	No                 string            `gorm:"size:64;not null;comment:转账单号"`
	AppID              int64             `gorm:"not null;comment:应用编号"`
	ChannelID          int64             `gorm:"not null;comment:转账渠道编号"`
	ChannelCode        string            `gorm:"size:32;not null;comment:转账渠道编码"`
	MerchantTransferID string            `gorm:"size:64;comment:商户转账单编号"`
	Subject            string            `gorm:"size:512;not null;comment:转账标题"`
	Price              int               `gorm:"not null;comment:转账金额，单位：分"`
	UserAccount        string            `gorm:"size:64;not null;comment:收款人账号"`
	UserName           string            `gorm:"size:64;comment:收款人姓名"`
	Status             int               `gorm:"not null;comment:转账状态"`
	SuccessTime        *time.Time        `gorm:"comment:订单转账成功时间"`
	NotifyURL          string            `gorm:"size:128;comment:异步通知地址"`
	UserIP             string            `gorm:"size:50;comment:用户 IP"`
	ChannelExtras      map[string]string `gorm:"serializer:json;comment:渠道的额外参数"`
	ChannelTransferNo  string            `gorm:"size:64;comment:渠道转账单号"`
	ChannelErrorCode   string            `gorm:"size:128;comment:调用渠道的错误码"`
	ChannelErrorMsg    string            `gorm:"size:256;comment:调用渠道的错误提示"`
	ChannelNotifyData  string            `gorm:"size:4096;comment:渠道的同步/异步通知的内容"`
	ChannelPackageInfo string            `gorm:"size:4096;comment:渠道 package 信息"`
}

func (PayTransfer) TableName() string {
	return "pay_transfer"
}
