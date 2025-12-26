package promotion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type RewardActivityService struct {
	q *query.Query
}

func NewRewardActivityService(q *query.Query) *RewardActivityService {
	return &RewardActivityService{q: q}
}

// CreateRewardActivity 创建活动
func (s *RewardActivityService) CreateRewardActivity(ctx context.Context, r *req.PromotionRewardActivityCreateReq) (int64, error) {
	rules, _ := json.Marshal(r.Rules)
	scopeValues, _ := json.Marshal(r.ProductScopeValues)

	activity := &promotion.PromotionRewardActivity{
		Name:               r.Name,
		Status:             0, // Default Open
		StartTime:          r.StartTime,
		EndTime:            r.EndTime,
		ProductScope:       r.ProductScope,
		ProductScopeValues: string(scopeValues),
		ConditionType:      r.ConditionType,
		Rules:              string(rules),
		Remark:             r.Remark,
	}
	err := s.q.PromotionRewardActivity.WithContext(ctx).Create(activity)
	return activity.ID, err
}

// UpdateRewardActivity 更新活动
func (s *RewardActivityService) UpdateRewardActivity(ctx context.Context, r *req.PromotionRewardActivityUpdateReq) error {
	_, err := s.q.PromotionRewardActivity.WithContext(ctx).Where(s.q.PromotionRewardActivity.ID.Eq(r.ID)).First()
	if err != nil {
		return errors.NewBizError(1004002000, "活动不存在")
	}

	rules, _ := json.Marshal(r.Rules)
	scopeValues, _ := json.Marshal(r.ProductScopeValues)

	_, err = s.q.PromotionRewardActivity.WithContext(ctx).Where(s.q.PromotionRewardActivity.ID.Eq(r.ID)).Updates(promotion.PromotionRewardActivity{
		Name:               r.Name,
		StartTime:          r.StartTime,
		EndTime:            r.EndTime,
		ProductScope:       r.ProductScope,
		ProductScopeValues: string(scopeValues),
		ConditionType:      r.ConditionType,
		Rules:              string(rules),
		Remark:             r.Remark,
	})
	return err
}

// DeleteRewardActivity 删除活动
func (s *RewardActivityService) DeleteRewardActivity(ctx context.Context, id int64) error {
	_, err := s.q.PromotionRewardActivity.WithContext(ctx).Where(s.q.PromotionRewardActivity.ID.Eq(id)).Delete()
	return err
}

// CloseRewardActivity 关闭活动
// Java: RewardActivityServiceImpl#closeRewardActivity
func (s *RewardActivityService) CloseRewardActivity(ctx context.Context, id int64) error {
	// 1. 校验存在
	activity, err := s.q.PromotionRewardActivity.WithContext(ctx).Where(s.q.PromotionRewardActivity.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1004002000, "活动不存在")
	}

	// 2. 检查状态：已关闭的活动不能再关闭
	if activity.Status == consts.CommonStatusEnable { // 已启用，需要关闭
		return errors.NewBizError(1004002003, "活动已关闭，不能重复关闭")
	}

	// 3. 更新状态为关闭
	_, err = s.q.PromotionRewardActivity.WithContext(ctx).
		Where(s.q.PromotionRewardActivity.ID.Eq(id)).
		Update(s.q.PromotionRewardActivity.Status, consts.CommonStatusDisable) // 禁用
	return err
}

