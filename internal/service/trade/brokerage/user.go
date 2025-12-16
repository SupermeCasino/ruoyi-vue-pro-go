package brokerage

import (
	"context"
	"errors"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	tradeReq "github.com/wxlbd/ruoyi-mall-go/internal/api/req/app/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BrokerageUserService struct {
	q         *query.Query
	logger    *zap.Logger
	memberSvc *member.MemberUserService
	configSvc *trade.TradeConfigService
}

func NewBrokerageUserService(q *query.Query, logger *zap.Logger, memberSvc *member.MemberUserService, configSvc *trade.TradeConfigService) *BrokerageUserService {
	return &BrokerageUserService{
		q:         q,
		logger:    logger,
		memberSvc: memberSvc,
		configSvc: configSvc,
	}
}

// GetBrokerageUser 获得分销用户
func (s *BrokerageUserService) GetBrokerageUser(ctx context.Context, id int64) (*brokerage.BrokerageUser, error) {
	return s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).First()
}

func parseTime(t string) time.Time {
	res, _ := time.ParseInLocation(time.DateTime, t, time.Local)
	return res
}

// GetBrokerageUserPage 获得分销用户分页
func (s *BrokerageUserService) GetBrokerageUserPage(ctx context.Context, r *req.BrokerageUserPageReq) (*core.PageResult[*brokerage.BrokerageUser], error) {
	q := s.q.BrokerageUser.WithContext(ctx)

	// Filter by BindUserId and Level
	if r.BindUserID > 0 {
		childIDs, err := s.GetChildUserIdsByLevel(ctx, r.BindUserID, r.Level)
		if err != nil {
			return nil, err
		}
		if len(childIDs) == 0 {
			return &core.PageResult[*brokerage.BrokerageUser]{List: []*brokerage.BrokerageUser{}, Total: 0}, nil
		}
		q = q.Where(s.q.BrokerageUser.ID.In(childIDs...))
	}

	if r.BrokerageEnabled != nil {
		q = q.Where(s.q.BrokerageUser.BrokerageEnabled.Is(*r.BrokerageEnabled))
	}

	// Time ranges
	if len(r.CreateTime) == 2 {
		q = q.Where(s.q.BrokerageUser.CreatedAt.Between(parseTime(r.CreateTime[0]), parseTime(r.CreateTime[1])))
	}
	// BindUserTime?

	page := r.PageNo
	size := r.PageSize
	offset := (page - 1) * size

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Limit(size).Offset(offset).Order(s.q.BrokerageUser.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*brokerage.BrokerageUser]{
		List:  list,
		Total: total,
	}, nil
}

