# 路由认证配置重构总结

## 问题背景
用户反馈 `/admin-api/system/tenant` 的三个接口（simple-list, get-by-website, get-id-by-name）应该是公共接口，不需要登录和权限认证。检查后发现，原路由配置中 `systemGroup.Use(middleware.Auth())` 被应用到所有后续路由，导致这些公共接口也需要认证。此外，还发现其他模块（SMS、Config、Infra）的路由完全没有认证保护。

## 修改内容

### 1. system.go 路由重构

#### 公共接口（无需认证）
以下接口移到 `middleware.Auth()` 之前，无需认证即可访问：

- **Auth 相关**
  - `/system/auth/login`
  - `/system/auth/logout`
  - `/system/auth/register`
  - `/system/auth/refresh-token`
  - `/system/auth/sms-login`
  - `/system/auth/send-sms-code`
  - `/system/auth/reset-password`
  - `/system/auth/social-auth-redirect`
  - `/system/auth/social-login`

- **Tenant 公共接口**
  - `/system/tenant/simple-list`
  - `/system/tenant/get-by-website`
  - `/system/tenant/get-id-by-name`

- **Dict 公共接口**
  - `/system/dict-type/simple-list`
  - `/system/dict-data/simple-list`
  - `/system/dict-data/list-all-simple`

- **Dept 公共接口**
  - `/system/dept/list`
  - `/system/dept/list-all-simple`
  - `/system/dept/simple-list`

- **Post 公共接口**
  - `/system/post/simple-list`

- **User 公共接口**
  - `/system/user/list-all-simple`
  - `/system/user/simple-list`

- **Role 公共接口**
  - `/system/role/list-all-simple`
  - `/system/role/simple-list`

- **Menu 公共接口**
  - `/system/menu/simple-list`

- **SMS 公共接口**
  - `/system/sms-channel/simple-list`

- **Area 公共接口**
  - `/system/area/tree`
  - `/system/area/get-by-ip`

#### 受保护接口（需要认证）
所有其他接口都被移到 `systemGroup.Use(middleware.Auth())` 之后，包括：

- 所有 CRUD 操作（create, update, delete）
- 所有分页查询（page）
- 所有单个资源查询（get）
- Auth 的 `/get-permission-info`
- Tenant 的 create, update, delete, get, page, export-excel
- Dict 的 create, update, delete, get, page, export-excel
- 等等所有管理操作

### 2. SMS/Config/Infra 路由认证加固

原先这些路由在 `systemGroup` 之外，没有任何认证保护。现已添加 `middleware.Auth()`：

```go
// SMS Protected Routes
smsChannelGroup := api.Group("/system/sms-channel", middleware.Auth())
smsTemplateGroup := api.Group("/system/sms-template", middleware.Auth())
smsLogGroup := api.Group("/system/sms-log", middleware.Auth())

// Infra Protected Routes  
configGroup := api.Group("/infra/config", middleware.Auth())
infraGroup := api.Group("/infra", middleware.Auth())
```

### 3. Social 路由位置修正

修正前，Social Client 和 Social User 路由错误地放在 `infraGroup` 内部却引用 `systemGroup`：

```go
// 错误: 在 infraGroup 内部引用 systemGroup
socialClientGroup := systemGroup.Group("/social-client")
```

修正后，移到正确的位置：

```go
// 正确: 在 systemGroup 内部，middleware.Auth() 之后
socialClientProtectedGroup := systemGroup.Group("/social-client")
socialUserProtectedGroup := systemGroup.Group("/social-user")
```

## 设计原则

1. **公共接口**：`simple-list`, `list-all-simple` 等列表接口通常用于前端下拉框、选择器等，应该是公共接口
2. **认证接口**：所有 CRUD 操作（create, update, delete）必须需要认证
3. **管理接口**：page, get 等管理查询接口需要认证
4. **回调接口**：第三方支付回调等特殊接口不需要认证（如 pay notify）

## 注意事项

1. **pay.go** 路由当前未添加认证保护，但 Pay Notify 的回调接口应保持公共（供第三方调用）
2. **权限控制**：认证（Auth）和权限（Permission）是两个不同的层次，本次修改只处理认证层，部分接口还需要通过 `casbinMiddleware.RequirePermission()` 进行权限控制
3. **Area 路由**：地区查询接口被设为公共接口，便于前端获取省市区数据

## 测试建议

1. 测试公共接口无需 token 即可访问
2. 测试受保护接口没有 token 时返回 401
3. 测试 simple-list 等列表接口的功能正常
