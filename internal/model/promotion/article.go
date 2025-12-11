package promotion

import (
	"time"
)

// PromotionArticleCategory 对应的数据库表：promotion_article_category
type PromotionArticleCategory struct {
	ID     int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"column:name;type:varchar(64);not null;comment:分类名称" json:"name"`
	PicURL string `gorm:"column:pic_url;type:varchar(255);comment:图标地址" json:"picUrl"`
	Sort   int    `gorm:"column:sort;type:int;not null;default:0;comment:排序" json:"sort"`
	Status int    `gorm:"column:status;type:tinyint;not null;default:0;comment:状态" json:"status"` // 0-开启 1-关闭

	Creator    string    `gorm:"column:creator" json:"creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Updater    string    `gorm:"column:updater" json:"updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"column:deleted" json:"deleted"`
	TenantID   int64     `gorm:"column:tenant_id" json:"tenantId"`
}

func (PromotionArticleCategory) TableName() string {
	return "promotion_article_category"
}

// PromotionArticle 对应的数据库表：promotion_article
type PromotionArticle struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CategoryID      int64  `gorm:"column:category_id;type:bigint;not null;comment:分类ID" json:"categoryId"`
	Title           string `gorm:"column:title;type:varchar(64);not null;comment:文章标题" json:"title"`
	Author          string `gorm:"column:author;type:varchar(64);comment:文章作者" json:"author"`
	PicURL          string `gorm:"column:pic_url;type:varchar(255);comment:封面图片" json:"picUrl"`
	Introduction    string `gorm:"column:introduction;type:varchar(255);comment:文章简介" json:"introduction"`
	BrowseCount     int    `gorm:"column:browse_count;type:int;not null;default:0;comment:浏览次数" json:"browseCount"`
	Sort            int    `gorm:"column:sort;type:int;not null;default:0;comment:排序" json:"sort"`
	Status          int    `gorm:"column:status;type:tinyint;not null;default:0;comment:状态" json:"status"` // 0-开启 1-关闭
	RecommendHot    bool   `gorm:"column:recommend_hot;type:tinyint(1);not null;default:0;comment:是否热门(小程序)" json:"recommendHot"`
	RecommendBanner bool   `gorm:"column:recommend_banner;type:tinyint(1);not null;default:0;comment:是否轮播图(小程序)" json:"recommendBanner"`
	Content         string `gorm:"column:content;type:longtext;comment:文章内容" json:"content"`

	Creator    string    `gorm:"column:creator" json:"creator"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Updater    string    `gorm:"column:updater" json:"updater"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"column:deleted" json:"deleted"`
	TenantID   int64     `gorm:"column:tenant_id" json:"tenantId"`
}

func (PromotionArticle) TableName() string {
	return "promotion_article"
}
