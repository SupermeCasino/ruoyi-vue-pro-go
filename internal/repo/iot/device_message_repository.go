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

type DeviceMessageRepositoryImpl struct {
	q *query.Query
}

func NewDeviceMessageRepository(q *query.Query) iotsvc.DeviceMessageRepository {
	return &DeviceMessageRepositoryImpl{q: q}
}

// Create 创建设备消息记录
func (r *DeviceMessageRepositoryImpl) Create(ctx context.Context, message *model.IotDeviceMessageDO) error {
	return r.q.IotDeviceMessageDO.WithContext(ctx).Create(message)
}

func (r *DeviceMessageRepositoryImpl) Count(ctx context.Context, startTime *time.Time) (int64, error) {
	m := r.q.IotDeviceMessageDO
	db := m.WithContext(ctx)
	if startTime != nil {
		db = db.Where(m.CreateTime.Gte(*startTime))
	}
	return db.Count()
}

func (r *DeviceMessageRepositoryImpl) GetSummaryByDate(ctx context.Context, interval int, startTime, endTime *time.Time) ([]*iot.IotStatisticsDeviceMessageSummaryByDateRespVO, error) {
	m := r.q.IotDeviceMessageDO
	db := m.WithContext(ctx)
	if startTime != nil {
		db = db.Where(m.CreateTime.Gte(*startTime))
	}
	if endTime != nil {
		db = db.Where(m.CreateTime.Lte(*endTime))
	}

	// 这里的 interval 逻辑暂时简化处理。Java 中通常是按天、周、月分组。
	// 这里我们默认按天分组，格式 'YYYY-MM-DD'

	var results []struct {
		Time            string
		UpstreamCount   int64
		DownstreamCount int64
	}

	err := db.UnderlyingDB().Select(
		"DATE_FORMAT(create_time, '%Y%m%d') AS time",
		"SUM(CASE WHEN upstream = 1 THEN 1 ELSE 0 END) AS upstream_count",
		"SUM(CASE WHEN upstream = 0 THEN 1 ELSE 0 END) AS downstream_count",
	).Group("time").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	res := make([]*iot.IotStatisticsDeviceMessageSummaryByDateRespVO, 0, len(results))
	for _, r := range results {
		res = append(res, &iot.IotStatisticsDeviceMessageSummaryByDateRespVO{
			Time:            r.Time,
			UpstreamCount:   r.UpstreamCount,
			DownstreamCount: r.DownstreamCount,
		})
	}
	return res, nil
}
func (r *DeviceMessageRepositoryImpl) GetPage(ctx context.Context, req *iot.IotDeviceMessagePageReqVO) (*pagination.PageResult[*model.IotDeviceMessageDO], error) {
	m := r.q.IotDeviceMessageDO
	db := m.WithContext(ctx).Where(m.DeviceID.Eq(req.DeviceID))

	if req.Method != "" {
		db = db.Where(m.Method.Eq(req.Method))
	}
	if req.Upstream != nil {
		db = db.Where(m.Upstream.Is(*req.Upstream))
	}
	if req.Reply != nil {
		db = db.Where(m.Reply.Is(*req.Reply))
	}
	if req.Identifier != "" {
		db = db.Where(m.Identifier.Eq(req.Identifier))
	}
	if len(req.Times) == 2 {
		if req.Times[0] != nil {
			db = db.Where(m.CreateTime.Gte(*req.Times[0]))
		}
		if req.Times[1] != nil {
			db = db.Where(m.CreateTime.Lte(*req.Times[1]))
		}
	}

	list, total, err := db.Order(m.CreateTime.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotDeviceMessageDO]{List: list, Total: total}, err
}

func (r *DeviceMessageRepositoryImpl) GetListByRequestIdsAndReply(ctx context.Context, deviceID int64, requestIDs []string, reply bool) ([]*model.IotDeviceMessageDO, error) {
	m := r.q.IotDeviceMessageDO
	return m.WithContext(ctx).
		Where(m.DeviceID.Eq(deviceID), m.RequestID.In(requestIDs...), m.Reply.Is(reply)).
		Find()
}
