package trade

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
)

// TradePriceCalculateReqBO 价格计算请求业务对象
type TradePriceCalculateReqBO struct {
	UserID                int64                       `json:"userId"`                // 用户ID
	CouponID              *int64                      `json:"couponId"`              // 优惠券ID
	PointStatus           bool                        `json:"pointStatus"`           // 是否使用积分
	DeliveryType          int                         `json:"deliveryType"`          // 配送方式
	AddressID             *int64                      `json:"addressId"`             // 收货地址ID
	PickUpStoreID         *int64                      `json:"pickUpStoreId"`         // 自提门店ID
	SeckillActivityId     int64                       `json:"seckillActivityId"`     // 秒杀活动ID
	CombinationActivityId int64                       `json:"combinationActivityId"` // 拼团活动ID
	CombinationHeadId     int64                       `json:"combinationHeadId"`     // 拼团团长ID
	BargainRecordId       int64                       `json:"bargainRecordId"`       // 砍价记录ID
	PointActivityId       int64                       `json:"pointActivityId"`       // 积分活动ID
	CartIDs               []int64                     `json:"cartIds"`               // 购物车ID数组
	Items                 []TradePriceCalculateItemBO `json:"items"`                 // 商品项数组
}

// TradePriceCalculateItemBO 价格计算商品项业务对象
type TradePriceCalculateItemBO struct {
	SkuID    int64 `json:"skuId"`    // 商品SKU ID
	Count    int   `json:"count"`    // 商品数量
	CartID   int64 `json:"cartId"`   // 购物车ID
	Selected bool  `json:"selected"` // 是否选中
}

// TradePriceCalculateRespBO 价格计算响应业务对象
type TradePriceCalculateRespBO struct {
	Type       int                              `json:"type"`       // 订单类型
	Price      TradePriceCalculatePriceBO       `json:"price"`      // 价格信息
	Items      []TradePriceCalculateItemRespBO  `json:"items"`      // 商品项数组
	CouponID   int64                            `json:"couponId"`   // 使用的优惠券ID
	TotalPoint int                              `json:"totalPoint"` // 用户总积分
	UsePoint   int                              `json:"usePoint"`   // 使用的积分
	GivePoint  int                              `json:"givePoint"`  // 赠送的积分
	Success    bool                             `json:"success"`    // 计算是否成功
	Coupons    []TradePriceCalculateCouponBO    `json:"coupons"`    // 可用优惠券数组
	Promotions []TradePriceCalculatePromotionBO `json:"promotions"` // 营销活动数组
}

// TradePriceCalculatePromotionBO 促销活动业务对象
type TradePriceCalculatePromotionBO struct {
	ID            int64                                `json:"id"`            // 活动ID
	Name          string                               `json:"name"`          // 活动名称
	Type          int                                  `json:"type"`          // 活动类型
	TotalPrice    int                                  `json:"totalPrice"`    // 总价格
	DiscountPrice int                                  `json:"discountPrice"` // 折扣金额
	Items         []TradePriceCalculatePromotionItemBO `json:"items"`         // 商品项明细
	Match         bool                                 `json:"match"`         // 是否匹配
	Description   string                               `json:"description"`   // 活动描述
}

// TradePriceCalculatePromotionItemBO 促销活动商品项业务对象
type TradePriceCalculatePromotionItemBO struct {
	SkuID         int64 `json:"skuId"`         // 商品SKU ID
	TotalPrice    int   `json:"totalPrice"`    // 总价格
	DiscountPrice int   `json:"discountPrice"` // 折扣金额
	PayPrice      int   `json:"payPrice"`      // 应付金额
}

// TradePriceCalculateCouponBO 优惠券业务对象
type TradePriceCalculateCouponBO struct {
	ID                 int64  `json:"id"`                 // 优惠券ID
	Name               string `json:"name"`               // 优惠券名称
	UsePrice           int    `json:"usePrice"`           // 使用门槛价格
	ValidStartTime     string `json:"validStartTime"`     // 有效开始时间
	ValidEndTime       string `json:"validEndTime"`       // 有效结束时间
	DiscountType       int    `json:"discountType"`       // 折扣类型
	DiscountPercent    int    `json:"discountPercent"`    // 折扣百分比
	DiscountPrice      int    `json:"discountPrice"`      // 折扣金额
	DiscountLimitPrice int    `json:"discountLimitPrice"` // 折扣限制价格
	Match              bool   `json:"match"`              // 是否匹配
	MismatchReason     string `json:"mismatchReason"`     // 不匹配原因
}

// TradePriceCalculatePriceBO 价格信息业务对象
type TradePriceCalculatePriceBO struct {
	TotalPrice    int `json:"totalPrice"`    // 总价格
	DiscountPrice int `json:"discountPrice"` // 折扣金额
	DeliveryPrice int `json:"deliveryPrice"` // 运费
	CouponPrice   int `json:"couponPrice"`   // 优惠券折扣
	PointPrice    int `json:"pointPrice"`    // 积分抵扣
	VipPrice      int `json:"vipPrice"`      // VIP折扣
	PayPrice      int `json:"payPrice"`      // 应付金额
}

// TradePriceCalculateItemRespBO 价格计算商品项响应业务对象
type TradePriceCalculateItemRespBO struct {
	SpuID         int64                         `json:"spuId"`         // 商品SPU ID
	SkuID         int64                         `json:"skuId"`         // 商品SKU ID
	Count         int                           `json:"count"`         // 商品数量
	CartID        int64                         `json:"cartId"`        // 购物车ID
	Selected      bool                          `json:"selected"`      // 是否选中
	Price         int                           `json:"price"`         // 商品单价
	DiscountPrice int                           `json:"discountPrice"` // 折扣金额
	DeliveryPrice int                           `json:"deliveryPrice"` // 运费
	CouponPrice   int                           `json:"couponPrice"`   // 优惠券折扣
	PointPrice    int                           `json:"pointPrice"`    // 积分抵扣
	UsePoint      int                           `json:"usePoint"`      // 使用的积分数量
	VipPrice      int                           `json:"vipPrice"`      // VIP折扣
	PayPrice      int                           `json:"payPrice"`      // 应付金额
	SpuName       string                        `json:"spuName"`       // 商品名称
	PicURL        string                        `json:"picUrl"`        // 商品图片
	CategoryID    int64                         `json:"categoryId"`    // 分类ID
	DeliveryTypes []int                         `json:"deliveryTypes"` // 配送方式
	GivePoint     int                           `json:"givePoint"`     // 赠送积分
	Properties    []resp.ProductSkuPropertyResp `json:"properties"`    // 商品属性
}
