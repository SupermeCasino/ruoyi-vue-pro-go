package calculators

import (
	"context"

	apiResp "github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	tradeSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"go.uber.org/zap"
)

// DeliveryPriceCalculator 运费计算器
type DeliveryPriceCalculator struct {
	*tradeSvc.BasePriceCalculator
	deliveryFreightSvc *tradeSvc.DeliveryExpressTemplateService
	memberAddressSvc   *memberSvc.MemberAddressService
	productSpuSvc      *product.ProductSpuService
}

// NewDeliveryPriceCalculator 创建运费计算器
func NewDeliveryPriceCalculator(
	deliveryFreightSvc *tradeSvc.DeliveryExpressTemplateService,
	memberAddressSvc *memberSvc.MemberAddressService,
	productSpuSvc *product.ProductSpuService,
	helper *tradeSvc.PriceCalculatorHelper,
	logger *zap.Logger,
) *DeliveryPriceCalculator {
	return &DeliveryPriceCalculator{
		BasePriceCalculator: tradeSvc.NewBasePriceCalculator(
			tradeModel.CalculatorNameDelivery,
			tradeModel.OrderDelivery,
			helper,
			logger,
		),
		deliveryFreightSvc: deliveryFreightSvc,
		memberAddressSvc:   memberAddressSvc,
		productSpuSvc:      productSpuSvc,
	}
}

// Calculate 执行运费计算
func (c *DeliveryPriceCalculator) Calculate(ctx context.Context, req *tradeSvc.TradePriceCalculateReqBO, respBO *tradeSvc.TradePriceCalculateRespBO) error {
	// 只有快递配送才需要计算运费
	if req.DeliveryType != 1 {
		c.LogCalculation(ctx, req, "非快递配送，跳过运费计算",
			zap.Int("deliveryType", req.DeliveryType),
		)
		return nil
	}

	c.LogCalculation(ctx, req, "开始执行运费计算")

	// 检查地址 ID 是否存在（商品结算预览时可能没有地址）
	if req.AddressID == nil || *req.AddressID == 0 {
		c.LogCalculation(ctx, req, "未提供地址ID，跳过运费计算")
		return nil
	}

	// 获取收货地址
	address, err := c.memberAddressSvc.GetAddress(ctx, req.UserID, *req.AddressID)
	if err != nil || address == nil {
		c.LogCalculation(ctx, req, "获取收货地址失败，跳过运费计算", zap.Error(err))
		return nil
	}

	// 1. 获取所有选中的商品项，并按运费模板分组
	// 同时统计每个模板的总件数/重量/体积，以及总价格
	templateGroups := make(map[int64]struct {
		totalCount float64
		totalPrice int
	})

	// 获取所有选中的商品 SPU 信息
	var selectedSpuIDs []int64
	for _, item := range respBO.Items {
		if item.Selected {
			selectedSpuIDs = append(selectedSpuIDs, item.SpuID)
		}
	}

	if len(selectedSpuIDs) == 0 {
		return nil
	}

	// 批量获取 SPU 详情以获取运费模板和计费方式
	spuList, err := c.productSpuSvc.GetSpuList(ctx, selectedSpuIDs)
	if err != nil {
		c.LogCalculation(ctx, req, "获取 SPU 信息失败", zap.Error(err))
		return nil
	}

	spuMap := make(map[int64]*apiResp.ProductSpuResp)
	for _, spu := range spuList {
		spuMap[spu.ID] = spu
	}

	// 分组统计
	for _, item := range respBO.Items {
		if !item.Selected {
			continue
		}
		spu := spuMap[item.SpuID]
		if spu == nil || spu.DeliveryTemplateID == 0 {
			continue
		}

		group := templateGroups[spu.DeliveryTemplateID]
		// TODO: 根据模板的计费方式（件/重量/体积）累加 count
		// 目前简化处理，假设都是按件数
		group.totalCount += float64(item.Count)
		group.totalPrice += item.PayPrice
		templateGroups[spu.DeliveryTemplateID] = group
	}

	// 2. 针对每个模板计算运费
	totalDeliveryPrice := 0
	for templateID, group := range templateGroups {
		freight, err := c.deliveryFreightSvc.CalculateFreight(ctx, templateID, int(address.AreaID), group.totalCount, group.totalPrice)
		if err != nil {
			c.LogCalculation(ctx, req, "计算模板运费失败", zap.Int64("templateId", templateID), zap.Error(err))
			continue
		}
		totalDeliveryPrice += freight
	}

	// 3. 更新响应
	respBO.Price.DeliveryPrice = totalDeliveryPrice
	// 分摊运费到商品项（可选，Java 版通常在主服务中处理或分摊）
	if totalDeliveryPrice > 0 {
		c.Helper.DividePrice(respBO.Items, totalDeliveryPrice) // 注意：这里 DividePrice 需要支持分摊运费，目前 Helper 可能只分摊折扣
		// 临时手动分摊到第一个选中的项目
		for i := range respBO.Items {
			if respBO.Items[i].Selected {
				respBO.Items[i].DeliveryPrice = totalDeliveryPrice
				break
			}
		}
	}

	c.LogCalculation(ctx, req, "运费计算完成", zap.Int("totalDeliveryPrice", totalDeliveryPrice))

	return nil
}

// IsApplicable 判断是否适用于当前订单类型
func (c *DeliveryPriceCalculator) IsApplicable(orderType int) bool {
	// 运费计算适用于所有订单类型
	return true
}
