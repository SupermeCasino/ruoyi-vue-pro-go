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

type MemberPointRecordHandler struct {
	svc           *memberSvc.MemberPointRecordService
	memberUserSvc *memberSvc.MemberUserService
}

func NewMemberPointRecordHandler(svc *memberSvc.MemberPointRecordService, memberUserSvc *memberSvc.MemberUserService) *MemberPointRecordHandler {
	return &MemberPointRecordHandler{svc: svc, memberUserSvc: memberUserSvc}
}

// GetPointRecordPage 获得用户积分记录分页
func (h *MemberPointRecordHandler) GetPointRecordPage(c *gin.Context) {
	var r req.MemberPointRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	pageResult, err := h.svc.GetPointRecordPage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	// Fetch User info for response
	userIds := lo.Map(pageResult.List, func(item *memberModel.MemberPointRecord, _ int) int64 {
		return item.UserID
	})
	userMap, err := h.memberUserSvc.GetUserMap(c, userIds)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	core.WriteSuccess(c, core.NewPageResult(lo.Map(pageResult.List, func(item *memberModel.MemberPointRecord, _ int) *resp.MemberPointRecordResp {
		nickname := ""
		if user, ok := userMap[item.UserID]; ok {
			nickname = user.Nickname
		}
		return &resp.MemberPointRecordResp{
			ID:          item.ID,
			UserID:      item.UserID,
			Nickname:    nickname,
			BizID:       item.BizID,
			BizType:     item.BizType,
			Title:       item.Title,
			Description: item.Description,
			Point:       item.Point,
			TotalPoint:  item.TotalPoint,
			CreatedAt:   item.CreatedAt,
		}
	}), pageResult.Total))
}
