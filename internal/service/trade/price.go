package trade

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	productModel "github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	pkgErrors "github.com/wxlbd/ruoyi-mall-go/pkg/errors"
)

// TradePriceService 价格计算服务
type TradePriceService struct {
	productSkuSvc      *product.ProductSkuService
	productSpuSvc      *product.ProductSpuService
	couponSvc          *promotion.CouponUserService
	rewardActivitySvc  *promotion.RewardActivityService
	seckillActivitySvc *promotion.SeckillActivityService // 已添加
	memberUserSvc      *memberSvc.MemberUserService
	memberLevelSvc     *memberSvc.MemberLevelService
	deliveryFreightSvc *DeliveryExpressTemplateService
	memberAddressSvc   *memberSvc.MemberAddressService
	memberConfigSvc    *memberSvc.MemberConfigService
}

func NewTradePriceService(
	productSkuSvc *product.ProductSkuService,
	productSpuSvc *product.ProductSpuService,
	couponSvc *promotion.CouponUserService,
	rewardActivitySvc *promotion.RewardActivityService,
	seckillActivitySvc *promotion.SeckillActivityService, // 已添加
	memberUserSvc *memberSvc.MemberUserService,
	memberLevelSvc *memberSvc.MemberLevelService,
	deliveryFreightSvc *DeliveryExpressTemplateService,
	memberAddressSvc *memberSvc.MemberAddressService,
	memberConfigSvc *memberSvc.MemberConfigService,
) *TradePriceService {
	return &TradePriceService{
		productSkuSvc:      productSkuSvc,
		productSpuSvc:      productSpuSvc,
		couponSvc:          couponSvc,
		rewardActivitySvc:  rewardActivitySvc,
		seckillActivitySvc: seckillActivitySvc, // 已添加
		memberUserSvc:      memberUserSvc,
		memberLevelSvc:     memberLevelSvc,
		deliveryFreightSvc: deliveryFreightSvc,
		memberAddressSvc:   memberAddressSvc,
		memberConfigSvc:    memberConfigSvc,
	}
}

// TradePriceCalculateReqBO 价格计算请求BO
type TradePriceCalculateReqBO struct {
	UserID            int64
	CouponID          *int64
	PointStatus       bool
	DeliveryType      int
	AddressID         *int64
	PickUpStoreID     *int64
	SeckillActivityId int64 // 已添加
	Items             []TradePriceCalculateItemBO
}

type TradePriceCalculateItemBO struct {
	SkuID    int64
	Count    int
	CartID   int64
	Selected bool
}

// TradePriceCalculateRespBO 价格计算响应BO
type TradePriceCalculateRespBO struct {
	Type       int
	Price      TradePriceCalculatePriceBO
	Items      []TradePriceCalculateItemRespBO
	CouponID   int64
	TotalPoint int
	UsePoint   int
	GivePoint  int
	Success    bool
}

type TradePriceCalculatePriceBO struct {
	TotalPrice    int
	DiscountPrice int
	DeliveryPrice int
	CouponPrice   int
	PointPrice    int
	VipPrice      int
	PayPrice      int
}

type TradePriceCalculateItemRespBO struct {
	SpuID         int64
	SkuID         int64
	Count         int
	CartID        int64
	Selected      bool
	Price         int
	DiscountPrice int
	DeliveryPrice int
	CouponPrice   int
	PointPrice    int
	UsePoint      int
	VipPrice      int
	PayPrice      int
	SpuName       string
	PicURL        string
	CategoryID    int64
	DeliveryTypes []int
	GivePoint     int
	Properties    []resp.ProductSkuPropertyResp
}

