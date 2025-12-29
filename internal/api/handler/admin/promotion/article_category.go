package promotion

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type ArticleCategoryHandler struct {
	svc promotion.ArticleCategoryService
}

func NewArticleCategoryHandler(svc promotion.ArticleCategoryService) *ArticleCategoryHandler {
	return &ArticleCategoryHandler{svc: svc}
}

func (h *ArticleCategoryHandler) CreateArticleCategory(c *gin.Context) {
	var r req.ArticleCategoryCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateArticleCategory(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *ArticleCategoryHandler) UpdateArticleCategory(c *gin.Context) {
	var r req.ArticleCategoryUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateArticleCategory(c, r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *ArticleCategoryHandler) DeleteArticleCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteArticleCategory(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *ArticleCategoryHandler) GetArticleCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetArticleCategory(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

func (h *ArticleCategoryHandler) GetArticleCategoryList(c *gin.Context) {
	var r req.ArticleCategoryListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetArticleCategoryList(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetSimpleList 获得文章分类精简列表 (Only Enabled)
func (h *ArticleCategoryHandler) GetSimpleList(c *gin.Context) {
	res, err := h.svc.GetArticleCategorySimpleList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetArticleCategoryPage 获得文章分类分页
func (h *ArticleCategoryHandler) GetArticleCategoryPage(c *gin.Context) {
	var r req.ArticleCategoryPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetArticleCategoryPage(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
