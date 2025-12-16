package promotion

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
)

type AppArticleHandler struct {
	articleSvc  promotion.ArticleService
	categorySvc promotion.ArticleCategoryService
}

func NewAppArticleHandler(articleSvc promotion.ArticleService, categorySvc promotion.ArticleCategoryService) *AppArticleHandler {
	return &AppArticleHandler{articleSvc: articleSvc, categorySvc: categorySvc}
}

// GetArticleCategoryList 获得文章分类列表
func (h *AppArticleHandler) GetArticleCategoryList(c *gin.Context) {
	res, err := h.categorySvc.GetArticleCategorySimpleList(c)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// GetArticlePage 获得文章分页
func (h *AppArticleHandler) GetArticlePage(c *gin.Context) {
	var r req.ArticlePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	res, err := h.articleSvc.GetArticlePageApp(c, r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// GetArticle 获得文章详情
func (h *AppArticleHandler) GetArticle(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}

	// 1. Get Detail
	res, err := h.articleSvc.GetArticle(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	// 2. Increment Browse Count (Async or Sync? Sync for now as per plan/Java usually)
	// Ignore error for view count update to avoid blocking read? Or log it?
	// Java: articleService.addArticleBrowseCount(id);
	_ = h.articleSvc.AddArticleBrowseCount(c, id)

	core.WriteSuccess(c, res)
}
