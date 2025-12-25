package trade

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/area"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	tradeRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	pkgContext "github.com/wxlbd/ruoyi-mall-go/pkg/context"

	"gorm.io/gorm"
)

type TradeOrderUpdateService struct {
	q          *query.Query
	skuSvc     *product.ProductSkuService
	cartSvc    *CartService
	priceSvc   *TradePriceService
	addressSvc *member.MemberAddressService
	couponSvc  *promotion.CouponUserService
	logSvc     *TradeOrderLogService
	noDAO      *tradeRepo.TradeNoRedisDAO
	paySvc     *pay.PayOrderService
	configSvc  *TradeConfigService
}

func NewTradeOrderUpdateService(
	q *query.Query,
	skuSvc *product.ProductSkuService,
	cartSvc *CartService,
	priceSvc *TradePriceService,
	addressSvc *member.MemberAddressService,
	couponSvc *promotion.CouponUserService,
	logSvc *TradeOrderLogService,
	noDAO *tradeRepo.TradeNoRedisDAO,
	paySvc *pay.PayOrderService,
	configSvc *TradeConfigService,
) *TradeOrderUpdateService {
	return &TradeOrderUpdateService{
		q:          q,
		skuSvc:     skuSvc,
		cartSvc:    cartSvc,
		priceSvc:   priceSvc,
		addressSvc: addressSvc,
		couponSvc:  couponSvc,
		logSvc:     logSvc,
		noDAO:      noDAO,
		paySvc:     paySvc,
		configSvc:  configSvc,
	}
}

// SettlementOrder 获得订单结算信息
func (s *TradeOrderUpdateService) SettlementOrder(ctx context.Context, uId int64, req *req.AppTradeOrderSettlementReq) (*resp.AppTradeOrderSettlementResp, error) {
	// 1. Calculate Price
	calcReq := &TradePriceCalculateReqBO{
		UserID:        uId,
		CouponID:      req.CouponID,
		PointStatus:   *req.PointStatus,
		DeliveryType:  req.DeliveryType,
		AddressID:     req.AddressID,
		PickUpStoreID: req.PickUpStoreID,
		Items:         make([]TradePriceCalculateItemBO, len(req.Items)),
	}

	// 设置各种促销活动参数（对齐Java版本TradeOrderConvert第219行）
	if req.SeckillActivityID != nil {
		calcReq.SeckillActivityId = *req.SeckillActivityID
	}
	if req.CombinationActivityID != nil {
		calcReq.CombinationActivityId = *req.CombinationActivityID
	}
	if req.CombinationHeadID != nil {
		calcReq.CombinationHeadId = *req.CombinationHeadID
	}
	if req.BargainRecordID != nil {
		calcReq.BargainRecordId = *req.BargainRecordID
	}
	if req.PointActivityID != nil {
		calcReq.PointActivityId = *req.PointActivityID
	}

	for i, item := range req.Items {
		calcReq.Items[i] = TradePriceCalculateItemBO{
			SkuID:    item.SkuID,
			Count:    item.Count,
			CartID:   item.CartID,
			Selected: true,
		}
	}

	priceResp, err := s.priceSvc.CalculateOrderPrice(ctx, calcReq)
	if err != nil {
		return nil, err
	}

	// 2. Fetch Address (if delivery) - 对齐 Java TradeOrderConvert 第 244-249 行
	var address *resp.AppTradeOrderSettlementAddress
	if req.AddressID != nil {
		addr, err := s.addressSvc.GetAddress(ctx, uId, *req.AddressID)
		if err == nil && addr != nil {
			address = &resp.AppTradeOrderSettlementAddress{
				ID:            addr.ID,
				Name:          addr.Name,
				Mobile:        addr.Mobile,
				AreaID:        int64(addr.AreaID),
				DetailAddress: addr.DetailAddress,
				DefaultStatus: addr.DefaultStatus,
				AreaName:      area.Format(int(addr.AreaID)), // 使用 area.Format() 获取地区名称
			}
		}
	}

	// 3. Assemble Response
	r := &resp.AppTradeOrderSettlementResp{
		Type:       priceResp.Type,
		Items:      make([]resp.AppTradeOrderSettlementItem, len(priceResp.Items)),
		Price:      resp.AppTradeOrderSettlementPrice(priceResp.Price),
		Address:    address,
		UsePoint:   priceResp.UsePoint,
		TotalPoint: priceResp.TotalPoint,
	}

	for i, item := range priceResp.Items {
		r.Items[i] = resp.AppTradeOrderSettlementItem{
			SpuID:      item.SpuID,
			SkuID:      item.SkuID,
			Count:      item.Count,
			CartID:     item.CartID,
			Price:      item.Price,
			PicURL:     item.PicURL,
			Properties: item.Properties,
			SpuName:    item.SpuName, // 补充 SpuName 字段（来自价格计算结果）
		}
	}

	return r, nil
}

