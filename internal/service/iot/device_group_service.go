package iot

import (
	"context"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type DeviceGroupService struct {
	deviceGroupRepo DeviceGroupRepository
}

func NewDeviceGroupService(deviceGroupRepo DeviceGroupRepository) *DeviceGroupService {
	return &DeviceGroupService{
		deviceGroupRepo: deviceGroupRepo,
	}
}

func (s *DeviceGroupService) Create(ctx context.Context, r *iot2.IotDeviceGroupSaveReqVO) (int64, error) {
	group := &model.IotDeviceGroupDO{
		Name:        r.Name,
		Status:      r.Status,
		Description: r.Description,
	}
	if err := s.deviceGroupRepo.Create(ctx, group); err != nil {
		return 0, err
	}
	return group.ID, nil
}

func (s *DeviceGroupService) Update(ctx context.Context, r *iot2.IotDeviceGroupSaveReqVO) error {
	group, err := s.deviceGroupRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if group == nil {
		return errors.NewBizError(1050001000, "设备分组不存在")
	}
	group.Name = r.Name
	group.Status = r.Status
	group.Description = r.Description
	return s.deviceGroupRepo.Update(ctx, group)
}

func (s *DeviceGroupService) Delete(ctx context.Context, id int64) error {
	return s.deviceGroupRepo.Delete(ctx, id)
}

func (s *DeviceGroupService) Get(ctx context.Context, id int64) (*model.IotDeviceGroupDO, error) {
	return s.deviceGroupRepo.GetByID(ctx, id)
}

func (s *DeviceGroupService) GetPage(ctx context.Context, r *iot2.IotDeviceGroupPageReqVO) (*pagination.PageResult[*model.IotDeviceGroupDO], error) {
	return s.deviceGroupRepo.GetPage(ctx, r)
}

// GetListByStatus 获取指定状态的设备分组列表
func (s *DeviceGroupService) GetListByStatus(ctx context.Context, status int8) ([]*model.IotDeviceGroupDO, error) {
	return s.deviceGroupRepo.ListByStatus(ctx, status)
}
