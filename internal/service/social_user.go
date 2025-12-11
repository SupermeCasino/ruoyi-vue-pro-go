package service

import (
	"context"
	"errors"

	"backend-go/internal/api/req"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/social/client"

	"gorm.io/gorm"
)

type SocialUserService struct {
	q       *query.Query
	factory client.SocialPlatformFactory
}

func NewSocialUserService(q *query.Query) *SocialUserService {
	return &SocialUserService{
		q:       q,
		factory: client.NewSocialPlatformFactory(q),
	}
}

// GetAuthorizeUrl 获取授权链接
func (s *SocialUserService) GetAuthorizeUrl(ctx context.Context, socialType int, userType int, redirectUri string) (string, error) {
	// 1. 获得社交平台客户端
	platform, err := s.factory.GetPlatform(ctx, socialType, userType)
	if err != nil {
		return "", err
	}

	// 2. 生成 State (这里简化，暂不存储校验，实际应存 Redis)
	state := "test-state" // TODO: Use UUID and Store in Redis for validation

	// 3. 生成链接
	return platform.GetAuthUrl(state, redirectUri), nil
}

// BindSocialUser 绑定社交用户
func (s *SocialUserService) BindSocialUser(ctx context.Context, userID int64, userType int, r *req.SocialUserBindReq) error {
	// 1. 获得社交平台客户端
	platform, err := s.factory.GetPlatform(ctx, r.Type, userType)
	if err != nil {
		return err
	}

	// 2. 获得社交用户信息
	authUser, err := platform.GetAuthUser(ctx, r.Code, r.State)
	if err != nil {
		return err
	}

	// 3. 查找或创建社交用户
	socialUser, err := s.authSocialUser(ctx, r.Type, authUser)
	if err != nil {
		return err
	}

	// 4. 绑定
	// 检查是否已经绑定了该账号
	count, err := s.q.SocialUserBind.WithContext(ctx).Where(
		s.q.SocialUserBind.UserID.Eq(userID),
		s.q.SocialUserBind.UserType.Eq(userType),
		s.q.SocialUserBind.SocialType.Eq(r.Type),
		// s.q.SocialUserBind.SocialUserID.Eq(socialUser.ID), // 这里的逻辑应该是：同一个 User 对同一个 SocialType 只能绑定一个? 或者只能绑定特定的?
		// 通常允许绑定多个吗？RuoYi 默认可以绑定多个吗？应该是一对一的 (User <-> SocialType instance).
		// 但为了健壮性，我们可以检查是否已经绑定了这个特定的 SocialUser
	).Count()
	if err != nil {
		return err
	}

	// 此外，还要检查该 സോഷ്യ셜用户是否已经被 *其他* 用户绑定
	bindCount, err := s.q.SocialUserBind.WithContext(ctx).Where(
		s.q.SocialUserBind.SocialType.Eq(r.Type),
		s.q.SocialUserBind.SocialUserID.Eq(socialUser.ID),
	).Count()
	if err != nil {
		return err
	}
	if bindCount > 0 {
		return core.NewBizError(1002004005, "该社交账号已被绑定")
	}

	if count == 0 {
		bind := &model.SocialUserBind{
			UserID:       userID,
			UserType:     userType,
			SocialType:   r.Type,
			SocialUserID: socialUser.ID,
		}
		if err := s.q.SocialUserBind.WithContext(ctx).Create(bind); err != nil {
			return err
		}
	}
	// 如果已经绑定，是否更新？不需要。

	return nil
}

// UnbindSocialUser 解绑社交用户
func (s *SocialUserService) UnbindSocialUser(ctx context.Context, userID int64, userType int, socialType int, openid string) error {
	// 查找社交用户
	socialUser, err := s.q.SocialUser.WithContext(ctx).Where(
		s.q.SocialUser.Type.Eq(socialType),
		s.q.SocialUser.Openid.Eq(openid),
	).First()
	if err != nil {
		return err // Not found
	}

	// 删除绑定关系
	_, err = s.q.SocialUserBind.WithContext(ctx).Where(
		s.q.SocialUserBind.UserID.Eq(userID),
		s.q.SocialUserBind.UserType.Eq(userType),
		s.q.SocialUserBind.SocialType.Eq(socialType),
		s.q.SocialUserBind.SocialUserID.Eq(socialUser.ID),
	).Delete()
	return err
}

