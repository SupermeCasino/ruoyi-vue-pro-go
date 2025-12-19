package model

// UserType 用户类型枚举
// 对应 Java: UserTypeEnum
// 统一的用户类型常量定义，避免魔法数字分散在各个文件中
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