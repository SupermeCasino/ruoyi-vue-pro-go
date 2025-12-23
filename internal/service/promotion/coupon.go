package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type CouponService struct {
	q *query.Query
}

func NewCouponService(q *query.Query) *CouponService {
	return &CouponService{
		q: q,
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
func (s *CouponService) GetCouponTemplatePage(ctx context.Context, req *req.CouponTemplatePageReq) (*pagination.PageResult[promotion.PromotionCouponTemplate], error) {
	q := s.q.PromotionCouponTemplate.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(s.q.PromotionCouponTemplate.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		q = q.Where(s.q.PromotionCouponTemplate.Status.Eq(int(*req.Status)))
	}
	if len(req.CreateTime) == 2 && req.CreateTime[0] != nil && req.CreateTime[1] != nil {
		q = q.Where(s.q.PromotionCouponTemplate.CreateTime.Between(*req.CreateTime[0], *req.CreateTime[1]))
	}

	result, count, err := q.FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	list := make([]promotion.PromotionCouponTemplate, len(result))
	for i, v := range result {
		list[i] = *v
	}

	return &pagination.PageResult[promotion.PromotionCouponTemplate]{
		List:  list,
		Total: count,
	}, nil
}

// UpdateCouponTemplateStatus 更新优惠券模板状态 (Admin)
// 对应 Java: CouponTemplateService.updateCouponTemplateStatus
func (s *CouponService) UpdateCouponTemplateStatus(ctx context.Context, id int64, status int32) error {
	t := s.q.PromotionCouponTemplate
	// 校验存在
	template, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}

	// 更新状态
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Update(t.Status, status)
	return err
}

// DeleteCouponTemplate 删除优惠券模板 (Admin)
// 对应 Java: CouponTemplateService.deleteCouponTemplate
func (s *CouponService) DeleteCouponTemplate(ctx context.Context, id int64) error {
	t := s.q.PromotionCouponTemplate
	// 校验存在
	template, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}

	// 删除
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

// GetCouponTemplate 获取优惠券模板详情 (Admin)
// 对应 Java: CouponTemplateService.getCouponTemplate
func (s *CouponService) GetCouponTemplate(ctx context.Context, id int64) (*promotion.PromotionCouponTemplate, error) {
	t := s.q.PromotionCouponTemplate
	return t.WithContext(ctx).Where(t.ID.Eq(id)).First()
}

// GetCouponTemplateList 获取优惠券模板列表 (Admin)
// 对应 Java: CouponTemplateService.getCouponTemplateList(ids)
func (s *CouponService) GetCouponTemplateList(ctx context.Context, ids []int64) ([]*promotion.PromotionCouponTemplate, error) {
	if len(ids) == 0 {
		return []*promotion.PromotionCouponTemplate{}, nil
	}
	t := s.q.PromotionCouponTemplate
	return t.WithContext(ctx).Where(t.ID.In(ids...)).Find()
}

// GetCouponPage 获得优惠券分页 (Admin)
func (s *CouponService) GetCouponPage(ctx context.Context, req *req.CouponPageReq) (*pagination.PageResult[promotion.PromotionCoupon], error) {
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

	return &pagination.PageResult[promotion.PromotionCoupon]{
		List:  list,
		Total: count,
	}, nil
}

// DeleteCoupon 删除/回收优惠券 (Admin)
// 对应 Java: CouponService.deleteCoupon
func (s *CouponService) DeleteCoupon(ctx context.Context, id int64) error {
	c := s.q.PromotionCoupon
	// 校验存在
	coupon, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1006002000, "优惠券不存在")
	}
	if coupon == nil {
		return errors.NewBizError(1006002000, "优惠券不存在")
	}

	// 删除
	_, err = c.WithContext(ctx).Where(c.ID.Eq(id)).Delete()
	return err
}