// GetSocialUserList 获取用户绑定的社交账号列表
func (s *SocialUserService) GetSocialUserList(ctx context.Context, userID int64, userType int) ([]*model.SocialUser, error) {
	// 查找绑定关系
	binds, err := s.q.SocialUserBind.WithContext(ctx).Where(
		s.q.SocialUserBind.UserID.Eq(userID),
		s.q.SocialUserBind.UserType.Eq(userType),
	).Find()
	if err != nil {
		return nil, err
	}

	if len(binds) == 0 {
		return []*model.SocialUser{}, nil
	}

	// 提取社交用户ID列表
	socialUserIDs := make([]int64, len(binds))
	for i, bind := range binds {
		socialUserIDs[i] = bind.SocialUserID
	}

	// 查询社交用户
	return s.q.SocialUser.WithContext(ctx).Where(s.q.SocialUser.ID.In(socialUserIDs...)).Find()
}

// GetSocialUser 获取社交用户
func (s *SocialUserService) GetSocialUser(ctx context.Context, id int64) (*model.SocialUser, error) {
	return s.q.SocialUser.WithContext(ctx).Where(s.q.SocialUser.ID.Eq(id)).First()
}

// GetSocialUserPage 获取社交用户分页
func (s *SocialUserService) GetSocialUserPage(ctx context.Context, r *req.SocialUserPageReq) (*core.PageResult[*model.SocialUser], error) {
	q := s.q.SocialUser.WithContext(ctx)

	if r.Type != nil {
		q = q.Where(s.q.SocialUser.Type.Eq(*r.Type))
	}
	if r.Nickname != "" {
		q = q.Where(s.q.SocialUser.Nickname.Like("%" + r.Nickname + "%"))
	}
	if r.Openid != "" {
		q = q.Where(s.q.SocialUser.Openid.Like("%" + r.Openid + "%"))
	}

	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.SocialUser.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*model.SocialUser]{
		List:  list,
		Total: total,
	}, nil
}

// GetSocialUserByCode 根据 Code 获得社交用户 (用于登录)
func (s *SocialUserService) GetSocialUserByCode(ctx context.Context, userType int, socialType int, code string, state string) (*model.SocialUser, int64, error) {
	// 1. 获得社交平台客户端
	platform, err := s.factory.GetPlatform(ctx, socialType, userType)
	if err != nil {
		return nil, 0, err
	}

	// 2. 获得社交用户信息
	authUser, err := platform.GetAuthUser(ctx, code, state)
	if err != nil {
		return nil, 0, err
	}

	// 3. 查找或创建社交用户
	socialUser, err := s.authSocialUser(ctx, socialType, authUser)
	if err != nil {
		return nil, 0, err
	}

	// 4. 检查绑定
	bind, err := s.q.SocialUserBind.WithContext(ctx).
		Where(s.q.SocialUserBind.UserType.Eq(userType), s.q.SocialUserBind.SocialUserID.Eq(socialUser.ID)).
		First()

	var userID int64
	if err == nil {
		userID = bind.UserID
	}

	return socialUser, userID, nil
}

// authSocialUser 内部方法：查找或创建社交用户
func (s *SocialUserService) authSocialUser(ctx context.Context, socialType int, authUser *client.AuthUser) (*model.SocialUser, error) {
	// 1. 查找社交用户
	socialUser, err := s.q.SocialUser.WithContext(ctx).
		Where(s.q.SocialUser.Type.Eq(socialType), s.q.SocialUser.Openid.Eq(authUser.Openid)).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		socialUser = &model.SocialUser{
			Type:   socialType,
			Openid: authUser.Openid,
		}
	} else if err != nil {
		return nil, err
	}

	// 2. 更新属性
	socialUser.Token = authUser.Token
	socialUser.RawTokenInfo = authUser.RawTokenInfo
	socialUser.Nickname = authUser.Nickname
	socialUser.Avatar = authUser.Avatar
	socialUser.RawUserInfo = authUser.RawUserInfo
	// socialUser.Code, State 没保留在 AuthUser 里，如果需要可以加上，或者直接传
	// 既然是绑定流程，最新的 Code/State 应该更新
	// 这里简化，暂不更新 Code/State 为请求中的参数，因为 binding 是一次性的

	// 3. 保存
	if socialUser.ID == 0 {
		if err := s.q.SocialUser.WithContext(ctx).Create(socialUser); err != nil {
			return nil, err
		}
	} else {
		if _, err := s.q.SocialUser.WithContext(ctx).Where(s.q.SocialUser.ID.Eq(socialUser.ID)).Updates(socialUser); err != nil {
			return nil, err
		}
	}

	return socialUser, nil
}
