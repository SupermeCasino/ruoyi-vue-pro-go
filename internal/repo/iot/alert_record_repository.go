package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type AlertRecordRepositoryImpl struct {
	q *query.Query
}

func NewAlertRecordRepository(q *query.Query) iotsvc.AlertRecordRepository {
	return &AlertRecordRepositoryImpl{q: q}
}

func (r *AlertRecordRepositoryImpl) Create(ctx context.Context, record *model.IotAlertRecordDO) error {
	return r.q.IotAlertRecordDO.WithContext(ctx).Create(record)
}

func (r *AlertRecordRepositoryImpl) Update(ctx context.Context, record *model.IotAlertRecordDO) error {
	_, err := r.q.IotAlertRecordDO.WithContext(ctx).Where(r.q.IotAlertRecordDO.ID.Eq(record.ID)).Updates(record)
	return err
}

func (r *AlertRecordRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotAlertRecordDO, error) {
	return r.q.IotAlertRecordDO.WithContext(ctx).Where(r.q.IotAlertRecordDO.ID.Eq(id)).First()
}

func (r *AlertRecordRepositoryImpl) GetPage(ctx context.Context, req *iot.IotAlertRecordPageReqVO) (*pagination.PageResult[*model.IotAlertRecordDO], error) {
	ar := r.q.IotAlertRecordDO
	db := ar.WithContext(ctx)
	if req.ConfigID != 0 {
		db = db.Where(ar.ConfigID.Eq(req.ConfigID))
	}
	if req.ProcessStatus != nil {
		db = db.Where(ar.ProcessStatus.Is(*req.ProcessStatus))
	}
	list, total, err := db.Order(ar.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotAlertRecordDO]{List: list, Total: total}, err
}

func (r *AlertRecordRepositoryImpl) GetListBySceneRuleId(ctx context.Context, sceneRuleID int64, deviceID *int64, processStatus *bool) ([]*model.IotAlertRecordDO, error) {
	ar := r.q.IotAlertRecordDO
	db := ar.WithContext(ctx).Where(ar.SceneRuleID.Eq(sceneRuleID))
	if deviceID != nil {
		db = db.Where(ar.DeviceID.Eq(*deviceID))
	}
	if processStatus != nil {
		db = db.Where(ar.ProcessStatus.Is(*processStatus))
	}
	return db.Find()
}
