package product

import (
	stdcontext "context"
	"fmt"
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	productSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AppProductSpuHandler struct {
	spuSvc         *productSvc.ProductSpuService
	propertySvc    *productSvc.ProductPropertyService
	historySvc     *productSvc.ProductBrowseHistoryService
	memberUserSvc  *memberSvc.MemberUserService
	memberLevelSvc *memberSvc.MemberLevelService
}

func NewAppProductSpuHandler(spuSvc *productSvc.ProductSpuService, propertySvc *productSvc.ProductPropertyService, historySvc *productSvc.ProductBrowseHistoryService, memberUserSvc *memberSvc.MemberUserService, memberLevelSvc *memberSvc.MemberLevelService) *AppProductSpuHandler {
	return &AppProductSpuHandler{
		spuSvc:         spuSvc,
		propertySvc:    propertySvc,
		historySvc:     historySvc,
		memberUserSvc:  memberUserSvc,
		memberLevelSvc: memberLevelSvc,
	}
}

// GetSpuDetail 获得 SPU 详情 (Trigger History) - 对齐Java版本逻辑
// @Summary 获得 SPU 详情
// @Tags 用户 APP - 商品 SPU
// @Produce json
// @Param id query int true "SPU 编号"
// @Success 200 {object} resp.AppProductSpuDetailResp
// @Router /app-api/product/spu/get-detail [get]
func (h *AppProductSpuHandler) GetSpuDetail(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		if !middleware.ValidateProductParams(c, err) {
			return
		}
	}

	// 调用Service层获取SPU详情，包含状态检查和错误处理
	spu, skus, err := h.spuSvc.GetSpuDetail(c, id)
	if err != nil {
		middleware.HandleProductError(c, err)
		return
	}

	// 增加浏览量，对齐Java版本逻辑
	userID := context.GetLoginUserID(c)
	if userID > 0 {
		_ = h.historySvc.CreateBrowseHistory(c, userID, id)
	}
	_ = h.spuSvc.UpdateBrowseCount(c, id, 1)

	// 计算VIP价格，对齐Java版本的会员折扣逻辑
	discountPercent := 100
	if userID > 0 {
		user, _ := h.memberUserSvc.GetUser(c, userID)
		if user != nil && user.LevelID > 0 {
			level, _ := h.memberLevelSvc.GetLevel(c, user.LevelID)
			if level != nil {
				discountPercent = level.DiscountPercent
			}
		}
	}

	// 在Handler层组装VO响应数据
	res := h.convertSpuDetailResp(spu, skus)

	// 为每个SKU计算VIP价格
	for i := range res.Skus {
		if discountPercent < 100 {
			res.Skus[i].VipPrice = int(int64(res.Skus[i].Price) * int64(discountPercent) / 100)
		} else {
			res.Skus[i].VipPrice = res.Skus[i].Price
		}
	}

	response.WriteSuccess(c, res)
}

