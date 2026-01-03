package system

import (
	system2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type OperateLogHandler struct {
	svc *system.OperateLogService
}

func NewOperateLogHandler(svc *system.OperateLogService) *OperateLogHandler {
	return &OperateLogHandler{svc: svc}
}

// GetOperateLogPage 获取操作日志分页
func (h *OperateLogHandler) GetOperateLogPage(c *gin.Context) {
	var r system2.OperateLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	pageResult, err := h.svc.GetOperateLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// Convert to Response DTO
	// Note: userName is derived from userId. For now we leave it empty or can join with user table later.
	list := make([]system2.OperateLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = system2.OperateLogResp{
			ID:            log.ID,
			TraceID:       log.TraceID,
			UserID:        log.UserID,
			UserName:      "", // TODO: Join with user table to get name
			Type:          log.Type,
			SubType:       log.SubType,
			BizID:         log.BizID,
			Action:        log.Action,
			Extra:         log.Extra,
			RequestMethod: log.RequestMethod,
			RequestURL:    log.RequestURL,
			UserIP:        log.UserIP,
			UserAgent:     log.UserAgent,
			CreateTime:    log.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[system2.OperateLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
