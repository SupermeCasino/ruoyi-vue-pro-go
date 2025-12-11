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

type MemberSignInRecordHandler struct {
	svc     *memberSvc.MemberSignInRecordService
	userSvc *memberSvc.MemberUserService
}

func NewMemberSignInRecordHandler(svc *memberSvc.MemberSignInRecordService, userSvc *memberSvc.MemberUserService) *MemberSignInRecordHandler {
	return &MemberSignInRecordHandler{svc: svc, userSvc: userSvc}
}

// GetSignInRecordPage 获得签到记录分页
func (h *MemberSignInRecordHandler) GetSignInRecordPage(c *gin.Context) {
	var r req.MemberSignInRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}

	pageResult, err := h.svc.GetSignInRecordPage(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	// Fetch users for nicknames
	userIds := lo.Map(pageResult.List, func(item *memberModel.MemberSignInRecord, _ int) int64 { return item.UserID })
	userMap, _ := h.userSvc.GetUserMap(c, userIds)

	respList := lo.Map(pageResult.List, func(item *memberModel.MemberSignInRecord, _ int) resp.MemberSignInRecordResp {
		nickname := ""
		if u, ok := userMap[item.UserID]; ok {
			nickname = u.Nickname
		}
		return resp.MemberSignInRecordResp{
			ID:         item.ID,
			UserID:     item.UserID,
			Nickname:   nickname,
			Day:        item.Day,
			Point:      item.Point,
			Experience: item.Experience,
			CreatedAt:  item.CreatedAt,
		}
	})

	core.WriteSuccess(c, core.NewPageResult(respList, pageResult.Total))
}
