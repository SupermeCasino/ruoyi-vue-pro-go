package member

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	query "github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils" // Added utils
)

type MemberUserService struct {
	q          *query.Query
	smsCodeSvc *service.SmsCodeService
	levelSvc   *MemberLevelService // Injection
}

func NewMemberUserService(q *query.Query, smsCodeSvc *service.SmsCodeService, levelSvc *MemberLevelService) *MemberUserService {
	return &MemberUserService{
		q:          q,
		smsCodeSvc: smsCodeSvc,
		levelSvc:   levelSvc,
	}
}

// GetUserInfo 获取用户个人信息
func (s *MemberUserService) GetUserInfo(ctx context.Context, id int64) (*resp.AppMemberUserInfoResp, error) {
	u := s.q.MemberUser
	user, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	var levelResp *resp.AppMemberUserLevelResp
	if user.LevelID > 0 {
		level, _ := s.levelSvc.GetLevel(ctx, user.LevelID)
		if level != nil {
			levelResp = &resp.AppMemberUserLevelResp{
				ID:    level.ID,
				Name:  level.Name,
				Level: level.Level,
				Icon:  level.Icon,
			}
		}
	}

	return &resp.AppMemberUserInfoResp{
		ID:               user.ID,
		Nickname:         user.Nickname,
		Avatar:           user.Avatar,
		Mobile:           user.Mobile,
		Sex:              user.Sex,
		Point:            user.Point,
		Experience:       user.Experience,
		Level:            levelResp,
		BrokerageEnabled: bool(user.BrokerageEnabled),
	}, nil
}

// CreateUser 创建会员用户
func (s *MemberUserService) CreateUser(ctx context.Context, nickname, avatar, regIp string, terminal int32) (*member.MemberUser, error) {
	// TODO: Handle TenantID if needed
	user := &member.MemberUser{
		Nickname:         nickname,
		Avatar:           avatar,
		RegisterIP:       regIp,
		RegisterTerminal: terminal,
		Status:           0, // Enabled
		Point:            0,
		Experience:       0,
	}
	// Generate random name if nickname empty?
	if user.Nickname == "" {
		user.Nickname = "user_" + utils.GenerateRandomString(6)
	}

	if err := s.q.MemberUser.WithContext(ctx).Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser 获得用户信息 (Internal)
func (s *MemberUserService) GetUser(ctx context.Context, id int64) (*member.MemberUser, error) {
	return s.q.MemberUser.WithContext(ctx).Where(s.q.MemberUser.ID.Eq(id)).First()
}

// UpdateUser 修改用户基本信息
func (s *MemberUserService) UpdateUser(ctx context.Context, id int64, req *req.AppMemberUserUpdateReq) error {
	// 校验用户存在
	u := s.q.MemberUser
	count, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user not found")
	}

	// 更新字段
	_, err = u.WithContext(ctx).Where(u.ID.Eq(id)).
		Select(u.Nickname, u.Avatar, u.Sex, u.Birthday, u.AreaID).
		Updates(&member.MemberUser{
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
			Sex:      req.Sex,
			Birthday: req.Birthday,
			AreaID:   req.AreaID,
		})
	return err
}

// UpdateUserMobile 修改用户手机
// TODO: Implement SMS Code Verification
// UpdateUserMobile 修改用户手机
func (s *MemberUserService) UpdateUserMobile(ctx context.Context, id int64, req *req.AppMemberUserUpdateMobileReq) error {
	// 1. 校验验证码
	if err := s.smsCodeSvc.ValidateSmsCode(ctx, req.Mobile, req.Scene, req.Code); err != nil {
		return err
	}

	// 2. Check if mobile already used by another user
	u := s.q.MemberUser
	count, err := u.WithContext(ctx).Where(u.Mobile.Eq(req.Mobile), u.ID.Neq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("mobile already used")
	}

	_, err = u.WithContext(ctx).Where(u.ID.Eq(id)).Update(u.Mobile, req.Mobile)
	return err
}

