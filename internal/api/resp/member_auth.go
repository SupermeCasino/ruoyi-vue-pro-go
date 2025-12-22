package resp

import "time"

// AppAuthLoginResp 登录响应
type AppAuthLoginResp struct {
	UserID       int64     `json:"userId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresTime  time.Time `json:"expiresTime"`
	OpenID       string    `json:"openid"`
}

// AppAuthWeixinJsapiSignatureResp 微信 JSAPI 签名响应
type AppAuthWeixinJsapiSignatureResp struct {
	AppID     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Timestamp int64  `json:"timestamp"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
}
