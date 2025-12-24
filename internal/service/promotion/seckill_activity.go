package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product" // Import Product services
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
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
	// 1.3 校验商品是否存在
	if _, err := s.spuSvc.GetSpu(ctx, r.SpuID); err != nil {
		return 0, err
	}
	for _, p := range r.Products {
		if _, err := s.skuSvc.GetSku(ctx, p.SkuID); err != nil {
			return 0, err
		}
	}

	// 2. Transaction
	var activityID int64
	err := s.q.Transaction(func(tx *query.Query) error {
		// Calculate Stock
		totalStock := lo.SumBy(r.Products, func(p req.SeckillProductBaseVO) int {
			return p.Stock
		})

		// Insert Activity
		activity := &promotion.PromotionSeckillActivity{
			SpuID:            r.SpuID,
			Name:             r.Name,
			Status:           model.CommonStatusEnable, // 使用 CommonStatusEnable 常量
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
		return errors.NewBizError(1001002000, "秒杀活动不存在")
	}
	if oldActivity.Status == model.CommonStatusDisable { // 使用 CommonStatusDisable 常量替代魔法数字
		return errors.NewBizError(1001002003, "秒杀活动已关闭，不能修改")
	}

	if err := s.validateProductConflict(ctx, r.ConfigIds, r.SpuID, r.ID); err != nil {
		return err
	}
	// 校验商品是否存在
	if _, err := s.spuSvc.GetSpu(ctx, r.SpuID); err != nil {
		return err
	}
	for _, p := range r.Products {
		if _, err := s.skuSvc.GetSku(ctx, p.SkuID); err != nil {
			return err
		}
	}

	return s.q.Transaction(func(tx *query.Query) error {
		totalStock := lo.SumBy(r.Products, func(p req.SeckillProductBaseVO) int {
			return p.Stock
		})

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

		// Update products
		if _, err := tx.PromotionSeckillProduct.WithContext(ctx).Where(tx.PromotionSeckillProduct.ActivityID.Eq(r.ID)).Delete(); err != nil {
			return err
		}

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
		return errors.NewBizError(1001002000, "秒杀活动不存在")
	}
	if act.Status == model.CommonStatusEnable { // 使用 CommonStatusEnable 常量替代魔法数字
		return errors.NewBizError(1001002004, "活动未关闭，不能删除")
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
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Status, model.CommonStatusDisable) // 使用 CommonStatusDisable 常量
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
func (s *SeckillActivityService) GetSeckillActivityPage(ctx context.Context, r *req.SeckillActivityPageReq) (*pagination.PageResult[*promotion.PromotionSeckillActivity], error) {
	q := s.q.PromotionSeckillActivity
	do := q.WithContext(ctx)
	if r.Name != "" {
		do = do.Where(q.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		do = do.Where(q.Status.Eq(*r.Status))
	}
	if len(r.CreateTime) == 2 && r.CreateTime[0] != nil && r.CreateTime[1] != nil {
		do = do.Where(q.CreateTime.Between(*r.CreateTime[0], *r.CreateTime[1]))
	}
	do = do.Order(q.Sort.Desc(), q.ID.Desc())
	list, count, err := do.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[*promotion.PromotionSeckillActivity]{List: list, Total: count}, nil
}

// GetSeckillActivityListByIds 获得秒杀活动列表
func (s *SeckillActivityService) GetSeckillActivityListByIds(ctx context.Context, ids []int64) ([]*promotion.PromotionSeckillActivity, error) {
	if len(ids) == 0 {
		return []*promotion.PromotionSeckillActivity{}, nil
	}
	q := s.q.PromotionSeckillActivity
	return q.WithContext(ctx).Where(q.ID.In(ids...)).Find()
}

// GetSeckillActivityListByConfigId 按秒杀时段获取活动列表 (App 端)
func (s *SeckillActivityService) GetSeckillActivityListByConfigId(ctx context.Context, configId int64, limit int) ([]*promotion.PromotionSeckillActivity, error) {
	q := s.q.PromotionSeckillActivity
	now := time.Now()
	list, err := q.WithContext(ctx).Where(
		q.Status.Eq(model.CommonStatusEnable), // 使用 CommonStatusEnable 常量替代魔法数字 1
		q.StartTime.Lte(now),
		q.EndTime.Gte(now),
	).Order(q.Sort.Desc()).Find()
	if err != nil {
		return nil, err
	}

	// 过滤包含此时段的活动
	filtered := lo.Filter(list, func(item *promotion.PromotionSeckillActivity, _ int) bool {
		return lo.Contains(item.ConfigIds, configId)
	})
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered, nil
}

// GetSeckillActivityPageForApp App 端秒杀活动分页 (别名)
func (s *SeckillActivityService) GetSeckillActivityPageForApp(ctx context.Context, configId *int64, pageNo, pageSize int) (*pagination.PageResult[*promotion.PromotionSeckillActivity], error) {
	q := s.q.PromotionSeckillActivity
	now := time.Now()
	list, err := q.WithContext(ctx).Where(
		q.Status.Eq(model.CommonStatusEnable), // 使用 CommonStatusEnable 常量替代魔法数字 1
		q.StartTime.Lte(now),
		q.EndTime.Gte(now),
	).Order(q.Sort.Desc()).Find()
	if err != nil {
		return nil, err
	}

	// 如果指定了 configId，过滤
	var filtered []*promotion.PromotionSeckillActivity
	if configId != nil {
		filtered = lo.Filter(list, func(item *promotion.PromotionSeckillActivity, _ int) bool {
			return lo.Contains(item.ConfigIds, *configId)
		})
	} else {
		filtered = list
	}

	// 手动分页
	total := int64(len(filtered))
	start := (pageNo - 1) * pageSize
	if start >= len(filtered) {
		return &pagination.PageResult[*promotion.PromotionSeckillActivity]{List: []*promotion.PromotionSeckillActivity{}, Total: total}, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return &pagination.PageResult[*promotion.PromotionSeckillActivity]{
		List:  filtered[start:end],
		Total: total,
	}, nil
}

// validateProductConflict 校验商品冲突
func (s *SeckillActivityService) validateProductConflict(ctx context.Context, configIds []int64, spuID int64, activityID int64) error {
	q := s.q.PromotionSeckillActivity
	// Find all ENABLED activities for this SPU
	conds := []gen.Condition{
		q.SpuID.Eq(spuID),
		q.Status.Eq(model.CommonStatusEnable), // 使用 CommonStatusEnable 常量替代魔法数字 (Enable)
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
		if len(lo.Intersect(act.ConfigIds, configIds)) > 0 {
			return errors.NewBizError(1001002002, "该商品已参加其它秒杀活动")
		}
	}
	return nil
}

// GetSeckillActivityAppPage 获得 App 端秒杀活动分页
func (s *SeckillActivityService) GetSeckillActivityAppPage(ctx context.Context, pageNo, pageSize int, configId int64) (*pagination.PageResult[*promotion.PromotionSeckillActivity], error) {
	// Java logic: filter by configId, status=ENABLE, now between startTime/endTime
	q := s.q.PromotionSeckillActivity

	// Fetch candidates (Status=Enable, Time Valid)
	now := time.Now()
	list, err := q.WithContext(ctx).Where(
		q.Status.Eq(model.CommonStatusEnable), // 使用 CommonStatusEnable 常量替代魔法数字 1
		q.StartTime.Lte(now),
		q.EndTime.Gte(now),
	).Order(q.Sort.Desc()).Find() // Fetch all active first

	if err != nil {
		return nil, err
	}

	filtered := lo.Filter(list, func(item *promotion.PromotionSeckillActivity, _ int) bool {
		return lo.Contains(item.ConfigIds, configId)
	})

	// Manual Pagination
	total := int64(len(filtered))
	start := (pageNo - 1) * pageSize
	if start >= len(filtered) {
		return &pagination.PageResult[*promotion.PromotionSeckillActivity]{List: []*promotion.PromotionSeckillActivity{}, Total: total}, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return &pagination.PageResult[*promotion.PromotionSeckillActivity]{
		List:  filtered[start:end],
		Total: total,
	}, nil
}

// ValidateJoinSeckill 校验是否参与秒杀
func (s *SeckillActivityService) ValidateJoinSeckill(ctx context.Context, activityId, skuId int64, count int) (*promotion.PromotionSeckillActivity, *promotion.PromotionSeckillProduct, error) {
	// 1. Get Activity
	act, err := s.GetSeckillActivity(ctx, activityId)
	if err != nil || act == nil {
		return nil, nil, errors.NewBizError(1001002000, "秒杀活动不存在")
	}
	if act.Status != model.CommonStatusEnable { // 使用 CommonStatusEnable 常量
		return nil, nil, errors.NewBizError(1001002003, "秒杀活动已关闭")
	}
	now := time.Now()
	if now.Before(act.StartTime) || now.After(act.EndTime) {
		return nil, nil, errors.NewBizError(1001002005, "秒杀活动时间不符")
	}

	// 3. Check Config (Time Segment)
	config, err := s.configSvc.GetCurrentSeckillConfig(ctx)
	if err != nil || config == nil || !lo.Contains(act.ConfigIds, config.ID) {
		return nil, nil, errors.NewBizError(1001002005, "秒杀活动时间不符")
	}

	// 4. Check Single Limit
	if act.SingleLimitCount > 0 && count > act.SingleLimitCount {
		return nil, nil, errors.NewBizError(1001002008, "超出单次限购数量")
	}

	// 5. Get Product
	q := s.q.PromotionSeckillProduct
	prod, err := q.WithContext(ctx).Where(q.ActivityID.Eq(activityId), q.SkuID.Eq(skuId)).First()
	if err != nil {
		return nil, nil, errors.NewBizError(1001002006, "秒杀商品不存在")
	}

	// 6. Check Stock
	if prod.Stock < count {
		return nil, nil, errors.NewBizError(1001002007, "秒杀库存不足")
	}

	return act, prod, nil
}

// UpdateSeckillStockDecr 扣减秒杀库存 (针对订单提交)
func (s *SeckillActivityService) UpdateSeckillStockDecr(ctx context.Context, id int64, skuId int64, count int) error {
	return s.q.Transaction(func(tx *query.Query) error {
		// 1.1 校验活动库存是否充足
		act, err := tx.PromotionSeckillActivity.WithContext(ctx).Where(tx.PromotionSeckillActivity.ID.Eq(id)).First()
		if err != nil || act.Stock < count {
			return errors.NewBizError(1001002007, "秒杀库存不足")
		}
		// 1.2 校验商品库存是否充足
		prod, err := tx.PromotionSeckillProduct.WithContext(ctx).Where(
			tx.PromotionSeckillProduct.ActivityID.Eq(id),
			tx.PromotionSeckillProduct.SkuID.Eq(skuId),
		).First()
		if err != nil || prod.Stock < count {
			return errors.NewBizError(1001002007, "秒杀库存不足")
		}

		res, err := tx.PromotionSeckillProduct.WithContext(ctx).Where(
			tx.PromotionSeckillProduct.ID.Eq(prod.ID),
			tx.PromotionSeckillProduct.Stock.Gte(count),
		).Update(tx.PromotionSeckillProduct.Stock, tx.PromotionSeckillProduct.Stock.Add(-count))
		if err != nil || res.RowsAffected == 0 {
			return errors.NewBizError(1001002007, "秒杀库存不足")
		}

		// 2.2 更新活动库存
		res, err = tx.PromotionSeckillActivity.WithContext(ctx).Where(
			tx.PromotionSeckillActivity.ID.Eq(id),
			tx.PromotionSeckillActivity.Stock.Gte(count),
		).Update(tx.PromotionSeckillActivity.Stock, tx.PromotionSeckillActivity.Stock.Add(-count))
		if err != nil || res.RowsAffected == 0 {
			return errors.NewBizError(1001002007, "秒杀库存不足")
		}
		return nil
	})
}

// UpdateSeckillStockIncr 增加秒杀库存 (针对订单取消/退款)
func (s *SeckillActivityService) UpdateSeckillStockIncr(ctx context.Context, id int64, skuId int64, count int) error {
	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 更新活动商品库存
		_, err := tx.PromotionSeckillProduct.WithContext(ctx).Where(
			tx.PromotionSeckillProduct.ActivityID.Eq(id),
			tx.PromotionSeckillProduct.SkuID.Eq(skuId),
		).Update(tx.PromotionSeckillProduct.Stock, tx.PromotionSeckillProduct.Stock.Add(count))
		if err != nil {
			return err
		}

		// 2. 更新活动库存
		_, err = tx.PromotionSeckillActivity.WithContext(ctx).Where(
			tx.PromotionSeckillActivity.ID.Eq(id),
		).Update(tx.PromotionSeckillActivity.Stock, tx.PromotionSeckillActivity.Stock.Add(count))
		return err
	})
}

// GetMatchSeckillActivityBySpuId 获取指定 SPU 的进行中的秒杀活动
func (s *SeckillActivityService) GetMatchSeckillActivityBySpuId(ctx context.Context, spuId int64) (*promotion.PromotionSeckillActivity, error) {
	now := time.Now()
	q := s.q.PromotionSeckillActivity
	return q.WithContext(ctx).
		Where(q.SpuID.Eq(spuId)).
		Where(q.Status.Eq(model.CommonStatusEnable)).
		Where(q.StartTime.Lt(now)).
		Where(q.EndTime.Gt(now)).
		First()
}
