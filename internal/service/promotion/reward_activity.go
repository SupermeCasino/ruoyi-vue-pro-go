package promotion

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/types"

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
		Status:             consts.CommonStatusEnable, // Default Open
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
	// 解析 productScopeValues（兼容逗号分隔和 JSON 格式）
	scopeValues, _ := types.ParseListFromCSV[int64](item.ProductScopeValues)

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
		Where(s.q.PromotionRewardActivity.Status.Eq(consts.CommonStatusEnable)). // Open
		Where(s.q.PromotionRewardActivity.StartTime.Lt(now)).                    // Started
		Where(s.q.PromotionRewardActivity.EndTime.Gt(now)).                      // Not Ended
		Order(s.q.PromotionRewardActivity.ID.Desc()).Find()                      // High Priority First (ID as fallback)
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

		// 解析 productScopeValues（兼容逗号分隔和 JSON 格式）
		scopeValues, _ := types.ParseListFromCSV[int64](activity.ProductScopeValues)

		for _, item := range items {
			if skuTaken[item.SkuID] {
				continue
			}

			// Check Scope
			isMatch := false
			switch activity.ProductScope {
			case consts.ProductScopeAll: // 全部商品
				isMatch = true
			case consts.ProductScopeSpu: // 指定商品
				if lo.Contains(scopeValues, item.SpuID) {
					isMatch = true
				}
			case consts.ProductScopeCategory: // 指定品类
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
			switch activity.ConditionType {
			case consts.ConditionTypePrice: // 满金额
				if matchedPrice >= rule.Limit {
					// Use the biggest valid limit? Usually rules are ascending or descending.
					// Assume simplified: Find MAX satisfied limit.
					// If rules are not sorted, need to find max.
					if rule.DiscountPrice > discount { // Simple assumption: bigger reduction is better
						discount = rule.DiscountPrice
					}
				}
			case consts.ConditionTypeCount: // 满数量
				if matchedCount >= rule.Limit {
					if rule.DiscountPrice > discount {
						discount = rule.DiscountPrice
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
		q.Status.Eq(consts.CommonStatusEnable), // ENABLE
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
	// 解析 productScopeValues（兼容逗号分隔和 JSON 格式）
	scopeValues, _ := types.ParseListFromCSV[int64](activity.ProductScopeValues)
	if activity.ProductScope == consts.ProductScopeSpu {
		return lo.Contains(scopeValues, spuID)
	}
	return false
}

// GetRewardActivityForApp 获得满减送活动（App 端）
// Java: AppRewardActivityController#getRewardActivity
func (s *RewardActivityService) GetRewardActivityForApp(ctx context.Context, id int64) (*resp.AppRewardActivityResp, error) {
	activity, err := s.q.PromotionRewardActivity.WithContext(ctx).Where(s.q.PromotionRewardActivity.ID.Eq(id)).First()
	if err != nil {
		return nil, nil // 活动不存在返回 null
	}

	// 解析规则
	var rules []resp.Rule
	_ = json.Unmarshal([]byte(activity.Rules), &rules)

	// 解析 productScopeValues（兼容逗号分隔和 JSON 格式）
	scopeValues, _ := types.ParseListFromCSV[int64](activity.ProductScopeValues)

	// 构建响应，包含规则描述
	appRules := make([]resp.AppRewardActivityRule, 0, len(rules))
	for _, rule := range rules {
		appRule := resp.AppRewardActivityRule{
			Limit:                    rule.Limit,
			DiscountPrice:            rule.DiscountPrice,
			FreeDelivery:             rule.FreeDelivery,
			Point:                    rule.Point,
			GiveCouponTemplateCounts: rule.GiveCouponTemplateCounts,
			Description:              s.GetRewardActivityRuleDescription(activity.ConditionType, &rule),
		}
		appRules = append(appRules, appRule)
	}

	return &resp.AppRewardActivityResp{
		ID:                 activity.ID,
		Status:             activity.Status,
		Name:               activity.Name,
		StartTime:          activity.StartTime.UnixMilli(), // 转换为毫秒时间戳
		EndTime:            activity.EndTime.UnixMilli(),   // 转换为毫秒时间戳
		ConditionType:      activity.ConditionType,
		ProductScope:       activity.ProductScope,
		ProductScopeValues: scopeValues,
		Rules:              appRules,
	}, nil
}

// GetRewardActivityRuleDescription 获取满减送活动规则描述
// Java: RewardActivityService#getRewardActivityRuleDescription
func (s *RewardActivityService) GetRewardActivityRuleDescription(conditionType int, rule *resp.Rule) string {
	description := ""

	// 构建条件描述
	if conditionType == consts.ConditionTypePrice {
		// 满 N 元（带空格和小数点）
		description = "满 " + formatMoneyWithDecimal(rule.Limit) + " 元"
	} else {
		// 满 N 件
		description = "满 " + formatInt(rule.Limit) + " 件"
	}

	// 构建优惠描述
	tips := make([]string, 0, 4)
	if rule.DiscountPrice > 0 {
		tips = append(tips, "减 "+formatMoneyWithDecimal(rule.DiscountPrice))
	}
	if rule.FreeDelivery {
		tips = append(tips, "包邮")
	}
	if rule.Point > 0 {
		tips = append(tips, "送 "+formatInt(rule.Point)+" 积分")
	}
	if len(rule.GiveCouponTemplateCounts) > 0 {
		totalCoupons := 0
		for _, count := range rule.GiveCouponTemplateCounts {
			totalCoupons += count
		}
		tips = append(tips, "送 "+formatInt(totalCoupons)+" 张优惠券")
	}

	if len(tips) > 0 {
		description = description + joinStrings(tips, "、")
	}

	return description
}

// formatMoneyWithDecimal 格式化金额（分转元，带小数点）
func formatMoneyWithDecimal(fen int) string {
	yuan := float64(fen) / 100.0
	return strconv.FormatFloat(yuan, 'f', 2, 64)
}

// formatInt 格式化整数
func formatInt(n int) string {
	if n < 0 {
		return "0"
	}
	return strconv.Itoa(n)
}

// formatFloat 格式化浮点数
func formatFloat(f float64) string {
	if f < 0 {
		return "0"
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// joinStrings 连接字符串数组
func joinStrings(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
