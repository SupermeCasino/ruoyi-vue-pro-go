package handler

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type OperateLogHandler struct {
	svc *service.OperateLogService
}

func NewOperateLogHandler(svc *service.OperateLogService) *OperateLogHandler {
	return &OperateLogHandler{svc: svc}
}

// GetOperateLogPage 获取操作日志分页
func (h *OperateLogHandler) GetOperateLogPage(c *gin.Context) {
	var r req.OperateLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	pageResult, err := h.svc.GetOperateLogPage(c, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	// Convert to Response DTO
	// Note: userName is derived from userId. For now we leave it empty or can join with user table later.
	list := make([]resp.OperateLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = resp.OperateLogResp{
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

	response.WriteSuccess(c, pagination.PageResult[resp.OperateLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