// CreateOrder 创建交易订单
func (s *TradeOrderUpdateService) CreateOrder(ctx context.Context, uId int64, reqVO *req.AppTradeOrderCreateReq) (*trade.TradeOrder, error) {
	// 1. Price Calculation
	calcReq := &TradePriceCalculateReqBO{
		UserID:            uId,
		CouponID:          reqVO.CouponID,
		PointStatus:       *reqVO.PointStatus,
		DeliveryType:      reqVO.DeliveryType,
		AddressID:         reqVO.AddressID,
		PickUpStoreID:     reqVO.PickUpStoreID,
		SeckillActivityId: 0, // 默认值，表示非秒杀订单
		Items:             make([]TradePriceCalculateItemBO, len(reqVO.Items)),
	}

	// 设置各种促销活动参数（对齐Java版本TradeOrderConvert第219行）
	if reqVO.SeckillActivityID != nil {
		calcReq.SeckillActivityId = *reqVO.SeckillActivityID
	}
	if reqVO.CombinationActivityID != nil {
		calcReq.CombinationActivityId = *reqVO.CombinationActivityID
	}
	if reqVO.CombinationHeadID != nil {
		calcReq.CombinationHeadId = *reqVO.CombinationHeadID
	}
	if reqVO.BargainRecordID != nil {
		calcReq.BargainRecordId = *reqVO.BargainRecordID
	}
	if reqVO.PointActivityID != nil {
		calcReq.PointActivityId = *reqVO.PointActivityID
	}

	for i, item := range reqVO.Items {
		calcReq.Items[i] = TradePriceCalculateItemBO{
			SkuID:    item.SkuID,
			Count:    item.Count,
			CartID:   item.CartID,
			Selected: true,
		}
	}
	priceResp, err := s.priceSvc.CalculateOrderPrice(ctx, calcReq)
	if err != nil {
		return nil, err
	}

	// 生成订单号 (使用 Redis DAO 确保全局唯一，在 Transaction 前生成以减少持有时间)
	orderNo, err := s.noDAO.GenerateOrderNo(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate order no failed: %w", err)
	}

	// 2. Transaction - 对齐 Java TradeOrderUpdateServiceImpl 第 169-204 行
	var order *trade.TradeOrder
	err = s.q.Transaction(func(tx *query.Query) error {
		// 2.1 Create Order (对齐 Java 实现)
		// 初始状态：0 = UNPAID (待支付)
		// 支付后：10 = UNDELIVERED (待发货)
		// 发货后：20 = DELIVERED (已发货)
		// 完成：30 = COMPLETED (已完成)
		// 取消：40 = CANCELED (已取消)
		order = &trade.TradeOrder{
			No:             orderNo,
			Type:           trade.OrderTypeNormal, // 普通订单
			Terminal:       getTerminal(ctx),      // 从请求头获取终端信息
			UserID:         uId,
			UserIP:         getClientIP(ctx),             // 从请求获取用户 IP
			Status:         trade.TradeOrderStatusUnpaid, // 待支付
			ProductCount:   len(reqVO.Items),
			Remark:         reqVO.Remark,
			PayStatus:      false,
			TotalPrice:     priceResp.Price.TotalPrice,
			DiscountPrice:  priceResp.Price.DiscountPrice,
			PayPrice:       priceResp.Price.PayPrice,
			CouponID:       priceResp.CouponID,
			CouponPrice:    priceResp.Price.CouponPrice,
			DeliveryType:   reqVO.DeliveryType,
			ReceiverName:   reqVO.ReceiverName,
			ReceiverMobile: reqVO.ReceiverMobile,
		}

		// 如果提供了地址 ID，则从地址服务获取详细信息
		if reqVO.AddressID != nil {
			addr, err := s.addressSvc.GetAddress(ctx, uId, *reqVO.AddressID)
			if err == nil && addr != nil {
				order.ReceiverName = addr.Name
				order.ReceiverMobile = addr.Mobile
				order.ReceiverAreaID = int(addr.AreaID)
				order.ReceiverDetailAddress = addr.DetailAddress
			}
		}

		// 如果是自提订单，设置自提门店 ID 和核销码
		if reqVO.DeliveryType == trade.DeliveryTypePickUp {
			if reqVO.PickUpStoreID != nil {
				order.PickUpStoreID = *reqVO.PickUpStoreID
			}
			// 生成随机核销码（对齐 Java RandomUtil.randomNumbers(8)）
			order.PickUpVerifyCode = generateRandomNumbers(trade.PickUpVerifyCodeLength)
		}

		if err := tx.TradeOrder.WithContext(ctx).Create(order); err != nil {
			return err
		}

		// 2.2 Create Order Items
		items := make([]*trade.TradeOrderItem, len(priceResp.Items))
		for i, item := range priceResp.Items {
			items[i] = &trade.TradeOrderItem{
				UserID:      uId,
				OrderID:     order.ID,
				SpuID:       item.SpuID,
				SpuName:     item.SpuName, // 从价格计算结果获取 SpuName
				SkuID:       item.SkuID,
				Count:       item.Count,
				Price:       item.Price,
				PayPrice:    item.PayPrice,
				PicURL:      item.PicURL,
				CouponPrice: item.CouponPrice,
				// Properties: item.Properties (need serialize),
			}
		}
		if err := tx.TradeOrderItem.WithContext(ctx).Create(items...); err != nil {
			return err
		}

		// 2.3 Clear Cart (if cart items)
		var cartIds []int64
		for _, item := range calcReq.Items {
			if item.CartID > 0 {
				cartIds = append(cartIds, item.CartID)
			}
		}
		if len(cartIds) > 0 {
			if err := s.cartSvc.DeleteCart(ctx, uId, cartIds); err != nil {
				return err
			}
		}

		// 2.4 Decrease Stock
		var stockItems []req.ProductSkuUpdateStockItemReq
		for _, item := range priceResp.Items {
			stockItems = append(stockItems, req.ProductSkuUpdateStockItemReq{
				ID:        item.SkuID,
				IncrCount: -item.Count,
			})
		}
		if err := s.skuSvc.UpdateSkuStock(ctx, &req.ProductSkuUpdateStockReq{Items: stockItems}); err != nil {
			return err
		}

		// 2.5 Use Coupon
		if priceResp.CouponID > 0 {
			if err := s.couponSvc.UseCoupon(ctx, uId, priceResp.CouponID, order.ID); err != nil {
				return err
			}
		}

		// 2.6 Log
		if err := s.createOrderLog(ctx, order, "Create Order", trade.OrderOperateTypeCreate); err != nil {
			return err
		}

		// 2.7 Create Pay Order (对齐 Java)
		if order.PayPrice > 0 {
			if err := s.createPayOrder(ctx, tx, order, priceResp.Items); err != nil {
				return err
			}
		}

		return nil
	})
	return order, err
}

