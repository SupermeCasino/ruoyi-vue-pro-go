package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	pageResult, err := h.svc.GetSignInRecordPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
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

	response.WriteSuccess(c, pagination.NewPageResult(respList, pageResult.Total))
}
