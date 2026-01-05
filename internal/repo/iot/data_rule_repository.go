package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type DataRuleRepositoryImpl struct {
	q *query.Query
}

func NewDataRuleRepository(q *query.Query) iotsvc.DataRuleRepository {
	return &DataRuleRepositoryImpl{q: q}
}

func (r *DataRuleRepositoryImpl) Create(ctx context.Context, rule *model.IotDataRuleDO) error {
	return r.q.IotDataRuleDO.WithContext(ctx).Create(rule)
}

func (r *DataRuleRepositoryImpl) Update(ctx context.Context, rule *model.IotDataRuleDO) error {
	_, err := r.q.IotDataRuleDO.WithContext(ctx).Where(r.q.IotDataRuleDO.ID.Eq(rule.ID)).Updates(rule)
	return err
}

func (r *DataRuleRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotDataRuleDO.WithContext(ctx).Where(r.q.IotDataRuleDO.ID.Eq(id)).Delete()
	return err
}

func (r *DataRuleRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotDataRuleDO, error) {
	return r.q.IotDataRuleDO.WithContext(ctx).Where(r.q.IotDataRuleDO.ID.Eq(id)).First()
}

func (r *DataRuleRepositoryImpl) GetPage(ctx context.Context, req *iot.IotDataRulePageReqVO) (*pagination.PageResult[*model.IotDataRuleDO], error) {
	dr := r.q.IotDataRuleDO
	db := dr.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(dr.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != 0 {
		db = db.Where(dr.Status.Eq(req.Status))
	}
	list, total, err := db.Order(dr.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotDataRuleDO]{List: list, Total: total}, err
}
