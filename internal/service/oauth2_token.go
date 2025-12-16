package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/cache"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

const (
	// Redis Key 前缀：访问令牌，与 Java 保持一致
	RedisKeyOAuth2AccessToken = "oauth2_access_token:%s"

	// 用户类型常量，与 Java UserTypeEnum 保持一致
	UserTypeMember = 1 // 会员
	UserTypeAdmin  = 2 // 管理员

	// 默认过期时间
	DefaultAccessTokenExpireSeconds  = 30 * 24 * 3600 // 30 天
	DefaultRefreshTokenExpireSeconds = 60 * 24 * 3600 // 60 天
)

// OAuth2AccessToken 访问令牌结构，与 Java OAuth2AccessTokenDO 对齐
type OAuth2AccessToken struct {
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	UserID       int64             `json:"userId"`
	UserType     int               `json:"userType"`
	TenantID     int64             `json:"tenantId"`
	UserInfo     map[string]string `json:"userInfo"`
	ClientID     string            `json:"clientId"`
	Scopes       []string          `json:"scopes"`
	ExpiresTime  time.Time         `json:"expiresTime"`
}

// OAuth2TokenService OAuth2 Token 服务
type OAuth2TokenService struct{}

func NewOAuth2TokenService() *OAuth2TokenService {
	return &OAuth2TokenService{}
}

// CreateAccessToken 创建访问令牌（使用 JWT 格式）
func (s *OAuth2TokenService) CreateAccessToken(ctx context.Context, userId int64, userType int, tenantId int64, userInfo map[string]string) (*OAuth2AccessToken, error) {
	// 1. 计算过期时间
	expireDuration := time.Duration(DefaultAccessTokenExpireSeconds) * time.Second
	refreshDuration := time.Duration(DefaultRefreshTokenExpireSeconds) * time.Second
	expiresTime := time.Now().Add(expireDuration)

	// 2. 获取昵称
	nickname := ""
	if userInfo != nil {
		nickname = userInfo["nickname"]
	}

	// 3. 使用 JWT 生成令牌（包含完整用户信息）
	accessToken, err := utils.GenerateTokenWithInfo(userId, userType, tenantId, nickname, expireDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateTokenWithInfo(userId, userType, tenantId, nickname, refreshDuration)
	if err != nil {
		return nil, err
	}

	// 4. 构建令牌对象
	tokenDO := &OAuth2AccessToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userId,
		UserType:     userType,
		TenantID:     tenantId,
		UserInfo:     userInfo,
		ClientID:     "default",
		Scopes:       []string{},
		ExpiresTime:  expiresTime,
	}

	// 5. 存储到 Redis（白名单机制）
	if err := s.setAccessTokenToRedis(ctx, tokenDO); err != nil {
		return nil, err
	}

	// 6. 同时存储 RefreshToken 到 Redis
	refreshTokenDO := &OAuth2AccessToken{
		AccessToken:  refreshToken,
		RefreshToken: refreshToken,
		UserID:       userId,
		UserType:     userType,
		TenantID:     tenantId,
		UserInfo:     userInfo,
		ClientID:     "default",
		Scopes:       []string{},
		ExpiresTime:  time.Now().Add(refreshDuration),
	}
	if err := s.setAccessTokenToRedis(ctx, refreshTokenDO); err != nil {
		return nil, err
	}

	return tokenDO, nil
}

// GetAccessToken 获取访问令牌
func (s *OAuth2TokenService) GetAccessToken(ctx context.Context, accessToken string) (*OAuth2AccessToken, error) {
	if cache.RDB == nil {
		return nil, nil
	}

	redisKey := fmt.Sprintf(RedisKeyOAuth2AccessToken, accessToken)
	data, err := cache.RDB.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, nil // Token 不存在
	}

	var tokenDO OAuth2AccessToken
	if err := json.Unmarshal([]byte(data), &tokenDO); err != nil {
		return nil, err
	}

	// 检查是否过期
	if time.Now().After(tokenDO.ExpiresTime) {
		return nil, nil
	}

	return &tokenDO, nil
}

// CheckAccessToken 校验访问令牌
func (s *OAuth2TokenService) CheckAccessToken(ctx context.Context, accessToken string) (*OAuth2AccessToken, error) {
	tokenDO, err := s.GetAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if tokenDO == nil {
		return nil, errors.NewBizError(401, "访问令牌不存在或已过期")
	}
	return tokenDO, nil
}

// RemoveAccessToken 删除访问令牌
func (s *OAuth2TokenService) RemoveAccessToken(ctx context.Context, accessToken string) (*OAuth2AccessToken, error) {
	// 1. 获取令牌信息
	tokenDO, _ := s.GetAccessToken(ctx, accessToken)

	// 2. 从 Redis 删除
	if cache.RDB != nil {
		redisKey := fmt.Sprintf(RedisKeyOAuth2AccessToken, accessToken)
		cache.RDB.Del(ctx, redisKey)
	}

	return tokenDO, nil
}

// RefreshAccessToken 刷新访问令牌
func (s *OAuth2TokenService) RefreshAccessToken(ctx context.Context, refreshToken string, userId int64, userType int, tenantId int64, userInfo map[string]string) (*OAuth2AccessToken, error) {
	// 直接创建新的访问令牌
	return s.CreateAccessToken(ctx, userId, userType, tenantId, userInfo)
}

// setAccessTokenToRedis 将令牌存储到 Redis
func (s *OAuth2TokenService) setAccessTokenToRedis(ctx context.Context, tokenDO *OAuth2AccessToken) error {
	if cache.RDB == nil {
		return nil // Redis 不可用时跳过
	}

	redisKey := fmt.Sprintf(RedisKeyOAuth2AccessToken, tokenDO.AccessToken)

	// 序列化
	data, err := json.Marshal(tokenDO)
	if err != nil {
		return err
	}

	// 计算剩余过期时间
	ttl := time.Until(tokenDO.ExpiresTime)
	if ttl <= 0 {
		return nil
	}

	return cache.RDB.Set(ctx, redisKey, string(data), ttl).Err()
}
