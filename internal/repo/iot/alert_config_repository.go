package iot

import (
	"context"
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type AlertConfigRepositoryImpl struct {
	q *query.Query
}

func NewAlertConfigRepository(q *query.Query) iotsvc.AlertConfigRepository {
	return &AlertConfigRepositoryImpl{q: q}
}

func (r *AlertConfigRepositoryImpl) Create(ctx context.Context, config *model.IotAlertConfigDO) error {
	return r.q.IotAlertConfigDO.WithContext(ctx).Create(config)
}

func (r *AlertConfigRepositoryImpl) Update(ctx context.Context, config *model.IotAlertConfigDO) error {
	_, err := r.q.IotAlertConfigDO.WithContext(ctx).Where(r.q.IotAlertConfigDO.ID.Eq(config.ID)).Updates(config)
	return err
}

func (r *AlertConfigRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotAlertConfigDO.WithContext(ctx).Where(r.q.IotAlertConfigDO.ID.Eq(id)).Delete()
	return err
}

func (r *AlertConfigRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotAlertConfigDO, error) {
	return r.q.IotAlertConfigDO.WithContext(ctx).Where(r.q.IotAlertConfigDO.ID.Eq(id)).First()
}

func (r *AlertConfigRepositoryImpl) GetPage(ctx context.Context, req *iot.IotAlertConfigPageReqVO) (*pagination.PageResult[*model.IotAlertConfigDO], error) {
	ac := r.q.IotAlertConfigDO
	db := ac.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(ac.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != 0 {
		db = db.Where(ac.Status.Eq(req.Status))
	}
	list, total, err := db.Order(ac.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotAlertConfigDO]{List: list, Total: total}, err
}

func (r *AlertConfigRepositoryImpl) GetListByStatus(ctx context.Context, status int8) ([]*model.IotAlertConfigDO, error) {
	ac := r.q.IotAlertConfigDO
	return ac.WithContext(ctx).Where(ac.Status.Eq(status)).Find()
}

func (r *AlertConfigRepositoryImpl) GetListBySceneRuleIdAndStatus(ctx context.Context, sceneRuleID int64, status int8) ([]*model.IotAlertConfigDO, error) {
	ac := r.q.IotAlertConfigDO
	var list []*model.IotAlertConfigDO
	err := ac.WithContext(ctx).Where(ac.Status.Eq(status)).UnderlyingDB().
		Where("JSON_CONTAINS(scene_rule_ids, ?)", strconv.FormatInt(sceneRuleID, 10)).
		Find(&list).Error
	return list, err
}
