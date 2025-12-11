package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CouponUserService struct {
	q *query.Query
}

func NewCouponUserService() *CouponUserService {
	return &CouponUserService{
		q: query.Q,
	}
}

// TakeCoupon 用户领取优惠券
func (s *CouponUserService) TakeCoupon(ctx context.Context, userId int64, req *req.AppCouponTakeReq) (int64, error) {
	// 1. Check Template
	template, err := s.q.PromotionCouponTemplate.WithContext(ctx).Where(s.q.PromotionCouponTemplate.ID.Eq(req.TemplateID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("优惠券模板不存在")
		}
		return 0, err
	}

	if template.Status != 1 { // 1: Enable
		return 0, errors.New("优惠券模板已禁用")
	}
	if template.TotalCount > 0 && template.TakeCount >= template.TotalCount {
		return 0, errors.New("优惠券已领完")
	}

	// 2. Check Limit
	if template.TakeLimitCount > 0 {
		count, err := s.q.PromotionCoupon.WithContext(ctx).Where(s.q.PromotionCoupon.TemplateID.Eq(template.ID), s.q.PromotionCoupon.UserID.Eq(userId)).Count()
		if err != nil {
			return 0, err
		}
		if int(count) >= template.TakeLimitCount {
			return 0, errors.New("超出领取限制")
		}
	}

	// 3. Calculate Validity
	var startTime, endTime time.Time
	now := time.Now()
	if template.ValidityType == 1 { // Fixed Date
		if template.ValidEndTime == nil || now.After(*template.ValidEndTime) {
			return 0, errors.New("优惠券已过期")
		}
		startTime = *template.ValidStartTime
		endTime = *template.ValidEndTime
	} else if template.ValidityType == 2 { // Term
		startTime = now.AddDate(0, 0, template.FixedStartTerm)
		endTime = startTime.AddDate(0, 0, template.FixedEndTerm)
	}

	// 4. Create Coupon
	coupon := &promotion.PromotionCoupon{
		TemplateID:      template.ID,
		Name:            template.Name,
		Status:          1, // Unused
		UserID:          userId,
		ValidStartTime:  startTime,
		ValidEndTime:    endTime,
		DiscountType:    template.DiscountType,
		DiscountPrice:   template.DiscountPrice,
		DiscountPercent: template.DiscountPercent,
		DiscountLimit:   template.DiscountLimit,
		UsePriceMin:     template.UsePriceMin,
	}

	err = s.q.Transaction(func(tx *query.Query) error {
		// Increment Template Take Count
		if _, err := tx.PromotionCouponTemplate.WithContext(ctx).Where(tx.PromotionCouponTemplate.ID.Eq(template.ID)).UpdateSimple(tx.PromotionCouponTemplate.TakeCount.Add(1)); err != nil {
			return err
		}
		// Save Coupon
		if err := tx.PromotionCoupon.WithContext(ctx).Create(coupon); err != nil {
			return err
		}
		return nil
	})

	return coupon.ID, err
}

// GetCouponPage 用户优惠券分页
func (s *CouponUserService) GetCouponPage(ctx context.Context, userId int64, req *req.AppCouponPageReq) (*core.PageResult[promotion.PromotionCoupon], error) {
	q := s.q.PromotionCoupon.WithContext(ctx).Where(s.q.PromotionCoupon.UserID.Eq(userId))
	if req.Status != nil {
		q = q.Where(s.q.PromotionCoupon.Status.Eq(*req.Status))
	}

	result, count, err := q.FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	list := make([]promotion.PromotionCoupon, len(result))
	for i, v := range result {
		list[i] = *v
	}

	return &core.PageResult[promotion.PromotionCoupon]{
		List:  list,
		Total: count,
	}, nil
}

// CalculateCoupon 计算优惠券金额
func (s *CouponUserService) CalculateCoupon(ctx context.Context, userId int64, couponId int64, price int64, spuIDs []int64, categoryIDs []int64) (int64, error) {
	coupon, err := s.q.PromotionCoupon.WithContext(ctx).Where(s.q.PromotionCoupon.ID.Eq(couponId), s.q.PromotionCoupon.UserID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("优惠券不存在")
		}
		return 0, err
	}

	// Check Validity
	if coupon.Status != 1 {
		return 0, errors.New("优惠券不可用")
	}
	now := time.Now()
	if now.Before(coupon.ValidStartTime) || now.After(coupon.ValidEndTime) {
		return 0, errors.New("优惠券不在有效期内")
	}

	// Retrieve Template for Product Scope
	template, err := s.q.PromotionCouponTemplate.WithContext(ctx).Where(s.q.PromotionCouponTemplate.ID.Eq(coupon.TemplateID)).First()
	if err != nil {
		return 0, err
	}

	// Check Min Price
	if price < int64(template.UsePriceMin) {
		return 0, errors.New("未满足使用金额")
	}

	// Check Scope
	if !s.checkScope(template, spuIDs, categoryIDs) {
		return 0, errors.New("不满足商品适用范围")
	}

	// Calculate Amount
	discount := int64(0)
	if coupon.DiscountType == 1 { // Price
		discount = int64(coupon.DiscountPrice)
	} else if coupon.DiscountType == 2 { // Percent
		discount = price * int64(coupon.DiscountPercent) / 100
		if coupon.DiscountLimit > 0 && discount > int64(coupon.DiscountLimit) {
			discount = int64(coupon.DiscountLimit)
		}
	}

	if discount > price {
		discount = price
	}

	return discount, nil
}