// GetSpuListByIds 获得商品 SPU 列表 - 对齐Java版本/app-api/product/spu/list-by-ids接口
// @Summary 获得商品 SPU 列表
// @Tags 用户 APP - 商品 SPU
// @Produce json
// @Param ids query string true "编号列表，逗号分隔"
// @Success 200 {object} []resp.AppProductSpuResp
// @Router /app-api/product/spu/list-by-ids [get]
func (h *AppProductSpuHandler) GetSpuListByIds(c *gin.Context) {
	// 参数解析和验证 - 对齐Java版本的@RequestParam("ids") Set<Long> ids
	idsStr := c.Query("ids")
	if idsStr == "" {
		// 返回空数组而不是null，对齐Java版本Collections.emptyList()
		response.WriteSuccess(c, []resp.AppProductSpuResp{})
		return
	}

	ids := utils.SplitToInt64(idsStr)
	if len(ids) == 0 {
		// 返回空数组而不是null，对齐Java版本Collections.emptyList()
		response.WriteSuccess(c, []resp.AppProductSpuResp{})
		return
	}

	// 调用Service层获取SPU列表
	list, err := h.spuSvc.GetSpuList(c, ids)
	if err != nil {
		middleware.HandleProductError(c, err)
		return
	}

	// 如果列表为空，返回空数组 - 对齐Java版本CollUtil.isEmpty(list)处理
	if len(list) == 0 {
		response.WriteSuccess(c, []resp.AppProductSpuResp{})
		return
	}

	// 转换为App端响应格式，对齐Java版本的AppProductSpuRespVO结构
	// 注意：Java版本中没有VIP价格计算，这里移除VIP价格相关逻辑以完全对齐
	resList := make([]resp.AppProductSpuResp, len(list))
	for i, spu := range list {
		// 确保字段完全对齐Java版本的AppProductSpuRespVO
		// Java版本：list.forEach(spu -> spu.setSalesCount(spu.getSalesCount() + spu.getVirtualSalesCount()));
		resList[i] = resp.AppProductSpuResp{
			ID:            spu.ID,
			Name:          spu.Name,
			Introduction:  spu.Introduction, // 对齐Java版本字段
			CategoryID:    spu.CategoryID,   // 对齐Java版本字段
			PicURL:        spu.PicURL,
			SliderPicURLs: spu.SliderPicURLs,  // 对齐Java版本字段，确保返回[]而不是null
			SpecType:      bool(spu.SpecType), // 转换BitBool为bool
			Price:         spu.Price,
			MarketPrice:   spu.MarketPrice,
			Stock:         spu.Stock,
			SalesCount:    spu.SalesCount,    // 已在Service层合并了虚拟销量
			DeliveryTypes: spu.DeliveryTypes, // 对齐Java版本字段，确保返回[]而不是null
		}
	}
	response.WriteSuccess(c, resList)
}

