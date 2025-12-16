# 订单支付全生命周期流转图 (Order Payment Lifecycle)

本图展示了从用户下单到支付完成及状态同步的完整流程。

```mermaid
sequenceDiagram
    autonumber
    participant Users as 用户 (App端)
    participant Trade as 交易模块 (商城)
    participant Pay as 支付模块 (基础服务)
    participant 3rd as 第三方支付 (微信/支付宝)
    participant DB as 数据库

    Note over Users, Trade: 1. 下单流程
    Users->>Trade: 请求创建订单 (POST /trade/order/create)
    Trade->>DB: 插入交易订单 (状态: 未支付)
    Trade->>Pay: 调用 PayOrderService.CreatePayOrder()
    Note right of Trade: 设置回调地址 NotifyURL = <br/>/app-api/trade/order/update-paid
    Pay->>DB: 插入支付订单 (状态: 等待支付)
    Pay-->>Trade: 返回支付订单号 (PayOrderID)
    Trade-->>Users: 返回订单号 & 支付订单号

    Note over Users, Pay: 2. 支付提交
    Users->>Pay: 请求提交支付 (POST /pay/order/submit) <br/>(选择渠道: 微信/支付宝)
    Pay->>3rd: 统一下单处理 (UnifiedOrder)
    3rd-->>Pay: 返回预支付信息 (如 prepay_id)
    Pay-->>Users: 返回客户端调起支付所需参数

    Note over Users, 3rd: 3. 用户支付操作
    Users->>3rd: 在应用内/小程序唤起收银台支付
    3rd-->>Users: 支付成功界面提示

    Note over 3rd, Pay: 4. 支付结果回调 (Webhook)
    3rd->>Pay: 异步通知支付结果 (POST /pay/notify/order/{channelId})
    Pay->>Pay: 校验签名 & 校验金额
    Pay->>DB: 更新支付订单状态 -> 支付成功
    Pay->>DB: 插入支付通知任务 (PayNotifyTask, 状态: 待通知)
    Pay-->>3rd: 返回 SUCCESS (ACK)

    Note over Pay, Trade: 5. 内部状态同步 (解耦核心)
    loop 异步任务 / 消息重试
        Pay->>DB: 扫描待执行的通知任务
        Pay->>Trade: HTTP POST请求调用业务回调地址 <br/>(/app-api/trade/order/update-paid)
        Note right of Pay: 此时 Pay 模块充当 Client <br/>调用 Trade 模块提供的 API
        Trade->>Trade: 执行业务逻辑 (TradeOrderUpdateService)
        Trade->>DB: 更新交易订单状态 -> 已支付
        Trade-->>Pay: 返回 "success" (HTTP 200)
        Pay->>DB: 更新通知任务状态 -> 通知成功
    end
```

## 关键说明
1.  **物理隔离**: `Trade` (交易) 模块和 `Pay` (支付) 模块在代码逻辑上是完全解耦的。`Pay` 模块通过 HTTP 回调（步骤 5）通知 `Trade` 模块，而不是直接去操作交易模块的数据库表。
2.  **可靠性保证**: 如果步骤 5 中的 HTTP 回调因为网络或业务系统繁忙而失败，`PayNotifyTask` 会在数据库中保持“待通知”状态，并按照指数退避策略（如 15s, 30s, 3m...）自动重试，确保交易状态最终一定会被同步。
3.  **核心接口映射**:
    - **下单**: `TradeOrderUpdateService.CreateOrder`
    - **提交支付**: `PayOrderService.SubmitOrder`
    - **支付回调处理**: `PayOrderService.NotifyOrder`
    - **业务状态更新**: `AppTradeOrderHandler.UpdateOrderPaid`
