package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
)

type CouponService struct {
	q *query.Query
}

func NewCouponService() *CouponService {
	return &CouponService{
		q: query.Q,
	}
}

// CreateCouponTemplate 创建优惠券模板 (Admin)
func (s *CouponService) CreateCouponTemplate(ctx context.Context, req *req.CouponTemplateCreateReq) (int64, error) {
	t := &promotion.PromotionCouponTemplate{
		Name:               req.Name,
		Status:             req.Status,
		TotalCount:         req.TotalCount,
		TakeLimitCount:     req.TakeLimitCount,
		TakeType:           req.TakeType,
		UsePriceMin:        req.UsePriceMin,
		ProductScope:       req.ProductScope,
		ProductScopeValues: req.ProductScopeValues,
		ValidityType:       req.ValidityType,
		ValidStartTime:     req.ValidStartTime,
		ValidEndTime:       req.ValidEndTime,
		FixedStartTerm:     req.FixedStartTerm,
		FixedEndTerm:       req.FixedEndTerm,
		DiscountType:       req.DiscountType,
		DiscountPrice:      req.DiscountPrice,
		DiscountPercent:    req.DiscountPercent,
		DiscountLimit:      req.DiscountLimit,
	}
	err := s.q.PromotionCouponTemplate.WithContext(ctx).Create(t)
	return t.ID, err
}

// UpdateCouponTemplate 更新优惠券模板 (Admin)
func (s *CouponService) UpdateCouponTemplate(ctx context.Context, req *req.CouponTemplateUpdateReq) error {
	_, err := s.q.PromotionCouponTemplate.WithContext(ctx).Where(s.q.PromotionCouponTemplate.ID.Eq(req.ID)).Updates(promotion.PromotionCouponTemplate{
		Name:               req.Name,
		Status:             req.Status,
		TotalCount:         req.TotalCount,
		TakeLimitCount:     req.TakeLimitCount,
		TakeType:           req.TakeType,
		UsePriceMin:        req.UsePriceMin,
		ProductScope:       req.ProductScope,
		ProductScopeValues: req.ProductScopeValues,
		ValidityType:       req.ValidityType,
		ValidStartTime:     req.ValidStartTime,
		ValidEndTime:       req.ValidEndTime,
		FixedStartTerm:     req.FixedStartTerm,
		FixedEndTerm:       req.FixedEndTerm,
		DiscountType:       req.DiscountType,
		DiscountPrice:      req.DiscountPrice,
		DiscountPercent:    req.DiscountPercent,
		DiscountLimit:      req.DiscountLimit,
	})
	return err
}

// GetCouponTemplatePage 获得优惠券模板分页 (Admin)
func (s *CouponService) GetCouponTemplatePage(ctx context.Context, req *req.CouponTemplatePageReq) (*core.PageResult[promotion.PromotionCouponTemplate], error) {
	q := s.q.PromotionCouponTemplate.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(s.q.PromotionCouponTemplate.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		q = q.Where(s.q.PromotionCouponTemplate.Status.Eq(*req.Status))
	}

	result, count, err := q.FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	list := make([]promotion.PromotionCouponTemplate, len(result))
	for i, v := range result {
		list[i] = *v
	}

	return &core.PageResult[promotion.PromotionCouponTemplate]{
		List:  list,
		Total: count,
	}, nil
}

// GetCouponPage 获得优惠券分页 (Admin)
func (s *CouponService) GetCouponPage(ctx context.Context, req *req.CouponPageReq) (*core.PageResult[promotion.PromotionCoupon], error) {
	q := s.q.PromotionCoupon.WithContext(ctx)
	if req.UserID != nil {
		q = q.Where(s.q.PromotionCoupon.UserID.Eq(*req.UserID))
	}
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