func (s *TradeOrderUpdateService) createPayOrder(ctx context.Context, tx *query.Query, order *trade.TradeOrder, items []TradePriceCalculateItemRespBO) error {
	// 1. Get Config
	conf, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		return err
	}

	// 2. Build Subject
	subject := ""
	if len(items) > 0 {
		subject = items[0].SpuName
		if len(items) > 1 {
			subject += fmt.Sprintf("等 %d 件商品", len(items))
		}
	}
	// 对齐 Java: 限制长度最大为 32
	if len([]rune(subject)) > 32 {
		subject = string([]rune(subject)[:32])
	}

	// 3. Create Pay Order
	payTimeout := conf.PayTimeoutMinutes
	if payTimeout <= 0 {
		payTimeout = 120 // 兜底：默认 120 分钟
	}
	expireTime := time.Now().Add(time.Duration(payTimeout) * time.Minute)
	payOrderID, err := s.paySvc.CreateOrder(ctx, &req.PayOrderCreateReq{
		AppKey:          "mall", // 对齐 Java: 默认使用 "mall"
		UserIP:          order.UserIP,
		MerchantOrderId: strconv.FormatInt(order.ID, 10),
		Subject:         subject,
		Body:            subject,
		Price:           order.PayPrice,
		ExpireTime:      expireTime,
	})
	if err != nil {
		return err
	}

	// 4. Update Trade Order
	_, err = tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(order.ID)).Update(tx.TradeOrder.PayOrderID, payOrderID)
	if err != nil {
		return err
	}
	order.PayOrderID = &payOrderID

	return nil
}

