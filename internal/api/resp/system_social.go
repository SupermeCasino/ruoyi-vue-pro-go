package resp

import "time"

// SocialClientResp 社交客户端响应
type SocialClientResp struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	SocialType   int       `json:"socialType"`
	UserType     int       `json:"userType"`
	ClientId     string    `json:"clientId"`
	ClientSecret string    `json:"clientSecret"`
	AgentId      string    `json:"agentId"`
	Status       int       `json:"status"`
	CreateTime   time.Time `json:"createTime"`
}

// SocialUserResp 社交用户响应
type SocialUserResp struct {
	ID          int64     `json:"id"`
	Type        int       `json:"type"`
	Openid      string    `json:"openid"`
	Token       string    `json:"token"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	RawUserInfo string    `json:"rawUserInfo"`
	Code        string    `json:"code"`
	State       string    `json:"state"`
	CreateTime  time.Time `json:"createTime"`
}
