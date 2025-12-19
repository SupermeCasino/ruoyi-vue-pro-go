package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	prodSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type CombinationActivityService interface {
	// Admin
	CreateCombinationActivity(ctx context.Context, req req.CombinationActivityCreateReq) (int64, error)
	UpdateCombinationActivity(ctx context.Context, req req.CombinationActivityUpdateReq) error
	CloseCombinationActivity(ctx context.Context, id int64) error
	DeleteCombinationActivity(ctx context.Context, id int64) error
	GetCombinationActivity(ctx context.Context, id int64) (*resp.CombinationActivityRespVO, error)
	GetCombinationActivityPage(ctx context.Context, req req.CombinationActivityPageReq) (*pagination.PageResult[*resp.CombinationActivityRespVO], error)
	GetCombinationActivityMap(ctx context.Context, ids []int64) (map[int64]*promotion.PromotionCombinationActivity, error)
	GetCombinationActivityListByIds(ctx context.Context, ids []int64) ([]*resp.CombinationActivityRespVO, error)

	// App
	GetCombinationActivityList(ctx context.Context, count int) ([]*resp.AppCombinationActivityRespVO, error)
	GetCombinationActivityPageForApp(ctx context.Context, req pagination.PageParam) (*pagination.PageResult[*resp.AppCombinationActivityRespVO], error)
	GetCombinationActivityDetail(ctx context.Context, id int64) (*resp.AppCombinationActivityDetailRespVO, error)
	ValidateCombinationActivityCanJoin(ctx context.Context, activityID int64) (*promotion.PromotionCombinationActivity, error)
	// GetMatchCombinationActivityBySpuId 获取指定 SPU 的进行中的拼团活动
	GetMatchCombinationActivityBySpuId(ctx context.Context, spuId int64) (*promotion.PromotionCombinationActivity, error)
}

type combinationActivityService struct {
	q      *query.Query
	spuSvc *prodSvc.ProductSpuService
	skuSvc *prodSvc.ProductSkuService
}

func NewCombinationActivityService(q *query.Query, spuSvc *prodSvc.ProductSpuService, skuSvc *prodSvc.ProductSkuService) CombinationActivityService {
	return &combinationActivityService{
		q:      q,
		spuSvc: spuSvc,
		skuSvc: skuSvc,
	}
}

func (s *combinationActivityService) CreateCombinationActivity(ctx context.Context, req req.CombinationActivityCreateReq) (int64, error) {
	// 1.1 校验商品
	if err := s.validateProducts(ctx, req.Products); err != nil {
		return 0, err
	}
	// 1.2 校验商品冲突
	if err := s.validateProductConflict(ctx, req.SpuID, 0); err != nil {
		return 0, err
	}

	// 2. 插入活动
	activity := &promotion.PromotionCombinationActivity{
		Name:             req.Name,
		SpuID:            req.SpuID,
		TotalLimitCount:  req.TotalLimitCount,
		SingleLimitCount: req.SingleLimitCount,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		UserSize:         req.UserSize,
		VirtualGroup:     req.VirtualGroup,
		Status:           1, // 1: Enable
		LimitDuration:    req.LimitDuration,
	}

	err := s.q.Transaction(func(tx *query.Query) error {
		if err := tx.PromotionCombinationActivity.WithContext(ctx).Create(activity); err != nil {
			return err
		}

		// 3. 插入商品
		products := make([]*promotion.PromotionCombinationProduct, len(req.Products))
		for i, p := range req.Products {
			products[i] = &promotion.PromotionCombinationProduct{
				ActivityID:        activity.ID,
				SpuID:             p.SpuID,
				SkuID:             p.SkuID,
				CombinationPrice:  p.CombinationPrice,
				ActivityStatus:    activity.Status,
				ActivityStartTime: activity.StartTime,
				ActivityEndTime:   activity.EndTime,
			}
		}
		if err := tx.PromotionCombinationProduct.WithContext(ctx).Create(products...); err != nil {
			return err
		}
		return nil
	})

	return activity.ID, err
}

