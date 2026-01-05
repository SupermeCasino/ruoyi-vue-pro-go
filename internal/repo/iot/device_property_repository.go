package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
)

type DevicePropertyRepositoryImpl struct {
	q *query.Query
}

func NewDevicePropertyRepository(q *query.Query) iotsvc.DevicePropertyRepository {
	return &DevicePropertyRepositoryImpl{q: q}
}

func (r *DevicePropertyRepositoryImpl) GetHistoryList(ctx context.Context, req *iot.IotDevicePropertyHistoryListReqVO) ([]*model.IotDevicePropertyDO, error) {
	m := r.q.IotDevicePropertyDO
	db := m.WithContext(ctx).Where(m.DeviceID.Eq(req.DeviceID), m.Identifier.Eq(req.Identifier))

	if len(req.Times) == 2 {
		if req.Times[0] != nil {
			db = db.Where(m.UpdateTime.Gte(*req.Times[0]))
		}
		if req.Times[1] != nil {
			db = db.Where(m.UpdateTime.Lte(*req.Times[1]))
		}
	}

	return db.Order(m.UpdateTime.Desc()).Find()
}

func (r *DevicePropertyRepositoryImpl) GetLatestProperties(ctx context.Context, deviceID int64) ([]*model.IotDevicePropertyDO, error) {
	m := r.q.IotDevicePropertyDO
	// 简化方案：查询该设备的所有属性，按时间倒序排列，后续在 service 层取每个 identifier 的第一条
	return m.WithContext(ctx).Where(m.DeviceID.Eq(deviceID)).Order(m.UpdateTime.Desc()).Find()
}

func (r *DevicePropertyRepositoryImpl) SaveProperty(ctx context.Context, property *model.IotDevicePropertyDO) error {
	return r.q.IotDevicePropertyDO.WithContext(ctx).Create(property)
}
