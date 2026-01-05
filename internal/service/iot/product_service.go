package iot

import (
	"context"

	"github.com/google/uuid"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ProductService struct {
	productRepo ProductRepository
	deviceRepo  DeviceRepository
}

func NewProductService(productRepo ProductRepository, deviceRepo DeviceRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		deviceRepo:  deviceRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, r *iot2.IotProductSaveReqVO) (int64, error) {
	if r.ProductKey != "" {
		exists, _ := s.productRepo.GetByKey(ctx, r.ProductKey)
		if exists != nil {
			return 0, model.ErrProductKeyExists
		}
	} else {
		r.ProductKey = uuid.New().String()[0:8]
	}

	product := &model.IotProductDO{
		Name:         r.Name,
		ProductKey:   r.ProductKey,
		CategoryID:   r.CategoryID,
		Icon:         r.Icon,
		PicURL:       r.PicURL,
		Description:  r.Description,
		Status:       consts.IotProductStatusUnpublished,
		DeviceType:   r.DeviceType,
		NetType:      r.NetType,
		LocationType: r.LocationType,
		CodecType:    r.CodecType,
	}
	if err := s.productRepo.Create(ctx, product); err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (s *ProductService) Update(ctx context.Context, r *iot2.IotProductSaveReqVO) error {
	product, err := s.productRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if product == nil {
		return model.ErrProductNotExists
	}
	if product.Status == consts.IotProductStatusPublished {
		return model.ErrProductStatusNotDelete
	}

	product.Name = r.Name
	product.CategoryID = r.CategoryID
	product.Icon = r.Icon
	product.PicURL = r.PicURL
	product.Description = r.Description
	product.DeviceType = r.DeviceType
	product.NetType = r.NetType
	product.LocationType = r.LocationType
	product.CodecType = r.CodecType

	return s.productRepo.Update(ctx, product)
}

func (s *ProductService) UpdateStatus(ctx context.Context, id int64, status int8) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return model.ErrProductNotExists
	}

	product.Status = status
	return s.productRepo.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return model.ErrProductNotExists
	}
	if product.Status == consts.IotProductStatusPublished {
		return model.ErrProductStatusNotDelete
	}
	count, err := s.deviceRepo.CountByProductID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return model.ErrProductDeleteFailHasDevice
	}

	return s.productRepo.Delete(ctx, id)
}

func (s *ProductService) Get(ctx context.Context, id int64) (*model.IotProductDO, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) GetByKey(ctx context.Context, productKey string) (*model.IotProductDO, error) {
	return s.productRepo.GetByKey(ctx, productKey)
}

func (s *ProductService) GetSimpleList(ctx context.Context) ([]*model.IotProductDO, error) {
	return s.productRepo.ListAll(ctx)
}

func (s *ProductService) GetPage(ctx context.Context, r *iot2.IotProductPageReqVO) (*pagination.PageResult[*model.IotProductDO], error) {
	return s.productRepo.GetPage(ctx, r)
}
