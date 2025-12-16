package handler

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"

	"github.com/gin-gonic/gin"
)

type ApiErrorLogHandler struct {
	svc *service.ApiErrorLogService
}

func NewApiErrorLogHandler(svc *service.ApiErrorLogService) *ApiErrorLogHandler {
	return &ApiErrorLogHandler{svc: svc}
}

// GetApiErrorLogPage 获取API错误日志分页
func (h *ApiErrorLogHandler) GetApiErrorLogPage(c *gin.Context) {
	var r req.ApiErrorLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	pageResult, err := h.svc.GetApiErrorLogPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	list := make([]resp.ApiErrorLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = resp.ApiErrorLogResp{
			ID:                        log.ID,
			TraceID:                   log.TraceID,
			UserID:                    log.UserID,
			UserType:                  log.UserType,
			ApplicationName:           log.ApplicationName,
			RequestMethod:             log.RequestMethod,
			RequestURL:                log.RequestURL,
			RequestParams:             log.RequestParams,
			UserIP:                    log.UserIP,
			UserAgent:                 log.UserAgent,
			ExceptionTime:             log.ExceptionTime,
			ExceptionName:             log.ExceptionName,
			ExceptionMessage:          log.ExceptionMessage,
			ExceptionRootCauseMessage: log.ExceptionRootCauseMessage,
			ExceptionStackTrace:       log.ExceptionStackTrace,
			ExceptionClassName:        log.ExceptionClassName,
			ExceptionFileName:         log.ExceptionFileName,
			ExceptionMethodName:       log.ExceptionMethodName,
			ExceptionLineNumber:       log.ExceptionLineNumber,
			ProcessStatus:             log.ProcessStatus,
			ProcessTime:               log.ProcessTime,
			ProcessUserID:             log.ProcessUserID,
			CreateTime:                log.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.ApiErrorLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// UpdateApiErrorLogProcess 更新API错误日志处理状态
func (h *ApiErrorLogHandler) UpdateApiErrorLogProcess(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	processStatus := int(core.ParseInt64(c.Query("processStatus")))

	// TODO: Get login user ID from context
	processUserID := int64(1)

	if err := h.svc.UpdateApiErrorLogProcess(c, id, processStatus, processUserID); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	core.WriteSuccess(c, true)
}
