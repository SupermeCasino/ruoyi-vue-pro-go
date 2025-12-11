# Backend Go - 芋道商城 Go 实现

## 项目简介

这是芋道商城（ruoyi-vue-pro）的 Go 语言实现版本，用于提供与 Java 实现完全对齐的 API 服务。项目采用 Go + Gin + GORM 技术栈，确保 API 返回结构、数据类型、逻辑实现与 Java 版本保持一致。

## 技术栈

- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT (本地验证)
- **依赖注入**: Wire

## 项目结构

```
backend-go/
├── cmd/
│   └── server/              # 服务启动入口
├── internal/
│   ├── api/
│   │   ├── handler/         # HTTP 请求处理器
│   │   ├── req/             # 请求对象 (VO)
│   │   ├── resp/            # 响应对象 (VO)
│   │   └── router/          # 路由定义
│   ├── middleware/          # 中间件 (鉴权、日志、错误处理等)
│   ├── model/               # 数据模型 (DO)
│   ├── service/             # 业务逻辑服务
│   ├── repository/          # 数据访问层
│   └── pkg/
│       ├── core/            # 核心包 (错误码、响应结构等)
│       └── utils/           # 工具函数 (JWT、加密等)
└── go.mod                   # Go 模块定义
```

---

## 鉴权机制

### 概述

Go 版本实现了与 Java 版本对齐的完整鉴权机制，支持用户类型区分和租户隔离。

### 认证流程

1. **Token 获取**
   - 支持三种方式获取 Token：
     - `Authorization: Bearer <token>` (Header)
     - `?Authorization=<token>` (Query Parameter)
     - `Authorization=<token>` (Form Parameter)

2. **Token 验证**
   - 使用 JWT 本地验证
   - 验证签名和过期时间
   - 提取用户信息

3. **用户信息存储**
   - 将完整的用户信息存储到 Gin Context
   - 支持在处理器中获取用户信息

### JWT Token 结构

```go
type Claims struct {
    UserID   int64  `json:"userId"`      // 用户 ID
    UserType int    `json:"userType"`    // 用户类型: 0=Member, 1=Admin
    TenantID int64  `json:"tenantId"`    // 租户 ID
    Nickname string `json:"nickname"`    // 用户昵称
    jwt.RegisteredClaims
}
```

### 使用示例

#### 生成 Token

```go
import "backend-go/internal/pkg/utils"

// 简单方式（仅包含 UserID）
token, err := utils.GenerateToken(userID, 24*time.Hour)

// 完整方式（包含所有用户信息）
token, err := utils.GenerateTokenWithInfo(
    userID,      // 用户 ID
    0,           // 用户类型 (0: Member, 1: Admin)
    tenantID,    // 租户 ID
    nickname,    // 用户昵称
    24*time.Hour, // 过期时间
)
```

#### 获取用户信息

```go
import "backend-go/internal/pkg/core"

// 获取用户 ID
userID := core.GetLoginUserID(c)

// 获取完整的用户信息
loginUser := core.GetLoginUser(c)
if loginUser != nil {
    userID := loginUser.UserID
    userType := loginUser.UserType
    tenantID := loginUser.TenantID
    nickname := loginUser.Nickname
}
```

### 鉴权中间件

在路由中使用鉴权中间件：

```go
import "backend-go/internal/middleware"

// 为特定路由组启用鉴权
authGroup := router.Group("/api/app")
authGroup.Use(middleware.Auth())
{
    // 需要鉴权的路由
    authGroup.POST("/cart/add", handler.AddCart)
    authGroup.GET("/order/list", handler.GetOrderPage)
}
```

---

## 错误码体系

### 错误码定义

Go 版本实现了完整的 HTTP 标准错误码体系，与 Java 版本对齐。

| 错误码 | 含义 | 使用场景 |
|--------|------|---------|
| `0` | 成功 | 请求成功 |
| `400` | 参数错误 | 请求参数验证失败 |
| `401` | 未授权 | 未登录或 Token 无效 |
| `403` | 禁止访问 | 无权限访问资源 |
| `404` | 资源不存在 | 请求的资源不存在 |
| `409` | 冲突 | 资源冲突（如重复创建） |
| `500` | 系统异常 | 服务器内部错误 |
| `501` | 未实现 | 功能未实现 |
| `503` | 服务不可用 | 服务暂时不可用 |

### 响应格式

所有 API 响应都遵循统一的格式：

#### 成功响应

```json
{
    "code": 0,
    "msg": "",
    "data": {
        // 实际数据
    }
}
```

#### 错误响应

```json
{
    "code": 400,
    "msg": "参数错误",
    "data": null
}
```

### 使用示例

#### 返回成功响应

```go
import "backend-go/internal/pkg/core"

// 返回数据
core.WriteSuccess(c, data)

// 或使用 Success 方法
c.JSON(200, core.Success(data))
```

#### 返回错误响应

```go
import "backend-go/internal/pkg/core"

// 参数错误
core.WriteError(c, core.ParamErrCode, "参数错误")

// 未授权
core.WriteError(c, core.UnauthorizedCode, "未登录")

// 禁止访问
core.WriteError(c, core.ForbiddenCode, "无权限访问")

// 资源不存在
core.WriteError(c, core.NotFoundCode, "资源不存在")

// 系统异常
core.WriteError(c, core.ServerErrCode, "系统异常")
```

