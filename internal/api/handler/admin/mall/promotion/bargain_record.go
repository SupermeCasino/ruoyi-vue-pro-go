package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type BargainRecordHandler struct {
	svc         *promotion.BargainRecordService
	activitySvc *promotion.BargainActivityService
	userSvc     *member.MemberUserService
}

func NewBargainRecordHandler(
	svc *promotion.BargainRecordService,
	activitySvc *promotion.BargainActivityService,
	userSvc *member.MemberUserService,
) *BargainRecordHandler {
	return &BargainRecordHandler{
		svc:         svc,
		activitySvc: activitySvc,
		userSvc:     userSvc,
	}
}

// GetBargainRecordPage 获得砍价记录分页
func (h *BargainRecordHandler) GetBargainRecordPage(c *gin.Context) {
	var r req.BargainRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 1. Get Page of DOs
	pageResult, err := h.svc.GetBargainRecordPageAdmin(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if len(pageResult.List) == 0 {
		response.WriteSuccess(c, pagination.PageResult[resp.BargainRecordResp]{
			List:  []resp.BargainRecordResp{},
			Total: pageResult.Total,
		})
		return
	}

	// 2. Collect IDs
	userIds := make([]int64, 0, len(pageResult.List))
	activityIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		userIds = append(userIds, item.UserID)
		activityIds = append(activityIds, item.ActivityID)
	}

	// 3. Fetch Enriched Data
	userMap, _ := h.userSvc.GetUserMap(c, userIds)
	activityMap, _ := h.activitySvc.GetBargainActivityMap(c, activityIds)

	// 4. Assemble VOs
	list := make([]resp.BargainRecordResp, len(pageResult.List))
	for i, item := range pageResult.List {
		nickname := ""
		avatar := ""
		if u, ok := userMap[item.UserID]; ok {
			nickname = u.Nickname
			avatar = u.Avatar
		}
		activityName := ""
		if act, ok := activityMap[item.ActivityID]; ok {
			activityName = act.Name
		}

		list[i] = resp.BargainRecordResp{
			ID:                item.ID,
			UserID:            item.UserID,
			UserNickname:      nickname,
			UserAvatar:        avatar,
			ActivityID:        item.ActivityID,
			ActivityName:      activityName,
			SpuID:             item.SpuID,
			SkuID:             item.SkuID,
			BargainFirstPrice: item.BargainFirstPrice,
			BargainPrice:      item.BargainPrice,
			Status:            item.Status,
			EndTime:           item.EndTime,
			OrderID:           item.OrderID,
			CreateTime:        item.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.BargainRecordResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
