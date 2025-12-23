package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	promotionModel "github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

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
	count := int(utils.ParseInt64(c.DefaultQuery("count", "6")))
	if count <= 0 {
		count = 6
	}

	list, err := h.activitySvc.GetBargainActivityListByCount(c.Request.Context(), count)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	if len(list) == 0 {
		response.WriteSuccess(c, []resp.AppBargainActivityRespVO{})
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

	// Convert to Response (匹配 Java BargainActivityConvert.convertAppList)并过滤
	result := make([]resp.AppBargainActivityRespVO, 0, len(list))
	for _, item := range list {
		spu, ok := spuMap[item.SpuID]
		if !ok || spu.Status != model.ProductSpuStatusEnable {
			continue
		}
		result = append(result, h.convertActivityResp(item, spu))
	}
	response.WriteSuccess(c, result)
}

// GetBargainActivityPage 获得砍价活动分页
// Java: GET /page, @PermitAll, 使用 PageParam
func (h *AppBargainActivityHandler) GetBargainActivityPage(c *gin.Context) {
	var p pagination.PageParam
	if err := c.ShouldBindQuery(&p); err != nil {
		response.WriteError(c, 1001004001, "参数校验失败")
		return
	}

	page, err := h.activitySvc.GetBargainActivityPageForApp(c.Request.Context(), &p)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	if page.Total == 0 {
		response.WriteSuccess(c, pagination.PageResult[resp.AppBargainActivityRespVO]{List: []resp.AppBargainActivityRespVO{}, Total: 0})
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

	// Convert to Response 并过滤
	result := make([]resp.AppBargainActivityRespVO, 0, len(page.List))
	for _, item := range page.List {
		spu, ok := spuMap[item.SpuID]
		if !ok || spu.Status != model.ProductSpuStatusEnable {
			continue
		}
		result = append(result, h.convertActivityResp(item, spu))
	}
	response.WriteSuccess(c, pagination.PageResult[resp.AppBargainActivityRespVO]{List: result, Total: page.Total})
}

// GetBargainActivityDetail 获得砍价活动详情
// Java: GET /get-detail, @PermitAll
func (h *AppBargainActivityHandler) GetBargainActivityDetail(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteError(c, 1001004001, "参数校验失败")
		return
	}

	act, err := h.activitySvc.GetBargainActivity(c.Request.Context(), id)
	if err != nil {
		response.WriteSuccess(c, nil)
		return
	}

	// Fetch SPU Info
	spu, _ := h.spuSvc.GetSpuDetail(c.Request.Context(), act.SpuID)
	if spu == nil || spu.Status != model.ProductSpuStatusEnable {
		response.WriteBizError(c, errors.NewBizError(1001004003, "砍价活动已结束或商品已下架"))
		return
	}

	// Fetch Success Count (Status = 1 = SUCCESS)
	successCount, _ := h.recordSvc.GetBargainRecordUserCount(c.Request.Context(), id, model.BargainRecordStatusSuccess)

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
	response.WriteSuccess(c, detail)
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
