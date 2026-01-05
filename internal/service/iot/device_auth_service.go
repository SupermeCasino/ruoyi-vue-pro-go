package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	iotcore "github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
)

const (
	iotDeviceAuthCacheKey = "iot:device:auth:%s:%s" // iot:device:auth:{productKey}:{deviceName}
	iotDeviceCacheTTL     = time.Hour
)

type IotDeviceCommonApiImpl struct {
	authUtils  *iotcore.DeviceAuthUtils
	rdb        *redis.Client
	deviceRepo DeviceRepository
}

func NewIotDeviceCommonApiImpl(rdb *redis.Client, deviceRepo DeviceRepository) iotcore.IotDeviceCommonApi {
	return &IotDeviceCommonApiImpl{
		authUtils:  iotcore.NewDeviceAuthUtils(),
		rdb:        rdb,
		deviceRepo: deviceRepo,
	}
}

// AuthDevice 设备认证
func (api *IotDeviceCommonApiImpl) AuthDevice(ctx context.Context, req *iotcore.DeviceAuthReq) (*iotcore.DeviceAuthResp, error) {
	// 1. 解析用户名获取设备信息
	deviceInfo, err := api.authUtils.ParseUsername(req.Username)
	if err != nil {
		return &iotcore.DeviceAuthResp{Success: false, Message: err.Error()}, nil
	}

	// 2. 获取设备信息（带 Redis 缓存）
	device, err := api.GetDevice(ctx, &iotcore.DeviceGetReq{
		ProductKey: deviceInfo.ProductKey,
		DeviceName: deviceInfo.DeviceName,
	})
	if err != nil {
		return &iotcore.DeviceAuthResp{Success: false, Message: err.Error()}, nil
	}
	if device == nil {
		return &iotcore.DeviceAuthResp{Success: false, Message: "device not found"}, nil
	}

	// 3. 验证密码
	if !api.authUtils.ValidatePassword(device.DeviceSecret, deviceInfo.DeviceName, deviceInfo.ProductKey, req.Password) {
		return &iotcore.DeviceAuthResp{Success: false, Message: "invalid password"}, nil
	}

	return &iotcore.DeviceAuthResp{Success: true}, nil
}

// GetDevice 获取设备信息（带 Redis 缓存回退逻辑）
func (api *IotDeviceCommonApiImpl) GetDevice(ctx context.Context, req *iotcore.DeviceGetReq) (*iotcore.DeviceResp, error) {
	key := fmt.Sprintf(iotDeviceAuthCacheKey, req.ProductKey, req.DeviceName)

	// 1. 从 Redis 读取
	val, err := api.rdb.Get(ctx, key).Result()
	if err == nil {
		var cacheResp iotcore.DeviceResp
		if err := json.Unmarshal([]byte(val), &cacheResp); err == nil {
			return &cacheResp, nil
		}
	}

	// 2. 缓存未命中，从数据库读取
	device, err := api.deviceRepo.GetByProductKeyAndName(ctx, req.ProductKey, req.DeviceName)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, nil
	}

	// 3. 转换为 Core 模型
	resp := &iotcore.DeviceResp{
		ID:           device.ID,
		ProductKey:   device.ProductKey,
		DeviceName:   device.DeviceName,
		DeviceSecret: device.DeviceSecret,
		State:        int(device.State),
		TenantID:     device.TenantID,
	}

	// 4. 写入 Redis 缓存
	if b, err := json.Marshal(resp); err == nil {
		api.rdb.Set(ctx, key, string(b), iotDeviceCacheTTL)
	}

	return resp, nil
}
