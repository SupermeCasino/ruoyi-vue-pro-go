package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type OtaTaskRecordRepositoryImpl struct {
	q *query.Query
}

func NewOtaTaskRecordRepository(q *query.Query) iotsvc.OtaTaskRecordRepository {
	return &OtaTaskRecordRepositoryImpl{q: q}
}

func (r *OtaTaskRecordRepositoryImpl) Create(ctx context.Context, records []*model.IotOtaTaskRecordDO) error {
	return r.q.IotOtaTaskRecordDO.WithContext(ctx).Create(records...)
}

func (r *OtaTaskRecordRepositoryImpl) GetPage(ctx context.Context, req *iot.IotOtaTaskRecordPageReqVO) (*pagination.PageResult[*model.IotOtaTaskRecordDO], error) {
	tr := r.q.IotOtaTaskRecordDO
	db := tr.WithContext(ctx).Where(tr.TaskID.Eq(req.TaskID))
	if req.Status != 0 {
		db = db.Where(tr.Status.Eq(req.Status))
	}
	list, total, err := db.Order(tr.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotOtaTaskRecordDO]{List: list, Total: total}, err
}

func (r *OtaTaskRecordRepositoryImpl) CreateBatch(ctx context.Context, records []*model.IotOtaTaskRecordDO) error {
	return r.q.IotOtaTaskRecordDO.WithContext(ctx).Create(records...)
}

// Update 更新 OTA 升级记录
func (r *OtaTaskRecordRepositoryImpl) Update(ctx context.Context, record *model.IotOtaTaskRecordDO) error {
	_, err := r.q.IotOtaTaskRecordDO.WithContext(ctx).Where(r.q.IotOtaTaskRecordDO.ID.Eq(record.ID)).Updates(record)
	return err
}

// GetListByDeviceIdAndStatus 根据设备ID和状态列表查询升级记录
func (r *OtaTaskRecordRepositoryImpl) GetListByDeviceIdAndStatus(ctx context.Context, deviceID int64, statuses []int) ([]*model.IotOtaTaskRecordDO, error) {
	tr := r.q.IotOtaTaskRecordDO
	// 将 int 切片转换为 int8 切片
	statusInt8 := make([]int8, len(statuses))
	for i, s := range statuses {
		statusInt8[i] = int8(s)
	}
	return tr.WithContext(ctx).Where(tr.DeviceID.Eq(deviceID), tr.Status.In(statusInt8...)).Find()
}

// GetListByTaskIdAndStatus 根据任务ID和状态列表查询升级记录
func (r *OtaTaskRecordRepositoryImpl) GetListByTaskIdAndStatus(ctx context.Context, taskID int64, statuses []int) ([]*model.IotOtaTaskRecordDO, error) {
	tr := r.q.IotOtaTaskRecordDO
	statusInt8 := make([]int8, len(statuses))
	for i, s := range statuses {
		statusInt8[i] = int8(s)
	}
	return tr.WithContext(ctx).Where(tr.TaskID.Eq(taskID), tr.Status.In(statusInt8...)).Find()
}
