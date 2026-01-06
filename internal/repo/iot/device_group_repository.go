package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type DeviceGroupRepositoryImpl struct {
	q *query.Query
}

func NewDeviceGroupRepository(q *query.Query) iotsvc.DeviceGroupRepository {
	return &DeviceGroupRepositoryImpl{q: q}
}

func (r *DeviceGroupRepositoryImpl) Create(ctx context.Context, group *model.IotDeviceGroupDO) error {
	return r.q.IotDeviceGroupDO.WithContext(ctx).Create(group)
}

func (r *DeviceGroupRepositoryImpl) Update(ctx context.Context, group *model.IotDeviceGroupDO) error {
	_, err := r.q.IotDeviceGroupDO.WithContext(ctx).Where(r.q.IotDeviceGroupDO.ID.Eq(group.ID)).Updates(group)
	return err
}

func (r *DeviceGroupRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotDeviceGroupDO.WithContext(ctx).Where(r.q.IotDeviceGroupDO.ID.Eq(id)).Delete()
	return err
}

func (r *DeviceGroupRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotDeviceGroupDO, error) {
	return r.q.IotDeviceGroupDO.WithContext(ctx).Where(r.q.IotDeviceGroupDO.ID.Eq(id)).First()
}

func (r *DeviceGroupRepositoryImpl) GetPage(ctx context.Context, req *iot.IotDeviceGroupPageReqVO) (*pagination.PageResult[*model.IotDeviceGroupDO], error) {
	dg := r.q.IotDeviceGroupDO
	db := dg.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(dg.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != 0 {
		db = db.Where(dg.Status.Eq(req.Status))
	}
	list, total, err := db.Order(dg.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotDeviceGroupDO]{List: list, Total: total}, err
}

func (r *DeviceGroupRepositoryImpl) ListByStatus(ctx context.Context, status int8) ([]*model.IotDeviceGroupDO, error) {
	return r.q.IotDeviceGroupDO.WithContext(ctx).Where(r.q.IotDeviceGroupDO.Status.Eq(status)).Find()
}
