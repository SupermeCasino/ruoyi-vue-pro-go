# Java 与 Go 订单金额计算对比分析

## 核心公式对齐

### 订单项 PayPrice 计算

**Java** (TradePriceCalculatorHelper.java:138-145)
```java
orderItem.setPayPrice(
    orderItem.getPrice() * orderItem.getCount()
    - orderItem.getDiscountPrice()
    + orderItem.getDeliveryPrice()
    - orderItem.getCouponPrice()
    - orderItem.getPointPrice()
    - orderItem.getVipPrice()
);
```

**Go** (price.go:238)
```go
PayPrice: itemPayPrice  // 初始值 = price * count
// 后续步骤中逐步应用各项折扣
```

**对齐状态**：✅ 最终结果一致，但处理方式不同
- Java：在每个计算器中实时更新
- Go：在最后统一计算

---

### 订单总价 PayPrice 计算

**Java** (TradePriceCalculatorHelper.java:104-122)
```java
// 先重置
price.setTotalPrice(0).setDiscountPrice(0).setDeliveryPrice(0)
        .setCouponPrice(0).setPointPrice(0).setVipPrice(0).setPayPrice(0);
// 再合计 item
result.getItems().forEach(item -> {
    if (!item.getSelected()) return;
    price.setTotalPrice(price.getTotalPrice() + item.getPrice() * item.getCount());
    price.setDiscountPrice(price.getDiscountPrice() + item.getDiscountPrice());
    price.setDeliveryPrice(price.getDeliveryPrice() + item.getDeliveryPrice());
    price.setCouponPrice(price.getCouponPrice() + item.getCouponPrice());
    price.setPointPrice(price.setPointPrice() + item.getPointPrice());
    price.setVipPrice(price.getVipPrice() + item.getVipPrice());
    price.setPayPrice(price.getPayPrice() + item.getPayPrice());
});
```

**Go** (price.go:285-295)
```go
respBO.Price.TotalPrice = totalPrice
respBO.Price.DiscountPrice = activityDiscount + seckillTotalDiscount
respBO.Price.VipPrice = totalVipPrice

payPrice := respBO.Price.TotalPrice - respBO.Price.DiscountPrice - respBO.Price.VipPrice
if payPrice < 0 {
    payPrice = 0
}
respBO.Price.PayPrice = payPrice
```

**对齐状态**：✅ 逻辑等价
- Java：逐项累加
- Go：直接计算

---

## 折扣计算流程对比

### 1. VIP 折扣（订单10）

**Java** (TradeDiscountActivityPriceCalculator.java:66-94)
```java
// 2.1 计算限时折扣的优惠金额
Integer discountPrice = calculateActivityPrice(discountProduct, orderItem);
// 2.2 计算 VIP 优惠金额
Integer vipPrice = calculateVipPrice(level, orderItem);
if (discountPrice <= 0 && vipPrice <= 0) {
    return;
}

// 3. 选择优惠金额多的
if (discountPrice > vipPrice) {
    // 使用限时折扣
    orderItem.setDiscountPrice(orderItem.getDiscountPrice() + discountPrice);
} else {
    // 使用 VIP 折扣
    orderItem.setVipPrice(vipPrice);
}

// 4. 分摊优惠
TradePriceCalculatorHelper.recountPayPrice(orderItem);
TradePriceCalculatorHelper.recountAllPrice(result);
```

**Go** (price.go:224-228)
```go
if levelDiscountPercent < 100 {
    vipTotal := int(int64(itemPrice) * int64(item.Count) * int64(levelDiscountPercent) / 100)
    itemVipSavings = itemPrice*item.Count - vipTotal
    itemPayPrice = vipTotal  // ✅ 已修复
}
```

**对齐状态**：✅ 已修复
- 修复前：VIP 折扣未应用到 PayPrice
- 修复后：正确更新 itemPayPrice

---

### 2. 秒杀折扣（订单8）

**Java** (TradeOrderTypeEnum)
```java
if (respBO.Type == TradeOrderTypeEnum.SECKILL.getType()) {
    // 秒杀逻辑
}
```

**Go** (price.go:213-220)
```go
if respBO.Type == 1 { // 秒杀
    _, seckillProd, err := s.seckillActivitySvc.ValidateJoinSeckill(ctx, req.SeckillActivityId, sku.ID, item.Count)
    if err != nil {
        return nil, err
    }
    seckillTotal := seckillProd.SeckillPrice * item.Count
    seckillDiscount = itemPayPrice - seckillTotal
    itemPayPrice = seckillTotal
}
```

**对齐状态**：✅ 正确
- 秒杀和 VIP 互斥
- 都直接修改 itemPayPrice

---

### 3. 满减活动折扣（订单20）

**Java** (TradeRewardActivityPriceCalculator.java:72-94)
```java
// 2.1 计算可以优惠的金额
Integer newDiscountPrice = rule.getDiscountPrice();
// 2.2 计算分摊的优惠金额
List<Integer> divideDiscountPrices = TradePriceCalculatorHelper.dividePrice(orderItems, newDiscountPrice);

// 3.3 更新 SKU 优惠金额
for (int i = 0; i < orderItems.size(); i++) {
    TradePriceCalculateRespBO.OrderItem orderItem = orderItems.get(i);
    orderItem.setDiscountPrice(orderItem.getDiscountPrice() + divideDiscountPrices.get(i));
    TradePriceCalculatorHelper.recountPayPrice(orderItem);
}
TradePriceCalculatorHelper.recountAllPrice(result);
```

