# 芋道商城 Go 版本 - 模块深度解析

> 本文档深入解析项目的各个核心业务模块，包括数据模型、业务流程、关键算法和扩展点。

## 目录

- [系统管理模块](#系统管理模块)
- [会员中心模块](#会员中心模块)
- [商品中心模块](#商品中心模块)
- [交易中心模块](#交易中心模块)
- [支付中心模块](#支付中心模块)
- [促销中心模块](#促销中心模块)
- [分销模块](#分销模块)

---

## 系统管理模块

### 模块概述

系统管理模块提供企业级后台管理系统的基础功能，包括用户、角色、权限、菜单、部门等。

### 核心数据模型

```
SystemUser (用户表)
├── ID: 用户ID
├── Username: 用户名
├── Password: 密码（加密存储）
├── Nickname: 昵称
├── DeptID: 部门ID
├── Email: 邮箱
├── Mobile: 手机号
├── Status: 状态 (0=启用, 1=禁用)
├── TenantID: 租户ID
└── CreatedAt/UpdatedAt: 时间戳

SystemRole (角色表)
├── ID: 角色ID
├── Name: 角色名称
├── Code: 角色编码
├── Status: 状态
└── TenantID: 租户ID

SystemRoleMenu (角色菜单关联表)
├── RoleID: 角色ID
├── MenuID: 菜单ID
└── Permissions: 权限标识

SystemUserRole (用户角色关联表)
├── UserID: 用户ID
└── RoleID: 角色ID

SystemMenu (菜单表)
├── ID: 菜单ID
├── ParentID: 父菜单ID
├── Name: 菜单名称
├── Path: 路由路径
├── Component: 组件名称
├── Permissions: 权限标识
└── Status: 状态
```

### 业务流程

#### 1. 用户登录流程

```
用户输入用户名和密码
    ↓
POST /admin-api/system/auth/login
    ↓
AuthHandler.Login()
    ├─ 参数验证
    └─ 调用 AuthService.Login()
    ↓
AuthService.Login()
    ├─ 1. 查询用户
    │   └─ SELECT * FROM system_user WHERE username = ?
    │
    ├─ 2. 验证密码
    │   └─ utils.CheckPassword(inputPwd, dbPwd)
    │
    ├─ 3. 检查用户状态
    │   └─ IF user.Status != 0 THEN 用户已禁用
    │
    ├─ 4. 查询用户角色
    │   └─ SELECT role_id FROM system_user_role WHERE user_id = ?
    │
    ├─ 5. 查询角色权限
    │   └─ SELECT menu_id FROM system_role_menu WHERE role_id IN (...)
    │
    ├─ 6. 生成 JWT Token
    │   └─ utils.GenerateTokenWithInfo(userID, userType, tenantID, nickname)
    │
    ├─ 7. 存储 Token 到 Redis（白名单）
    │   └─ SET oauth2_access_token:{token} {userInfo} EX 86400
    │
    ├─ 8. 记录登录日志
    │   └─ INSERT INTO system_login_log (...)
    │
    └─ 9. 返回 Token 和用户信息
    ↓
Handler 返回响应
{
    "code": 0,
    "msg": "success",
    "data": {
        "token": "eyJhbGc...",
        "user": {
            "id": 1,
            "username": "admin",
            "nickname": "管理员",
            "roles": ["admin"],
            "permissions": ["system:user:query", "system:user:create"]
        }
    }
}
```

#### 2. 权限检查流程

```
用户请求受保护的资源
    ↓
Auth 中间件验证 Token
    ├─ 验证 JWT 签名
    ├─ 检查 Token 有效期
    └─ 检查 Redis 白名单
    ↓
提取用户信息到上下文
    ↓
Handler 执行业务逻辑
    ├─ 获取登录用户信息
    │   └─ loginUser := core.GetLoginUser(c)
    │
    ├─ 检查权限
    │   └─ permissionService.CheckPermission(userID, "system:user:delete")
    │
    └─ 执行业务逻辑或返回 403
    ↓
返回结果
```

#### 3. 权限检查实现

```go
// PermissionService 权限检查
func (s *PermissionService) CheckPermission(ctx context.Context, userID int64, permission string) (bool, error) {
    // 1. 获取用户的所有角色
    userRoles, err := s.getUserRoles(ctx, userID)
    if err != nil {
        return false, err
    }
    
    // 2. 遍历角色，检查权限
    for _, roleID := range userRoles {
        permissions, err := s.getRolePermissions(ctx, roleID)
        if err != nil {
            continue
        }
        
        // 3. 检查权限是否存在
        for _, perm := range permissions {
            if perm == permission || perm == "*:*:*" {
                return true, nil
            }
        }
    }
    
    return false, nil
}

// 在 Handler 中使用
func (h *UserHandler) DeleteUser(c *gin.Context) {
    loginUser := core.GetLoginUser(c)
    
    // 检查权限
    hasPermission, err := h.permissionService.CheckPermission(
        c.Request.Context(),
        loginUser.UserID,
        "system:user:delete",
    )
    
    if !hasPermission {
        core.WriteError(c, core.ForbiddenCode, "无权限执行此操作")
        return
    }
    
    // 执行删除逻辑
    // ...
}
```

### 扩展点

1. **自定义权限检查** - 在 PermissionService 中添加更复杂的权限逻辑
2. **审计日志** - 在关键操作中记录操作日志
3. **数据权限** - 实现基于部门、岗位的数据权限控制
4. **菜单动态生成** - 根据用户权限动态生成菜单

---

## 会员中心模块

### 模块概述

会员中心模块管理平台的普通用户（会员），包括用户信息、等级、积分、签到等功能。

### 核心数据模型

```
MemberUser (会员用户表)
├── ID: 用户ID
├── Username: 用户名
├── Password: 密码
├── Nickname: 昵称
├── Avatar: 头像
├── Mobile: 手机号
├── Email: 邮箱
├── LevelID: 会员等级ID
├── Points: 积分
├── Balance: 余额
├── Status: 状态
└── CreatedAt/UpdatedAt: 时间戳

MemberLevel (会员等级表)
├── ID: 等级ID
├── Name: 等级名称
├── Icon: 等级图标
├── RequiredPoints: 升级所需积分
├── Discount: 折扣率
└── Benefits: 等级权益

MemberPointRecord (积分记录表)
├── ID: 记录ID
├── UserID: 用户ID
├── Points: 积分数量（正数增加，负数扣除）
├── Type: 类型 (1=购物获得, 2=签到获得, 3=兑换消耗)
├── Reason: 原因
└── CreatedAt: 创建时间

MemberSignInRecord (签到记录表)
├── ID: 记录ID
├── UserID: 用户ID
├── SignInDate: 签到日期
├── ContinuousDays: 连续签到天数
├── Points: 获得积分
└── CreatedAt: 创建时间
```

### 业务流程

#### 1. 会员注册流程

```
用户输入注册信息
    ↓
POST /app-api/member/auth/register
    ↓
AppAuthHandler.Register()
    ├─ 参数验证
    └─ 调用 AppAuthService.Register()
    ↓
AppAuthService.Register()
    ├─ 1. 检查用户名唯一性
    │   └─ SELECT COUNT(*) FROM member_user WHERE username = ?
    │
    ├─ 2. 检查手机号唯一性
    │   └─ SELECT COUNT(*) FROM member_user WHERE mobile = ?
    │
    ├─ 3. 加密密码
    │   └─ utils.HashPassword(password)
    │
    ├─ 4. 创建用户
    │   └─ INSERT INTO member_user (...)
    │
    ├─ 5. 初始化用户积分（可选）
    │   └─ INSERT INTO member_point_record (...)
    │
    ├─ 6. 生成 JWT Token
    │   └─ utils.GenerateTokenWithInfo(userID, 0, tenantID, nickname)
    │
    └─ 7. 返回 Token
    ↓
返回响应
```

#### 2. 积分系统流程

```
用户完成积分获取行为（如购物、签到）
    ↓
Service 调用 MemberPointRecordService.AddPoints()
    ├─ 1. 验证积分数量
    │   └─ IF points <= 0 THEN 返回错误
    │
    ├─ 2. 创建积分记录
    │   └─ INSERT INTO member_point_record (...)
    │
    ├─ 3. 更新用户积分
    │   └─ UPDATE member_user SET points = points + ? WHERE id = ?
    │
    ├─ 4. 检查等级升级
    │   └─ IF user.points >= level.required_points THEN 升级
    │
    └─ 5. 返回结果
    ↓
业务逻辑继续
```

#### 3. 签到系统流程

```
用户点击签到按钮
    ↓
POST /app-api/member/sign-in/sign-in
    ↓
AppMemberSignInRecordHandler.SignIn()
    ├─ 参数验证
    └─ 调用 MemberSignInRecordService.SignIn()
    ↓
MemberSignInRecordService.SignIn()
    ├─ 1. 检查今天是否已签到
    │   └─ SELECT * FROM member_sign_in_record 
    │       WHERE user_id = ? AND sign_in_date = TODAY()
    │
    ├─ 2. 如果已签到，返回错误
    │   └─ RETURN 今天已签到
    │
    ├─ 3. 获取签到配置
    │   └─ SELECT * FROM member_sign_in_config WHERE id = 1
    │
    ├─ 4. 计算连续签到天数
    │   └─ SELECT MAX(continuous_days) FROM member_sign_in_record 
    │       WHERE user_id = ? AND sign_in_date >= DATE_SUB(TODAY(), INTERVAL 1 DAY)
    │
    ├─ 5. 确定签到奖励
    │   └─ IF 连续签到 THEN 奖励 = 基础奖励 * 倍数
    │
    ├─ 6. 创建签到记录
    │   └─ INSERT INTO member_sign_in_record (...)
    │
    ├─ 7. 增加用户积分
    │   └─ memberPointRecordService.AddPoints(userID, points, "签到获得")
    │
    └─ 8. 返回签到结果
    ↓
返回响应
{
    "code": 0,
    "msg": "success",
    "data": {
        "points": 10,
        "continuousDays": 5,
        "totalPoints": 150
    }
}
```

### 扩展点

1. **等级权益系统** - 不同等级享受不同的折扣和权益
2. **积分兑换商城** - 用户可以用积分兑换商品或优惠券
3. **会员分组** - 根据消费行为或标签对会员进行分组
4. **推荐系统** - 基于会员行为的个性化推荐

---

## 商品中心模块

### 模块概述

商品中心模块管理平台的所有商品信息，包括分类、品牌、属性、SPU/SKU 等。

### 核心数据模型

```
ProductCategory (商品分类表)
├── ID: 分类ID
├── ParentID: 父分类ID
├── Name: 分类名称
├── Icon: 分类图标
├── Sort: 排序
└── Status: 状态

ProductBrand (商品品牌表)
├── ID: 品牌ID
├── Name: 品牌名称
├── Logo: 品牌logo
├── Description: 品牌描述
└── Status: 状态

ProductProperty (商品属性表)
├── ID: 属性ID
├── CategoryID: 分类ID
├── Name: 属性名称
├── Type: 属性类型 (1=规格, 2=参数)
├── Values: 属性值列表 (JSON)
└── Status: 状态

ProductSPU (商品SPU表 - 标准产品单元)
├── ID: SPUID
├── CategoryID: 分类ID
├── BrandID: 品牌ID
├── Name: 商品名称
├── Description: 商品描述
├── MainPicture: 主图
├── Pictures: 图片列表 (JSON)
├── Price: 价格
├── Status: 状态
└── CreatedAt/UpdatedAt: 时间戳

ProductSKU (商品SKU表 - 库存单位)
├── ID: SKUID
├── SPUID: SPUID
├── SkuCode: SKU编码
├── Properties: 属性值 (JSON)
├── Price: 价格
├── Stock: 库存
├── SoldCount: 销售数量
└── Status: 状态

ProductComment (商品评论表)
├── ID: 评论ID
├── SPUID: SPUID
├── UserID: 用户ID
├── Rating: 评分 (1-5)
├── Content: 评论内容
├── Pictures: 评论图片 (JSON)
├── Status: 状态
└── CreatedAt: 创建时间

ProductFavorite (商品收藏表)
├── ID: 收藏ID
├── UserID: 用户ID
├── SPUID: SPUID
└── CreatedAt: 创建时间

ProductBrowseHistory (浏览历史表)
├── ID: 记录ID
├── UserID: 用户ID
├── SPUID: SPUID
└── CreatedAt: 浏览时间
```

### 业务流程

#### 1. 商品展示流程

```
用户浏览商品列表
    ↓
GET /app-api/product/spu/list
    ↓
AppProductSpuHandler.List()
    ├─ 参数验证（分类ID、排序、分页）
    └─ 调用 ProductSpuService.GetSpuList()
    ↓
ProductSpuService.GetSpuList()
    ├─ 1. 构建查询条件
    │   ├─ 分类ID 过滤
    │   ├─ 品牌ID 过滤
    │   ├─ 价格范围 过滤
    │   └─ 关键词搜索
    │
    ├─ 2. 查询 SPU 列表
    │   └─ SELECT * FROM product_spu 
    │       WHERE category_id = ? AND status = 1
    │       ORDER BY sort DESC, created_at DESC
    │       LIMIT ? OFFSET ?
    │
    ├─ 3. 查询 SKU 信息（库存、价格）
    │   └─ SELECT * FROM product_sku WHERE spu_id IN (...)
    │
    ├─ 4. 查询商品评论统计
    │   └─ SELECT COUNT(*), AVG(rating) FROM product_comment 
    │       WHERE spu_id IN (...)
    │
    ├─ 5. 查询用户收藏状态（如果已登录）
    │   └─ SELECT spu_id FROM product_favorite 
    │       WHERE user_id = ? AND spu_id IN (...)
    │
    └─ 6. 组装响应数据
    ↓
返回商品列表
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
                "rating": 4.5,
                "commentCount": 100,
                "isFavorite": true
            }
        ],
        "total": 1000
    }
}
```

#### 2. 商品详情流程

```
用户点击商品
    ↓
GET /app-api/product/spu/detail/:id
    ↓
AppProductSpuHandler.GetDetail()
    ├─ 参数验证
    └─ 调用 ProductSpuService.GetSpuDetail()
    ↓
ProductSpuService.GetSpuDetail()
    ├─ 1. 查询 SPU 信息
    │   └─ SELECT * FROM product_spu WHERE id = ?
    │
    ├─ 2. 查询 SKU 列表
    │   └─ SELECT * FROM product_sku WHERE spu_id = ?
    │
    ├─ 3. 查询商品属性
    │   └─ SELECT * FROM product_property 
    │       WHERE category_id = (SELECT category_id FROM product_spu WHERE id = ?)
    │
    ├─ 4. 查询商品评论
    │   └─ SELECT * FROM product_comment WHERE spu_id = ?
    │       ORDER BY created_at DESC LIMIT 10
    │
    ├─ 5. 查询用户收藏状态
    │   └─ SELECT * FROM product_favorite 
    │       WHERE user_id = ? AND spu_id = ?
    │
    ├─ 6. 记录浏览历史
    │   └─ INSERT INTO product_browse_history (user_id, spu_id, created_at)
    │       ON DUPLICATE KEY UPDATE created_at = NOW()
    │
    └─ 7. 组装详情数据
    ↓
返回商品详情
```

#### 3. 商品评论流程

```
用户提交评论
    ↓
POST /app-api/product/comment/create
    ↓
AppProductCommentHandler.Create()
    ├─ 参数验证
    └─ 调用 ProductCommentService.CreateComment()
    ↓
ProductCommentService.CreateComment()
    ├─ 1. 检查用户是否购买过该商品
    │   └─ SELECT COUNT(*) FROM trade_order_item 
    │       WHERE user_id = ? AND spu_id = ? AND order_status = 已完成
    │
    ├─ 2. 检查是否已评论
    │   └─ SELECT * FROM product_comment 
    │       WHERE user_id = ? AND spu_id = ?
    │
    ├─ 3. 创建评论
    │   └─ INSERT INTO product_comment (...)
    │
    ├─ 4. 更新商品评分
    │   └─ UPDATE product_spu SET rating = (
    │       SELECT AVG(rating) FROM product_comment WHERE spu_id = ?
    │       )
    │
    └─ 5. 返回评论ID
    ↓
返回成功响应
```

### 扩展点

1. **商品搜索** - 集成 Elasticsearch 实现全文搜索
2. **推荐系统** - 基于用户行为的个性化推荐
3. **库存管理** - 实现库存预警和自动补货
4. **商品评分** - 更复杂的评分算法（考虑时间、有用性等）

---

## 交易中心模块

### 模块概述

交易中心模块管理订单、购物车、售后等交易流程。

### 核心数据模型

```
TradeCart (购物车表)
├── ID: 购物车ID
├── UserID: 用户ID
├── SKUID: SKUID
├── Quantity: 数量
├── Selected: 是否选中
└── CreatedAt/UpdatedAt: 时间戳

TradeOrder (订单表)
├── ID: 订单ID
├── OrderNo: 订单号
├── UserID: 用户ID
├── TotalAmount: 订单总额
├── PayAmount: 实付金额
├── Status: 订单状态 (1=待支付, 2=已支付, 3=待发货, 4=已发货, 5=已完成, 6=已取消)
├── PaymentTime: 支付时间
├── DeliveryTime: 发货时间
├── ReceiveTime: 收货时间
├── CancelTime: 取消时间
├── CancelReason: 取消原因
└── CreatedAt/UpdatedAt: 时间戳

TradeOrderItem (订单项表)
├── ID: 项ID
├── OrderID: 订单ID
├── SPUID: SPUID
├── SKUID: SKUID
├── Quantity: 数量
├── Price: 单价
├── Amount: 小计
└── CreatedAt: 创建时间

TradeAfterSale (售后表)
├── ID: 售后ID
├── OrderID: 订单ID
├── OrderItemID: 订单项ID
├── Type: 售后类型 (1=退货, 2=退款, 3=换货)
├── Reason: 原因
├── Status: 状态 (1=待审核, 2=已同意, 3=待退货, 4=已收货, 5=已完成, 6=已拒绝)
├── RefundAmount: 退款金额
└── CreatedAt/UpdatedAt: 时间戳

DeliveryExpress (快递公司表)
├── ID: 快递ID
├── Name: 快递名称
├── Code: 快递编码
└── Status: 状态

DeliveryExpressTemplate (运费模板表)
├── ID: 模板ID
├── Name: 模板名称
├── ChargeType: 计费方式 (1=按重量, 2=按件数)
├── Rules: 运费规则 (JSON)
└── Status: 状态
```

### 业务流程

#### 1. 下单流程

```
用户点击结算
    ↓
POST /app-api/trade/order/create
    ↓
AppTradeOrderHandler.Create()
    ├─ 参数验证（收货地址、优惠券等）
    └─ 调用 TradeOrderService.CreateOrder()
    ↓
TradeOrderService.CreateOrder()
    ├─ 1. 开启事务
    │   └─ tx := db.Begin()
    │
    ├─ 2. 查询购物车商品
    │   └─ SELECT * FROM trade_cart WHERE user_id = ? AND selected = 1
    │
    ├─ 3. 检查库存
    │   ├─ FOR EACH item IN cart
    │   │   └─ SELECT stock FROM product_sku WHERE id = ? FOR UPDATE
    │   │       IF stock < quantity THEN 库存不足
    │   └─ END FOR
    │
    ├─ 4. 计算订单金额
    │   ├─ 商品总额 = SUM(sku.price * quantity)
    │   ├─ 运费 = calculateShipping(items, address)
    │   ├─ 优惠券折扣 = calculateCouponDiscount(coupon)
    │   └─ 实付金额 = 商品总额 + 运费 - 优惠券折扣
    │
    ├─ 5. 创建订单
    │   └─ INSERT INTO trade_order (...)
    │
    ├─ 6. 创建订单项
    │   ├─ FOR EACH item IN cart
    │   │   └─ INSERT INTO trade_order_item (...)
    │   └─ END FOR
    │
    ├─ 7. 扣减库存
    │   ├─ FOR EACH item IN cart
    │   │   └─ UPDATE product_sku SET stock = stock - ? WHERE id = ?
    │   └─ END FOR
    │
    ├─ 8. 使用优惠券
    │   └─ UPDATE promotion_coupon SET used_count = used_count + 1 WHERE id = ?
    │
    ├─ 9. 清空购物车
    │   └─ DELETE FROM trade_cart WHERE user_id = ? AND selected = 1
    │
    ├─ 10. 提交事务
    │   └─ tx.Commit()
    │
    └─ 11. 返回订单ID
    ↓
返回响应
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

#### 2. 订单支付流程

```
用户点击支付
    ↓
POST /app-api/trade/order/pay
    ↓
AppTradeOrderHandler.Pay()
    ├─ 参数验证（订单ID、支付方式）
    └─ 调用 TradeOrderService.PayOrder()
    ↓
TradeOrderService.PayOrder()
    ├─ 1. 查询订单
    │   └─ SELECT * FROM trade_order WHERE id = ? AND user_id = ?
    │
    ├─ 2. 检查订单状态
    │   └─ IF order.status != 待支付 THEN 订单状态错误
    │
    ├─ 3. 创建支付订单
    │   └─ payOrderService.CreatePayOrder(orderId, payAmount, paymentMethod)
    │
    ├─ 4. 调用支付渠道
    │   ├─ IF paymentMethod == 支付宝
    │   │   └─ 调用支付宝 API
    │   ├─ ELSE IF paymentMethod == 微信
    │   │   └─ 调用微信 API
    │   └─ END IF
    │
    └─ 5. 返回支付信息
    ↓
返回支付链接或二维码
```

#### 3. 订单发货流程

```
商家点击发货
    ↓
POST /admin-api/trade/order/delivery
    ↓
TradeOrderHandler.Delivery()
    ├─ 参数验证（订单ID、快递公司、快递单号）
    └─ 调用 TradeOrderService.DeliveryOrder()
    ↓
TradeOrderService.DeliveryOrder()
    ├─ 1. 查询订单
    │   └─ SELECT * FROM trade_order WHERE id = ?
    │
    ├─ 2. 检查订单状态
    │   └─ IF order.status != 已支付 THEN 订单状态错误
    │
    ├─ 3. 更新订单状态
    │   └─ UPDATE trade_order SET status = 已发货, delivery_time = NOW()
    │
    ├─ 4. 保存快递信息
    │   └─ INSERT INTO trade_order_delivery (...)
    │
    └─ 5. 返回成功
    ↓
返回成功响应
```

#### 4. 售后流程

```
用户申请售后
    ↓
POST /app-api/trade/after-sale/apply
    ↓
AppTradeAfterSaleHandler.Apply()
    ├─ 参数验证
    └─ 调用 TradeAfterSaleService.ApplyAfterSale()
    ↓
TradeAfterSaleService.ApplyAfterSale()
    ├─ 1. 检查订单和订单项
    │   └─ SELECT * FROM trade_order_item WHERE id = ? AND order_id = ?
    │
    ├─ 2. 检查是否已申请售后
    │   └─ SELECT * FROM trade_after_sale WHERE order_item_id = ?
    │
    ├─ 3. 创建售后单
    │   └─ INSERT INTO trade_after_sale (...)
    │
    └─ 4. 返回售后单ID
    ↓
商家审核售后
    ├─ 1. 查询售后单
    │   └─ SELECT * FROM trade_after_sale WHERE id = ?
    │
    ├─ 2. 审核通过/拒绝
    │   └─ UPDATE trade_after_sale SET status = 已同意/已拒绝
    │
    └─ 3. 如果通过，等待用户退货
    ↓
用户退货
    ├─ 1. 快递上门取件
    │   └─ 用户填写快递单号
    │
    └─ 2. 更新售后单状态
        └─ UPDATE trade_after_sale SET status = 待收货
    ↓
商家收货确认
    ├─ 1. 检查商品
    │   └─ 确认商品完好
    │
    ├─ 2. 更新售后单状态
    │   └─ UPDATE trade_after_sale SET status = 已完成
    │
    ├─ 3. 处理退款
    │   └─ payRefundService.CreateRefund(...)
    │
    └─ 4. 恢复库存
        └─ UPDATE product_sku SET stock = stock + quantity
```

### 扩展点

1. **订单推荐** - 基于订单历史的推荐
2. **订单预测** - 预测用户可能购买的商品
3. **物流跟踪** - 实时物流信息推送
4. **订单分析** - 订单数据分析和报表

---

## 支付中心模块

### 模块概述

支付中心模块集成多种支付渠道，管理支付订单、退款等。

### 核心数据模型

```
PayApp (支付应用表)
├── ID: 应用ID
├── Name: 应用名称
├── AppID: 应用ID
├── AppSecret: 应用密钥
└── Status: 状态

PayChannel (支付渠道表)
├── ID: 渠道ID
├── Code: 渠道编码 (alipay, wechat, balance)
├── Name: 渠道名称
├── AppID: 应用ID
├── Config: 渠道配置 (JSON)
└── Status: 状态

PayOrder (支付订单表)
├── ID: 支付订单ID
├── OrderNo: 订单号
├── TradeOrderID: 交易订单ID
├── Amount: 支付金额
├── ChannelID: 支付渠道ID
├── Status: 状态 (1=待支付, 2=已支付, 3=支付失败, 4=已关闭)
├── PaymentTime: 支付时间
├── ChannelOrderNo: 渠道订单号
└── CreatedAt/UpdatedAt: 时间戳

PayRefund (退款表)
├── ID: 退款ID
├── PayOrderID: 支付订单ID
├── RefundNo: 退款号
├── Amount: 退款金额
├── Reason: 退款原因
├── Status: 状态 (1=待退款, 2=已退款, 3=退款失败)
├── RefundTime: 退款时间
└── CreatedAt/UpdatedAt: 时间戳

PayNotify (支付回调表)
├── ID: 回调ID
├── OrderNo: 订单号
├── ChannelID: 渠道ID
├── Content: 回调内容 (JSON)
├── Status: 处理状态
└── CreatedAt: 创建时间
```

### 业务流程

#### 1. 支付流程

```
用户点击支付
    ↓
POST /app-api/trade/order/pay
    ↓
TradeOrderHandler.Pay()
    ├─ 参数验证
    └─ 调用 PayOrderService.CreatePayOrder()
    ↓
PayOrderService.CreatePayOrder()
    ├─ 1. 创建支付订单
    │   └─ INSERT INTO pay_order (...)
    │
    ├─ 2. 查询支付渠道配置
    │   └─ SELECT * FROM pay_channel WHERE id = ?
    │
    ├─ 3. 调用渠道支付接口
    │   ├─ IF channel == 支付宝
    │   │   └─ alipayClient.Pay(...)
    │   ├─ ELSE IF channel == 微信
    │   │   └─ wechatClient.Pay(...)
    │   └─ END IF
    │
    └─ 4. 返回支付信息
    ↓
返回支付链接或二维码
{
    "code": 0,
    "msg": "success",
    "data": {
        "payUrl": "https://...",
        "payOrderNo": "PAY20241215123456"
    }
}
```

#### 2. 支付回调流程

```
支付渠道异步通知支付结果
    ↓
POST /app-api/pay/notify/{channel}
    ↓
PayNotifyHandler.Notify()
    ├─ 参数验证
    └─ 调用 PayNotifyService.HandleNotify()
    ↓
PayNotifyService.HandleNotify()
    ├─ 1. 验证签名
    │   └─ IF 签名验证失败 THEN 返回失败
    │
    ├─ 2. 查询支付订单
    │   └─ SELECT * FROM pay_order WHERE channel_order_no = ?
    │
    ├─ 3. 检查订单状态
    │   └─ IF order.status != 待支付 THEN 已处理，返回成功
    │
    ├─ 4. 更新支付订单状态
    │   └─ UPDATE pay_order SET status = 已支付, payment_time = NOW()
    │
    ├─ 5. 更新交易订单状态
    │   └─ UPDATE trade_order SET status = 已支付, payment_time = NOW()
    │
    ├─ 6. 触发订单支付成功事件
    │   └─ eventBus.Publish(OrderPaidEvent)
    │
    └─ 7. 返回成功
    ↓
返回 200 OK（告诉渠道已处理）
```

#### 3. 退款流程

```
商家或用户申请退款
    ↓
POST /admin-api/pay/refund/create
    ↓
PayRefundHandler.Create()
    ├─ 参数验证
    └─ 调用 PayRefundService.CreateRefund()
    ↓
PayRefundService.CreateRefund()
    ├─ 1. 查询支付订单
    │   └─ SELECT * FROM pay_order WHERE id = ?
    │
    ├─ 2. 检查支付订单状态
    │   └─ IF order.status != 已支付 THEN 无法退款
    │
    ├─ 3. 创建退款单
    │   └─ INSERT INTO pay_refund (...)
    │
    ├─ 4. 调用渠道退款接口
    │   ├─ IF channel == 支付宝
    │   │   └─ alipayClient.Refund(...)
    │   ├─ ELSE IF channel == 微信
    │   │   └─ wechatClient.Refund(...)
    │   └─ END IF
    │
    ├─ 5. 更新退款单状态
    │   └─ UPDATE pay_refund SET status = 已退款, refund_time = NOW()
    │
    └─ 6. 返回退款ID
    ↓
返回成功响应
```

### 扩展点

1. **支付渠道扩展** - 添加新的支付渠道（如银行卡、数字钱包）
2. **支付风控** - 实现支付风险控制和反欺诈
3. **对账系统** - 与支付渠道对账
4. **支付分析** - 支付数据分析和报表

---

## 促销中心模块

### 模块概述

促销中心模块管理各种营销活动，包括优惠券、秒杀、拼团、砍价等。

### 核心数据模型

```
PromotionCoupon (优惠券表)
├── ID: 优惠券ID
├── Name: 优惠券名称
├── Type: 类型 (1=满减, 2=折扣, 3=代金券)
├── DiscountType: 折扣类型 (1=固定金额, 2=百分比)
├── DiscountValue: 折扣值
├── MinAmount: 最小消费金额
├── MaxAmount: 最大优惠金额
├── TotalCount: 总数量
├── UsedCount: 已使用数量
├── StartTime: 开始时间
├── EndTime: 结束时间
└── Status: 状态

PromotionSeckillActivity (秒杀活动表)
├── ID: 活动ID
├── Name: 活动名称
├── SPUID: 商品ID
├── OriginalPrice: 原价
├── SeckillPrice: 秒杀价
├── SeckillStock: 秒杀库存
├── SeckillSoldCount: 已秒杀数量
├── StartTime: 开始时间
├── EndTime: 结束时间
└── Status: 状态

PromotionCombinationActivity (拼团活动表)
├── ID: 活动ID
├── Name: 活动名称
├── SPUID: 商品ID
├── OriginalPrice: 原价
├── CombinationPrice: 拼团价
├── RequiredCount: 成团人数
├── LimitCount: 每人限购数量
├── StartTime: 开始时间
├── EndTime: 结束时间
└── Status: 状态

PromotionCombinationRecord (拼团记录表)
├── ID: 记录ID
├── ActivityID: 活动ID
├── GroupID: 团ID
├── UserID: 用户ID
├── Quantity: 数量
├── Status: 状态 (1=待成团, 2=已成团, 3=已失败)
└── CreatedAt: 创建时间

PromotionBargainActivity (砍价活动表)
├── ID: 活动ID
├── Name: 活动名称
├── SPUID: 商品ID
├── OriginalPrice: 原价
├── MinPrice: 最低价
├── StartTime: 开始时间
├── EndTime: 结束时间
└── Status: 状态

PromotionBargainRecord (砍价记录表)
├── ID: 记录ID
├── ActivityID: 活动ID
├── UserID: 用户ID
├── CurrentPrice: 当前价格
├── Status: 状态 (1=砍价中, 2=已完成)
└── CreatedAt: 创建时间

PromotionBargainHelp (砍价助力表)
├── ID: 助力ID
├── BargainRecordID: 砍价记录ID
├── HelpUserID: 助力用户ID
├── ReducePrice: 砍价金额
└── CreatedAt: 创建时间
```

### 业务流程

#### 1. 优惠券使用流程

```
用户在下单时选择优惠券
    ↓
POST /app-api/trade/order/create
    ↓
TradeOrderService.CreateOrder()
    ├─ 1. 查询优惠券
    │   └─ SELECT * FROM promotion_coupon WHERE id = ? AND user_id = ?
    │
    ├─ 2. 检查优惠券状态
    │   ├─ IF coupon.status != 可用 THEN 优惠券已失效
    │   ├─ IF coupon.start_time > NOW() THEN 优惠券未开始
    │   ├─ IF coupon.end_time < NOW() THEN 优惠券已过期
    │   └─ IF coupon.used_count >= coupon.total_count THEN 优惠券已用完
    │
    ├─ 3. 检查优惠券使用条件
    │   └─ IF order_amount < coupon.min_amount THEN 不满足最小消费
    │
    ├─ 4. 计算优惠金额
    │   ├─ IF coupon.discount_type == 固定金额
    │   │   └─ discount = coupon.discount_value
    │   ├─ ELSE IF coupon.discount_type == 百分比
    │   │   └─ discount = order_amount * coupon.discount_value / 100
    │   └─ END IF
    │   └─ discount = MIN(discount, coupon.max_amount)
    │
    ├─ 5. 更新优惠券使用状态
    │   └─ UPDATE promotion_coupon SET used_count = used_count + 1
    │
    └─ 6. 在订单中记录优惠券
        └─ INSERT INTO trade_order (coupon_id, coupon_discount, ...)
    ↓
继续下单流程
```

#### 2. 秒杀流程

```
用户进入秒杀页面
    ↓
GET /app-api/promotion/seckill/activity/:id
    ↓
AppSeckillActivityHandler.GetDetail()
    ├─ 1. 查询秒杀活动
    │   └─ SELECT * FROM promotion_seckill_activity WHERE id = ?
    │
    ├─ 2. 检查活动状态
    │   ├─ IF activity.status != 进行中 THEN 活动未开始或已结束
    │   └─ IF activity.start_time > NOW() THEN 活动未开始
    │
    ├─ 3. 计算剩余库存
    │   └─ remaining = activity.seckill_stock - activity.seckill_sold_count
    │
    └─ 4. 返回活动信息
    ↓
用户点击秒杀按钮
    ↓
POST /app-api/promotion/seckill/buy
    ↓
AppSeckillActivityHandler.Buy()
    ├─ 参数验证
    └─ 调用 SeckillActivityService.BuySeckill()
    ↓
SeckillActivityService.BuySeckill()
    ├─ 1. 使用分布式锁防止超卖
    │   └─ lock := redis.SetNX("seckill:activity:{id}", "1", 1s)
    │       IF !lock THEN 秒杀已结束
    │
    ├─ 2. 检查库存
    │   └─ SELECT seckill_stock - seckill_sold_count AS remaining
    │       IF remaining <= 0 THEN 库存不足
    │
    ├─ 3. 扣减库存
    │   └─ UPDATE promotion_seckill_activity 
    │       SET seckill_sold_count = seckill_sold_count + 1
    │       WHERE id = ? AND seckill_sold_count < seckill_stock
    │
    ├─ 4. 创建订单
    │   └─ 调用 TradeOrderService.CreateOrder(...)
    │
    └─ 5. 返回订单ID
    ↓
返回成功响应
```

#### 3. 拼团流程

```
用户发起拼团
    ↓
POST /app-api/promotion/combination/create
    ↓
AppCombinationActivityHandler.Create()
    ├─ 参数验证
    └─ 调用 CombinationActivityService.CreateGroup()
    ↓
CombinationActivityService.CreateGroup()
    ├─ 1. 查询拼团活动
    │   └─ SELECT * FROM promotion_combination_activity WHERE id = ?
    │
    ├─ 2. 检查活动状态
    │   └─ IF activity.status != 进行中 THEN 活动已结束
    │
    ├─ 3. 生成团ID
    │   └─ groupID = generateGroupID()
    │
    ├─ 4. 创建拼团记录
    │   └─ INSERT INTO promotion_combination_record (...)
    │
    └─ 5. 返回团ID
    ↓
其他用户加入拼团
    ↓
POST /app-api/promotion/combination/join
    ↓
AppCombinationActivityHandler.Join()
    ├─ 参数验证
    └─ 调用 CombinationActivityService.JoinGroup()
    ↓
CombinationActivityService.JoinGroup()
    ├─ 1. 查询拼团记录
    │   └─ SELECT * FROM promotion_combination_record WHERE group_id = ?
    │
    ├─ 2. 检查拼团状态
    │   └─ IF record.status != 待成团 THEN 拼团已成团或已失败
    │
    ├─ 3. 检查人数
    │   └─ SELECT COUNT(*) FROM promotion_combination_record 
    │       WHERE group_id = ?
    │       IF count >= activity.required_count THEN 拼团已满
    │
    ├─ 4. 添加拼团成员
    │   └─ INSERT INTO promotion_combination_record (...)
    │
    ├─ 5. 检查是否成团
    │   ├─ SELECT COUNT(*) FROM promotion_combination_record 
    │   │   WHERE group_id = ?
    │   └─ IF count >= activity.required_count THEN 成团
    │       └─ UPDATE promotion_combination_record SET status = 已成团
    │
    └─ 6. 返回拼团信息
    ↓
返回成功响应
```

### 扩展点

1. **优惠券推荐** - 基于用户行为推荐优惠券
2. **活动分析** - 活动效果分析和优化
3. **限时秒杀** - 更复杂的秒杀规则和库存管理
4. **社交分享** - 拼团和砍价的社交分享功能

---

## 分销模块

### 模块概述

分销模块实现社交分销体系，支持分销商管理、佣金计算、提现等功能。

### 核心数据模型

```
BrokerageUser (分销商表)
├── ID: 分销商ID
├── UserID: 用户ID
├── Level: 分销商等级 (1=一级, 2=二级)
├── ParentID: 上级分销商ID
├── TotalBrokerage: 总佣金
├── AvailableBrokerage: 可用佣金
├── WithdrawnBrokerage: 已提现佣金
├── Status: 状态
└── CreatedAt/UpdatedAt: 时间戳

BrokerageRecord (佣金记录表)
├── ID: 记录ID
├── UserID: 分销商ID
├── OrderID: 订单ID
├── BrokerageAmount: 佣金金额
├── Type: 类型 (1=销售佣金, 2=推荐佣金)
├── Status: 状态 (1=待结算, 2=已结算)
└── CreatedAt: 创建时间

BrokerageWithdraw (提现记录表)
├── ID: 提现ID
├── UserID: 分销商ID
├── Amount: 提现金额
├── BankAccount: 银行账户
├── Status: 状态 (1=待审核, 2=已审核, 3=已提现, 4=已拒绝)
├── AuditTime: 审核时间
├── WithdrawTime: 提现时间
└── CreatedAt: 创建时间
```

### 业务流程

#### 1. 分销商注册流程

```
用户通过分销链接进入
    ↓
GET /app-api/brokerage/user/register?referrer_id=123
    ↓
AppBrokerageUserHandler.Register()
    ├─ 参数验证
    └─ 调用 BrokerageUserService.Register()
    ↓
BrokerageUserService.Register()
    ├─ 1. 检查推荐人
    │   └─ SELECT * FROM brokerage_user WHERE user_id = ?
    │
    ├─ 2. 确定分销商等级
    │   ├─ IF 推荐人是一级分销商 THEN 新用户为二级
    │   └─ ELSE 新用户为一级
    │
    ├─ 3. 创建分销商记录
    │   └─ INSERT INTO brokerage_user (...)
    │
    └─ 4. 返回成功
    ↓
返回成功响应
```

#### 2. 佣金计算流程

```
用户通过分销链接下单
    ↓
POST /app-api/trade/order/create
    ↓
TradeOrderService.CreateOrder()
    ├─ 1. 检查是否来自分销链接
    │   └─ IF referrer_id 存在 THEN 记录推荐人
    │
    ├─ 2. 创建订单
    │   └─ INSERT INTO trade_order (referrer_id, ...)
    │
    └─ 3. 继续下单流程
    ↓
订单支付成功
    ↓
PayNotifyService.HandleNotify()
    ├─ 1. 更新订单状态
    │   └─ UPDATE trade_order SET status = 已支付
    │
    ├─ 2. 计算佣金
    │   ├─ 查询订单
    │   │   └─ SELECT * FROM trade_order WHERE id = ?
    │   │
    │   ├─ 查询推荐人
    │   │   └─ SELECT * FROM brokerage_user WHERE user_id = ?
    │   │
    │   ├─ 计算佣金金额
    │   │   ├─ 一级佣金 = order_amount * 10%
    │   │   ├─ 二级佣金 = order_amount * 5%
    │   │   └─ 佣金 = 一级佣金 + 二级佣金
    │   │
    │   └─ 创建佣金记录
    │       └─ INSERT INTO brokerage_record (...)
    │
    ├─ 3. 更新分销商可用佣金
    │   ├─ UPDATE brokerage_user 
    │   │   SET available_brokerage = available_brokerage + ?
    │   │   WHERE user_id = ?
    │   │
    │   └─ 如果有上级，也更新上级佣金
    │
    └─ 4. 返回成功
    ↓
继续支付流程
```

#### 3. 提现流程

```
分销商申请提现
    ↓
POST /app-api/brokerage/withdraw/apply
    ↓
AppBrokerageWithdrawHandler.Apply()
    ├─ 参数验证
    └─ 调用 BrokerageWithdrawService.ApplyWithdraw()
    ↓
BrokerageWithdrawService.ApplyWithdraw()
    ├─ 1. 查询分销商
    │   └─ SELECT * FROM brokerage_user WHERE user_id = ?
    │
    ├─ 2. 检查可用佣金
    │   └─ IF available_brokerage < amount THEN 余额不足
    │
    ├─ 3. 创建提现记录
    │   └─ INSERT INTO brokerage_withdraw (...)
    │
    ├─ 4. 冻结可用佣金
    │   └─ UPDATE brokerage_user 
    │       SET available_brokerage = available_brokerage - ?
    │       WHERE user_id = ?
    │
    └─ 5. 返回提现ID
    ↓
商家审核提现
    ├─ 1. 查询提现记录
    │   └─ SELECT * FROM brokerage_withdraw WHERE id = ?
    │
    ├─ 2. 审核通过/拒绝
    │   ├─ IF 审核通过
    │   │   └─ UPDATE brokerage_withdraw SET status = 已审核
    │   └─ ELSE
    │       └─ UPDATE brokerage_withdraw SET status = 已拒绝
    │           UPDATE brokerage_user 
    │           SET available_brokerage = available_brokerage + ?
    │
    └─ 3. 返回结果
    ↓
系统处理提现
    ├─ 1. 调用支付接口转账
    │   └─ paymentGateway.Transfer(bankAccount, amount)
    │
    ├─ 2. 更新提现状态
    │   └─ UPDATE brokerage_withdraw SET status = 已提现, withdraw_time = NOW()
    │
    ├─ 3. 更新分销商已提现佣金
    │   └─ UPDATE brokerage_user 
    │       SET withdrawn_brokerage = withdrawn_brokerage + ?
    │
    └─ 4. 返回成功
    ↓
返回成功响应
```

### 扩展点

1. **分销等级** - 实现更复杂的分销等级体系
2. **佣金规则** - 支持自定义佣金计算规则
3. **分销报表** - 分销数据分析和报表
4. **分销推广** - 分销推广工具和素材

---

## 总结

本文档深入解析了项目的各个核心业务模块，包括：

- **系统管理模块** - 用户、角色、权限、菜单管理
- **会员中心模块** - 会员信息、等级、积分、签到管理
- **商品中心模块** - 商品分类、品牌、属性、SPU/SKU 管理
- **交易中心模块** - 购物车、订单、售后管理
- **支付中心模块** - 支付订单、退款、回调管理
- **促销中心模块** - 优惠券、秒杀、拼团、砍价管理
- **分销模块** - 分销商、佣金、提现管理

每个模块都包含：

✅ 核心数据模型和关系
✅ 详细的业务流程和交互
✅ 关键算法和实现细节
✅ 扩展点和优化方向

通过学习本文档，你可以：

✅ 深入理解各个业务模块的设计思想
✅ 掌握复杂业务流程的实现方式
✅ 学会如何扩展和优化各个模块
✅ 为项目的功能扩展提供参考

祝你学习愉快！🚀
