# GORM 审计插件使用说明

## 功能概述

GORM 审计插件（`AuditPlugin`）会自动为数据库操作填充以下字段：

- **Creator**: 创建记录时，自动填充当前登录用户 ID（string 类型）
- **Updater**: 更新记录时，自动填充当前登录用户 ID（string 类型）
- **TenantID**: 创建记录时，自动填充当前租户 ID（int64 类型）

## 工作原理

1. **Auth 中间件**：用户登录后，`middleware.Auth()` 将 `LoginUser`（包含 UserID、TenantID）存储到 `gin.Context`
2. **InjectContext 中间件**：`middleware.InjectContext()` 将 `gin.Context` 注入到 `request.Context` 中
3. **GORM Hook**：在数据库操作前，`AuditPlugin` 的 `beforeCreate`/`beforeUpdate` Hook 从 context 获取用户信息，并填充相应字段

## 字段存在性检查

插件会在设置字段前检查模型是否有对应字段：

- 如果表中没有 `tenant_id` 字段，不会尝试设置 TenantID
- 如果表中没有 `creator` 或 `updater` 字段，不会尝试设置这些字段

这避免了对不需要租户隔离的表或没有审计字段的表报错。

## 手动测试步骤

### 1. 启动应用
\`\`\`bash
cd /Users/wxl/GolandProjects/yudao/backend-go
go run cmd/server/main.go
\`\`\`

### 2. 登录获取 Token
\`\`\`bash
curl -X POST http://localhost:8080/admin-api/system/auth/login \\
  -H "Content-Type: application/json" \\
  -d '{"username":"admin","password":"admin123"}'
\`\`\`

保存返回的 `accessToken`。

### 3. 创建角色测试
\`\`\`bash
curl -X POST http://localhost:8080/admin-api/system/role/create \\
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{"status":0,"name":"测试审计","code":"test_audit","sort":0}'
\`\`\`

### 4. 检查数据库

连接到 MySQL 数据库，查询新创建的角色：

\`\`\`sql
SELECT id, name, code, creator, tenant_id, create_time 
FROM system_role 
WHERE code = 'test_audit';
\`\`\`

**预期结果**：
- `creator` 字段应该是当前登录用户 ID（如 "1"）
- `tenant_id` 字段应该是当前租户 ID（如果用户有租户）

### 5. 更新角色测试

\`\`\`bash
curl -X PUT http://localhost:8080/admin-api/system/role/update \\
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{"id":ROLE_ID,"status":0,"name":"测试审计更新","code":"test_audit","sort":1}'
\`\`\`

再次查询数据库：

\`\`\`sql
SELECT id, name, updater, update_time 
FROM system_role 
WHERE id = ROLE_ID;
\`\`\`

**预期结果**：
- `updater` 字段应该是当前登录用户 ID
- `update_time` 已更新

## 故障排查

### 问题：creator/updater 字段仍然为空

**可能原因**：
1. 用户未登录（没有 token）
2. Token 无效
3. Auth 中间件未正确设置 LoginUser

**检查方法**：
1. 确认请求带了有效的 Authorization header
2. 检查日志是否有 Auth 中间件相关错误
3. 在 `gorm_plugin.go` 的 `beforeCreate` 中添加调试日志：
   \`\`\`go
   user := GetLoginUser(ginCtx)
   if user == nil {
       fmt.Println("WARNING: No LoginUser found in context")
       return
   }
   fmt.Printf("Setting Creator: %d\\n", user.UserID)
   \`\`\`

### 问题：某些表报错 "unknown column 'tenant_id'"

**可能原因**：表确实没有 `tenant_id` 字段

**解决方案**：
- 插件已经包含字段存在性检查，这个错误不应该发生
- 如果仍然发生，检查 `hasField` 函数是否正常工作
- 确认 GORM Schema 正确加载

## 代码位置

- **插件实现**：`internal/pkg/core/gorm_plugin.go`
- **Context 中间件**：`internal/middleware/context.go`
- **插件注册**：`internal/pkg/core/db.go` (第 44-47 行)
- **中间件注册**：`internal/api/router/router.go` (第 154-155 行)
