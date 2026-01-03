package promotion

import (
	"context"

	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, req promotion2.ArticleCreateReq) (int64, error)
	UpdateArticle(ctx context.Context, req promotion2.ArticleUpdateReq) error
	DeleteArticle(ctx context.Context, id int64) error
	GetArticle(ctx context.Context, id int64) (*promotion2.ArticleRespVO, error)
	GetArticlePage(ctx context.Context, req promotion2.ArticlePageReq) (*pagination.PageResult[*promotion2.ArticleRespVO], error)
	GetArticlePageApp(ctx context.Context, req promotion2.ArticlePageReq) (*pagination.PageResult[*promotion2.ArticleRespVO], error)
	GetLastArticleByTitle(ctx context.Context, title string) (*promotion2.ArticleRespVO, error)
	AddArticleBrowseCount(ctx context.Context, id int64) error
}

type articleService struct {
	q *query.Query
}

func NewArticleService(q *query.Query) ArticleService {
	return &articleService{q: q}
}

func (s *articleService) CreateArticle(ctx context.Context, req promotion2.ArticleCreateReq) (int64, error) {
	// Validate Category
	if err := s.validateArticleCategory(ctx, req.CategoryID); err != nil {
		return 0, err
	}

	article := &promotion.PromotionArticle{
		CategoryID:      req.CategoryID,
		Title:           req.Title,
		Author:          req.Author,
		PicURL:          req.PicURL,
		Introduction:    req.Introduction,
		BrowseCount:     req.BrowseCount, // Can set initial browse count
		Sort:            req.Sort,
		Status:          req.Status,
		RecommendHot:    model.BitBool(req.RecommendHot),
		RecommendBanner: model.BitBool(req.RecommendBanner),
		Content:         req.Content,
	}
	err := s.q.PromotionArticle.WithContext(ctx).Create(article)
	return article.ID, err
}

func (s *articleService) UpdateArticle(ctx context.Context, req promotion2.ArticleUpdateReq) error {
	// Validate Exists
	if _, err := s.validateArticleExists(ctx, req.ID); err != nil {
		return err
	}
	// Validate Category
	if err := s.validateArticleCategory(ctx, req.CategoryID); err != nil {
		return err
	}

	_, err := s.q.PromotionArticle.WithContext(ctx).Where(s.q.PromotionArticle.ID.Eq(req.ID)).Updates(promotion.PromotionArticle{
		CategoryID:      req.CategoryID,
		Title:           req.Title,
		Author:          req.Author,
		PicURL:          req.PicURL,
		Introduction:    req.Introduction,
		BrowseCount:     req.BrowseCount,
		Sort:            req.Sort,
		Status:          req.Status,
		RecommendHot:    model.BitBool(req.RecommendHot),
		RecommendBanner: model.BitBool(req.RecommendBanner),
		Content:         req.Content,
	})
	return err
}

func (s *articleService) DeleteArticle(ctx context.Context, id int64) error {
	if _, err := s.validateArticleExists(ctx, id); err != nil {
		return err
	}
	_, err := s.q.PromotionArticle.WithContext(ctx).Where(s.q.PromotionArticle.ID.Eq(id)).Delete()
	return err
}

func (s *articleService) GetArticle(ctx context.Context, id int64) (*promotion2.ArticleRespVO, error) {
	article, err := s.validateArticleExists(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.convertArticleToResp(article), nil
}

func (s *articleService) GetLastArticleByTitle(ctx context.Context, title string) (*promotion2.ArticleRespVO, error) {
	q := s.q.PromotionArticle
	article, err := q.WithContext(ctx).Where(q.Title.Eq(title)).Order(q.ID.Desc()).First()
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.NewBizError(404, "文章不存在")
	}
	return s.convertArticleToResp(article), nil
}

func (s *articleService) GetArticlePage(ctx context.Context, req promotion2.ArticlePageReq) (*pagination.PageResult[*promotion2.ArticleRespVO], error) {
	q := s.q.PromotionArticle
	do := q.WithContext(ctx)
	if req.Title != "" {
		do = do.Where(q.Title.Like("%" + req.Title + "%"))
	}
	if req.CategoryID > 0 {
		do = do.Where(q.CategoryID.Eq(req.CategoryID))
	}
	if req.Status != nil {
		do = do.Where(q.Status.Eq(*req.Status))
	}
	// Admin page usually sorts by Sort desc, then ID desc
	list, total, err := do.Order(q.Sort.Asc(), q.ID.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	result := make([]*promotion2.ArticleRespVO, len(list))
	for i, item := range list {
		result[i] = s.convertArticleToResp(item)
	}
	return &pagination.PageResult[*promotion2.ArticleRespVO]{List: result, Total: total}, nil
}

func (s *articleService) GetArticlePageApp(ctx context.Context, req promotion2.ArticlePageReq) (*pagination.PageResult[*promotion2.ArticleRespVO], error) {
	q := s.q.PromotionArticle
	do := q.WithContext(ctx).Where(q.Status.Eq(consts.CommonStatusEnable)) // 使用启用状态常量替代魔法数字 0

	if req.Title != "" {
		do = do.Where(q.Title.Like("%" + req.Title + "%"))
	}
	if req.CategoryID > 0 {
		do = do.Where(q.CategoryID.Eq(req.CategoryID))
	}
	// Add other app filters if needed (e.g., RecommendHot)

	list, total, err := do.Order(q.Sort.Asc(), q.ID.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	result := make([]*promotion2.ArticleRespVO, len(list))
	for i, item := range list {
		result[i] = s.convertArticleToResp(item)
	}
	return &pagination.PageResult[*promotion2.ArticleRespVO]{List: result, Total: total}, nil
}

func (s *articleService) AddArticleBrowseCount(ctx context.Context, id int64) error {
	_, err := s.q.PromotionArticle.WithContext(ctx).Where(s.q.PromotionArticle.ID.Eq(id)).
		Update(s.q.PromotionArticle.BrowseCount, s.q.PromotionArticle.BrowseCount.Add(1))
	return err
}

// Helpers

func (s *articleService) validateArticleExists(ctx context.Context, id int64) (*promotion.PromotionArticle, error) {
	article, err := s.q.PromotionArticle.WithContext(ctx).Where(s.q.PromotionArticle.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(404, "文章不存在")
	}
	return article, nil
}

func (s *articleService) validateArticleCategory(ctx context.Context, categoryID int64) error {
	count, err := s.q.PromotionArticleCategory.WithContext(ctx).Where(s.q.PromotionArticleCategory.ID.Eq(categoryID), s.q.PromotionArticleCategory.Status.Eq(consts.CommonStatusEnable)).Count() // 使用启用状态常量替代魔法数字 0
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.NewBizError(400, "文章分类不存在或已关闭")
	}
	return nil
}

func (s *articleService) convertArticleToResp(item *promotion.PromotionArticle) *promotion2.ArticleRespVO {
	return &promotion2.ArticleRespVO{
		ID:              item.ID,
		CategoryID:      item.CategoryID,
		Title:           item.Title,
		Author:          item.Author,
		PicURL:          item.PicURL,
		Introduction:    item.Introduction,
		BrowseCount:     item.BrowseCount,
		Sort:            item.Sort,
		Status:          item.Status,
		RecommendHot:    bool(item.RecommendHot),
		RecommendBanner: bool(item.RecommendBanner),
		Content:         item.Content,
		CreateTime:      item.CreateTime,
	}
}
