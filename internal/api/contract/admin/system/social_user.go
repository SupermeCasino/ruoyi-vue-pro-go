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

type AppSocialWxaSubscribeTemplateResp struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Example string `json:"example"`
	Type    int    `json:"type"` // 2 为一次性订阅，3 为长期订阅
}