func (s *combinationActivityService) UpdateCombinationActivity(ctx context.Context, req req.CombinationActivityUpdateReq) error {
	// 1. 校验是否存在
	old, err := s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.NewBizError(1001006000, "拼团活动不存在")
	}
	if old.Status == 0 { // Disable
		return errors.NewBizError(1001006010, "拼团活动已关闭，不能修改")
	}

	// 2.1 校验商品
	if err := s.validateProducts(ctx, req.Products); err != nil {
		return err
	}
	// 2.2 校验商品冲突
	if err := s.validateProductConflict(ctx, req.SpuID, req.ID); err != nil {
		return err
	}

	// 3. 更新
	activity := &promotion.PromotionCombinationActivity{
		ID:               req.ID,
		Name:             req.Name,
		SpuID:            req.SpuID,
		TotalLimitCount:  req.TotalLimitCount,
		SingleLimitCount: req.SingleLimitCount,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		UserSize:         req.UserSize,
		VirtualGroup:     req.VirtualGroup,
		LimitDuration:    req.LimitDuration,
	}

	return s.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.PromotionCombinationActivity.WithContext(ctx).Where(tx.PromotionCombinationActivity.ID.Eq(req.ID)).Updates(activity); err != nil {
			return err
		}

		// 删除旧商品
		if _, err := tx.PromotionCombinationProduct.WithContext(ctx).Where(tx.PromotionCombinationProduct.ActivityID.Eq(req.ID)).Delete(); err != nil {
			return err
		}

		// 插入新商品
		products := make([]*promotion.PromotionCombinationProduct, len(req.Products))
		for i, p := range req.Products {
			products[i] = &promotion.PromotionCombinationProduct{
				ActivityID:        activity.ID,
				SpuID:             p.SpuID,
				SkuID:             p.SkuID,
				CombinationPrice:  p.CombinationPrice,
				ActivityStatus:    old.Status,
				ActivityStartTime: activity.StartTime,
				ActivityEndTime:   activity.EndTime,
			}
		}
		if err := tx.PromotionCombinationProduct.WithContext(ctx).Create(products...); err != nil {
			return err
		}
		return nil
	})
}

func (s *combinationActivityService) DeleteCombinationActivity(ctx context.Context, id int64) error {
	activity, err := s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1001006000, "拼团活动不存在")
	}
	if activity.Status == 1 { // Enable
		return errors.NewBizError(1001006011, "拼团活动进行中，无法删除")
	}
	_, err = s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.Eq(id)).Delete()
	return err
}

// CloseCombinationActivity 关闭拼团活动
func (s *combinationActivityService) CloseCombinationActivity(ctx context.Context, id int64) error {
	q := s.q.PromotionCombinationActivity
	activity, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1001006000, "拼团活动不存在")
	}
	if activity.Status == 0 { // 已禁用
		return errors.NewBizError(1001006012, "拼团活动已关闭")
	}
	_, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Status, 0) // 0 = Disable
	return err
}

func (s *combinationActivityService) validateProductConflict(ctx context.Context, spuID int64, activityID int64) error {
	q := s.q.PromotionCombinationActivity
	query := q.WithContext(ctx).Where(q.Status.Eq(1), q.SpuID.Eq(spuID)) // Enable & SpuID match
	if activityID > 0 {
		query = query.Where(q.ID.Neq(activityID))
	}
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.NewBizError(1001006008, "该商品已存在于其他拼团活动中")
	}
	return nil
}

