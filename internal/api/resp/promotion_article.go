package resp

import (
	"time"
)

// ArticleCategoryRespVO 文章分类 Response
type ArticleCategoryRespVO struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	PicURL     string    `json:"picUrl"`
	Sort       int       `json:"sort"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"createTime"`
}

// ArticleCategorySimpleRespVO 文章分类精简 Response
type ArticleCategorySimpleRespVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ArticleRespVO 文章 Response
type ArticleRespVO struct {
	ID              int64     `json:"id"`
	CategoryID      int64     `json:"categoryId"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	PicURL          string    `json:"picUrl"`
	Introduction    string    `json:"introduction"`
	BrowseCount     int       `json:"browseCount"`
	Sort            int       `json:"sort"`
	Status          int       `json:"status"`
	RecommendHot    bool      `json:"recommendHot"`
	RecommendBanner bool      `json:"recommendBanner"`
	Content         string    `json:"content"`
	CreateTime      time.Time `json:"createTime"`
}