// GetRewardActivity 获得活动详情
func (s *RewardActivityService) GetRewardActivity(ctx context.Context, id int64) (*resp.PromotionRewardActivityResp, error) {
	item, err := s.q.PromotionRewardActivity.WithContext(ctx).Where(s.q.PromotionRewardActivity.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetRewardActivityPage 获得活动分页
func (s *RewardActivityService) GetRewardActivityPage(ctx context.Context, r *req.PromotionRewardActivityPageReq) (*pagination.PageResult[*resp.PromotionRewardActivityResp], error) {
	q := s.q.PromotionRewardActivity.WithContext(ctx)
	if r.Name != "" {
		q = q.Where(s.q.PromotionRewardActivity.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.PromotionRewardActivity.Status.Eq(*r.Status))
	}
	if len(r.CreateTime) == 2 && r.CreateTime[0] != nil && r.CreateTime[1] != nil {
		q = q.Where(s.q.PromotionRewardActivity.CreateTime.Between(*r.CreateTime[0], *r.CreateTime[1]))
	}

	list, total, err := q.Order(s.q.PromotionRewardActivity.ID.Desc()).FindByPage(r.GetOffset(), r.PageSize)
	if err != nil {
		return nil, err
	}

	resList := lo.Map(list, func(item *promotion.PromotionRewardActivity, _ int) *resp.PromotionRewardActivityResp {
		return s.convertResp(item)
	})

	return &pagination.PageResult[*resp.PromotionRewardActivityResp]{
		List:  resList,
		Total: total,
	}, nil
}

func (s *RewardActivityService) convertResp(item *promotion.PromotionRewardActivity) *resp.PromotionRewardActivityResp {
	var rules []resp.Rule
	_ = json.Unmarshal([]byte(item.Rules), &rules)
	var scopeValues []int64
	_ = json.Unmarshal([]byte(item.ProductScopeValues), &scopeValues)

	return &resp.PromotionRewardActivityResp{
		ID:                 item.ID,
		Name:               item.Name,
		Status:             item.Status,
		StartTime:          item.StartTime,
		EndTime:            item.EndTime,
		ProductScope:       item.ProductScope,
		ProductScopeValues: scopeValues,
		ConditionType:      item.ConditionType,
		Rules:              rules,
		Remark:             item.Remark,
		CreateTime:         item.CreateTime,
	}
}

type ActivityMatchItem struct {
	SkuID      int64
	SpuID      int64
	CategoryID int64
	Price      int
	Count      int
}

type ActivityMatchResult struct {
	TotalPrice    int
	TotalDiscount int
	ActivityID    int64
	ActivityName  string
	SkuIDs        []int64
}

// CalculateRewardActivity 计算满减送活动优惠
func (s *RewardActivityService) CalculateRewardActivity(ctx context.Context, items []ActivityMatchItem) (int, []ActivityMatchResult, error) {
	// 1. Fetch All Active Activities
	now := time.Now()
	activities, err := s.q.PromotionRewardActivity.WithContext(ctx).
		Where(s.q.PromotionRewardActivity.Status.Eq(0)).      // Open
		Where(s.q.PromotionRewardActivity.StartTime.Lt(now)). // Started
		Where(s.q.PromotionRewardActivity.EndTime.Gt(now)).   // Not Ended
		Order(s.q.PromotionRewardActivity.ID.Desc()).Find()   // High Priority First (ID as fallback)
	if err != nil {
		return 0, nil, err
	}

	if len(activities) == 0 {
		return 0, nil, nil
	}

	activityMap := make(map[int64]*promotion.PromotionRewardActivity)
	for _, a := range activities {
		activityMap[a.ID] = a
	}

	// 2. Iterate Activities and Match Items
	// Strategy: High priority activity grabs items.
	// Map SkuID -> Taken(bool)
	skuTaken := make(map[int64]bool)
	var results []ActivityMatchResult
	totalDiscount := 0

	for _, activity := range activities {
		var matchedItems []ActivityMatchItem
		var matchedPrice int
		var matchedCount int

		var scopeValues []int64
		_ = json.Unmarshal([]byte(activity.ProductScopeValues), &scopeValues)

		for _, item := range items {
			if skuTaken[item.SkuID] {
				continue
			}

			// Check Scope
			isMatch := false
			if activity.ProductScope == consts.ProductScopeAll { // 全部商品
				isMatch = true
			} else if activity.ProductScope == consts.ProductScopeSpu { // 指定商品
				if lo.Contains(scopeValues, item.SpuID) {
					isMatch = true
				}
			} else if activity.ProductScope == consts.ProductScopeCategory { // 指定品类
				if lo.Contains(scopeValues, item.CategoryID) {
					isMatch = true
				}
			}

			if isMatch {
				matchedItems = append(matchedItems, item)
				matchedPrice += item.Price * item.Count
				matchedCount += item.Count
			}
		}

		if len(matchedItems) == 0 {
			continue
		}

		// 3. Check Rules
		var rules []resp.Rule
		_ = json.Unmarshal([]byte(activity.Rules), &rules)
		// Find best rule: Usually sorted or check all.
		// Rule Limit: Price or Count.
		discount := 0
		for _, rule := range rules {
			if activity.ConditionType == consts.ConditionTypePrice { // 满金额
				if matchedPrice >= rule.Limit {
					// Use the biggest valid limit? Usually rules are ascending or descending.
					// Assume simplified: Find MAX satisfied limit.
					// If rules are not sorted, need to find max.
					if rule.ReducePrice > discount { // Simple assumption: bigger reduction is better
						discount = rule.ReducePrice
					}
				}
			} else if activity.ConditionType == consts.ConditionTypeCount { // 满数量
				if matchedCount >= rule.Limit {
					if rule.ReducePrice > discount {
						discount = rule.ReducePrice
					}
				}
			}
		}

		if discount > 0 {
			// Apply Activity
			// Cap discount at total price?
			if discount > matchedPrice {
				discount = matchedPrice
			}
			totalDiscount += discount

			// Mark items as taken
			var skuIDs []int64
			for _, item := range matchedItems {
				skuTaken[item.SkuID] = true
				skuIDs = append(skuIDs, item.SkuID)
			}

			results = append(results, ActivityMatchResult{
				TotalPrice:    matchedPrice,
				TotalDiscount: discount,
				ActivityID:    activity.ID,
				ActivityName:  activity.Name,
				SkuIDs:        skuIDs,
			})
		}
	}

	return totalDiscount, results, nil
}

// GetRewardActivityMapBySpuIds 获得指定 SPU 编号数组的满减送活动 Map
func (s *RewardActivityService) GetRewardActivityMapBySpuIds(ctx context.Context, spuIDs []int64) (map[int64]*promotion.PromotionRewardActivity, error) {
	if len(spuIDs) == 0 {
		return map[int64]*promotion.PromotionRewardActivity{}, nil
	}
	now := time.Now()
	q := s.q.PromotionRewardActivity
	activities, err := q.WithContext(ctx).Where(
		q.Status.Eq(0), // ENABLE
		q.StartTime.Lte(now),
		q.EndTime.Gte(now),
	).Find()
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*promotion.PromotionRewardActivity)
	for _, activity := range activities {
		for _, spuID := range spuIDs {
			if _, ok := res[spuID]; ok {
				continue
			}
			if s.isSpuMatchActivity(activity, spuID) {
				res[spuID] = activity
			}
		}
	}
	return res, nil
}

func (s *RewardActivityService) isSpuMatchActivity(activity *promotion.PromotionRewardActivity, spuID int64) bool {
	if activity.ProductScope == consts.ProductScopeAll {
		return true
	}
	var scopeValues []int64
	_ = json.Unmarshal([]byte(activity.ProductScopeValues), &scopeValues)
	if activity.ProductScope == consts.ProductScopeSpu {
		return lo.Contains(scopeValues, spuID)
	}
	return false
}
