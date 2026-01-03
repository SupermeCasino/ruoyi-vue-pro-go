package member

import (
	"context"
	"errors"

	member2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	pkgErrors "github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"gorm.io/gorm"
)

// 会员等级相关错误码 (对齐 Java: ErrorCodeConstants)
const (
	ErrorCodeLevelNotExists     = 1004011000
	ErrorCodeLevelNameExists    = 1004011001
	ErrorCodeLevelValueExists   = 1004011002
	ErrorCodeLevelExperienceMin = 1004011003
	ErrorCodeLevelExperienceMax = 1004011004
	ErrorCodeLevelHasUser       = 1004011005
)

// MemberLevelService 会员等级 Service
type MemberLevelService struct {
	q *query.Query
}

func NewMemberLevelService(q *query.Query) *MemberLevelService {
	return &MemberLevelService{q: q}
}

// CreateLevel 创建会员等级
func (s *MemberLevelService) CreateLevel(ctx context.Context, r *member2.MemberLevelCreateReq) (int64, error) {
	// 校验配置是否有效
	if err := s.validateConfigValid(ctx, 0, r.Name, r.Level, r.Experience); err != nil {
		return 0, err
	}

	// 插入
	level := &member.MemberLevel{
		Name:            r.Name,
		Level:           r.Level,
		Experience:      r.Experience,
		DiscountPercent: r.DiscountPercent,
		Icon:            r.Icon,
		BackgroundURL:   r.BackgroundURL,
		Status:          r.Status,
	}
	err := s.q.MemberLevel.WithContext(ctx).Create(level)
	return level.ID, err
}

// UpdateLevel 更新会员等级
func (s *MemberLevelService) UpdateLevel(ctx context.Context, r *member2.MemberLevelUpdateReq) error {
	// 校验存在
	if _, err := s.validateLevelExists(ctx, r.ID); err != nil {
		return err
	}
	// 校验配置是否有效
	if err := s.validateConfigValid(ctx, r.ID, r.Name, r.Level, r.Experience); err != nil {
		return err
	}

	// 更新
	_, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(r.ID)).Updates(member.MemberLevel{
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

// DeleteLevel 删除会员等级
func (s *MemberLevelService) DeleteLevel(ctx context.Context, id int64) error {
	// 校验存在
	if _, err := s.validateLevelExists(ctx, id); err != nil {
		return err
	}
	// 校验分组下是否有用户
	if err := s.validateLevelHasUser(ctx, id); err != nil {
		return err
	}
	// 删除
	_, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(id)).Delete()
	return err
}

func (s *MemberLevelService) validateLevelExists(ctx context.Context, id int64) (*member.MemberLevel, error) {
	level, err := s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewBizError(ErrorCodeLevelNotExists, "用户等级不存在")
		}
		return nil, err
	}
	return level, nil
}

func (s *MemberLevelService) validateNameUnique(list []*member.MemberLevel, id int64, name string) error {
	for _, l := range list {
		if l.Name != name {
			continue
		}
		if id == 0 || id != l.ID {
			return pkgErrors.NewBizError(ErrorCodeLevelNameExists, "用户等级名称已经被使用")
		}
	}
	return nil
}

func (s *MemberLevelService) validateLevelUnique(list []*member.MemberLevel, id int64, level int) error {
	for _, l := range list {
		if l.Level != level {
			continue
		}
		if id == 0 || id != l.ID {
			return pkgErrors.NewBizError(ErrorCodeLevelValueExists, "用户等级值已被使用")
		}
	}
	return nil
}

func (s *MemberLevelService) validateExperienceOutRange(list []*member.MemberLevel, id int64, level int, experience int) error {
	for _, l := range list {
		if l.ID == id {
			continue
		}
		if l.Level < level {
			// 经验大于前一个等级
			if experience <= l.Experience {
				return pkgErrors.NewBizError(ErrorCodeLevelExperienceMin, "升级经验必须大于上一个等级设置的升级经验")
			}
		} else if l.Level > level {
			// 小于下一个级别
			if experience >= l.Experience {
				return pkgErrors.NewBizError(ErrorCodeLevelExperienceMax, "升级经验必须小于下一个等级设置的升级经验")
			}
		}
	}
	return nil
}

func (s *MemberLevelService) validateConfigValid(ctx context.Context, id int64, name string, level int, experience int) error {
	list, err := s.q.MemberLevel.WithContext(ctx).Find()
	if err != nil {
		return err
	}
	// 校验名称唯一
	if err := s.validateNameUnique(list, id, name); err != nil {
		return err
	}
	// 校验等级唯一
	if err := s.validateLevelUnique(list, id, level); err != nil {
		return err
	}
	// 校验升级所需经验是否有效: 大于前一个等级，小于下一个级别
	if err := s.validateExperienceOutRange(list, id, level, experience); err != nil {
		return err
	}
	return nil
}

func (s *MemberLevelService) validateLevelHasUser(ctx context.Context, id int64) error {
	u := s.q.MemberUser
	count, err := u.WithContext(ctx).Where(u.LevelID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return pkgErrors.NewBizError(ErrorCodeLevelHasUser, "用户等级下存在用户，无法删除")
	}
	return nil
}

// GetLevel 获得会员等级
func (s *MemberLevelService) GetLevel(ctx context.Context, id int64) (*member.MemberLevel, error) {
	if id <= 0 {
		return nil, nil
	}
	return s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.Eq(id)).First()
}

