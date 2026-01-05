package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type OtaTaskRepositoryImpl struct {
	q *query.Query
}

func NewOtaTaskRepository(q *query.Query) iotsvc.OtaTaskRepository {
	return &OtaTaskRepositoryImpl{q: q}
}

func (r *OtaTaskRepositoryImpl) Create(ctx context.Context, task *model.IotOtaTaskDO) error {
	return r.q.IotOtaTaskDO.WithContext(ctx).Create(task)
}

func (r *OtaTaskRepositoryImpl) Update(ctx context.Context, task *model.IotOtaTaskDO) error {
	_, err := r.q.IotOtaTaskDO.WithContext(ctx).Where(r.q.IotOtaTaskDO.ID.Eq(task.ID)).Updates(task)
	return err
}

func (r *OtaTaskRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotOtaTaskDO.WithContext(ctx).Where(r.q.IotOtaTaskDO.ID.Eq(id)).Delete()
	return err
}

func (r *OtaTaskRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotOtaTaskDO, error) {
	return r.q.IotOtaTaskDO.WithContext(ctx).Where(r.q.IotOtaTaskDO.ID.Eq(id)).First()
}

func (r *OtaTaskRepositoryImpl) GetPage(ctx context.Context, req *iot.IotOtaTaskPageReqVO) (*pagination.PageResult[*model.IotOtaTaskDO], error) {
	t := r.q.IotOtaTaskDO
	db := t.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(t.Name.Like("%" + req.Name + "%"))
	}
	list, total, err := db.Order(t.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotOtaTaskDO]{List: list, Total: total}, err
}
