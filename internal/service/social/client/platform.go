package client

import (
	"context"
)

// AuthUser 社交平台返回的用户信息
type AuthUser struct {
	Openid       string
	Token        string
	RawTokenInfo string
	Nickname     string
	Avatar       string
	RawUserInfo  string
}

// SocialPlatform 社交平台接口
type SocialPlatform interface {
	// GetAuthUser 使用 code 换取用户信息
	GetAuthUser(ctx context.Context, code string, state string) (*AuthUser, error)
}

// SocialPlatformFactory 社交平台工厂接口
type SocialPlatformFactory interface {
	// GetPlatform 获得社交平台客户端
	GetPlatform(ctx context.Context, socialType int, userType int) (SocialPlatform, error)
}