// dividePrice 按支付金额比例分摊折扣金额
// 对应 Java：TradePriceCalculatorHelper#dividePrice
func dividePrice(items []TradePriceCalculateItemRespBO, totalDiscount int) []int {
	if len(items) == 0 {
		return []int{}
	}

	// 计算所有【已选中】项的总支付金额
	// Java: Integer total = calculateTotalPayPrice(orderItems); // 只计算 selected=true 的
	totalPayPrice := 0
	for _, item := range items {
		if item.Selected {
			totalPayPrice += item.PayPrice
		}
	}

	if totalPayPrice == 0 {
		return make([]int, len(items))
	}

	// 按比例分摊
	dividedPrices := make([]int, len(items))
	remainPrice := totalDiscount
	lastSelectedIndex := -1

	// 找到最后一个选中项的索引
	for i := len(items) - 1; i >= 0; i-- {
		if items[i].Selected {
			lastSelectedIndex = i
			break
		}
	}

	for i := 0; i < len(items); i++ {
		// 1. 如果是未选中，则分摊为 0
		if !items[i].Selected {
			dividedPrices[i] = 0
			continue
		}
		// 2. 如果选中，则按照百分比进行分摊
		if i < lastSelectedIndex {
			// 前 n-1 项按比例计算
			dividedPrices[i] = int(int64(totalDiscount) * int64(items[i].PayPrice) / int64(totalPayPrice))
			remainPrice -= dividedPrices[i]
		} else {
			// 最后一项用剩余金额（避免舍入误差）
			dividedPrices[i] = remainPrice
		}
	}

	return dividedPrices
}