// DeliveryOrder 订单发货
func (s *TradeOrderUpdateService) DeliveryOrder(ctx context.Context, reqVO *req.TradeOrderDeliveryReq) error {
	// 使用 GORM 事务包装（对齐 Java @Transactional）
	return s.q.Transaction(func(tx *query.Query) error {
		// 1. Check Order Exists
		order, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(reqVO.ID)).First()
		if err != nil {
			return err
		}
		if order.Status != trade.TradeOrderStatusUndelivered { // 待发货
			// return fmt.Errorf("order status error")
		}

		now := time.Now()
		// 2. Update Order (in transaction)
		_, err = tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(reqVO.ID)).Updates(trade.TradeOrder{
			Status:       trade.TradeOrderStatusDelivered, // 已发货
			LogisticsID:  reqVO.LogisticsID,
			LogisticsNo:  reqVO.LogisticsNo,
			DeliveryTime: &now,
		})
		if err != nil {
			return err
		}

		// 3. Log (in transaction)
		logOrder := *order
		logOrder.Status = trade.TradeOrderStatusDelivered
		return s.createOrderLog(ctx, &logOrder, "Order Delivered", trade.OrderOperateTypeDelivery)
	})
}

// UpdateOrderPaid 更新订单为已支付
func (s *TradeOrderUpdateService) UpdateOrderPaid(ctx context.Context, id int64, payOrderId int64) error {
	// 1. Get Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if order.Status != trade.TradeOrderStatusUnpaid { // 待支付
		return fmt.Errorf("order status is not unpaid")
	}
	if order.PayStatus {
		return fmt.Errorf("order is already paid")
	}

	// 2. Update
	now := time.Now()
	err = s.q.Transaction(func(tx *query.Query) error {
		// Update Order
		updateMap := map[string]interface{}{
			"status":       trade.TradeOrderStatusUndelivered, // 待发货
			"pay_status":   true,
			"pay_time":     &now,
			"pay_order_id": payOrderId,
		}
		if _, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(id)).Updates(updateMap); err != nil {
			return err
		}

		// Log
		logOrder := *order
		logOrder.Status = trade.TradeOrderStatusUndelivered
		if err := s.createOrderLog(ctx, &logOrder, "Order Paid", trade.OrderOperateTypePay); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *TradeOrderUpdateService) createOrderLog(ctx context.Context, order *trade.TradeOrder, content string, operateType int) error {
	uid := int64(0)

	log := &trade.TradeOrderLog{
		UserID:       uid,
		UserType:     model.UserTypeMember, // 会员用户
		OrderID:      order.ID,
		BeforeStatus: 0, // Simplified, ideally pass old status
		AfterStatus:  order.Status,
		OperateType:  operateType,
		Content:      content,
	}
	return s.q.TradeOrderLog.WithContext(ctx).Create(log)
}

