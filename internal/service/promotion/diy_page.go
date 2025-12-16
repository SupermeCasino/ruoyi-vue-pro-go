package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/types"
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
	GetDiyPageByTemplateId(ctx context.Context, templateId int64) ([]*promotion.PromotionDiyPage, error)
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
	if err := s.validateNameUnique(ctx, 0, req.TemplateID, req.Name); err != nil {
		return 0, err
	}

	page := &promotion.PromotionDiyPage{
		TemplateID:     req.TemplateID,
		Name:           req.Name,
		Remark:         req.Remark,
		PreviewPicUrls: types.StringListFromCSV(req.PreviewPicUrls),
		Property:       req.Property,
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
	if err := s.validateNameUnique(ctx, req.ID, req.TemplateID, req.Name); err != nil {
		return err
	}

	_, err = s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.ID.Eq(req.ID)).Updates(promotion.PromotionDiyPage{
		TemplateID:     req.TemplateID,
		Name:           req.Name,
		Remark:         req.Remark,
		PreviewPicUrls: types.StringListFromCSV(req.PreviewPicUrls),
		Property:       req.Property,
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
	if len(req.CreateTime) == 2 {
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", req.CreateTime[0], time.Local)
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", req.CreateTime[1], time.Local)
		do = do.Where(q.CreateTime.Between(startTime, endTime))
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

func (s *diyPageService) GetDiyPageByTemplateId(ctx context.Context, templateId int64) ([]*promotion.PromotionDiyPage, error) {
	return s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.TemplateID.Eq(templateId)).Find()
}

// Helpers

func (s *diyPageService) validateNameUnique(ctx context.Context, id int64, templateId int64, name string) error {
	q := s.q.PromotionDiyPage
	do := q.WithContext(ctx).Where(q.TemplateID.Eq(templateId), q.Name.Eq(name))
	if id > 0 {
		do = do.Where(q.ID.Neq(id))
	}
	count, err := do.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.NewBizError(400, "页面名称已存在")
	}
	return nil
}

func (s *diyPageService) validateDiyPageExists(ctx context.Context, id int64) (*promotion.PromotionDiyPage, error) {
	page, err := s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(404, "装修页面不存在")
	}
	return page, nil
}

func (s *diyPageService) convertDiyPageToResp(item *promotion.PromotionDiyPage) *resp.DiyPageResp {
	return &resp.DiyPageResp{
		ID:             item.ID,
		TemplateID:     item.TemplateID,
		Name:           item.Name,
		Remark:         item.Remark,
		PreviewPicUrls: []string(item.PreviewPicUrls),
		Property:       item.Property,
		CreateTime:     item.CreateTime,
	}
}
