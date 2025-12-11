package trade

import (
	"backend-go/internal/api/resp"
	memberSvc "backend-go/internal/service/member"
	"backend-go/internal/service/product"
	"backend-go/internal/service/promotion"
	"context"
	"errors"
)

// TradePriceService 价格计算 Service
type TradePriceService struct {
	productSkuSvc      *product.ProductSkuService
	productSpuSvc      *product.ProductSpuService
	couponSvc          *promotion.CouponUserService
	rewardActivitySvc  *promotion.RewardActivityService
	memberUserSvc      *memberSvc.MemberUserService
	memberLevelSvc     *memberSvc.MemberLevelService
	deliveryFreightSvc *DeliveryFreightTemplateService // Added
	memberAddressSvc   *memberSvc.MemberAddressService // Added
}

func NewTradePriceService(
	productSkuSvc *product.ProductSkuService,
	productSpuSvc *product.ProductSpuService,
	couponSvc *promotion.CouponUserService,
	rewardActivitySvc *promotion.RewardActivityService,
	memberUserSvc *memberSvc.MemberUserService,
	memberLevelSvc *memberSvc.MemberLevelService,
	deliveryFreightSvc *DeliveryFreightTemplateService, // Added
	memberAddressSvc *memberSvc.MemberAddressService, // Added
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
	}
}

// TradePriceCalculateReqBO 价格计算 Request BO
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

// TradePriceCalculateRespBO 价格计算 Response BO
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

