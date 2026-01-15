package iot

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type DeviceRepositoryImpl struct {
	q *query.Query
}

func NewDeviceRepository(q *query.Query) iotsvc.DeviceRepository {
	return &DeviceRepositoryImpl{q: q}
}

func (r *DeviceRepositoryImpl) Create(ctx context.Context, device *model.IotDeviceDO) error {
	return r.q.IotDeviceDO.WithContext(ctx).Create(device)
}

func (r *DeviceRepositoryImpl) Update(ctx context.Context, device *model.IotDeviceDO) error {
	_, err := r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ID.Eq(device.ID)).Updates(device)
	return err
}

func (r *DeviceRepositoryImpl) UpdateActiveTime(ctx context.Context, id int64, activeTime time.Time) error {
	_, err := r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ID.Eq(id)).Update(r.q.IotDeviceDO.ActiveTime, activeTime)
	return err
}

func (r *DeviceRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ID.Eq(id)).Delete()
	return err
}

func (r *DeviceRepositoryImpl) DeleteList(ctx context.Context, ids []int64) error {
	_, err := r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ID.In(ids...)).Delete()
	return err
}

func (r *DeviceRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotDeviceDO, error) {
	return r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ID.Eq(id)).First()
}

func (r *DeviceRepositoryImpl) GetPage(ctx context.Context, req *iot.IotDevicePageReqVO) (*pagination.PageResult[*model.IotDeviceDO], error) {
	d := r.q.IotDeviceDO
	db := d.WithContext(ctx)
	if req.DeviceName != "" {
		db = db.Where(d.DeviceName.Like("%" + req.DeviceName + "%"))
	}
	if req.Nickname != "" {
		db = db.Where(d.Nickname.Like("%" + req.Nickname + "%"))
	}
	if req.ProductID != 0 {
		db = db.Where(d.ProductID.Eq(req.ProductID))
	}
	if req.DeviceType != 0 {
		db = db.Where(d.DeviceType.Eq(req.DeviceType))
	}
	if req.Status != 0 {
		db = db.Where(d.State.Eq(req.Status))
	}
	list, total, err := db.Order(d.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotDeviceDO]{List: list, Total: total}, err
}

func (r *DeviceRepositoryImpl) CountByProductID(ctx context.Context, productID int64) (int64, error) {
	return r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ProductID.Eq(productID)).Count()
}

func (r *DeviceRepositoryImpl) CountByGatewayID(ctx context.Context, gatewayID int64) (int64, error) {
	return r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.GatewayID.Eq(gatewayID)).Count()
}

func (r *DeviceRepositoryImpl) Count(ctx context.Context, startTime *time.Time) (int64, error) {
	d := r.q.IotDeviceDO
	db := d.WithContext(ctx)
	if startTime != nil {
		db = db.Where(d.CreateTime.Gte(*startTime))
	}
	return db.Count()
}

func (r *DeviceRepositoryImpl) GetStateCountMap(ctx context.Context) (map[int8]int64, error) {
	d := r.q.IotDeviceDO
	var results []struct {
		State int8
		Count int64
	}
	err := d.WithContext(ctx).Select(d.State, d.State.Count().As("count")).Group(d.State).Scan(&results)
	if err != nil {
		return nil, err
	}
	res := make(map[int8]int64)
	for _, r := range results {
		res[r.State] = r.Count
	}
	return res, nil
}

func (r *DeviceRepositoryImpl) ListByCondition(ctx context.Context, deviceType *int8, productID *int64) ([]*model.IotDeviceDO, error) {
	d := r.q.IotDeviceDO
	db := d.WithContext(ctx)
	if deviceType != nil {
		db = db.Where(d.DeviceType.Eq(*deviceType))
	}
	if productID != nil {
		db = db.Where(d.ProductID.Eq(*productID))
	}
	return db.Find()
}

func (r *DeviceRepositoryImpl) ListByProductKeyAndNames(ctx context.Context, productKey string, names []string) ([]*model.IotDeviceDO, error) {
	return r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ProductKey.Eq(productKey), r.q.IotDeviceDO.DeviceName.In(names...)).Find()
}

func (r *DeviceRepositoryImpl) GetByProductKeyAndName(ctx context.Context, productKey string, name string) (*model.IotDeviceDO, error) {
	return r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.ProductKey.Eq(productKey), r.q.IotDeviceDO.DeviceName.Eq(name)).First()
}