// ResetUserPassword 重置用户密码 (忘记密码)
func (s *MemberUserService) ResetUserPassword(ctx context.Context, req *req.AppMemberUserResetPasswordReq) error {
	// 1. 校验验证码 (场景: 重置密码)
	// TODO: Replace magic number with Enum. MEMBER_RESET_PASSWORD = 4
	if err := s.smsCodeSvc.ValidateSmsCode(ctx, req.Mobile, 4, req.Code); err != nil {
		return err
	}

	// 2. 查询用户
	u := s.q.MemberUser
	user, err := u.WithContext(ctx).Where(u.Mobile.Eq(req.Mobile)).First()
	if err != nil {
		return errors.New("mobile not registered")
	}

	// 3. Hash Password
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// 4. Update Password
	_, err = u.WithContext(ctx).Where(u.ID.Eq(user.ID)).Update(u.Password, hashedPwd)
	return err
}

// UpdateUserPassword 修改用户密码
func (s *MemberUserService) UpdateUserPassword(ctx context.Context, id int64, req *req.AppMemberUserUpdatePasswordReq) error {
	// 1. 校验验证码
	if err := s.smsCodeSvc.ValidateSmsCode(ctx, req.Mobile, req.Scene, req.Code); err != nil {
		return err
	}

	// 2. Hash Password
	hashedPwd, err := utils.HashPassword(req.Password) // Use utils package directly or maybe better to define alias if utils not imported in this file? Utils is usually common.
	if err != nil {
		return err
	}

	u := s.q.MemberUser
	_, err = u.WithContext(ctx).Where(u.ID.Eq(id)).Update(u.Password, hashedPwd)
	return err
}

// GetUserCountByTagId 获得标签下的用户数量
func (s *MemberUserService) GetUserCountByTagId(ctx context.Context, tagId int64) (int64, error) {
	// MySQL JSON_CONTAINS equivalent
	// gorm-gen doesn't support JSON_CONTAINS directly via generic API well, using raw query or Where condition string
	// But `TagIds` is defined as `[]int64`, gorm handles serialization.
	// We need to query where tagId is in the JSON array.
	// Using `JSON_CONTAINS(tag_ids, CAST(tagId AS CHAR))` or just `JSON_CONTAINS(tag_ids, 'tagId')` depending on DB type.
	// For simplicity and safety, we can try using the gorm-gen provided capabilities or fallback to WithContext clauses.
	//
	// However, gen might generate helper for JSON array field if configured properly.
	// If not, we use:
	// return s.q.MemberUser.WithContext(ctx).Where(gorm.Expr("JSON_CONTAINS(tag_ids, CAST(? AS CHAR))", tagId)).Count()
	// Using UnderlyingDB() to execute raw condition for JSON_CONTAINS as gorm-gen specific Where expects gen.Condition
	var count int64
	err := s.q.MemberUser.WithContext(ctx).UnderlyingDB().Where("JSON_CONTAINS(tag_ids, ?)", tagId).Count(&count).Error
	return count, err
}

// UpdateUserPoint 更新用户积分
// point: 增加积分 (正数) 或 消费积分 (负数)
func (s *MemberUserService) UpdateUserPoint(ctx context.Context, id int64, point int) bool {
	if point == 0 {
		return true
	}
	// GORM update with expression: point = point + ?
	// result := s.q.MemberUser.WithContext(ctx).Where(s.q.MemberUser.ID.Eq(id)).UpdateSimple(s.q.MemberUser.Point.Add(point))
	// Since gorm-gen might need specific syntax for simple update expressions or we can use generic update.
	// Safe way with atomic update:
	u := s.q.MemberUser
	info, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Update(u.Point, u.Point.Add(int32(point)))
	if err != nil {
		return false
	}
	return info.RowsAffected > 0
}

func (s *MemberUserService) GetUserListByNickname(ctx context.Context, nickname string) ([]*member.MemberUser, error) {
	return s.q.MemberUser.WithContext(ctx).Where(s.q.MemberUser.Nickname.Like("%" + nickname + "%")).Find()
}

