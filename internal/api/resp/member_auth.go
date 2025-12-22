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