// TakeCouponByAdmin 管理员发送优惠券给用户 (Admin)
// 对应 Java: CouponService.takeCouponByAdmin
func (s *CouponService) TakeCouponByAdmin(ctx context.Context, templateId int64, userIds []int64) error {
	if len(userIds) == 0 {
		return nil
	}

	// 1. 获取优惠券模板
	t := s.q.PromotionCouponTemplate
	template, err := t.WithContext(ctx).Where(t.ID.Eq(templateId)).First()
	if err != nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}
	if template == nil {
		return errors.NewBizError(1006001000, "优惠券模板不存在")
	}

	// 2. 计算有效期
	var validStartTime, validEndTime time.Time
	if template.ValidStartTime != nil {
		validStartTime = *template.ValidStartTime
	} else {
		validStartTime = time.Now()
	}
	if template.ValidEndTime != nil {
		validEndTime = *template.ValidEndTime
	} else {
		// 默认30天后过期
		validEndTime = time.Now().AddDate(0, 0, 30)
	}

	// 3. 为每个用户创建优惠券
	coupons := make([]*promotion.PromotionCoupon, 0, len(userIds))
	for _, userId := range userIds {
		coupon := &promotion.PromotionCoupon{
			TemplateID:      templateId,
			Name:            template.Name,
			UserID:          userId,
			Status:          1, // 未使用
			UsePriceMin:     template.UsePriceMin,
			ValidStartTime:  validStartTime,
			ValidEndTime:    validEndTime,
			DiscountType:    template.DiscountType,
			DiscountPrice:   template.DiscountPrice,
			DiscountPercent: template.DiscountPercent,
			DiscountLimit:   template.DiscountLimit,
		}
		coupons = append(coupons, coupon)
	}

	// 4. 批量创建
	c := s.q.PromotionCoupon
	return c.WithContext(ctx).Create(coupons...)
}

// ========== App 端方法 ==========

// GetCouponTemplateForApp 获取单个优惠券模板 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplate
func (s *CouponService) GetCouponTemplateForApp(ctx context.Context, id int64, userId int64) (*resp.AppCouponTemplateResp, error) {
	t := s.q.PromotionCouponTemplate
	template, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, nil // 返回 null 对齐 Java
	}

	// 处理是否可领取 (对齐 Java: couponService.getUserCanCanTakeMap)
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, []int64{template.ID})
	if err != nil {
		return nil, err
	}

	return s.convertToAppCouponTemplateResp(template, canTakeMap[template.ID]), nil
}

// GetCouponTemplateListForApp 获取优惠券模板列表 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplateList (带条件)
func (s *CouponService) GetCouponTemplateListForApp(ctx context.Context, spuId *int64, productScope *int, count int, userId int64) ([]*resp.AppCouponTemplateResp, error) {
	t := s.q.PromotionCouponTemplate
	q := t.WithContext(ctx).Where(t.Status.Eq(1)) // 只查询启用状态
	q = q.Where(t.TakeType.Eq(1))                 // 领取方式 = 直接领取 (CouponTakeTypeEnum.USER)

	if productScope != nil {
		q = q.Where(t.ProductScope.Eq(*productScope))
	}

	if count <= 0 {
		count = 10
	}
	q = q.Limit(count)

	templates, err := q.Find()
	if err != nil {
		return nil, err
	}

	// 获取用户是否可领取
	templateIds := make([]int64, len(templates))
	for i, tmpl := range templates {
		templateIds[i] = tmpl.ID
	}
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, templateIds)
	if err != nil {
		return nil, err
	}

	// 转换响应
	result := make([]*resp.AppCouponTemplateResp, len(templates))
	for i, tmpl := range templates {
		result[i] = s.convertToAppCouponTemplateResp(tmpl, canTakeMap[tmpl.ID])
	}
	return result, nil
}

// GetCouponTemplateListByIdsForApp 按 ID 获取优惠券模板列表 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplateList (按ids)
func (s *CouponService) GetCouponTemplateListByIdsForApp(ctx context.Context, ids []int64, userId int64) ([]*resp.AppCouponTemplateResp, error) {
	if len(ids) == 0 {
		return []*resp.AppCouponTemplateResp{}, nil
	}
	t := s.q.PromotionCouponTemplate
	templates, err := t.WithContext(ctx).Where(t.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	// 获取用户是否可领取
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, ids)
	if err != nil {
		return nil, err
	}

	// 转换响应
	result := make([]*resp.AppCouponTemplateResp, len(templates))
	for i, tmpl := range templates {
		result[i] = s.convertToAppCouponTemplateResp(tmpl, canTakeMap[tmpl.ID])
	}
	return result, nil
}

