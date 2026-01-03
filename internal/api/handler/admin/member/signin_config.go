package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

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
	var r member.MemberSignInConfigCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateSignInConfig(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateSignInConfig 更新签到规则
func (h *MemberSignInConfigHandler) UpdateSignInConfig(c *gin.Context) {
	var r member.MemberSignInConfigUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateSignInConfig(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteSignInConfig 删除签到规则
func (h *MemberSignInConfigHandler) DeleteSignInConfig(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteSignInConfig(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetSignInConfig 获得签到规则
func (h *MemberSignInConfigHandler) GetSignInConfig(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	config, err := h.svc.GetSignInConfig(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, toConfigResp(config))
}

// GetSignInConfigList 获得签到规则列表
func (h *MemberSignInConfigHandler) GetSignInConfigList(c *gin.Context) {
	var status *int
	if val, ok := c.GetQuery("status"); ok {
		s := int(utils.ParseInt64(val))
		status = &s
	}

	list, err := h.svc.GetSignInConfigList(c, status)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	respList := lo.Map(list, func(item *memberModel.MemberSignInConfig, _ int) member.MemberSignInConfigResp {
		return toConfigResp(item)
	})
	response.WriteSuccess(c, respList)
}

func toConfigResp(config *memberModel.MemberSignInConfig) member.MemberSignInConfigResp {
	return member.MemberSignInConfigResp{
		ID:         config.ID,
		Day:        config.Day,
		Point:      config.Point,
		Experience: config.Experience,
		Status:     config.Status,
		CreateTime: config.CreateTime,
	}
}
