package resp

type AppSocialUserResp struct {
	Openid   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
