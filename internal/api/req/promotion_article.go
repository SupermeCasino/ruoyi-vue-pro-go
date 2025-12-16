package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// ArticleCategoryCreateReq 文章分类创建请求
type ArticleCategoryCreateReq struct {
	Name   string `json:"name" binding:"required"`
	PicURL string `json:"picUrl"`
	Sort   int    `json:"sort"`
	Status int    `json:"status" binding:"required"` // 0-开启 1-关闭
}

// ArticleCategoryUpdateReq 文章分类更新请求
type ArticleCategoryUpdateReq struct {
	ID     int64  `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	PicURL string `json:"picUrl"`
	Sort   int    `json:"sort"`
	Status int    `json:"status" binding:"required"`
}

// ArticleCategoryListReq 文章分类列表请求
type ArticleCategoryListReq struct {
	Name   string `form:"name"`
	Status *int   `form:"status"` // 0-开启 1-关闭
}

// ArticleCreateReq 文章创建请求
type ArticleCreateReq struct {
	CategoryID      int64  `json:"categoryId" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author"`
	PicURL          string `json:"picUrl"`
	Introduction    string `json:"introduction"`
	BrowseCount     int    `json:"browseCount"`
	Sort            int    `json:"sort"`
	Status          int    `json:"status" binding:"required"`
	RecommendHot    bool   `json:"recommendHot"`
	RecommendBanner bool   `json:"recommendBanner"`
	Content         string `json:"content"`
}

// ArticleUpdateReq 文章更新请求
type ArticleUpdateReq struct {
	ID              int64  `json:"id" binding:"required"`
	CategoryID      int64  `json:"categoryId" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author"`
	PicURL          string `json:"picUrl"`
	Introduction    string `json:"introduction"`
	BrowseCount     int    `json:"browseCount"`
	Sort            int    `json:"sort"`
	Status          int    `json:"status" binding:"required"`
	RecommendHot    bool   `json:"recommendHot"`
	RecommendBanner bool   `json:"recommendBanner"`
	Content         string `json:"content"`
}

// ArticlePageReq 文章分页请求
type ArticlePageReq struct {
	pagination.PageParam
	Title      string `form:"title"`
	CategoryID int64  `form:"categoryId"`
	Status     *int   `form:"status"`
}
