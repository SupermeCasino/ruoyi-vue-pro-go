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

type MemberTagHandler struct {
	svc *memberSvc.MemberTagService
}

func NewMemberTagHandler(svc *memberSvc.MemberTagService) *MemberTagHandler {
	return &MemberTagHandler{svc: svc}
}

// CreateTag 创建用户标签
func (h *MemberTagHandler) CreateTag(c *gin.Context) {
	var r req.MemberTagCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateTag(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateTag 更新用户标签
func (h *MemberTagHandler) UpdateTag(c *gin.Context) {
	var r req.MemberTagUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.UpdateTag(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteTag 删除用户标签
func (h *MemberTagHandler) DeleteTag(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.DeleteTag(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetTag 获得用户标签详情
func (h *MemberTagHandler) GetTag(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	item, err := h.svc.GetTag(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, h.convertResp(item))
}

// GetTagPage 获得用户标签分页
func (h *MemberTagHandler) GetTagPage(c *gin.Context) {
	var r req.MemberTagPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	pageResult, err := h.svc.GetTagPage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WritePage(c, pageResult.Total, lo.Map(pageResult.List, func(item *memberModel.MemberTag, _ int) *resp.MemberTagResp {
		return h.convertResp(item)
	}))
}

// GetSimpleTagList 获得精简用户标签列表
func (h *MemberTagHandler) GetSimpleTagList(c *gin.Context) {
	list, err := h.svc.GetTagList(c)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, lo.Map(list, func(item *memberModel.MemberTag, _ int) *resp.MemberTagResp {
		return h.convertResp(item)
	}))
}

func (h *MemberTagHandler) convertResp(item *memberModel.MemberTag) *resp.MemberTagResp {
	if item == nil {
		return nil
	}
	return &resp.MemberTagResp{
		ID:        item.ID,
		Name:      item.Name,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
	}
}
