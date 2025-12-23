package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

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
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateGroup(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateGroup 更新用户分组
func (h *MemberGroupHandler) UpdateGroup(c *gin.Context) {
	var r req.MemberGroupUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.UpdateGroup(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteGroup 删除用户分组
func (h *MemberGroupHandler) DeleteGroup(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.DeleteGroup(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetGroup 获得用户分组详情
func (h *MemberGroupHandler) GetGroup(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	item, err := h.svc.GetGroup(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, h.convertResp(item))
}

// GetGroupPage 获得用户分组分页
func (h *MemberGroupHandler) GetGroupPage(c *gin.Context) {
	var r req.MemberGroupPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	pageResult, err := h.svc.GetGroupPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WritePage(c, pageResult.Total, lo.Map(pageResult.List, func(item *memberModel.MemberGroup, _ int) *resp.MemberGroupResp {
		return h.convertResp(item)
	}))
}

// GetSimpleGroupList 获得精简用户分组列表
func (h *MemberGroupHandler) GetSimpleGroupList(c *gin.Context) {
	list, err := h.svc.GetEnableGroupList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, lo.Map(list, func(item *memberModel.MemberGroup, _ int) *resp.MemberGroupSimpleResp {
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
		CreateTime: item.CreateTime,
	}
}
