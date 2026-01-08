package calculators

import (
	"context"

	member2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
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

	// 获取收货地址：缺少 addressId 时尝试使用默认地址（对齐 Java 结算行为）
	var address *member2.AppAddressResp
	var err error
	if req.AddressID != nil && *req.AddressID > 0 {
		address, err = c.memberAddressSvc.GetAddress(ctx, req.UserID, *req.AddressID)
		if err != nil {
			c.LogCalculation(ctx, req, "获取收货地址失败，跳过运费计算", zap.Error(err))
			return nil
		}
	}
	if address == nil {
		address, err = c.memberAddressSvc.GetDefaultAddress(ctx, req.UserID)
		if err != nil {
			c.LogCalculation(ctx, req, "获取默认收货地址失败，跳过运费计算", zap.Error(err))
			return nil
		}
	}
	if address == nil {
		c.LogCalculation(ctx, req, "未找到可用收货地址，跳过运费计算")
		return nil
	}

	// 1. 获取所有选中的商品项，并按运费模板分组
	// 同时统计每个模板的总件数/重量/体积，以及总价格
	templateGroups := make(map[int64]struct {
		totalCount float64
		totalPrice int
	})

	// 1. 收集所有使用的运费模板 ID
	templateIDs := make(map[int64]bool)
	for _, item := range respBO.Items {
		if item.Selected && item.DeliveryTemplateID > 0 {
			templateIDs[item.DeliveryTemplateID] = true
		}
	}

	if len(templateIDs) == 0 {
		c.LogCalculation(ctx, req, "未找到可用运费模板，所有商品均未设置运费模板")
		return nil
	}

	// 2. 批量获取模板信息（含 chargeMode）
	templateChargeModes := make(map[int64]int)
	for templateID := range templateIDs {
		template, err := c.deliveryFreightSvc.GetDeliveryExpressTemplate(ctx, templateID)
		if err != nil || template == nil {
			c.LogCalculation(ctx, req, "获取运费模板失败", zap.Int64("templateId", templateID), zap.Error(err))
			continue
		}
		templateChargeModes[templateID] = template.ChargeMode
	}

	// 3. 按运费模板分组统计，根据 chargeMode 计算 chargeValue
	// 对齐 Java: TradeDeliveryPriceCalculator.getChargeValue
	// chargeMode = 1: 按件数
	// chargeMode = 2: 按重量 (weight * count)
	// chargeMode = 3: 按体积 (volume * count)
	for _, item := range respBO.Items {
		if !item.Selected {
			continue
		}
		if item.DeliveryTemplateID == 0 {
			c.LogCalculation(ctx, req, "商品未设置运费模板，跳过",
				zap.Int64("spuId", item.SpuID),
				zap.Int64("skuId", item.SkuID),
			)
			continue
		}

		chargeMode, ok := templateChargeModes[item.DeliveryTemplateID]
		if !ok {
			continue
		}

		// 根据计费方式计算 chargeValue（对齐 Java getChargeValue）
		var chargeValue float64
		switch chargeMode {
		case 1: // COUNT: 按件数
			chargeValue = float64(item.Count)
		case 2: // WEIGHT: 按重量
			if item.Weight > 0 {
				chargeValue = item.Weight * float64(item.Count)
			}
		case 3: // VOLUME: 按体积
			if item.Volume > 0 {
				chargeValue = item.Volume * float64(item.Count)
			}
		default:
			chargeValue = float64(item.Count) // 默认按件数
		}

		group := templateGroups[item.DeliveryTemplateID]
		group.totalCount += chargeValue
		group.totalPrice += item.PayPrice
		templateGroups[item.DeliveryTemplateID] = group
	}
	if len(templateGroups) == 0 {
		c.LogCalculation(ctx, req, "未找到可用运费模板，所有商品均未设置运费模板",
			zap.Int64("areaId", address.AreaID),
		)
		return tradeSvc.NewTradeError(tradeSvc.ErrorCodeDeliveryTemplateNotExists)
	}

	// 2. 针对每个模板计算运费
	totalDeliveryPrice := 0
	for templateID, group := range templateGroups {
		freight, err := c.deliveryFreightSvc.CalculateFreight(ctx, templateID, int(address.AreaID), group.totalCount, group.totalPrice)
		if err != nil {
			c.LogCalculation(ctx, req, "计算模板运费失败",
				zap.Int64("templateId", templateID),
				zap.Int64("areaId", address.AreaID),
				zap.Float64("totalCount", group.totalCount),
				zap.Int("totalPrice", group.totalPrice),
				zap.Error(err),
			)
			if err.Error() == "该区域不支持配送" {
				return tradeSvc.NewTradeError(tradeSvc.ErrorCodeDeliveryNotSupport)
			}
			if err.Error() == "运费模板不存在" {
				return tradeSvc.NewTradeError(tradeSvc.ErrorCodeDeliveryTemplateNotExists)
			}
			return tradeSvc.NewTradeError(tradeSvc.ErrorCodeDeliveryCalculateError)
		}
		totalDeliveryPrice += freight
	}

	// 3. 更新响应
	respBO.Price.DeliveryPrice = totalDeliveryPrice

	// 4. 分摊运费到商品项，并重新计算 PayPrice
	if totalDeliveryPrice > 0 {
		// 分摊运费
		dividedDeliveryPrices := c.Helper.DividePrice(respBO.Items, totalDeliveryPrice)

		for i := 0; i < len(respBO.Items); i++ {
			if !respBO.Items[i].Selected {
				continue
			}

			// 更新商品项运费
			respBO.Items[i].DeliveryPrice = dividedDeliveryPrices[i]
			// 重新计算商品项 PayPrice (PayPrice = Price * Count - Discount + Delivery ...)
			c.Helper.RecountPayPrice(&respBO.Items[i])
		}

		// 5. 更新整体响应价格 (PayPrice = Sum(Item.PayPrice))
		c.Helper.UpdateResponsePrice(respBO)
	}

	c.LogCalculation(ctx, req, "运费计算完成", zap.Int("totalDeliveryPrice", totalDeliveryPrice))

	return nil
}

// IsApplicable 判断是否适用于当前订单类型
func (c *DeliveryPriceCalculator) IsApplicable(orderType int) bool {
	// 运费计算适用于所有订单类型
	return true
}
