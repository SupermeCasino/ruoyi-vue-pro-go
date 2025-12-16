package promotion

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type DiyTemplateService interface {
	CreateDiyTemplate(ctx context.Context, req req.DiyTemplateCreateReq) (int64, error)
	UpdateDiyTemplate(ctx context.Context, req req.DiyTemplateUpdateReq) error
	UseDiyTemplate(ctx context.Context, id int64) error
	DeleteDiyTemplate(ctx context.Context, id int64) error
	GetDiyTemplate(ctx context.Context, id int64) (*resp.DiyTemplateResp, error)
	GetDiyTemplatePage(ctx context.Context, req req.DiyTemplatePageReq) (*core.PageResult[*resp.DiyTemplateResp], error)
	GetDiyTemplateProperty(ctx context.Context, id int64) (string, error)
	UpdateDiyTemplateProperty(ctx context.Context, req req.DiyTemplatePropertyUpdateReq) error
}

type diyTemplateService struct {
	q *query.Query
}

func NewDiyTemplateService(q *query.Query) DiyTemplateService {
	return &diyTemplateService{q: q}
}

func (s *diyTemplateService) CreateDiyTemplate(ctx context.Context, req req.DiyTemplateCreateReq) (int64, error) {
	// Name duplicate check? Usually not strict for templates.
	template := &promotion.PromotionDiyTemplate{
		Name:         req.Name,
		CoverImage:   req.CoverImage,
		PreviewImage: req.PreviewImage,
		Status:       req.Status,
		Property:     req.Property,
		Sort:         req.Sort,
		Remark:       req.Remark,
	}
	err := s.q.PromotionDiyTemplate.WithContext(ctx).Create(template)
	return template.ID, err
}

func (s *diyTemplateService) UpdateDiyTemplate(ctx context.Context, req req.DiyTemplateUpdateReq) error {
	_, err := s.validateDiyTemplateExists(ctx, req.ID)
	if err != nil {
		return err
	}

	_, err = s.q.PromotionDiyTemplate.WithContext(ctx).Where(s.q.PromotionDiyTemplate.ID.Eq(req.ID)).Updates(promotion.PromotionDiyTemplate{
		Name:         req.Name,
		CoverImage:   req.CoverImage,
		PreviewImage: req.PreviewImage,
		Status:       req.Status,
		Property:     req.Property,
		Sort:         req.Sort,
		Remark:       req.Remark,
	})
	return err
}

func (s *diyTemplateService) DeleteDiyTemplate(ctx context.Context, id int64) error {
	_, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return err
	}
	// Check if used by Pages?
	count, err := s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.TemplateID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return core.NewBizError(400, "该模板已被页面使用，无法删除")
	}

	_, err = s.q.PromotionDiyTemplate.WithContext(ctx).Where(s.q.PromotionDiyTemplate.ID.Eq(id)).Delete()
	return err
}

func (s *diyTemplateService) GetDiyTemplate(ctx context.Context, id int64) (*resp.DiyTemplateResp, error) {
	template, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.convertDiyTemplateToResp(template), nil
}

func (s *diyTemplateService) GetDiyTemplatePage(ctx context.Context, req req.DiyTemplatePageReq) (*core.PageResult[*resp.DiyTemplateResp], error) {
	q := s.q.PromotionDiyTemplate
	do := q.WithContext(ctx)
	if req.Name != "" {
		do = do.Where(q.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		do = do.Where(q.Status.Eq(*req.Status))
	}

	list, total, err := do.Order(q.Sort.Asc(), q.ID.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	result := make([]*resp.DiyTemplateResp, len(list))
	for i, item := range list {
		result[i] = s.convertDiyTemplateToResp(item)
	}
	return &core.PageResult[*resp.DiyTemplateResp]{List: result, Total: total}, nil
}

func (s *diyTemplateService) GetDiyTemplateProperty(ctx context.Context, id int64) (string, error) {
	template, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return "", err
	}
	return template.Property, nil
}

// UseDiyTemplate 使用装修模板
// Java: DiyTemplateServiceImpl#useDiyTemplate
// NOTE: Java 使用 used/usedTime 字段，Go Model 当前未包含此字段
// 简化实现：仅校验存在性，实际的使用状态由前端维护
func (s *diyTemplateService) UseDiyTemplate(ctx context.Context, id int64) error {
	// 校验存在
	_, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return err
	}
	// TODO: 完整实现需要在 PromotionDiyTemplate Model 中添加 Used/UsedAt 字段
	// 参考 Java: DiyTemplateDO.used, DiyTemplateDO.usedTime
	return nil
}

// UpdateDiyTemplateProperty 更新装修模板属性
// Java: DiyTemplateServiceImpl#updateDiyTemplateProperty
func (s *diyTemplateService) UpdateDiyTemplateProperty(ctx context.Context, req req.DiyTemplatePropertyUpdateReq) error {
	// 校验存在
	_, err := s.validateDiyTemplateExists(ctx, req.ID)
	if err != nil {
		return err
	}
	// 更新属性
	_, err = s.q.PromotionDiyTemplate.WithContext(ctx).
		Where(s.q.PromotionDiyTemplate.ID.Eq(req.ID)).
		Updates(promotion.PromotionDiyTemplate{
			Property: req.Property,
		})
	return err
}

// Helpers

func (s *diyTemplateService) validateDiyTemplateExists(ctx context.Context, id int64) (*promotion.PromotionDiyTemplate, error) {
	template, err := s.q.PromotionDiyTemplate.WithContext(ctx).Where(s.q.PromotionDiyTemplate.ID.Eq(id)).First()
	if err != nil {
		return nil, core.NewBizError(404, "装修模板不存在")
	}
	return template, nil
}

func (s *diyTemplateService) convertDiyTemplateToResp(item *promotion.PromotionDiyTemplate) *resp.DiyTemplateResp {
	return &resp.DiyTemplateResp{
		ID:           item.ID,
		Name:         item.Name,
		CoverImage:   item.CoverImage,
		PreviewImage: item.PreviewImage,
		Status:       item.Status,
		Property:     item.Property,
		Sort:         item.Sort,
		Remark:       item.Remark,
		CreateTime:   item.CreateTime,
	}
}
