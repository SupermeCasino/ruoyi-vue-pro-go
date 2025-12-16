package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type MemberLevelHandler struct {
	svc *memberSvc.MemberLevelService
}

func NewMemberLevelHandler(svc *memberSvc.MemberLevelService) *MemberLevelHandler {
	return &MemberLevelHandler{svc: svc}
}

// CreateLevel 创建等级
func (h *MemberLevelHandler) CreateLevel(c *gin.Context) {
	var r req.MemberLevelCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	id, err := h.svc.CreateLevel(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateLevel 更新等级
func (h *MemberLevelHandler) UpdateLevel(c *gin.Context) {
	var r req.MemberLevelUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.UpdateLevel(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteLevel 删除等级
func (h *MemberLevelHandler) DeleteLevel(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.DeleteLevel(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetLevel 获得等级详情
func (h *MemberLevelHandler) GetLevel(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	item, err := h.svc.GetLevel(c, id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, h.convertResp(item))
}

// GetLevelPage 获得等级分页
func (h *MemberLevelHandler) GetLevelPage(c *gin.Context) {
	var r req.MemberLevelPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	pageResult, err := h.svc.GetLevelPage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WritePage(c, pageResult.Total, lo.Map(pageResult.List, func(item *member.MemberLevel, _ int) *resp.MemberLevelResp {
		return h.convertResp(item)
	}))
}

// GetLevelListSimple 获得开启的等级列表 (用于下拉)
func (h *MemberLevelHandler) GetLevelListSimple(c *gin.Context) {
	list, err := h.svc.GetLevelSimpleList(c)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, lo.Map(list, func(item *member.MemberLevel, _ int) *resp.MemberLevelResp {
		return h.convertResp(item)
	}))
}

func (h *MemberLevelHandler) convertResp(item *member.MemberLevel) *resp.MemberLevelResp {
	return &resp.MemberLevelResp{
		ID:              item.ID,
		Name:            item.Name,
		Level:           item.Level,
		Experience:      item.Experience,
		DiscountPercent: item.DiscountPercent,
		Icon:            item.Icon,
		BackgroundURL:   item.BackgroundURL,
		Status:          item.Status,
		Remark:          item.Remark,
		CreatedAt:       item.CreatedAt,
	}
}
