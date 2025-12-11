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

type MemberGroupHandler struct {
	svc *memberSvc.MemberGroupService
}

func NewMemberGroupHandler(svc *memberSvc.MemberGroupService) *MemberGroupHandler {
	return &MemberGroupHandler{svc: svc}
}

// CreateGroup 创建用户分组
func (h *MemberGroupHandler) CreateGroup(c *gin.Context) {
	var r req.MemberGroupCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateGroup(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateGroup 更新用户分组
func (h *MemberGroupHandler) UpdateGroup(c *gin.Context) {
	var r req.MemberGroupUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.UpdateGroup(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteGroup 删除用户分组
func (h *MemberGroupHandler) DeleteGroup(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.DeleteGroup(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetGroup 获得用户分组详情
func (h *MemberGroupHandler) GetGroup(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	item, err := h.svc.GetGroup(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, h.convertResp(item))
}

// GetGroupPage 获得用户分组分页
func (h *MemberGroupHandler) GetGroupPage(c *gin.Context) {
	var r req.MemberGroupPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	pageResult, err := h.svc.GetGroupPage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WritePage(c, pageResult.Total, lo.Map(pageResult.List, func(item *memberModel.MemberGroup, _ int) *resp.MemberGroupResp {
		return h.convertResp(item)
	}))
}

// GetSimpleGroupList 获得精简用户分组列表
func (h *MemberGroupHandler) GetSimpleGroupList(c *gin.Context) {
	list, err := h.svc.GetEnableGroupList(c)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, lo.Map(list, func(item *memberModel.MemberGroup, _ int) *resp.MemberGroupSimpleResp {
		return &resp.MemberGroupSimpleResp{
			ID:   item.ID,
			Name: item.Name,
		}
	}))
}

func (h *MemberGroupHandler) convertResp(item *memberModel.MemberGroup) *resp.MemberGroupResp {
	if item == nil {
		return nil
	}
	return &resp.MemberGroupResp{
		ID:        item.ID,
		Name:      item.Name,
		Remark:    item.Remark,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
	}
}
