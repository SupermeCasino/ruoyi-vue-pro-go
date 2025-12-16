package service

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type JobLogService struct {
	q *query.Query
}

func NewJobLogService(q *query.Query) *JobLogService {
	return &JobLogService{q: q}
}

// GetJobLog 获取定时任务日志
func (s *JobLogService) GetJobLog(ctx context.Context, id int64) (*model.InfraJobLog, error) {
	return s.q.InfraJobLog.WithContext(ctx).Where(s.q.InfraJobLog.ID.Eq(id)).First()
}

// GetJobLogPage 获取定时任务日志分页
func (s *JobLogService) GetJobLogPage(ctx context.Context, r *req.JobLogPageReq) (*pagination.PageResult[*model.InfraJobLog], error) {
	q := s.q.InfraJobLog.WithContext(ctx)

	if r.JobID != nil {
		q = q.Where(s.q.InfraJobLog.JobID.Eq(*r.JobID))
	}
	if r.HandlerName != "" {
		q = q.Where(s.q.InfraJobLog.HandlerName.Like("%" + r.HandlerName + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.InfraJobLog.Status.Eq(*r.Status))
	}
	if len(r.BeginTime) == 2 {
		q = q.Where(s.q.InfraJobLog.BeginTime.Between(r.BeginTime[0], r.BeginTime[1]))
	}

	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.InfraJobLog.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.InfraJobLog]{
		List:  list,
		Total: total,
	}, nil
}
