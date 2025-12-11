package member

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	memberModel "backend-go/internal/model/member"
	"backend-go/internal/pkg/core"
	memberSvc "backend-go/internal/service/member"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type MemberSignInConfigHandler struct {
	svc *memberSvc.MemberSignInConfigService
}

func NewMemberSignInConfigHandler(svc *memberSvc.MemberSignInConfigService) *MemberSignInConfigHandler {
	return &MemberSignInConfigHandler{svc: svc}
}

// CreateSignInConfig 创建签到规则
func (h *MemberSignInConfigHandler) CreateSignInConfig(c *gin.Context) {
	var r req.MemberSignInConfigCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateSignInConfig(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateSignInConfig 更新签到规则
func (h *MemberSignInConfigHandler) UpdateSignInConfig(c *gin.Context) {
	var r req.MemberSignInConfigUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.UpdateSignInConfig(c, &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteSignInConfig 删除签到规则
func (h *MemberSignInConfigHandler) DeleteSignInConfig(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	if err := h.svc.DeleteSignInConfig(c, id); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetSignInConfig 获得签到规则
func (h *MemberSignInConfigHandler) GetSignInConfig(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	config, err := h.svc.GetSignInConfig(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, toConfigResp(config))
}

// GetSignInConfigList 获得签到规则列表
func (h *MemberSignInConfigHandler) GetSignInConfigList(c *gin.Context) {
	var status *int
	if val, ok := c.GetQuery("status"); ok {
		s := int(core.ParseInt64(val))
		status = &s
	}

	list, err := h.svc.GetSignInConfigList(c, status)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	respList := lo.Map(list, func(item *memberModel.MemberSignInConfig, _ int) resp.MemberSignInConfigResp {
		return toConfigResp(item)
	})
	core.WriteSuccess(c, respList)
}

func toConfigResp(config *memberModel.MemberSignInConfig) resp.MemberSignInConfigResp {
	return resp.MemberSignInConfigResp{
		ID:         config.ID,
		Day:        config.Day,
		Point:      config.Point,
		Experience: config.Experience,
		Status:     config.Status,
		CreateTime: config.CreatedAt,
	}
}
