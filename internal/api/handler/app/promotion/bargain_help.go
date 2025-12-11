package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	memberModel "backend-go/internal/model/member"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/member"
	"backend-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type AppBargainHelpHandler struct {
	helpSvc *promotion.BargainHelpService
	userSvc *member.MemberUserService
}

func NewAppBargainHelpHandler(helpSvc *promotion.BargainHelpService, userSvc *member.MemberUserService) *AppBargainHelpHandler {
	return &AppBargainHelpHandler{
		helpSvc: helpSvc,
		userSvc: userSvc,
	}
}

// CreateBargainHelp 砍价助力
// Java: POST /create, returns ReducePrice
func (h *AppBargainHelpHandler) CreateBargainHelp(c *gin.Context) {
	var r req.AppBargainHelpCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 1001004001, "参数校验失败")
		return
	}
	help, err := h.helpSvc.CreateBargainHelp(c.Request.Context(), c.GetInt64(core.CtxUserIDKey), &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, help.ReducePrice)
}

// GetBargainHelpList 获得砍价助力列表
// Java: GET /list
func (h *AppBargainHelpHandler) GetBargainHelpList(c *gin.Context) {
	recordId := core.ParseInt64(c.Query("recordId"))
	if recordId == 0 {
		core.WriteSuccess(c, []resp.AppBargainHelpRespVO{})
		return
	}

	list, err := h.helpSvc.GetBargainHelpList(c.Request.Context(), recordId)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	// Fetch User Info
	userIds := make([]int64, len(list))
	for i, item := range list {
		userIds[i] = item.UserID
	}
	userMap := make(map[int64]*memberModel.MemberUser)
	if len(userIds) > 0 {
		um, err := h.userSvc.GetUserMap(c.Request.Context(), userIds)
		if err == nil {
			for k, v := range um {
				userMap[k] = v
			}
		}
	}

	resList := make([]resp.AppBargainHelpRespVO, len(list))
	for i, item := range list {
		vo := resp.AppBargainHelpRespVO{
			ReducePrice: item.ReducePrice,
			CreateTime:  item.CreatedAt,
		}
		if u, ok := userMap[item.UserID]; ok {
			vo.Nickname = u.Nickname
			vo.Avatar = u.Avatar
		}
		resList[i] = vo
	}
	core.WriteSuccess(c, resList)
}