// CalculateOrderPrice 订单价格计算
// Java：TradePriceServiceImpl#calculateOrderPrice
func (s *TradePriceService) CalculateOrderPrice(ctx context.Context, req *TradePriceCalculateReqBO) (*TradePriceCalculateRespBO, error) {
	// 1. 获得 SKU 编号数组
	var skuIDs []int64
	for _, item := range req.Items {
		skuIDs = append(skuIDs, item.SkuID)
	}

	// 2. 获得 SKU 列表
	skus, err := s.productSkuSvc.GetSkuList(ctx, skuIDs)
	if err != nil {
		return nil, err
	}
	skuMap := make(map[int64]*resp.ProductSkuResp)
	var spuIDs []int64
	spuIdMapKeys := make(map[int64]bool) // 避免查询重复
	for _, sku := range skus {
		skuMap[sku.ID] = sku
		if !spuIdMapKeys[sku.SpuID] {
			spuIDs = append(spuIDs, sku.SpuID)
			spuIdMapKeys[sku.SpuID] = true
		}
	}

	// 2.1 获得 SPU 列表
	spus, err := s.productSpuSvc.GetSpuList(ctx, spuIDs)
	if err != nil {
		return nil, err
	}
	spuMap := make(map[int64]*resp.ProductSpuResp)
	for _, spu := range spus {
		spuMap[spu.ID] = spu
	}

	// 3. 初始化结果
	respBO := &TradePriceCalculateRespBO{
		Price:   TradePriceCalculatePriceBO{},
		Items:   make([]TradePriceCalculateItemRespBO, 0),
		Success: true,
	}

	var totalPrice, totalPayPrice int

	// 4. 计算秒杀优惠（订单8）
	// 让我们使用 `respBO.Type` 作为标志。
	if req.SeckillActivityId > 0 {
		respBO.Type = 1 // 秒杀
		if len(req.Items) != 1 {
			return nil, errors.New("秒杀时，只允许选择一个商品")
		}
		item := req.Items[0]
		// 验证秒杀参与
		_, _, err := s.seckillActivitySvc.ValidateJoinSeckill(ctx, req.SeckillActivityId, item.SkuID, item.Count)
		if err != nil {
			return nil, err
		}
	} else {
		respBO.Type = 0 // 正常
	}

	// 4. 计算 VIP 会员折扣
	levelDiscountPercent := 100
	if req.UserID > 0 {
		user, _ := s.memberUserSvc.GetUser(ctx, req.UserID)
		if user != nil && user.LevelID > 0 {
			level, _ := s.memberLevelSvc.GetLevel(ctx, user.LevelID)
			if level != nil && level.Status == 0 {
				levelDiscountPercent = level.DiscountPercent
			}
		}
	}

	var totalVipPrice int

	// 4. 循环计算商品项
	for _, item := range req.Items {
		if !item.Selected {
			continue
		}
		sku, ok := skuMap[item.SkuID]
		if !ok {
			return nil, errors.New("商品不存在")
		}
		spu, ok := spuMap[sku.SpuID]
		if !ok {
			return nil, errors.New("商品 SPU 不存在")
		}

		// 计算商品价格
		itemPrice := sku.Price
		itemPayPrice := itemPrice * item.Count

		itemVipSavings := 0
		seckillDiscount := 0

		// 秒杀逻辑（订单8）
		if respBO.Type == 1 { // 秒杀
			_, seckillProd, err := s.seckillActivitySvc.ValidateJoinSeckill(ctx, req.SeckillActivityId, sku.ID, item.Count)
			if err != nil {
				return nil, err
			}
			seckillTotal := seckillProd.SeckillPrice * item.Count
			seckillDiscount = itemPayPrice - seckillTotal
			itemPayPrice = seckillTotal
		} else {
			// 正常订单逻辑（VIP）（订单10）
			// 严格对齐 Java：TradeDiscountActivityPriceCalculator.calculateVipPrice
			// Java: Integer newPrice = calculateRatePrice(orderItem.getPayPrice(), level.getDiscountPercent())
			// Java: return orderItem.getPayPrice() - newPrice;
			if levelDiscountPercent < 100 {
				// 基于 PayPrice 计算，不是基于 Price！
				vipTotal := int(int64(itemPayPrice) * int64(levelDiscountPercent) / 100)
				itemVipSavings = itemPayPrice - vipTotal
				itemPayPrice = vipTotal
			}
		}

		itemResp := TradePriceCalculateItemRespBO{
			SpuID:         sku.SpuID,
			SkuID:         sku.ID,
			Count:         item.Count,
			CartID:        item.CartID,
			Selected:      item.Selected,
			Price:         itemPrice,
			PayPrice:      itemPayPrice,
			PicURL:        sku.PicURL,
			Properties:    sku.Properties,
			SpuName:       spu.Name,
			CategoryID:    spu.CategoryID,
			DiscountPrice: seckillDiscount,
			VipPrice:      itemVipSavings,
		}

		totalPrice += itemPrice * item.Count
		totalPayPrice += itemPayPrice
		totalVipPrice += itemVipSavings

		respBO.Items = append(respBO.Items, itemResp)
	}

	// 5. 计算满减/满折活动（订单20）
	// 仅适用于正常订单
	// Java TradeRewardActivityPriceCalculator 检查：
	// if (!TradeOrderTypeEnum.isNormal(result.getType())) { return; }

	activityDiscount := 0
	if respBO.Type == 0 { // 仅正常订单
		matchItems := make([]promotion.ActivityMatchItem, 0)
		for _, item := range respBO.Items {
			matchItems = append(matchItems, promotion.ActivityMatchItem{
				SkuID:      item.SkuID,
				SpuID:      item.SpuID,
				CategoryID: item.CategoryID,
				Price:      item.Price,
				Count:      item.Count,
			})
		}
		var err error
		activityDiscount, _, err = s.rewardActivitySvc.CalculateRewardActivity(ctx, matchItems)
		if err != nil {
			return nil, err
		}

		// 分摊满减折扣到各项
		if activityDiscount > 0 {
			divideActivityDiscounts := dividePrice(respBO.Items, activityDiscount)
			for i := range respBO.Items {
				respBO.Items[i].DiscountPrice += divideActivityDiscounts[i]
				// Java: TradePriceCalculatorHelper.recountPayPrice(orderItem)
				item := &respBO.Items[i]
				item.PayPrice = item.Price*item.Count - item.DiscountPrice + item.DeliveryPrice - item.CouponPrice - item.PointPrice - item.VipPrice
				if item.PayPrice < 0 {
					item.PayPrice = 0
				}
			}
		}
	}

	// 汇总所有折扣
	var totalDiscountPrice int
	for _, it := range respBO.Items {
		totalDiscountPrice += it.DiscountPrice
	}

	// 设置订单总价
	respBO.Price.TotalPrice = totalPrice
	respBO.Price.DiscountPrice = totalDiscountPrice
	respBO.Price.VipPrice = totalVipPrice

	// 计算最终支付金额
	// PayPrice = TotalPrice - DiscountPrice - VipPrice
	payPrice := respBO.Price.TotalPrice - respBO.Price.DiscountPrice - respBO.Price.VipPrice
	if payPrice < 0 {
		payPrice = 0
	}
	respBO.Price.PayPrice = payPrice

	// 6. 计算优惠券优惠（订单30）
	// Java: TradeCouponPriceCalculator: 只有【普通】订单，才允许使用优惠劵
	if respBO.Type == 0 && req.CouponID != nil && *req.CouponID > 0 {
		var spuIDs []int64
		var categoryIDs []int64
		// 辅助Map避免重复
		spuMapForCoupon := make(map[int64]bool)
		catMapForCoupon := make(map[int64]bool)

		for _, item := range respBO.Items {
			if !spuMapForCoupon[item.SpuID] {
				spuIDs = append(spuIDs, item.SpuID)
				spuMapForCoupon[item.SpuID] = true
			}
			if !catMapForCoupon[item.CategoryID] {
				categoryIDs = append(categoryIDs, item.CategoryID)
				catMapForCoupon[item.CategoryID] = true
			}
		}

		couponPrice, err := s.couponSvc.CalculateCoupon(ctx, req.UserID, *req.CouponID, int64(respBO.Price.PayPrice), spuIDs, categoryIDs)
		if err != nil {
			return nil, err
		}

		couponPriceInt := int(couponPrice)
		if couponPriceInt > 0 {
			// 分摊优惠券折扣到各项
			divideCouponPrices := dividePrice(respBO.Items, couponPriceInt)
			for i := range respBO.Items {
				respBO.Items[i].CouponPrice += divideCouponPrices[i]
				// Java: TradePriceCalculatorHelper.recountPayPrice(orderItem)
				item := &respBO.Items[i]
				item.PayPrice = item.Price*item.Count - item.DiscountPrice + item.DeliveryPrice - item.CouponPrice - item.PointPrice - item.VipPrice
				if item.PayPrice < 0 {
					item.PayPrice = 0
				}
			}
		}

		respBO.CouponID = *req.CouponID
		respBO.Price.CouponPrice = couponPriceInt
		respBO.Price.PayPrice -= couponPriceInt
		if respBO.Price.PayPrice < 0 {
			respBO.Price.PayPrice = 0
		}
	}

	// 7. 计算积分抵扣（订单40）
	if req.PointStatus && req.UserID > 0 {
		config, _ := s.memberConfigSvc.GetConfig(ctx)
		user, _ := s.memberUserSvc.GetUser(ctx, req.UserID)
		if config != nil && config.PointTradeDeductEnable > 0 && user != nil && user.Point > 0 {
			// 7.1 计算积分抵扣金额
			// 配置：PointTradeDeductUnitPrice（抵扣单位价格，单位：分）
			deductUnitPrice := config.PointTradeDeductUnitPrice

			if deductUnitPrice > 0 {
				canUsePoints := int(user.Point)

				// 7.2 限制最大积分抵扣数量
				// 配置：PointTradeDeductMaxPrice（注意：Java 中此字段名虽然叫 MaxPrice，但实际逻辑是限制积分数量 MaxPoints）
				if config.PointTradeDeductMaxPrice > 0 {
					if canUsePoints > config.PointTradeDeductMaxPrice {
						canUsePoints = config.PointTradeDeductMaxPrice
					}
				}

				// 7.3 计算抵扣金额
				pointTotalValue := canUsePoints * deductUnitPrice

				// 7.4 限制不超过应付金额
				// Java 逻辑：TradePointUsePriceCalculator -> if (payPrice <= pointPrice) throw exception
				// 严格对齐：禁止 0 元购
				if pointTotalValue >= respBO.Price.PayPrice {
					return nil, pkgErrors.NewBizError(1004003005, "支付金额不能小于等于 0") // PRICE_CALCULATE_PAY_PRICE_ILLEGAL
				}

				// 7.5 分摊积分抵扣到各项
				dividePointPrices := dividePrice(respBO.Items, pointTotalValue)
				for i := range respBO.Items {
					respBO.Items[i].PointPrice += dividePointPrices[i]
					// Java: TradePriceCalculatorHelper.recountPayPrice(orderItem)
					item := &respBO.Items[i]
					item.PayPrice = item.Price*item.Count - item.DiscountPrice + item.DeliveryPrice - item.CouponPrice - item.PointPrice - item.VipPrice
					if item.PayPrice < 0 {
						item.PayPrice = 0
					}
				}

				// 7.6 更新响应
				respBO.UsePoint = canUsePoints
				respBO.TotalPoint = int(user.Point)
				respBO.Price.PointPrice = pointTotalValue
				respBO.Price.PayPrice -= respBO.Price.PointPrice
			}
		}
	}

	// 8. 重新计算每个订单项的 PayPrice
	// 对齐 Java：TradePriceCalculatorHelper#recountPayPrice
	// PayPrice = Price * Count - DiscountPrice + DeliveryPrice - CouponPrice - PointPrice - VipPrice
	for i := range respBO.Items {
		item := &respBO.Items[i]
		item.PayPrice = item.Price*item.Count - item.DiscountPrice + item.DeliveryPrice - item.CouponPrice - item.PointPrice - item.VipPrice
		if item.PayPrice < 0 {
			item.PayPrice = 0
		}
	}

	// 9. 计算运费（订单50）
	// 逻辑：基于商品项计算运费
	deliveryPrice := 0
	if req.DeliveryType == 1 && req.AddressID != nil && *req.AddressID > 0 {
		address, err := s.memberAddressSvc.GetAddress(ctx, req.UserID, *req.AddressID)
		if err != nil {
			return nil, err
		}
		if address != nil {
			templateCountMap := make(map[int64]int)
			templatePriceMap := make(map[int64]int)
			for _, item := range respBO.Items {
				spu := spuMap[item.SpuID]
				if spu != nil {
					templateCountMap[spu.DeliveryTemplateID] += item.Count
					templatePriceMap[spu.DeliveryTemplateID] += item.PayPrice
				}
			}

			for tplID, count := range templateCountMap {
				if tplID > 0 {
					price := templatePriceMap[tplID]
					p, err := s.deliveryFreightSvc.CalculateFreight(ctx, tplID, int(address.AreaID), float64(count), price)
					if err != nil {
						return nil, err
					}
					deliveryPrice += p
				}
			}
		}
	}
	respBO.Price.DeliveryPrice = deliveryPrice
	respBO.Price.PayPrice += deliveryPrice

	return respBO, nil
}

