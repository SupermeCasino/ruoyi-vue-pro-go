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

type ProductPropertyValueService struct {
	q      *query.Query
	skuSvc *ProductSkuService
}

func NewProductPropertyValueService(q *query.Query) *ProductPropertyValueService {
	return &ProductPropertyValueService{q: q}
}

func (s *ProductPropertyValueService) SetSkuService(skuSvc *ProductSkuService) {
	s.skuSvc = skuSvc
}

// CreatePropertyValue 创建属性值
func (s *ProductPropertyValueService) CreatePropertyValue(ctx context.Context, req *product2.ProductPropertyValueCreateReq) (int64, error) {
	u := s.q.ProductPropertyValue
	// 如果已经添加过该属性值，直接返回
	exist, err := u.WithContext(ctx).Where(u.PropertyID.Eq(req.PropertyID), u.Name.Eq(req.Name)).First()
	if err == nil && exist != nil {
		return exist.ID, nil
	}

	value := &product.ProductPropertyValue{
		PropertyID: req.PropertyID,
		Name:       req.Name,
		Remark:     req.Remark,
	}
	err = u.WithContext(ctx).Create(value)
	return value.ID, err
}

// UpdatePropertyValue 更新属性值
func (s *ProductPropertyValueService) UpdatePropertyValue(ctx context.Context, req *product2.ProductPropertyValueUpdateReq) error {
	if err := s.validatePropertyValueExists(ctx, req.ID); err != nil {
		return err
	}

	// 校验名字唯一
	u := s.q.ProductPropertyValue
	exist, err := u.WithContext(ctx).Where(u.PropertyID.Eq(req.PropertyID), u.Name.Eq(req.Name)).First()
	if err == nil && exist != nil && exist.ID != req.ID {
		return errors.NewBizError(1006002003, "属性值名称已存在") // PROPERTY_VALUE_EXISTS
	}

	_, err = u.WithContext(ctx).Where(u.ID.Eq(req.ID)).Updates(&product.ProductPropertyValue{
		PropertyID: req.PropertyID,
		Name:       req.Name,
		Remark:     req.Remark,
	})
	if err != nil {
		return err
	}
	if s.skuSvc != nil {
		_, _ = s.skuSvc.UpdateSkuPropertyValue(ctx, req.ID, req.Name)
	}
	return nil
}

// DeletePropertyValue 删除属性值
func (s *ProductPropertyValueService) DeletePropertyValue(ctx context.Context, id int64) error {
	if err := s.validatePropertyValueExists(ctx, id); err != nil {
		return err
	}
	_, err := s.q.ProductPropertyValue.WithContext(ctx).Where(s.q.ProductPropertyValue.ID.Eq(id)).Delete()
	return err
}

// GetPropertyValue 获得属性值
func (s *ProductPropertyValueService) GetPropertyValue(ctx context.Context, id int64) (*product2.ProductPropertyValueResp, error) {
	u := s.q.ProductPropertyValue
	value, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, nil
	}
	return s.convertResp(value), nil
}

// GetPropertyValuePage 获得属性值分页
func (s *ProductPropertyValueService) GetPropertyValuePage(ctx context.Context, req *product2.ProductPropertyValuePageReq) (*pagination.PageResult[*product2.ProductPropertyValueResp], error) {
	u := s.q.ProductPropertyValue
	q := u.WithContext(ctx)
	if req.PropertyID != 0 {
		q = q.Where(u.PropertyID.Eq(req.PropertyID))
	}
	if req.Name != "" {
		q = q.Where(u.Name.Like("%" + req.Name + "%"))
	}

	list, total, err := q.Order(u.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	resList := lo.Map(list, func(item *product.ProductPropertyValue, _ int) *product2.ProductPropertyValueResp {
		return s.convertResp(item)
	})
	return &pagination.PageResult[*product2.ProductPropertyValueResp]{
		List:  resList,
		Total: total,
	}, nil
}

func (s *ProductPropertyValueService) GetPropertyValueCountByPropertyId(ctx context.Context, propertyId int64) int64 {
	u := s.q.ProductPropertyValue
	count, _ := u.WithContext(ctx).Where(u.PropertyID.Eq(propertyId)).Count()
	return count
}

func (s *ProductPropertyValueService) DeletePropertyValueByPropertyID(ctx context.Context, propertyId int64) error {
	u := s.q.ProductPropertyValue
	_, err := u.WithContext(ctx).Where(u.PropertyID.Eq(propertyId)).Delete()
	return err
}

func (s *ProductPropertyValueService) validatePropertyValueExists(ctx context.Context, id int64) error {
	u := s.q.ProductPropertyValue
	count, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.NewBizError(1006002004, "属性值不存在") // PROPERTY_VALUE_NOT_EXISTS
	}
	return nil
}

func (s *ProductPropertyValueService) GetPropertyValueListByPropertyIds(ctx context.Context, propertyIds []int64) ([]*product.ProductPropertyValue, error) {
	if len(propertyIds) == 0 {
		return []*product.ProductPropertyValue{}, nil
	}
	u := s.q.ProductPropertyValue
	return u.WithContext(ctx).Where(u.PropertyID.In(propertyIds...)).Find()
}

func (s *ProductPropertyValueService) convertResp(item *product.ProductPropertyValue) *product2.ProductPropertyValueResp {
	return &product2.ProductPropertyValueResp{
		ID:         item.ID,
		PropertyID: item.PropertyID,
		Name:       item.Name,
		Remark:     item.Remark,
		CreateTime: item.CreateTime,
	}
}
