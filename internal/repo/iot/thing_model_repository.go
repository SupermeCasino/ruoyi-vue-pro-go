package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ThingModelRepositoryImpl struct {
	q *query.Query
}

func NewThingModelRepository(q *query.Query) iotsvc.ThingModelRepository {
	return &ThingModelRepositoryImpl{q: q}
}

func (r *ThingModelRepositoryImpl) Create(ctx context.Context, tm *model.IotThingModelDO) error {
	return r.q.IotThingModelDO.WithContext(ctx).Create(tm)
}

func (r *ThingModelRepositoryImpl) Update(ctx context.Context, tm *model.IotThingModelDO) error {
	_, err := r.q.IotThingModelDO.WithContext(ctx).Where(r.q.IotThingModelDO.ID.Eq(tm.ID)).Updates(tm)
	return err
}

func (r *ThingModelRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotThingModelDO.WithContext(ctx).Where(r.q.IotThingModelDO.ID.Eq(id)).Delete()
	return err
}

func (r *ThingModelRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotThingModelDO, error) {
	return r.q.IotThingModelDO.WithContext(ctx).Where(r.q.IotThingModelDO.ID.Eq(id)).First()
}

func (r *ThingModelRepositoryImpl) ListByProductID(ctx context.Context, productID int64) ([]*model.IotThingModelDO, error) {
	return r.q.IotThingModelDO.WithContext(ctx).Where(r.q.IotThingModelDO.ProductID.Eq(productID)).Find()
}

func (r *ThingModelRepositoryImpl) ListByProductIDAndType(ctx context.Context, productID int64, tmType int8) ([]*model.IotThingModelDO, error) {
	return r.q.IotThingModelDO.WithContext(ctx).Where(r.q.IotThingModelDO.ProductID.Eq(productID), r.q.IotThingModelDO.Type.Eq(tmType)).Find()
}

func (r *ThingModelRepositoryImpl) GetPage(ctx context.Context, req *iot.IotThingModelPageReqVO) (*pagination.PageResult[*model.IotThingModelDO], error) {
	tm := r.q.IotThingModelDO
	list, total, err := tm.WithContext(ctx).Where(tm.ProductID.Eq(req.ProductID)).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotThingModelDO]{List: list, Total: total}, err
}
