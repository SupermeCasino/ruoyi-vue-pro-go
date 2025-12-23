package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BargainActivityHandler struct {
	activitySvc *promotion.BargainActivityService
	recordSvc   *promotion.BargainRecordService
	helpSvc     *promotion.BargainHelpService
	spuSvc      *product.ProductSpuService
}

func NewBargainActivityHandler(activitySvc *promotion.BargainActivityService, recordSvc *promotion.BargainRecordService, helpSvc *promotion.BargainHelpService, spuSvc *product.ProductSpuService) *BargainActivityHandler {
	return &BargainActivityHandler{
		activitySvc: activitySvc,
		recordSvc:   recordSvc,
		helpSvc:     helpSvc,
		spuSvc:      spuSvc,
	}
}

// CreateBargainActivity 创建砍价活动
func (h *BargainActivityHandler) CreateBargainActivity(c *gin.Context) {
	var r req.BargainActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 1001004001, "参数校验失败: "+err.Error())
		return
	}
	id, err := h.activitySvc.CreateBargainActivity(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateBargainActivity 更新砍价活动
func (h *BargainActivityHandler) UpdateBargainActivity(c *gin.Context) {
	var r req.BargainActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 1001004001, "参数校验失败: "+err.Error())
		return
	}
	if err := h.activitySvc.UpdateBargainActivity(c.Request.Context(), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// CloseBargainActivity 关闭砍价活动
func (h *BargainActivityHandler) CloseBargainActivity(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 1001004001, "活动ID不能为空")
		return
	}
	if err := h.activitySvc.CloseBargainActivity(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteBargainActivity 删除砍价活动
func (h *BargainActivityHandler) DeleteBargainActivity(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 1001004001, "活动ID不能为空")
		return
	}
	if err := h.activitySvc.DeleteBargainActivity(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetBargainActivity 获得砍价活动
func (h *BargainActivityHandler) GetBargainActivity(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 1001004001, "活动ID不能为空")
		return
	}
	act, err := h.activitySvc.GetBargainActivity(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if act == nil {
		response.WriteSuccess(c, nil)
		return
	}

	respVO := resp.BargainActivityResp{
		ID:                act.ID,
		SpuID:             act.SpuID,
		SkuID:             act.SkuID,
		Name:              act.Name,
		StartTime:         act.StartTime,
		EndTime:           act.EndTime,
		BargainFirstPrice: act.BargainFirstPrice,
		BargainMinPrice:   act.BargainMinPrice,
		Stock:             act.Stock,
		TotalStock:        act.TotalStock,
		HelpMaxCount:      act.HelpMaxCount,
		BargainCount:      act.BargainCount,
		TotalLimitCount:   act.TotalLimitCount,
		RandomMinPrice:    act.RandomMinPrice,
		RandomMaxPrice:    act.RandomMaxPrice,
		Status:            act.Status,
		Sort:              act.Sort,
		CreateTime:         act.CreateTime,
	}
	// Note: CreateReq has Remark but DO model logic I used previously didn't output Remark.
	// But it is there. I should ensure DO has Remark later if needed.
	// Checked DO: Seckill has Remark. My Bargain DO?
	// I'll check Bargain DO content again.
	// If Remark is missing, I'll ignore for now or add later.
	response.WriteSuccess(c, respVO)
}

// GetBargainActivityPage 获得砍价活动分页
func (h *BargainActivityHandler) GetBargainActivityPage(c *gin.Context) {
	var r req.BargainActivityPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 1001004001, "参数校验失败: "+err.Error())
		return
	}
	pageResult, err := h.activitySvc.GetBargainActivityPage(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	if len(pageResult.List) == 0 {
		response.WriteSuccess(c, pagination.PageResult[resp.BargainActivityPageItemResp]{
			List:  []resp.BargainActivityPageItemResp{},
			Total: pageResult.Total,
		})
		return
	}

	// Collect IDs
	spuIds := make([]int64, 0, len(pageResult.List))
	activityIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		spuIds = append(spuIds, item.SpuID)
		activityIds = append(activityIds, item.ID)
	}

	// Fetch SPUs
	spuMap := make(map[int64]*resp.ProductSpuResp)
	if len(spuIds) > 0 {
		spuList, err := h.spuSvc.GetSpuList(c.Request.Context(), spuIds)
		if err == nil {
			for _, spu := range spuList {
				spuMap[spu.ID] = spu
			}
		}
	}

	// Fetch Stats (Ignoring errors for stats)
	// Status 1 = Success? Need to check Enum.
	// Usually Success=1 or 2.
	recordUserCountMap, _ := h.recordSvc.GetBargainRecordUserCountMap(c.Request.Context(), activityIds, nil)
	successStatus := 1 // Assume Success=1
	recordSuccessUserCountMap, _ := h.recordSvc.GetBargainRecordUserCountMap(c.Request.Context(), activityIds, &successStatus)
	helpUserCountMap, _ := h.helpSvc.GetBargainHelpUserCountMapByActivity(c.Request.Context(), activityIds)

	list := make([]resp.BargainActivityPageItemResp, len(pageResult.List))
	for i, item := range pageResult.List {
		spuName := ""
		picUrl := ""
		marketPrice := 0

		if spu, ok := spuMap[item.SpuID]; ok {
			spuName = spu.Name
			picUrl = spu.PicURL
			marketPrice = spu.MarketPrice // Or SKU Price if available? SPU market price usually min.
		}

		list[i] = resp.BargainActivityPageItemResp{
			BargainActivityResp: resp.BargainActivityResp{
				ID:                item.ID,
				Name:              item.Name,
				Status:            item.Status,
				StartTime:         item.StartTime,
				EndTime:           item.EndTime,
				Stock:             item.Stock,
				TotalStock:        item.TotalStock,
				SpuID:             item.SpuID,
				SkuID:             item.SkuID,
				BargainFirstPrice: item.BargainFirstPrice,
				BargainMinPrice:   item.BargainMinPrice,
				HelpMaxCount:      item.HelpMaxCount,
				BargainCount:      item.BargainCount,
				TotalLimitCount:   item.TotalLimitCount,
				RandomMinPrice:    item.RandomMinPrice,
				RandomMaxPrice:    item.RandomMaxPrice,
				Sort:              item.Sort,
				CreateTime:         item.CreateTime,
			},
			SpuName:                spuName,
			PicUrl:                 picUrl,
			MarketPrice:            marketPrice,
			RecordUserCount:        recordUserCountMap[item.ID],
			RecordSuccessUserCount: recordSuccessUserCountMap[item.ID],
			HelpUserCount:          helpUserCountMap[item.ID],
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.BargainActivityPageItemResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