// GetSpuPage 获得商品 SPU 分页 - 对齐Java版本/app-api/product/spu/page接口
// @Summary 获得商品 SPU 分页
// @Tags 用户 APP - 商品 SPU
// @Produce json
// @Success 200 {object} pagination.PageResult[resp.AppProductSpuResp]
// @Router /app-api/product/spu/page [get]
func (h *AppProductSpuHandler) GetSpuPage(c *gin.Context) {
	var r req.AppProductSpuPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		if !middleware.ValidateProductParams(c, err) {
			return
		}
	}

	// 调用 Service
	pageResult, err := h.spuSvc.GetSpuPageForApp(c, &r)
	if err != nil {
		middleware.HandleProductError(c, err)
		return
	}

	// 如果列表为空，返回空分页结果 - 对齐Java版本CollUtil.isEmpty(pageResult.getList())处理
	if len(pageResult.List) == 0 {
		response.WriteSuccess(c, pagination.PageResult[resp.AppProductSpuResp]{
			List:  []resp.AppProductSpuResp{},
			Total: pageResult.Total,
		})
		return
	}

	// 转换为App端响应格式，对齐Java版本的AppProductSpuRespVO结构
	// Java版本：pageResult.getList().forEach(spu -> spu.setSalesCount(spu.getSalesCount() + spu.getVirtualSalesCount()));
	list := make([]resp.AppProductSpuResp, len(pageResult.List))
	for i, spu := range pageResult.List {
		list[i] = resp.AppProductSpuResp{
			ID:            spu.ID,
			Name:          spu.Name,
			Introduction:  spu.Introduction, // 对齐Java版本字段
			CategoryID:    spu.CategoryID,   // 对齐Java版本字段
			PicURL:        spu.PicURL,
			SliderPicURLs: spu.SliderPicURLs,  // 对齐Java版本字段，确保返回[]而不是null
			SpecType:      bool(spu.SpecType), // 转换BitBool为bool
			Price:         spu.Price,
			MarketPrice:   spu.MarketPrice,
			Stock:         spu.Stock,
			SalesCount:    spu.SalesCount,    // 已在Service层合并了虚拟销量
			DeliveryTypes: spu.DeliveryTypes, // 对齐Java版本字段，确保返回[]而不是null
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.AppProductSpuResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

// collectPropertyIDs 收集所有SKU中的PropertyID (复用Admin端逻辑)
func (h *AppProductSpuHandler) collectPropertyIDs(skus []*product.ProductSku) []int64 {
	propertyIDSet := make(map[int64]bool)
	for _, sku := range skus {
		for _, prop := range sku.Properties {
			if prop.PropertyID > 0 { // 排除默认规格(ID=0)
				propertyIDSet[prop.PropertyID] = true
			}
		}
	}

	propertyIDs := make([]int64, 0, len(propertyIDSet))
	for id := range propertyIDSet {
		propertyIDs = append(propertyIDs, id)
	}
	return propertyIDs
}

// getPropertyNameMap 获取属性ID到名称的映射 (复用Admin端逻辑)
func (h *AppProductSpuHandler) getPropertyNameMap(ctx stdcontext.Context, propertyIDs []int64) map[int64]string {
	if len(propertyIDs) == 0 {
		return make(map[int64]string)
	}

	// 调用PropertyService批量获取属性信息
	properties, err := h.propertySvc.GetPropertyListByIds(ctx, propertyIDs)
	if err != nil {
		// 记录错误但不影响其他字段，使用默认值
		propertyMap := make(map[int64]string)
		for _, id := range propertyIDs {
			propertyMap[id] = h.getPropertyName(id, propertyMap)
		}
		return propertyMap
	}

	// 构建映射
	propertyMap := make(map[int64]string, len(properties))
	for _, prop := range properties {
		propertyMap[prop.ID] = prop.Name
	}
	return propertyMap
}

// getPropertyName 获取属性名称，处理默认规格 (复用Admin端逻辑)
func (h *AppProductSpuHandler) getPropertyName(propertyID int64, propertyMap map[int64]string) string {
	if propertyID == 0 {
		return "默认" // 默认规格
	}

	if name, exists := propertyMap[propertyID]; exists {
		return name
	}

	// 属性不存在时的默认处理
	return fmt.Sprintf("属性%d", propertyID)
}

// convertSpuDetailResp 在Handler层组装SPU详情响应VO
func (h *AppProductSpuHandler) convertSpuDetailResp(spu *product.ProductSpu, skus []*product.ProductSku) *resp.AppProductSpuDetailResp {
	// 确保数组字段返回[]而不是null，对齐Java版本处理逻辑
	sliderPicURLs := spu.SliderPicURLs
	if sliderPicURLs == nil {
		sliderPicURLs = []string{}
	}

	// 1. 收集所有需要查询的PropertyID
	propertyIDs := h.collectPropertyIDs(skus)

	// 2. 批量查询属性名称
	propertyMap := h.getPropertyNameMap(stdcontext.Background(), propertyIDs)

	// 3. 转换SKU数组，填充PropertyName
	skuResps := make([]resp.AppProductSpuDetailSkuResp, 0, len(skus))
	for _, sku := range skus {
		// 转换SKU属性数组
		properties := make([]resp.AppProductPropertyValueDetail, 0)
		if sku.Properties != nil {
			for _, prop := range sku.Properties {
				propertyName := h.getPropertyName(prop.PropertyID, propertyMap)
				properties = append(properties, resp.AppProductPropertyValueDetail{
					PropertyID:   prop.PropertyID,
					PropertyName: propertyName,
					ValueID:      prop.ValueID,
					ValueName:    prop.ValueName,
				})
			}
		}

		skuResps = append(skuResps, resp.AppProductSpuDetailSkuResp{
			ID:          sku.ID,
			Properties:  properties,
			Price:       int(sku.Price),
			MarketPrice: int(sku.MarketPrice),
			VipPrice:    int(sku.Price), // 默认VIP价格等于原价，在Handler层计算折扣
			PicURL:      sku.PicURL,
			Stock:       int(sku.Stock),
			Weight:      sku.Weight,
			Volume:      sku.Volume,
		})
	}

	return &resp.AppProductSpuDetailResp{
		ID:            spu.ID,
		Name:          spu.Name,
		Introduction:  spu.Introduction, // 确保返回完整介绍文本
		Description:   spu.Description,  // 确保返回完整描述文本
		CategoryID:    spu.CategoryID,   // 确保返回正确分类ID
		PicURL:        spu.PicURL,
		SliderPicURLs: sliderPicURLs, // 确保返回[]而不是null
		SpecType:      bool(spu.SpecType),
		Price:         int(spu.Price),
		MarketPrice:   int(spu.MarketPrice),
		Stock:         int(spu.Stock),
		SalesCount:    int(spu.SalesCount), // 已在Service层合并了虚拟销量
		Skus:          skuResps,
	}
}
