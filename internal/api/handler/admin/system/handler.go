package system

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAreaHandler,
	NewAuthHandler,
	NewDeptHandler,
	NewDictHandler,
	NewLoginLogHandler,
	NewMailHandler,
	NewMenuHandler,
	NewNoticeHandler,
	NewNotifyHandler,
	NewOAuth2ClientHandler,
	NewOperateLogHandler,
	NewPermissionHandler,
	NewPostHandler,
	NewRoleHandler,
	NewSmsChannelHandler,
	NewSmsLogHandler,
	NewSmsTemplateHandler,
	NewSocialClientHandler,
	NewSocialUserHandler,
	NewTenantHandler,
	NewTenantPackageHandler,
	NewUserHandler,
	NewHandlers,
)

type Handlers struct {
	Area          *AreaHandler
	Auth          *AuthHandler
	Dept          *DeptHandler
	Dict          *DictHandler
	LoginLog      *LoginLogHandler
	Mail          *MailHandler
	Menu          *MenuHandler
	Notice        *NoticeHandler
	Notify        *NotifyHandler
	OAuth2Client  *OAuth2ClientHandler
	OperateLog    *OperateLogHandler
	Permission    *PermissionHandler
	Post          *PostHandler
	Role          *RoleHandler
	SmsChannel    *SmsChannelHandler
	SmsLog        *SmsLogHandler
	SmsTemplate   *SmsTemplateHandler
	SocialClient  *SocialClientHandler
	SocialUser    *SocialUserHandler
	Tenant        *TenantHandler
	TenantPackage *TenantPackageHandler
	User          *UserHandler
}

func NewHandlers(
	area *AreaHandler,
	auth *AuthHandler,
	dept *DeptHandler,
	dict *DictHandler,
	loginLog *LoginLogHandler,
	mail *MailHandler,
	menu *MenuHandler,
	notice *NoticeHandler,
	notify *NotifyHandler,
	oauth2Client *OAuth2ClientHandler,
	operateLog *OperateLogHandler,
	permission *PermissionHandler,
	post *PostHandler,
	role *RoleHandler,
	smsChannel *SmsChannelHandler,
	smsLog *SmsLogHandler,
	smsTemplate *SmsTemplateHandler,
	socialClient *SocialClientHandler,
	socialUser *SocialUserHandler,
	tenant *TenantHandler,
	tenantPackage *TenantPackageHandler,
	user *UserHandler,
) *Handlers {
	return &Handlers{
		Area:          area,
		Auth:          auth,
		Dept:          dept,
		Dict:          dict,
		LoginLog:      loginLog,
		Mail:          mail,
		Menu:          menu,
		Notice:        notice,
		Notify:        notify,
		OAuth2Client:  oauth2Client,
		OperateLog:    operateLog,
		Permission:    permission,
		Post:          post,
		Role:          role,
		SmsChannel:    smsChannel,
		SmsLog:        smsLog,
		SmsTemplate:   smsTemplate,
		SocialClient:  socialClient,
		SocialUser:    socialUser,
		Tenant:        tenant,
		TenantPackage: tenantPackage,
		User:          user,
	}
}
