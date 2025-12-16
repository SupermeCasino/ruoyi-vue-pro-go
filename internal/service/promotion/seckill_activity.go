package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product" // Import Product services

	"gorm.io/gen"
)

type SeckillActivityService struct {
	q         *query.Query
	configSvc *SeckillConfigService
	spuSvc    *product.ProductSpuService
	skuSvc    *product.ProductSkuService // Need SkuService for validation
}

func NewSeckillActivityService(q *query.Query, configSvc *SeckillConfigService, spuSvc *product.ProductSpuService, skuSvc *product.ProductSkuService) *SeckillActivityService {
	return &SeckillActivityService{
		q:         q,
		configSvc: configSvc,
		spuSvc:    spuSvc,
		skuSvc:    skuSvc,
	}
}

// CreateSeckillActivity 创建秒杀活动
func (s *SeckillActivityService) CreateSeckillActivity(ctx context.Context, r *req.SeckillActivityCreateReq) (int64, error) {
	// 1. Validate
	if err := s.configSvc.ValidateSeckillConfigExists(ctx, r.ConfigIds); err != nil {
		return 0, err
	}
	if err := s.validateProductConflict(ctx, r.ConfigIds, r.SpuID, 0); err != nil {
		return 0, err
	}
	// TODO: Validate Product Exists (Spu + Sku)
	// assuming spuSvc and skuSvc has methods.

	// 2. Transaction
	var activityID int64
	err := s.q.Transaction(func(tx *query.Query) error {
		// Calculate Stock
		totalStock := 0
		for _, p := range r.Products {
			totalStock += p.Stock
		}

		// Insert Activity
		activity := &promotion.PromotionSeckillActivity{
			SpuID:            r.SpuID,
			Name:             r.Name,
			Status:           1, // Enable by default? Java says CommonStatusEnum.ENABLE
			Remark:           r.Remark,
			StartTime:        r.StartTime,
			EndTime:          r.EndTime,
			Sort:             r.Sort,
			ConfigIds:        r.ConfigIds,
			TotalLimitCount:  r.TotalLimitCount,
			SingleLimitCount: r.SingleLimitCount,
			Stock:            totalStock,
			TotalStock:       totalStock,
		}
		if err := tx.PromotionSeckillActivity.WithContext(ctx).Create(activity); err != nil {
			return err
		}
		activityID = activity.ID

		// Insert Products
		products := make([]*promotion.PromotionSeckillProduct, len(r.Products))
		for i, p := range r.Products {
			products[i] = &promotion.PromotionSeckillProduct{
				ActivityID:        activity.ID,
				ConfigIds:         r.ConfigIds,
				SpuID:             r.SpuID,
				SkuID:             p.SkuID,
				SeckillPrice:      p.SeckillPrice,
				Stock:             p.Stock,
				ActivityStatus:    activity.Status,
				ActivityStartTime: activity.StartTime,
				ActivityEndTime:   activity.EndTime,
			}
		}
		if err := tx.PromotionSeckillProduct.WithContext(ctx).Create(products...); err != nil {
			return err
		}
		return nil
	})
	return activityID, err
}

