package product

import (
	"context"

	product2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type ProductPropertyService struct {
	q            *query.Query
	valueService *ProductPropertyValueService
	skuSvc       *ProductSkuService
}

func NewProductPropertyService(q *query.Query, valueService *ProductPropertyValueService) *ProductPropertyService {
	return &ProductPropertyService{
		q:            q,
		valueService: valueService,
	}
}

func (s *ProductPropertyService) SetSkuService(skuSvc *ProductSkuService) {
	s.skuSvc = skuSvc
}

// CreateProperty 创建属性项
func (s *ProductPropertyService) CreateProperty(ctx context.Context, req *product2.ProductPropertyCreateReq) (int64, error) {
	// 校验名字重复
	u := s.q.ProductProperty
	exist, err := u.WithContext(ctx).Where(u.Name.Eq(req.Name)).First()
	if err == nil && exist != nil {
		return exist.ID, nil // 如果已存在，直接返回 ID (Java 逻辑)
	}

	property := &product.ProductProperty{
		Name:   req.Name,
		Remark: req.Remark,
	}
	err = u.WithContext(ctx).Create(property)
	return property.ID, err
}

// UpdateProperty 更新属性项
func (s *ProductPropertyService) UpdateProperty(ctx context.Context, req *product2.ProductPropertyUpdateReq) error {
	// 校验存在
	if err := s.validatePropertyExists(ctx, req.ID); err != nil {
		return err
	}
	// 校验名字重复
	u := s.q.ProductProperty
	exist, err := u.WithContext(ctx).Where(u.Name.Eq(req.Name)).First()
	if err == nil && exist != nil && exist.ID != req.ID {
		return errors.NewBizError(1006002000, "属性项名称已存在") // PROPERTY_EXISTS
	}

	_, err = u.WithContext(ctx).Where(u.ID.Eq(req.ID)).Updates(&product.ProductProperty{
		Name:   req.Name,
		Remark: req.Remark,
	})
	if err != nil {
		return err
	}

	// 更新 SPU 相关属性
	if s.skuSvc != nil {
		_, _ = s.skuSvc.UpdateSkuProperty(ctx, req.ID, req.Name)
	}
	return nil
}

// DeleteProperty 删除属性项
func (s *ProductPropertyService) DeleteProperty(ctx context.Context, id int64) error {
	// 校验存在
	if err := s.validatePropertyExists(ctx, id); err != nil {
		return err
	}
	// 校验其下是否有属性值
	count := s.valueService.GetPropertyValueCountByPropertyId(ctx, id)
	if count > 0 {
		return errors.NewBizError(1006002002, "属性项下存在属性值，无法删除") // PROPERTY_DELETE_FAIL_VALUE_EXISTS
	}

	// 删除
	_, err := s.q.ProductProperty.WithContext(ctx).Where(s.q.ProductProperty.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	// 同步删除属性值 (Java: deletePropertyValueByPropertyId)
	return s.valueService.DeletePropertyValueByPropertyID(ctx, id)
}

// GetProperty 获得属性项
func (s *ProductPropertyService) GetProperty(ctx context.Context, id int64) (*product2.ProductPropertyResp, error) {
	u := s.q.ProductProperty
	property, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, nil // Or return error if strict
	}
	return s.convertResp(property), nil
}

// GetPropertyPage 获得属性项分页
func (s *ProductPropertyService) GetPropertyPage(ctx context.Context, req *product2.ProductPropertyPageReq) (*pagination.PageResult[*product2.ProductPropertyResp], error) {
	u := s.q.ProductProperty
	q := u.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}

	list, total, err := q.Order(u.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	resList := lo.Map(list, func(item *product.ProductProperty, _ int) *product2.ProductPropertyResp {
		return s.convertResp(item)
	})
	return &pagination.PageResult[*product2.ProductPropertyResp]{
		List:  resList,
		Total: total,
	}, nil
}

// GetPropertyList 获得属性项列表
func (s *ProductPropertyService) GetPropertyList(ctx context.Context, req *product2.ProductPropertyListReq) ([]*product2.ProductPropertyResp, error) {
	u := s.q.ProductProperty
	q := u.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}
	list, err := q.Order(u.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *product.ProductProperty, _ int) *product2.ProductPropertyResp {
		return s.convertResp(item)
	}), nil
}

// GetPropertyListByIds 获得属性项列表 (按 ID)
func (s *ProductPropertyService) GetPropertyListByIds(ctx context.Context, ids []int64) ([]*product2.ProductPropertyResp, error) {
	if len(ids) == 0 {
		return []*product2.ProductPropertyResp{}, nil
	}
	u := s.q.ProductProperty
	list, err := u.WithContext(ctx).Where(u.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *product.ProductProperty, _ int) *product2.ProductPropertyResp {
		return s.convertResp(item)
	}), nil
}

func (s *ProductPropertyService) validatePropertyExists(ctx context.Context, id int64) error {
	u := s.q.ProductProperty
	count, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.NewBizError(1006002001, "属性项不存在") // PROPERTY_NOT_EXISTS
	}
	return nil
}

func (s *ProductPropertyService) convertResp(item *product.ProductProperty) *product2.ProductPropertyResp {
	return &product2.ProductPropertyResp{
		ID:         item.ID,
		Name:       item.Name,
		Remark:     item.Remark,
		CreateTime: item.CreateTime,
	}
}
