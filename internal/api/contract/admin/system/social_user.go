package system

type AppSocialUserBindReq struct {
	Type  int    `json:"type" binding:"required"`
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

type AppSocialUserUnbindReq struct {
	Type   int    `json:"type" binding:"required"`
	OpenID string `json:"openid" binding:"required"`
}

type AppSocialUserResp struct {
	Openid   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