// CancelOrder 取消交易订单
func (s *TradeOrderUpdateService) CancelOrder(ctx context.Context, uId int64, id int64) error {
	// 1. Check Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id), s.q.TradeOrder.UserID.Eq(uId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	if order.Status != trade.TradeOrderStatusUnpaid { // 待支付，只允许未支付订单取消
		// For now, restrict to Unpaid.
		return errors.New("订单状态不允许取消")
	}

	// 2. Transaction
	err = s.q.Transaction(func(tx *query.Query) error {
		// 2.1 Update Order Status
		now := time.Now()
		if _, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(id)).Updates(trade.TradeOrder{
			Status:     trade.TradeOrderStatusCanceled, // 已取消
			CancelTime: &now,
			CancelType: trade.OrderCancelTypeUser, // 用户取消
		}); err != nil {
			return err
		}

		// 2.2 Release Stock
		items, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.OrderID.Eq(id)).Find()
		if err != nil {
			return err
		}
		var stockItems []req.ProductSkuUpdateStockItemReq
		for _, item := range items {
			stockItems = append(stockItems, req.ProductSkuUpdateStockItemReq{
				ID:        item.SkuID,
				IncrCount: item.Count, // Positive to restore stock
			})
		}
		if err := s.skuSvc.UpdateSkuStock(ctx, &req.ProductSkuUpdateStockReq{Items: stockItems}); err != nil {
			return err
		}

		// 2.3 Refund Coupon
		if order.CouponID > 0 {
			if err := s.couponSvc.ReturnCoupon(ctx, uId, order.CouponID); err != nil {
				return err
			}
		}

		// 2.4 Log
		logOrder := *order
		logOrder.Status = trade.TradeOrderStatusCanceled
		if err := s.createOrderLog(ctx, &logOrder, "User Cancelled Order", trade.OrderOperateTypeCancel); err != nil {
			return err
		}

		return nil
	})

	return err
}

// DeleteOrder 删除订单
func (s *TradeOrderUpdateService) DeleteOrder(ctx context.Context, uId int64, id int64) error {
	// 1. Check Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id), s.q.TradeOrder.UserID.Eq(uId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	// Java: Check status (Cancelled or Completed can be deleted?)
	// Usually only Cancelled or Completed.
	if order.Status != trade.TradeOrderStatusCanceled && order.Status != trade.TradeOrderStatusCompleted {
		return errors.New("只有取消或完成的订单可以删除")
	}

	// 2. Delete (Soft Delete)
	_, err = s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).Delete()
	return err
}

// CancelPaidOrder 取消已支付订单 (对齐 Java TradeOrderApi.cancelPaidOrder)
func (s *TradeOrderUpdateService) CancelPaidOrder(ctx context.Context, uId int64, id int64, cancelType int) error {
	// 1. Check Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id), s.q.TradeOrder.UserID.Eq(uId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	if !order.PayStatus {
		return errors.New("订单未支付，请使用 CancelOrder")
	}
	// 只有待发货状态允许取消（对应 Java 逻辑：已支付但未发货）
	if order.Status != trade.TradeOrderStatusUndelivered {
		return errors.New("订单状态不允许取消")
	}

	// 2. Transaction
	err = s.q.Transaction(func(tx *query.Query) error {
		// 2.1 Update Order Status
		now := time.Now()
		if _, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(id)).Updates(trade.TradeOrder{
			Status:       trade.TradeOrderStatusCanceled, // 已取消
			CancelTime:   &now,
			CancelType:   cancelType,                   // 系统或拼团关闭取消
			RefundStatus: trade.OrderRefundStatusApply, // 标记为申请退款 (待人工或自动退款)
		}); err != nil {
			return err
		}

		// 2.2 Release Stock
		items, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.OrderID.Eq(id)).Find()
		if err != nil {
			return err
		}
		var stockItems []req.ProductSkuUpdateStockItemReq
		for _, item := range items {
			stockItems = append(stockItems, req.ProductSkuUpdateStockItemReq{
				ID:        item.SkuID,
				IncrCount: item.Count, // Positive to restore stock
			})
		}
		if err := s.skuSvc.UpdateSkuStock(ctx, &req.ProductSkuUpdateStockReq{Items: stockItems}); err != nil {
			return err
		}

		// 2.3 Refund Coupon
		if order.CouponID > 0 {
			if err := s.couponSvc.ReturnCoupon(ctx, uId, order.CouponID); err != nil {
				return err
			}
		}

		// 2.4 Log
		logOrder := *order
		logOrder.Status = trade.TradeOrderStatusCanceled
		if err := s.createOrderLog(ctx, &logOrder, "Paid Order Cancelled", trade.OrderOperateTypeCancel); err != nil {
			return err
		}

		return nil
	})

	return err
}

// UpdateOrderRemark 订单备注
func (s *TradeOrderUpdateService) UpdateOrderRemark(ctx context.Context, req *req.TradeOrderRemarkReq) error {
	// 使用 GORM 事务包装（对齐 Java @Transactional）
	return s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(req.ID)).Update(tx.TradeOrder.Remark, req.Remark)
		return err
	})
}

