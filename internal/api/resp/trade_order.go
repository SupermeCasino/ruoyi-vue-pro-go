package resp

import "time"

// AppTradeOrderCreateResp 交易订单创建 Response
type AppTradeOrderCreateResp struct {
	ID         int64 `json:"id"`
	PayOrderID int64 `json:"payOrderId"`
}

// AppTradeOrderSettlementResp 交易订单结算信息 Response
type AppTradeOrderSettlementResp struct {
	Type       int                              `json:"type"`
	Items      []AppTradeOrderSettlementItem    `json:"items"`
	Coupons    []AppTradeOrderSettlementCoupon  `json:"coupons"`
	Price      AppTradeOrderSettlementPrice     `json:"price"`
	Address    *AppTradeOrderSettlementAddress  `json:"address"`
	UsePoint   int                              `json:"usePoint"`
	TotalPoint int                              `json:"totalPoint"`
	Promotions []AppTradeOrderSettlementPromotion `json:"promotions"`
}

type AppTradeOrderSettlementItem struct {
	CategoryID int64                    `json:"categoryId"`
	SpuID      int64                    `json:"spuId"`
	SpuName    string                   `json:"spuName"`
	SkuID      int64                    `json:"skuId"`
	Price      int                      `json:"price"`
	PicURL     string                   `json:"picUrl"`
	Properties []ProductSkuPropertyResp `json:"properties"`
	CartID     int64                    `json:"cartId"`
	Count      int                      `json:"count"`
}

type AppTradeOrderSettlementCoupon struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	UsePrice           int        `json:"usePrice"`
	ValidStartTime     *time.Time `json:"validStartTime"`
	ValidEndTime       *time.Time `json:"validEndTime"`
	DiscountType       int        `json:"discountType"`
	DiscountPercent    int        `json:"discountPercent"`
	DiscountPrice      int        `json:"discountPrice"`
	DiscountLimitPrice int        `json:"discountLimitPrice"`
	Match              bool       `json:"match"`
	MismatchReason     string     `json:"mismatchReason"`
}

type AppTradeOrderSettlementPrice struct {
	TotalPrice    int `json:"totalPrice"`
	DiscountPrice int `json:"discountPrice"`
	DeliveryPrice int `json:"deliveryPrice"`
	CouponPrice   int `json:"couponPrice"`
	PointPrice    int `json:"pointPrice"`
	VipPrice      int `json:"vipPrice"`
	PayPrice      int `json:"payPrice"`
}

type AppTradeOrderSettlementAddress struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Mobile        string `json:"mobile"`
	AreaID        int64  `json:"areaId"`
	AreaName      string `json:"areaName"`
	DetailAddress string `json:"detailAddress"`
	DefaultStatus bool   `json:"defaultStatus"`
}

// AppTradeOrderSettlementPromotion 交易订单结算 - 促销活动信息
type AppTradeOrderSettlementPromotion struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Type          int    `json:"type"`          // 活动类型
	DiscountPrice int    `json:"discountPrice"` // 折扣金额
}

// AppTradeOrderDetailResp 交易订单详情 Response
type AppTradeOrderDetailResp struct {
	ID                    int64                   `json:"id"`
	No                    string                  `json:"no"`
	Type                  int                     `json:"type"`
	CreateTime            time.Time               `json:"createTime"`
	UserRemark            string                  `json:"userRemark"`
	Status                int                     `json:"status"`
	ProductCount          int                     `json:"productCount"`
	FinishTime            *time.Time              `json:"finishTime"`
	CancelTime            *time.Time              `json:"cancelTime"`
	CommentStatus         bool                    `json:"commentStatus"`
	PayStatus             bool                    `json:"payStatus"`
	PayOrderID            int64                   `json:"payOrderId"`
	PayTime               *time.Time              `json:"payTime"`
	PayExpireTime         *time.Time              `json:"payExpireTime"`
	PayChannelCode        string                  `json:"payChannelCode"`
	PayChannelName        string                  `json:"payChannelName"`
	TotalPrice            int                     `json:"totalPrice"`
	DiscountPrice         int                     `json:"discountPrice"`
	DeliveryPrice         int                     `json:"deliveryPrice"`
	AdjustPrice           int                     `json:"adjustPrice"`
	PayPrice              int                     `json:"payPrice"`
	DeliveryType          int                     `json:"deliveryType"`
	LogisticsID           int64                   `json:"logisticsId"`
	LogisticsName         string                  `json:"logisticsName"`
	LogisticsNo           string                  `json:"logisticsNo"`
	DeliveryTime          *time.Time              `json:"deliveryTime"`
	ReceiveTime           *time.Time              `json:"receiveTime"`
	ReceiverName          string                  `json:"receiverName"`
	ReceiverMobile        string                  `json:"receiverMobile"`
	ReceiverAreaID        int                     `json:"receiverAreaId"`
	ReceiverAreaName      string                  `json:"receiverAreaName"`
	ReceiverDetailAddress string                  `json:"receiverDetailAddress"`
	PickUpStoreID         int64                   `json:"pickUpStoreId"`
	PickUpVerifyCode      string                  `json:"pickUpVerifyCode"`
	RefundStatus          int                     `json:"refundStatus"`
	RefundPrice           int                     `json:"refundPrice"`
	CouponID              int64                   `json:"couponId"`
	CouponPrice           int                     `json:"couponPrice"`
	PointPrice            int                     `json:"pointPrice"`
	VipPrice              int                     `json:"vipPrice"`
	CombinationRecordID   int64                   `json:"combinationRecordId"`
	Items                 []AppTradeOrderItemResp `json:"items"`
}

