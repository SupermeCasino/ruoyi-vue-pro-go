package iot

import (
	"context"
	"encoding/json"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

type DataSinkService struct {
	dataSinkRepo DataSinkRepository
}

func NewDataSinkService(dataSinkRepo DataSinkRepository) *DataSinkService {
	return &DataSinkService{
		dataSinkRepo: dataSinkRepo,
	}
}

func (s *DataSinkService) Create(ctx context.Context, r *iot2.IotDataSinkSaveReqVO) (int64, error) {
	config, _ := json.Marshal(r.Config)
	sink := &model.IotDataSinkDO{
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		Type:        r.Type,
		Config:      datatypes.JSON(config),
	}
	if err := s.dataSinkRepo.Create(ctx, sink); err != nil {
		return 0, err
	}
	return sink.ID, nil
}

func (s *DataSinkService) Update(ctx context.Context, r *iot2.IotDataSinkSaveReqVO) error {
	sink, err := s.dataSinkRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if sink == nil {
		return model.ErrDataSinkNotExists
	}

	config, _ := json.Marshal(r.Config)
	sink.Name = r.Name
	sink.Description = r.Description
	sink.Status = r.Status
	sink.Type = r.Type
	sink.Config = datatypes.JSON(config)

	return s.dataSinkRepo.Update(ctx, sink)
}

func (s *DataSinkService) Delete(ctx context.Context, id int64) error {
	count, err := s.dataSinkRepo.CountBySinkID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return model.ErrDataSinkUsedByRule
	}
	return s.dataSinkRepo.Delete(ctx, id)
}

func (s *DataSinkService) Get(ctx context.Context, id int64) (*model.IotDataSinkDO, error) {
	return s.dataSinkRepo.GetByID(ctx, id)
}

func (s *DataSinkService) GetPage(ctx context.Context, r *iot2.IotDataSinkPageReqVO) (*pagination.PageResult[*model.IotDataSinkDO], error) {
	return s.dataSinkRepo.GetPage(ctx, r)
}

func (s *DataSinkService) GetListByStatus(ctx context.Context, status int8) ([]*model.IotDataSinkDO, error) {
	return s.dataSinkRepo.GetListByStatus(ctx, status)
}