### 错误码常量

```go
const (
    SuccessCode        = 0
    ParamErrCode       = 400
    UnauthorizedCode   = 401
    ForbiddenCode      = 403
    NotFoundCode       = 404
    ConflictCode       = 409
    ServerErrCode      = 500
    NotImplementCode   = 501
    ServiceUnavailCode = 503
)
```

---

## 中间件

### 已实现的中间件

#### 1. 鉴权中间件 (Auth)

```go
middleware.Auth()
```

- 验证 JWT Token
- 提取用户信息
- 支持三种 Token 获取方式

#### 2. 错误处理中间件 (ErrorHandler)

```go
middleware.ErrorHandler()
```

- 捕获业务错误
- 统一错误响应格式
- 记录错误日志

#### 3. 恢复中间件 (Recovery)

```go
middleware.Recovery()
```

- 捕获 panic
- 返回 500 错误响应
- 记录堆栈跟踪

#### 4. API 访问日志中间件 (APIAccessLogMiddleware)

```go
middleware.APIAccessLogMiddleware()
```

- 记录所有 API 访问
- 记录请求参数、请求体、响应体
- 清理敏感数据
- 异步记录日志

#### 5. 参数验证中间件 (ValidatorMiddleware)

```go
middleware.ValidatorMiddleware()
```

- 提供参数验证功能
- 与 Java 的 @Valid 注解对齐

### 中间件使用示例

```go
import "backend-go/internal/middleware"

// 全局中间件
router.Use(middleware.ErrorHandler())
router.Use(middleware.Recovery())
router.Use(middleware.APIAccessLogMiddleware())

// 路由组中间件
authGroup := router.Group("/api/app")
authGroup.Use(middleware.Auth())
```

---

## API 响应结构

### 通用响应结构

```go
type Result[T any] struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data T      `json:"data"`
}
```

### 分页响应结构

```go
type PageResult[T any] struct {
    List  []T   `json:"list"`
    Total int64 `json:"total"`
}
```

### 使用示例

```go
// 返回单个对象
core.WriteSuccess(c, user)

// 返回分页数据
pageResult := core.PageResult[User]{
    List:  users,
    Total: total,
}
core.WriteSuccess(c, pageResult)
```

---

## 参数验证

### 验证标签

使用 Gin 的 binding 标签进行参数验证：

```go
type AppCartAddReq struct {
    SkuID int64 `json:"skuId" binding:"required"`
    Count int   `json:"count" binding:"required,min=1"`
}
```

### 常用验证标签

| 标签 | 含义 |
|------|------|
| `required` | 必填 |
| `min=N` | 最小值 |
| `max=N` | 最大值 |
| `len=N` | 长度 |
| `email` | 邮箱格式 |
| `url` | URL 格式 |
| `dive` | 嵌套验证 |

### 验证示例

```go
var req AppCartAddReq
if err := c.ShouldBindJSON(&req); err != nil {
    core.WriteError(c, core.ParamErrCode, err.Error())
    return
}
```

---

## 与 Java 版本的对齐情况

### 已对齐项

- ✅ API 返回结构 (CommonResult, PageResult)
- ✅ 错误码体系 (HTTP 标准错误码)
- ✅ 鉴权机制 (JWT Token + 用户信息)
- ✅ 用户类型区分 (Member/Admin)
- ✅ 租户隔离 (TenantID)
- ✅ API 访问日志
- ✅ 参数验证
- ✅ Token 获取方式 (Header/Query/Form)

### 部分对齐项

- 🟡 API 端点 (缺少 4 个端点的完整实现)
- 🟡 VO/DO/BO 结构 (基本对齐，部分字段差异)

### 对齐度

**整体对齐度: 97%**

详见 `ALIGNMENT_VERIFICATION_REPORT.md`

---

## 快速开始

### 环境要求

- Go 1.20+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖

```bash
go mod download
```

### 配置文件

创建 `.env` 文件或设置环境变量：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=yudao

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=yudao-backend-go-secret
```

### 启动服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

---

## 常见问题

### 1. Token 过期如何处理？

返回 401 错误码，前端需要重新登录获取新 Token。

### 2. 如何区分用户类型？

通过 `loginUser.UserType` 字段：
- `0`: 普通用户 (Member)
- `1`: 管理员 (Admin)

### 3. 如何实现租户隔离？

在查询时使用 `loginUser.TenantID` 过滤数据：

```go
loginUser := core.GetLoginUser(c)
orders := querySvc.GetOrdersByTenant(c, loginUser.TenantID)
```

### 4. 如何添加新的错误码？

在 `internal/pkg/core/error.go` 中添加常量和错误变量：

```go
const NewErrorCode = 4xx

var ErrNewError = NewBizError(NewErrorCode, "错误描述")
```

---

## 相关文档

- [对齐检查清单](./ALIGNMENT_CHECKLIST.md) - 详细的对齐检查项
- [修复总结](./ALIGNMENT_FIX_SUMMARY.md) - 修复内容和对比
- [验证报告](./ALIGNMENT_VERIFICATION_REPORT.md) - 自查验证结果

---

## 贡献指南

1. 确保代码与 Java 版本对齐
2. 遵循现有的代码风格
3. 添加必要的注释和文档
4. 提交前运行测试

---

## 许可证

MIT License

---

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。

