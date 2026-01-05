package iot

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

type DevicePropertyService struct {
	devicePropertyRepo DevicePropertyRepository
}

func NewDevicePropertyService(devicePropertyRepo DevicePropertyRepository) *DevicePropertyService {
	return &DevicePropertyService{
		devicePropertyRepo: devicePropertyRepo,
	}
}

func (s *DevicePropertyService) GetLatestDeviceProperties(ctx context.Context, deviceID int64) (map[string]*model.IotDevicePropertyDO, error) {
	list, err := s.devicePropertyRepo.GetLatestProperties(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// 由于 repository 返回的是按时间倒序的所有属性，我们需要在内存中取每个 identifier 的第一条
	res := make(map[string]*model.IotDevicePropertyDO)
	for _, p := range list {
		if _, ok := res[p.Identifier]; !ok {
			res[p.Identifier] = p
		}
	}
	return res, nil
}

func (s *DevicePropertyService) GetHistoryDevicePropertyList(ctx context.Context, req *iot.IotDevicePropertyHistoryListReqVO) ([]*model.IotDevicePropertyDO, error) {
	return s.devicePropertyRepo.GetHistoryList(ctx, req)
}

// SaveDeviceProperty 保存设备属性（用于属性上报）
func (s *DevicePropertyService) SaveDeviceProperty(ctx context.Context, deviceID int64, properties map[string]any) error {
	now := time.Now()
	for identifier, value := range properties {
		valueStr := ""
		switch v := value.(type) {
		case string:
			valueStr = v
		default:
			// 其他类型转为 JSON
			if b, err := json.Marshal(v); err == nil {
				valueStr = string(b)
			}
		}

		property := &model.IotDevicePropertyDO{
			DeviceID:   deviceID,
			Identifier: identifier,
			Value:      valueStr,
			UpdateTime: &now,
		}
		if err := s.devicePropertyRepo.SaveProperty(ctx, property); err != nil {
			return err
		}
	}
	return nil
}

// GetDeviceServerId 获取设备连接的网关服务器 ID
func (s *DevicePropertyService) GetDeviceServerId(ctx context.Context, deviceID int64) string {
	// 从属性中获取 serverId（如果之前存储过）
	props, err := s.GetLatestDeviceProperties(ctx, deviceID)
	if err != nil {
		return ""
	}
	if prop, ok := props["serverId"]; ok {
		return prop.Value
	}
	return ""
}

// UpdateDeviceServerId 更新设备连接的网关服务器 ID
func (s *DevicePropertyService) UpdateDeviceServerId(ctx context.Context, deviceID int64, serverID string) error {
	now := time.Now()
	property := &model.IotDevicePropertyDO{
		DeviceID:   deviceID,
		Identifier: "serverId",
		Value:      serverID,
		UpdateTime: &now,
	}
	return s.devicePropertyRepo.SaveProperty(ctx, property)
}
