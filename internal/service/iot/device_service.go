package iot

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	iotcore "github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

type DeviceService struct {
	productRepo ProductRepository
	deviceRepo  DeviceRepository
	authUtils   *iotcore.DeviceAuthUtils
}

func NewDeviceService(productRepo ProductRepository, deviceRepo DeviceRepository, authUtils *iotcore.DeviceAuthUtils) *DeviceService {
	return &DeviceService{
		productRepo: productRepo,
		deviceRepo:  deviceRepo,
		authUtils:   authUtils,
	}
}

func (s *DeviceService) Create(ctx context.Context, r *iot2.IotDeviceSaveReqVO) (int64, error) {
	product, err := s.productRepo.GetByID(ctx, r.ProductID)
	if err != nil {
		return 0, err
	}
	if product == nil {
		return 0, model.ErrProductNotExists
	}

	exists, _ := s.deviceRepo.GetByProductKeyAndName(ctx, product.ProductKey, r.DeviceName)
	if exists != nil {
		return 0, model.ErrDeviceNameExists
	}

	groupIdsJson, _ := json.Marshal(r.GroupIDs)
	device := &model.IotDeviceDO{
		DeviceName:   r.DeviceName,
		Nickname:     r.Nickname,
		SerialNumber: r.SerialNumber,
		PicURL:       r.PicURL,
		GroupIDs:     datatypes.JSON(groupIdsJson),
		ProductID:    r.ProductID,
		ProductKey:   product.ProductKey,
		DeviceType:   product.DeviceType,
		GatewayID:    r.GatewayID,
		State:        0,
		DeviceSecret: uuid.New().String(),
		Config:       datatypes.JSON(r.Config),
		LocationType: r.LocationType,
		Latitude:     r.Latitude,
		Longitude:    r.Longitude,
	}
	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return 0, err
	}
	return device.ID, nil
}

func (s *DeviceService) Update(ctx context.Context, r *iot2.IotDeviceSaveReqVO) error {
	device, err := s.deviceRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if device == nil {
		return model.ErrDeviceNotExists
	}

	groupIdsJson, _ := json.Marshal(r.GroupIDs)
	device.Nickname = r.Nickname
	device.SerialNumber = r.SerialNumber
	device.PicURL = r.PicURL
	device.GroupIDs = datatypes.JSON(groupIdsJson)
	device.GatewayID = r.GatewayID
	device.Config = datatypes.JSON(r.Config)
	device.LocationType = r.LocationType
	device.Latitude = r.Latitude
	device.Longitude = r.Longitude

	return s.deviceRepo.Update(ctx, device)
}

func (s *DeviceService) Delete(ctx context.Context, id int64) error {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if device == nil {
		return model.ErrDeviceNotExists
	}
	// Gateway sub-device check
	if device.DeviceType == consts.IotProductDeviceTypeGateway { // 网关
		count, err := s.deviceRepo.CountByGatewayID(ctx, id)
		if err != nil {
			return err
		}
		if count > 0 {
			return model.ErrDeviceHasChildren
		}
	}
	return s.deviceRepo.Delete(ctx, id)
}

func (s *DeviceService) DeleteList(ctx context.Context, ids []int64) error {
	return s.deviceRepo.DeleteList(ctx, ids)
}

func (s *DeviceService) Get(ctx context.Context, id int64) (*model.IotDeviceDO, error) {
	return s.deviceRepo.GetByID(ctx, id)
}

func (s *DeviceService) GetByProductKeyAndName(ctx context.Context, productKey, deviceName string) (*model.IotDeviceDO, error) {
	return s.deviceRepo.GetByProductKeyAndName(ctx, productKey, deviceName)
}

func (s *DeviceService) GetAuthInfo(ctx context.Context, id int64) (*iot2.IotDeviceAuthInfoRespVO, error) {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, model.ErrDeviceNotExists
	}

	return &iot2.IotDeviceAuthInfoRespVO{
		ProductKey:   device.ProductKey,
		DeviceName:   device.DeviceName,
		DeviceSecret: device.DeviceSecret,
		MqttHost:     "127.0.0.1", // TODO: 从配置获取
		MqttPort:     1883,        // TODO: 从配置获取
	}, nil
}

func (s *DeviceService) Auth(ctx context.Context, productKey, deviceName, password string) (*model.IotDeviceDO, error) {
	device, err := s.deviceRepo.GetByProductKeyAndName(ctx, productKey, deviceName)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, model.ErrDeviceNotExists
	}

	// 校验密码
	content := s.authUtils.BuildAuthContent(deviceName, productKey)
	expectPassword := s.authUtils.BuildPassword(device.DeviceSecret, content)
	if expectPassword != password {
		return nil, model.ErrDeviceSecretInvalid
	}

	return device, nil
}

func (s *DeviceService) GetListByProductKeyAndNames(ctx context.Context, productKey string, deviceNames []string) ([]*model.IotDeviceDO, error) {
	return s.deviceRepo.ListByProductKeyAndNames(ctx, productKey, deviceNames)
}

func (s *DeviceService) GetPage(ctx context.Context, r *iot2.IotDevicePageReqVO) (*pagination.PageResult[*model.IotDeviceDO], error) {
	return s.deviceRepo.GetPage(ctx, r)
}

func (s *DeviceService) UpdateGroup(ctx context.Context, r *iot2.IotDeviceUpdateGroupReqVO) error {
	groupIdsJson, _ := json.Marshal(r.GroupIDs)
	for _, id := range r.IDs {
		device, _ := s.deviceRepo.GetByID(ctx, id)
		if device != nil {
			device.GroupIDs = datatypes.JSON(groupIdsJson)
			s.deviceRepo.Update(ctx, device)
		}
	}
	return nil
}

// UpdateDeviceFirmware 更新设备固件版本
func (s *DeviceService) UpdateDeviceFirmware(ctx context.Context, id int64, firmwareID int64) error {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if device == nil {
		return model.ErrDeviceNotExists
	}
	device.FirmwareID = firmwareID
	return s.deviceRepo.Update(ctx, device)
}

// GetCountByProductID 获取指定产品的设备数量
func (s *DeviceService) GetCountByProductID(ctx context.Context, productID int64) (int64, error) {
	return s.deviceRepo.CountByProductID(ctx, productID)
}