func (s *combinationActivityService) GetCombinationActivity(ctx context.Context, id int64) (*resp.CombinationActivityRespVO, error) {
	activity, err := s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(1001006000, "拼团活动不存在")
	}
	prods, err := s.q.PromotionCombinationProduct.WithContext(ctx).Where(s.q.PromotionCombinationProduct.ActivityID.Eq(id)).Find()
	if err != nil {
		return nil, err
	}

	vo := &resp.CombinationActivityRespVO{
		ID:               activity.ID,
		Name:             activity.Name,
		SpuID:            activity.SpuID,
		TotalLimitCount:  activity.TotalLimitCount,
		SingleLimitCount: activity.SingleLimitCount,
		StartTime:        activity.StartTime,
		EndTime:          activity.EndTime,
		UserSize:         activity.UserSize,
		VirtualGroup:     activity.VirtualGroup,
		LimitDuration:    activity.LimitDuration,
		Status:           activity.Status,
		CreateTime:       activity.CreatedAt,
		Products:         make([]resp.CombinationProductRespVO, len(prods)),
	}

	for i, p := range prods {
		vo.Products[i] = resp.CombinationProductRespVO{
			SpuID:             p.SpuID,
			SkuID:             p.SkuID,
			CombinationPrice:  p.CombinationPrice,
			ActivityStatus:    p.ActivityStatus,
			ActivityStartTime: p.ActivityStartTime,
			ActivityEndTime:   p.ActivityEndTime,
		}
	}
	return vo, nil
}

func (s *combinationActivityService) GetCombinationActivityPage(ctx context.Context, req req.CombinationActivityPageReq) (*pagination.PageResult[*resp.CombinationActivityRespVO], error) {
	q := s.q.PromotionCombinationActivity.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(s.q.PromotionCombinationActivity.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		q = q.Where(s.q.PromotionCombinationActivity.Status.Eq(*req.Status))
	}
	if len(req.CreateTime) == 2 && req.CreateTime[0] != nil && req.CreateTime[1] != nil {
		q = q.Where(s.q.PromotionCombinationActivity.CreatedAt.Between(*req.CreateTime[0], *req.CreateTime[1]))
	}

	list, total, err := q.Order(s.q.PromotionCombinationActivity.ID.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	result := make([]*resp.CombinationActivityRespVO, len(list))
	for i, item := range list {
		result[i] = &resp.CombinationActivityRespVO{
			ID:               item.ID,
			Name:             item.Name,
			SpuID:            item.SpuID,
			TotalLimitCount:  item.TotalLimitCount,
			SingleLimitCount: item.SingleLimitCount,
			StartTime:        item.StartTime,
			EndTime:          item.EndTime,
			UserSize:         item.UserSize,
			VirtualGroup:     item.VirtualGroup,
			LimitDuration:    item.LimitDuration,
			Status:           item.Status,
			CreateTime:       item.CreatedAt,
		}
	}
	return &pagination.PageResult[*resp.CombinationActivityRespVO]{
		List:  result,
		Total: total,
	}, nil
}

func (s *combinationActivityService) GetCombinationActivityMap(ctx context.Context, ids []int64) (map[int64]*promotion.PromotionCombinationActivity, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	list, err := s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*promotion.PromotionCombinationActivity, len(list))
	for _, item := range list {
		result[item.ID] = item
	}
	return result, nil
}

