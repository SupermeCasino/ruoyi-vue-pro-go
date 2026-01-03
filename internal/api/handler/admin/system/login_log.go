package system

import (
	system2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginLogHandler struct {
	svc *system.LoginLogService
}

func NewLoginLogHandler(svc *system.LoginLogService) *LoginLogHandler {
	return &LoginLogHandler{svc: svc}
}

// GetLoginLogPage 获取登录日志分页
func (h *LoginLogHandler) GetLoginLogPage(c *gin.Context) {
	var r system2.LoginLogPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	pageResult, err := h.svc.GetLoginLogPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// Convert to Response DTO
	list := make([]system2.LoginLogResp, len(pageResult.List))
	for i, log := range pageResult.List {
		list[i] = system2.LoginLogResp{
			ID:         log.ID,
			LogType:    log.LogType,
			UserID:     log.UserID,
			UserType:   log.UserType,
			TraceID:    log.TraceID,
			Username:   log.Username,
			Result:     log.Result,
			UserIP:     log.UserIP,
			UserAgent:  log.UserAgent,
			CreateTime: log.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[system2.LoginLogResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