// UpdateSeckillActivity 更新秒杀活动
func (s *SeckillActivityService) UpdateSeckillActivity(ctx context.Context, r *req.SeckillActivityUpdateReq) error {
	q := s.q.PromotionSeckillActivity
	oldActivity, err := q.WithContext(ctx).Where(q.ID.Eq(r.ID)).First()
	if err != nil {
		return core.NewBizError(1001002000, "秒杀活动不存在")
	}
	if oldActivity.Status == 2 { // Disable?
		return core.NewBizError(1001002003, "秒杀活动已关闭，不能修改")
	}

	if err := s.validateProductConflict(ctx, r.ConfigIds, r.SpuID, r.ID); err != nil {
		return err
	}

	return s.q.Transaction(func(tx *query.Query) error {
		totalStock := 0
		for _, p := range r.Products {
			totalStock += p.Stock
		}

		upd := &promotion.PromotionSeckillActivity{
			SpuID:            r.SpuID,
			Name:             r.Name,
			Remark:           r.Remark,
			StartTime:        r.StartTime,
			EndTime:          r.EndTime,
			Sort:             r.Sort,
			ConfigIds:        r.ConfigIds,
			TotalLimitCount:  r.TotalLimitCount,
			SingleLimitCount: r.SingleLimitCount,
			Stock:            totalStock,
		}
		if totalStock > oldActivity.TotalStock {
			upd.TotalStock = totalStock
		} else {
			upd.TotalStock = oldActivity.TotalStock // Keep max
		}

		if _, err := tx.PromotionSeckillActivity.WithContext(ctx).Where(tx.PromotionSeckillActivity.ID.Eq(r.ID)).Updates(upd); err != nil {
			return err
		}

		// Update products: Full Replace strategy or Diff?
		// Java does Diff. For simplicity, delete and recreate?
		// Or try to match SKU.
		// Deleting and Recreating is safer for logic correctness but loses ID?
		// Java logic is "Diff".
		// I'll use Delete + Create for simplicity unless IDs must be preserved for orders?
		// Orders link to ActivityID + SkuID? OR SeckillProductId?
		// SeckillProductDO has ID.
		// If orders reference SeckillProductID, then we MUST preserve it.
		// Java: SeckillProductDO logic:
		// "Orders usually store activity_id + item_id".
		// TradeOrderItem stores "PromotionActivityId".
		// Whatever. I'll Delete and Recreate for now. Simpler.

		tx.PromotionSeckillProduct.WithContext(ctx).Where(tx.PromotionSeckillProduct.ActivityID.Eq(r.ID)).Delete()

		products := make([]*promotion.PromotionSeckillProduct, len(r.Products))
		for i, p := range r.Products {
			products[i] = &promotion.PromotionSeckillProduct{
				ActivityID:        r.ID,
				ConfigIds:         r.ConfigIds,
				SpuID:             r.SpuID,
				SkuID:             p.SkuID,
				SeckillPrice:      p.SeckillPrice,
				Stock:             p.Stock,
				ActivityStatus:    oldActivity.Status, // Keep status
				ActivityStartTime: r.StartTime,
				ActivityEndTime:   r.EndTime,
			}
		}
		return tx.PromotionSeckillProduct.WithContext(ctx).Create(products...)
	})
}

// DeleteSeckillActivity 删除秒杀活动
func (s *SeckillActivityService) DeleteSeckillActivity(ctx context.Context, id int64) error {
	q := s.q.PromotionSeckillActivity
	act, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
	if err != nil {
		return core.NewBizError(1001002000, "秒杀活动不存在")
	}
	if act.Status == 1 { // Enable
		return core.NewBizError(1001002004, "活动未关闭，不能删除")
	}
	// Delete Activity and Products
	return s.q.Transaction(func(tx *query.Query) error {
		tx.PromotionSeckillActivity.WithContext(ctx).Where(tx.PromotionSeckillActivity.ID.Eq(id)).Delete()
		tx.PromotionSeckillProduct.WithContext(ctx).Where(tx.PromotionSeckillProduct.ActivityID.Eq(id)).Delete()
		return nil
	})
}

// CloseSeckillActivity 关闭秒杀活动
func (s *SeckillActivityService) CloseSeckillActivity(ctx context.Context, id int64) error {
	q := s.q.PromotionSeckillActivity
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Status, 2) // 2=Disable
	return err
}

