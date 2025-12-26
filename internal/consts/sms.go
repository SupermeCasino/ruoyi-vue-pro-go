package consts

// SmsSendStatus SMS 短信发送状态常量
// 对应 Java: SmsSendStatusEnum
// 与数据库 system_sms_log 表的 send_status 字段对应
const (
	// SmsSendStatusInit 发送中/初始化
	SmsSendStatusInit int32 = 0
	// SmsSendStatusSuccess 发送成功
	SmsSendStatusSuccess int32 = 10
	// SmsSendStatusFailure 发送失败
	SmsSendStatusFailure int32 = 20
	// SmsSendStatusIgnore 忽略（模板或渠道禁用，不发送）
	SmsSendStatusIgnore int32 = 30
)

// SmsReceiveStatus SMS 短信接收状态常量
// 对应 Java: SmsReceiveStatusEnum
// 与数据库 system_sms_log 表的 receive_status 字段对应
const (
	// SmsReceiveStatusInit 初始化/未接收
	SmsReceiveStatusInit int32 = 0
	// SmsReceiveStatusSuccess 接收成功
	SmsReceiveStatusSuccess int32 = 10
	// SmsReceiveStatusFailure 接收失败
	SmsReceiveStatusFailure int32 = 20
)

// SmsSendStatusNames SMS 发送状态名称映射
var SmsSendStatusNames = map[int32]string{
	SmsSendStatusInit:    "初始化",
	SmsSendStatusSuccess: "成功",
	SmsSendStatusFailure: "失败",
	SmsSendStatusIgnore:  "忽略",
}

// SmsReceiveStatusNames SMS 接收状态名称映射
var SmsReceiveStatusNames = map[int32]string{
	SmsReceiveStatusInit:    "初始化",
	SmsReceiveStatusSuccess: "成功",
	SmsReceiveStatusFailure: "失败",
}

// GetSmsSendStatusName 获取短信发送状态名称
func GetSmsSendStatusName(status int32) string {
	if name, exists := SmsSendStatusNames[status]; exists {
		return name
	}
	return "未知"
}

// GetSmsReceiveStatusName 获取短信接收状态名称
func GetSmsReceiveStatusName(status int32) string {
	if name, exists := SmsReceiveStatusNames[status]; exists {
		return name
	}
	return "未知"
}

// SMSChannelCode SMS 短信渠道代码
// 与数据库 system_sms_channel 表的 code 字段对应
const (
	SMSChannelCodeAliyun        = "ALIYUN"
	SMSChannelCodeTencent       = "TENCENT"
	SMSChannelCodeHuawei        = "HUAWEI"
	SMSChannelCodeQiniu         = "QINIU"
	SMSChannelCodeDebugDingTalk = "DEBUG_DING_TALK"
	SMSChannelCodeDebug         = "DEBUG"
)