// GetCombinationActivityListByIds 获得拼团活动列表，基于活动编号数组
// Java: CombinationActivityController#getCombinationActivityListByIds
func (s *combinationActivityService) GetCombinationActivityListByIds(ctx context.Context, ids []int64) ([]*resp.CombinationActivityRespVO, error) {
	if len(ids) == 0 {
		return []*resp.CombinationActivityRespVO{}, nil
	}

	// 1. 获得开启的活动列表
	list, err := s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	// 过滤 disabled 状态 (Status == 0)
	enabledList := make([]*promotion.PromotionCombinationActivity, 0)
	for _, activity := range list {
		if activity.Status != 0 { // 0 = Disable
			enabledList = append(enabledList, activity)
		}
	}

	if len(enabledList) == 0 {
		return []*resp.CombinationActivityRespVO{}, nil
	}

	// 2. 获取 Product 列表
	activityIds := make([]int64, len(enabledList))
	for i, activity := range enabledList {
		activityIds[i] = activity.ID
	}
	productList, err := s.q.PromotionCombinationProduct.WithContext(ctx).
		Where(s.q.PromotionCombinationProduct.ActivityID.In(activityIds...)).
		Find()
	if err != nil {
		return nil, err
	}

	// 3. 获取 SPU 列表
	spuIds := make([]int64, len(enabledList))
	for i, activity := range enabledList {
		spuIds[i] = activity.SpuID
	}
	spuList, err := s.spuSvc.GetSpuList(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := make(map[int64]*resp.ProductSpuResp)
	for _, spu := range spuList {
		spuMap[spu.ID] = spu
	}

	// 4. 组合返回数据
	productMap := make(map[int64][]resp.CombinationProductRespVO)
	for _, prod := range productList {
		productMap[prod.ActivityID] = append(productMap[prod.ActivityID], resp.CombinationProductRespVO{
			SpuID:             prod.SpuID,
			SkuID:             prod.SkuID,
			CombinationPrice:  prod.CombinationPrice,
			ActivityStatus:    prod.ActivityStatus,
			ActivityStartTime: prod.ActivityStartTime,
			ActivityEndTime:   prod.ActivityEndTime,
		})
	}

	result := make([]*resp.CombinationActivityRespVO, len(enabledList))
	for i, activity := range enabledList {
		result[i] = &resp.CombinationActivityRespVO{
			ID:               activity.ID,
			Name:             activity.Name,
			SpuID:            activity.SpuID,
			TotalLimitCount:  activity.TotalLimitCount,
			SingleLimitCount: activity.SingleLimitCount,
			StartTime:        activity.StartTime,
			EndTime:          activity.EndTime,
			UserSize:         activity.UserSize,
			VirtualGroup:     activity.VirtualGroup,
			LimitDuration:    activity.LimitDuration,
			Status:           activity.Status,
			CreateTime:       activity.CreatedAt,
			Products:         productMap[activity.ID],
		}
	}

	return result, nil
}

func (s *combinationActivityService) GetCombinationActivityList(ctx context.Context, count int) ([]*resp.AppCombinationActivityRespVO, error) {
	q := s.q.PromotionCombinationActivity
	list, err := q.WithContext(ctx).
		Where(q.Status.Eq(1)). // Enable
		Order(q.ID.Desc()).    // Usually Sort desc
		Limit(count).
		Find()
	if err != nil {
		return nil, err
	}

	return s.buildAppActivityList(ctx, list)
}

func (s *combinationActivityService) GetCombinationActivityPageForApp(ctx context.Context, p pagination.PageParam) (*pagination.PageResult[*resp.AppCombinationActivityRespVO], error) {
	q := s.q.PromotionCombinationActivity
	list, total, err := q.WithContext(ctx).
		Where(q.Status.Eq(1)). // Enable
		Order(q.ID.Desc()).
		FindByPage(p.GetOffset(), p.GetLimit())
	if err != nil {
		return nil, err
	}

	vos, err := s.buildAppActivityList(ctx, list)
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[*resp.AppCombinationActivityRespVO]{
		List:  vos,
		Total: total,
	}, nil
}

func (s *combinationActivityService) GetCombinationActivityDetail(ctx context.Context, id int64) (*resp.AppCombinationActivityDetailRespVO, error) {
	activity, err := s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.NewBizError(1001006000, "拼团活动不存在")
	}
	if activity.Status != 1 { // Enable
		return nil, errors.NewBizError(1001006001, "拼团活动已关闭")
	}

	prods, err := s.q.PromotionCombinationProduct.WithContext(ctx).Where(s.q.PromotionCombinationProduct.ActivityID.Eq(id)).Find()
	if err != nil {
		return nil, err
	}

	// Basic Info
	spu, err := s.spuSvc.GetSpu(ctx, activity.SpuID)
	if err != nil {
		return nil, err
	}

	minPrice := 0
	if len(prods) > 0 {
		minPrice = prods[0].CombinationPrice
		for _, p := range prods {
			if p.CombinationPrice < minPrice {
				minPrice = p.CombinationPrice
			}
		}
	}

	_ = minPrice // 用于后续扩展
	_ = spu      // SPU 信息用于后续扩展

	detailVo := &resp.AppCombinationActivityDetailRespVO{
		ID:               activity.ID,
		Name:             activity.Name,
		Status:           activity.Status,
		StartTime:        &activity.StartTime,
		EndTime:          &activity.EndTime,
		UserSize:         activity.UserSize,
		SuccessCount:     0, // TODO: 需要从 RecordService 获取
		SpuID:            activity.SpuID,
		TotalLimitCount:  activity.TotalLimitCount,
		SingleLimitCount: activity.SingleLimitCount,
		Products:         make([]resp.AppCombinationActivityDetailProduct, len(prods)),
	}

	for i, p := range prods {
		detailVo.Products[i] = resp.AppCombinationActivityDetailProduct{
			SkuID:            p.SkuID,
			CombinationPrice: p.CombinationPrice,
		}
	}

	return detailVo, nil
}

