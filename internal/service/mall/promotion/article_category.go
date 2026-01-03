package promotion

import (
	"context"
	"time"

	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ArticleCategoryService interface {
	CreateArticleCategory(ctx context.Context, req promotion2.ArticleCategoryCreateReq) (int64, error)
	UpdateArticleCategory(ctx context.Context, req promotion2.ArticleCategoryUpdateReq) error
	DeleteArticleCategory(ctx context.Context, id int64) error
	GetArticleCategory(ctx context.Context, id int64) (*promotion2.ArticleCategoryRespVO, error)
	GetArticleCategoryList(ctx context.Context, req promotion2.ArticleCategoryListReq) ([]*promotion2.ArticleCategoryRespVO, error)
	GetArticleCategorySimpleList(ctx context.Context) ([]*promotion2.ArticleCategorySimpleRespVO, error)
	GetArticleCategoryPage(ctx context.Context, req promotion2.ArticleCategoryPageReq) (*pagination.PageResult[*promotion2.ArticleCategoryRespVO], error)
}

type articleCategoryService struct {
	q *query.Query
}

func NewArticleCategoryService(q *query.Query) ArticleCategoryService {
	return &articleCategoryService{q: q}
}

func (s *articleCategoryService) CreateArticleCategory(ctx context.Context, req promotion2.ArticleCategoryCreateReq) (int64, error) {
	category := &promotion.PromotionArticleCategory{
		Name:   req.Name,
		PicURL: req.PicURL,
		Sort:   req.Sort,
		Status: req.Status,
	}
	err := s.q.PromotionArticleCategory.WithContext(ctx).Create(category)
	return category.ID, err
}

func (s *articleCategoryService) UpdateArticleCategory(ctx context.Context, req promotion2.ArticleCategoryUpdateReq) error {
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

func (s *articleCategoryService) GetArticleCategory(ctx context.Context, id int64) (*promotion2.ArticleCategoryRespVO, error) {
	category, err := s.q.PromotionArticleCategory.WithContext(ctx).Where(s.q.PromotionArticleCategory.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(404, "文章分类不存在")
	}
	return &promotion2.ArticleCategoryRespVO{
		ID:         category.ID,
		Name:       category.Name,
		PicURL:     category.PicURL,
		Sort:       category.Sort,
		Status:     category.Status,
		CreateTime: category.CreateTime,
	}, nil
}

func (s *articleCategoryService) GetArticleCategoryList(ctx context.Context, req promotion2.ArticleCategoryListReq) ([]*promotion2.ArticleCategoryRespVO, error) {
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

	result := make([]*promotion2.ArticleCategoryRespVO, len(list))
	for i, item := range list {
		// Manual Filter for Status if needed? No.
		result[i] = &promotion2.ArticleCategoryRespVO{
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

func (s *articleCategoryService) GetArticleCategorySimpleList(ctx context.Context) ([]*promotion2.ArticleCategorySimpleRespVO, error) {
	list, err := s.q.PromotionArticleCategory.WithContext(ctx).
		Where(s.q.PromotionArticleCategory.Status.Eq(consts.CommonStatusEnable)). // 使用启用状态常量替代魔法数字 0
		Order(s.q.PromotionArticleCategory.Sort.Desc()).
		Find()
	if err != nil {
		return nil, err
	}

	result := make([]*promotion2.ArticleCategorySimpleRespVO, len(list))
	for i, item := range list {
		result[i] = &promotion2.ArticleCategorySimpleRespVO{
			ID:   item.ID,
			Name: item.Name,
		}
	}
	return result, nil
}

func (s *articleCategoryService) GetArticleCategoryPage(ctx context.Context, req promotion2.ArticleCategoryPageReq) (*pagination.PageResult[*promotion2.ArticleCategoryRespVO], error) {
	q := s.q.PromotionArticleCategory
	do := q.WithContext(ctx)

	// 过滤条件
	if req.Name != "" {
		do = do.Where(q.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		do = do.Where(q.Status.Eq(*req.Status))
	}
	if len(req.CreateTime) == 2 {
		do = do.Where(q.CreateTime.Between(
			parseTime(req.CreateTime[0]),
			parseTime(req.CreateTime[1]),
		))
	}

	// 统计总数
	total, err := do.Count()
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.PageNo - 1) * req.PageSize
	list, err := do.Order(q.Sort.Desc(), q.ID.Desc()).
		Offset(offset).
		Limit(req.PageSize).
		Find()
	if err != nil {
		return nil, err
	}

	// 转换为响应 VO
	result := make([]*promotion2.ArticleCategoryRespVO, len(list))
	for i, item := range list {
		result[i] = &promotion2.ArticleCategoryRespVO{
			ID:         item.ID,
			Name:       item.Name,
			PicURL:     item.PicURL,
			Sort:       item.Sort,
			Status:     item.Status,
			CreateTime: item.CreateTime,
		}
	}

	return &pagination.PageResult[*promotion2.ArticleCategoryRespVO]{
		List:  result,
		Total: total,
	}, nil
}

// parseTime 辅助函数：解析时间字符串
func parseTime(tStr string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", tStr)
	return t
}
