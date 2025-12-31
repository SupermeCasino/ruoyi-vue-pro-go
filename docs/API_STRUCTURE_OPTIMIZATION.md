# API 目录结构优化方案

## 问题分析

### 当前问题
1. **参数爆炸**: `InitRouter` 有 166 个 Handler 参数
2. **导入混乱**: 需要大量包别名 (`adminHandler`, `payAdmin`, `memberHandler`)
3. **类型分散**: req/resp 按模块分散,难以维护契约一致性
4. **缺乏抽象**: Handler 未按模块分组

## 优化方案: 分层模块化

### 核心理念
- **契约优先**: API 定义(req/resp)与实现(handler)分离
- **聚合注入**: 使用 Handler 聚合减少参数数量
- **接口抽象**: 利用 Go 接口实现松耦合

---

## 新目录结构

```
internal/api/
├── api/                  # API 契约层 (Request/Response)
│   ├── admin/
│   │   ├── mall/         # ✅ Admin Mall 模块契约
│   │   │   ├── product.go
│   │   │   └── trade.go
│   │   ├── system/       # ✅ Admin System 模块契约
│   │   │   └── user.go
│   │   └── pay/
│   └── app/
│       └── mall/
│
├── handler/              # Handler 实现层
│   ├── admin/
│   │   ├── handler.go    // type AdminHandlers struct { Mall, System, Pay ... }
│   │   ├── mall/         # ✅ Mall 模块
│   │   │   ├── product/
│   │   │   └── trade/
│   │   ├── system/       # ✅ System 模块
│   │   │   ├── user/
│   │   │   └── dept/
│   │   ├── pay/          # ✅ Pay 模块
│   │   └── member/
│   └── app/
│       ├── handler.go    // type AppHandlers struct { Mall, System ... }
│       └── mall/
│
└── router/
    ├── router.go         
    ├── admin.go
    └── app.go
```

---

## 实现示例

### 1. Handler 聚合 (handler/admin/handler.go)

```go
package admin

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/member"
)

// AdminHandlers 顶级聚合
type AdminHandlers struct {
	Mall    *mall.MallHandlers       // ✅ 明确的模块划分
	System  *system.SystemHandlers
	Pay     *pay.PayHandlers
	Member  *member.MemberHandlers
}

// NewAdminHandlers Wire 注入
func NewAdminHandlers(
    mallHandlers *mall.MallHandlers,
    systemHandlers *system.SystemHandlers,
    // ...
) *AdminHandlers {
    return &AdminHandlers{
        Mall:   mallHandlers,
        System: systemHandlers,
        // ...
    }
}
```

### 2. 子模块聚合 (handler/admin/mall/handler.go)

```go
package mall

type MallHandlers struct {
    Product    *ProductHandlers
    Trade      *TradeHandlers
    Promotion  *PromotionHandlers
}

type ProductHandlers struct {
    Spu     *product.SpuHandler
    Sku     *product.SkuHandler
    // ...
}
```

### 3. 调用方式对比

```go
// ❌ 旧方式
h.Product.Spu.GetPage(c)

// ✅ 新方式 (层级明确)
h.Mall.Product.Spu.GetPage(c)
h.System.User.GetPage(c)
```

**优势**:
- ✅ 模块化清晰:Member/Product/Trade 分组明确
- ✅ 类型安全:通过结构体字段访问,避免参数混淆
- ✅ 易于扩展:新增handler只需在对应模块结构体添加字段

---

### 3. 简化的 Router (router/router.go)

```go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"
)

// InitRouter 简化版 - 只需3个聚合参数!
func InitRouter(
	systemHandlers *system.SystemHandlers,
	adminHandlers *admin.AdminHandlers,
	appHandlers *app.AppHandlers,
	casbinMiddleware *middleware.CasbinMiddleware,
) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())
	
	// System 路由
	registerSystemRoutes(r, systemHandlers, casbinMiddleware)
	
	// Admin 路由
	registerAdminRoutes(r, adminHandlers, casbinMiddleware)
	
	// App 路由
	registerAppRoutes(r, appHandlers)
	
	return r
}
```

**优势**:
- ✅ 参数从 166个 → 4个!
- ✅ 可读性极大提升
- ✅ 易于测试和Mock

---

### 4. 模块化路由注册 (router/admin.go)

```go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"
)

func registerAdminRoutes(
	r *gin.Engine,
	h *admin.AdminHandlers,
	casbin *middleware.CasbinMiddleware,
) {
	adminAPI := r.Group("/admin-api")
	adminAPI.Use(casbin.CheckPermission())
	
	// Member 模块
	memberAPI := adminAPI.Group("/system/member")
	{
		memberAPI.GET("/user/page", h.Member.User.GetPage)
		memberAPI.POST("/user/create", h.Member.User.Create)
		memberAPI.PUT("/user/update", h.Member.User.Update)
		memberAPI.DELETE("/user/delete", h.Member.User.Delete)
		
		memberAPI.GET("/level/list", h.Member.Level.GetList)
		memberAPI.POST("/level/create", h.Member.Level.Create)
		// ...
	}
	
	// Product 模块
	productAPI := adminAPI.Group("/product")
	{
		productAPI.GET("/spu/page", h.Product.Spu.GetPage)
		productAPI.POST("/spu/create", h.Product.Spu.Create)
		// ...
	}
	
	// Trade 模块
共tradeAPI := adminAPI.Group("/trade")
	{
		tradeAPI.GET("/order/page", h.Trade.Order.GetPage)
		tradeAPI.POST("/order/delivery", h.Trade.Order.Delivery)
		// ...
	}
}
```

**优势**:
- ✅ 路由分组清晰,一目了然
- ✅ 通过 `h.Member.User.GetPage` 调用,语义明确
- ✅ 易于维护和重构

---

## 迁移策略

### Phase 1: 创建契约层(不破坏现有代码)
```bash
# 1. 创建 contract 目录
mkdir -p internal/api/contract/{admin,app}

# 2. 逐步迁移 req/resp 到契约层
# 保留原文件,新增import指向契约
```

### Phase 2: 创建Handler聚合
```bash
# 1. 创建 handler.go 聚合文件
# 2. 更新 wire.go 提供 AdminHandlers/AppHandlers
```

### Phase 3: 简化Router
```bash
# 1. 更新 router.go 使用聚合
# 2. 删除旧的独立参数
```

---

## 对比总结

| 维度 | 当前方案 | 优化方案 |
|------|---------|---------|
| Router 参数数量 | 166个 | 4个 |
| Handler 访问 | `payWalletHandler.GetPage(...)` | `adminHandlers.Pay.Wallet.GetPage(...)` |
| 新增 Handler | 需修改 InitRouter 签名 + wire.go | 只需修改聚合结构体 |
| API 契约查找 | 分散在73个req文件 | 集中在contract/{admin,app}下 |
| 包导入复杂度 | 需10+个别名 | 只需3个聚合包 |
| 测试复杂度 | 需mock 166个参数 | 只需mock 3个聚合 |

---

## 符合 Go 哲学的设计原则

1. **组合优于继承**: 使用结构体组合聚合handlers
2. **接口隔离**: 每个handler只关注自己的职责
3. **依赖注入**: Wire自动注入聚合结构
4. **清晰即文档**: 通过类型和命名自解释

这个方案让代码**简洁且符合Go惯用法**!
