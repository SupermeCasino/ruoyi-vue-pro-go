package promotion

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// ArticleRespVO 文章响应 VO
type ArticleRespVO struct {
	ID              int64     `json:"id"`
	CategoryID      int64     `json:"categoryId"`
	SpuID           int64     `json:"spuId"`
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

// ArticleCreateReq 文章创建请求
type ArticleCreateReq struct {
	CategoryID      int64  `json:"categoryId" binding:"required"`
	SpuID           int64  `json:"spuId"`
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author"`
	PicURL          string `json:"picUrl"`
	Introduction    string `json:"introduction"`
	BrowseCount     int    `json:"browseCount"`
	Sort            int    `json:"sort"`
	Status          *int   `json:"status" binding:"required"`
	RecommendHot    bool   `json:"recommendHot"`
	RecommendBanner bool   `json:"recommendBanner"`
	Content         string `json:"content"`
}

// ArticleUpdateReq 文章更新请求
type ArticleUpdateReq struct {
	ID              int64  `json:"id" binding:"required"`
	CategoryID      int64  `json:"categoryId" binding:"required"`
	SpuID           int64  `json:"spuId"`
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author"`
	PicURL          string `json:"picUrl"`
	Introduction    string `json:"introduction"`
	BrowseCount     int    `json:"browseCount"`
	Sort            int    `json:"sort"`
	Status          *int   `json:"status" binding:"required"`
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