**Go** (price.go:259-276)
```go
activityDiscount, _, err = s.rewardActivitySvc.CalculateRewardActivity(ctx, matchItems)
if err != nil {
    return nil, err
}
// ⚠️ 缺少分摊逻辑
respBO.Price.DiscountPrice = activityDiscount + seckillTotalDiscount
```

**对齐状态**：❌ 缺失分摊逻辑
- Java：按项分摊折扣
- Go：只计算总折扣，未分摊到各项

**影响**：
- 订单项的 `DiscountPrice` 为 0（应该被分摊）
- 订单项的 `PayPrice` 未正确减少

---

### 4. 优惠券折扣（订单30）

**Java** (TradeCouponPriceCalculator.java:70-91)
```java
// 3.1 计算可以优惠的金额
Integer totalPayPrice = TradePriceCalculatorHelper.calculateTotalPayPrice(orderItems);
Integer couponPrice = getCouponPrice(coupon, totalPayPrice);
// 3.2 计算分摊的优惠金额
List<Integer> divideCouponPrices = TradePriceCalculatorHelper.dividePrice(orderItems, couponPrice);

// 4.3 更新 SKU 优惠金额
for (int i = 0; i < orderItems.size(); i++) {
    TradePriceCalculateRespBO.OrderItem orderItem = orderItems.get(i);
    orderItem.setCouponPrice(divideCouponPrices.get(i));
    TradePriceCalculatorHelper.recountPayPrice(orderItem);
}
TradePriceCalculatorHelper.recountAllPrice(result);
```

**Go** (price.go:347-357)
```go
couponPrice, err := s.couponSvc.CalculateCoupon(ctx, req.UserID, *req.CouponID, int64(respBO.Price.PayPrice), spuIDs, categoryIDs)
if err != nil {
    return nil, err
}
respBO.CouponID = *req.CouponID
respBO.Price.CouponPrice = int(couponPrice)
respBO.Price.PayPrice -= int(couponPrice)
```

**对齐状态**：⚠️ 部分缺失
- Java：分摊到各项
- Go：只更新总价，未分摊到各项

---

### 5. 积分抵扣（订单40）

**Java** (TradePointUsePriceCalculator)
```java
if (payPrice <= pointPrice) throw exception  // 禁止 0 元购
```

**Go** (price.go:385)
```go
if pointTotalValue >= respBO.Price.PayPrice {
    return nil, core.NewBizError(1004003005, "支付金额不能小于等于 0")
}
```

**对齐状态**：✅ 正确
- 条件等价：`payPrice <= pointPrice` 等同于 `pointTotalValue >= respBO.Price.PayPrice`

---

### 6. 运费计算（订单50）

**Java** (DeliveryExpressTemplateService)
```java
// 基于 item.getPayPrice() 计算运费
// PayPrice 已经包含了所有折扣
```

**Go** (price.go:385)
```go
templatePriceMap[spu.DeliveryTemplateID] += item.PayPrice
```

**对齐状态**：✅ 正确
- 都使用 PayPrice 作为基数

---

## 计算顺序对比

### Java 计算顺序
1. 初始化订单项（原价）
2. 计算 VIP/限时折扣（选择较大的）
3. 计算满减活动折扣
4. 计算优惠券折扣
5. 计算积分抵扣
6. 计算运费
7. 每步都调用 `recountPayPrice()` 和 `recountAllPrice()`

### Go 计算顺序
1. 初始化订单项（原价）
2. 计算 VIP 折扣（✅ 已修复）
3. 计算秒杀折扣
4. 计算满减活动折扣（⚠️ 未分摊）
5. 计算优惠券折扣
6. 计算积分抵扣
7. 计算运费
8. 最后统一计算总价

**差异**：
- Java：实时更新，每步都重新计算
- Go：延迟计算，最后统一处理

---

## 数据一致性检查

### 订单项字段

| 字段 | Java | Go | 状态 |
|------|------|----|----|
| Price | ✅ 原价 | ✅ 原价 | ✅ |
| Count | ✅ 数量 | ✅ 数量 | ✅ |
| DiscountPrice | ✅ 折扣 | ⚠️ 仅秒杀 | ⚠️ |
| VipPrice | ✅ VIP折扣 | ✅ VIP折扣 | ✅ |
| CouponPrice | ✅ 优惠券 | ⚠️ 未分摊 | ⚠️ |
| PointPrice | ✅ 积分 | ⚠️ 未分摊 | ⚠️ |
| DeliveryPrice | ✅ 运费 | ⚠️ 未分摊 | ⚠️ |
| PayPrice | ✅ 支付价 | ✅ 支付价 | ✅ |

---

## 关键发现

### ✅ 已修复
1. VIP 折扣正确应用到 PayPrice

### ⚠️ 需要改进
1. 满减活动折扣未分摊到各项
2. 优惠券折扣未分摊到各项
3. 积分抵扣未分摊到各项
4. 运费未分摊到各项

### 🔴 高风险
- 订单项的 `DiscountPrice`、`CouponPrice`、`PointPrice`、`DeliveryPrice` 可能为 0
- 创建订单时，这些字段会被保存为 0，导致数据不一致
- 后续查询订单时，无法准确还原各项折扣

---

## 建议修复优先级

1. **P0 - 立即修复**：VIP 折扣（✅ 已完成）
2. **P1 - 高优先级**：满减活动折扣分摊
3. **P2 - 中优先级**：优惠券和积分分摊
4. **P3 - 低优先级**：运费分摊
