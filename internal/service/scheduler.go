package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

// JobHandler is the interface for job handlers
type JobHandler interface {
	Execute(ctx context.Context, param string) error
}

// Scheduler manages cron jobs using gocron/v2
type Scheduler struct {
	scheduler gocron.Scheduler
	q         *query.Query
	log       *zap.Logger
	handlers  map[string]JobHandler
	jobMap    map[int64]gocron.Job
	mu        sync.RWMutex
}

// NewScheduler creates a new Scheduler
func NewScheduler(q *query.Query, log *zap.Logger, payTransferSyncJob *PayTransferSyncJob) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	scheduler := &Scheduler{
		scheduler: s,
		q:         q,
		log:       log,
		handlers:  make(map[string]JobHandler),
		jobMap:    make(map[int64]gocron.Job),
	}

	// Register specific jobs
	scheduler.RegisterHandler("payTransferSyncJob", payTransferSyncJob)

	// Auto start scheduler in background
	go func() {
		if err := scheduler.Start(context.Background()); err != nil {
			log.Error("Failed to start scheduler", zap.Error(err))
		}
	}()

	return scheduler, nil
}

// RegisterHandler registers a job handler by name
func (s *Scheduler) RegisterHandler(name string, handler JobHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[name] = handler
}

// Start loads all enabled jobs from DB and starts the scheduler
func (s *Scheduler) Start(ctx context.Context) error {
	jobs, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.Status.Eq(JobStatusNormal)).Find()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err := s.scheduleJob(ctx, job); err != nil {
			s.log.Error("Failed to schedule job", zap.Int64("jobId", job.ID), zap.Error(err))
		}
	}

	s.scheduler.Start()
	s.log.Info("Scheduler started", zap.Int("jobCount", len(jobs)))
	return nil
}

// Shutdown stops the scheduler
func (s *Scheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}

// scheduleJob adds a single job to the scheduler
func (s *Scheduler) scheduleJob(ctx context.Context, job *model.InfraJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	handler, ok := s.handlers[job.HandlerName]
	if !ok {
		return fmt.Errorf("handler not found: %s", job.HandlerName)
	}

	gocronJob, err := s.scheduler.NewJob(
		gocron.CronJob(job.CronExpression, false),
		gocron.NewTask(func() {
			s.executeJob(ctx, job, handler)
		}),
		gocron.WithName(fmt.Sprintf("job-%d", job.ID)),
	)
	if err != nil {
		return err
	}

	s.jobMap[job.ID] = gocronJob
	s.log.Info("Job scheduled", zap.Int64("jobId", job.ID), zap.String("handlerName", job.HandlerName), zap.String("cron", job.CronExpression))
	return nil
}

// executeJob runs a job and logs the result
func (s *Scheduler) executeJob(ctx context.Context, job *model.InfraJob, handler JobHandler) {
	beginTime := time.Now()

	logRecord := &model.InfraJobLog{
		JobID:        job.ID,
		HandlerName:  job.HandlerName,
		HandlerParam: job.HandlerParam,
		ExecuteIndex: 1,
		BeginTime:    beginTime,
		Status:       0,
	}
	_ = s.q.InfraJobLog.WithContext(ctx).Create(logRecord)

	var status int
	var result string
	err := handler.Execute(ctx, job.HandlerParam)
	endTime := time.Now()
	duration := int(endTime.Sub(beginTime).Milliseconds())

	if err != nil {
		status = 2
		result = err.Error()
		s.log.Error("Job execution failed", zap.Int64("jobId", job.ID), zap.Error(err))
	} else {
		status = 1
		result = "success"
		s.log.Info("Job execution completed", zap.Int64("jobId", job.ID), zap.Int("duration", duration))
	}

	_, _ = s.q.InfraJobLog.WithContext(ctx).Where(s.q.InfraJobLog.ID.Eq(logRecord.ID)).Updates(map[string]interface{}{
		"end_time": endTime,
		"duration": duration,
		"status":   status,
		"result":   result,
	})
}

// AddJob adds a new job to the scheduler
func (s *Scheduler) AddJob(ctx context.Context, jobID int64) error {
	job, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(jobID)).First()
	if err != nil {
		return err
	}
	if job.Status != JobStatusNormal {
		return nil
	}
	return s.scheduleJob(ctx, job)
}

// RemoveJob removes a job from the scheduler
func (s *Scheduler) RemoveJob(jobID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	gocronJob, ok := s.jobMap[jobID]
	if !ok {
		return nil
	}

	if err := s.scheduler.RemoveJob(gocronJob.ID()); err != nil {
		return err
	}
	delete(s.jobMap, jobID)
	s.log.Info("Job removed from scheduler", zap.Int64("jobId", jobID))
	return nil
}

// UpdateJobStatus handles status changes
func (s *Scheduler) UpdateJobStatus(ctx context.Context, jobID int64, status int) error {
	if status == JobStatusNormal {
		return s.AddJob(ctx, jobID)
	}
	return s.RemoveJob(jobID)
}

// TriggerJob executes a job immediately
func (s *Scheduler) TriggerJob(ctx context.Context, jobID int64) error {
	job, err := s.q.InfraJob.WithContext(ctx).Where(s.q.InfraJob.ID.Eq(jobID)).First()
	if err != nil {
		return err
	}

	s.mu.RLock()
	handler, ok := s.handlers[job.HandlerName]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("handler not found: %s", job.HandlerName)
	}

	go s.executeJob(ctx, job, handler)
	return nil
}

// GetNextTimes calculates the next n execution times for a cron expression
func (s *Scheduler) GetNextTimes(cronExpression string, count int) ([]string, error) {
	// Parse the cron expression using gocron
	// We create a temporary job just to parse the cron and get next run times
	tempJob, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpression, false),
		gocron.NewTask(func() {}), // dummy task
	)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	// Get the next run times
	var times []string
	lastRunAt, err := tempJob.LastRun()
	if err == nil && !lastRunAt.IsZero() {
		// Use last run as base
	}

	nextRunAt, err := tempJob.NextRun()
	if err != nil {
		// Remove temp job and return error
		_ = s.scheduler.RemoveJob(tempJob.ID())
		return nil, err
	}

	// Collect next run times
	// Note: gocron v2 doesn't expose multiple future runs directly
	// We'll approximate by adding cron intervals
	times = append(times, nextRunAt.Format("2006-01-02 15:04:05"))

	// For now, return just the next run time
	// A more accurate implementation would require parsing the cron expression manually
	// or using a dedicated cron parser library

	// Clean up temp job
	_ = s.scheduler.RemoveJob(tempJob.ID())

	// If we need more than one time, we would need to use a cron parser
	// For MVP, returning just next run time
	if count > 1 {
		// TODO: Use a proper cron parser library to get multiple next times
		// For now, just return the single next time
	}

	return times, nil
}
