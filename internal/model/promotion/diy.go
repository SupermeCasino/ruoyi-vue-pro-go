package promotion

import (
	"time"
)

// PromotionDiyTemplate 装修模板表
type PromotionDiyTemplate struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"column:name;type:varchar(64);not null;comment:模板名称" json:"name"`
	CoverImage   string `gorm:"column:cover_image;type:varchar(255);comment:封面图片" json:"coverImage"`
	PreviewImage string `gorm:"column:preview_image;type:varchar(255);comment:预览图片" json:"previewImage"`
	Status       int    `gorm:"column:status;type:tinyint;not null;default:0;comment:状态" json:"status"` // 0-开启 1-关闭
	Property     string `gorm:"column:property;type:longtext;comment:模板属性" json:"property"`             // JSON
	Sort         int    `gorm:"column:sort;type:int;not null;default:0;comment:排序" json:"sort"`
	Remark       string `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`

	Creator    string    `gorm:"column:creator" json:"creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Updater    string    `gorm:"column:updater" json:"updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"column:deleted" json:"deleted"`
	TenantID   int64     `gorm:"column:tenant_id" json:"tenantId"`
}

func (PromotionDiyTemplate) TableName() string {
	return "promotion_diy_template"
}

// PromotionDiyPage 装修页面表
type PromotionDiyPage struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateID int64  `gorm:"column:template_id;type:bigint;not null;comment:模板ID" json:"templateId"`
	Name       string `gorm:"column:name;type:varchar(64);not null;comment:页面名称" json:"name"`
	Remark     string `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	Status     int    `gorm:"column:status;type:tinyint;not null;default:0;comment:状态" json:"status"` // 0-开启 1-关闭
	Property   string `gorm:"column:property;type:longtext;comment:页面属性" json:"property"`             // JSON

	Creator    string    `gorm:"column:creator" json:"creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Updater    string    `gorm:"column:updater" json:"updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"column:deleted" json:"deleted"`
	TenantID   int64     `gorm:"column:tenant_id" json:"tenantId"`
}

func (PromotionDiyPage) TableName() string {
	return "promotion_diy_page"
}
