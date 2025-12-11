package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type JobLogHandler struct {
	svc *service.JobLogService
}

func NewJobLogHandler(svc *service.JobLogService) *JobLogHandler {
	return &JobLogHandler{svc: svc}
}

// GetJobLog 获取定时任务日志
func (h *JobLogHandler) GetJobLog(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	log, err := h.svc.GetJobLog(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if log == nil {
		core.WriteError(c, 404, "日志不存在")
		return
	}
	core.WriteSuccess(c, resp.JobLogResp{
		ID:           log.ID,
		JobID:        log.JobID,
		HandlerName:  log.HandlerName,
		HandlerParam: log.HandlerParam,
		ExecuteIndex: log.ExecuteIndex,
		BeginTime:    log.BeginTime,
		EndTime:      log.EndTime,
		Duration:     log.Duration,
		Status:       log.Status,
		Result:       log.Result,
		CreateTime:   log.CreatedAt,
	})
}

// GetJobLogPage 获取定时任务日志分页
func (h *JobLogHandler) GetJobLogPage(c *gin.Context) {
	var r req.JobLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	pageResult, err := h.svc.GetJobLogPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	list := make([]resp.JobLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = resp.JobLogResp{
			ID:           log.ID,
			JobID:        log.JobID,
			HandlerName:  log.HandlerName,
			HandlerParam: log.HandlerParam,
			ExecuteIndex: log.ExecuteIndex,
			BeginTime:    log.BeginTime,
			EndTime:      log.EndTime,
			Duration:     log.Duration,
			Status:       log.Status,
			Result:       log.Result,
			CreateTime:   log.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.JobLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
