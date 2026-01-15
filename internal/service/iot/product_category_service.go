package iot

import (
	"context"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ProductCategoryService struct {
	productCategoryRepo ProductCategoryRepository
}

func NewProductCategoryService(productCategoryRepo ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{
		productCategoryRepo: productCategoryRepo,
	}
}

func (s *ProductCategoryService) CreateProductCategory(ctx context.Context, r *iot2.IotProductCategorySaveReqVO) (int64, error) {
	category := &model.IotProductCategoryDO{
		Name:        r.Name,
		Sort:        r.Sort,
		Status:      r.Status,
		Description: r.Description,
	}
	if err := s.productCategoryRepo.Create(ctx, category); err != nil {
		return 0, err
	}
	return category.ID, nil
}

func (s *ProductCategoryService) UpdateProductCategory(ctx context.Context, r *iot2.IotProductCategorySaveReqVO) error {
	category, err := s.productCategoryRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if category == nil {
		return model.ErrProductCategoryNotExists
	}

	category.Name = r.Name
	category.Sort = r.Sort
	category.Status = r.Status
	category.Description = r.Description

	return s.productCategoryRepo.Update(ctx, category)
}

func (s *ProductCategoryService) DeleteProductCategory(ctx context.Context, id int64) error {
	category, err := s.productCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return model.ErrProductCategoryNotExists
	}
	return s.productCategoryRepo.Delete(ctx, id)
}

func (s *ProductCategoryService) GetProductCategory(ctx context.Context, id int64) (*model.IotProductCategoryDO, error) {
	return s.productCategoryRepo.GetByID(ctx, id)
}

func (s *ProductCategoryService) GetProductCategoryPage(ctx context.Context, r *iot2.IotProductCategoryPageReqVO) (*pagination.PageResult[*model.IotProductCategoryDO], error) {
	return s.productCategoryRepo.GetPage(ctx, r)
}

func (s *ProductCategoryService) GetProductCategoryListByStatus(ctx context.Context, status int8) ([]*model.IotProductCategoryDO, error) {
	return s.productCategoryRepo.GetListByStatus(ctx, status)
}

func (s *ProductCategoryService) GetProductCategoryListByIDs(ctx context.Context, ids []int64) ([]*model.IotProductCategoryDO, error) {
	return s.productCategoryRepo.GetListByIDs(ctx, ids)
}

func (s *ProductCategoryService) GetProductCategoryCount(ctx context.Context) (int64, error) {
	return s.productCategoryRepo.Count(ctx, nil)
}