// GetLevelList 获得会员等级列表
func (s *MemberLevelService) GetLevelList(ctx context.Context, ids []int64) ([]*member.MemberLevel, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.ID.In(ids...)).Find()
}

// GetLevelListByStatus 获得指定状态的会员等级列表
func (s *MemberLevelService) GetLevelListByStatus(ctx context.Context, status int) ([]*member.MemberLevel, error) {
	return s.q.MemberLevel.WithContext(ctx).Where(s.q.MemberLevel.Status.Eq(status)).Find()
}

// GetEnableLevelList 获得开启状态的会员等级列表
func (s *MemberLevelService) GetEnableLevelList(ctx context.Context) ([]*member.MemberLevel, error) {
	return s.GetLevelListByStatus(ctx, consts.CommonStatusEnable)
}

// GetLevelSimpleList 获得开启的等级列表 (对齐旧版代码，供 Handler 使用)
func (s *MemberLevelService) GetLevelSimpleList(ctx context.Context) ([]*member.MemberLevel, error) {
	return s.GetEnableLevelList(ctx)
}

// GetLevelPage 获得会员等级分页
func (s *MemberLevelService) GetLevelPage(ctx context.Context, r *member2.MemberLevelPageReq) (*pagination.PageResult[*member.MemberLevel], error) {
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

// UpdateUserLevel 修改会员的等级
func (s *MemberLevelService) UpdateUserLevel(ctx context.Context, userId int64, levelId *int64, reason string) error {
	u := s.q.MemberUser
	user, err := u.WithContext(ctx).Where(u.ID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.NewBizError(1004001000, "用户不存在")
		}
		return err
	}

	// 等级未发生变化
	var newLevelId int64
	if levelId != nil {
		newLevelId = *levelId
	}
	if user.LevelID == newLevelId {
		return nil
	}

	var experience int
	var newLevel *member.MemberLevel
	if levelId == nil || *levelId == 0 {
		// 取消用户等级时，需要扣减经验
		experience = -int(user.Experience)
	} else {
		// 复制等级配置
		newLevel, err = s.validateLevelExists(ctx, *levelId)
		if err != nil {
			return err
		}
		// 变动经验值 = 等级的升级经验 - 会员当前的经验；正数为增加经验，负数为扣减经验
		experience = newLevel.Experience - int(user.Experience)
	}

	// 事务处理
	return s.q.Transaction(func(tx *query.Query) error {
		userExp := 0
		if newLevel != nil {
			userExp = newLevel.Experience
		}
		// 1. 记录等级变动
		// TODO: MemberLevelRecordServiceImpl.createLevelRecord

		// 2. 记录会员经验变动
		// TODO: MemberExperienceRecordServiceImpl.createExperienceRecord

		// 3. 更新会员表上的等级编号、经验值
		_, err = tx.MemberUser.WithContext(ctx).Where(tx.MemberUser.ID.Eq(userId)).Updates(&member.MemberUser{
			LevelID:    newLevelId,
			Experience: int32(userExp),
		})
		if err != nil {
			return err
		}

		// 4. 给会员发送等级变动消息
		// TODO: notifyMemberLevelChange(userId, newLevel)

		_ = experience // 避免 unused 报错

		return nil
	})
}

// AddExperience 增加会员经验
func (s *MemberLevelService) AddExperience(ctx context.Context, userId int64, experience int, bizType int, bizId string) error {
	if experience == 0 {
		return nil
	}

	// 1. 创建经验记录
	user, err := s.q.MemberUser.WithContext(ctx).Where(s.q.MemberUser.ID.Eq(userId)).First()
	if err != nil {
		return err
	}
	userExperience := int(user.Experience)
	userExperience = userExperience + experience
	if userExperience < 0 {
		userExperience = 0 // 防止扣出负数
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// 1.1 创建经验记录
		// TODO: memberExperienceRecordService.createExperienceRecord

		// 2.1 保存等级变更记录
		newLevel, err := s.calculateNewLevel(ctx, user.LevelID, userExperience)
		if err != nil {
			return err
		}

		newLevelId := user.LevelID
		if newLevel != nil {
			newLevelId = newLevel.ID
			// 2.2 创建等级变更记录
			// TODO: memberLevelRecordService.createLevelRecord

			// 2.3 给会员发送等级变动消息
			// TODO: notifyMemberLevelChange(userId, newLevel)
		}

		// 3. 更新会员表上的等级编号、经验值
		_, err = tx.MemberUser.WithContext(ctx).Where(tx.MemberUser.ID.Eq(userId)).Updates(&member.MemberUser{
			LevelID:    newLevelId,
			Experience: int32(userExperience),
		})
		return err
	})
}

// calculateNewLevel 计算会员等级
func (s *MemberLevelService) calculateNewLevel(ctx context.Context, currentLevelId int64, userExperience int) (*member.MemberLevel, error) {
	list, err := s.GetEnableLevelList(ctx)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	var matchLevel *member.MemberLevel
	for _, l := range list {
		if userExperience >= l.Experience {
			if matchLevel == nil || l.Level > matchLevel.Level {
				matchLevel = l
			}
		}
	}

	if matchLevel == nil {
		return nil, nil
	}

	// 等级没有变化
	if matchLevel.ID == currentLevelId {
		return nil, nil
	}

	return matchLevel, nil
}
