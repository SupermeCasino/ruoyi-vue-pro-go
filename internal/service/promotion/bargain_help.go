package promotion

import (
	"context"

	"backend-go/internal/api/req"
	"backend-go/internal/model/promotion"
	"backend-go/internal/repo/query"
)

type BargainHelpService struct {
	q *query.Query
}

func NewBargainHelpService(q *query.Query) *BargainHelpService {
	return &BargainHelpService{q: q}
}

// GetBargainHelpUserCountMapByActivity 获得砍价活动的助力用户数量 Map
func (s *BargainHelpService) GetBargainHelpUserCountMapByActivity(ctx context.Context, activityIds []int64) (map[int64]int, error) {
	if len(activityIds) == 0 {
		return make(map[int64]int), nil
	}
	q := s.q.PromotionBargainHelp
	rows, err := q.WithContext(ctx).
		Select(q.ActivityID, q.UserID.Count()).
		Where(q.ActivityID.In(activityIds...)).
		Group(q.ActivityID).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]int)
	for rows.Next() {
		var activityID int64
		var count int
		if err := rows.Scan(&activityID, &count); err == nil {
			result[activityID] = count
		}
	}
	return result, nil
}

// GetBargainHelpList 获得砍价助力列表 (按 recordId)
func (s *BargainHelpService) GetBargainHelpList(ctx context.Context, recordID int64) ([]*promotion.PromotionBargainHelp, error) {
	q := s.q.PromotionBargainHelp
	return q.WithContext(ctx).Where(q.RecordID.Eq(recordID)).Order(q.CreatedAt.Desc()).Find()
}

// GetBargainHelp 获得指定记录和用户的助力记录
func (s *BargainHelpService) GetBargainHelp(ctx context.Context, recordID int64, userID int64) (*promotion.PromotionBargainHelp, error) {
	q := s.q.PromotionBargainHelp
	return q.WithContext(ctx).Where(q.RecordID.Eq(recordID), q.UserID.Eq(userID)).First()
}

// GetBargainHelpCountByActivity 获得用户在指定活动的助力次数
func (s *BargainHelpService) GetBargainHelpCountByActivity(ctx context.Context, activityID int64, userID int64) (int64, error) {
	q := s.q.PromotionBargainHelp
	return q.WithContext(ctx).Where(q.ActivityID.Eq(activityID), q.UserID.Eq(userID)).Count()
}

// CreateBargainHelp 砍价助力
func (s *BargainHelpService) CreateBargainHelp(ctx context.Context, userID int64, r *req.AppBargainHelpCreateReq) (*promotion.PromotionBargainHelp, error) {
	// TODO: Implement full business logic
	// 1. 校验砍价记录存在
	// 2. 校验用户不能帮自己砍
	// 3. 校验用户是否已经助力过
	// 4. 校验活动助力次数限制
	// 5. 随机计算砍价金额
	// 6. 创建助力记录
	// 7. 更新砍价记录价格
	return nil, nil
}