func (s *combinationActivityService) ValidateCombinationActivityCanJoin(ctx context.Context, activityID int64) (*promotion.PromotionCombinationActivity, error) {
	activity, err := s.q.PromotionCombinationActivity.WithContext(ctx).Where(s.q.PromotionCombinationActivity.ID.Eq(activityID)).First()
	if err != nil {
		return nil, errors.NewBizError(1001006000, "拼团活动不存在")
	}
	if activity.Status != 1 {
		return nil, errors.NewBizError(1001006001, "拼团活动已关闭")
	}
	now := time.Now()
	if now.Before(activity.StartTime) {
		return nil, errors.NewBizError(1001006002, "拼团活动未开始")
	}
	if now.After(activity.EndTime) {
		return nil, errors.NewBizError(1001006003, "拼团活动已结束")
	}
	return activity, nil
}

func (s *combinationActivityService) GetMatchCombinationActivityBySpuId(ctx context.Context, spuId int64) (*promotion.PromotionCombinationActivity, error) {
	now := time.Now()
	q := s.q.PromotionCombinationActivity
	return q.WithContext(ctx).
		Where(q.SpuID.Eq(spuId)).
		Where(q.Status.Eq(1)). // 1 = Enable
		Where(q.StartTime.Lt(now)).
		Where(q.EndTime.Gt(now)).
		First()
}

func (s *combinationActivityService) validateProducts(ctx context.Context, products []req.CombinationProductBaseVO) error {
	for _, p := range products {
		if _, err := s.spuSvc.GetSpu(ctx, p.SpuID); err != nil {
			return err
		}
		if _, err := s.skuSvc.GetSku(ctx, p.SkuID); err != nil {
			return err
		}
	}
	return nil
}

func (s *combinationActivityService) buildAppActivityList(ctx context.Context, list []*promotion.PromotionCombinationActivity) ([]*resp.AppCombinationActivityRespVO, error) {
	if len(list) == 0 {
		return []*resp.AppCombinationActivityRespVO{}, nil
	}
	spuIds := make([]int64, len(list))
	for i, item := range list {
		spuIds[i] = item.SpuID
	}
	spuList, err := s.spuSvc.GetSpuList(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := make(map[int64]*resp.ProductSpuResp)
	for _, spu := range spuList {
		spuMap[spu.ID] = spu
	}

	activityIds := make([]int64, len(list))
	for i, item := range list {
		activityIds[i] = item.ID
	}
	products, _ := s.q.PromotionCombinationProduct.WithContext(ctx).Where(s.q.PromotionCombinationProduct.ActivityID.In(activityIds...)).Find()
	priceMap := make(map[int64]int) // activityId -> minPrice
	for _, p := range products {
		min, ok := priceMap[p.ActivityID]
		if !ok || p.CombinationPrice < min {
			priceMap[p.ActivityID] = p.CombinationPrice
		}
	}

	result := make([]*resp.AppCombinationActivityRespVO, len(list))
	for i, item := range list {
		spu := spuMap[item.SpuID]
		if spu == nil {
			continue // Should not happen
		}
		result[i] = &resp.AppCombinationActivityRespVO{
			ID:               item.ID,
			Name:             item.Name,
			UserSize:         item.UserSize,
			SpuID:            item.SpuID,
			SpuName:          spu.Name,
			PicUrl:           spu.PicURL,
			MarketPrice:      spu.MarketPrice,
			CombinationPrice: priceMap[item.ID],
		}
	}
	return result, nil
}