// UpdateOrderPrice 订单调价
func (s *TradeOrderUpdateService) UpdateOrderPrice(ctx context.Context, req *req.TradeOrderUpdatePriceReq) error {
	// 使用 GORM 事务包装（对齐 Java @Transactional）
	return s.q.Transaction(func(tx *query.Query) error {
		order, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(req.ID)).First()
		if err != nil {
			return err
		}
		if order.PayStatus {
			return errors.New("已支付订单不允许改价")
		}

		// New Price Calculation
		newPayPrice := order.PayPrice + req.AdjustPrice
		if newPayPrice < 0 {
			return errors.New("调价后金额不能小于 0")
		}

		_, err = tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(req.ID)).Updates(map[string]interface{}{
			"adjust_price": order.AdjustPrice + req.AdjustPrice,
			"pay_price":    newPayPrice,
		})
		return err
	})
}

// UpdateOrderAddress 修改订单收货地址
func (s *TradeOrderUpdateService) UpdateOrderAddress(ctx context.Context, req *req.TradeOrderUpdateAddressReq) error {
	// 使用 GORM 事务包装（对齐 Java @Transactional）
	return s.q.Transaction(func(tx *query.Query) error {
		// Check status (only undelivered?)
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(req.ID)).Updates(map[string]interface{}{
			"receiver_name":           req.ReceiverName,
			"receiver_mobile":         req.ReceiverMobile,
			"receiver_area_id":        req.ReceiverAreaID,
			"receiver_detail_address": req.ReceiverDetailAddress,
		})
		return err
	})
}

// PickUpOrderByAdmin 核销订单 (By ID)
func (s *TradeOrderUpdateService) PickUpOrderByAdmin(ctx context.Context, adminUserId int64, id int64) error {
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	return s.pickUpOrder(ctx, order)
}

// PickUpOrderByVerifyCode 核销订单 (By Code)
func (s *TradeOrderUpdateService) PickUpOrderByVerifyCode(ctx context.Context, adminUserId int64, verifyCode string) error {
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.PickUpVerifyCode.Eq(verifyCode)).First()
	if err != nil {
		return errors.New("核销码无效")
	}
	return s.pickUpOrder(ctx, order)
}

func (s *TradeOrderUpdateService) pickUpOrder(ctx context.Context, order *trade.TradeOrder) error {
	if order.DeliveryType != trade.DeliveryTypePickUp {
		return errors.New("非自提订单")
	}
	if order.Status != trade.TradeOrderStatusUndelivered {
		return errors.New("订单状态不正确")
	}

	now := time.Now()
	err := s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(order.ID)).Updates(trade.TradeOrder{
			Status:      trade.TradeOrderStatusCompleted, // 已完成
			ReceiveTime: &now,
		})
		if err != nil {
			return err
		}
		// Log
		return s.createOrderLog(ctx, order, "Admin PickUp", trade.OrderOperateTypePickUp) // 自提核销
	})
	return err
}

// GetByPickUpVerifyCode 查询核销码对应的订单
func (s *TradeOrderUpdateService) GetByPickUpVerifyCode(ctx context.Context, verifyCode string) (*trade.TradeOrder, error) {
	return s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.PickUpVerifyCode.Eq(verifyCode)).First()
}

// CreateOrderItemCommentByMember 创建订单项评价
func (s *TradeOrderUpdateService) CreateOrderItemCommentByMember(ctx context.Context, uId int64, req *req.AppTradeOrderItemCommentCreateReq) (int64, error) {
	// 1. Get Order Item
	item, err := s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.ID.Eq(req.OrderItemID), s.q.TradeOrderItem.UserID.Eq(uId)).First()
	if err != nil {
		return 0, err
	}
	if item.CommentStatus {
		return 0, errors.New("该商品已评价")
	}

	// 3. Update Order Item Comment Status
	_, err = s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.ID.Eq(item.ID)).Update(s.q.TradeOrderItem.CommentStatus, true)
	if err != nil {
		return 0, err
	}

	// 4. Update Order Comment Status if all items commented
	// Count uncommented items
	count, _ := s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.OrderID.Eq(item.OrderID), s.q.TradeOrderItem.CommentStatus.Eq(model.NewBitBool(false))).Count()
	if count == 0 {
		_, _ = s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(item.OrderID)).Update(s.q.TradeOrder.CommentStatus, true)
	}

	return 0, nil // Return Comment ID
}

