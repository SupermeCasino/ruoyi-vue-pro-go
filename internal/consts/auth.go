package consts

// LoginLogTypeEnum 登录日志类型
const (
	LoginLogTypeUsername = 100 // 使用账号登录
	LoginLogTypeSocial   = 101 // 使用社交登录
	LoginLogTypeMobile   = 103 // 使用手机登录
	LoginLogTypeSms      = 104 // 使用短信登录

	LogoutLogTypeSelf   = 200 // 主动登出
	LogoutLogTypeDelete = 202 // 强制退出
)

// LoginResultEnum 登录结果
const (
	LoginResultSuccess          = 0  // 成功
	LoginResultBadCredentials   = 10 // 账号或密码不正确
	LoginResultUserDisabled     = 20 // 用户被禁用
	LoginResultCaptchaNotFound  = 30 // 验证码不存在
	LoginResultCaptchaCodeError = 31 // 验证码不正确
)

// UserType 用户类型枚举
// 对应 Java: UserTypeEnum
const (
	// UserTypeMember 会员用户
	UserTypeMember = 1
	// UserTypeAdmin 管理员用户
	UserTypeAdmin = 2
	// UserTypeUnknown 未知用户类型（用于默认值或错误处理）
	UserTypeUnknown = 0
)

// UserTypeNames 用户类型名称映射
var UserTypeNames = map[int]string{
	UserTypeUnknown: "未知",
	UserTypeMember:  "会员",
	UserTypeAdmin:   "管理员",
}

// GetUserTypeName 获取用户类型名称
func GetUserTypeName(userType int) string {
	if name, exists := UserTypeNames[userType]; exists {
		return name
	}
	return UserTypeNames[UserTypeUnknown]
}

// IsValidUserType 验证用户类型是否有效
func IsValidUserType(userType int) bool {
	_, exists := UserTypeNames[userType]
	return exists && userType != UserTypeUnknown
}

// UserTypes 所有有效的用户类型
var UserTypes = []int{
	UserTypeMember,
	UserTypeAdmin,
}

// NoticeTypeEnum 通知类型 (对齐 Java: NoticeTypeEnum)
const (
	// NoticeTypeNotice 通知
	NoticeTypeNotice = 1
	// NoticeTypeAnnouncement 公告
	NoticeTypeAnnouncement = 2
)

// NotifyTemplateTypeEnum 通知模板类型 (对齐 Java: NotifyTemplateTypeEnum)
const (
	// NotifyTemplateTypeNotificationMessage 通知消息
	NotifyTemplateTypeNotificationMessage = 1
	// NotifyTemplateTypeSystemMessage 系统消息
	NotifyTemplateTypeSystemMessage = 2
)

// SocialTypeEnum 社交平台类型枚举 (对齐 Java: SocialTypeEnum)
const (
	// SocialTypeGitee Gitee
	SocialTypeGitee = 10
	// SocialTypeDingTalk 钉钉
	SocialTypeDingTalk = 20
	// SocialTypeWechatEnterprise 企业微信
	SocialTypeWechatEnterprise = 30
	// SocialTypeWechatMP 微信公众平台 - 移动端 H5
	SocialTypeWechatMP = 31
	// SocialTypeWechatOpen 微信开放平台 - 网站应用 PC 端扫码授权登录
	SocialTypeWechatOpen = 32
	// SocialTypeWechatMiniProgram 微信小程序
	SocialTypeWechatMiniProgram = 34
)

// SocialTypeValues 社交平台类型值数组
var SocialTypeValues = []int{
	SocialTypeGitee,
	SocialTypeDingTalk,
	SocialTypeWechatEnterprise,
	SocialTypeWechatMP,
	SocialTypeWechatOpen,
	SocialTypeWechatMiniProgram,
}

// IsValidSocialType 验证社交平台类型是否有效
func IsValidSocialType(socialType int) bool {
	for _, v := range SocialTypeValues {
		if v == socialType {
			return true
		}
	}
	return false
}
