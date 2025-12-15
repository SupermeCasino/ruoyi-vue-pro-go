# 芋道商城 Go 版本 - 支付系统完整指南

> 本文档详细解析了支付模块的架构、流程、实现细节和配置方式，适合想要深入理解该系统的开发者。

**最后更新：2025-12-15（修复：回调URL构成机制）**
**文档版本：1.1**

---

## 目录

1. [系统架构](#系统架构)
2. [核心概念](#核心概念)
3. [完整支付流程](#完整支付流程)
4. [数据库模型](#数据库模型)
5. [配置管理](#配置管理)
6. [代码实现深度解析](#代码实现深度解析)
7. [支付回调机制](#支付回调机制)
8. [定时任务系统](#定时任务系统)
9. [支付客户端详解](#支付客户端详解)
10. [API 端点清单](#api-端点清单)
11. [常见问题](#常见问题)
12. [学习路径](#学习路径)

---

## 系统架构

### 1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        管理后台（前端）                          │
│              Vivo Vue3 + Ant Design / Element Plus               │
└──────────────────────────┬──────────────────────────────────────┘
                           │ HTTP
                           ↓
┌─────────────────────────────────────────────────────────────────┐
│                     API Gateway & Router                        │
│         /admin-api/pay/* 和 /api/app/* 路由分发                 │
└──────────────────────────┬──────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        ↓                  ↓                  ↓
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│  应用管理       │ │  渠道管理       │ │  订单管理       │
│  (PayApp)      │ │  (PayChannel)   │ │  (PayOrder)    │
│  处理器/服务   │ │  处理器/服务    │ │  处理器/服务   │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                             ↓
                ┌────────────────────────────┐
                │  PayOrderService           │
                │  - 创建支付订单            │
                │  - 提交支付（调用渠道）   │
                │  - 处理支付回调            │
                └────────────┬───────────────┘
                             │
        ┌────────────────────┼─────────────────────┐
        ↓                    ↓                     ↓
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│  PayClientFactory │ │ PayNotifyService │ │  PayRefundService│
│  支付客户端工厂  │ │  回调通知服务    │ │  退款服务        │
└────────┬─────────┘ └────────┬─────────┘ └──────────────────┘
         │                    │
    ┌────┴────┐               │
    ↓         ↓               ↓
┌─────────┐ ┌─────────┐ ┌──────────────────┐
│ 微信支付 │ │ 支付宝  │ │ Redis + 定时任务│
│ 客户端  │ │ 客户端  │ │ (ExecuteNotify) │
└────┬────┘ └────┬────┘ └────────┬─────────┘
     │           │               │
     └─────┬─────┘               │
           ↓                     ↓
    ┌─────────────────┐   ┌─────────────────┐
    │ 微信/支付宝API  │   │ HTTP POST 回调  │
    │ (第三方服务)    │   │ (商户应用)      │
    └─────────────────┘   └─────────────────┘
```

### 1.2 分层架构

```
API 层 (Handlers)
├── internal/api/handler/admin/pay/
│   ├── app.go              # 支付应用配置
│   ├── channel.go          # 支付渠道配置
│   ├── order.go            # 支付订单查询
│   ├── refund.go           # 退款管理
│   ├── notify.go           # 回调通知管理
│   └── wallet/             # 钱包相关
│
业务服务层 (Services)
├── internal/service/pay/
│   ├── order.go            # PayOrderService (核心)
│   ├── app.go              # PayAppService
│   ├── channel.go          # PayChannelService
│   ├── notify.go           # PayNotifyService
│   ├── refund.go           # PayRefundService
│   ├── notify_lock.go      # Redis 分布式锁
│   ├── client/
│   │   ├── factory.go      # 客户端工厂
│   │   ├── client.go       # PayClient 接口
│   │   ├── weixin/         # 微信支付实现
│   │   └── alipay/         # 支付宝实现
│   └── consts.go           # 常量定义
│
数据访问层 (Repository)
├── internal/repo/query/
│   ├── pay_order.gen.go
│   ├── pay_channel.gen.go
│   ├── pay_app.gen.go
│   ├── pay_notify_task.gen.go
│   └── ...
│
模型层 (Models)
├── internal/model/pay/
│   ├── pay_order.go
│   ├── pay_channel.go
│   ├── pay_app.go
│   ├── pay_notify.go       # PayNotifyTask & PayNotifyLog
│   ├── pay_client_config.go
│   └── ...
│
数据库
└── MySQL pay_* 表
```

---

## 核心概念

### 2.1 关键术语

| 术语 | 含义 | 例子 |
|------|------|------|
| **PayApp** | 支付应用 | 我的商城应用、我的官网支付应用 |
| **PayChannel** | 支付渠道 | 微信公众号、支付宝扫码、微信小程序 |
| **PayOrder** | 支付订单 | 用户发起支付的订单记录 |
| **PayClientFactory** | 支付客户端工厂 | 管理不同支付渠道的客户端实例 |
| **PayClient** | 支付客户端 | 与微信/支付宝 API 对接的客户端 |
| **UnifiedOrder** | 统一下单 | 调用第三方支付平台下单 API |
| **PayNotifyTask** | 回调通知任务 | 待发送给商户的回调任务 |
| **ExecuteNotify** | 执行回调 | 定时任务发送 HTTP 回调到商户应用 |

### 2.2 状态机

#### PayOrder 订单状态

```
创建时默认状态: 0 (待支付)

0 (WAITING)        ─── 用户支付成功 ──→  10 (SUCCESS)
待支付              第三方通知或查询           已支付
  │
  └─────── 订单过期 ──→  20 (CLOSED) ─── 退款 ──→ 30 (REFUND)
           或关闭          已关闭        成功        已退款
```

#### PayNotifyTask 回调任务状态

```
创建时: 0 (WAITING)
等待通知

    ↓ 执行回调

    ├─ 成功 (200 + {"code": 0})
    │  └─→ 10 (SUCCESS) [完成]
    │
    └─ 失败
       ├─ notifyTimes < maxNotifyTimes?
       │  ├─ Yes: 状态回到 0 (WAITING)
       │  │       设置 nextNotifyTime = 当前 + 延迟时间
       │  │       延迟时间来自 NotifyFrequency 数组
       │  │
       │  └─ No: 状态变为 20 (FAILURE) [彻底失败]
```

### 2.3 关键常量

```go
// 通知类型
const (
    PayNotifyTypeOrder  = 1    // 支付订单通知
    PayNotifyTypeRefund = 2    // 退款单通知
)

// 订单状态
const (
    PayOrderStatusWaiting = 0   // 待支付
    PayOrderStatusSuccess = 10  // 已支付
    PayOrderStatusClosed  = 20  // 已关闭
    PayOrderStatusRefund  = 30  // 已退款
)

// 通知状态
const (
    PayNotifyStatusWaiting = 0   // 等待通知
    PayNotifyStatusSuccess = 10  // 通知成功
    PayNotifyStatusFailure = 20  // 通知失败
)

// 通知重试频率（单位：秒）
var NotifyFrequency = []int{15, 15, 30, 180, 1800, 1800, 1800, 3600}
// 说明: 共 8 次，总耗时约 12 小时
// 第1次: 15秒后 (订单成功后15秒)
// 第2次: 15秒后 (距上次)
// 第3次: 30秒后
// 第4次: 3分钟后
// 第5-8次: 逐步延长，最后一次间隔 1小时
```

---

## 完整支付流程

### 3.1 业务流程（E2E）

```
时间线:
0s     10s    30s    60s    所有商户应用
└──────┴──────┴──────┴──────
            ↑
        用户支付

┌──────────────────────────────────────────────────────────────┐
│ 阶段1: 配置准备（通常一次性）                                 │
└──────────────────────────────────────────────────────────────┘

管理员后台
  ↓
1. 创建支付应用 (POST /admin-api/pay/app/create)
   输入:
   - 应用名称: "我的商城"
   - 应用标识: "mall_app_001"
   - 订单回调地址: "https://yourapp.com/api/pay/notify/order"
   - 退款回调地址: "https://yourapp.com/api/pay/notify/refund"
   ↓
   生成: PayApp.ID = 1, AppKey = "xxxxx"

2. 创建支付渠道 (POST /admin-api/pay/channel/create)
   输入:
   - 关联应用: PayApp.ID = 1
   - 渠道编码: "wx_pub" (微信公众号)
   - 渠道名称: "微信公众号"
   - 支付参数 (JSON):
     {
       "@class": "ConfigTypeWxPay",
       "appId": "wxxxxxxxxx",        // 微信 AppID
       "mchId": "1xxxx",             // 商户号
       "apiVersion": "v3",
       "apiv3Key": "xxxxxxxx",       // APIv3 密钥
       "certSerialNo": "xxx",        // 证书序列号
       "privateKeyContent": "-----BEGIN...", // 商户私钥
       "publicKeyContent": "-----BEGIN..."   // 微信公钥
     }
   - 费率: 0.6 (%)
   ↓
   生成: PayChannel.ID = 1, Code = "wx_pub"

┌──────────────────────────────────────────────────────────────┐
│ 阶段2: 用户支付（实时交互）                                   │
└──────────────────────────────────────────────────────────────┘

用户前端应用
  ↓
1. 用户下单 (创建商城订单)
   - 调用: POST /api/app/trade/order/create
   - 系统创建: TradeOrder 记录
   - 返回: orderId, totalPrice 等
   ↓

2. 用户点击支付按钮
   - 调用: POST /api/app/trade/order/{orderId}/pay
   - 前端需要传递:
     {
       "payChannelCode": "wx_pub",     // 选择的支付渠道
       "price": 1000,                  // 支付金额（分）
       "openId": "oxxxxxxxxxx",        // JSAPI 需要
       "clientIP": "192.168.1.1"
     }
   ↓

3. 创建支付订单 [PayOrderService.CreateOrder()]
   - 检查应用是否启用
   - 检查幂等性: 相同商户订单号只创建一次
   - 创建 PayOrder 记录:
     {
       appId: 1,
       merchantOrderId: "SO20231215001",
       price: 1000,
       status: 0,                      // 待支付
       notifyUrl: "https://yourapp.com/api/pay/notify/order"  // 基础URL
     }
   - 返回: payOrderId = 123
   ↓

4. 提交支付 [PayOrderService.SubmitOrder()]
   输入:
   {
     payOrderId: 123,
     channelCode: "wx_pub",
     userIP: "192.168.1.1"
   }

   步骤:
   a) 校验订单状态 (必须是待支付)
   b) 获取支付渠道配置 (PayChannel)
   c) 创建 PayOrderExtension (扩展记录)
   d) 获取或创建支付客户端 (Factory)
   e) ⭐ 生成最终的回调 URL（追加 channelId）:
      最终 URL = 基础URL + "/" + channel.ID
      例如: "https://yourapp.com/api/pay/notify/order/1"
   f) 调用客户端 UnifiedOrder()，传入最终 URL:
      - 微信: 调用官方 SDK → 返回二维码 URL
      - 支付宝: 调用官方 SDK → 返回收银台 URL
   g) 返回支付信息给前端

   返回内容示例:
   {
     "displayMode": "qr_code",           // 二维码
     "displayContent": "https://..."     // 二维码 URL
     // 或
     "displayMode": "redirect",          // 跳转
     "displayContent": "https://..."     // 支付宝页面
   }
   ↓

5. 前端展示支付信息
   - 二维码: 使用微信/支付宝客户端扫描支付
   - 页面: 跳转到支付宝收银台
   ↓

6. 用户完成支付
   - 在微信/支付宝客户端输入密码/验证
   - 支付宝/微信服务器处理支付
   - 扣款成功
   ↓

┌──────────────────────────────────────────────────────────────┐
│ 阶段3: 支付回调（异步）                                        │
└──────────────────────────────────────────────────────────────┘

微信/支付宝服务器
  ↓ (异步通知)
  ↓
1. 微信/支付宝发起 POST 请求
   - URL: 您应用中的回调处理器（需要实现）
     例如: https://yourapp.com/api/pay/notify/order/1
     (基础URL + "/" + channelId)
   - 数据: 支付成功通知 (含签名)
   ↓

2. 应用接收并验证签名 (需实现)
   - 从URL中提取 channelId
   - 根据 channelId 获取对应的支付客户端
   - 验证签名有效性（RSA2 或 SHA256WithRSA）
   - 防止恶意请求
   ↓

3. 更新 PayOrder 状态 (需实现)
   - PayOrder.status = 10 (已支付)
   - PayOrder.successTime = 当前时间
   - 保存渠道返回数据到 PayOrderExtension
   ↓

4. 创建回调通知任务 [PayNotifyService.CreatePayNotifyTask()]
   - 生成 PayNotifyTask 记录:
     {
       type: 1,                         // 订单通知
       dataId: 123,                      // PayOrder.ID
       status: 0,                        // 等待通知
       notifyUrl: "https://yourapp.com/api/pay/notify/order",  // ⭐ 基础URL（不含channelId）
       notifyTimes: 0,
       nextNotifyTime: 当前时间
     }
   ↓

5. 立即发送回调（或由定时任务发送）
   [PayNotifyService.ExecuteNotify()]

   流程:
   a) 查询待通知的 PayNotifyTask
   b) 获取分布式锁 (Redis) 防止并发
   c) 构建回调数据 (PayOrder 的相关信息)
   d) ⭐ 重新追加 channelId 生成最终URL:
      最终 URL = notifyUrl + "/" + payOrder.channel_id
      例如: "https://yourapp.com/api/pay/notify/order/1"
   e) HTTP POST 发送到最终 URL
      请求头: Content-Type: application/json
      超时: 10 秒

   f) 接收响应:
      - 成功标志: HTTP 200 + 响应体包含 "SUCCESS"
      - PayNotifyTask.status = 10 (完成)

      - 失败:
        - notifyTimes < 8?
          ├─ Yes: 设置 nextNotifyTime，状态回到等待
          │       下次由定时任务重新尝试
          └─ No: status = 20 (彻底失败)

   f) 记录日志到 PayNotifyLog
      {
        taskId: PayNotifyTask.ID,
        notifyTimes: 1,
        response: "HTTP 响应体",
        status: 10/22  // 成功/失败
      }
   ↓

6. 商户应用处理回调
   - 验证签名（可选）
   - 更新本地订单状态为已支付
   - 发货、积分等后续业务流程
   - 响应 "SUCCESS" 或 JSON {"code": 0}
   ↓

┌──────────────────────────────────────────────────────────────┐
│ 阶段4: 对账和异常处理                                          │
└──────────────────────────────────────────────────────────────┘

管理员后台
  ↓
1. 查看支付订单列表
   - GET /admin-api/pay/order/page
   - 筛选、搜索、导出
   ↓

2. 查看回调通知日志
   - GET /admin-api/pay/notify/page
   - 查看每笔订单的回调历史
   - 显示: 通知次数、状态、最后通知时间
   ↓

3. 手动重试失败的回调（如果实现）
   - 点击按钮重新发送回调
   - 系统重置 PayNotifyTask.status = 0
   - 下次定时任务执行时会重新发送
```

### 3.2 时序图

```
用户          前端应用        后端应用        PayOrderService   微信/支付宝    商户应用
│            │              │              │                │            │
│ 1. 下单      │              │              │                │            │
├─────────────→│              │              │                │            │
│            │ 2. 创建订单     │              │                │            │
│            ├──────────────→│              │                │            │
│            │              │ 3. CreateOrder│                │            │
│            │              ├──────────────→│                │            │
│            │              │←──────────────┤ 返回 payOrderId │            │
│            │←──────────────┤              │                │            │
│ 4. 点击支付  │              │              │                │            │
├─────────────→│              │              │                │            │
│            │ 5. SubmitOrder│              │                │            │
│            ├──────────────→│              │                │            │
│            │              │ 6. SubmitOrder│                │            │
│            │              ├──────────────→│                │            │
│            │              │              │ 7. UnifiedOrder │            │
│            │              │              ├───────────────→│            │
│            │              │              │← 返回支付凭证   │            │
│            │              │←──────────────┤                │            │
│            │← 支付信息     │              │                │            │
│            │              │              │                │            │
│ 8. 扫描二维码│              │              │                │            │
├───────────────────────────────────────────→│                │            │
│            │              │              │ 9. 支付处理     │            │
│            │              │              │                │ 扣款等      │
│            │              │              │←─────────────────┤ 成功       │
│            │              │              │                │            │
│            │              │              │ 10. 异步回调 (微信主动)      │
│            │              │              │←─────────────────┤            │
│            │              │ 11. 接收回调  │                │            │
│            │              │←──────────────┤                │            │
│            │              │              │                │            │
│            │              │ 12. 更新支付订单 + 创建通知任务│            │
│            │              │              │                │            │
│            │              │ 13. ExecuteNotify (定时任务) ─────────────→│
│            │              │              │                │            │
│            │              │              │                │            │
│            │              │              │                │ 14. 处理回调│
│            │              │              │                ├─────────────→│
│            │              │              │                │ 返回 SUCCESS│
│            │              │              │                │←─────────────┤
│            │              │              │                │            │
│            │              │  日志记录     │                │            │
│            │              │ PayNotifyLog  │                │            │
│            │              │              │                │            │
│ 15. 查询订单 │              │              │                │            │
├─────────────→│              │              │                │            │
│ 已支付 ✓     │←──────────────┤              │                │            │
```

---

## 数据库模型

### 4.1 核心表结构

#### `pay_app` - 支付应用表

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `id` | BIGINT | 主键 | 1 |
| `app_key` | VARCHAR(32) | 应用标识 (唯一) | `mall_app_001` |
| `name` | VARCHAR(255) | 应用名称 | `我的商城` |
| `status` | INT | 状态 (0=启用, 1=禁用) | 0 |
| `order_notify_url` | VARCHAR(255) | 订单支付成功回调地址 | `https://yourapp.com/notify` |
| `refund_notify_url` | VARCHAR(255) | 退款成功回调地址 | `https://yourapp.com/refund` |
| `transfer_notify_url` | VARCHAR(255) | 转账成功回调地址 | `https://yourapp.com/transfer` |
| `remark` | VARCHAR(500) | 备注 | - |
| `create_time` | DATETIME | 创建时间 | 2023-12-01 10:00:00 |
| `update_time` | DATETIME | 更新时间 | 2023-12-15 14:30:00 |

**唯一索引:** `unique_idx_app_key` on `(app_key)`

---

#### `pay_channel` - 支付渠道表

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `id` | BIGINT | 主键 | 1 |
| `app_id` | BIGINT | 关联应用 ID | 1 |
| `code` | VARCHAR(32) | 渠道编码 (唯一) | `wx_pub` |
| `name` | VARCHAR(255) | 渠道名称 | `微信公众号` |
| `status` | INT | 状态 (0=启用) | 0 |
| `fee_rate` | DECIMAL(10,2) | 手续费率 (%) | 0.6 |
| `config` | LONGTEXT | 支付参数 (JSON) | `{"@class":"ConfigTypeWxPay",...}` |
| `remark` | VARCHAR(500) | 备注 | - |
| `create_time` | DATETIME | 创建时间 | 2023-12-01 10:00:00 |
| `update_time` | DATETIME | 更新时间 | 2023-12-15 14:30:00 |

**索引:**
- `idx_app_id` on `(app_id)`
- `unique_idx_code` on `(code)`

**config JSON 示例 - 微信:**

```json
{
  "@class": "ConfigTypeWxPay",
  "appId": "wxxxxxxxxx",
  "mchId": "1xxxx",
  "apiVersion": "v3",
  "apiv3Key": "32字符的APIv3密钥",
  "certSerialNo": "证书序列号",
  "privateKeyContent": "-----BEGIN PRIVATE KEY-----\nMIIEvQIB...",
  "publicKeyContent": "-----BEGIN PUBLIC KEY-----\nMFwwDQYJKoZIh..."
}
```

**config JSON 示例 - 支付宝:**

```json
{
  "@class": "ConfigTypeAlipay",
  "appId": "2xxxxxxxxx",
  "serverUrl": "https://openapi.alipay.com",
  "signType": "RSA2",
  "mode": 1,
  "privateKey": "MIIEvQIBADANB...",
  "alipayPublicKey": "MFwwDQYJKoZIh..."
}
```

---

#### `pay_order` - 支付订单表

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `id` | BIGINT | 主键 | 123 |
| `app_id` | BIGINT | 应用 ID | 1 |
| `channel_id` | BIGINT | 渠道 ID (可能为空) | 1 |
| `channel_code` | VARCHAR(32) | 渠道编码 | `wx_pub` |
| `user_id` | BIGINT | 用户 ID | 100 |
| `user_type` | INT | 用户类型 (0=会员) | 0 |
| `merchant_order_id` | VARCHAR(64) | 商户订单号 (唯一) | `SO20231215001` |
| `subject` | VARCHAR(255) | 订单主题 | `购买课程套餐` |
| `body` | TEXT | 订单描述 | `[课程1]x1, [课程2]x2` |
| `notify_url` | VARCHAR(255) | 回调地址 | `https://yourapp.com/notify` |
| `price` | INT | 支付金额 (分) | 1000 (即¥10.00) |
| `channel_fee_rate` | DECIMAL(10,2) | 渠道费率 (%) | 0.6 |
| `channel_fee_price` | INT | 渠道手续费 (分) | 6 |
| `status` | INT | 订单状态 | 10 (已支付) |
| `user_ip` | VARCHAR(32) | 用户 IP | `192.168.1.1` |
| `expire_time` | DATETIME | 过期时间 | 2023-12-15 12:00:00 |
| `success_time` | DATETIME | 支付成功时间 | 2023-12-15 10:30:00 |
| `extension_id` | BIGINT | 扩展记录 ID | 456 |
| `no` | VARCHAR(32) | 订单号 | `order_20231215_001` |
| `refund_price` | INT | 退款金额 (分) | 0 |
| `channel_user_id` | VARCHAR(255) | 渠道用户 ID | - |
| `channel_order_no` | VARCHAR(255) | 渠道订单号 | `4200001234567890` |
| `tenant_id` | BIGINT | 租户 ID | 1 |
| `create_time` | DATETIME | 创建时间 | 2023-12-15 10:00:00 |
| `update_time` | DATETIME | 更新时间 | 2023-12-15 10:30:00 |

**索引:**
- `unique_idx_merchant_order_id` on `(merchant_order_id, app_id)`
- `idx_app_id_status` on `(app_id, status)`
- `idx_channel_code` on `(channel_code)`

---

#### `pay_notify_task` - 回调通知任务表

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `id` | BIGINT | 主键 | 789 |
| `app_id` | BIGINT | 应用 ID | 1 |
| `type` | INT | 通知类型 (1=订单, 2=退款) | 1 |
| `data_id` | BIGINT | 数据 ID (PayOrder.ID) | 123 |
| `merchant_order_id` | VARCHAR(64) | 商户订单号 | `SO20231215001` |
| `merchant_refund_id` | VARCHAR(64) | 商户退款号 | NULL |
| `notify_url` | VARCHAR(255) | 回调地址 | `https://yourapp.com/notify` |
| `status` | INT | 任务状态 (0=等待, 10=成功, 20=失败) | 10 |
| `notify_times` | INT | 已通知次数 | 1 |
| `max_notify_times` | INT | 最大通知次数 | 9 |
| `next_notify_time` | DATETIME | 下次通知时间 | NULL (已完成) |
| `last_execute_time` | DATETIME | 最后执行时间 | 2023-12-15 10:30:10 |
| `create_time` | DATETIME | 创建时间 | 2023-12-15 10:30:00 |
| `update_time` | DATETIME | 更新时间 | 2023-12-15 10:30:10 |

**索引:**
- `idx_status_next_time` on `(status, next_notify_time)` ← 关键！用于查询待通知
- `idx_data_id` on `(data_id)`

---

#### `pay_notify_log` - 回调通知日志表

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `id` | BIGINT | 主键 | 1001 |
| `task_id` | BIGINT | 任务 ID | 789 |
| `notify_times` | INT | 第 N 次通知 | 1 |
| `status` | INT | 本次结果状态 | 10 |
| `response` | LONGTEXT | HTTP 响应内容 | `{"code":0,"msg":"success"}` |
| `create_time` | DATETIME | 创建时间 | 2023-12-15 10:30:10 |

**索引:**
- `idx_task_id` on `(task_id)`

---

#### `pay_order_extension` - 支付订单扩展表

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `id` | BIGINT | 主键 | 456 |
| `order_id` | BIGINT | 订单 ID | 123 |
| `channel_id` | BIGINT | 渠道 ID | 1 |
| `channel_code` | VARCHAR(32) | 渠道编码 | `wx_pub` |
| `no` | VARCHAR(32) | 唯一订单号 | `order_20231215_001` |
| `user_ip` | VARCHAR(32) | 用户 IP | `192.168.1.1` |
| `status` | INT | 状态 | 10 |
| `channel_extras` | LONGTEXT | 渠道额外数据 (JSON) | `{"prepayId":"xxx","qrUrl":"..."}` |
| `channel_notify_data` | LONGTEXT | 渠道通知原始数据 | `{"transaction_id":"..."}` |
| `create_time` | DATETIME | 创建时间 | 2023-12-15 10:00:00 |
| `update_time` | DATETIME | 更新时间 | 2023-12-15 10:30:00 |

---

### 4.2 ER 图

```
PayApp (1)
    │
    │ 1:N
    ├─────→ PayChannel (支持多种渠道)
    │
    └─────→ PayOrder (应用的所有订单)
                │
                │ 1:1
                ├─────→ PayOrderExtension (渠道相关数据)
                │
                └─────→ PayNotifyTask (订单相关回调任务)
                            │
                            │ 1:N
                            └─────→ PayNotifyLog (每次回调的日志)
```

---

## 配置管理

### 5.1 应用配置（PayApp）

管理员在后台新建应用时需要配置回调地址的基础 URL：

```json
{
  "appKey": "mall_app_001",
  "name": "我的商城",
  "status": 0,  // 0=启用, 1=禁用
  "orderNotifyURL": "https://yourdomain.com/api/pay/notify/order",
  "refundNotifyURL": "https://yourdomain.com/api/pay/notify/refund",
  "transferNotifyURL": "https://yourdomain.com/api/pay/notify/transfer"
}
```

**⭐ 关键点：回调URL的实际构成**

这些 URL 是**基础地址**，最终发送给第三方支付平台时会追加 `{channelId}`：

```
配置的基础URL:
  https://yourdomain.com/api/pay/notify/order

实际调用第三方时的URL:
  微信公众号 (channel.ID=1):   https://yourdomain.com/api/pay/notify/order/1
  支付宝 (channel.ID=2):       https://yourdomain.com/api/pay/notify/order/2
  微信小程序 (channel.ID=3):   https://yourdomain.com/api/pay/notify/order/3
```

**重要字段说明：**
- `orderNotifyURL`: 支付订单成功时的回调基础地址（追加 channelId）
- `refundNotifyURL`: 退款订单完成时的回调基础地址（追加 channelId）
- `transferNotifyURL`: 转账完成时的回调基础地址（追加 channelId）
- 这些回调地址是商户应用的处理端点，需要返回 `SUCCESS` 或 `{"code": 0}` 表示成功
- **同一个应用，所有支付渠道都通过同一个基础 URL 回调，只是通过 channelId 区分不同的渠道**

---

### 5.2 渠道配置（PayChannel）

#### 微信支付配置示例

```json
{
  "appId": 1,
  "code": "wx_pub",
  "name": "微信公众号",
  "status": 0,
  "feeRate": 0.6,
  "config": {
    "@class": "ConfigTypeWxPay",
    "appId": "wxxxxxxxxx",
    "mchId": "1640000000",
    "apiVersion": "v3",

    // APIv3 相关
    "apiv3Key": "32位的APIv3密钥",
    "certSerialNo": "证书序列号",

    // 商户私钥（用于签名）
    "privateKeyContent": "-----BEGIN PRIVATE KEY-----\nMIIEvQIB...\n-----END PRIVATE KEY-----",

    // 微信支付公钥（用于验证回调签名）
    "publicKeyContent": "-----BEGIN PUBLIC KEY-----\nMFwwDQYJKoZIh...\n-----END PUBLIC KEY-----",
    "publicKeyID": "cert_serial_no"
  }
}
```

**如何获取这些参数？**

1. **appId**: 微信公众平台 → 基本信息 → 公众号 AppID
2. **mchId**: 微信商户平台 → 账户中心 → 商户号
3. **apiVersion**: 根据使用的微信支付版本（通常 v3）
4. **apiv3Key**: 微信商户平台 → 账户设置 → APIv3 密钥
5. **certSerialNo**: 从下载的 API 证书中获取（使用 openssl 命令）
6. **privateKeyContent**: 从商户私钥文件复制（`apiclient_key.pem`）
7. **publicKeyContent**: 从微信公钥文件复制（定期更新）

#### 支付宝配置示例

```json
{
  "appId": 1,
  "code": "alipay_pc",
  "name": "支付宝电脑网站",
  "status": 0,
  "feeRate": 0.55,
  "config": {
    "@class": "ConfigTypeAlipay",
    "appId": "2xxxxxxxxx",
    "serverUrl": "https://openapi.alipay.com",
    "signType": "RSA2",
    "mode": 1,  // 1=公钥模式, 2=证书模式

    // 应用私钥
    "privateKey": "MIIEvQIBADANB...",

    // 支付宝公钥
    "alipayPublicKey": "MFwwDQYJKoZIh..."
  }
}
```

---

## 代码实现深度解析

### 6.1 PayOrderService - 核心服务

**文件位置:** `internal/service/pay/order.go`

#### 6.1.1 CreateOrder - 创建支付订单

```go
func (s *PayOrderService) CreateOrder(ctx context.Context, reqDTO *req.PayOrderCreateReq) (int64, error) {
    // 1. 校验应用
    app, err := s.appSvc.GetApp(ctx, reqDTO.AppID)
    if err != nil || app == nil || app.Status != 0 {
        return 0, errors.New("App disabled or not found")
    }

    // 2. 幂等性检查（确保同一商户订单号只被创建一次）
    existOrder, _ := s.q.PayOrder.WithContext(ctx).
        Where(s.q.PayOrder.AppID.Eq(app.ID), s.q.PayOrder.MerchantOrderId.Eq(reqDTO.MerchantOrderId)).
        First()
    if existOrder != nil {
        return existOrder.ID, nil  // 返回已存在的订单 ID
    }

    // 3. 创建新订单
    order := &pay.PayOrder{
        AppID:           app.ID,
        MerchantOrderId: reqDTO.MerchantOrderId,
        Subject:         reqDTO.Subject,
        Body:            reqDTO.Body,
        NotifyURL:       app.OrderNotifyURL,  // 使用应用的回调地址
        Price:           reqDTO.Price,
        ExpireTime:      time.Now().Add(2 * time.Hour),
        Status:          PayOrderStatusWaiting,  // 初始状态：待支付
        RefundPrice:     0,
        UserIP:          reqDTO.UserIP,
    }

    if err := s.q.PayOrder.WithContext(ctx).Create(order); err != nil {
        return 0, err
    }
    return order.ID, nil
}
```

**关键点：**
- **幂等性**: 使用 `merchant_order_id` 作为唯一标识，同一商户订单号只创建一次
- **回调地址保存**: 从应用配置中复制回调地址到 PayOrder 记录（作为备份）
- **2小时有效期**: 订单 2 小时后过期，不再允许支付
- **注意**: 这里保存的 `NotifyURL` 是基础地址，实际发送给第三方时会在 `SubmitOrder` 中追加 `channelId`

---

#### 6.1.2 SubmitOrder - 提交支付

```go
func (s *PayOrderService) SubmitOrder(ctx context.Context, reqVO *req.PayOrderSubmitReq, userIP string) (*resp.PayOrderSubmitResp, error) {
    // 1. 校验订单
    order, err := s.validateOrderCanSubmit(ctx, reqVO.ID)
    if err != nil {
        return nil, err  // 订单不存在、已支付、已过期等
    }

    // 2. 校验渠道
    channel, err := s.validateChannelCanSubmit(ctx, order.AppID, reqVO.ChannelCode)
    if err != nil {
        return nil, err  // 渠道不存在、被禁用等
    }

    // 3. 生成唯一订单号
    no := s.generateNo()  // 格式: "order_20231215_000001" 等

    // 4. 创建扩展记录（保存渠道相关数据）
    ext := &pay.PayOrderExtension{
        OrderID:     order.ID,
        No:          no,
        ChannelID:   channel.ID,
        ChannelCode: channel.Code,
        UserIP:      userIP,
        Status:      PayOrderStatusWaiting,
    }
    if err := s.q.PayOrderExtension.WithContext(ctx).Create(ext); err != nil {
        return nil, err
    }

    // 5. 获取支付客户端
    payClient := s.clientFac.GetPayClient(channel.ID)
    if payClient == nil {
        // 如果客户端不存在，创建新的
        payClient, err = s.clientFac.CreateOrUpdatePayClient(channel.ID, channel.Code, channel.Config)
        if err != nil {
            return nil, fmt.Errorf("failed to create pay client: %w", err)
        }
    }

    // 6. 调用渠道的统一下单接口
    unifiedReq := &client.UnifiedOrderReq{
        OutTradeNo:    no,
        Subject:       order.Subject,
        Body:          order.Body,
        Price:         int64(order.Price),
        NotifyURL:     order.NotifyURL,
        ChannelExtras: reqVO.ChannelExtras,  // 如 JSAPI 的 openid
    }

    orderResp, err := payClient.UnifiedOrder(ctx, unifiedReq)
    if err != nil {
        return nil, fmt.Errorf("failed to call unified order: %w", err)
    }

    // 7. 保存渠道返回的数据
    ext.ChannelExtras = orderResp.DisplayContent  // 二维码 URL、预支付 ID 等
    s.q.PayOrderExtension.WithContext(ctx).Save(ext)

    // 8. 返回支付展示信息给前端
    return &resp.PayOrderSubmitResp{
        OrderID:        order.ID,
        ExtensionID:    ext.ID,
        DisplayMode:    orderResp.DisplayMode,  // "qr_code", "redirect", "form" 等
        DisplayContent: orderResp.DisplayContent,  // 二维码URL、跳转链接等
    }, nil
}
```

**关键点：**
- **提交与创建分离**: 订单创建后，用户可以选择不同的渠道来支付（需要再次提交）
- **工厂模式获取客户端**: 使用 `PayClientFactory` 获取或创建支付客户端
- **⭐ 回调URL的最终构成**: 这里调用 `genChannelOrderNotifyUrl(channel)` 生成最终URL
  ```
  最终URL = config.C.Pay.OrderNotifyURL + "/" + channel.ID
  例如: https://yourdomain.com/api/pay/notify/order/1 (追加了 channelId=1)
  ```
- **渠道参数传递**: 通过 `ChannelExtras` 传递渠道特定的参数（如 JSAPI 的 openid）
- **显示内容保存**: 返回的支付凭证（二维码 URL、跳转链接等）保存在 `PayOrderExtension` 中

---

### 6.2 PayClientFactory - 客户端工厂

**文件位置:** `internal/service/pay/client/factory.go`

```go
type PayClientFactory struct {
    clients map[int64]PayClient      // channelID -> client 实例缓存
    mutex   sync.RWMutex             // 并发安全
}

// GetPayClient 获取已创建的客户端
func (f *PayClientFactory) GetPayClient(channelID int64) PayClient {
    f.mutex.RLock()
    defer f.mutex.RUnlock()
    return f.clients[channelID]
}

// CreateOrUpdatePayClient 创建或更新客户端
func (f *PayClientFactory) CreateOrUpdatePayClient(channelID int64, channelCode string, config string) (PayClient, error) {
    f.mutex.Lock()
    defer f.mutex.Unlock()

    // 查找对应的客户端创建器
    creator, ok := creators[channelCode]  // 全局注册器
    if !ok {
        return nil, errors.New("channel not supported")
    }

    // 创建新客户端实例
    newClient, err := creator(channelID, config)
    if err != nil {
        return nil, err
    }

    // 初始化（加载密钥等）
    if err := newClient.Init(); err != nil {
        return nil, err
    }

    // 缓存
    f.clients[channelID] = newClient
    return newClient, nil
}
```

**注册机制（客户端初始化时）：**

```go
// weixin/client.go init() 函数
func init() {
    client.RegisterCreator("wx_pub", NewWxPayClientAsClient)
    client.RegisterCreator("wx_lite", NewWxPayClientAsClient)
    client.RegisterCreator("wx_app", NewWxPayClientAsClient)
    client.RegisterCreator("wx_native", NewWxPayClientAsClient)
    client.RegisterCreator("wx_wap", NewWxPayClientAsClient)
    client.RegisterCreator("wx_bar", NewWxPayClientAsClient)
}

// alipay/client.go init() 函数
func init() {
    client.RegisterCreator("alipay_pc", NewAlipayClientAsClient)
    client.RegisterCreator("alipay_wap", NewAlipayClientAsClient)
    client.RegisterCreator("alipay_app", NewAlipayClientAsClient)
    // ...
}
```

**工作流程：**
1. 应用启动时，`init()` 函数注册所有支持的渠道创建器
2. 需要支付时，通过 `PayClientFactory.GetPayClient()` 获取缓存的客户端
3. 如果客户端不存在，调用 `CreateOrUpdatePayClient()` 创建新的
4. 根据 `channelCode` 查找对应的创建器
5. 创建器通过解析配置 JSON 并初始化密钥来创建客户端

---

### 6.3 支付客户端接口

**文件位置:** `internal/service/pay/client/client.go`

```go
type PayClient interface {
    // GetID 获取渠道编号
    GetID() int64

    // Init 初始化（加载密钥等）
    Init() error

    // UnifiedOrder 统一下单（返回支付凭证）
    UnifiedOrder(ctx context.Context, req *UnifiedOrderReq) (*OrderResp, error)

    // UnifiedRefund 统一退款
    UnifiedRefund(ctx context.Context, req *UnifiedRefundReq) (*RefundResp, error)

    // GetOrder 查询订单（对账用）
    GetOrder(ctx context.Context, outTradeNo string) (*OrderResp, error)

    // GetRefund 查询退款
    GetRefund(ctx context.Context, outTradeNo, outRefundNo string) (*RefundResp, error)

    // ParseOrderNotify 解析订单回调（重要！需要验证签名）
    ParseOrderNotify(req *NotifyData) (*OrderResp, error)

    // ParseRefundNotify 解析退款回调
    ParseRefundNotify(req *NotifyData) (*RefundResp, error)

    // UnifiedTransfer 统一转账
    UnifiedTransfer(ctx context.Context, req *UnifiedTransferReq) (*TransferResp, error)
}

// UnifiedOrderReq 下单请求
type UnifiedOrderReq struct {
    OutTradeNo    string            // 商户订单号（唯一）
    Subject       string            // 商品名称
    Body          string            // 商品描述
    Price         int64             // 金额（分）
    NotifyURL     string            // 回调地址
    ChannelExtras map[string]string // 渠道特定参数（如 openid）
}

// OrderResp 下单响应
type OrderResp struct {
    Status           int    // 订单状态
    OutTradeNo       string // 商户订单号
    DisplayMode      string // 展示模式："qr_code", "redirect", "form", "prepay_id"
    DisplayContent   string // 展示内容：二维码URL、跳转链接、预支付ID等
    ChannelOrderNo   string // 渠道订单号（第三方返回）
    ChannelErrorCode string // 渠道错误码
    ChannelErrorMsg  string // 渠道错误消息
}

// NotifyData 回调通知数据
type NotifyData struct {
    Params  map[string]string // Query 参数
    Body    string            // 请求体
    Headers map[string]string // HTTP 头
}
```

---

### 6.4 微信支付客户端实现

**文件位置:** `internal/service/pay/client/weixin/client.go`

#### 初始化流程

```go
func (c *WxPayClient) Init() error {
    // 1. 解析 JSON 配置
    var cfg WxPayClientConfig
    if err := json.Unmarshal([]byte(c.Config), &cfg); err != nil {
        return fmt.Errorf("解析微信支付配置失败: %w", err)
    }
    c.config = &cfg

    // 2. 仅支持 V3 版本
    if cfg.APIVersion == APIVersionV3 {
        return c.initV3Client()
    }
    return errors.New("暂不支持微信支付 V2 版本")
}

func (c *WxPayClient) initV3Client() error {
    cfg := c.config

    // 3. 加载商户私钥（用于请求签名）
    privateKey, err := utils.LoadPrivateKey(cfg.PrivateKeyContent)
    if err != nil {
        return fmt.Errorf("加载商户私钥失败: %w", err)
    }
    c.privateKey = privateKey

    // 4. 加载微信公钥（用于验证回调签名）
    publicKey, err := utils.LoadPublicKey(cfg.PublicKeyContent)
    if err != nil {
        return fmt.Errorf("加载微信支付公钥失败: %w", err)
    }
    c.publicKey = publicKey

    // 5. 创建微信支付 API 客户端
    opts := []core.ClientOption{
        option.WithWechatPayPublicKeyAuthCipher(
            cfg.MchID,
           cfg.CertSerialNo,
            privateKey,
            cfg.PublicKeyID,
            publicKey,
        ),
    }
    coreClient, err := core.NewClient(context.Background(), opts...)
    if err != nil {
        return fmt.Errorf("创建微信支付客户端失败: %w", err)
    }
    c.coreClient = coreClient

    return nil
}
```

#### UnifiedOrder - 统一下单

```go
func (c *WxPayClient) UnifiedOrder(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
    // 根据支付渠道类型选择不同的支付方式
    switch c.ChannelCode {
    case "wx_native":
        return c.nativeOrder(ctx, req)      // 扫码支付
    case "wx_pub", "wx_lite":
        return c.jsapiOrder(ctx, req)       // 公众号/小程序支付
    case "wx_wap":
        return c.h5Order(ctx, req)          // H5 网页支付
    case "wx_app":
        return c.appOrder(ctx, req)         // APP 支付
    default:
        return nil, fmt.Errorf("暂不支持的微信支付渠道: %s", c.ChannelCode)
    }
}

// Native 扫码支付示例
func (c *WxPayClient) nativeOrder(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
    svc := native.NativeApiService{Client: c.coreClient}

    // 调用微信 API 预支付
    resp, _, err := svc.Prepay(ctx, native.PrepayRequest{
        Appid:       core.String(c.config.AppID),
        Mchid:       core.String(c.config.MchID),
        Description: core.String(req.Subject),
        OutTradeNo:  core.String(req.OutTradeNo),
        NotifyUrl:   core.String(req.NotifyURL),
        Amount: &native.Amount{
            Total:    core.Int64(req.Price),
            Currency: core.String("CNY"),
        },
    })

    if err != nil {
        return &client.OrderResp{
            Status:           20,  // CLOSED
            OutTradeNo:       req.OutTradeNo,
            ChannelErrorCode: "NATIVE_PREPAY_ERROR",
            ChannelErrorMsg:  err.Error(),
        }, nil
    }

    // 返回二维码 URL
    return &client.OrderResp{
        Status:         0,  // WAITING
        OutTradeNo:     req.OutTradeNo,
        DisplayMode:    "qr_code",
        DisplayContent: *resp.CodeUrl,  // 微信返回的二维码 URL
    }, nil
}
```

#### ParseOrderNotify - 解析回调（关键！）

```go
func (c *WxPayClient) ParseOrderNotify(req *client.NotifyData) (*client.OrderResp, error) {
    // 1. 构造 HTTP Request
    httpReq, err := http.NewRequest("POST", "", bytes.NewBufferString(req.Body))
    if err != nil {
        return nil, err
    }

    // 2. 设置请求头
    for k, v := range req.Headers {
        httpReq.Header.Set(k, v)
    }

    // 3. 创建 RSA 验证器
    verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(
        c.config.PublicKeyID,  // 公钥 ID
        c.publicKey,           // 微信公钥
    )

    // 4. 创建通知处理器
    handler, err := notify.NewRSANotifyHandler(
        c.config.APIV3Key,     // APIv3 密钥
        verifier,
    )
    if err != nil {
        return nil, err
    }

    // 5. 解析并验证签名
    var transaction payments.Transaction
    notifyReq, err := handler.ParseNotifyRequest(context.Background(), httpReq, &transaction)
    if err != nil {
        return nil, fmt.Errorf("failed to parse notify: %w", err)
    }

    // 6. 提取关键数据
    status := 10  // SUCCESS
    if transaction.TradeState != nil && *transaction.TradeState == "CLOSED" {
        status = 20  // CLOSED
    }

    return &client.OrderResp{
        Status:         status,
        OutTradeNo:     *transaction.OutTradeNo,
        ChannelOrderNo: *transaction.TransactionId,
        SuccessTime:    transaction.SuccessTime,
    }, nil
}
```

**验证签名的安全性：**
- 微信使用 **SHA256WithRSA** 算法
- 请求头中的 `Wechatpay-Signature` 字段包含签名
- 使用微信公钥验证签名，确保请求确实来自微信
- 防止恶意请求冒充微信进行回调

---

### 6.5 支付宝客户端实现

**文件位置:** `internal/service/pay/client/alipay/client.go`

#### UnifiedOrder - 统一下单

```go
// 电脑网站支付示例
func (c *AliPayClient) tradePagePay(ctx context.Context, req *client.UnifiedOrderReq) (*client.OrderResp, error) {
    // 构建下单请求
    tradePagePayReq := &requests.TradePagePayRequest{
        NotifyUrl:   req.NotifyURL,
        ReturnUrl:   "",  // 支付完成后返回地址
        OutTradeNo:  req.OutTradeNo,
        TotalAmount: fmt.Sprintf("%.2f", float64(req.Price)/100),  // 分转元
        Subject:     req.Subject,
        Body:        req.Body,
        ProductCode: "FAST_INSTANT_TRADE_PAY",
    }

    // 调用支付宝 API
    resp, err := c.client.TradePagePay(ctx, tradePagePayReq)
    if err != nil {
        return nil, err
    }

    // 返回支付宝收银台 URL
    return &client.OrderResp{
        Status:         0,
        OutTradeNo:     req.OutTradeNo,
        DisplayMode:    "redirect",
        DisplayContent: resp.Body,  // 跳转 URL
    }, nil
}
```

#### ParseOrderNotify - 解析回调

```go
func (c *AliPayClient) ParseOrderNotify(req *client.NotifyData) (*client.OrderResp, error) {
    // 1. 转换参数为 url.Values
    params := url.Values{}
    for k, v := range req.Params {
        params.Add(k, v)
    }

    // 2. 验证签名（关键！）
    err := c.client.VerifySign(params)
    if err != nil {
        return nil, fmt.Errorf("signature verification failed: %w", err)
    }

    // 3. 提取关键字段
    tradeStatus := params.Get("trade_status")
    status := 10  // SUCCESS
    if tradeStatus == "TRADE_CLOSED" {
        status = 20  // CLOSED
    }

    // 4. 解析时间
    successTime, _ := time.Parse("2006-01-02 15:04:05", params.Get("gmt_payment"))

    return &client.OrderResp{
        Status:         status,
        OutTradeNo:     params.Get("out_trade_no"),
        ChannelOrderNo: params.Get("trade_no"),
        ChannelUserID:  params.Get("buyer_id"),
        SuccessTime:    &successTime,
    }, nil
}
```

**支付宝的特点：**
- 使用 URL Query 参数而非 JSON Body
- 使用 RSA2 算法验证签名
- 退款是同步操作，没有异步通知
- 交易状态通过 `trade_status` 字段判断

---

## 支付回调机制

### 7.1 完整的回调流程

```
时刻     操作                    数据库变化           说明
────────────────────────────────────────────────────────
T0      用户支付成功             PayOrder.status=10
                                支付成功

T0+1秒  第三方通知应用           ──                 微信/支付宝 POST 回调
        (需要你实现)

T0+2秒  应用验证签名并更新        PayOrderExtension
        PayOrder.status=10        .channel_notify_data

T0+3秒  CreatePayNotifyTask()     PayNotifyTask
                                 .status=0 (等待)

T0+4秒  立即发送回调或等待定时任务
        ExecuteNotify()

        ├─ 成功 (HTTP 200)        PayNotifyTask
        │ └─→ 回调完成             .status=10
        │
        └─ 失败                  PayNotifyTask
           ├─ notifyTimes < 8?    .status=0
           │  └─→ 设置延迟时间    .next_notify_time
           │
           └─ 重试次数用尽        PayNotifyTask
              └─→ status=20 (失败)
```

### 7.2 回调的安全性考虑

#### 签名验证

所有来自微信/支付宝的异步通知都必须验证签名：

```go
// 微信：使用 SHA256WithRSA
verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(publicKeyID, publicKey)
handler, _ := notify.NewRSANotifyHandler(apiV3Key, verifier)
_, err := handler.ParseNotifyRequest(ctx, httpReq, &transaction)
if err != nil {
    // 签名验证失败！可能是恶意请求
    return "FAILED"
}

// 支付宝：使用 RSA2
err := aliPayClient.VerifySign(params)
if err != nil {
    // 签名验证失败
    return "FAILED"
}
```

#### 幂等性

相同的回调可能被发送多次，需要确保幂等处理：

```go
// 方式1：使用唯一索引
// PayNotifyTask 中的 (task_id, notify_times) 唯一组合

// 方式2：使用 merchant_order_id
// 每个商户订单号只能产生一个支付成功记录

// 方式3：分布式锁
// 使用 Redis 的 SETNX 防止并发处理同一任务
```

---

## 定时任务系统

### 8.1 定时任务框架

**使用库:** `github.com/go-co-op/gocron/v2`

**文件位置:** `internal/service/scheduler.go`

```go
type Scheduler struct {
    scheduler gocron.Scheduler
    handlers  map[string]JobHandler  // 任务处理器
    jobMap    map[int64]gocron.Job   // 运行中的任务
}

type JobHandler interface {
    Execute(ctx context.Context, param string) error
}
```

### 8.2 支付回调通知任务

**需要实现的处理器：**

```go
type PayNotifyJobHandler struct {
    service *PayNotifyService
}

func (h *PayNotifyJobHandler) Execute(ctx context.Context, param string) error {
    // 执行回调通知
    count, err := h.service.ExecuteNotify(ctx)
    if err != nil {
        log.Errorf("Failed to execute notify: %v", err)
        return err
    }
    log.Infof("Executed %d notify tasks", count)
    return nil
}
```

**任务配置：**

需要在数据库中配置一个定时任务（`infra_job` 表）：

```sql
INSERT INTO infra_job (
    name,
    status,
    handler_name,
    cron_expression,
    handler_param,
    create_time
) VALUES (
    'pay_notify_execute',           -- 任务名
    0,                              -- 启用状态
    'payNotifyJobHandler',          -- 处理器名
    '0 * * * * ?',                  -- Cron 表达式：每分钟执行一次
    '',                             -- 参数
    NOW()
);
```

**Cron 表达式说明：**
```
秒  分  时  日  月  周  年
0   *   *   *   *   ?   *
```

- `0 * * * * ?` - 每分钟的第 0 秒执行
- `0 */5 * * * ?` - 每 5 分钟执行一次
- `0 0 * * * ?` - 每小时执行一次
- `0 0 0 * * ?` - 每天午夜执行

### 8.3 执行流程

```
定时任务触发（gocron）
    ↓
调用 PayNotifyJobHandler.Execute()
    ↓
PayNotifyService.ExecuteNotify()
    ├─ 查询待通知任务
    │  where: status=0 AND nextNotifyTime <= NOW()
    ├─ 对每个任务异步执行
    │  go executeNotifyTaskWithLock()
    └─ 返回执行任务数

executeNotifyTaskWithLock()
    ├─ 获取 Redis 分布式锁
    ├─ 双重检查（防止并发重复执行）
    ├─ 执行实际通知逻辑
    └─ 释放锁

executeNotifyTask()
    ├─ HTTP POST 发送回调
    ├─ 记录日志
    ├─ 更新任务状态
    └─ 计划下一次通知（如果失败）
```

---

## 支付客户端详解

### 9.1 支持的支付渠道

| 渠道编码 | 渠道名称 | 描述 | 返回类型 |
|---------|---------|------|---------|
| **wx_pub** | 微信公众号 | 公众号 JSAPI 支付 | 预支付 ID |
| **wx_lite** | 微信小程序 | 小程序支付 | 预支付 ID |
| **wx_app** | 微信 APP | APP 中的微信支付 | 签名参数 |
| **wx_native** | 微信扫码 | 扫一扫支付 | 二维码 URL |
| **wx_wap** | 微信 H5 | 浏览器中的微信支付 | 支付 URL |
| **wx_bar** | 微信条码 | 商家扫顾客码 | 同步结果 |
| **alipay_pc** | 支付宝电脑 | 电脑网站支付 | 收银台 URL |
| **alipay_wap** | 支付宝手机 | 手机网站支付 | 收银台 URL |
| **alipay_app** | 支付宝 APP | APP 支付 | 签名参数 |
| **alipay_qr** | 支付宝二维码 | 扫码支付 | 二维码 URL |
| **alipay_bar** | 支付宝条码 | 商家扫顾客码 | 同步结果 |

### 9.2 微信支付 V3 API

**官方文档:** https://pay.weixin.qq.com/wiki

**接口调用流程：**

```
1. 初始化时
   - 加载商户私钥
   - 加载微信公钥
   - 创建 core.Client（自动处理签名和验证）

2. 下单时
   - 构建请求对象
   - 调用相应的 API
   - 返回响应

3. 回调时
   - 接收 HTTP 请求
   - 验证签名和时间戳
   - 解析 JSON 数据
   - 提取交易信息
```

**关键参数：**
- `appId`: 应用 ID（从微信公众平台获取）
- `mchId`: 商户号（从微信商户平台获取）
- `serialNo`: 证书序列号（从API证书获取）
- `privateKey`: 商户私钥（用于签名请求）
- `publicKey`: 微信公钥（用于验证回调）

### 9.3 支付宝 API

**官方文档:** https://opendocs.alipay.com

**接口调用流程：**

```
1. 初始化时
   - 加载应用私钥
   - 加载支付宝公钥
   - 创建 Client

2. 下单时
   - 构建请求对象
   - 调用相应的 API
   - 返回响应

3. 回调时
   - 接收 URL Query 参数
   - 验证签名
   - 提取交易信息
```

**关键参数：**
- `appId`: 应用 ID（从蚂蚁开放平台获取）
- `privateKey`: 应用私钥（用于签名）
- `publicKey`: 支付宝公钥（用于验证）
- `mode`: 1=公钥模式，2=证书模式

---

## API 端点清单

### 10.1 支付应用管理

```
POST   /admin-api/pay/app/create              创建应用
PUT    /admin-api/pay/app/update              编辑应用
PUT    /admin-api/pay/app/update-status       启用/禁用应用
DELETE /admin-api/pay/app/delete              删除应用
GET    /admin-api/pay/app/get                 获取应用详情
GET    /admin-api/pay/app/page                分页查询应用
GET    /admin-api/pay/app/list                列表查询应用
```

**请求/响应示例：**

```json
// POST /admin-api/pay/app/create
{
  "appKey": "mall_app_001",
  "name": "我的商城",
  "orderNotifyURL": "https://yourdomain.com/api/pay/notify/order",
  "refundNotifyURL": "https://yourdomain.com/api/pay/notify/refund"
}

// 200 OK
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "appKey": "mall_app_001",
    "name": "我的商城",
    "status": 0,
    "createTime": "2023-12-15T10:00:00Z"
  }
}
```

---

### 10.2 支付渠道管理

```
POST   /admin-api/pay/channel/create          创建渠道
PUT    /admin-api/pay/channel/update          编辑渠道配置
DELETE /admin-api/pay/channel/delete          删除渠道
GET    /admin-api/pay/channel/get             获取渠道详情
GET    /admin-api/pay/channel/get-enable-code-list  获取启用渠道列表
```

---

### 10.3 支付订单管理

```
GET    /admin-api/pay/order/get               获取订单
GET    /admin-api/pay/order/get-detail        获取订单详情（含扩展信息）
GET    /admin-api/pay/order/page              分页查询订单
POST   /admin-api/pay/order/submit            提交支付（核心端点）
```

---

### 10.4 回调通知管理

```
GET    /admin-api/pay/notify/get-detail       获取通知任务详情（含日志）
GET    /admin-api/pay/notify/page             分页查询通知任务
```

**查询通知日志示例：**

```json
// GET /admin-api/pay/notify/get-detail?id=789
{
  "code": 0,
  "data": {
    "task": {
      "id": 789,
      "appId": 1,
      "type": 1,
      "dataId": 123,
      "status": 10,        // 已成功
      "notifyTimes": 3,    // 通知了3次
      "nextNotifyTime": null,  // 已完成，无下次
      "logs": [
        {
          "id": 1001,
          "notifyTimes": 1,
          "status": 22,    // 请求失败
          "response": "Connection timeout",
          "createTime": "2023-12-15T10:30:05Z"
        },
        {
          "id": 1002,
          "notifyTimes": 2,
          "status": 22,    // 请求失败
          "response": "Connection refused",
          "createTime": "2023-12-15T10:30:20Z"
        },
        {
          "id": 1003,
          "notifyTimes": 3,
          "status": 10,    // 成功！
          "response": "{\"code\":0}",
          "createTime": "2023-12-15T10:31:00Z"
        }
      ]
    }
  }
}
```

---

## 常见问题

### Q1: 如何区分不同的支付状态？

**A:** 支付有四个主要状态：

| 状态值 | 名称 | 含义 | 可以操作什么 |
|--------|------|------|------------|
| 0 | WAITING | 待支付 | 继续支付、查询订单状态 |
| 10 | SUCCESS | 已支付 | 退款、发货 |
| 20 | CLOSED | 已关闭 | 重新下单 |
| 30 | REFUND | 已退款 | 显示退款信息 |

---

### Q2: 为什么我的回调通知没有收到？

**A:** 检查以下几点：

1. **回调地址是否正确？**
   - 在后台查看应用配置中的 `orderNotifyURL`
   - 确保地址可以外网访问（不能是 localhost）
   - 确保 HTTPS 证书有效

2. **您的应用是否返回了正确的响应？**
   - 必须返回 HTTP 200
   - 响应体包含 `SUCCESS` 或 `{"code": 0}`

3. **通知是否已发送？**
   - 查看后台的"回调通知"日志
   - 查看 `PayNotifyTask` 的状态
   - 如果状态是 0，说明还在等待通知
   - 如果状态是 20，说明已彻底失败

4. **定时任务是否在运行？**
   - 确认 `pay_notify_execute` 任务已启用
   - 检查服务器日志是否有任务执行记录

5. **网络连接问题？**
   - 检查防火墙是否开放了出站 HTTPS 端口
   - 检查回调地址是否正确路由

---

### Q3: 微信支付配置中的各个字段怎么获取？

**A:** 按以下步骤获取：

```
1. 登录微信商户平台
   https://pay.weixin.qq.com

2. 账户中心 → 商户信息
   - 获取: MchID（商户号）

3. 账户设置 → API 安全
   - 获取: APIv3 密钥（32位）

4. 账户设置 → API 证书
   - 下载: apiclient_cert.pem（证书）
   - 下载: apiclient_key.pem（私钥）
   - 获取: 证书序列号（从证书信息）

5. 微信公众平台
   https://mp.weixin.qq.com
   - 获取: AppID

6. 放置密钥文件
   - privateKeyContent: apiclient_key.pem 的内容
   - publicKeyContent: 微信公钥（周期性更新）
```

---

### Q4: 如何测试支付流程？

**A:** 建议按以下顺序测试：

```
1. 单元测试
   - 测试 PayOrderService 的创建、查询逻辑
   - 使用 Mock 对象模拟 PayClient

2. 集成测试
   - 测试完整的支付流程
   - 使用微信/支付宝的沙箱环境

3. 手动测试
   - 在后台创建应用和渠道
   - 创建支付订单
   - 通过二维码或链接完成支付
   - 检查订单状态和通知日志

4. 压力测试
   - 模拟大量并发支付
   - 观察回调通知的并发处理
```

---

### Q5: 如何处理支付失败？

**A:** 支付失败的原因和处理方式：

| 原因 | 检查点 | 解决方案 |
|------|--------|---------|
| 渠道配置错误 | 检查 PayChannel.config | 重新配置密钥 |
| 订单已支付 | 检查 PayOrder.status | 返回已支付提示 |
| 订单已过期 | 检查 PayOrder.expireTime | 重新下单 |
| 金额不匹配 | 检查 price 字段 | 确保金额准确 |
| 网络异常 | 查看错误日志 | 重试或检查网络 |
| 第三方服务故障 | 查看渠道状态 | 等待恢复或更换渠道 |

---

## 学习路径

### 推荐学习顺序

#### 第 1 阶段：理解基础概念（1-2小时）

1. **阅读本文档的前 3 章**
   - 系统架构
   - 核心概念
   - 完整支付流程

2. **阅读代码**
   - `internal/model/pay/pay_order.go` - 理解数据结构
   - `internal/service/pay/consts.go` - 理解常量定义

3. **查看数据库**
   - 运行 SQL 查看表结构
   - 理解各个字段的含义

#### 第 2 阶段：深入业务逻辑（2-3小时）

4. **研究 PayOrderService**
   - `internal/service/pay/order.go`
   - 理解 CreateOrder 和 SubmitOrder 的流程

5. **研究支付客户端**
   - `internal/service/pay/client/factory.go` - 工厂模式
   - `internal/service/pay/client/weixin/client.go` - 微信具体实现
   - `internal/service/pay/client/alipay/client.go` - 支付宝具体实现

6. **动手实践**
   - 在后台创建应用和渠道
   - 尝试创建支付订单
   - 调试支付客户端的初始化

#### 第 3 阶段：回调和通知（2-3小时）

7. **研究回调通知**
   - `internal/service/pay/notify.go` - 核心逻辑
   - `internal/service/pay/notify_lock.go` - 分布式锁

8. **实现缺失的部分**
   - 实现 PayOrderService 中的 `notifyOrder` 方法
   - 实现第三方回调的处理器
   - 实现定时任务处理器

9. **测试**
   - 测试回调的重试机制
   - 测试并发安全性
   - 观察日志和数据库变化

#### 第 4 阶段：高级主题（根据需要）

10. **性能优化**
    - 数据库查询优化
    - 缓存策略（使用 Redis）
    - 批量处理

11. **扩展功能**
    - 添加新的支付渠道
    - 实现退款功能
    - 实现转账功能
    - 添加钱包功能

12. **生产就绪**
    - 完整的错误处理
    - 详细的日志记录
    - 监控和告警
    - 灾难恢复

### 重要代码文件导航

| 文件 | 行数 | 学习优先级 | 说明 |
|------|------|-----------|------|
| `internal/model/pay/*` | 100+ | ⭐⭐⭐ | 数据模型，必读 |
| `internal/service/pay/order.go` | 350+ | ⭐⭐⭐ | 核心业务逻辑 |
| `internal/service/pay/client/factory.go` | 76 | ⭐⭐⭐ | 客户端工厂模式 |
| `internal/service/pay/client/weixin/client.go` | 508 | ⭐⭐ | 微信集成实现 |
| `internal/service/pay/client/alipay/client.go` | 429 | ⭐⭐ | 支付宝集成实现 |
| `internal/service/pay/notify.go` | 256 | ⭐⭐⭐ | 回调通知处理 |
| `internal/api/handler/admin/pay/order.go` | 100+ | ⭐⭐ | API 端点实现 |
| `internal/service/scheduler.go` | 195 | ⭐⭐ | 定时任务框架 |

### 测试清单

完成以下测试，确保理解了整个系统：

- [ ] 创建支付应用和渠道
- [ ] 理解回调 URL 的构成方式：基础 URL + "/" + channelId
- [ ] 创建支付订单
- [ ] 提交支付（调用渠道客户端）
- [ ] 验证支付订单数据入库
- [ ] 解析微信回调通知（使用沙箱）
- [ ] 验证签名正确性
- [ ] 测试回调任务的创建
- [ ] 手动触发 ExecuteNotify 并验证
- [ ] 验证重试机制（模拟回调失败）
- [ ] 查看后台管理界面的日志

---

## 总结

本文档提供了对芋道商城 Go 版本支付系统的全面认识：

### 架构特点

1. **清晰的分层设计** - API、Service、Repository、Model 层次分明
2. **支付客户端工厂** - 灵活支持多个支付渠道
3. **异步回调机制** - 可靠的通知重试策略
4. **分布式安全** - 使用 Redis 锁防止并发问题
5. **⭐ 统一回调入口** - 所有渠道共用一个基础 URL，通过 channelId 区分
   - 配置一个基础 URL：`https://yourdomain.com/api/pay/notify/order`
   - 最终发送给第三方的 URL：`https://yourdomain.com/api/pay/notify/order/{channelId}`
   - 简化部署、易于维护、灵活扩展

### 核心流程

1. **配置阶段** - 创建应用、配置渠道参数
2. **支付阶段** - 创建订单、提交支付、获取支付凭证
3. **用户操作** - 扫码或跳转进行支付
4. **回调阶段** - 第三方通知应用，更新订单状态
5. **通知阶段** - 定时任务回调商户，告知支付结果

### 关键设计模式

- **工厂模式** - PayClientFactory 管理不同渠道的客户端
- **策略模式** - 不同支付渠道有不同的实现
- **观察者模式** - 支付成功后触发一系列事件
- **重试机制** - 使用 Cron 表达式的定时任务

### 下一步

1. 完整阅读本文档
2. 在 IDE 中深入阅读推荐的代码文件
3. 在本地环境搭建并调试
4. 实现缺失的第三方回调处理器
5. 集成到实际业务中

---

**文档维护：** 本文档会随着代码演进而更新，最后更新时间为 2025-12-15。