// ReceiveOrder 用户确认收货
func (s *TradeOrderUpdateService) ReceiveOrder(ctx context.Context, uId int64, orderId int64) error {
	// 1. Get Order
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(orderId), s.q.TradeOrder.UserID.Eq(uId)).First()
	if err != nil {
		return errors.New("订单不存在")
	}

	// 2. Validate Status - 只有已发货状态才能确认收货
	if order.Status != trade.TradeOrderStatusDelivered { // 已发货
		return errors.New("订单状态不正确，无法确认收货")
	}

	// 3. Update Order Status
	now := time.Now()
	err = s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(order.ID)).Updates(trade.TradeOrder{
			Status:      trade.TradeOrderStatusCompleted, // 已完成
			ReceiveTime: &now,
		})
		if err != nil {
			return err
		}
		// Log
		return s.createOrderLog(ctx, order, "用户确认收货", trade.OrderOperateTypeReceive) // 确认收货
	})
	return err
}

// UpdatePaidOrderRefunded 更新支付订单为已退款 (Callback from Pay)
func (s *TradeOrderUpdateService) UpdatePaidOrderRefunded(ctx context.Context, orderId int64, payRefundId int64) error {
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(orderId)).First()
	if err != nil {
		return err
	}

	return s.q.Transaction(func(tx *query.Query) error {
		_, err := tx.TradeOrder.WithContext(ctx).Where(tx.TradeOrder.ID.Eq(orderId)).Updates(map[string]interface{}{
			"refund_status": trade.OrderRefundStatusRefunded, // 已退款
			// "pay_refund_id": payRefundId, // No field
		})
		if err != nil {
			return err
		}
		// Log
		return s.createOrderLog(ctx, order, "Order Refunded (Pay Callback)", trade.OrderOperateTypeRefund)
	})
}

// getTerminal 从上下文获取终端类型 (对齐 Java WebFrameworkUtils.getTerminal())
// 终端类型：0=未知, 10=微信小程序, 11=微信公众号, 20=H5网页, 31=手机App
func getTerminal(ctx context.Context) int {
	// 从 context 中获取 gin.Context（中间件注入）
	ginCtx, ok := ctx.Value(pkgContext.CtxGinContextKey).(*gin.Context)
	if !ok {
		return 0
	}

	// 1. 优先从请求头读取终端类型（Terminal 标准头）
	terminalStr := ginCtx.GetHeader("Terminal")
	if terminalStr == "" {
		// 2. 其次从查询参数读取 terminal
		terminalStr = ginCtx.Query("terminal")
	}

	if terminalStr != "" {
		// 映射常见终端值 (可以根据 project schema 调整)
		switch strings.ToLower(terminalStr) {
		case "10", "wechat_mini_program":
			return 10
		case "11", "wechat_official_account":
			return 11
		case "20", "h5":
			return 20
		case "31", "app":
			return 31
		}
	}
	return 0 // TerminalEnum.UNKNOWN
}

// getClientIP 从上下文获取用户 IP (对齐 Java ServletUtils.getClientIP())
func getClientIP(ctx context.Context) string {
	// 从 context 中获取 gin.Context（中间件注入）
	if ginCtx, ok := ctx.Value(pkgContext.CtxGinContextKey).(*gin.Context); ok {
		// 1. 检查 X-Forwarded-For 头（来自代理）
		if forwardedFor := ginCtx.GetHeader("X-Forwarded-For"); forwardedFor != "" {
			// X-Forwarded-For 可能包含多个IP，取第一个（真实用户IP）
			ips := strings.Split(forwardedFor, ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}

		// 2. 检查 X-Real-IP 头（单级代理）
		if realIP := ginCtx.GetHeader("X-Real-IP"); realIP != "" {
			return strings.TrimSpace(realIP)
		}

		// 3. 使用 Gin 的 ClientIP() 方法（自动处理代理）
		return ginCtx.ClientIP()
	}

	// 非 Gin context 时返回默认值
	return "127.0.0.1"
}

// generateRandomNumbers 生成指定位数的随机数字字符串 (对齐 Java RandomUtil.randomNumbers(n))
func generateRandomNumbers(n int) string {
	const digits = "0123456789"
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = digits[rand.Intn(10)]
	}
	return string(result)
}
