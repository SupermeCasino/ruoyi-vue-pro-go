package product

import (
	"context"
	"fmt"
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	productSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductSpuHandler struct {
	svc         *productSvc.ProductSpuService
	propertySvc *productSvc.ProductPropertyService
}

func NewProductSpuHandler(svc *productSvc.ProductSpuService, propertySvc *productSvc.ProductPropertyService) *ProductSpuHandler {
	return &ProductSpuHandler{
		svc:         svc,
		propertySvc: propertySvc,
	}
}

// CreateSpu 创建 SPU
// @Router /admin-api/product/spu/create [post]
func (h *ProductSpuHandler) CreateSpu(c *gin.Context) {
	var r req.ProductSpuSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateSpu(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(id))
}

// UpdateSpu 更新 SPU
// @Router /admin-api/product/spu/update [put]
func (h *ProductSpuHandler) UpdateSpu(c *gin.Context) {
	var r req.ProductSpuSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateSpu(c, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// UpdateSpuStatus 更新 SPU 状态
// @Router /admin-api/product/spu/update-status [put]
func (h *ProductSpuHandler) UpdateSpuStatus(c *gin.Context) {
	var r req.ProductSpuUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateSpuStatus(c, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// DeleteSpu 删除 SPU
// @Router /admin-api/product/spu/delete [delete]
func (h *ProductSpuHandler) DeleteSpu(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteSpu(c, id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// GetSpuDetail 获得 SPU 详情
// @Router /admin-api/product/spu/get-detail [get]
func (h *ProductSpuHandler) GetSpuDetail(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	spu, skus, err := h.svc.GetSpuDetail(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	// 在Handler层组装响应数据
	res := h.convertSpuDetailResp(spu, skus)
	c.JSON(200, response.Success(res))
}

// GetSpuPage 获得 SPU 分页
// @Router /admin-api/product/spu/page [get]
func (h *ProductSpuHandler) GetSpuPage(c *gin.Context) {
	var r req.ProductSpuPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	res, err := h.svc.GetSpuPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(res))
}

// GetTabsCount 获得 SPU Tab 统计
// @Router /admin-api/product/spu/get-count [get]
func (h *ProductSpuHandler) GetTabsCount(c *gin.Context) {
	res, err := h.svc.GetTabsCount(c)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(res))
}

// GetSpuSimpleList 获得 SPU 精简列表
// @Router /admin-api/product/spu/list-all-simple [get]
func (h *ProductSpuHandler) GetSpuSimpleList(c *gin.Context) {
	res, err := h.svc.GetSpuSimpleList(c)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(res))
}

// GetSpuList 根据 ID 列表获得 SPU 详情列表
// @Router /admin-api/product/spu/list [get]
func (h *ProductSpuHandler) GetSpuList(c *gin.Context) {
	var r req.ProductSpuListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	res, err := h.svc.GetSpuList(c, r.SpuIDs)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(res))
}

// ExportSpuList 导出商品列表
// @Router /admin-api/product/spu/export [get]
func (h *ProductSpuHandler) ExportSpuList(c *gin.Context) {
	var r req.ProductSpuPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	// 导出不分页
	r.PageNo = 1
	r.PageSize = 10000

	res, err := h.svc.GetSpuPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// 转换为导出 VO
	exportList := make([]*resp.ProductSpuExportVO, len(res.List))
	for i, spu := range res.List {
		exportList[i] = &resp.ProductSpuExportVO{
			ID:         spu.ID,
			Name:       spu.Name,
			CategoryID: spu.CategoryID,
			Price:      spu.Price,
			Stock:      spu.Stock,
			Status:     spu.Status,
			SalesCount: spu.SalesCount,
		}
	}

	if err := excel.WriteExcel(c, "商品列表.xlsx", "数据", exportList); err != nil {
		c.Error(err)
		return
	}
}

// collectPropertyIDs 收集所有SKU中的PropertyID
func (h *ProductSpuHandler) collectPropertyIDs(skus []*product.ProductSku) []int64 {
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

// getPropertyNameMap 获取属性ID到名称的映射
func (h *ProductSpuHandler) getPropertyNameMap(ctx context.Context, propertyIDs []int64) map[int64]string {
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

// getPropertyName 获取属性名称，处理默认规格
func (h *ProductSpuHandler) getPropertyName(propertyID int64, propertyMap map[int64]string) string {
	if propertyID == 0 {
		return "默认" // 默认规格
	}

	if name, exists := propertyMap[propertyID]; exists {
		return name
	}

	// 属性不存在时的默认处理
	return fmt.Sprintf("属性%d", propertyID)
}

// convertSpuDetailResp 在Admin Handler层组装SPU详情响应VO
func (h *ProductSpuHandler) convertSpuDetailResp(spu *product.ProductSpu, skus []*product.ProductSku) *resp.ProductSpuResp {
	// 确保数组字段返回[]而不是null
	sliderPicURLs := spu.SliderPicURLs
	if sliderPicURLs == nil {
		sliderPicURLs = []string{}
	}

	deliveryTypes := spu.DeliveryTypes
	if deliveryTypes == nil {
		deliveryTypes = []int{}
	}

	// 1. 收集所有需要查询的PropertyID
	propertyIDs := h.collectPropertyIDs(skus)

	// 2. 批量查询属性名称
	propertyMap := h.getPropertyNameMap(context.Background(), propertyIDs)

	// 3. 转换SKU数组，填充PropertyName
	skuResps := make([]*resp.ProductSkuResp, 0, len(skus))
	for _, sku := range skus {
		// 转换SKU属性数组
		properties := make([]resp.ProductSkuPropertyResp, 0)
		if sku.Properties != nil {
			for _, prop := range sku.Properties {
				propertyName := h.getPropertyName(prop.PropertyID, propertyMap)
				properties = append(properties, resp.ProductSkuPropertyResp{
					PropertyID:   prop.PropertyID,
					PropertyName: propertyName,
					ValueID:      prop.ValueID,
					ValueName:    prop.ValueName,
				})
			}
		}

		skuResps = append(skuResps, &resp.ProductSkuResp{
			ID:          sku.ID,
			SpuID:       sku.SpuID,
			Properties:  properties,
			Price:       sku.Price,
			MarketPrice: sku.MarketPrice,
			CostPrice:   sku.CostPrice,
			BarCode:     sku.BarCode,
			PicURL:      sku.PicURL,
			Stock:       sku.Stock,
			Weight:      sku.Weight,
			Volume:      sku.Volume,
		})
	}

	return &resp.ProductSpuResp{
		ID:                 spu.ID,
		Name:               spu.Name,
		Keyword:            spu.Keyword,
		Introduction:       spu.Introduction,
		Description:        spu.Description,
		CategoryID:         spu.CategoryID,
		BrandID:            spu.BrandID,
		PicURL:             spu.PicURL,
		SliderPicURLs:      sliderPicURLs,
		Sort:               spu.Sort,
		Status:             spu.Status,
		SpecType:           bool(spu.SpecType),
		Price:              spu.Price,
		MarketPrice:        spu.MarketPrice,
		CostPrice:          spu.CostPrice,
		Stock:              spu.Stock,
		DeliveryTypes:      deliveryTypes,
		DeliveryTemplateID: spu.DeliveryTemplateID,
		GiveIntegral:       spu.GiveIntegral,
		SubCommissionType:  bool(spu.SubCommissionType),
		SalesCount:         spu.SalesCount,
		VirtualSalesCount:  spu.VirtualSalesCount,
		BrowseCount:        spu.BrowseCount,
		CreateTime:         spu.CreateTime,
		Skus:               skuResps,
	}
}
