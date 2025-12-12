package trade

import (
	"backend-go/internal/api/resp"
	memberSvc "backend-go/internal/service/member"
	"backend-go/internal/service/product"
	"backend-go/internal/service/promotion"
	"context"
	"errors"
)

// TradePriceService 价格计算服务
type TradePriceService struct {
	productSkuSvc      *product.ProductSkuService
	productSpuSvc      *product.ProductSpuService
	couponSvc          *promotion.CouponUserService
	rewardActivitySvc  *promotion.RewardActivityService
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
		memberUserSvc:      memberUserSvc,
		memberLevelSvc:     memberLevelSvc,
		deliveryFreightSvc: deliveryFreightSvc,
		memberAddressSvc:   memberAddressSvc,
		memberConfigSvc:    memberConfigSvc, // 已添加
	}
}

// TradePriceCalculateReqBO 价格计算请求BO
type TradePriceCalculateReqBO struct {
	UserID        int64
	CouponID      *int64
	PointStatus   bool
	DeliveryType  int
	AddressID     *int64
	PickUpStoreID *int64
	Items         []TradePriceCalculateItemBO
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

	var totalVipPrice int // 仅声明新变量

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

		// 计算 VIP 优惠金额
		// 逻辑：VipPrice (ItemResp) 代表该商品的 VIP 优惠总额
		itemVipSavings := 0
		if levelDiscountPercent < 100 {
			// 优惠金额 = 原价 * 数量 * (1 - 折扣率)
			// 避免精度问题：Price * Count - (Price * Count * Discount / 100)
			vipTotal := int(int64(itemPrice) * int64(item.Count) * int64(levelDiscountPercent) / 100)
			itemVipSavings = itemPayPrice - vipTotal
		}

		itemResp := TradePriceCalculateItemRespBO{
			SpuID:      sku.SpuID,
			SkuID:      sku.ID,
			Count:      item.Count,
			CartID:     item.CartID,
			Selected:   item.Selected,
			Price:      itemPrice,
			PayPrice:   itemPayPrice, // 基础应付金额，后续扣减
			PicURL:     sku.PicURL,
			Properties: sku.Properties,
			SpuName:    spu.Name,
			CategoryID: spu.CategoryID,
			VipPrice:   itemVipSavings, // 存储优惠金额
		}

		totalPrice += itemPayPrice
		totalPayPrice += itemPayPrice
		totalVipPrice += itemVipSavings

		respBO.Items = append(respBO.Items, itemResp)
	}

	// 5. 设置合计
	// 5. 计算满减/满折活动
	matchItems := make([]promotion.ActivityMatchItem, 0)
	for _, item := range respBO.Items {
		matchItems = append(matchItems, promotion.ActivityMatchItem{
			SkuID:      item.SkuID,
			SpuID:      item.SpuID,
			CategoryID: item.CategoryID,
			Price:      item.Price, // 单价（原价）
			Count:      item.Count,
		})
	}
	activityDiscount, _, err := s.rewardActivitySvc.CalculateRewardActivity(ctx, matchItems)
	if err != nil {
		return nil, err
	}

	// 5. 设置合计
	respBO.Price.TotalPrice = totalPrice
	respBO.Price.DiscountPrice = activityDiscount // 设置活动优惠（满减/满折）

	// 3. 计算 VIP 优惠 (Order 10 - 部分 2)
	// 逻辑：限时折扣 (Seckill) 和 VIP 是互斥的 (Order 10)。
	// 目前 Seckill 尚未迁移，因此仅计算 VIP 优惠。后续迁移 Seckill 时需在此处添加互斥判定 (取优惠最大值)。
	// TODO: 迁移 Seckill 后计算 Seckill 优惠
	respBO.Price.VipPrice = totalVipPrice

	// 当前应付金额 (商品 - 优惠)
	payPrice := totalPrice - activityDiscount - totalVipPrice
	if payPrice < 0 {
		payPrice = 0
	}
	respBO.Price.PayPrice = payPrice

	// 4. 计算优惠券优惠（订单30）
	if req.CouponID != nil && *req.CouponID > 0 {
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
		respBO.CouponID = *req.CouponID
		respBO.Price.CouponPrice = int(couponPrice)
		respBO.Price.PayPrice -= int(couponPrice)
		if respBO.Price.PayPrice < 0 {
			respBO.Price.PayPrice = 0
		}
	}

	// 5. 计算积分抵扣 (Order 40)
	if req.PointStatus && req.UserID > 0 {
		config, _ := s.memberConfigSvc.GetConfig(ctx)
		user, _ := s.memberUserSvc.GetUser(ctx, req.UserID)
		if config != nil && config.PointTradeDeductEnable > 0 && user != nil && user.Point > 0 {
			// 5.1 计算积分抵扣金额
			// Conf: PointTradeDeductUnitPrice (抵扣单位价格，单位：分)
			deductUnitPrice := config.PointTradeDeductUnitPrice

			if deductUnitPrice > 0 {
				canUsePoints := int(user.Point)

				// 5.2 限制最大积分抵扣数量
				// Conf: PointTradeDeductMaxPrice (注意：Java 中此字段名虽然叫 MaxPrice，但实际逻辑是限制积分数量 MaxPoints)
				if config.PointTradeDeductMaxPrice > 0 {
					if canUsePoints > config.PointTradeDeductMaxPrice {
						canUsePoints = config.PointTradeDeductMaxPrice
					}
				}

				// 5.3 计算抵扣金额
				pointTotalValue := canUsePoints * deductUnitPrice

				// 5.4 限制不超过应付金额
				// 注意：Java 逻辑中如果 PayPrice <= pointPrice 会抛出异常 (禁止0元购)。
				// 这里为了更好的用户体验，我们做自动截断：最大抵扣金额 = 应付金额。
				if pointTotalValue >= respBO.Price.PayPrice {
					pointTotalValue = respBO.Price.PayPrice
					canUsePoints = pointTotalValue / deductUnitPrice // 重新计算对应积分消耗
					pointTotalValue = canUsePoints * deductUnitPrice
				}

				// 5.5 更新响应
				respBO.UsePoint = canUsePoints
				respBO.TotalPoint = int(user.Point)
				respBO.Price.PointPrice = pointTotalValue
				respBO.Price.PayPrice -= respBO.Price.PointPrice
			}
		}
	}

	// 6. 计算运费（订单50）
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
			// 如果严格需要，重新计算spuMap或假设之前定义的spuMap有效。
			// spuMap在步骤1中定义。
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
	// 重要：运费必须加到PayPrice
	respBO.Price.PayPrice += deliveryPrice

	return respBO, nil
}
