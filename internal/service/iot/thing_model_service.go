package iot

import (
	"context"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

type ThingModelService struct {
	productRepo    ProductRepository
	thingModelRepo ThingModelRepository
}

func NewThingModelService(productRepo ProductRepository, thingModelRepo ThingModelRepository) *ThingModelService {
	return &ThingModelService{
		productRepo:    productRepo,
		thingModelRepo: thingModelRepo,
	}
}

func (s *ThingModelService) Create(ctx context.Context, r *iot2.IotThingModelSaveReqVO) (int64, error) {
	product, err := s.productRepo.GetByID(ctx, r.ProductID)
	if err != nil {
		return 0, err
	}
	if product == nil {
		return 0, model.ErrProductNotExists
	}
	if product.Status == consts.IotProductStatusPublished {
		return 0, model.ErrProductStatusNotAllowThingModel
	}

	thingModel := &model.IotThingModelDO{
		Identifier:  r.Identifier,
		Name:        r.Name,
		Description: r.Description,
		ProductID:   r.ProductID,
		ProductKey:  product.ProductKey,
		Type:        r.Type,
		Property:    datatypes.JSON(r.Property),
		Event:       datatypes.JSON(r.Event),
		Service:     datatypes.JSON(r.Service),
	}
	if err := s.thingModelRepo.Create(ctx, thingModel); err != nil {
		return 0, err
	}
	return thingModel.ID, nil
}

func (s *ThingModelService) Update(ctx context.Context, r *iot2.IotThingModelSaveReqVO) error {
	tm, err := s.thingModelRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if tm == nil {
		return model.ErrThingModelNotExists
	}

	product, err := s.productRepo.GetByID(ctx, tm.ProductID)
	if err != nil {
		return err
	}
	if product != nil && product.Status == consts.IotProductStatusPublished {
		return model.ErrProductStatusNotAllowThingModel
	}

	tm.Name = r.Name
	tm.Description = r.Description
	tm.Type = r.Type
	tm.Property = datatypes.JSON(r.Property)
	tm.Event = datatypes.JSON(r.Event)
	tm.Service = datatypes.JSON(r.Service)

	return s.thingModelRepo.Update(ctx, tm)
}

func (s *ThingModelService) Delete(ctx context.Context, id int64) error {
	tm, err := s.thingModelRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if tm == nil {
		return nil
	}
	product, err := s.productRepo.GetByID(ctx, tm.ProductID)
	if err != nil {
		return err
	}
	if product != nil && product.Status == consts.IotProductStatusPublished {
		return model.ErrProductStatusNotAllowThingModel
	}
	return s.thingModelRepo.Delete(ctx, id)
}

func (s *ThingModelService) Get(ctx context.Context, id int64) (*model.IotThingModelDO, error) {
	return s.thingModelRepo.GetByID(ctx, id)
}

func (s *ThingModelService) GetList(ctx context.Context, r *iot2.IotThingModelListReqVO) ([]*model.IotThingModelDO, error) {
	return s.thingModelRepo.ListByProductID(ctx, r.ProductID)
}

func (s *ThingModelService) GetPage(ctx context.Context, r *iot2.IotThingModelPageReqVO) (*pagination.PageResult[*model.IotThingModelDO], error) {
	return s.thingModelRepo.GetPage(ctx, r)
}

func (s *ThingModelService) GetTSL(ctx context.Context, productId int64) (*iot2.IotThingModelTSLRespVO, error) {
	product, err := s.productRepo.GetByID(ctx, productId)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, model.ErrProductNotExists
	}

	thingModels, err := s.thingModelRepo.ListByProductID(ctx, productId)
	if err != nil {
		return nil, err
	}

	tsl := &iot2.IotThingModelTSLRespVO{
		ProductID:  product.ID,
		ProductKey: product.ProductKey,
	}

	for _, m := range thingModels {
		switch m.Type {
		case consts.IotThingModelTypeProperty: // 属性
			tsl.Properties = append(tsl.Properties, string(m.Property))
		case consts.IotThingModelTypeService: // 服务
			tsl.Services = append(tsl.Services, string(m.Service))
		case consts.IotThingModelTypeEvent: // 事件
			tsl.Events = append(tsl.Events, string(m.Event))
		}
	}
	return tsl, nil
}

func (s *ThingModelService) GetThingModelListByProductIdAndType(ctx context.Context, productId int64, tmType int8) ([]*model.IotThingModelDO, error) {
	return s.thingModelRepo.ListByProductIDAndType(ctx, productId, tmType)
}
