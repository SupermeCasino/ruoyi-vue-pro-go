package datascope

import (
	"context"
)

// DataScope 数据范围类型
type DataScope int

const (
	DataScopeAll          DataScope = 1 // 全部数据权限
	DataScopeDeptCustom   DataScope = 2 // 指定部门数据权限
	DataScopeDeptOnly     DataScope = 3 // 本部门数据权限
	DataScopeDeptAndChild DataScope = 4 // 本部门及以下数据权限
	DataScopeSelf         DataScope = 5 // 仅本人数据权限
)

// Context Keys
type contextKey string

const (
	// SkipDataScopeKey 跳过数据权限检查的Context Key
	SkipDataScopeKey contextKey = "skip_data_scope"
)

// SkipDataScope 返回一个跳过数据权限检查的Context
func SkipDataScope(ctx context.Context) context.Context {
	return context.WithValue(ctx, SkipDataScopeKey, true)
}

// ShouldSkipDataScope 检查是否应该跳过数据权限检查
func ShouldSkipDataScope(ctx context.Context) bool {
	if v := ctx.Value(SkipDataScopeKey); v != nil {
		if skip, ok := v.(bool); ok {
			return skip
		}
	}
	return false
}
