package model

import (
	"time"
)

// SystemMailLog 邮件日志
type SystemMailLog struct {
	ID               int64             `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID           int64             `gorm:"column:user_id;comment:用户编号" json:"userId"`
	UserType         int               `gorm:"column:user_type;comment:用户类型" json:"userType"`
	ToMails          StringListFromCSV `gorm:"column:to_mails;not null;comment:接收邮箱" json:"toMails"`
	CcMails          StringListFromCSV `gorm:"column:cc_mails;comment:抄送邮箱" json:"ccMails"`
	BccMails         StringListFromCSV `gorm:"column:bcc_mails;comment:密送邮箱" json:"bccMails"`
	AccountID        int64             `gorm:"column:account_id;not null;comment:邮箱账号编号" json:"accountId"`
	FromMail         string            `gorm:"column:from_mail;not null;comment:发送邮箱" json:"fromMail"`
	TemplateID       int64             `gorm:"column:template_id;not null;comment:模板编号" json:"templateId"`
	TemplateCode     string            `gorm:"column:template_code;not null;comment:模板编码" json:"templateCode"`
	TemplateNickname string            `gorm:"column:template_nickname;comment:模版发送人名称" json:"templateNickname"`
	TemplateTitle    string            `gorm:"column:template_title;not null;comment:邮件标题" json:"templateTitle"`
	TemplateContent  string            `gorm:"column:template_content;not null;comment:邮件内容" json:"templateContent"`
	TemplateParams   string            `gorm:"column:template_params;not null;comment:邮件参数" json:"templateParams"` // JSON map
	SendStatus       int               `gorm:"column:send_status;not null;comment:发送状态" json:"sendStatus"`
	SendTime         *time.Time        `gorm:"column:send_time;comment:发送时间" json:"sendTime"`
	SendMessageID    string            `gorm:"column:send_message_id;comment:发送返回的消息 ID" json:"sendMessageId"`
	SendException    string            `gorm:"column:send_exception;comment:发送异常" json:"sendException"`

	TenantBaseDO
}

func (SystemMailLog) TableName() string {
	return "system_mail_log"
}
