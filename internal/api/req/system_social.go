package req

// SocialClientSaveReq 社交客户端创建/更新请求
type SocialClientSaveReq struct {
	ID           *int64 `json:"id"`
	Name         string `json:"name" binding:"required"`
	SocialType   int    `json:"socialType" binding:"required"`
	UserType     int    `json:"userType" binding:"required"`
	ClientId     string `json:"clientId" binding:"required"`
	ClientSecret string `json:"clientSecret" binding:"required"`
	AgentId      string `json:"agentId"`
	Status       int    `json:"status"`
}

// SocialClientPageReq 社交客户端分页请求
type SocialClientPageReq struct {
	PageNo     int    `form:"pageNo" json:"pageNo"`
	PageSize   int    `form:"pageSize" json:"pageSize"`
	Name       string `form:"name" json:"name"`
	SocialType *int   `form:"socialType" json:"socialType"`
	UserType   *int   `form:"userType" json:"userType"`
	ClientId   string `form:"clientId" json:"clientId"`
}

// SocialUserBindReq 社交用户绑定请求
type SocialUserBindReq struct {
	Type  int    `json:"type" binding:"required"`
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

// SocialUserUnbindReq 社交用户解绑请求
type SocialUserUnbindReq struct {
	Type   int    `json:"type" binding:"required"`
	Openid string `json:"openid" binding:"required"`
}

// SocialUserPageReq 社交用户分页请求
type SocialUserPageReq struct {
	PageNo   int    `form:"pageNo" json:"pageNo"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Type     *int   `form:"type" json:"type"`
	Nickname string `form:"nickname" json:"nickname"`
	Openid   string `form:"openid" json:"openid"`
}

// SocialWxaSubscribeMessageSendReq 微信小程序订阅消息发送请求
type SocialWxaSubscribeMessageSendReq struct {
	UserID        int64                  `json:"userId"`
	UserType      int                    `json:"userType"`
	TemplateTitle string                 `json:"templateTitle"`
	Page          string                 `json:"page"`
	Messages      map[string]interface{} `json:"messages"`
}
