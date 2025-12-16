package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	promotionModel "github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type AppBargainActivityHandler struct {
	activitySvc *promotion.BargainActivityService
	recordSvc   *promotion.BargainRecordService
	spuSvc      *product.ProductSpuService
}

func NewAppBargainActivityHandler(activitySvc *promotion.BargainActivityService, recordSvc *promotion.BargainRecordService, spuSvc *product.ProductSpuService) *AppBargainActivityHandler {
	return &AppBargainActivityHandler{
		activitySvc: activitySvc,
		recordSvc:   recordSvc,
		spuSvc:      spuSvc,
	}
}

// GetBargainActivityList 获得砍价活动列表 (首页推荐)
// Java: GET /list, @PermitAll
func (h *AppBargainActivityHandler) GetBargainActivityList(c *gin.Context) {
	count := int(core.ParseInt64(c.DefaultQuery("count", "6")))
	if count <= 0 {
		count = 6
	}

	list, err := h.activitySvc.GetBargainActivityListByCount(c.Request.Context(), count)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	if len(list) == 0 {
		core.WriteSuccess(c, []resp.AppBargainActivityRespVO{})
		return
	}

	// Fetch SPU Info
	spuIds := make([]int64, len(list))
	for i, item := range list {
		spuIds[i] = item.SpuID
	}
	spuMap := make(map[int64]*resp.ProductSpuResp)
	spuList, err := h.spuSvc.GetSpuList(c.Request.Context(), spuIds)
	if err == nil {
		for _, spu := range spuList {
			spuMap[spu.ID] = spu
		}
	}

	// Convert to Response (匹配 Java BargainActivityConvert.convertAppList)
	result := make([]resp.AppBargainActivityRespVO, len(list))
	for i, item := range list {
		result[i] = h.convertActivityResp(item, spuMap[item.SpuID])
	}
	core.WriteSuccess(c, result)
}

// GetBargainActivityPage 获得砍价活动分页
// Java: GET /page, @PermitAll, 使用 PageParam
func (h *AppBargainActivityHandler) GetBargainActivityPage(c *gin.Context) {
	var p core.PageParam
	if err := c.ShouldBindQuery(&p); err != nil {
		core.WriteError(c, 1001004001, "参数校验失败")
		return
	}

	page, err := h.activitySvc.GetBargainActivityPageForApp(c.Request.Context(), &p)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	if page.Total == 0 {
		core.WriteSuccess(c, core.PageResult[resp.AppBargainActivityRespVO]{List: []resp.AppBargainActivityRespVO{}, Total: 0})
		return
	}

	// Fetch SPU Info
	spuIds := make([]int64, len(page.List))
	for i, item := range page.List {
		spuIds[i] = item.SpuID
	}
	spuMap := make(map[int64]*resp.ProductSpuResp)
	spuList, err := h.spuSvc.GetSpuList(c.Request.Context(), spuIds)
	if err == nil {
		for _, spu := range spuList {
			spuMap[spu.ID] = spu
		}
	}

	// Convert to Response
	result := make([]resp.AppBargainActivityRespVO, len(page.List))
	for i, item := range page.List {
		result[i] = h.convertActivityResp(item, spuMap[item.SpuID])
	}
	core.WriteSuccess(c, core.PageResult[resp.AppBargainActivityRespVO]{List: result, Total: page.Total})
}

// GetBargainActivityDetail 获得砍价活动详情
// Java: GET /get-detail, @PermitAll
func (h *AppBargainActivityHandler) GetBargainActivityDetail(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 1001004001, "参数校验失败")
		return
	}

	act, err := h.activitySvc.GetBargainActivity(c.Request.Context(), id)
	if err != nil {
		core.WriteSuccess(c, nil)
		return
	}

	// Fetch SPU Info
	spu, _ := h.spuSvc.GetSpuDetail(c.Request.Context(), act.SpuID)

	// Fetch Success Count (Status = 1 = SUCCESS)
	successCount, _ := h.recordSvc.GetBargainRecordUserCount(c.Request.Context(), id, 1)

	// 匹配 Java BargainActivityConvert.convert(activity, successUserCount, spu)
	detail := resp.AppBargainActivityDetailRespVO{
		AppBargainActivityRespVO: h.convertActivityResp(act, spu),
		BargainFirstPrice:        act.BargainFirstPrice,
		HelpMaxCount:             act.HelpMaxCount,
		BargainCount:             act.BargainCount,
		TotalLimitCount:          act.TotalLimitCount,
		RandomMinPrice:           act.RandomMinPrice,
		RandomMaxPrice:           act.RandomMaxPrice,
		SuccessUserCount:         int(successCount),
		Remark:                   act.Remark,
	}
	core.WriteSuccess(c, detail)
}

// convertActivityResp 转换活动响应 (匹配 Java BargainActivityConvert)
func (h *AppBargainActivityHandler) convertActivityResp(item *promotionModel.PromotionBargainActivity, spu *resp.ProductSpuResp) resp.AppBargainActivityRespVO {
	r := resp.AppBargainActivityRespVO{
		ID:              item.ID,
		Name:            item.Name,
		StartTime:       item.StartTime,
		EndTime:         item.EndTime,
		SpuID:           item.SpuID,
		SkuID:           item.SkuID,
		Stock:           item.Stock,
		BargainMinPrice: item.BargainMinPrice,
	}
	if spu != nil {
		r.PicUrl = spu.PicURL
		r.MarketPrice = spu.MarketPrice
	}
	return r
}
