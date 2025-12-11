package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type LoginLogHandler struct {
	svc *service.LoginLogService
}

func NewLoginLogHandler(svc *service.LoginLogService) *LoginLogHandler {
	return &LoginLogHandler{svc: svc}
}

// GetLoginLogPage 获取登录日志分页
func (h *LoginLogHandler) GetLoginLogPage(c *gin.Context) {
	var r req.LoginLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}

	pageResult, err := h.svc.GetLoginLogPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// Convert to Response DTO
	list := make([]resp.LoginLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = resp.LoginLogResp{
			ID:         log.ID,
			LogType:    log.LogType,
			UserID:     log.UserID,
			UserType:   log.UserType,
			TraceID:    log.TraceID,
			Username:   log.Username,
			Result:     log.Result,
			UserIP:     log.UserIP,
			UserAgent:  log.UserAgent,
			CreateTime: log.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.LoginLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