type AppTradeOrderItemResp struct {
	ID              int64                    `json:"id"`
	OrderID         int64                    `json:"orderId"`
	SpuID           int64                    `json:"spuId"`
	SpuName         string                   `json:"spuName"`
	SkuID           int64                    `json:"skuId"`
	Properties      []ProductSkuPropertyResp `json:"properties"`
	PicURL          string                   `json:"picUrl"`
	Count           int                      `json:"count"`
	CommentStatus   bool                     `json:"commentStatus"`
	Price           int                      `json:"price"`
	PayPrice        int                      `json:"payPrice"`
	AfterSaleID     int64                    `json:"afterSaleId"`
	AfterSaleStatus int                      `json:"afterSaleStatus"`
}

// AppTradeOrderPageItemResp 分页项 Response
type AppTradeOrderPageItemResp struct {
	ID                  int64                   `json:"id"`
	No                  string                  `json:"no"`
	Type                int                     `json:"type"`
	Status              int                     `json:"status"`
	ProductCount        int                     `json:"productCount"`
	CommentStatus       bool                    `json:"commentStatus"`
	CreateTime          time.Time               `json:"createTime"`
	PayOrderID          int64                   `json:"payOrderId"`
	PayPrice            int                     `json:"payPrice"`
	DeliveryType        int                     `json:"deliveryType"`
	Items               []AppTradeOrderItemResp `json:"items"`
	CombinationRecordID int64                   `json:"combinationRecordId"`
}

// Admin DTOs
type TradeOrderDetailResp struct {
	TradeOrderBase
	Items            []TradeOrderItemBase `json:"items"`
	User             interface{}          `json:"user"` // MemberUserRespVO
	BrokerageUser    interface{}          `json:"brokerageUser"`
	Logs             []TradeOrderLogResp  `json:"logs"`
	ReceiverAreaName string               `json:"receiverAreaName"`
}

type TradeOrderBase struct {
	ID                    int64      `json:"id"`
	No                    string     `json:"no"`
	CreateTime            time.Time  `json:"createTime"`
	Type                  int        `json:"type"`
	Terminal              int        `json:"terminal"`
	UserID                int64      `json:"userId"`
	UserIP                string     `json:"userIp"`
	UserRemark            string     `json:"userRemark"`
	Status                int        `json:"status"`
	ProductCount          int        `json:"productCount"`
	FinishTime            *time.Time `json:"finishTime"`
	CancelTime            *time.Time `json:"cancelTime"`
	CancelType            int        `json:"cancelType"`
	Remark                string     `json:"remark"`
	PayOrderID            int64      `json:"payOrderId"`
	PayStatus             bool       `json:"payStatus"`
	PayTime               *time.Time `json:"payTime"`
	PayChannelCode        string     `json:"payChannelCode"`
	TotalPrice            int        `json:"totalPrice"`
	DiscountPrice         int        `json:"discountPrice"`
	DeliveryPrice         int        `json:"deliveryPrice"`
	AdjustPrice           int        `json:"adjustPrice"`
	PayPrice              int        `json:"payPrice"`
	DeliveryType          int        `json:"deliveryType"`
	PickUpStoreID         int64      `json:"pickUpStoreId"`
	PickUpVerifyCode      string     `json:"pickUpVerifyCode"`
	DeliveryTemplateID    int64      `json:"deliveryTemplateId"`
	LogisticsID           int64      `json:"logisticsId"`
	LogisticsNo           string     `json:"logisticsNo"`
	DeliveryTime          *time.Time `json:"deliveryTime"`
	ReceiveTime           *time.Time `json:"receiveTime"`
	ReceiverName          string     `json:"receiverName"`
	ReceiverMobile        string     `json:"receiverMobile"`
	ReceiverAreaID        int        `json:"receiverAreaId"`
	ReceiverDetailAddress string     `json:"receiverDetailAddress"`
	AfterSaleStatus       int        `json:"afterSaleStatus"`
	RefundPrice           int        `json:"refundPrice"`
	CouponID              int64      `json:"couponId"`
	CouponPrice           int        `json:"couponPrice"`
	PointPrice            int        `json:"pointPrice"`
	VipPrice              int        `json:"vipPrice"`
	BrokerageUserID       int64      `json:"brokerageUserId"`
}

