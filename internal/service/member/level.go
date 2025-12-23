package member

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	pkgErrors "github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"gorm.io/gorm"
)

type MemberLevelService struct {
	q *query.Query
}

func NewMemberLevelService(q *query.Query) *MemberLevelService {
	return &MemberLevelService{q: q}
}

// CreateLevel 创建等级
func (s *MemberLevelService) CreateLevel(ctx context.Context, r *req.MemberLevelCreateReq) (int64, error) {
	// Check Name Unique? Usually required.
	count, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.Name.Eq(r.Name)).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, pkgErrors.NewBizError(1004014000, "等级名称已存在")
	}
	// Check Level Value Unique?
	count, err = s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.Level.Eq(r.Level)).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, pkgErrors.NewBizError(1004014001, "等级值已存在")
	}

	level := &member.MemberLevel{
		Name:            r.Name,
		Level:           r.Level,
		Experience:      r.Experience,
		DiscountPercent: r.DiscountPercent,
		Icon:            r.Icon,
		BackgroundURL:   r.BackgroundURL,
		Status:          r.Status,
	}
	err = s.q.MemberLevel.WithContext(ctx).Create(level)
	return level.ID, err
}

// UpdateLevel 更新等级
func (s *MemberLevelService) UpdateLevel(ctx context.Context, r *req.MemberLevelUpdateReq) error {
	l, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(r.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.NewBizError(1004014002, "等级不存在")
		}
		return err
	}

	// Check Name conflict
	if l.Name != r.Name {
		count, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.Name.Eq(r.Name)).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return pkgErrors.NewBizError(1004014000, "等级名称已存在")
		}
	}
	// Check Level Value conflict
	if l.Level != r.Level {
		count, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.Level.Eq(r.Level)).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return pkgErrors.NewBizError(1004014001, "等级值已存在")
		}
	}

	_, err = s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(r.ID)).Updates(member.MemberLevel{
		Name:            r.Name,
		Level:           r.Level,
		Experience:      r.Experience,
		DiscountPercent: r.DiscountPercent,
		Icon:            r.Icon,
		BackgroundURL:   r.BackgroundURL,
		Status:          r.Status,
	})
	return err
}

// DeleteLevel 删除等级
func (s *MemberLevelService) DeleteLevel(ctx context.Context, id int64) error {
	// Check if users exist in this level? Java code usually checks.
	// For now, skip user check.
	_, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(id)).Delete()
	return err
}

// GetLevel 获得等级
func (s *MemberLevelService) GetLevel(ctx context.Context, id int64) (*member.MemberLevel, error) {
	return s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(id)).First()
}

// GetLevelSimpleList 获得开启的等级列表 (Simple)
func (s *MemberLevelService) GetLevelSimpleList(ctx context.Context) ([]*member.MemberLevel, error) {
	return s.q.MemberLevel.WithContext(ctx).
		Where(s.q.MemberLevel.Status.Eq(0)). // Enabled
		Order(s.q.MemberLevel.Level.Asc()).
		Find()
}

// GetLevelPage 获得等级分页
func (s *MemberLevelService) GetLevelPage(ctx context.Context, r *req.MemberLevelPageReq) (*pagination.PageResult[*member.MemberLevel], error) {
	q := s.q.MemberLevel.WithContext(ctx)
	if r.Name != "" {
		q = q.Where(s.q.MemberLevel.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.MemberLevel.Status.Eq(*r.Status))
	}
	list, total, err := q.Order(s.q.MemberLevel.Level.Asc()).FindByPage(r.GetOffset(), r.PageSize)
	return &pagination.PageResult[*member.MemberLevel]{
		List:  list,
		Total: total,
	}, err
}

// AddExperience 增加或减少经验
func (s *MemberLevelService) AddExperience(ctx context.Context, userId int64, experience int, bizType int, bizId string) error {
	if experience == 0 {
		return nil
	}
	// Logic to handle negative experience is implicit? Java handles absolute values based on bizType.
	// We assume `experience` here is the delta (positive or negative).

	// Transaction
	return s.q.Transaction(func(tx *query.Query) error {
		u := tx.MemberUser
		user, err := u.WithContext(ctx).Where(u.ID.Eq(userId)).First()
		if err != nil {
			return err
		}

		newExperience := int(user.Experience) + experience
		if newExperience < 0 {
			newExperience = 0 // Prevent negative experience
		}

		// 1. Update User Experience
		_, err = u.WithContext(ctx).Where(u.ID.Eq(userId)).Update(u.Experience, newExperience)
		if err != nil {
			return err
		}

		// 2. Check Level Upgrade
		newLevel, err := s.calculateNewLevel(ctx, newExperience)
		if err != nil {
			return err
		}

		if newLevel != nil && newLevel.ID != user.LevelID {
			// Level Changed
			// Update User Level
			_, err = u.WithContext(ctx).Where(u.ID.Eq(userId)).Update(u.LevelID, newLevel.ID)
			if err != nil {
				return err
			}

			// TODO: Log MemberLevelRecord
			// TODO: Notify Member
		}

		// TODO: Log MemberExperienceRecord

		return nil
	})
}

// calculateNewLevel 计算新等级
func (s *MemberLevelService) calculateNewLevel(ctx context.Context, experience int) (*member.MemberLevel, error) {
	// Get all enabled levels sorted by level value
	levels, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.Status.Eq(0)).Order(s.q.MemberLevel.Level.Desc()).Find()
	if err != nil {
		return nil, err
	}

	for _, l := range levels {
		if experience >= int(l.Experience) {
			return l, nil
		}
	}
	return nil, nil // No level matched (maybe experience is too low for any level, or keep current? assume lowest level is 0 exp?)
}

// UpdateUserLevel Admin 更新用户等级
func (s *MemberLevelService) UpdateUserLevel(ctx context.Context, userId int64, levelId *int64, reason string) error {
	u := s.q.MemberUser
	user, err := u.WithContext(ctx).Where(u.ID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.NewBizError(1004014003, "用户不存在")
		}
		return err
	}

	var newLevelId int64 = 0
	if levelId != nil {
		// 校验等级是否存在
		_, err := s.GetLevel(ctx, *levelId)
		if err != nil {
			return pkgErrors.NewBizError(1004014002, "等级不存在")
		}
		newLevelId = *levelId
	}

	// 如果等级没变化，直接返回
	if user.LevelID == newLevelId {
		return nil
	}

	// 更新用户等级
	_, err = u.WithContext(ctx).Where(u.ID.Eq(userId)).Update(u.LevelID, newLevelId)
	return err
}

// GetLevelListByIds 根据 ID 列表获得等级列表
func (s *MemberLevelService) GetLevelListByIds(ctx context.Context, ids []int64) ([]*member.MemberLevel, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.In(ids...)).Find()
}
