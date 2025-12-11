package member

import (
	"context"
	"strings"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/member"
	"backend-go/internal/pkg/core"
	"backend-go/internal/pkg/utils"
	"backend-go/internal/repo/query"
	"backend-go/internal/service"
)

// 确保 utils 包被使用（用于密码校验）
var _ = utils.CheckPasswordHash

type MemberAuthService struct {
	repo       *query.Query
	smsCodeSvc *service.SmsCodeService
	userSvc    *MemberUserService
	socialSvc  *service.SocialUserService
	tokenSvc   *service.OAuth2TokenService
}

func NewMemberAuthService(repo *query.Query, smsCodeSvc *service.SmsCodeService, userSvc *MemberUserService, socialSvc *service.SocialUserService, tokenSvc *service.OAuth2TokenService) *MemberAuthService {
	return &MemberAuthService{
		repo:       repo,
		smsCodeSvc: smsCodeSvc,
		userSvc:    userSvc,
		socialSvc:  socialSvc,
		tokenSvc:   tokenSvc,
	}
}

// Login 手机+密码登录
func (s *MemberAuthService) Login(ctx context.Context, r *req.AppAuthLoginReq) (*resp.AppAuthLoginResp, error) {
	// 1. 查询用户
	userRepo := s.repo.MemberUser
	user, err := userRepo.WithContext(ctx).Where(userRepo.Mobile.Eq(r.Mobile)).First()
	if err != nil {
		return nil, core.NewBizError(1004003002, "账号或密码不正确") // 参考 Java MemberErrorCodeConstants
	}

	// 2. 校验状态. 0:开启, 1:关闭
	if user.Status != 0 {
		return nil, core.NewBizError(1004003001, "用户已被禁用")
	}

	// 3. 校验密码
	if !utils.CheckPasswordHash(r.Password, user.Password) {
		return nil, core.NewBizError(1004003002, "账号或密码不正确")
	}

	// 4. Check Social Bind need?
	if r.SocialType != 0 {
		bindReq := &req.SocialUserBindReq{
			Type:  r.SocialType,
			Code:  r.SocialCode,
			State: r.SocialState,
		}
		if err := s.socialSvc.BindSocialUser(ctx, user.ID, 1, bindReq); err != nil { // 1=Member
			return nil, err
		}
	}

	// 5. 生成 Token（使用 OAuth2TokenService，UserType=1 表示会员）
	return s.createToken(ctx, user)
}

// SmsLogin 手机+验证码登录
func (s *MemberAuthService) SmsLogin(ctx context.Context, r *req.AppAuthSmsLoginReq) (*resp.AppAuthLoginResp, error) {
	// 1. 校验验证码
	if err := s.smsCodeSvc.ValidateSmsCode(ctx, r.Mobile, r.Scene, r.Code); err != nil {
		return nil, err
	}

	// 2. 查询用户，不存在则注册
	userRepo := s.repo.MemberUser
	user, err := userRepo.WithContext(ctx).Where(userRepo.Mobile.Eq(r.Mobile)).First()
	if err != nil {
		// Auto Register
		createdUser, err := s.userSvc.CreateUser(ctx, "手机用户"+r.Mobile[len(r.Mobile)-4:], "", "", 0)
		if err != nil {
			return nil, err
		}
		user = createdUser
	}

	// 3. 校验状态
	if user.Status != 0 {
		return nil, core.NewBizError(1004003001, "用户已被禁用")
	}

	// 4. Bind Social if needed
	if r.SocialType != 0 {
		bindReq := &req.SocialUserBindReq{
			Type:  r.SocialType,
			Code:  r.SocialCode,
			State: r.SocialState,
		}
		if err := s.socialSvc.BindSocialUser(ctx, user.ID, 1, bindReq); err != nil {
			return nil, err
		}
	}

	// 5. 生成 Token（使用 OAuth2TokenService）
	return s.createToken(ctx, user)
}

