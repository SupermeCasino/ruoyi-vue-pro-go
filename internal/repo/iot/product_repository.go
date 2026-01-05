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

type ProductRepositoryImpl struct {
	q *query.Query
}

func NewProductRepository(q *query.Query) iotsvc.ProductRepository {
	return &ProductRepositoryImpl{q: q}
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *model.IotProductDO) error {
	return r.q.IotProductDO.WithContext(ctx).Create(product)
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *model.IotProductDO) error {
	_, err := r.q.IotProductDO.WithContext(ctx).Where(r.q.IotProductDO.ID.Eq(product.ID)).Updates(product)
	return err
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotProductDO.WithContext(ctx).Where(r.q.IotProductDO.ID.Eq(id)).Delete()
	return err
}

func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotProductDO, error) {
	return r.q.IotProductDO.WithContext(ctx).Where(r.q.IotProductDO.ID.Eq(id)).First()
}

func (r *ProductRepositoryImpl) GetByKey(ctx context.Context, key string) (*model.IotProductDO, error) {
	return r.q.IotProductDO.WithContext(ctx).Where(r.q.IotProductDO.ProductKey.Eq(key)).First()
}

func (r *ProductRepositoryImpl) GetPage(ctx context.Context, req *iot.IotProductPageReqVO) (*pagination.PageResult[*model.IotProductDO], error) {
	p := r.q.IotProductDO
	db := p.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(p.Name.Like("%" + req.Name + "%"))
	}
	if req.ProductKey != "" {
		db = db.Where(p.ProductKey.Eq(req.ProductKey))
	}
	list, total, err := db.Order(p.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotProductDO]{List: list, Total: total}, err
}

func (r *ProductRepositoryImpl) ListAll(ctx context.Context) ([]*model.IotProductDO, error) {
	return r.q.IotProductDO.WithContext(ctx).Find()
}

func (r *ProductRepositoryImpl) Count(ctx context.Context, startTime *time.Time) (int64, error) {
	p := r.q.IotProductDO
	db := p.WithContext(ctx)
	if startTime != nil {
		db = db.Where(p.CreateTime.Gte(*startTime))
	}
	return db.Count()
}
