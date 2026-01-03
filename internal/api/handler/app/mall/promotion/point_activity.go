package promotion

import (
	"strconv"
	"strings"

	productContract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	appPromotionContract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	productSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	promotionSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
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

// GetPointActivityPage 获得积分商城活动分页
// 对齐 Java: AppPointActivityController.getPointActivityPage
func (h *AppPointActivityHandler) GetPointActivityPage(c *gin.Context) {
	var reqVO promotion2.PointActivityPageReq
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	pageResult, err := h.svc.GetPointActivityPage(c.Request.Context(), &reqVO)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	if len(pageResult.List) == 0 {
		response.WriteSuccess(c, pagination.PageResult[*appPromotionContract.AppPointActivityRespVO]{
			List:  []*appPromotionContract.AppPointActivityRespVO{},
			Total: pageResult.Total,
		})
		return
	}

	resultList, err := h.buildAppPointActivityRespVOList(c, lo.Map(pageResult.List, func(item promotion.PromotionPointActivity, _ int) *promotion.PromotionPointActivity {
		return &item
	}))
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	response.WriteSuccess(c, pagination.PageResult[*appPromotionContract.AppPointActivityRespVO]{
		List:  resultList,
		Total: pageResult.Total,
	})
}

// GetPointActivity 获得积分商城活动明细
// 对齐 Java: AppPointActivityController.getPointActivity (路径对应 get-detail)
func (h *AppPointActivityHandler) GetPointActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteError(c, 400, "参数错误")
		return
	}

	activity, products, err := h.svc.GetPointActivity(c.Request.Context(), id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	if activity == nil || activity.Status == 0 {
		response.WriteSuccess(c, nil)
		return
	}

	// 拼接数据
	respVO := &appPromotionContract.AppPointActivityDetailRespVO{
		ID:         activity.ID,
		SpuID:      activity.SpuID,
		Status:     activity.Status,
		Stock:      activity.Stock,
		TotalStock: activity.TotalStock,
		Remark:     activity.Remark,
	}

	// 商品列表
	respVO.Products = make([]appPromotionContract.AppPointProductRespVO, len(products))
	for i, p := range products {
		respVO.Products[i] = appPromotionContract.AppPointProductRespVO{
			ID:    p.ID,
			SkuID: p.SkuID,
			Count: p.Count,
			Point: p.Point,
			Price: p.Price,
			Stock: p.Stock,
		}
	}

	// 设置最低积分/价格
	if len(products) > 0 {
		minProduct := lo.MinBy(products, func(a, b *promotion.PromotionPointProduct) bool {
			return a.Point < b.Point
		})
		if minProduct != nil {
			respVO.Point = minProduct.Point
			respVO.Price = minProduct.Price
		}
	}

	response.WriteSuccess(c, respVO)
}

// GetPointActivityListByIds 获得积分商城活动列表
// 对齐 Java: AppPointActivityController.getPointActivityListByIds
func (h *AppPointActivityHandler) GetPointActivityListByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		response.WriteSuccess(c, []*appPromotionContract.AppPointActivityRespVO{})
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
		response.WriteSuccess(c, []*appPromotionContract.AppPointActivityRespVO{})
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

func (h *AppPointActivityHandler) buildAppPointActivityRespVOList(c *gin.Context, activityList []*promotion.PromotionPointActivity) ([]*appPromotionContract.AppPointActivityRespVO, error) {
	if len(activityList) == 0 {
		return []*appPromotionContract.AppPointActivityRespVO{}, nil
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
	spuMap := lo.KeyBy(spuList, func(item *productContract.ProductSpuResp) int64 {
		return item.ID
	})

	// 3. 组装结果
	result := make([]*appPromotionContract.AppPointActivityRespVO, 0, len(activityList))
	for _, activity := range activityList {
		// ✅ 核心修复: 过滤无效 SPU (不存在或非上架状态)
		spu, ok := spuMap[activity.SpuID]
		if !ok || spu.Status != consts.ProductSpuStatusEnable {
			continue
		}

		// ✅ 核心修复: 过滤非开启状态的活动 (CommonStatusEnable = 0)
		if activity.Status != consts.CommonStatusEnable {
			continue
		}

		vo := &appPromotionContract.AppPointActivityRespVO{
			ID:          activity.ID,
			SpuID:       activity.SpuID,
			Status:      activity.Status,
			Stock:       activity.Stock,
			TotalStock:  activity.TotalStock,
			SpuName:     spu.Name,
			PicUrl:      spu.PicURL,
			MarketPrice: spu.MarketPrice,
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

		result = append(result, vo)
	}
	return result, nil
}
