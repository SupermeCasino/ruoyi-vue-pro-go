package promotion

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
)

type ArticleCategoryService interface {
	CreateArticleCategory(ctx context.Context, req req.ArticleCategoryCreateReq) (int64, error)
	UpdateArticleCategory(ctx context.Context, req req.ArticleCategoryUpdateReq) error
	DeleteArticleCategory(ctx context.Context, id int64) error
	GetArticleCategory(ctx context.Context, id int64) (*resp.ArticleCategoryRespVO, error)
	GetArticleCategoryList(ctx context.Context, req req.ArticleCategoryListReq) ([]*resp.ArticleCategoryRespVO, error)
	GetArticleCategorySimpleList(ctx context.Context) ([]*resp.ArticleCategorySimpleRespVO, error)
}

type articleCategoryService struct {
	q *query.Query
}

func NewArticleCategoryService(q *query.Query) ArticleCategoryService {
	return &articleCategoryService{q: q}
}

func (s *articleCategoryService) CreateArticleCategory(ctx context.Context, req req.ArticleCategoryCreateReq) (int64, error) {
	category := &promotion.PromotionArticleCategory{
		Name:   req.Name,
		PicURL: req.PicURL,
		Sort:   req.Sort,
		Status: req.Status,
	}
	err := s.q.PromotionArticleCategory.WithContext(ctx).Create(category)
	return category.ID, err
}

func (s *articleCategoryService) UpdateArticleCategory(ctx context.Context, req req.ArticleCategoryUpdateReq) error {
	_, err := s.q.PromotionArticleCategory.WithContext(ctx).Where(s.q.PromotionArticleCategory.ID.Eq(req.ID)).Updates(promotion.PromotionArticleCategory{
		Name:   req.Name,
		PicURL: req.PicURL,
		Sort:   req.Sort,
		Status: req.Status,
	})
	return err
}

func (s *articleCategoryService) DeleteArticleCategory(ctx context.Context, id int64) error {
	// Check if used by articles
	count, err := s.q.PromotionArticle.WithContext(ctx).Where(s.q.PromotionArticle.CategoryID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.NewBizError(400, "该分类下有文章，无法删除")
	}

	_, err = s.q.PromotionArticleCategory.WithContext(ctx).Where(s.q.PromotionArticleCategory.ID.Eq(id)).Delete()
	return err
}

func (s *articleCategoryService) GetArticleCategory(ctx context.Context, id int64) (*resp.ArticleCategoryRespVO, error) {
	category, err := s.q.PromotionArticleCategory.WithContext(ctx).Where(s.q.PromotionArticleCategory.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(404, "文章分类不存在")
	}
	return &resp.ArticleCategoryRespVO{
		ID:         category.ID,
		Name:       category.Name,
		PicURL:     category.PicURL,
		Sort:       category.Sort,
		Status:     category.Status,
		CreateTime: category.CreateTime,
	}, nil
}

func (s *articleCategoryService) GetArticleCategoryList(ctx context.Context, req req.ArticleCategoryListReq) ([]*resp.ArticleCategoryRespVO, error) {
	q := s.q.PromotionArticleCategory
	do := q.WithContext(ctx)
	if req.Name != "" {
		do = do.Where(q.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		do = do.Where(q.Status.Eq(*req.Status))
	}
	// Re-evaluating DTO.
	// Let's implement basics without status filter for List for now to show all.
	// Or check if Status is used.
	// Actually, usually admin queries all.

	list, err := do.Order(q.Sort.Asc()).Find()
	if err != nil {
		return nil, err
	}

	result := make([]*resp.ArticleCategoryRespVO, len(list))
	for i, item := range list {
		// Manual Filter for Status if needed? No.
		result[i] = &resp.ArticleCategoryRespVO{
			ID:         item.ID,
			Name:       item.Name,
			PicURL:     item.PicURL,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreateTime,
		}
	}
	return result, nil
}

func (s *articleCategoryService) GetArticleCategorySimpleList(ctx context.Context) ([]*resp.ArticleCategorySimpleRespVO, error) {
	list, err := s.q.PromotionArticleCategory.WithContext(ctx).
		Where(s.q.PromotionArticleCategory.Status.Eq(0)). // Only Enable
		Order(s.q.PromotionArticleCategory.Sort.Asc()).
		Find()
	if err != nil {
		return nil, err
	}

	result := make([]*resp.ArticleCategorySimpleRespVO, len(list))
	for i, item := range list {
		result[i] = &resp.ArticleCategorySimpleRespVO{
			ID:   item.ID,
			Name: item.Name,
		}
	}
	return result, nil
}
