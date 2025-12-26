package consts

// CommonStatus 通用状态枚举
// 对应 Java: CommonStatusEnum
const (
	// CommonStatusEnable 开启
	CommonStatusEnable = 0
	// CommonStatusDisable 禁用
	CommonStatusDisable = 1
)

// CommonStatusValues 通用状态值数组
var CommonStatusValues = []int{CommonStatusEnable, CommonStatusDisable}

// IsValidCommonStatus 验证通用状态是否有效
func IsValidCommonStatus(status int) bool {
	for _, v := range CommonStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// IsCommonStatusEnable 判断是否为启用状态
func IsCommonStatusEnable(status int) bool {
	return status == CommonStatusEnable
}

// IsCommonStatusDisable 判断是否为禁用状态
func IsCommonStatusDisable(status int) bool {
	return status == CommonStatusDisable
}

// CommonStatusNames 通用状态名称映射
var CommonStatusNames = map[int32]string{
	CommonStatusEnable:  "启用",
	CommonStatusDisable: "禁用",
}

// GetCommonStatusName 获取通用状态名称
func GetCommonStatusName(status int32) string {
	if name, exists := CommonStatusNames[status]; exists {
		return name
	}
	return "未知"
}

// General Numeric Limits and Defaults (通用数值限制和默认值)
const (
	// 分页默认值
	DefaultPageSize = 10  // 默认分页大小
	MaxPageSize     = 100 // 最大分页大小

	// 价格相关
	MinPrice = 1         // 最小价格（分）
	MaxPrice = 999999999 // 最大价格（分）

	// 折扣相关
	MinDiscountPercent = 1    // 最小折扣百分比
	MaxDiscountPercent = 9999 // 最大折扣百分比

	// HTTP状态码常量
	HTTPStatusOK                  = 200 // 成功
	HTTPStatusBadRequest          = 400 // 请求参数错误
	HTTPStatusInternalServerError = 500 // 服务器内部错误
)

// HTTPStatusValues HTTP状态码值数组
var HTTPStatusValues = []int{HTTPStatusOK, HTTPStatusBadRequest, HTTPStatusInternalServerError}

// IsValidHTTPStatus 验证HTTP状态码是否有效
func IsValidHTTPStatus(status int) bool {
	for _, v := range HTTPStatusValues {
		if v == status {
			return true
		}
	}
	return false
}

// Sender Type Constants (发送者类型常量，对齐客服系统)
const (
	SenderTypeMember = 1 // 用户发送
	SenderTypeAdmin  = 2 // 客服发送
)

// SenderTypeValues 发送者类型值数组
var SenderTypeValues = []int{SenderTypeMember, SenderTypeAdmin}

// IsValidSenderType 验证发送者类型是否有效
func IsValidSenderType(senderType int) bool {
	for _, v := range SenderTypeValues {
		if v == senderType {
			return true
		}
	}
	return false
}

// IsSenderTypeMember 判断是否为用户发送
func IsSenderTypeMember(senderType int) bool {
	return senderType == SenderTypeMember
}

// IsSenderTypeAdmin 判断是否为客服发送
func IsSenderTypeAdmin(senderType int) bool {
	return senderType == SenderTypeAdmin
}

// SexEnum 性别枚举 (对齐 Java: SexEnum)
const (
	// SexUnknown 未知
	SexUnknown = 0
	// SexMale 男
	SexMale = 1
	// SexFemale 女
	SexFemale = 2
)

// SexValues 性别值数组
var SexValues = []int{SexUnknown, SexMale, SexFemale}

// IsValidSex 验证性别是否有效
func IsValidSex(sex int) bool {
	for _, v := range SexValues {
		if v == sex {
			return true
		}
	}
	return false
}

// SexNames 性别名称映射
var SexNames = map[int]string{
	SexUnknown: "未知",
	SexMale:    "男",
	SexFemale:  "女",
}

// GetSexName 获取性别名称
func GetSexName(sex int) string {
	if name, exists := SexNames[sex]; exists {
		return name
	}
	return "未知"
}

// MenuTypeEnum 菜单类型枚举 (对齐 Java: MenuTypeEnum)
const (
	// MenuTypeDir 目录
	MenuTypeDir = 1
	// MenuTypeMenu 菜单
	MenuTypeMenu = 2
	// MenuTypeButton 按钮
	MenuTypeButton = 3
)

// MenuTypeValues 菜单类型值数组
var MenuTypeValues = []int{MenuTypeDir, MenuTypeMenu, MenuTypeButton}

// IsValidMenuType 验证菜单类型是否有效
func IsValidMenuType(menuType int) bool {
	for _, v := range MenuTypeValues {
		if v == menuType {
			return true
		}
	}
	return false
}
