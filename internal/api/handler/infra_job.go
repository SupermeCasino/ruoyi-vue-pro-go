package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	svc *service.JobService
}

func NewJobHandler(svc *service.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

// CreateJob 创建定时任务
func (h *JobHandler) CreateJob(c *gin.Context) {
	var r req.JobSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateJob(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateJob 更新定时任务
func (h *JobHandler) UpdateJob(c *gin.Context) {
	var r req.JobSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateJob(c, &r); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// UpdateJobStatus 更新定时任务状态
func (h *JobHandler) UpdateJobStatus(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	status := int(core.ParseInt64(c.Query("status")))
	if err := h.svc.UpdateJobStatus(c, id, status); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteJob 删除定时任务
func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if err := h.svc.DeleteJob(c, id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetJob 获取定时任务
func (h *JobHandler) GetJob(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	job, err := h.svc.GetJob(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if job == nil {
		core.WriteError(c, 404, "任务不存在")
		return
	}
	core.WriteSuccess(c, resp.JobResp{
		ID:             job.ID,
		Name:           job.Name,
		Status:         job.Status,
		HandlerName:    job.HandlerName,
		HandlerParam:   job.HandlerParam,
		CronExpression: job.CronExpression,
		RetryCount:     job.RetryCount,
		RetryInterval:  job.RetryInterval,
		MonitorTimeout: job.MonitorTimeout,
		CreateTime:     job.CreatedAt,
	})
}

// GetJobPage 获取定时任务分页
func (h *JobHandler) GetJobPage(c *gin.Context) {
	var r req.JobPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	pageResult, err := h.svc.GetJobPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	list := make([]resp.JobResp, len(pageResult.List))
	for i, job := range pageResult.List {
		list[i] = resp.JobResp{
			ID:             job.ID,
			Name:           job.Name,
			Status:         job.Status,
			HandlerName:    job.HandlerName,
			HandlerParam:   job.HandlerParam,
			CronExpression: job.CronExpression,
			RetryCount:     job.RetryCount,
			RetryInterval:  job.RetryInterval,
			MonitorTimeout: job.MonitorTimeout,
			CreateTime:     job.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.JobResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// TriggerJob 触发定时任务
func (h *JobHandler) TriggerJob(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if err := h.svc.TriggerJob(c, id); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}
