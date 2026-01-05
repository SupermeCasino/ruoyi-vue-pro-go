package iot

import (
	"context"
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

type DataSinkRepositoryImpl struct {
	q *query.Query
}

func NewDataSinkRepository(q *query.Query) iotsvc.DataSinkRepository {
	return &DataSinkRepositoryImpl{q: q}
}

func (r *DataSinkRepositoryImpl) Create(ctx context.Context, sink *model.IotDataSinkDO) error {
	return r.q.IotDataSinkDO.WithContext(ctx).Create(sink)
}

func (r *DataSinkRepositoryImpl) Update(ctx context.Context, sink *model.IotDataSinkDO) error {
	_, err := r.q.IotDataSinkDO.WithContext(ctx).Where(r.q.IotDataSinkDO.ID.Eq(sink.ID)).Updates(sink)
	return err
}

func (r *DataSinkRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotDataSinkDO.WithContext(ctx).Where(r.q.IotDataSinkDO.ID.Eq(id)).Delete()
	return err
}

func (r *DataSinkRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotDataSinkDO, error) {
	return r.q.IotDataSinkDO.WithContext(ctx).Where(r.q.IotDataSinkDO.ID.Eq(id)).First()
}

func (r *DataSinkRepositoryImpl) GetPage(ctx context.Context, req *iot.IotDataSinkPageReqVO) (*pagination.PageResult[*model.IotDataSinkDO], error) {
	ds := r.q.IotDataSinkDO
	db := ds.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(ds.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != 0 {
		db = db.Where(ds.Status.Eq(req.Status))
	}
	list, total, err := db.Order(ds.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotDataSinkDO]{List: list, Total: total}, err
}

func (r *DataSinkRepositoryImpl) CountBySinkID(ctx context.Context, id int64) (int64, error) {
	return r.q.IotDataRuleDO.WithContext(ctx).Where(r.q.IotDataRuleDO.SinkIDs.Like(datatypes.JSON("%" + strconv.FormatInt(id, 10) + "%"))).Count()
}

func (r *DataSinkRepositoryImpl) GetListByStatus(ctx context.Context, status int8) ([]*model.IotDataSinkDO, error) {
	ds := r.q.IotDataSinkDO
	return ds.WithContext(ctx).Where(ds.Status.Eq(status)).Order(ds.ID.Desc()).Find()
}
