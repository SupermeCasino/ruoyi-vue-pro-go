# 订单金额计算器系统 - 深度教学文档

## 目录
1. [系统概览](#系统概览)
2. [核心设计模式](#核心设计模式)
3. [架构设计](#架构设计)
4. [工作流程](#工作流程)
5. [关键组件详解](#关键组件详解)
6. [设计模式应用](#设计模式应用)
7. [最佳实践](#最佳实践)
8. [常见场景](#常见场景)

---

## 系统概览

### 什么是订单金额计算器？

订单金额计算器是一个**复杂的、可扩展的、分层的价格计算系统**，用于在电商平台中精确计算订单的最终支付金额。它需要考虑：

- ✅ 商品原价
- ✅ 优惠券折扣
- ✅ 限时折扣活动
- ✅ VIP会员折扣
- ✅ 积分抵扣
- ✅ 运费计算
- ✅ 满减送活动
- ✅ 秒杀/拼团/砍价等特殊订单类型

### 为什么需要这样的系统？

**问题背景：** 如果直接在一个方法中计算所有这些因素，代码会变成：

```go
// ❌ 反面例子：单一巨大方法
func CalculatePrice(order Order) Price {
    // 1000+ 行代码混合在一起
    // 添加新的折扣类型需要修改这个方法
    // 难以测试、难以维护、难以扩展
}
```

**解决方案：** 使用**策略模式 + 工厂模式 + 助手模式**，将计算逻辑分离成独立的、可组合的计算器。

---

## 核心设计模式

### 1. 策略模式（Strategy Pattern）

**定义：** 定义一族算法，将每一个算法封装起来，并让它们可以相互替换。

#### 在本系统中的应用：

```go
// 策略接口 - 定义所有计算器必须实现的契约
type PriceCalculator interface {
    Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error
    GetOrder() int      // 执行顺序（优先级）
    GetName() string    // 计算器名称
    IsApplicable(orderType int) bool  // 是否适用于当前订单类型
}
```

**为什么使用策略模式？**

| 优势 | 说明 |
|------|------|
| **开闭原则** | 新增计算器无需修改现有代码，只需实现接口 |
| **单一职责** | 每个计算器只负责一种折扣计算 |
| **易于测试** | 可以独立测试每个计算器 |
| **运行时选择** | 根据订单类型动态选择适用的计算器 |
| **灵活组合** | 可以灵活组合多个计算器 |

#### 具体实现示例：

```go
// 优惠券计算器
type CouponPriceCalculator struct {
    BasePriceCalculator
    couponSvc *promotion.CouponService
}

func (c *CouponPriceCalculator) Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error {
    // 只处理普通订单
    if resp.Type != consts.TradeOrderTypeNormal {
        return nil
    }
    // 计算优惠券折扣...
}

func (c *CouponPriceCalculator) GetOrder() int {
    return 10  // 优先级：10
}

func (c *CouponPriceCalculator) GetName() string {
    return "CouponCalculator"
}

func (c *CouponPriceCalculator) IsApplicable(orderType int) bool {
    return orderType == consts.TradeOrderTypeNormal
}
```

---

### 2. 工厂模式（Factory Pattern）

**定义：** 定义一个接口来创建对象，但让子类决定实例化哪个类。

#### 在本系统中的应用：

```go
type PriceCalculatorFactory struct {
    calculators []PriceCalculator
    logger      *zap.Logger
}

// 注册计算器
func (f *PriceCalculatorFactory) RegisterCalculator(calculator PriceCalculator) {
    f.calculators = append(f.calculators, calculator)
}

// 获取所有计算器（按优先级排序）
func (f *PriceCalculatorFactory) GetCalculators() []PriceCalculator {
    // 排序并返回
}

// 获取适用于特定订单类型的计算器
func (f *PriceCalculatorFactory) GetApplicableCalculators(orderType int) []PriceCalculator {
    // 过滤并返回
}
```

**为什么使用工厂模式？**

- **集中管理：** 所有计算器的创建和注册在一个地方
- **解耦创建逻辑：** 调用方不需要知道如何创建计算器
- **动态注册：** 可以在运行时动态注册新的计算器
- **统一初始化：** 确保所有计算器都被正确初始化

---

### 3. 模板方法模式（Template Method Pattern）

**定义：** 在基类中定义算法的骨架，让子类实现具体步骤。

#### 在本系统中的应用：

```go
// 基础计算器 - 提供通用功能
type BasePriceCalculator struct {
    name   string
    order  int
    Helper *PriceCalculatorHelper
    logger *zap.Logger
}

// 模板方法 - 子类可以复用
func (b *BasePriceCalculator) GetName() string {
    return b.name
}

func (b *BasePriceCalculator) GetOrder() int {
    return b.order
}

// 日志记录模板
func (b *BasePriceCalculator) LogCalculation(ctx context.Context, req *TradePriceCalculateReqBO, message string, fields ...zap.Field) {
    // 统一的日志格式
}
```

**为什么使用模板方法模式？**

- **代码复用：** 所有计算器共享通用的日志、验证等逻辑
- **一致性：** 确保所有计算器遵循相同的执行流程
- **易于扩展：** 新计算器只需实现 `Calculate` 方法

---

### 4. 助手模式（Helper/Utility Pattern）

**定义：** 将通用的、可复用的工具方法集中在一个助手类中。

#### 在本系统中的应用：

```go
type PriceCalculatorHelper struct {
    logger *zap.Logger
}

// 通用工具方法
func (h *PriceCalculatorHelper) DividePrice(items []TradePriceCalculateItemRespBO, totalDiscount int) []int {
    // 按支付金额比例分摊折扣
}

func (h *PriceCalculatorHelper) RecountPayPrice(item *TradePriceCalculateItemRespBO) {
    // 重新计算支付金额
}

func (h *PriceCalculatorHelper) CalculateRatePrice(price int, discountPercent int) int {
    // 计算折扣价格
}

func (h *PriceCalculatorHelper) UpdateResponsePrice(resp *TradePriceCalculateRespBO) {
    // 更新响应对象的价格信息
}
```

**为什么使用助手模式？**

- **避免重复：** 多个计算器都需要分摊折扣、重新计算价格等
- **集中维护：** 修改计算逻辑只需改一个地方
- **提高可读性：** 复杂的计算逻辑被封装成有意义的方法名

---

## 架构设计

### 系统分层

```
┌─────────────────────────────────────────────────────────────┐
│                    API 层（Handler）                         │
│              接收请求，调用 Service 层                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│              TradePriceService（服务层）                     │
│  - 协调整个计算流程                                          │
│  - 管理计算器执行顺序                                        │
│  - 参数验证和结果验证                                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
┌───────▼──────┐ ┌────▼─────┐ ┌─────▼──────────┐
│ 计算器1      │ │ 计算器2  │ │ 计算器3        │
│ (优惠券)     │ │ (限时折扣)│ │ (VIP折扣)      │
└───────┬──────┘ └────┬─────┘ └─────┬──────────┘
        │              │              │
        └──────────────┼──────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│         PriceCalculatorHelper（助手类）                      │
│  - 分摊折扣                                                  │
│  - 重新计算价格                                              │
│  - 通用工具方法                                              │
└──────────────────────────────────────────────────────────────┘
```

### 数据流向

```
请求 (TradePriceCalculateReqBO)
  │
  ├─ 参数验证
  │
  ├─ 初始化响应对象
  │
  ├─ 构建商品项响应
  │  ├─ 获取 SKU 信息
  │  ├─ 获取 SPU 信息
  │  └─ 验证库存
  │
  ├─ 按优先级执行计算器
  │  ├─ 计算器1：优惠券折扣
  │  ├─ 计算器2：限时折扣
  │  ├─ 计算器3：VIP折扣
  │  ├─ 计算器4：积分抵扣
  │  └─ 计算器N：...
  │
  ├─ 更新最终价格信息
  │
  ├─ 结果验证
  │
  └─ 返回响应 (TradePriceCalculateRespBO)
```

---

## 工作流程

### 完整的计算流程

#### 第一步：参数验证

```go
func (s *TradePriceService) validateRequest(req *TradePriceCalculateReqBO) error {
    // 1. 验证用户ID
    if req.UserID <= 0 {
        return pkgErrors.NewBizError(1004003001, "用户ID不能为空")
    }
    
    // 2. 验证商品项
    if len(req.Items) == 0 {
        return pkgErrors.NewBizError(1004003001, "计算价格时，商品不能为空")
    }
    
    // 3. 验证每个商品项
    for _, item := range req.Items {
        if item.SkuID <= 0 {
            return pkgErrors.NewBizError(1004003001, "商品SKU ID不能为空")
        }
        if item.Count <= 0 {
            return pkgErrors.NewBizError(1004003001, "商品数量必须大于0")
        }
    }
    
    return nil
}
```

**为什么要参数验证？**
- 防止垃圾数据进入计算流程
- 提前发现问题，避免后续步骤出错
- 提供清晰的错误提示

---

#### 第二步：初始化响应对象

```go
func (h *PriceCalculatorHelper) BuildCalculateResp(req *TradePriceCalculateReqBO) *TradePriceCalculateRespBO {
    resp := &TradePriceCalculateRespBO{
        Type:       tradeModel.TradeOrderTypeNormal,  // 默认为普通订单
        Price:      TradePriceCalculatePriceBO{},
        Items:      make([]TradePriceCalculateItemRespBO, 0),
        Success:    true,
        Coupons:    make([]TradePriceCalculateCouponBO, 0),
        Promotions: make([]TradePriceCalculatePromotionBO, 0),
    }
    
    // 根据请求参数确定订单类型
    if req.SeckillActivityId > 0 {
        resp.Type = tradeModel.TradeOrderTypeSeckill
    } else if req.BargainRecordId > 0 {
        resp.Type = tradeModel.TradeOrderTypeBargain
    } else if req.CombinationActivityId > 0 {
        resp.Type = tradeModel.TradeOrderTypeCombination
    } else if req.PointActivityId > 0 {
        resp.Type = tradeModel.TradeOrderTypePoint
    }
    
    return resp
}
```

**关键点：**
- 订单类型决定了哪些计算器会被执行
- 不同订单类型有不同的计算规则

---

#### 第三步：构建商品项响应

这是一个**关键步骤**，需要：

1. **批量获取 SKU 信息**
   ```go
   skuList, err := s.skuSvc.GetSkuList(ctx, skuIDs)
   ```
   - 获取商品的价格、库存、图片等信息

2. **批量获取 SPU 信息**
   ```go
   spuList, err := s.spuSvc.GetSpuList(ctx, spuIDs)
   ```
   - 获取商品的名称、分类、配送方式等信息

3. **验证商品状态**
   - 商品是否已上架
   - 库存是否充足

4. **构建响应商品项**
   ```go
   item := TradePriceCalculateItemRespBO{
       SpuID:      spu.ID,
       SkuID:      sku.ID,
       Count:      reqItem.Count,
       Price:      sku.Price,           // 原价
       PayPrice:   sku.Price * reqItem.Count,  // 初始支付金额
       SpuName:    spu.Name,
       PicURL:     sku.PicURL,
       // ... 其他字段
   }
   ```

**为什么这么复杂？**
- 需要验证商品的有效性
- 需要获取后续计算所需的基础信息
- 需要检查库存，防止超卖

---

#### 第四步：按顺序执行计算器

这是**核心步骤**，体现了策略模式的威力：

```go
// 4. 按顺序执行计算器
for _, calculator := range s.calculators {
    // 检查计算器是否适用于当前订单类型
    if !calculator.IsApplicable(resp.Type) {
        s.logger.Debug("跳过不适用的计算器",
            zap.String("calculator", calculator.GetName()),
            zap.Int("orderType", resp.Type),
        )
        continue
    }
    
    s.logger.Info("执行价格计算器",
        zap.String("calculator", calculator.GetName()),
        zap.Int("order", calculator.GetOrder()),
    )
    
    // 执行计算
    if err := calculator.Calculate(ctx, req, resp); err != nil {
        s.logger.Error("价格计算器执行失败",
            zap.String("calculator", calculator.GetName()),
            zap.Error(err),
        )
        return nil, err
    }
}
```

**计算器执行顺序很重要：**

| 顺序 | 计算器 | 说明 |
|------|--------|------|
| 1 | SKU 优惠计算器 | 计算 SKU 级别的优惠（限时折扣、VIP折扣等） |
| 2 | 优惠券计算器 | 计算优惠券折扣 |
| 3 | 积分计算器 | 计算积分抵扣 |
| 4 | 运费计算器 | 计算运费 |
| 5 | 满减送计算器 | 计算满减送活动 |

**为什么顺序很重要？**
- 优惠券通常是在所有折扣之后计算的
- 积分抵扣应该在运费之前
- 满减送应该最后计算，因为它可能触发额外的优惠

---

#### 第五步：更新最终价格信息

```go
func (h *PriceCalculatorHelper) UpdateResponsePrice(resp *TradePriceCalculateRespBO) {
    // 重新计算总价格
    totalPrice := h.CalculateTotalPrice(resp.Items, true)
    totalPayPrice := h.CalculateTotalPayPrice(resp.Items, true)
    
    // 计算各种折扣总额
    totalDiscountPrice := 0
    totalCouponPrice := 0
    totalPointPrice := 0
    totalVipPrice := 0
    totalDeliveryPrice := 0
    
    for _, item := range resp.Items {
        if !item.Selected {
            continue
        }
        totalDiscountPrice += item.DiscountPrice
        totalCouponPrice += item.CouponPrice
        totalPointPrice += item.PointPrice
        totalVipPrice += item.VipPrice
        totalDeliveryPrice += item.DeliveryPrice
    }
    
    // 更新价格信息
    resp.Price.TotalPrice = totalPrice
    resp.Price.DiscountPrice = totalDiscountPrice
    resp.Price.CouponPrice = totalCouponPrice
    resp.Price.PointPrice = totalPointPrice
    resp.Price.VipPrice = totalVipPrice
    resp.Price.DeliveryPrice = totalDeliveryPrice
    resp.Price.PayPrice = totalPayPrice
}
```

**这一步做什么？**
- 汇总所有商品项的价格信息
- 计算各种折扣的总额
- 生成最终的支付金额

---

#### 第六步：结果验证

```go
func (s *TradePriceService) validateResponse(req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error {
    // 验证商品项数量
    if len(resp.Items) == 0 {
        return pkgErrors.NewBizError(1004003001, "商品项不能为空")
    }
    
    // 验证支付金额（积分订单允许支付金额为0）
    if req.PointActivityId == 0 && resp.Price.PayPrice <= 0 {
        return pkgErrors.NewBizError(1004003004, "支付金额不合法")
    }
    
    return nil
}
```

**为什么需要结果验证？**
- 确保计算结果的合理性
- 防止异常数据返回给客户端
- 提前发现计算逻辑的问题

---

## 关键组件详解

### 1. TradePriceService（服务层）

**职责：**
- 协调整个计算流程
- 管理计算器的执行顺序
- 参数和结果验证
- 日志记录

**关键方法：**

```go
// 计算订单价格（主入口）
func (s *TradePriceService) CalculateOrderPrice(ctx context.Context, req *TradePriceCalculateReqBO) (*TradePriceCalculateRespBO, error)

// 内部计算逻辑
func (s *TradePriceService) calculatePriceInternal(ctx context.Context, req *TradePriceCalculateReqBO, checkStock bool) (*TradePriceCalculateRespBO, error)

// 构建商品项响应
func (s *TradePriceService) buildItemsResponse(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO, checkStock bool) error

// 获取适用的计算器
func (s *TradePriceService) GetApplicableCalculators(orderType int) []PriceCalculator
```

---

### 2. PriceCalculator 接口

**定义：**
```go
type PriceCalculator interface {
    Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error
    GetOrder() int
    GetName() string
    IsApplicable(orderType int) bool
}
```

**实现示例：优惠券计算器**

```go
type CouponPriceCalculator struct {
    BasePriceCalculator
    couponSvc *promotion.CouponService
}

func (c *CouponPriceCalculator) Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error {
    // 只有【普通】订单，才允许使用优惠劵
    if resp.Type != consts.TradeOrderTypeNormal {
        if req.CouponID != nil && *req.CouponID > 0 {
            return pkgErrors.NewBizError(1004001004, "优惠券仅限普通订单使用")
        }
        return nil
    }
    
    // 1. 加载用户的优惠劵列表
    coupons, err := c.couponSvc.GetUnusedCouponList(ctx, req.UserID)
    if err != nil {
        c.LogError(ctx, req, err, "获取用户优惠券列表失败")
        return err
    }
    
    // 2. 过滤过期的优惠券
    now := time.Now()
    coupons = lo.Filter(coupons, func(coupon *promotionModel.PromotionCoupon, _ int) bool {
        return !now.After(coupon.ValidEndTime)
    })
    
    // 3. 计算优惠劵的使用条件
    resp.Coupons = c.calculateCoupons(coupons, resp)
    
    // 4. 校验优惠劵是否可用
    if req.CouponID == nil || *req.CouponID <= 0 {
        return nil
    }
    
    // ... 计算优惠券折扣
    
    return nil
}

func (c *CouponPriceCalculator) GetOrder() int {
    return 10  // 优先级
}

func (c *CouponPriceCalculator) GetName() string {
    return "CouponCalculator"
}

func (c *CouponPriceCalculator) IsApplicable(orderType int) bool {
    return orderType == consts.TradeOrderTypeNormal
}
```

---

### 3. PriceCalculatorHelper（助手类）

**核心方法：**

#### a) DividePrice - 按比例分摊折扣

```go
func (h *PriceCalculatorHelper) DividePrice(items []TradePriceCalculateItemRespBO, totalDiscount int) []int {
    // 计算所有【已选中】项的总支付金额
    totalPayPrice := 0
    for _, item := range items {
        if item.Selected {
            totalPayPrice += item.PayPrice
        }
    }
    
    if totalPayPrice == 0 {
        return make([]int, len(items))
    }
    
    // 按比例分摊
    dividedPrices := make([]int, len(items))
    remainPrice := totalDiscount
    lastSelectedIndex := -1
    
    // 找到最后一个选中项的索引
    for i := len(items) - 1; i >= 0; i-- {
        if items[i].Selected {
            lastSelectedIndex = i
            break
        }
    }
    
    for i := 0; i < len(items); i++ {
        if !items[i].Selected {
            dividedPrices[i] = 0
            continue
        }
        
        if i < lastSelectedIndex {
            // 前 n-1 项按比例计算
            dividedPrices[i] = int(int64(totalDiscount) * int64(items[i].PayPrice) / int64(totalPayPrice))
            remainPrice -= dividedPrices[i]
        } else {
            // 最后一项用剩余金额（避免舍入误差）
            dividedPrices[i] = remainPrice
        }
    }
    
    return dividedPrices
}
```

**为什么这样设计？**

假设有 3 个商品，总优惠 100 元：
- 商品1：支付金额 200 元
- 商品2：支付金额 300 元
- 商品3：支付金额 500 元
- 总支付金额：1000 元

分摊结果：
- 商品1：100 × (200/1000) = 20 元
- 商品2：100 × (300/1000) = 30 元
- 商品3：100 - 20 - 30 = 50 元（使用剩余金额，避免舍入误差）

**关键点：**
- 最后一项使用剩余金额，确保总额精确
- 避免浮点数舍入导致的误差

---

#### b) RecountPayPrice - 重新计算支付金额

```go
func (h *PriceCalculatorHelper) RecountPayPrice(item *TradePriceCalculateItemRespBO) {
    // PayPrice = Price * Count - DiscountPrice + DeliveryPrice - CouponPrice - PointPrice - VipPrice
    item.PayPrice = item.Price*item.Count - item.DiscountPrice + item.DeliveryPrice - item.CouponPrice - item.PointPrice - item.VipPrice
    if item.PayPrice < 0 {
        item.PayPrice = 0
    }
}
```

**公式解析：**

```
支付金额 = 原价 - 活动折扣 + 运费 - 优惠券 - 积分 - VIP折扣
```

**为什么这样设计？**
- 每次应用一个折扣后，都需要重新计算支付金额
- 确保各种折扣的累积效果正确

---

#### c) CalculateRatePrice - 计算折扣价

```go
func (h *PriceCalculatorHelper) CalculateRatePrice(price int, discountPercent int) int {
    if discountPercent <= 0 {
        return 0
    }
    // 计算折扣价：price * discountPercent / 100
    targetPrice := (price * discountPercent) / 100
    return targetPrice
}
```

**示例：**
- 原价：1000 元
- 折扣：80%（即 discountPercent = 80）
- 折扣后价格：1000 × 80 / 100 = 800 元

---

### 4. 业务对象（BO）

#### TradePriceCalculateReqBO - 请求对象

```go
type TradePriceCalculateReqBO struct {
    UserID                int64                       // 用户ID
    CouponID              *int64                      // 优惠券ID
    PointStatus           bool                        // 是否使用积分
    DeliveryType          int                         // 配送方式
    AddressID             *int64                      // 收货地址ID
    PickUpStoreID         *int64                      // 自提门店ID
    SeckillActivityId     int64                       // 秒杀活动ID
    CombinationActivityId int64                       // 拼团活动ID
    CombinationHeadId     int64                       // 拼团团长ID
    BargainRecordId       int64                       // 砍价记录ID
    PointActivityId       int64                       // 积分活动ID
    CartIDs               []int64                     // 购物车ID数组
    Items                 []TradePriceCalculateItemBO // 商品项数组
}
```

**设计特点：**
- 使用指针表示可选字段（CouponID、AddressID 等）
- 支持多种订单类型（秒杀、拼团、砍价等）
- Items 是核心数据，包含要计算的商品

---

#### TradePriceCalculateRespBO - 响应对象

```go
type TradePriceCalculateRespBO struct {
    Type       int                              // 订单类型
    Price      TradePriceCalculatePriceBO       // 价格信息
    Items      []TradePriceCalculateItemRespBO  // 商品项数组
    CouponID   int64                            // 使用的优惠券ID
    TotalPoint int                              // 用户总积分
    UsePoint   int                              // 使用的积分
    GivePoint  int                              // 赠送的积分
    Success    bool                             // 计算是否成功
    Coupons    []TradePriceCalculateCouponBO    // 可用优惠券数组
    Promotions []TradePriceCalculatePromotionBO // 营销活动数组
}
```

**设计特点：**
- 包含完整的价格信息（Price）
- 包含所有商品项的详细信息（Items）
- 包含可用优惠券列表（Coupons）
- 包含应用的促销活动明细（Promotions）

---

#### TradePriceCalculatePriceBO - 价格信息

```go
type TradePriceCalculatePriceBO struct {
    TotalPrice    int // 总价格（原价）
    DiscountPrice int // 折扣金额（活动折扣）
    DeliveryPrice int // 运费
    CouponPrice   int // 优惠券折扣
    PointPrice    int // 积分抵扣
    VipPrice      int // VIP折扣
    PayPrice      int // 应付金额
}
```

**价格关系：**

```
PayPrice = TotalPrice - DiscountPrice + DeliveryPrice - CouponPrice - PointPrice - VipPrice
```

---

## 设计模式应用

### 1. 策略模式的优势

**场景：添加新的折扣类型**

❌ **不使用策略模式的做法：**

```go
func CalculatePrice(order Order) Price {
    // ... 100+ 行代码
    if order.HasCoupon {
        // 计算优惠券
    }
    if order.HasDiscount {
        // 计算限时折扣
    }
    if order.HasVip {
        // 计算VIP折扣
    }
    if order.HasNewUserCoupon {
        // 计算新用户优惠券
    }
    // ... 继续添加新的折扣类型
}
```

**问题：**
- 每次添加新的折扣类型都要修改这个方法
- 方法变得越来越大，难以维护
- 难以测试单个折扣类型
- 违反开闭原则

✅ **使用策略模式的做法：**

```go
// 1. 创建新的计算器
type NewUserCouponCalculator struct {
    BasePriceCalculator
    // ... 依赖注入
}

func (c *NewUserCouponCalculator) Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error {
    // 计算新用户优惠券
    return nil
}

func (c *NewUserCouponCalculator) GetOrder() int {
    return 15  // 优先级
}

func (c *NewUserCouponCalculator) GetName() string {
    return "NewUserCouponCalculator"
}

func (c *NewUserCouponCalculator) IsApplicable(orderType int) bool {
    return orderType == consts.TradeOrderTypeNormal
}

// 2. 在初始化时注册
calculators := []PriceCalculator{
    // ... 其他计算器
    NewNewUserCouponCalculator(...),
}

service := NewTradePriceService(calculators, ...)
```

**优势：**
- 无需修改现有代码
- 新的计算器是独立的、可测试的
- 符合开闭原则

---

### 2. 工厂模式的优势

**场景：根据订单类型选择计算器**

❌ **不使用工厂模式的做法：**

```go
func GetCalculators(orderType int) []PriceCalculator {
    var calculators []PriceCalculator
    
    if orderType == TradeOrderTypeNormal {
        calculators = append(calculators, &CouponCalculator{})
        calculators = append(calculators, &DiscountCalculator{})
        calculators = append(calculators, &VipCalculator{})
    } else if orderType == TradeOrderTypeSeckill {
        calculators = append(calculators, &SeckillCalculator{})
    } else if orderType == TradeOrderTypeBargain {
        calculators = append(calculators, &BargainCalculator{})
    }
    
    return calculators
}
```

**问题：**
- 创建逻辑分散在各处
- 难以维护计算器的注册信息
- 难以添加新的计算器

✅ **使用工厂模式的做法：**

```go
type PriceCalculatorFactory struct {
    calculators []PriceCalculator
}

func (f *PriceCalculatorFactory) RegisterCalculator(calculator PriceCalculator) {
    f.calculators = append(f.calculators, calculator)
}

func (f *PriceCalculatorFactory) GetApplicableCalculators(orderType int) []PriceCalculator {
    applicable := make([]PriceCalculator, 0)
    for _, calculator := range f.calculators {
        if calculator.IsApplicable(orderType) {
            applicable = append(applicable, calculator)
        }
    }
    return applicable
}

// 初始化时
factory := NewPriceCalculatorFactory(logger)
factory.RegisterCalculator(NewCouponCalculator(...))
factory.RegisterCalculator(NewDiscountCalculator(...))
factory.RegisterCalculator(NewVipCalculator(...))
factory.RegisterCalculator(NewSeckillCalculator(...))
factory.RegisterCalculator(NewBargainCalculator(...))
```

**优势：**
- 创建逻辑集中在一个地方
- 易于添加新的计算器
- 易于管理计算器的优先级

---

### 3. 助手模式的优势

**场景：多个计算器都需要分摊折扣**

❌ **不使用助手模式的做法：**

```go
// CouponCalculator 中
func (c *CouponCalculator) Calculate(...) error {
    // 分摊优惠券折扣
    dividedPrices := make([]int, len(items))
    totalPayPrice := 0
    for _, item := range items {
        if item.Selected {
            totalPayPrice += item.PayPrice
        }
    }
    // ... 分摊逻辑
}

// DiscountCalculator 中
func (c *DiscountCalculator) Calculate(...) error {
    // 分摊活动折扣
    dividedPrices := make([]int, len(items))
    totalPayPrice := 0
    for _, item := range items {
        if item.Selected {
            totalPayPrice += item.PayPrice
        }
    }
    // ... 分摊逻辑（重复）
}
```

**问题：**
- 代码重复
- 修改分摊逻辑需要改多个地方
- 容易出现不一致

✅ **使用助手模式的做法：**

```go
// PriceCalculatorHelper 中
func (h *PriceCalculatorHelper) DividePrice(items []TradePriceCalculateItemRespBO, totalDiscount int) []int {
    // 统一的分摊逻辑
    // ...
}

// CouponCalculator 中
func (c *CouponCalculator) Calculate(...) error {
    dividedPrices := c.Helper.DividePrice(items, couponPrice)
    // ...
}

// DiscountCalculator 中
func (c *DiscountCalculator) Calculate(...) error {
    dividedPrices := c.Helper.DividePrice(items, discountPrice)
    // ...
}
```

**优势：**
- 代码复用
- 修改逻辑只需改一个地方
- 确保一致性

---

## 最佳实践

### 1. 计算器的优先级设计

**原则：** 优先级应该反映折扣的应用顺序

```go
const (
    OrderSkuPromotion = 1    // SKU 级别的优惠（最先）
    OrderCoupon       = 10   // 优惠券
    OrderPoint        = 20   // 积分
    OrderDelivery     = 30   // 运费
    OrderReward       = 40   // 满减送（最后）
)
```

**为什么？**
- SKU 级别的优惠应该最先应用（因为它影响商品的基础价格）
- 优惠券通常在其他折扣之后应用
- 运费应该在折扣之后计算
- 满减送应该最后，因为它可能触发额外的优惠

---

### 2. 计算器的职责分离

**原则：** 每个计算器只负责一种折扣类型

```go
// ✅ 好的设计
type CouponPriceCalculator struct {
    // 只负责优惠券折扣
}

type DiscountActivityPriceCalculator struct {
    // 只负责限时折扣
}

type VipPriceCalculator struct {
    // 只负责VIP折扣
}

// ❌ 不好的设计
type ComplexCalculator struct {
    // 同时负责优惠券、折扣、VIP等多种折扣
    // 这违反了单一职责原则
}
```

---

### 3. 错误处理

**原则：** 计算器中的错误应该立即返回，不应该继续执行

```go
func (s *TradePriceService) calculatePriceInternal(...) (*TradePriceCalculateRespBO, error) {
    // ...
    
    for _, calculator := range s.calculators {
        if err := calculator.Calculate(ctx, req, resp); err != nil {
            s.logger.Error("价格计算器执行失败",
                zap.String("calculator", calculator.GetName()),
                zap.Error(err),
            )
            return nil, err  // 立即返回错误
        }
    }
    
    // ...
}
```

**为什么？**
- 如果一个计算器出错，后续的计算可能会基于错误的数据
- 立即返回错误可以防止数据污染

---

### 4. 日志记录

**原则：** 记录关键的计算步骤，便于调试和审计

```go
func (c *CouponCalculator) Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error {
    c.LogCalculation(ctx, req, "开始执行优惠券价格计算",
        zap.Int64("couponId", *req.CouponID),
    )
    
    // ... 计算逻辑
    
    c.LogCalculation(ctx, req, "分摊优惠券折扣",
        zap.Int64("skuId", resp.Items[i].SkuID),
        zap.Int("dividedCouponPrice", discount),
    )
    
    // ...
}
```

**记录什么？**
- 计算器的开始和结束
- 关键的中间结果
- 错误和异常

---

### 5. 参数验证

**原则：** 在计算开始前验证所有参数

```go
func (s *TradePriceService) validateRequest(req *TradePriceCalculateReqBO) error {
    if req.UserID <= 0 {
        return pkgErrors.NewBizError(1004003001, "用户ID不能为空")
    }
    
    if len(req.Items) == 0 {
        return pkgErrors.NewBizError(1004003001, "计算价格时，商品不能为空")
    }
    
    for _, item := range req.Items {
        if item.SkuID <= 0 {
            return pkgErrors.NewBizError(1004003001, "商品SKU ID不能为空")
        }
        if item.Count <= 0 {
            return pkgErrors.NewBizError(1004003001, "商品数量必须大于0")
        }
    }
    
    return nil
}
```

---

### 6. 结果验证

**原则：** 在返回结果前验证结果的合理性

```go
func (s *TradePriceService) validateResponse(req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error {
    if len(resp.Items) == 0 {
        return pkgErrors.NewBizError(1004003001, "商品项不能为空")
    }
    
    if req.PointActivityId == 0 && resp.Price.PayPrice <= 0 {
        return pkgErrors.NewBizError(1004003004, "支付金额不合法")
    }
    
    return nil
}
```

---

## 常见场景

### 场景 1：计算普通订单的价格

```go
req := &TradePriceCalculateReqBO{
    UserID: 123,
    Items: []TradePriceCalculateItemBO{
        {
            SkuID:    456,
            Count:    2,
            Selected: true,
        },
    },
    CouponID: ptr.Int64(789),  // 使用优惠券
    PointStatus: true,          // 使用积分
    DeliveryType: 1,            // 快递配送
}

resp, err := priceService.CalculateOrderPrice(ctx, req)
if err != nil {
    // 处理错误
}

// resp.Price.PayPrice 就是最终的支付金额
```

**执行流程：**
1. 验证参数
2. 初始化响应对象（Type = TradeOrderTypeNormal）
3. 构建商品项响应
4. 执行计算器：
   - SKU 优惠计算器（限时折扣、VIP折扣）
   - 优惠券计算器
   - 积分计算器
   - 运费计算器
   - 满减送计算器
5. 更新最终价格信息
6. 验证结果

---

### 场景 2：计算秒杀订单的价格

```go
req := &TradePriceCalculateReqBO{
    UserID: 123,
    Items: []TradePriceCalculateItemBO{
        {
            SkuID:    456,
            Count:    1,
            Selected: true,
        },
    },
    SeckillActivityId: 999,  // 秒杀活动ID
}

resp, err := priceService.CalculateOrderPrice(ctx, req)
```

**执行流程：**
1. 验证参数
2. 初始化响应对象（Type = TradeOrderTypeSeckill）
3. 构建商品项响应
4. 执行计算器：
   - 只有 IsApplicable(TradeOrderTypeSeckill) = true 的计算器会被执行
   - 优惠券计算器会被跳过（因为秒杀订单不支持优惠券）
5. 更新最终价格信息
6. 验证结果

---

### 场景 3：添加新的折扣类型

假设需要添加"满赠优惠券"的折扣类型：

```go
// 1. 创建新的计算器
type FullGiftCouponCalculator struct {
    BasePriceCalculator
    couponSvc *promotion.CouponService
}

func (c *FullGiftCouponCalculator) Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error {
    // 只处理普通订单
    if resp.Type != consts.TradeOrderTypeNormal {
        return nil
    }
    
    // 计算满赠优惠券
    // ...
    
    return nil
}

func (c *FullGiftCouponCalculator) GetOrder() int {
    return 45  // 在满减送之后
}

func (c *FullGiftCouponCalculator) GetName() string {
    return "FullGiftCouponCalculator"
}

func (c *FullGiftCouponCalculator) IsApplicable(orderType int) bool {
    return orderType == consts.TradeOrderTypeNormal
}

// 2. 在初始化时注册
calculators := []PriceCalculator{
    // ... 其他计算器
    NewFullGiftCouponCalculator(...),
}

service := NewTradePriceService(calculators, ...)
```

**无需修改现有代码！**

---

## 总结

### 核心要点

1. **策略模式** - 使计算器可以灵活组合和扩展
2. **工厂模式** - 集中管理计算器的创建和注册
3. **模板方法模式** - 提供通用的计算器基类
4. **助手模式** - 避免代码重复，提供通用工具方法

### 为什么这样设计？

| 问题 | 解决方案 | 好处 |
|------|--------|------|
| 计算逻辑复杂 | 分离成多个计算器 | 易于理解和维护 |
| 难以扩展 | 使用策略模式 | 添加新折扣无需修改现有代码 |
| 代码重复 | 使用助手类 | 提高代码复用率 |
| 计算顺序重要 | 使用优先级机制 | 确保计算的正确性 |
| 难以测试 | 每个计算器独立 | 可以单独测试每个计算器 |

### 学习建议

1. **先理解整体流程** - 从 `CalculateOrderPrice` 开始
2. **再学习每个计算器** - 理解各个计算器的职责
3. **最后学习设计模式** - 理解为什么这样设计
4. **动手实践** - 尝试添加新的计算器

---

**文档完成时间：** 2024年12月26日
**适用版本：** Go 版本的订单金额计算系统
