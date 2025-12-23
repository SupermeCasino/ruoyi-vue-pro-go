# 芋道商城 Go 版本 - 快速参考指南

> 本文档提供项目开发中的快速查询和常用代码片段。

## 目录

- [常用命令](#常用命令)
- [代码片段](#代码片段)
- [API 调用示例](#api-调用示例)
- [数据库操作](#数据库操作)
- [错误处理](#错误处理)
- [调试技巧](#调试技巧)

---

## 常用命令

### 项目构建和运行

```bash
# 下载依赖
make deps
# 或
go mod tidy

# 直接运行
make run
# 或
go run cmd/server/main.go

# 热重载开发（推荐）
make dev

# 编译构建
make build

# 清理构建产物
make clean

# 重新生成 Wire 依赖注入代码
make wire

# 重新生成 GORM DAO 代码
make gen
```

### 数据库操作

```bash
# 连接 MySQL
mysql -h localhost -u root -p yudao

# 查看表结构
DESCRIBE table_name;

# 查看所有表
SHOW TABLES;

# 导出数据
mysqldump -h localhost -u root -p yudao > backup.sql

# 导入数据
mysql -h localhost -u root -p yudao < backup.sql
```

### Redis 操作

```bash
# 连接 Redis
redis-cli

# 查看所有 key
KEYS *

# 查看 key 类型
TYPE key_name

# 查看 key 值
GET key_name

# 删除 key
DEL key_name

# 查看 key 过期时间
TTL key_name

# 清空所有数据
FLUSHALL
```

---

## 代码片段

### 获取登录用户信息

```go
// 在 Handler 或 Service 中
loginUser := core.GetLoginUser(c)
if loginUser == nil {
    core.WriteError(c, core.UnauthorizedCode, "未登录")
    return
}

userID := loginUser.UserID
userType := loginUser.UserType      // 0=Member, 1=Admin
tenantID := loginUser.TenantID
nickname := loginUser.Nickname
```

### 写入响应

```go
// 成功响应
core.WriteSuccess(c, data)

// 错误响应
core.WriteError(c, core.ParamErrCode, "参数错误")

// 业务异常响应
core.WriteBizError(c, err)

// 分页响应
core.WritePage(c, total, list)
```

### 参数验证

```go
// 绑定 JSON 参数
var req req.CreateUserReq
if err := c.ShouldBindJSON(&req); err != nil {
    core.WriteError(c, core.ParamErrCode, err.Error())
    return
}

// 绑定 Query 参数
pageNo := c.DefaultQuery("pageNo", "1")
pageSize := c.DefaultQuery("pageSize", "10")

// 绑定 Path 参数
id := c.Param("id")
```

### 数据库查询

```go
// 单条查询
user, err := s.q.SystemUser.WithContext(ctx).
    Where(s.q.SystemUser.ID.Eq(id)).
    First()

// 列表查询
users, err := s.q.SystemUser.WithContext(ctx).
    Where(s.q.SystemUser.Status.Eq(0)).
    Order(s.q.SystemUser.CreateTime.Desc()).
    Offset(offset).
    Limit(limit).
    Find()

// 统计
count, err := s.q.SystemUser.WithContext(ctx).
    Where(s.q.SystemUser.Status.Eq(0)).
    Count()

// 更新
_, err := s.q.SystemUser.WithContext(ctx).
    Where(s.q.SystemUser.ID.Eq(id)).
    Updates(&model.SystemUser{
        Nickname: "新昵称",
        Status:   1,
    })

// 删除
_, err := s.q.SystemUser.WithContext(ctx).
    Where(s.q.SystemUser.ID.Eq(id)).
    Delete()
```

### 事务处理

```go
// 开启事务
tx := s.q.WithContext(ctx).Begin()

// 执行操作
if err := tx.Create(&entity).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
return tx.Commit().Error
```

### Redis 操作

```go
// 设置值
core.RDB.Set(ctx, "key", "value", time.Hour)

// 获取值
val, err := core.RDB.Get(ctx, "key").Result()

// 删除值
core.RDB.Del(ctx, "key")

// 检查是否存在
exists, err := core.RDB.Exists(ctx, "key").Result()

// 设置过期时间
core.RDB.Expire(ctx, "key", time.Hour)

// 获取过期时间
ttl, err := core.RDB.TTL(ctx, "key").Result()
```

### JWT Token 操作

```go
// 生成 Token
token, err := utils.GenerateTokenWithInfo(
    userID,      // 用户 ID
    userType,    // 用户类型 (0=Member, 1=Admin)
    tenantID,    // 租户 ID
    nickname,    // 用户昵称
    24*time.Hour, // 有效期
)

// 解析 Token
claims, err := utils.ParseToken(token)
if err != nil {
    // Token 无效
}

userID := claims.UserID
userType := claims.UserType
```

### 密码操作

```go
// 加密密码
hashedPwd, err := utils.HashPassword(plainPassword)

// 验证密码
isMatch := utils.CheckPassword(plainPassword, hashedPwd)
```

### 日志输出

```go
// Info 级别
logger.Info("操作成功", zap.String("user", username))

// Warn 级别
logger.Log.Warn("警告信息", zap.Error(err))

// Error 级别
logger.Log.Error("错误信息", zap.Error(err))

// Debug 级别
logger.Log.Debug("调试信息", zap.Any("data", data))
```

---

## API 调用示例

### 用户登录

```bash
curl -X POST http://localhost:48080/admin-api/system/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 响应
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员"
    }
  }
}
```

### 获取用户列表

```bash
curl -X GET "http://localhost:48080/admin-api/system/user/page?pageNo=1&pageSize=10" \
  -H "Authorization: Bearer eyJhbGc..."

# 响应
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "管理员",
        "email": "admin@example.com"
      }
    ],
    "total": 100
  }
}
```

### 创建用户

```bash
curl -X POST http://localhost:48080/admin-api/system/user/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{
    "username": "newuser",
    "password": "123456",
    "nickname": "新用户",
    "email": "newuser@example.com",
    "mobile": "13800138000",
    "deptId": 1,
    "status": 0
  }'

# 响应
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2
  }
}
```

### 更新用户

```bash
curl -X PUT http://localhost:48080/admin-api/system/user/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{
    "id": 2,
    "nickname": "更新昵称",
    "email": "newemail@example.com",
    "status": 0
  }'

# 响应
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

### 删除用户

```bash
curl -X DELETE http://localhost:48080/admin-api/system/user/delete/2 \
  -H "Authorization: Bearer eyJhbGc..."

# 响应
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

### 获取商品列表

```bash
curl -X GET "http://localhost:48080/app-api/product/spu/list?categoryId=1&pageNo=1&pageSize=10" \
  -H "Authorization: Bearer eyJhbGc..."

# 响应
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "商品名称",
        "price": 99.99,
        "pictures": ["url1", "url2"],
        "rating": 4.5
      }
    ],
    "total": 1000
  }
}
```

### 创建订单

```bash
curl -X POST http://localhost:48080/app-api/trade/order/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{
    "items": [
      {
        "skuId": 1,
        "quantity": 2
      }
    ],
    "addressId": 1,
    "couponId": null
  }'

# 响应
{
  "code": 0,
  "msg": "success",
  "data": {
    "orderId": 123456,
    "orderNo": "2024121512345678",
    "totalAmount": 999.99,
    "payAmount": 899.99
  }
}
```

---

## 数据库操作

### 常用 SQL 语句

```sql
-- 查询用户
SELECT * FROM system_user WHERE id = 1;

-- 查询用户列表（分页）
SELECT * FROM system_user 
WHERE status = 0 
ORDER BY created_at DESC 
LIMIT 10 OFFSET 0;

-- 统计用户数量
SELECT COUNT(*) FROM system_user WHERE status = 0;

-- 更新用户
UPDATE system_user SET nickname = '新昵称' WHERE id = 1;

-- 删除用户
DELETE FROM system_user WHERE id = 1;

-- 查询订单及订单项
SELECT o.*, oi.* 
FROM trade_order o 
LEFT JOIN trade_order_item oi ON o.id = oi.order_id 
WHERE o.user_id = 1;

-- 统计订单金额
SELECT SUM(total_amount) FROM trade_order WHERE user_id = 1;

-- 查询商品评论
SELECT * FROM product_comment 
WHERE spu_id = 1 
ORDER BY created_at DESC 
LIMIT 10;

-- 统计商品评分
SELECT AVG(rating), COUNT(*) FROM product_comment WHERE spu_id = 1;
```

### GORM 查询技巧

```go
// 条件查询
u := s.q.SystemUser
users, err := u.WithContext(ctx).
    Where(u.Status.Eq(0)).
    Where(u.DeptID.Eq(deptID)).
    Find()

// OR 条件
users, err := u.WithContext(ctx).
    Where(u.Status.Eq(0).Or(u.Status.Eq(1))).
    Find()

// IN 条件
users, err := u.WithContext(ctx).
    Where(u.ID.In(1, 2, 3)).
    Find()

// LIKE 条件
users, err := u.WithContext(ctx).
    Where(u.Username.Like("%" + keyword + "%")).
    Find()

// 范围查询
users, err := u.WithContext(ctx).
    Where(u.CreateTime.Between(startTime, endTime)).
    Find()

// 排序
users, err := u.WithContext(ctx).
    Order(u.CreateTime.Desc()).
    Order(u.ID.Asc()).
    Find()

// 分页
users, err := u.WithContext(ctx).
    Offset((pageNo - 1) * pageSize).
    Limit(pageSize).
    Find()

// 分组
type Result struct {
    DeptID int64
    Count  int64
}
var results []Result
err := s.q.SystemUser.WithContext(ctx).
    Select(s.q.SystemUser.DeptID, s.q.SystemUser.ID.Count()).
    Group(s.q.SystemUser.DeptID).
    Scan(&results)
```

---

## 错误处理

### 常见错误码

| 错误码 | 含义 | 处理方式 |
|--------|------|---------|
| 0 | 成功 | 正常返回 |
| 400 | 参数错误 | 检查请求参数 |
| 401 | 未授权 | 重新登录 |
| 403 | 禁止访问 | 检查权限 |
| 404 | 资源不存在 | 检查资源ID |
| 409 | 冲突 | 检查唯一性约束 |
| 500 | 系统异常 | 查看日志 |

### 错误处理模式

```go
// 模式 1：直接返回错误
if err != nil {
    core.WriteBizError(c, err)
    return
}

// 模式 2：自定义错误信息
if user == nil {
    core.WriteError(c, core.NotFoundCode, "用户不存在")
    return
}

// 模式 3：业务异常
if user.Status != 0 {
    core.WriteError(c, 1001001001, "用户已被禁用")
    return
}

// 模式 4：参数验证错误
if req.Username == "" {
    core.WriteError(c, core.ParamErrCode, "用户名不能为空")
    return
}
```

---

## 调试技巧

### 启用 SQL 日志

```go
// 在 GORM 查询中添加 Debug()
users, err := s.q.SystemUser.WithContext(ctx).
    Where(s.q.SystemUser.Status.Eq(0)).
    Debug().  // 启用调试，会打印 SQL
    Find()

// 输出示例：
// SELECT * FROM system_user WHERE status = 0
```

### 打印变量值

```go
// 使用 zap 日志
logger.Log.Info("用户信息", zap.Any("user", user))

// 使用 fmt
fmt.Printf("用户ID: %d\n", user.ID)

// 使用 JSON 序列化
data, _ := json.MarshalIndent(user, "", "  ")
fmt.Println(string(data))
```

### 设置断点调试

```bash
# 使用 Delve 调试器
dlv debug cmd/server/main.go

# 在 IDE 中设置断点并运行
# GoLand: Run → Debug
```

### 查看日志

```bash
# 查看实时日志
tail -f logs/app.log

# 搜索特定日志
grep "error" logs/app.log

# 查看最后 100 行日志
tail -100 logs/app.log
```

### 性能分析

```bash
# 使用 pprof 进行性能分析
go tool pprof http://localhost:6060/debug/pprof/profile

# 查看内存使用
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看 goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

---

## 快速检查清单

### 添加新 API 时

- [ ] 定义请求对象 (req/*.go)
- [ ] 定义响应对象 (resp/*.go)
- [ ] 实现 Handler 方法
- [ ] 实现 Service 方法
- [ ] 实现 Repository 方法（如需要）
- [ ] 注册路由
- [ ] 添加参数验证
- [ ] 添加错误处理
- [ ] 添加权限检查（如需要）
- [ ] 测试 API

### 修改数据模型时

- [ ] 修改 model/*.go
- [ ] 运行 `make gen` 重新生成 DAO 代码
- [ ] 更新 Repository 查询逻辑
- [ ] 更新 Service 业务逻辑
- [ ] 更新 Handler 请求/响应对象
- [ ] 测试数据库操作

### 部署前检查

- [ ] 更新配置文件 (config/config.local.yaml)
- [ ] 检查数据库连接
- [ ] 检查 Redis 连接
- [ ] 检查日志输出目录
- [ ] 运行单元测试
- [ ] 检查错误日志
- [ ] 验证 API 功能

---

## 常见问题快速解决

### Q: Token 过期如何处理？

```go
if errors.Is(err, jwt.ErrTokenExpired) {
    core.WriteError(c, core.UnauthorizedCode, "Token 已过期，请重新登录")
    return
}
```

### Q: 如何实现分页？

```go
pageNo := c.DefaultQuery("pageNo", "1")
pageSize := c.DefaultQuery("pageSize", "10")
offset := (pageNo - 1) * pageSize

users, err := s.q.SystemUser.WithContext(ctx).
    Offset(offset).
    Limit(pageSize).
    Find()

count, _ := s.q.SystemUser.WithContext(ctx).Count()

core.WritePage(c, count, users)
```

### Q: 如何处理并发请求？

```go
// 使用 Redis 分布式锁
lock := core.RDB.SetNX(ctx, "lock:key", "1", time.Second)
if !lock.Val() {
    core.WriteError(c, core.ServerErrCode, "请求过于频繁，请稍后再试")
    return
}
defer core.RDB.Del(ctx, "lock:key")

// 执行业务逻辑
```

### Q: 如何缓存数据？

```go
cacheKey := fmt.Sprintf("user:%d", id)

// 先查缓存
val, err := core.RDB.Get(ctx, cacheKey).Result()
if err == nil {
    // 缓存命中
    return val, nil
}

// 查数据库
user, err := s.q.SystemUser.WithContext(ctx).First()

// 写入缓存
core.RDB.Set(ctx, cacheKey, user, time.Hour)

return user, nil
```

### Q: 如何实现软删除？

```go
// 在 model 中添加 DeletedAt 字段
type SystemUser struct {
    // ... 其他字段
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// GORM 会自动处理软删除
// 删除时只更新 deleted_at
s.q.SystemUser.WithContext(ctx).Delete(&user)

// 查询时自动排除已删除的记录
users, _ := s.q.SystemUser.WithContext(ctx).Find()
```

---

## 总结

本快速参考指南提供了：

✅ 常用命令和快速操作
✅ 常用代码片段和模式
✅ API 调用示例
✅ 数据库操作技巧
✅ 错误处理方式
✅ 调试技巧和工具
✅ 常见问题快速解决方案

在开发过程中，可以快速查阅本文档找到所需的代码片段和解决方案。

祝你开发愉快！🚀