// GetCouponTemplatePageForApp 获取优惠券模板分页 (App 端)
// 对齐 Java: AppCouponTemplateController.getCouponTemplatePage
func (s *CouponService) GetCouponTemplatePageForApp(ctx context.Context, r *req.AppCouponTemplatePageReq, userId int64) (*pagination.PageResult[*resp.AppCouponTemplateResp], error) {
	t := s.q.PromotionCouponTemplate
	q := t.WithContext(ctx).Where(t.Status.Eq(1)) // 只查询启用状态
	q = q.Where(t.TakeType.Eq(1))                 // 领取方式 = 直接领取

	if r.ProductScope != nil {
		q = q.Where(t.ProductScope.Eq(*r.ProductScope))
	}

	templates, count, err := q.FindByPage((r.PageNo-1)*r.PageSize, r.PageSize)
	if err != nil {
		return nil, err
	}

	// 获取用户是否可领取
	templateIds := make([]int64, len(templates))
	for i, tmpl := range templates {
		templateIds[i] = tmpl.ID
	}
	canTakeMap, err := s.GetUserCanTakeMap(ctx, userId, templateIds)
	if err != nil {
		return nil, err
	}

	// 转换响应
	result := make([]*resp.AppCouponTemplateResp, len(templates))
	for i, tmpl := range templates {
		result[i] = s.convertToAppCouponTemplateResp(tmpl, canTakeMap[tmpl.ID])
	}

	return &pagination.PageResult[*resp.AppCouponTemplateResp]{
		List:  result,
		Total: count,
	}, nil
}

// convertToAppCouponTemplateResp 转换模板 Model 到 App 响应 VO
func (s *CouponService) convertToAppCouponTemplateResp(t *promotion.PromotionCouponTemplate, canTake bool) *resp.AppCouponTemplateResp {
	return &resp.AppCouponTemplateResp{
		ID:                 t.ID,
		Name:               t.Name,
		Description:        "", // 模型中暂无此字段
		TotalCount:         t.TotalCount,
		TakeLimitCount:     t.TakeLimitCount,
		UsePrice:           t.UsePriceMin,
		ProductScope:       int(t.ProductScope),
		ProductScopeValues: t.ProductScopeValues,
		ValidityType:       t.ValidityType,
		ValidStartTime:     t.ValidStartTime,
		ValidEndTime:       t.ValidEndTime,
		FixedStartTerm:     t.FixedStartTerm,
		FixedEndTerm:       t.FixedEndTerm,
		DiscountType:       t.DiscountType,
		DiscountPercent:    t.DiscountPercent,
		DiscountPrice:      t.DiscountPrice,
		DiscountLimitPrice: t.DiscountLimit,
		TakeCount:          t.TakeCount,
		CanTake:            canTake,
	}
}

// GetUserCanTakeMap 获取用户是否可领取某模板的 Map
// 对齐 Java: CouponService.getUserCanCanTakeMap
func (s *CouponService) GetUserCanTakeMap(ctx context.Context, userId int64, templateIds []int64) (map[int64]bool, error) {
	result := make(map[int64]bool)
	if userId == 0 || len(templateIds) == 0 {
		for _, id := range templateIds {
			result[id] = true // 未登录用户默认可领取
		}
		return result, nil
	}

	// 查询用户对每个模板已领取的数量
	c := s.q.PromotionCoupon
	coupons, err := c.WithContext(ctx).
		Where(c.UserID.Eq(userId)).
		Where(c.TemplateID.In(templateIds...)).
		Find()
	if err != nil {
		return nil, err
	}

	// 统计每个模板的领取数量
	takeCountMap := make(map[int64]int)
	for _, coupon := range coupons {
		takeCountMap[coupon.TemplateID]++
	}

	// 查询模板的领取限制
	t := s.q.PromotionCouponTemplate
	templates, err := t.WithContext(ctx).Where(t.ID.In(templateIds...)).Find()
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		takeCount := takeCountMap[template.ID]
		if template.TakeLimitCount <= 0 {
			result[template.ID] = true // 无限制
		} else {
			result[template.ID] = takeCount < template.TakeLimitCount
		}
	}

	return result, nil
}
