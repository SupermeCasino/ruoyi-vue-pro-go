package iot

import (
	"context"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type OtaFirmwareService struct {
	otaFirmwareRepo OtaFirmwareRepository
}

func NewOtaFirmwareService(otaFirmwareRepo OtaFirmwareRepository) *OtaFirmwareService {
	return &OtaFirmwareService{
		otaFirmwareRepo: otaFirmwareRepo,
	}
}

func (s *OtaFirmwareService) Create(ctx context.Context, r *iot2.IotOtaFirmwareSaveReqVO) (int64, error) {
	firmware := &model.IotOtaFirmwareDO{
		Name:                r.Name,
		Description:         r.Description,
		Version:             r.Version,
		ProductID:           r.ProductID,
		FileURL:             r.FileURL,
		FileSize:            r.FileSize,
		FileDigestAlgorithm: r.FileDigestAlgorithm,
		FileDigestValue:     r.FileDigestValue,
	}
	if err := s.otaFirmwareRepo.Create(ctx, firmware); err != nil {
		return 0, err
	}
	return firmware.ID, nil
}

func (s *OtaFirmwareService) Update(ctx context.Context, r *iot2.IotOtaFirmwareSaveReqVO) error {
	f, err := s.otaFirmwareRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if f == nil {
		return errors.NewBizError(1050003000, "固件不存在")
	}
	f.Name = r.Name
	f.Description = r.Description
	f.Version = r.Version
	f.ProductID = r.ProductID
	f.FileURL = r.FileURL
	f.FileSize = r.FileSize
	f.FileDigestAlgorithm = r.FileDigestAlgorithm
	f.FileDigestValue = r.FileDigestValue
	return s.otaFirmwareRepo.Update(ctx, f)
}

func (s *OtaFirmwareService) Delete(ctx context.Context, id int64) error {
	return s.otaFirmwareRepo.Delete(ctx, id)
}

func (s *OtaFirmwareService) Get(ctx context.Context, id int64) (*model.IotOtaFirmwareDO, error) {
	return s.otaFirmwareRepo.GetByID(ctx, id)
}

func (s *OtaFirmwareService) GetPage(ctx context.Context, r *iot2.IotOtaFirmwarePageReqVO) (*pagination.PageResult[*model.IotOtaFirmwareDO], error) {
	return s.otaFirmwareRepo.GetPage(ctx, r)
}
