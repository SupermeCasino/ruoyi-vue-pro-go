package req

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
)

// AppTradeOrderSettlementReq 交易订单结算请求
type AppTradeOrderSettlementReq struct {
	Items                 []AppTradeOrderSettlementItem `json:"items" binding:"required,dive"`
	CouponID              *int64                        `json:"couponId"`
	PointStatus           bool                          `json:"pointStatus" binding:"required"`
	DeliveryType          int                           `json:"deliveryType" binding:"required"` // 1: 快递, 2: 自提
	AddressID             *int64                        `json:"addressId"`
	PickUpStoreID         *int64                        `json:"pickUpStoreId"`
	ReceiverName          string                        `json:"receiverName"`
	ReceiverMobile        string                        `json:"receiverMobile"`
	SeckillActivityID     *int64                        `json:"seckillActivityId"`
	CombinationActivityID *int64                        `json:"combinationActivityId"`
	CombinationHeadID     *int64                        `json:"combinationHeadId"`
	BargainRecordID       *int64                        `json:"bargainRecordId"`
	PointActivityID       *int64                        `json:"pointActivityId"`
}

type AppTradeOrderSettlementItem struct {
	SkuID  int64 `json:"skuId"`
	Count  int   `json:"count"`
	CartID int64 `json:"cartId"`
}

// AppTradeOrderCreateReq 交易订单创建请求
type AppTradeOrderCreateReq struct {
	AppTradeOrderSettlementReq
	Remark string `json:"remark"`
}

// AppTradeOrderPageReq 交易订单分页请求
type AppTradeOrderPageReq struct {
	core.PageParam
	Status        *int  `form:"status"`
	CommentStatus *bool `form:"commentStatus"`
}

// TradeOrderPageReq 管理后台 - 交易订单分页请求
type TradeOrderPageReq struct {
	core.PageParam
	No               string   `form:"no"`
	UserID           *int64   `form:"userId"`
	UserNickname     string   `form:"userNickname"`
	UserMobile       string   `form:"userMobile"`
	DeliveryType     *int     `form:"deliveryType"`
	LogisticsID      *int64   `form:"logisticsId"`
	PickUpStoreIDs   []int64  `form:"pickUpStoreIds"`
	PickUpVerifyCode string   `form:"pickUpVerifyCode"`
	Type             *int     `form:"type"`
	Status           *int     `form:"status"`
	PayChannelCode   string   `form:"payChannelCode"`
	CreateTime       []string `form:"createTime[]"` // time range
	Terminal         *int     `form:"terminal"`
}

// TradeOrderDeliveryReq 订单发货请求
type TradeOrderDeliveryReq struct {
	ID          int64  `json:"id" binding:"required"`
	LogisticsID int64  `json:"logisticsId" binding:"required"`
	LogisticsNo string `json:"logisticsNo" binding:"required"`
}

// TradeOrderUpdateAddressReq 更新订单地址请求
type TradeOrderUpdateAddressReq struct {
	ID                    int64  `json:"id" binding:"required"`
	ReceiverName          string `json:"receiverName" binding:"required"`
	ReceiverMobile        string `json:"receiverMobile" binding:"required"`
	ReceiverAreaID        int    `json:"receiverAreaId" binding:"required"`
	ReceiverDetailAddress string `json:"receiverDetailAddress" binding:"required"`
}

// TradeOrderUpdatePriceReq 更新订单价格请求
type TradeOrderUpdatePriceReq struct {
	ID          int64 `json:"id" binding:"required"`
	AdjustPrice int   `json:"adjustPrice" binding:"required"` // 单位：分
}

// TradeOrderRemarkReq 订单备注请求
type TradeOrderRemarkReq struct {
	ID     int64  `json:"id" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}

// AppTradeOrderItemCommentCreateReq 用户 App - 商品评价创建请求
type AppTradeOrderItemCommentCreateReq struct {
	Anonymous         bool     `json:"anonymous"`
	OrderItemID       int64    `json:"orderItemId" binding:"required"`
	DescriptionScores int      `json:"descriptionScores" binding:"required,min=1,max=5"`
	BenefitScores     int      `json:"benefitScores" binding:"required,min=1,max=5"`
	Content           string   `json:"content" binding:"required"`
	PicUrls           []string `json:"picUrls" binding:"max=9"`
}

// AppTradeOrderSettlementProductReq 获得订单结算的商品信息请求
type AppTradeOrderSettlementProductReq struct {
	SkuID int64 `form:"skuId" binding:"required"`
	Count int   `form:"count" binding:"required,min=1"`
}

// AppTradeOrderUpdatePaidReq 更新订单为已支付请求
type AppTradeOrderUpdatePaidReq struct {
	ID int64 `json:"id" binding:"required"`
}

// AppTradeOrderReceiveReq 确认收货请求
type AppTradeOrderReceiveReq struct {
	ID int64 `json:"id" binding:"required"`
}