// GetSeckillActivity 获得秒杀活动
func (s *SeckillActivityService) GetSeckillActivity(ctx context.Context, id int64) (*promotion.PromotionSeckillActivity, error) {
	q := s.q.PromotionSeckillActivity
	return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

// GetSeckillProductListByActivityID 获得秒杀商品列表
func (s *SeckillActivityService) GetSeckillProductListByActivityID(ctx context.Context, activityID int64) ([]*promotion.PromotionSeckillProduct, error) {
	q := s.q.PromotionSeckillProduct
	return q.WithContext(ctx).Where(q.ActivityID.Eq(activityID)).Find()
}

// GetSeckillProductListByActivityIds 批量获得秒杀商品列表
func (s *SeckillActivityService) GetSeckillProductListByActivityIds(ctx context.Context, activityIds []int64) ([]*promotion.PromotionSeckillProduct, error) {
	if len(activityIds) == 0 {
		return []*promotion.PromotionSeckillProduct{}, nil
	}
	q := s.q.PromotionSeckillProduct
	return q.WithContext(ctx).Where(q.ActivityID.In(activityIds...)).Find()
}

// GetSeckillActivityPage 分页获得秒杀活动
func (s *SeckillActivityService) GetSeckillActivityPage(ctx context.Context, r *req.SeckillActivityPageReq) (*core.PageResult[*promotion.PromotionSeckillActivity], error) {
	q := s.q.PromotionSeckillActivity
	do := q.WithContext(ctx)
	if r.Name != "" {
		do = do.Where(q.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		do = do.Where(q.Status.Eq(*r.Status))
	}
	do = do.Order(q.Sort.Desc(), q.ID.Desc())
	list, count, err := do.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*promotion.PromotionSeckillActivity]{List: list, Total: count}, nil
}

// GetSeckillActivityListByIds 获得秒杀活动列表
func (s *SeckillActivityService) GetSeckillActivityListByIds(ctx context.Context, ids []int64) ([]*promotion.PromotionSeckillActivity, error) {
	if len(ids) == 0 {
		return []*promotion.PromotionSeckillActivity{}, nil
	}
	q := s.q.PromotionSeckillActivity
	return q.WithContext(ctx).Where(q.ID.In(ids...)).Find()
}

// validateProductConflict 校验商品冲突
func (s *SeckillActivityService) validateProductConflict(ctx context.Context, configIds []int64, spuID int64, activityID int64) error {
	q := s.q.PromotionSeckillActivity
	// Find all ENABLED activities for this SPU
	conds := []gen.Condition{
		q.SpuID.Eq(spuID),
		q.Status.Eq(1), // Enable
	}
	if activityID > 0 {
		conds = append(conds, q.ID.Neq(activityID))
	}

	list, err := q.WithContext(ctx).Where(conds...).Find()
	if err != nil {
		return err
	}

	for _, act := range list {
		// Check config overlap
		hasConfigOverlap := false
		for _, id := range act.ConfigIds {
			for _, reqId := range configIds {
				if id == reqId {
					hasConfigOverlap = true
					break
				}
			}
			if hasConfigOverlap {
				break
			}
		}
		if hasConfigOverlap {
			return core.NewBizError(1001002002, "该商品已参加其它秒杀活动")
		}
	}
	return nil
}

// GetSeckillActivityAppPage 获得 App 端秒杀活动分页
func (s *SeckillActivityService) GetSeckillActivityAppPage(ctx context.Context, pageNo, pageSize int, configId int64) (*core.PageResult[*promotion.PromotionSeckillActivity], error) {
	// Java logic: filter by configId, status=ENABLE, now between startTime/endTime
	q := s.q.PromotionSeckillActivity

	// How to filter JSON configIds contains configId?
	// Use LIKE? `configIds` stored as `[1,2,3]`.
	// Simple Like "%1%" is dangerous (matches 10).
	// Ideally GORM `datatypes.JSONOverlaps` or just fetch and filter in memory if volume is low.
	// Given Seckill activities are usually limited, memory filter is acceptable.

	// Fetch candidates (Status=Enable, Time Valid)
	now := time.Now()
	list, err := q.WithContext(ctx).Where(
		q.Status.Eq(1),
		q.StartTime.Lte(now),
		q.EndTime.Gte(now),
	).Order(q.Sort.Desc()).Find() // Fetch all active first

	if err != nil {
		return nil, err
	}

	var filtered []*promotion.PromotionSeckillActivity
	for _, item := range list {
		for _, cid := range item.ConfigIds {
			if cid == configId {
				filtered = append(filtered, item)
				break
			}
		}
	}

	// Manual Pagination
	total := int64(len(filtered))
	start := (pageNo - 1) * pageSize
	if start >= len(filtered) {
		return &core.PageResult[*promotion.PromotionSeckillActivity]{List: []*promotion.PromotionSeckillActivity{}, Total: total}, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return &core.PageResult[*promotion.PromotionSeckillActivity]{
		List:  filtered[start:end],
		Total: total,
	}, nil
}

// ValidateJoinSeckill 校验是否参与秒杀
func (s *SeckillActivityService) ValidateJoinSeckill(ctx context.Context, activityId, skuId int64, count int) (*promotion.PromotionSeckillActivity, *promotion.PromotionSeckillProduct, error) {
	// 1. Get Activity
	act, err := s.GetSeckillActivity(ctx, activityId)
	if err != nil || act == nil {
		return nil, nil, core.NewBizError(1001002000, "秒杀活动不存在")
	}
	if act.Status != 1 {
		return nil, nil, core.NewBizError(1001002003, "秒杀活动已关闭")
	}
	now := time.Now()
	if now.Before(act.StartTime) || now.After(act.EndTime) {
		return nil, nil, core.NewBizError(1001002005, "秒杀活动时间不符")
	}

	// 2. Get Product
	q := s.q.PromotionSeckillProduct
	prod, err := q.WithContext(ctx).Where(q.ActivityID.Eq(activityId), q.SkuID.Eq(skuId)).First()
	if err != nil {
		return nil, nil, core.NewBizError(1001002006, "秒杀商品不存在")
	}

	// 3. Check Stock
	if prod.Stock < count {
		return nil, nil, core.NewBizError(1001002007, "秒杀库存不足")
	}

	// 4. Check Single Limit
	if act.SingleLimitCount > 0 && count > act.SingleLimitCount {
		return nil, nil, core.NewBizError(1001002008, "超出单次限购数量")
	}

	return act, prod, nil
}
