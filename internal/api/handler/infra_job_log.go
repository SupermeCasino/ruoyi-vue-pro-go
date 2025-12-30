package handler

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

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
	id := utils.ParseInt64(c.Query("id"))
	log, err := h.svc.GetJobLog(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if log == nil {
		response.WriteBizError(c, errors.ErrNotFound)
		return
	}
	response.WriteSuccess(c, resp.JobLogResp{
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
		CreateTime:   log.CreateTime,
	})
}

// GetJobLogPage 获取定时任务日志分页
func (h *JobLogHandler) GetJobLogPage(c *gin.Context) {
	var r req.JobLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	pageResult, err := h.svc.GetJobLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
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
			CreateTime:   log.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.JobLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// ExportJobLogExcel 导出定时任务日志 Excel
func (h *JobLogHandler) ExportJobLogExcel(c *gin.Context) {
	var r req.JobLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 设置为导出所有数据
	r.PageSize = 0
	pageResult, err := h.svc.GetJobLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
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
			CreateTime:   log.CreateTime,
		}
	}

	if err := excel.WriteExcel(c, "任务日志.xls", "数据", list); err != nil {
		response.WriteError(c, 500, err.Error())
	}
}
