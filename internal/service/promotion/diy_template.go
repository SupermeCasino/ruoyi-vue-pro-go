package promotion

import (
	"context"
	stdErrors "errors"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/types"
	"gorm.io/gorm"
)

type DiyTemplateService interface {
	CreateDiyTemplate(ctx context.Context, req req.DiyTemplateCreateReq) (int64, error)
	UpdateDiyTemplate(ctx context.Context, req req.DiyTemplateUpdateReq) error
	UseDiyTemplate(ctx context.Context, id int64) error
	DeleteDiyTemplate(ctx context.Context, id int64) error
	GetDiyTemplate(ctx context.Context, id int64) (*resp.DiyTemplateResp, error)
	GetDiyTemplatePage(ctx context.Context, req req.DiyTemplatePageReq) (*pagination.PageResult[*resp.DiyTemplateResp], error)
	GetDiyTemplateProperty(ctx context.Context, id int64) (string, error)
	UpdateDiyTemplateProperty(ctx context.Context, req req.DiyTemplatePropertyUpdateReq) error
	GetUsedDiyTemplate(ctx context.Context) (*promotion.PromotionDiyTemplate, error)
}

type diyTemplateService struct {
	q *query.Query
}

func NewDiyTemplateService(q *query.Query) DiyTemplateService {
	return &diyTemplateService{q: q}
}

func (s *diyTemplateService) CreateDiyTemplate(ctx context.Context, req req.DiyTemplateCreateReq) (int64, error) {
	if err := s.validateNameUnique(ctx, 0, req.Name); err != nil {
		return 0, err
	}
	template := &promotion.PromotionDiyTemplate{
		Name:           req.Name,
		PreviewPicUrls: types.StringListFromCSV(req.PreviewPicUrls),
		Property:       req.Property,
		Remark:         req.Remark,
		Used:           false,
	}
	err := s.q.PromotionDiyTemplate.WithContext(ctx).Create(template)
	if err != nil {
		return 0, err
	}
	// 创建默认页面
	if err := s.createDefaultPage(ctx, template.ID); err != nil {
		return 0, err
	}
	return template.ID, nil
}

func (s *diyTemplateService) UpdateDiyTemplate(ctx context.Context, req req.DiyTemplateUpdateReq) error {
	_, err := s.validateDiyTemplateExists(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.validateNameUnique(ctx, req.ID, req.Name); err != nil {
		return err
	}

	_, err = s.q.PromotionDiyTemplate.WithContext(ctx).Where(s.q.PromotionDiyTemplate.ID.Eq(req.ID)).Updates(promotion.PromotionDiyTemplate{
		Name:           req.Name,
		PreviewPicUrls: types.StringListFromCSV(req.PreviewPicUrls),
		Property:       req.Property,
		Remark:         req.Remark,
	})
	return err
}

func (s *diyTemplateService) DeleteDiyTemplate(ctx context.Context, id int64) error {
	template, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return err
	}
	if template.Used {
		return errors.NewBizError(400, "该模板正在使用，无法删除")
	}

	_, err = s.q.PromotionDiyTemplate.WithContext(ctx).Where(s.q.PromotionDiyTemplate.ID.Eq(id)).Delete()
	// Logic to delete pages associated?
	// Java deletes pages too? Yes, usually cascade or manual delete.
	// Java doesn't show explicit page delete in controller snippet but usually Service has it.
	// We should probably delete pages too.
	_, err = s.q.PromotionDiyPage.WithContext(ctx).Where(s.q.PromotionDiyPage.TemplateID.Eq(id)).Delete()
	return err
}

func (s *diyTemplateService) GetDiyTemplate(ctx context.Context, id int64) (*resp.DiyTemplateResp, error) {
	template, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.convertDiyTemplateToResp(template), nil
}

func (s *diyTemplateService) GetDiyTemplatePage(ctx context.Context, req req.DiyTemplatePageReq) (*pagination.PageResult[*resp.DiyTemplateResp], error) {
	q := s.q.PromotionDiyTemplate
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

	result := make([]*resp.DiyTemplateResp, len(list))
	for i, item := range list {
		result[i] = s.convertDiyTemplateToResp(item)
	}
	return &pagination.PageResult[*resp.DiyTemplateResp]{List: result, Total: total}, nil
}

func (s *diyTemplateService) GetDiyTemplateProperty(ctx context.Context, id int64) (string, error) {
	template, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return "", err
	}
	return template.Property, nil
}

// UseDiyTemplate 使用装修模板
func (s *diyTemplateService) UseDiyTemplate(ctx context.Context, id int64) error {
	// 校验存在
	_, err := s.validateDiyTemplateExists(ctx, id)
	if err != nil {
		return err
	}

	// 开启事务
	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 将所有已使用的设置为未使用
		err := tx.PromotionDiyTemplate.WithContext(ctx).UnderlyingDB().Model(&promotion.PromotionDiyTemplate{}).Where("used = ?", true).Updates(map[string]interface{}{"used": false}).Error
		if err != nil {
			return err
		}

		// 2. 更新新的为使用
		now := time.Now()
		_, err = tx.PromotionDiyTemplate.WithContext(ctx).
			Where(tx.PromotionDiyTemplate.ID.Eq(id)).
			Updates(map[string]interface{}{
				"used":      true,
				"used_time": &now,
			})
		return err
	})
}

func (s *diyTemplateService) GetUsedDiyTemplate(ctx context.Context) (*promotion.PromotionDiyTemplate, error) {
	template := &promotion.PromotionDiyTemplate{}
	err := s.q.PromotionDiyTemplate.WithContext(ctx).UnderlyingDB().Where("used = ?", true).First(template).Error
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return template, nil
}

// UpdateDiyTemplateProperty 更新装修模板属性
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
func (s *diyTemplateService) createDefaultPage(ctx context.Context, templateID int64) error {
	pages := []*promotion.PromotionDiyPage{
		{
			TemplateID: templateID,
			Name:       "首页",
			Remark:     "默认首页",
			Property:   "{}",
		},
		{
			TemplateID: templateID,
			Name:       "我的",
			Remark:     "默认我的页面",
			Property:   "{}",
		},
	}
	return s.q.PromotionDiyPage.WithContext(ctx).Create(pages...)
}

func (s *diyTemplateService) validateNameUnique(ctx context.Context, id int64, name string) error {
	q := s.q.PromotionDiyTemplate
	do := q.WithContext(ctx).Where(q.Name.Eq(name))
	if id > 0 {
		do = do.Where(q.ID.Neq(id))
	}
	count, err := do.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.NewBizError(400, "模板名称已存在")
	}
	return nil
}

func (s *diyTemplateService) validateDiyTemplateExists(ctx context.Context, id int64) (*promotion.PromotionDiyTemplate, error) {
	template, err := s.q.PromotionDiyTemplate.WithContext(ctx).Where(s.q.PromotionDiyTemplate.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(404, "装修模板不存在")
	}
	return template, nil
}

func (s *diyTemplateService) convertDiyTemplateToResp(item *promotion.PromotionDiyTemplate) *resp.DiyTemplateResp {
	return &resp.DiyTemplateResp{
		ID:             item.ID,
		Name:           item.Name,
		PreviewPicUrls: []string(item.PreviewPicUrls),
		Property:       item.Property,
		Used:           bool(item.Used),
		UsedTime:       item.UsedTime,
		Remark:         item.Remark,
		CreateTime:     item.CreateTime,
	}
}
