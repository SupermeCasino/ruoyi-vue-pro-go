package promotion

import (
	"strconv"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	productSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	promotionSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AppPointActivityHandler struct {
	svc    *promotionSvc.PointActivityService
	spuSvc *productSvc.ProductSpuService
}

func NewAppPointActivityHandler(svc *promotionSvc.PointActivityService, spuSvc *productSvc.ProductSpuService) *AppPointActivityHandler {
	return &AppPointActivityHandler{
		svc:    svc,
		spuSvc: spuSvc,
	}
}

// GetPointActivityListByIds 获得积分商城活动列表
// 对齐 Java: AppPointActivityController.getPointActivityListByIds
func (h *AppPointActivityHandler) GetPointActivityListByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		response.WriteSuccess(c, []*resp.AppPointActivityRespVO{})
		return
	}

	parts := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		if id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		response.WriteSuccess(c, []*resp.AppPointActivityRespVO{})
		return
	}

	list, err := h.svc.GetPointActivityListByIds(c.Request.Context(), ids)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	result, err := h.buildAppPointActivityRespVOList(c, list)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, result)
}

func (h *AppPointActivityHandler) buildAppPointActivityRespVOList(c *gin.Context, activityList []*promotion.PromotionPointActivity) ([]*resp.AppPointActivityRespVO, error) {
	if len(activityList) == 0 {
		return []*resp.AppPointActivityRespVO{}, nil
	}

	// 1. 获取活动商品列表 (用于获取最低积分/价格)
	activityIds := lo.Map(activityList, func(item *promotion.PromotionPointActivity, _ int) int64 {
		return item.ID
	})
	products, err := h.svc.GetPointProductListByActivityIds(c, activityIds)
	if err != nil {
		return nil, err
	}
	productsMap := lo.GroupBy(products, func(item *promotion.PromotionPointProduct) int64 {
		return item.ActivityID
	})

	// 2. 获取 SPU 信息
	spuIds := lo.Map(activityList, func(item *promotion.PromotionPointActivity, _ int) int64 {
		return item.SpuID
	})
	spuList, err := h.spuSvc.GetSpuList(c, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := lo.KeyBy(spuList, func(item *resp.ProductSpuResp) int64 {
		return item.ID
	})

	// 3. 组装结果
	result := make([]*resp.AppPointActivityRespVO, len(activityList))
	for i, activity := range activityList {
		vo := &resp.AppPointActivityRespVO{
			ID:    activity.ID,
			SpuID: activity.SpuID,
		}

		// 设置 Product 信息 (Min Point/Price)
		if actProducts, ok := productsMap[activity.ID]; ok && len(actProducts) > 0 {
			minProduct := lo.MinBy(actProducts, func(a, b *promotion.PromotionPointProduct) bool {
				return a.Point < b.Point
			})
			if minProduct != nil {
				vo.Point = minProduct.Point
				vo.Price = minProduct.Price
			}
		}

		// 设置 SPU 信息
		if spu, ok := spuMap[activity.SpuID]; ok {
			vo.SpuName = spu.Name
			vo.PicUrl = spu.PicURL
			vo.MarketPrice = spu.MarketPrice
		}

		result[i] = vo
	}
	return result, nil
}
