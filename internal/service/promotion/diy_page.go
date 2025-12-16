package promotion

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type DiyPageService interface {
	CreateDiyPage(ctx context.Context, req req.DiyPageCreateReq) (int64, error)
	UpdateDiyPage(ctx context.Context, req req.DiyPageUpdateReq) error
	DeleteDiyPage(ctx context.Context, id int64) error
	GetDiyPage(ctx context.Context, id int64) (*resp.DiyPageResp, error)
	GetDiyPageList(ctx context.Context, ids []int64) ([]*resp.DiyPageResp, error)
	GetDiyPagePage(ctx context.Context, req req.DiyPagePageReq) (*pagination.PageResult[*resp.DiyPageResp], error)
	GetDiyPageProperty(ctx context.Context, id int64) (string, error)
	UpdateDiyPageProperty(ctx context.Context, req req.DiyPagePropertyUpdateReq) error
}

type diyPageService struct {
	q           *query.Query
	templateSvc DiyTemplateService
}

func NewDiyPageService(q *query.Query, templateSvc DiyTemplateService) DiyPageService {
	return &diyPageService{q: q, templateSvc: templateSvc}
}

func (s *diyPageService) CreateDiyPage(ctx context.Context, req req.DiyPageCreateReq) (int64, error) {
	// Validate Template Exists
	if _, err := s.templateSvc.GetDiyTemplate(ctx, req.TemplateID); err != nil {
		return 0, err
	}

	page := &promotion.PromotionDiyPage{
		TemplateID: req.TemplateID,
		Name:       req.Name,
		Remark:     req.Remark,
		Status:     req.Status,
		Property:   req.Property,
	}
	err := s.q.PromotionDiyPage.WithContext(ctx).Create(page)
	return page.ID, err
}

func (s *diyPageService) UpdateDiyPage(ctx context.Context, req req.DiyPageUpdateReq) error {
	_, err := s.validateDiyPageExists(ctx, req.ID)
	if err != nil {
		return err
	}
	// Validate Template Exists
	if _, err := s.templateSvc.GetDiyTemplate(ctx, req.TemplateID); err != nil {
		return err
	}

	_, err = s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.ID.Eq(req.ID)).Updates(promotion.PromotionDiyPage{
		TemplateID: req.TemplateID,
		Name:       req.Name,
		Remark:     req.Remark,
		Status:     req.Status,
		Property:   req.Property,
	})
	return err
}

func (s *diyPageService) DeleteDiyPage(ctx context.Context, id int64) error {
	_, err := s.validateDiyPageExists(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.ID.Eq(id)).Delete()
	return err
}

func (s *diyPageService) GetDiyPage(ctx context.Context, id int64) (*resp.DiyPageResp, error) {
	page, err := s.validateDiyPageExists(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.convertDiyPageToResp(page), nil
}

func (s *diyPageService) GetDiyPagePage(ctx context.Context, req req.DiyPagePageReq) (*pagination.PageResult[*resp.DiyPageResp], error) {
	q := s.q.PromotionDiyPage
	do := q.WithContext(ctx)
	if req.Name != "" {
		do = do.Where(q.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		do = do.Where(q.Status.Eq(*req.Status))
	}

	list, total, err := do.Order(q.ID.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	result := make([]*resp.DiyPageResp, len(list))
	for i, item := range list {
		result[i] = s.convertDiyPageToResp(item)
	}
	return &pagination.PageResult[*resp.DiyPageResp]{List: result, Total: total}, nil
}

func (s *diyPageService) GetDiyPageList(ctx context.Context, ids []int64) ([]*resp.DiyPageResp, error) {
	if len(ids) == 0 {
		return []*resp.DiyPageResp{}, nil
	}
	list, err := s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	result := make([]*resp.DiyPageResp, len(list))
	for i, item := range list {
		result[i] = s.convertDiyPageToResp(item)
	}
	return result, nil
}

func (s *diyPageService) GetDiyPageProperty(ctx context.Context, id int64) (string, error) {
	page, err := s.validateDiyPageExists(ctx, id)
	if err != nil {
		return "", err
	}
	return page.Property, nil
}

func (s *diyPageService) UpdateDiyPageProperty(ctx context.Context, req req.DiyPagePropertyUpdateReq) error {
	// 校验存在
	_, err := s.validateDiyPageExists(ctx, req.ID)
	if err != nil {
		return err
	}
	// 更新
	_, err = s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.ID.Eq(req.ID)).Updates(promotion.PromotionDiyPage{
		Property: req.Property,
	})
	return err
}

// Helpers

func (s *diyPageService) validateDiyPageExists(ctx context.Context, id int64) (*promotion.PromotionDiyPage, error) {
	page, err := s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(404, "装修页面不存在")
	}
	return page, nil
}

func (s *diyPageService) convertDiyPageToResp(item *promotion.PromotionDiyPage) *resp.DiyPageResp {
	return &resp.DiyPageResp{
		ID:         item.ID,
		TemplateID: item.TemplateID,
		Name:       item.Name,
		Remark:     item.Remark,
		Status:     item.Status,
		Property:   item.Property,
		CreateTime: item.CreateTime,
	}
}
