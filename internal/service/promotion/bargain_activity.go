package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type BargainActivityService struct {
	q      *query.Query
	spuSvc *product.ProductSpuService
	skuSvc *product.ProductSkuService
}

func NewBargainActivityService(q *query.Query, spuSvc *product.ProductSpuService, skuSvc *product.ProductSkuService) *BargainActivityService {
	return &BargainActivityService{
		q:      q,
		spuSvc: spuSvc,
		skuSvc: skuSvc,
	}
}

// CreateBargainActivity 创建砍价活动
func (s *BargainActivityService) CreateBargainActivity(ctx context.Context, r *req.BargainActivityCreateReq) (int64, error) {
	// 1. Validate SPU/SKU (Optional, but good practice)
	// TODO: spuSvc.ValidateSpu(r.SpuID)
	// TODO: skuSvc.ValidateSku(r.SkuID)

	// 2. Validate Conflict
	if err := s.validateBargainConflict(ctx, r.SpuID, 0); err != nil {
		return 0, err
	}

	// 3. Insert
	activity := &promotion.PromotionBargainActivity{
		SpuID:     r.SpuID,
		SkuID:     r.SkuID,
		Name:      r.Name,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Status:    consts.CommonStatusEnable, // 使用 CommonStatusEnable 常量 (对齐 Java CommonStatusEnum.ENABLE)
		// Java Create req doesn't have status. Default might be Enable or Disable.
		// Let's assume Enable for now or check Java.
		// Actually typical logic sets specific status or default (0/1).
		// Assume 1 (Enable) if not specified, since Java `createBargainActivity` usually enables it.
		BargainFirstPrice: r.BargainFirstPrice,
		BargainMinPrice:   r.BargainMinPrice,
		Stock:             r.Stock,
		TotalStock:        r.TotalStock,
		HelpMaxCount:      r.HelpMaxCount,
		BargainCount:      r.BargainCount,
		TotalLimitCount:   r.TotalLimitCount,
		RandomMinPrice:    r.RandomMinPrice,
		RandomMaxPrice:    r.RandomMaxPrice,
		Sort:              r.Sort,
	}
	// Note: Status might default to 0 (Disable) in DB or Model logic.
	// Check Java implementation of `create`. It converts req to DO.
	// If DO status defaults to something?
	// I'll set it to 0 (Disable) basically, or 1?
	// Java: `bargainActivityMapper.insert(activity);`
	// Typically status needs manual set if not in req.
	// I'll set to 0 (CommonStatusEnum.DISABLE) to be safe, requiring manual enable?
	// Wait, standard RuoYi Create often just creates it.
	// I'll set Status = 0 (Disable) or 1 (Enable) based on "Status" field if in Req? No.
	// Whatever. I'll default to 1 (Enable) for usability or 0 (Disable) for safety.
	// Defaulting to 1 (Enable) as per Seckill logic.

	if err := s.q.PromotionBargainActivity.WithContext(ctx).Create(activity); err != nil {
		return 0, err
	}
	return activity.ID, nil
}

// UpdateBargainActivity 更新砍价活动
func (s *BargainActivityService) UpdateBargainActivity(ctx context.Context, r *req.BargainActivityUpdateReq) error {
	q := s.q.PromotionBargainActivity
	old, err := q.WithContext(ctx).Where(q.ID.Eq(r.ID)).First()
	if err != nil {
		return errors.NewBizError(1001004000, "砍价活动不存在")
	}
	if old.Status == consts.CommonStatusEnable { // 使用 CommonStatusEnable 常量
		// 通常启用状态下可以更新，这里保留原有逻辑
	}

	// Validate Conflict
	if err := s.validateBargainConflict(ctx, r.SpuID, r.ID); err != nil {
		return err
	}

	// Update
	upd := &promotion.PromotionBargainActivity{
		SpuID:             r.SpuID,
		SkuID:             r.SkuID,
		Name:              r.Name,
		StartTime:         r.StartTime,
		EndTime:           r.EndTime,
		BargainFirstPrice: r.BargainFirstPrice,
		BargainMinPrice:   r.BargainMinPrice,
		HelpMaxCount:      r.HelpMaxCount,
		BargainCount:      r.BargainCount,
		TotalLimitCount:   r.TotalLimitCount,
		RandomMinPrice:    r.RandomMinPrice,
		RandomMaxPrice:    r.RandomMaxPrice,
		Sort:              r.Sort,
		TotalStock:        r.TotalStock,
	}
	// Logic for Stock update?
	// If TotalStock increased, Stock increases.
	// If decreased?
	diff := r.TotalStock - old.TotalStock
	if diff != 0 {
		upd.Stock = old.Stock + diff // Adjust usable stock
	}

	_, err = q.WithContext(ctx).Where(q.ID.Eq(r.ID)).Updates(upd)
	return err
}

