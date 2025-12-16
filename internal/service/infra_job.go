package service

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

// JobStatus 任务状态
const (
	JobStatusInit   = 0 // 初始化
	JobStatusNormal = 1 // 开启
	JobStatusStop   = 2 // 暂停
)

type JobService struct {
	q         *query.Query
	scheduler *Scheduler
}

func NewJobService(q *query.Query, scheduler *Scheduler) *JobService {
	return &JobService{q: q, scheduler: scheduler}
}

// CreateJob 创建定时任务
func (s *JobService) CreateJob(ctx context.Context, r *req.JobSaveReq) (int64, error) {
	job := &model.InfraJob{
		Name:           r.Name,
		Status:         JobStatusInit,
		HandlerName:    r.HandlerName,
		HandlerParam:   r.HandlerParam,
		CronExpression: r.CronExpression,
		RetryCount:     r.RetryCount,
		RetryInterval:  r.RetryInterval,
		MonitorTimeout: r.MonitorTimeout,
	}
	if err := s.q.InfraJob.WithContext(ctx).Create(job); err != nil {
		return 0, err
	}
	return job.ID, nil
}

// UpdateJob 更新定时任务
func (s *JobService) UpdateJob(ctx context.Context, r *req.JobSaveReq) error {
	if r.ID == nil {
		return errors.New("任务 ID 不能为空")
	}
	_, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(*r.ID)).Updates(map[string]interface{}{
		"name":            r.Name,
		"handler_name":    r.HandlerName,
		"handler_param":   r.HandlerParam,
		"cron_expression": r.CronExpression,
		"retry_count":     r.RetryCount,
		"retry_interval":  r.RetryInterval,
		"monitor_timeout": r.MonitorTimeout,
	})
	if err != nil {
		return err
	}
	// Reschedule if scheduler exists
	if s.scheduler != nil {
		_ = s.scheduler.RemoveJob(*r.ID)
		_ = s.scheduler.AddJob(ctx, *r.ID)
	}
	return nil
}

// DeleteJob 删除定时任务
func (s *JobService) DeleteJob(ctx context.Context, id int64) error {
	if s.scheduler != nil {
		_ = s.scheduler.RemoveJob(id)
	}
	_, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(id)).Delete()
	return err
}

// GetJob 获取定时任务
func (s *JobService) GetJob(ctx context.Context, id int64) (*model.InfraJob, error) {
	return s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(id)).First()
}

// GetJobPage 获取定时任务分页
func (s *JobService) GetJobPage(ctx context.Context, r *req.JobPageReq) (*core.PageResult[*model.InfraJob], error) {
	q := s.q.InfraJob.WithContext(ctx)

	if r.Name != "" {
		q = q.Where(s.q.InfraJob.Name.Like("%" + r.Name + "%"))
	}
	if r.HandlerName != "" {
		q = q.Where(s.q.InfraJob.HandlerName.Like("%" + r.HandlerName + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.InfraJob.Status.Eq(*r.Status))
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

	list, err := q.Order(s.q.InfraJob.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*model.InfraJob]{
		List:  list,
		Total: total,
	}, nil
}

// UpdateJobStatus 更新定时任务状态
func (s *JobService) UpdateJobStatus(ctx context.Context, id int64, status int) error {
	_, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(id)).Update(s.q.InfraJob.Status, status)
	if err != nil {
		return err
	}
	if s.scheduler != nil {
		return s.scheduler.UpdateJobStatus(ctx, id, status)
	}
	return nil
}

// TriggerJob 触发定时任务
func (s *JobService) TriggerJob(ctx context.Context, id int64) error {
	if s.scheduler != nil {
		return s.scheduler.TriggerJob(ctx, id)
	}
	return errors.New("调度器未初始化")
}
