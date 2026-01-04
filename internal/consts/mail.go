package consts

import "github.com/wxlbd/ruoyi-mall-go/pkg/errors"

// Mail 发送状态枚举
// 对应 Java: MailSendStatusEnum
const (
	MailSendStatusInit    = 0  // 初始化
	MailSendStatusSuccess = 10 // 发送成功
	MailSendStatusFailure = 20 // 发送失败
	MailSendStatusIgnore  = 30 // 忽略，即不发送
)

// Mail 业务错误码
// 对应 Java: ErrorCodeConstants
var (
	ErrMailAccountNotExists            = errors.NewBizError(1002023000, "邮箱账号不存在")
	ErrMailAccountRelateTemplateExists = errors.NewBizError(1002023001, "无法删除，该邮箱账号还有邮件模板")
	ErrMailTemplateNotExists           = errors.NewBizError(1002024000, "邮件模版不存在")
	ErrMailTemplateCodeExists          = errors.NewBizError(1002024001, "邮件模版 code 已存在")
	ErrMailSendTemplateParamMiss       = errors.NewBizError(1002025000, "模板参数缺失")
	ErrMailSendMailNotExists           = errors.NewBizError(1002025001, "邮箱不存在")
)