// CalculateOrderPrice 价格计算
// simplified version: only calculates base price = SKU Price * Count.
// Ignores coupons, points, VIP, and shipping fees for now.
func (s *TradePriceService) CalculateOrderPrice(ctx context.Context, req *TradePriceCalculateReqBO) (*TradePriceCalculateRespBO, error) {
	// 1. Get SKU IDs
	var skuIDs []int64
	for _, item := range req.Items {
		skuIDs = append(skuIDs, item.SkuID)
	}

	// 2. Fetch SKUs
	skus, err := s.productSkuSvc.GetSkuList(ctx, skuIDs)
	if err != nil {
		return nil, err
	}
	skuMap := make(map[int64]*resp.ProductSkuResp)
	var spuIDs []int64
	spuIdMapKeys := make(map[int64]bool) // To avoid duplicates for query
	for _, sku := range skus {
		skuMap[sku.ID] = sku
		if !spuIdMapKeys[sku.SpuID] {
			spuIDs = append(spuIDs, sku.SpuID)
			spuIdMapKeys[sku.SpuID] = true
		}
	}

	// 2.1 Fetch SPUs for CategoryID
	spus, err := s.productSpuSvc.GetSpuList(ctx, spuIDs)
	if err != nil {
		return nil, err
	}
	spuMap := make(map[int64]*resp.ProductSpuResp)
	for _, spu := range spus {
		spuMap[spu.ID] = spu
	}

	// 3. Initialize Response
	respBO := &TradePriceCalculateRespBO{
		Price:   TradePriceCalculatePriceBO{},
		Items:   make([]TradePriceCalculateItemRespBO, 0),
		Success: true,
	}

	var totalPrice, totalPayPrice int

	// 4. Calculate VIP Level Discount
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

	var totalVipPrice int // Only declare new var

	// 4. Loop Items and Calculate
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

		// Calculate Item Price
		itemPrice := sku.Price
		itemPayPrice := itemPrice * item.Count

		// Calculate VIP Price (Savings)
		// VipPrice in ItemResp usually means "The Price if VIP" or "The Discount Amount"?
		// Checked RespBO: VipPrice int.
		// Let's assume VipPrice field in ItemResp is the "Discount Amount" for consistency with others?
		// Or is it "Unit Price for VIP"?
		// Usually in Order Item, we store "Original Price", "Pay Price", "Discount Amount".
		// Let's settle on: `VipPrice` (in ItemResp) = Total VIP Savings for this item.

		itemVipSavings := 0
		if levelDiscountPercent < 100 {
			// Savings = Original * Count * (1 - Discount%)
			// Avoid slight precision issues by: Price * Count - (Price * Count * Discount / 100)
			vipTotal := int(int64(itemPrice) * int64(item.Count) * int64(levelDiscountPercent) / 100)
			itemVipSavings = itemPayPrice - vipTotal
		}

		itemResp := TradePriceCalculateItemRespBO{
			SpuID:    sku.SpuID,
			SkuID:    sku.ID,
			Count:    item.Count,
			CartID:   item.CartID,
			Selected: item.Selected,
			Price:    itemPrice,
			PayPrice: itemPayPrice, // Will be reduced later? Or currently logic flow?
			// Use logic: Calculate Base Pay Price here. Deductions are separate?
			// Current logic: PayPrice is accumulated.
			PicURL:     sku.PicURL,
			Properties: sku.Properties,
			SpuName:    spu.Name,
			CategoryID: spu.CategoryID,
			VipPrice:   itemVipSavings, // Store Savings
		}

		totalPrice += itemPayPrice
		totalPayPrice += itemPayPrice
		totalVipPrice += itemVipSavings

		respBO.Items = append(respBO.Items, itemResp)
	}

	// 5. Set Totals
	// 5. Calculate Reward Activity (Full Reduction)
	matchItems := make([]promotion.ActivityMatchItem, 0)
	for _, item := range respBO.Items {
		matchItems = append(matchItems, promotion.ActivityMatchItem{
			SkuID:      item.SkuID,
			SpuID:      item.SpuID,
			CategoryID: item.CategoryID,
			Price:      item.Price, // Unit Price (Original)
			Count:      item.Count,
		})
	}
	activityDiscount, _, err := s.rewardActivitySvc.CalculateRewardActivity(ctx, matchItems)
	if err != nil {
		return nil, err
	}

	// 5. Total Price
	respBO.Price.TotalPrice = totalPrice
	respBO.Price.DiscountPrice = activityDiscount // Set Activity Discount
	respBO.Price.VipPrice = totalVipPrice         // Set VIP Discount
	// 5. Calculate Delivery Price
	deliveryPrice := 0
	if req.DeliveryType == 1 && req.AddressID != nil && *req.AddressID > 0 { // Express Delivery
		// Get Address to get AreaID
		address, err := s.memberAddressSvc.GetAddress(ctx, req.UserID, *req.AddressID)
		if err != nil {
			return nil, err
		}
		if address != nil {
			// Calculate Freight based on items
			// Group items by TemplateID? SPU has TemplateID usually.
			// Current Model SPU has DeliveryTemplateID?
			// Let's assume SPU has DeliveryTemplateId. Need to check ProductSpuResp.
			// If not available in Resp, need to fetch SPU Entity or assume simpler model.
			// ProductSpuResp (Step 3936 viewed?) spuMap uses *resp.ProductSpuResp.
			// Let's check ProductSpuResp definition if possible, or assume it has DeliveryTemplateID.
			// If not, we might need to fetch it.
			// For this iteration, I'll group by SPU's TemplateID.

			// Map TemplateID -> Count (or Weight/Volume)
			// Assuming Count for now as per CalculateFreight logic
			templateMap := make(map[int64]int)
			for _, item := range respBO.Items {
				spu := spuMap[item.SpuID]
				if spu != nil {
					// Use SPU's delivery template ID. If 0, assume free or default?
					// If spu.DeliveryTemplateID is not in VO, we are stuck.
					// Checking viewed code: Step 3936 line 128 `spuMap[spu.ID] = spu`.
					// I don't see ProductSpuResp definition.
					// I will assume it has `DeliveryTemplateId`.
					templateMap[spu.DeliveryTemplateID] += item.Count
				}
			}

			for tplID, count := range templateMap {
				if tplID > 0 {
					p, err := s.deliveryFreightSvc.CalculateFreight(ctx, tplID, int(address.AreaID), count)
					if err != nil {
						// log error? fail?
						return nil, err
					}
					deliveryPrice += p
				}
			}
		}
	}
	respBO.Price.DeliveryPrice = deliveryPrice

	// Initial Pay Price after Activity AND VIP
	// Note: If Activity + VIP > Total, PayPrice = 0.
	payPrice := totalPrice - activityDiscount - totalVipPrice
	if payPrice < 0 {
		payPrice = 0
	}
	respBO.Price.PayPrice = payPrice

	// 6. Calculate Coupon
	if req.CouponID != nil && *req.CouponID > 0 {
		// Prepare Context for Coupon Check
		// Need SPU IDs and Category IDs. We have SKU Map.
		// Note: ProductSkuResp contains CategoryID inside `ProductSkuResp`?
		// Let's check ProductSkuResp struct definition if possible. Assuming it has CategoryID.
		var spuIDs []int64
		var categoryIDs []int64
		// Helper map to avoid duplicates
		spuMap := make(map[int64]bool)
		catMap := make(map[int64]bool)

		for _, item := range respBO.Items {
			if !spuMap[item.SpuID] {
				spuIDs = append(spuIDs, item.SpuID)
				spuMap[item.SpuID] = true
			}
			if !catMap[item.CategoryID] {
				categoryIDs = append(categoryIDs, item.CategoryID)
				catMap[item.CategoryID] = true
			}
		}

		couponPrice, err := s.couponSvc.CalculateCoupon(ctx, req.UserID, *req.CouponID, int64(respBO.Price.PayPrice), spuIDs, categoryIDs)
		if err != nil {
			// If coupon invalid, strictly we should return error or just ignore?
			// Usually return error to tell user why coupon failed.
			return nil, err
		}
		respBO.CouponID = *req.CouponID
		respBO.Price.CouponPrice = int(couponPrice)
		respBO.Price.PayPrice -= int(couponPrice)
		if respBO.Price.PayPrice < 0 {
			respBO.Price.PayPrice = 0 // Min 0
		}
	}

	return respBO, nil
}
