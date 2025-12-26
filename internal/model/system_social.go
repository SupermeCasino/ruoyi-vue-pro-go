package model

// SocialUser 社交用户
type SocialUser struct {
	ID           int64  `gorm:"primaryKey;autoIncrement;comment:编号"`
	Type         int    `gorm:"column:type;not null;comment:社交类型"`
	Openid       string `gorm:"column:openid;not null;comment:社交 openid"`
	Token        string `gorm:"column:token;comment:社交 token"`
	RawTokenInfo string `gorm:"column:raw_token_info;comment:原始 Token 数据"`
	Nickname     string `gorm:"column:nickname;not null;comment:用户昵称"`
	Avatar       string `gorm:"column:avatar;comment:用户头像"`
	RawUserInfo  string `gorm:"column:raw_user_info;comment:原始用户数据"`
	Code         string `gorm:"column:code;comment:最后一次的授权码"`
	State        string `gorm:"column:state;comment:最后一次的 state"`
	TenantBaseDO
}

func (*SocialUser) TableName() string {
	return "system_social_user"
}

// SocialUserBind 社交绑定
type SocialUserBind struct {
	ID           int64 `gorm:"primaryKey;autoIncrement;comment:编号"`
	UserID       int64 `gorm:"column:user_id;not null;comment:用户编号"`
	UserType     int   `gorm:"column:user_type;not null;comment:用户类型"`
	SocialType   int   `gorm:"column:social_type;not null;comment:社交类型"`
	SocialUserID int64 `gorm:"column:social_user_id;not null;comment:社交用户编号"`
	TenantBaseDO
}

func (*SocialUserBind) TableName() string {
	return "system_social_user_bind"
}

// SocialClient 社交客户端
type SocialClient struct {
	ID           int64  `gorm:"primaryKey;autoIncrement;comment:编号"`
	Name         string `gorm:"column:name;not null;comment:应用名"`
	SocialType   int    `gorm:"column:social_type;not null;comment:社交类型"`
	UserType     int    `gorm:"column:user_type;not null;comment:用户类型"`
	ClientId     string `gorm:"column:client_id;not null;comment:客户端 id"`
	ClientSecret string `gorm:"column:client_secret;not null;comment:客户端 secret"`
	AgentId      string `gorm:"column:agent_id;comment:应用 agentId"`
	Status       int    `gorm:"column:status;not null;default:0;comment:状态"`
	TenantBaseDO
}

func (*SocialClient) TableName() string {
	return "system_social_client"
}
