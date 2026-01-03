package promotion

import "time"

// ArticleCategoryRespVO 文章分类 Response VO
type ArticleCategoryRespVO struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	PicURL     string    `json:"picUrl"`
	Sort       int       `json:"sort"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"createTime"`
}

// ArticleCategorySimpleRespVO 文章分类精简 Response VO
type ArticleCategorySimpleRespVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ArticleCategoryCreateReq struct {
	Name   string `json:"name" binding:"required"`
	PicURL string `json:"picUrl"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

type ArticleCategoryUpdateReq struct {
	ID     int64  `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	PicURL string `json:"picUrl"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

type ArticleCategoryListReq struct {
	Name   string `form:"name"`
	Status *int   `form:"status"`
}

type ArticleCategoryPageReq struct {
	PageNo     int      `form:"pageNo"`
	PageSize   int      `form:"pageSize"`
	Name       string   `form:"name"`
	Status     *int     `form:"status"`
	CreateTime []string `form:"createTime[]"`
}