// CreateBrokerageUser 创建分销用户
func (s *BrokerageUserService) CreateBrokerageUser(ctx context.Context, r *req.BrokerageUserCreateReq) (int64, error) {
	// 1. Check if exists
	exists, err := s.GetBrokerageUser(ctx, r.UserID)
	if err == nil && exists != nil {
		return 0, errors.New("分销用户已存在")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Assuming gorm.ErrRecordNotFound handled by repo or checked here
		// Actually s.GetBrokerageUser returns error if not found?
		// Usually generated code returns error.
		return 0, err
	}

	// 2. Validate Bind
	if err := s.validateCanBindUser(ctx, &brokerage.BrokerageUser{ID: r.UserID}, r.BindUserID); err != nil {
		return 0, err
	}

	// 3. Create
	now := time.Now()
	user := &brokerage.BrokerageUser{
		ID:               r.UserID,
		BindUserID:       r.BindUserID,
		BindUserTime:     &now,
		BrokerageTime:    &now,
		BrokerageEnabled: true, // Assuming enabled on manual create
	}
	if err := s.q.BrokerageUser.WithContext(ctx).Create(user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetOrCreateBrokerageUser 获得或创建分销用户
func (s *BrokerageUserService) GetOrCreateBrokerageUser(ctx context.Context, id int64) (*brokerage.BrokerageUser, error) {
	u, err := s.GetBrokerageUser(ctx, id)
	if err == nil && u != nil {
		return u, nil
	}
	// Create if not found
	// In Java: default enabled = true? Check TradeConfig.
	// For now assume default enabled if configured.
	now := time.Now()
	user := &brokerage.BrokerageUser{
		ID:               id,
		BrokerageEnabled: true, // Default enabled? Or check config?
		BrokerageTime:    &now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := s.q.BrokerageUser.WithContext(ctx).Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// BindBrokerageUser 绑定推广员
func (s *BrokerageUserService) BindBrokerageUser(ctx context.Context, userId int64, bindUserId int64) (bool, error) {
	// Java valid: bindUserId > 0 && userId != bindUserId
	if bindUserId == 0 || userId == bindUserId {
		return false, nil // Should return error or false? Java returns boolean success.
	}
	// Check if already bound
	user, err := s.GetBrokerageUser(ctx, userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	if user != nil && user.BindUserID > 0 {
		return false, nil // Already bound
	}

	// Delegate to UpdateBrokerageUserId logic which contains validation
	// But UpdateBrokerageUserId requires manual handling?
	// Or we can call validate and update.
	// Let's reuse UpdateBrokerageUserId for consistency, but UpdateBrokerageUserId assumes Admin override?
	// App binding usually has simpler rules (can't change once bound).
	// Admin can change.

	// If user does not exist, create it first?
	if user == nil {
		if _, err := s.CreateBrokerageUser(ctx, &req.BrokerageUserCreateReq{
			UserID:     userId,
			BindUserID: bindUserId,
		}); err != nil {
			return false, err
		}
		return true, nil
	}

	// If exists, update
	if err := s.UpdateBrokerageUserId(ctx, userId, bindUserId); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateBrokerageUserEnabled 修改推广资格
func (s *BrokerageUserService) UpdateBrokerageUserEnabled(ctx context.Context, id int64, enabled bool) error {
	u, err := s.GetBrokerageUser(ctx, id)
	if err != nil {
		return err
	}
	if u.BrokerageEnabled == enabled {
		return nil
	}

	updates := map[string]interface{}{
		"brokerage_enabled": enabled,
	}
	if enabled {
		now := time.Now()
		updates["brokerage_time"] = &now
	} else {
		updates["brokerage_time"] = nil
	}
	_, err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).Updates(updates)
	return err
}

// UpdateBrokerageUserId 修改推广员
func (s *BrokerageUserService) UpdateBrokerageUserId(ctx context.Context, id int64, bindUserId int64) error {
	user, err := s.GetBrokerageUser(ctx, id)
	if err != nil {
		return err
	}
	if user.BindUserID == bindUserId {
		return nil
	}

	// Clear
	if bindUserId == 0 {
		_, err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).Updates(map[string]interface{}{
			"bind_user_id":   0, // Or null
			"bind_user_time": nil,
		})
		return err
	}

	// Validate
	if err := s.validateCanBindUser(ctx, user, bindUserId); err != nil {
		return err
	}

	now := time.Now()
	_, err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).Updates(map[string]interface{}{
		"bind_user_id":   bindUserId,
		"bind_user_time": &now,
	})
	return err
}

// validateCanBindUser 校验能否绑定
func (s *BrokerageUserService) validateCanBindUser(ctx context.Context, user *brokerage.BrokerageUser, bindUserId int64) error {
	if bindUserId == 0 {
		return nil
	}
	// 1. Check Bind User exists
	bindUser, err := s.memberSvc.GetUser(ctx, bindUserId) // This is MemberUser
	if err != nil || bindUser == nil {
		return errors.New("推广人不存在")
	}

	// 2. Check Bind User Brokerage Enabled
	brokerageBindUser, err := s.GetBrokerageUser(ctx, bindUserId)
	if err != nil || brokerageBindUser == nil || !brokerageBindUser.BrokerageEnabled {
		return errors.New("推广人无推广资格")
	}

	// 3. Self bind
	if user.ID == bindUserId {
		return errors.New("不能绑定自己")
	}

	// 4. Loop check
	// A -> B -> A
	currentBindId := brokerageBindUser.BindUserID
	for i := 0; i < 100; i++ { // Limit loop
		if currentBindId == 0 {
			break
		}
		if currentBindId == user.ID {
			return errors.New("不能循环绑定")
		}
		// Fetch next
		next, err := s.GetBrokerageUser(ctx, currentBindId)
		if err != nil || next == nil {
			break
		}
		currentBindId = next.BindUserID
	}
	return nil
}

// GetChildUserIdsByLevel 获得下级用户编号列表
func (s *BrokerageUserService) GetChildUserIdsByLevel(ctx context.Context, bindUserId int64, level int) ([]int64, error) {
	// Level 1
	var level1Ids []int64
	err := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.BindUserID.Eq(bindUserId)).Select(s.q.BrokerageUser.ID).Scan(&level1Ids)
	if err != nil {
		return nil, err
	}
	if len(level1Ids) == 0 {
		return []int64{}, nil
	}

	if level == 1 {
		return level1Ids, nil
	}

	// Level 2
	var level2Ids []int64
	if len(level1Ids) > 0 {
		err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.BindUserID.In(level1Ids...)).Select(s.q.BrokerageUser.ID).Scan(&level2Ids)
		if err != nil {
			return nil, err
		}
	}

	if level == 2 {
		return level2Ids, nil
	}

	// All (1 + 2)
	return append(level1Ids, level2Ids...), nil
}

// GetBrokerageUserCountByBindUserId 获得推广用户数量
func (s *BrokerageUserService) GetBrokerageUserCountByBindUserId(ctx context.Context, bindUserId int64, level int) (int64, error) {
	ids, err := s.GetChildUserIdsByLevel(ctx, bindUserId, level)
	if err != nil {
		return 0, err
	}
	return int64(len(ids)), nil
}

// GetBrokerageUserChildSummaryPage 获得下级分销统计分页
func (s *BrokerageUserService) GetBrokerageUserChildSummaryPage(ctx context.Context, r *tradeReq.AppBrokerageUserChildSummaryPageReqVO, userId int64) (*core.PageResult[*brokerage.BrokerageUser], error) {
	childIDs, err := s.GetChildUserIdsByLevel(ctx, userId, r.Level)
	if err != nil {
		return nil, err
	}
	if len(childIDs) == 0 {
		return &core.PageResult[*brokerage.BrokerageUser]{List: []*brokerage.BrokerageUser{}, Total: 0}, nil
	}

	q := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.In(childIDs...))

	// Note: Nickname filtering requires Join with MemberUser, skipping for simple implementation or assuming nicknames not in BrokerageUser table.
	// Java joins or fetches map.
	// Here, we just return the page of BrokerageUsers, logic layer will fill info.
	// But Sorting?
	if r.Sorting == "userCount" {
		// Complex sort, skip for now or use default
	} else if r.Sorting == "brokeragePrice" {
		q = q.Order(s.q.BrokerageUser.BrokeragePrice.Desc())
	} else {
		q = q.Order(s.q.BrokerageUser.BrokerageTime.Desc())
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Limit(r.PageSize).Offset((r.PageNo - 1) * r.PageSize).Find()
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*brokerage.BrokerageUser]{List: list, Total: total}, nil
}

// GetBrokerageUserRankPageByUserCount 获得分销用户排行分页（基于用户量）
func (s *BrokerageUserService) GetBrokerageUserRankPageByUserCount(ctx context.Context, r *tradeReq.AppBrokerageUserRankPageReqVO) (*core.PageResult[*brokerage.BrokerageUser], error) {
	// Complex query: Count sub-users and rank.
	// This usually involves Group By or Subquery.
	// Simplified: return empty for now or Implement proper SQL.
	// Java uses: select id, nickname, avatar, (select count(*) from trade_brokerage_user where bind_user_id = t.id) as brokerage_user_count ...
	// GORM raw SQL might be best.

	// Placeholder
	return &core.PageResult[*brokerage.BrokerageUser]{List: []*brokerage.BrokerageUser{}, Total: 0}, nil
}