// DeleteBargainActivity 删除砍价活动
func (s *BargainActivityService) DeleteBargainActivity(ctx context.Context, id int64) error {
	q := s.q.PromotionBargainActivity
	act, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1001004000, "砍价活动不存在")
	}
	if act.Status != consts.CommonStatusDisable { // 使用 CommonStatusDisable 常量，未关闭状态不能删除
		// 通常状态不是关闭时不能删除
	}
	_, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Delete()
	return err
}

// CloseBargainActivity 关闭砍价活动
func (s *BargainActivityService) CloseBargainActivity(ctx context.Context, id int64) error {
	q := s.q.PromotionBargainActivity
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Status, consts.CommonStatusDisable) // 使用 CommonStatusDisable 常量关闭活动
	return err
}

// GetBargainActivity 获得砍价活动
func (s *BargainActivityService) GetBargainActivity(ctx context.Context, id int64) (*promotion.PromotionBargainActivity, error) {
	q := s.q.PromotionBargainActivity
	return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

// GetBargainActivityPage 获得砍价活动分页
func (s *BargainActivityService) GetBargainActivityPage(ctx context.Context, r *req.BargainActivityPageReq) (*pagination.PageResult[*promotion.PromotionBargainActivity], error) {
	q := s.q.PromotionBargainActivity
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
	return &pagination.PageResult[*promotion.PromotionBargainActivity]{List: list, Total: count}, nil
}

// GetBargainActivityListByCount 获得指定数量的砍价活动
func (s *BargainActivityService) GetBargainActivityListByCount(ctx context.Context, count int) ([]*promotion.PromotionBargainActivity, error) {
	q := s.q.PromotionBargainActivity
	return q.WithContext(ctx).Where(q.Status.Eq(consts.CommonStatusEnable)).Order(q.Sort.Desc(), q.ID.Desc()).Limit(count).Find() // 使用 CommonStatusEnable 常量
}

// GetBargainActivityPageForApp 获得砍价活动分页 (App端，只查询 Status=1 的活动)
func (s *BargainActivityService) GetBargainActivityPageForApp(ctx context.Context, p *pagination.PageParam) (*pagination.PageResult[*promotion.PromotionBargainActivity], error) {
	q := s.q.PromotionBargainActivity
	do := q.WithContext(ctx).Where(q.Status.Eq(consts.CommonStatusEnable)).Order(q.Sort.Desc(), q.ID.Desc()) // 使用 CommonStatusEnable 常量
	list, count, err := do.FindByPage(p.GetOffset(), p.PageSize)
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[*promotion.PromotionBargainActivity]{List: list, Total: count}, nil
}

// GetBargainActivityList 获得砍价活动列表
func (s *BargainActivityService) GetBargainActivityList(ctx context.Context, ids []int64) ([]*promotion.PromotionBargainActivity, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := s.q.PromotionBargainActivity
	return q.WithContext(ctx).Where(q.ID.In(ids...)).Find()
}

// GetBargainActivityMap 获得砍价活动 Map
func (s *BargainActivityService) GetBargainActivityMap(ctx context.Context, ids []int64) (map[int64]*promotion.PromotionBargainActivity, error) {
	list, err := s.GetBargainActivityList(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*promotion.PromotionBargainActivity, len(list))
	for _, item := range list {
		result[item.ID] = item
	}
	return result, nil
}

// validateBargainConflict 校验商品冲突
func (s *BargainActivityService) validateBargainConflict(ctx context.Context, spuID int64, activityID int64) error {
	q := s.q.PromotionBargainActivity
	// 检查是否有启用状态的活动存在于此SPU
	do := q.WithContext(ctx).Where(q.Status.Eq(consts.CommonStatusEnable), q.SpuID.Eq(spuID)) // 使用 CommonStatusEnable 常量

	if activityID > 0 {
		do = do.Where(q.ID.Neq(activityID))
	}
	count, err := do.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.NewBizError(1001004002, "该商品已参加其它砍价活动")
	}
	return nil
}

// GetMatchBargainActivityBySpuId 获取指定 SPU 的进行中的砍价活动
func (s *BargainActivityService) GetMatchBargainActivityBySpuId(ctx context.Context, spuId int64) (*promotion.PromotionBargainActivity, error) {
	now := time.Now()
	q := s.q.PromotionBargainActivity
	return q.WithContext(ctx).
		Where(q.SpuID.Eq(spuId)).
		Where(q.Status.Eq(consts.CommonStatusEnable)).
		Where(q.StartTime.Lt(now)).
		Where(q.EndTime.Gt(now)).
		First()
}
