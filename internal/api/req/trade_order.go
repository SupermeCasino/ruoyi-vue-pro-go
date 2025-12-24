package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/types"
)

// AppTradeOrderSettlementReq 交易订单结算请求
type AppTradeOrderSettlementReq struct {
	Items                 []AppTradeOrderSettlementItem `json:"items" form:"items" binding:"dive"`
	CouponID              *int64                        `json:"couponId" form:"couponId"`
	PointStatus           *bool                         `json:"pointStatus" form:"pointStatus" binding:"required"`
	DeliveryType          int                           `json:"deliveryType" form:"deliveryType" binding:"required"` // 1: 快递, 2: 自提
	AddressID             *int64                        `json:"addressId" form:"addressId"`
	PickUpStoreID         *int64                        `json:"pickUpStoreId" form:"pickUpStoreId"`
	ReceiverName          string                        `json:"receiverName" form:"receiverName"`
	ReceiverMobile        string                        `json:"receiverMobile" form:"receiverMobile"`
	SeckillActivityID     *int64                        `json:"seckillActivityId" form:"seckillActivityId"`
	CombinationActivityID *int64                        `json:"combinationActivityId" form:"combinationActivityId"`
	CombinationHeadID     *int64                        `json:"combinationHeadId" form:"combinationHeadId"`
	BargainRecordID       *int64                        `json:"bargainRecordId" form:"bargainRecordId"`
	PointActivityID       *int64                        `json:"pointActivityId" form:"pointActivityId"`
}

type AppTradeOrderSettlementItem struct {
	SkuID  int64 `json:"skuId" form:"skuId"`
	Count  int   `json:"count" form:"count"`
	CartID int64 `json:"cartId" form:"cartId"`
}

// AppTradeOrderCreateReq 交易订单创建请求
type AppTradeOrderCreateReq struct {
	AppTradeOrderSettlementReq
	Remark string `json:"remark"`
}

// AppTradeOrderPageReq 交易订单分页请求
type AppTradeOrderPageReq struct {
	pagination.PageParam
	Status        *int  `form:"status"`
	CommentStatus *bool `form:"commentStatus"`
}

// TradeOrderPageReq 管理后台 - 交易订单分页请求
type TradeOrderPageReq struct {
	pagination.PageParam
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
	CommentStatus    *bool    `form:"commentStatus"`
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

// AppTradeOrderSettlementQueryReq 订单结算查询请求 (用于GET请求)
type AppTradeOrderSettlementQueryReq struct {
	SkuIds                types.ListFromCSV[int64] `form:"skuIds"`  // SKU ID列表，逗号分隔
	Counts                types.ListFromCSV[int]   `form:"counts"`  // 数量列表，逗号分隔
	CartIds               types.ListFromCSV[int64] `form:"cartIds"` // 购物车ID列表，逗号分隔
	CouponID              *int64                   `form:"couponId"`
	PointStatus           *bool                    `form:"pointStatus" binding:"required"`
	DeliveryType          int                      `form:"deliveryType" binding:"required"` // 1: 快递, 2: 自提
	AddressID             *int64                   `form:"addressId"`
	PickUpStoreID         *int64                   `form:"pickUpStoreId"`
	ReceiverName          string                   `form:"receiverName"`
	ReceiverMobile        string                   `form:"receiverMobile"`
	SeckillActivityID     *int64                   `form:"seckillActivityId"`
	CombinationActivityID *int64                   `form:"combinationActivityId"`
	CombinationHeadID     *int64                   `form:"combinationHeadId"`
	BargainRecordID       *int64                   `form:"bargainRecordId"`
	PointActivityID       *int64                   `form:"pointActivityId"`
}

func (q *AppTradeOrderSettlementQueryReq) ToSettlementItems() []AppTradeOrderSettlementItem {
	if len(q.SkuIds) == 0 {
		return nil
	}
	items := make([]AppTradeOrderSettlementItem, len(q.SkuIds))
	for i, skuId := range q.SkuIds {
		items[i].SkuID = skuId
		if i < len(q.Counts) {
			items[i].Count = q.Counts[i]
		}
		if i < len(q.CartIds) {
			items[i].CartID = q.CartIds[i]
		}
	}
	return items
}

// AppTradeOrderUpdatePaidReq 更新订单为已支付请求
type AppTradeOrderUpdatePaidReq struct {
	ID int64 `json:"id" binding:"required"`
}

// AppTradeOrderReceiveReq 确认收货请求
type AppTradeOrderReceiveReq struct {
	ID int64 `json:"id" binding:"required"`
}
