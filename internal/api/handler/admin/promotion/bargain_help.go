package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type BargainHelpHandler struct {
	svc     *promotion.BargainHelpService
	userSvc *member.MemberUserService
}

func NewBargainHelpHandler(svc *promotion.BargainHelpService, userSvc *member.MemberUserService) *BargainHelpHandler {
	return &BargainHelpHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

// GetBargainHelpPage 获得砍价助力分页
func (h *BargainHelpHandler) GetBargainHelpPage(c *gin.Context) {
	var r req.BargainHelpPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 1. Get Page
	pageResult, err := h.svc.GetBargainHelpPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if len(pageResult.List) == 0 {
		response.WriteSuccess(c, pagination.PageResult[resp.BargainHelpResp]{
			List:  []resp.BargainHelpResp{},
			Total: pageResult.Total,
		})
		return
	}

	// 2. Collect IDs
	userIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		userIds = append(userIds, item.UserID)
	}

	// 3. Fetch Data
	userMap, _ := h.userSvc.GetUserMap(c, userIds)

	// 4. Assemble
	list := make([]resp.BargainHelpResp, len(pageResult.List))
	for i, item := range pageResult.List {
		vo := resp.BargainHelpResp{
			ID:          item.ID,
			UserID:      item.UserID,
			ActivityID:  item.ActivityID,
			RecordID:    item.RecordID,
			ReducePrice: item.ReducePrice,
			CreateTime:  item.CreateTime,
		}
		if u, ok := userMap[item.UserID]; ok {
			vo.UserNickname = u.Nickname
			vo.UserAvatar = u.Avatar
		}
		list[i] = vo
	}

	response.WriteSuccess(c, pagination.PageResult[resp.BargainHelpResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
