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

type ProductCategoryRepositoryImpl struct {
	q *query.Query
}

func NewProductCategoryRepository(q *query.Query) iotsvc.ProductCategoryRepository {
	return &ProductCategoryRepositoryImpl{q: q}
}

func (r *ProductCategoryRepositoryImpl) Create(ctx context.Context, category *model.IotProductCategoryDO) error {
	return r.q.IotProductCategoryDO.WithContext(ctx).Create(category)
}

func (r *ProductCategoryRepositoryImpl) Update(ctx context.Context, category *model.IotProductCategoryDO) error {
	_, err := r.q.IotProductCategoryDO.WithContext(ctx).Where(r.q.IotProductCategoryDO.ID.Eq(category.ID)).Updates(category)
	return err
}

func (r *ProductCategoryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotProductCategoryDO.WithContext(ctx).Where(r.q.IotProductCategoryDO.ID.Eq(id)).Delete()
	return err
}

func (r *ProductCategoryRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotProductCategoryDO, error) {
	return r.q.IotProductCategoryDO.WithContext(ctx).Where(r.q.IotProductCategoryDO.ID.Eq(id)).First()
}

func (r *ProductCategoryRepositoryImpl) GetPage(ctx context.Context, req *iot.IotProductCategoryPageReqVO) (*pagination.PageResult[*model.IotProductCategoryDO], error) {
	pc := r.q.IotProductCategoryDO
	db := pc.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(pc.Name.Like("%" + req.Name + "%"))
	}
	list, total, err := db.Order(pc.Sort.Asc(), pc.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotProductCategoryDO]{List: list, Total: total}, err
}

func (r *ProductCategoryRepositoryImpl) GetListByStatus(ctx context.Context, status int8) ([]*model.IotProductCategoryDO, error) {
	pc := r.q.IotProductCategoryDO
	return pc.WithContext(ctx).Where(pc.Status.Eq(status)).Order(pc.Sort.Asc(), pc.ID.Desc()).Find()
}

func (r *ProductCategoryRepositoryImpl) Count(ctx context.Context, startTime *time.Time) (int64, error) {
	pc := r.q.IotProductCategoryDO
	db := pc.WithContext(ctx)
	if startTime != nil {
		db = db.Where(pc.CreateTime.Gte(*startTime))
	}
	return db.Count()
}

func (r *ProductCategoryRepositoryImpl) GetProductCategoryDeviceCountMap(ctx context.Context) (map[string]int64, error) {
	pc := r.q.IotProductCategoryDO
	prod := r.q.IotProductDO
	dev := r.q.IotDeviceDO

	var results []struct {
		CategoryName string
		DeviceCount  int64
	}

	err := dev.WithContext(ctx).
		LeftJoin(prod, prod.ID.EqCol(dev.ProductID)).
		LeftJoin(pc, pc.ID.EqCol(prod.CategoryID)).
		Select(pc.Name.As("category_name"), dev.ID.Count().As("device_count")).
		Group(pc.Name).
		Scan(&results)

	if err != nil {
		return nil, err
	}

	res := make(map[string]int64)
	for _, r := range results {
		if r.CategoryName != "" {
			res[r.CategoryName] = r.DeviceCount
		}
	}
	return res, nil
}
