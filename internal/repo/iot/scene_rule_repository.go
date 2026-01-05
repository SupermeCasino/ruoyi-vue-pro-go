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

type SceneRuleRepositoryImpl struct {
	q *query.Query
}

func NewSceneRuleRepository(q *query.Query) iotsvc.SceneRuleRepository {
	return &SceneRuleRepositoryImpl{q: q}
}

func (r *SceneRuleRepositoryImpl) Create(ctx context.Context, rule *model.IotSceneRuleDO) error {
	return r.q.IotSceneRuleDO.WithContext(ctx).Create(rule)
}

func (r *SceneRuleRepositoryImpl) Update(ctx context.Context, rule *model.IotSceneRuleDO) error {
	_, err := r.q.IotSceneRuleDO.WithContext(ctx).Where(r.q.IotSceneRuleDO.ID.Eq(rule.ID)).Updates(rule)
	return err
}

func (r *SceneRuleRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotSceneRuleDO.WithContext(ctx).Where(r.q.IotSceneRuleDO.ID.Eq(id)).Delete()
	return err
}

func (r *SceneRuleRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotSceneRuleDO, error) {
	return r.q.IotSceneRuleDO.WithContext(ctx).Where(r.q.IotSceneRuleDO.ID.Eq(id)).First()
}

func (r *SceneRuleRepositoryImpl) GetPage(ctx context.Context, req *iot.IotSceneRulePageReqVO) (*pagination.PageResult[*model.IotSceneRuleDO], error) {
	sr := r.q.IotSceneRuleDO
	db := sr.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(sr.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != 0 {
		db = db.Where(sr.Status.Eq(req.Status))
	}
	list, total, err := db.Order(sr.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotSceneRuleDO]{List: list, Total: total}, err
}

func (r *SceneRuleRepositoryImpl) CountBySceneRuleID(ctx context.Context, id int64) (int64, error) {
	return r.q.IotDeviceDO.WithContext(ctx).Where(r.q.IotDeviceDO.GroupIDs.Like(datatypes.JSON("%" + strconv.FormatInt(id, 10) + "%"))).Count()
}

func (r *SceneRuleRepositoryImpl) GetListByStatus(ctx context.Context, status int8) ([]*model.IotSceneRuleDO, error) {
	sr := r.q.IotSceneRuleDO
	return sr.WithContext(ctx).Where(sr.Status.Eq(status)).Order(sr.ID.Desc()).Find()
}
