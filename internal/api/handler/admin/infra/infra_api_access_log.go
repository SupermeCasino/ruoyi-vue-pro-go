package infra

import (
	infra2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ApiAccessLogHandler struct {
	svc *infra.ApiAccessLogService
}

func NewApiAccessLogHandler(svc *infra.ApiAccessLogService) *ApiAccessLogHandler {
	return &ApiAccessLogHandler{svc: svc}
}

// GetApiAccessLogPage 获取API访问日志分页
func (h *ApiAccessLogHandler) GetApiAccessLogPage(c *gin.Context) {
	var r infra2.ApiAccessLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	pageResult, err := h.svc.GetApiAccessLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]infra2.ApiAccessLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = infra2.ApiAccessLogResp{
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
			CreateTime:      log.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[infra2.ApiAccessLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