type TradeOrderItemBase struct {
	ID               int64                    `json:"id"`
	UserID           int64                    `json:"userId"`
	OrderID          int64                    `json:"orderId"`
	SpuID            int64                    `json:"spuId"`
	SpuName          string                   `json:"spuName"`
	SkuID            int64                    `json:"skuId"`
	PicURL           string                   `json:"picUrl"`
	Count            int                      `json:"count"`
	Price            int                      `json:"price"`
	DiscountPrice    int                      `json:"discountPrice"`
	DeliveryPrice    int                      `json:"deliveryPrice"`
	AdjustPrice      int                      `json:"adjustPrice"`
	PayPrice         int                      `json:"payPrice"`
	OrderPartPrice   int                      `json:"orderPartPrice"`   // 子订单分摊金额
	OrderDividePrice int                      `json:"orderDividePrice"` // 分摊后子订单实付金额
	CouponPrice      int                      `json:"couponPrice"`
	PointPrice       int                      `json:"pointPrice"`
	UsePoint         int                      `json:"usePoint"`
	GivePoint        int                      `json:"givePoint"`
	VipPrice         int                      `json:"vipPrice"`
	AfterSaleID      int64                    `json:"afterSaleId"`
	AfterSaleStatus  int                      `json:"afterSaleStatus"`
	Properties       []ProductSkuPropertyResp `json:"properties"`
}

type TradeOrderLogResp struct {
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createTime"`
	UserType   int       `json:"userType"`
}

type TradeOrderPageItemResp struct {
	TradeOrderBase
	Items            []TradeOrderItemBase `json:"items"`
	User             *MemberUserResp      `json:"user"`
	BrokerageUser    *MemberUserResp      `json:"brokerageUser"`
	ReceiverAreaName string               `json:"receiverAreaName"`
}

type TradeOrderSummaryResp struct {
	OrderCount     int64 `json:"orderCount"`
	OrderPayPrice  int64 `json:"orderPayPrice"`
	AfterSaleCount int64 `json:"afterSaleCount"`
	AfterSalePrice int64 `json:"afterSalePrice"`
}

// AppTradeProductSettlementResp 用户 App - 商品结算信息 Response VO
type AppTradeProductSettlementResp struct {
	SpuID          int64           `json:"spuId"`
	Skus           []Sku           `json:"skus"`
	RewardActivity *RewardActivity `json:"rewardActivity"`
}

type Sku struct {
	ID               int64  `json:"id"`
	PromotionPrice   int    `json:"promotionPrice"`
	PromotionType    int    `json:"promotionType"` // 对应 PromotionTypeEnum 枚举
	PromotionID      int64  `json:"promotionId"`
	PromotionEndTime string `json:"promotionEndTime"`
}

type RewardActivity struct {
	ID            int64                `json:"id"`
	ConditionType int                  `json:"conditionType"`
	Rules         []RewardActivityRule `json:"rules"`
}

type RewardActivityRule struct {
	Limit                    int           `json:"limit"`
	DiscountPrice            int           `json:"discountPrice"`
	FreeDelivery             bool          `json:"freeDelivery"`
	Point                    int           `json:"point"`
	GiveCouponTemplateCounts map[int64]int `json:"giveCouponTemplateCounts"`
}