// GetUserMap 获得用户 Map (Entity)
func (s *MemberUserService) GetUserMap(ctx context.Context, ids []int64) (map[int64]*member.MemberUser, error) {
	users, err := s.q.MemberUser.WithContext(ctx).Where(s.q.MemberUser.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]*member.MemberUser)
	for _, u := range users {
		userMap[u.ID] = u
	}
	return userMap, nil
}

// GetUserRespMap 获得用户 Map (Response VO)
func (s *MemberUserService) GetUserRespMap(ctx context.Context, ids []int64) (map[int64]*resp.MemberUserResp, error) {
	if len(ids) == 0 {
		return make(map[int64]*resp.MemberUserResp), nil
	}
	u := s.q.MemberUser
	list, err := u.WithContext(ctx).Where(u.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]*resp.MemberUserResp, len(list))
	for _, user := range list {
		// Fetch Level info if needed or just basic info
		// For now, basic info is enough as per TradeOrder requirements
		userMap[user.ID] = &resp.MemberUserResp{
			ID:        user.ID,
			Mobile:    user.Mobile,
			Status:    user.Status,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Sex:       user.Sex,
			AreaID:    int64(user.AreaID),
			Birthday:  user.Birthday,
			Mark:      user.Mark,
			LevelID:   user.LevelID,
			GroupID:   user.GroupID,
			CreatedAt: user.CreatedAt,
			Point:     user.Point,
		}
	}
	return userMap, nil
}

// ========== Admin API Service Methods ==========

// AdminUpdateUser Admin 更新会员用户
func (s *MemberUserService) AdminUpdateUser(ctx context.Context, r *req.MemberUserUpdateReq) error {
	u := s.q.MemberUser
	count, err := u.WithContext(ctx).Where(u.ID.Eq(r.ID)).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user not found")
	}

	// Convert []int64 to []int for model.IntListFromCSV
	tagIds := make([]int, len(r.TagIDs))
	for i, v := range r.TagIDs {
		tagIds[i] = int(v)
	}

	updateUser := &member.MemberUser{
		Nickname: r.Name,
		Mark:     r.Mark,
		TagIds:   tagIds,
	}
	if r.GroupID != nil {
		updateUser.GroupID = *r.GroupID
	}

	_, err = u.WithContext(ctx).Where(u.ID.Eq(r.ID)).
		Select(u.Nickname, u.Mark, u.TagIds, u.GroupID).
		Updates(updateUser)
	return err
}

// GetUserPage Admin 获得会员用户分页
func (s *MemberUserService) GetUserPage(ctx context.Context, r *req.MemberUserPageReq) (*pagination.PageResult[*member.MemberUser], error) {
	u := s.q.MemberUser
	q := u.WithContext(ctx)

	// 动态条件
	if r.Mobile != "" {
		q = q.Where(u.Mobile.Like("%" + r.Mobile + "%"))
	}
	if r.Nickname != "" {
		q = q.Where(u.Nickname.Like("%" + r.Nickname + "%"))
	}
	if r.LevelID != nil {
		q = q.Where(u.LevelID.Eq(*r.LevelID))
	}
	if r.GroupID != nil {
		q = q.Where(u.GroupID.Eq(*r.GroupID))
	}
	if len(r.LoginDate) == 2 {
		if r.LoginDate[0] != nil {
			q = q.Where(u.LoginDate.Gte(*r.LoginDate[0]))
		}
		if r.LoginDate[1] != nil {
			q = q.Where(u.LoginDate.Lte(*r.LoginDate[1]))
		}
	}
	if len(r.CreateTime) == 2 {
		if r.CreateTime[0] != nil {
			q = q.Where(u.CreatedAt.Gte(*r.CreateTime[0]))
		}
		if r.CreateTime[1] != nil {
			q = q.Where(u.CreatedAt.Lte(*r.CreateTime[1]))
		}
	}

	// 统计总数
	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (r.PageNo - 1) * r.PageSize
	list, err := q.Order(u.ID.Desc()).Offset(offset).Limit(r.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*member.MemberUser]{
		Total: total,
		List:  list,
	}, nil
}
