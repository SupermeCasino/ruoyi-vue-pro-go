package promotion

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/types"
)

// PromotionDiyTemplate 装修模板表
type PromotionDiyTemplate struct {
	ID             int64                   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name           string                  `gorm:"column:name;type:varchar(64);not null;comment:模板名称" json:"name"`
	PreviewPicUrls types.StringListFromCSV `gorm:"column:preview_pic_urls;type:varchar(2000);comment:预览图片" json:"previewPicUrls"`
	Property       string                  `gorm:"column:property;type:longtext;comment:模板属性" json:"property"` // JSON
	Remark         string                  `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	Used           model.BitBool           `gorm:"column:used;type:bit(1);not null;default:0;comment:是否使用" json:"used"`
	UsedTime       *time.Time              `gorm:"column:used_time;comment:使用时间" json:"usedTime"`

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
	ID             int64                   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateID     int64                   `gorm:"column:template_id;type:bigint;not null;comment:模板ID" json:"templateId"`
	Name           string                  `gorm:"column:name;type:varchar(64);not null;comment:页面名称" json:"name"`
	Remark         string                  `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	PreviewPicUrls types.StringListFromCSV `gorm:"column:preview_pic_urls;type:varchar(2000);comment:预览图片" json:"previewPicUrls"`
	Property       string                  `gorm:"column:property;type:longtext;comment:页面属性" json:"property"` // JSON

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
