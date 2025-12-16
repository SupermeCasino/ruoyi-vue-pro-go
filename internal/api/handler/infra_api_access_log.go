package handler

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ApiAccessLogHandler struct {
	svc *service.ApiAccessLogService
}

func NewApiAccessLogHandler(svc *service.ApiAccessLogService) *ApiAccessLogHandler {
	return &ApiAccessLogHandler{svc: svc}
}

// GetApiAccessLogPage 获取API访问日志分页
func (h *ApiAccessLogHandler) GetApiAccessLogPage(c *gin.Context) {
	var r req.ApiAccessLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	pageResult, err := h.svc.GetApiAccessLogPage(c, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	list := make([]resp.ApiAccessLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = resp.ApiAccessLogResp{
			ID:              log.ID,
			TraceID:         log.TraceID,
			UserID:          log.UserID,
			UserType:        log.UserType,
			ApplicationName: log.ApplicationName,
			RequestMethod:   log.RequestMethod,
			RequestURL:      log.RequestURL,
			RequestParams:   log.RequestParams,
			ResponseBody:    log.ResponseBody,
			UserIP:          log.UserIP,
			UserAgent:       log.UserAgent,
			OperateModule:   log.OperateModule,
			OperateName:     log.OperateName,
			OperateType:     log.OperateType,
			BeginTime:       log.BeginTime,
			EndTime:         log.EndTime,
			Duration:        log.Duration,
			ResultCode:      log.ResultCode,
			ResultMsg:       log.ResultMsg,
			CreateTime:      log.CreatedAt,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.ApiAccessLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
