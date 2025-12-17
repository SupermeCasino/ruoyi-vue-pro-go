package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type CombinationRecordHandler struct {
	svc         promotion.CombinationRecordService
	activitySvc promotion.CombinationActivityService
}

func NewCombinationRecordHandler(
	svc promotion.CombinationRecordService,
	activitySvc promotion.CombinationActivityService,
) *CombinationRecordHandler {
	return &CombinationRecordHandler{
		svc:         svc,
		activitySvc: activitySvc,
	}
}

// GetCombinationRecordPage 获得拼团记录分页
func (h *CombinationRecordHandler) GetCombinationRecordPage(c *gin.Context) {
	var r req.CombinationRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	// 1. Get Page
	pageResult, err := h.svc.GetCombinationRecordPageAdmin(c, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	if len(pageResult.List) == 0 {
		response.WriteSuccess(c, pagination.PageResult[resp.CombinationRecordPageItemRespVO]{
			List:  []resp.CombinationRecordPageItemRespVO{},
			Total: pageResult.Total,
		})
		return
	}

	// 2. Collection IDs for Activity
	activityIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		activityIds = append(activityIds, item.ActivityID)
	}

	// 3. Fetch Data (Only Activity needed)
	activityMap, _ := h.activitySvc.GetCombinationActivityMap(c, activityIds)

	// 4. Assemble
	list := make([]resp.CombinationRecordPageItemRespVO, len(pageResult.List))
	for i, item := range pageResult.List {
		vo := resp.CombinationRecordPageItemRespVO{
			ID:               item.ID,
			ActivityID:       item.ActivityID,
			UserID:           item.UserID,
			UserCount:        item.UserCount,
			Status:           item.Status,
			CombinationPrice: item.CombinationPrice,
			HeadID:           item.HeadID,
			OrderID:          item.OrderID,
			VirtualGroup:     bool(item.VirtualGroup), // VirtualGroup is BitBool usually? Model says bool.
			// Let's check model definition. It says `bool`.
			// If it's BitBool in DB but bool in Struct, cast might be needed if type differs.
			// Model line 81: VirtualGroup bool.
			ExpireTime: item.ExpireTime,
			StartTime:  item.StartTime,
			EndTime:    item.EndTime,
			CreateTime: item.CreatedAt,
			Nickname:   item.Nickname,
			Avatar:     item.Avatar,
			SpuID:      item.SpuID,
			SpuName:    item.SpuName,
			PicUrl:     item.PicUrl,
		}

		// Activity Name
		if act, ok := activityMap[item.ActivityID]; ok {
			vo.ActivityName = act.Name
			vo.UserSize = act.UserSize // Activity defines size
		} else {
			vo.UserSize = item.UserSize // Fallback to record's size
		}

		list[i] = vo
	}

	response.WriteSuccess(c, pagination.PageResult[resp.CombinationRecordPageItemRespVO]{
		List:  list,
		Total: pageResult.Total,
	})
}

// GetCombinationRecordSummary 获得拼团记录的概要信息
func (h *CombinationRecordHandler) GetCombinationRecordSummary(c *gin.Context) {
	summary, err := h.svc.GetCombinationRecordSummaryAdmin(c.Request.Context())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, summary)
}