// SocialLogin 社交快捷登录
func (s *MemberAuthService) SocialLogin(ctx context.Context, r *req.AppAuthSocialLoginReq) (*resp.AppAuthLoginResp, error) {
	// 1. 获得社交用户
	socialUser, bindUserId, err := s.socialSvc.GetSocialUserByCode(ctx, 1, int(r.Type), r.Code, r.State) // 1=Member
	if err != nil {
		return nil, err
	}
	if socialUser == nil {
		return nil, core.NewBizError(1002002000, "社交账号不存在") // AUTH_SOCIAL_USER_NOT_FOUND
	}

	var user *member.MemberUser
	if bindUserId != 0 {
		// Case 1: Already bound
		user, err = s.userSvc.GetUser(ctx, bindUserId)
		if err != nil {
			return nil, err
		}
	} else {
		// Case 2: Not bound -> Auto Register + Bind
		// Create User
		user, err = s.userSvc.CreateUser(ctx, socialUser.Nickname, socialUser.Avatar, "", 0)
		if err != nil {
			return nil, err
		}
		// Bind
		bindReq := &req.SocialUserBindReq{
			Type:  int(r.Type),
			Code:  r.Code,
			State: r.State,
		}
		if err := s.socialSvc.BindSocialUser(ctx, user.ID, 1, bindReq); err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, core.NewBizError(1004003005, "用户不存在")
	}

	// Create Token（使用 OAuth2TokenService）
	return s.createToken(ctx, user)
}

// SendSmsCode 发送验证码
func (s *MemberAuthService) SendSmsCode(ctx context.Context, r *req.AppAuthSmsSendReq) error {
	return s.smsCodeSvc.SendSmsCode(ctx, r.Mobile, r.Scene)
}

// ValidateSmsCode 校验验证码
func (s *MemberAuthService) ValidateSmsCode(ctx context.Context, r *req.AppAuthSmsValidateReq) error {
	return s.smsCodeSvc.ValidateSmsCode(ctx, r.Mobile, r.Scene, r.Code)
}

// RefreshToken 刷新访问令牌
func (s *MemberAuthService) RefreshToken(ctx context.Context, refreshToken string) (*resp.AppAuthLoginResp, error) {
	// 1. 验证 refreshToken（从 Redis 获取原令牌信息）
	oldToken, err := s.tokenSvc.GetAccessToken(ctx, refreshToken)
	if err != nil || oldToken == nil {
		return nil, core.NewBizError(401, "Token无效或已过期")
	}

	// 2. 获取用户信息
	userRepo := s.repo.MemberUser
	user, err := userRepo.WithContext(ctx).Where(userRepo.ID.Eq(oldToken.UserID)).First()
	if err != nil {
		return nil, core.NewBizError(401, "用户不存在")
	}

	// 3. 校验状态
	if user.Status != 0 {
		return nil, core.NewBizError(1004003001, "用户已被禁用")
	}

	// 4. 创建新的访问令牌
	return s.createToken(ctx, user)
}

// Logout 退出登录
func (s *MemberAuthService) Logout(ctx context.Context, token string) error {
	// 1. 处理 token，移除 Bearer 前缀
	if strings.HasPrefix(strings.ToUpper(token), "BEARER ") {
		token = token[7:]
	}
	if token == "" {
		return nil
	}

	// 2. 使用 OAuth2TokenService 删除访问令牌
	_, _ = s.tokenSvc.RemoveAccessToken(ctx, token)
	return nil
}

// createToken 创建访问令牌（使用 OAuth2TokenService，与 Java 对齐）
func (s *MemberAuthService) createToken(ctx context.Context, user *member.MemberUser) (*resp.AppAuthLoginResp, error) {
	// 构建用户信息
	userInfo := map[string]string{
		"nickname": user.Nickname,
	}

	// 创建访问令牌（UserType=1 表示会员，TenantID=0 表示默认租户）
	tokenDO, err := s.tokenSvc.CreateAccessToken(ctx, user.ID, service.UserTypeMember, 0, userInfo)
	if err != nil {
		return nil, core.ErrUnknown
	}

	return &resp.AppAuthLoginResp{
		UserID:       user.ID,
		AccessToken:  tokenDO.AccessToken,
		RefreshToken: tokenDO.RefreshToken,
		ExpiresTime:  tokenDO.ExpiresTime,
	}, nil
}
