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

// JsapiSignature JSAPI 签名信息
type JsapiSignature struct {
	AppID     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Timestamp int64  `json:"timestamp"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
}

// SocialPlatform 社交平台接口
type SocialPlatform interface {
	// GetAuthUser 使用 code 换取用户信息
	GetAuthUser(ctx context.Context, code string, state string) (*AuthUser, error)
	// GetAuthUrl 获得授权 URL
	GetAuthUrl(state string, redirectUri string) string
	// GetMobile 获得手机号
	GetMobile(ctx context.Context, code string) (string, error)
	// CreateJsapiSignature 创建 JSAPI 签名
	CreateJsapiSignature(ctx context.Context, url string) (*JsapiSignature, error)
	// GetWxaQrcode 获得微信小程序码
	GetWxaQrcode(ctx context.Context, path string, width int) ([]byte, error)
	// GetSubscribeTemplateList 获得订阅模板列表
	GetSubscribeTemplateList(ctx context.Context) ([]any, error)
}

// SocialPlatformFactory 社交平台工厂接口
type SocialPlatformFactory interface {
	// GetPlatform 获得社交平台客户端
	GetPlatform(ctx context.Context, socialType int, userType int) (SocialPlatform, error)
}