// CalculateProductPrice 计算商品的价格
// 对应 Java: TradePriceServiceImpl#calculateProductPrice
func (s *TradePriceService) CalculateProductPrice(ctx context.Context, userId int64, spuIds []int64) ([]resp.AppTradeProductSettlementResp, error) {
	// 1. 获得 SPU 列表
	spus, err := s.productSpuSvc.GetSpuList(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	if len(spus) == 0 {
		return []resp.AppTradeProductSettlementResp{}, nil
	}

	// 2. 获得 SKU 列表
	skus, err := s.productSkuSvc.GetSkuListBySpuIds(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	skuMap := make(map[int64][]*productModel.ProductSku)
	for _, sku := range skus {
		skuMap[sku.SpuID] = append(skuMap[sku.SpuID], sku)
	}

	// 3. 获得满减送活动
	// 暂时简化，Java 中会查询 RewardActivity
	// activityMap, _ := s.rewardActivitySvc.GetRewardActivityMapBySpuIds(ctx, spuIds)

	// 4. 拼装结果
	var results []resp.AppTradeProductSettlementResp
	for _, spu := range spus {
		spuSkus := skuMap[spu.ID]
		if len(spuSkus) == 0 {
			continue
		}

		var skuResps []resp.Sku
		for _, sku := range spuSkus {
			// 默认原价
			promotionPrice := sku.Price
			// TODO: 查询限时折扣等营销活动

			skuResps = append(skuResps, resp.Sku{
				ID:             sku.ID,
				PromotionPrice: promotionPrice,
				PromotionType:  0, // 无
			})
		}

		results = append(results, resp.AppTradeProductSettlementResp{
			SpuID: spu.ID,
			Skus:  skuResps,
			// RewardActivity: activityMap[spu.ID],
		})
	}
	return results, nil
}