// UseCoupon 核销优惠券
func (s *CouponUserService) UseCoupon(ctx context.Context, userId int64, couponId int64, orderId int64) error {
	coupon, err := s.q.PromotionCoupon.WithContext(ctx).Where(s.q.PromotionCoupon.ID.Eq(couponId), s.q.PromotionCoupon.UserID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("优惠券不存在")
		}
		return err
	}

	if coupon.Status != 1 {
		return errors.New("优惠券不可用")
	}
	now := time.Now()
	// Strict validity check again? Maybe lenient if already checked in Price Calc?
	// Consistent behavior: check again.
	if now.Before(coupon.ValidStartTime) || now.After(coupon.ValidEndTime) {
		return errors.New("优惠券不在有效期内")
	}

	// Update Status
	// Use map for updates to include UsedTime
	_, err = s.q.PromotionCoupon.WithContext(ctx).Where(s.q.PromotionCoupon.ID.Eq(couponId)).Updates(map[string]interface{}{
		"status":       2, // Used
		"use_order_id": orderId,
		"use_time":     now,
	})
	return err
}

// checkScope 检查适用范围
func (s *CouponUserService) checkScope(t *promotion.PromotionCouponTemplate, spuIDs []int64, categoryIDs []int64) bool {
	if t.ProductScope == 1 { // All
		return true
	}
	if t.ProductScope == 2 { // Category
		if len(t.ProductScopeValues) == 0 {
			return true
		}
		for _, v := range t.ProductScopeValues {
			for _, cid := range categoryIDs {
				if v == cid {
					return true
				}
			}
		}
		return false
	}
	if t.ProductScope == 3 { // SPU
		if len(t.ProductScopeValues) == 0 {
			return true
		}
		for _, v := range t.ProductScopeValues {
			for _, sid := range spuIDs {
				if v == sid {
					return true
				}
			}
		}
		return false
	}
	return true
}

// GetCouponMatchList 获取匹配的优惠券
func (s *CouponUserService) GetCouponMatchList(ctx context.Context, userId int64, price int64, spuIDs []int64, categoryIDs []int64) ([]promotion.PromotionCoupon, error) {
	// Fetch all valid coupons
	now := time.Now()
	coupons, err := s.q.PromotionCoupon.WithContext(ctx).
		Where(s.q.PromotionCoupon.UserID.Eq(userId)).
		Where(s.q.PromotionCoupon.Status.Eq(1)). // Unused
		Where(s.q.PromotionCoupon.ValidStartTime.Lt(now)).
		Where(s.q.PromotionCoupon.ValidEndTime.Gt(now)).
		Find()
	if err != nil {
		return nil, err
	}

	var matched []promotion.PromotionCoupon
	templateIDs := make([]int64, 0)
	for _, c := range coupons {
		templateIDs = append(templateIDs, c.TemplateID)
	}
	if len(templateIDs) == 0 {
		return matched, nil
	}

	templates, err := s.q.PromotionCouponTemplate.WithContext(ctx).Where(s.q.PromotionCouponTemplate.ID.In(templateIDs...)).Find()
	if err != nil {
		return nil, err
	}
	tmplMap := make(map[int64]*promotion.PromotionCouponTemplate)
	for _, t := range templates {
		tmplMap[t.ID] = t
	}

	for _, c := range coupons {
		t, ok := tmplMap[c.TemplateID]
		if !ok {
			continue
		}
		if price < int64(t.UsePriceMin) {
			continue
		}
		if !s.checkScope(t, spuIDs, categoryIDs) {
			continue
		}
		matched = append(matched, *c)
	}

	// ... Return matched coupons ...
	return matched, nil
}

// ReturnCoupon 退还优惠券
func (s *CouponUserService) ReturnCoupon(ctx context.Context, userId int64, couponId int64) error {
	coupon, err := s.q.PromotionCoupon.WithContext(ctx).Where(s.q.PromotionCoupon.ID.Eq(couponId), s.q.PromotionCoupon.UserID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("优惠券不存在")
		}
		return err
	}

	if coupon.Status != 2 { // 2: Used
		return errors.New("优惠券状态错误，无法退还")
	}

	// Update Status to Unused (1) and clear Usage info
	updates := map[string]interface{}{
		"status":       1,
		"use_order_id": nil,
		"use_time":     nil,
	}
	_, err = s.q.PromotionCoupon.WithContext(ctx).Where(s.q.PromotionCoupon.ID.Eq(couponId)).Updates(updates)
	return err
}
