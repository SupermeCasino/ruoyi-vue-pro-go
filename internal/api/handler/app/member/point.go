package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AppMemberPointRecordHandler struct {
	svc *memberSvc.MemberPointRecordService
}

func NewAppMemberPointRecordHandler(svc *memberSvc.MemberPointRecordService) *AppMemberPointRecordHandler {
	return &AppMemberPointRecordHandler{svc: svc}
}

// GetPointRecordPage 获得用户积分记录分页
func (h *AppMemberPointRecordHandler) GetPointRecordPage(c *gin.Context) {
	var r req.AppMemberPointRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	userId := core.GetLoginUserID(c)
	pageResult, err := h.svc.GetAppPointRecordPage(c, userId, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	core.WriteSuccess(c, core.NewPageResult(lo.Map(pageResult.List, func(item *memberModel.MemberPointRecord, _ int) *resp.AppMemberPointRecordResp {
		return &resp.AppMemberPointRecordResp{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			Point:       item.Point,
			CreatedAt:   item.CreatedAt,
		}
	}), pageResult.Total))
}
