package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// PromotionBanner 首页轮播图
// Table: promotion_banner
type PromotionBanner struct {
	ID       int64  `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Title    string `gorm:"column:title;type:varchar(64);not null;comment:标题" json:"title"`
	PicURL   string `gorm:"column:pic_url;type:varchar(255);not null;comment:图片地址" json:"picUrl"`
	Url      string `gorm:"column:url;type:varchar(255);not null;comment:跳转地址" json:"url"`
	Status   int    `gorm:"column:status;type:tinyint;not null;default:0;comment:状态" json:"status"` // 0: 开启, 1: 关闭
	Sort     int    `gorm:"column:sort;type:int;not null;default:0;comment:排序" json:"sort"`
	Position int    `gorm:"column:position;type:tinyint;not null;default:1;comment:位置" json:"position"` // 1: 首页
	Memo     string `gorm:"column:memo;type:varchar(255);comment:备注" json:"memo"`
	model.TenantBaseDO
}

func (PromotionBanner) TableName() string {
	return "promotion_banner"
}
