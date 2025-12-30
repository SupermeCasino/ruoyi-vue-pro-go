package consts

// RoleTypeEnum 角色类型枚举
// 对应 Java: RoleTypeEnum
const (
	// RoleTypeSystem 内置系统角色
	RoleTypeSystem = 1
	// RoleTypeCustom 自定义角色
	RoleTypeCustom = 2
)

// DataScopeEnum 数据范围枚举
// 对应 Java: DataScopeEnum
const (
	// DataScopeAll 全部数据权限
	DataScopeAll = 1
	// DataScopeDeptCustom 指定部门数据权限
	DataScopeDeptCustom = 2
	// DataScopeDeptOnly 部门数据权限
	DataScopeDeptOnly = 3
	// DataScopeDeptAndChild 部门及以下数据权限
	DataScopeDeptAndChild = 4
	// DataScopeSelf 仅本人数据权限
	DataScopeSelf = 5
)

// RoleCodeEnum 角色编码枚举
// 对应 Java: RoleCodeEnum
const (
	// RoleCodeSuperAdmin 超级管理员
	RoleCodeSuperAdmin = "super_admin"
	// RoleCodeTenantAdmin 租户管理员
	RoleCodeTenantAdmin = "tenant_admin"
)
