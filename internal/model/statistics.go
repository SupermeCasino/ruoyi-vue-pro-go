package model

import "time"

// TradeStatistics 交易统计数据对象
type TradeStatistics struct {
	ID               int64     `gorm:"primaryKey;column:id" json:"id"`
	StatisticsTime   time.Time `gorm:"column:statistics_time;index" json:"statisticsTime"`
	OrderCount       int64     `gorm:"column:order_count" json:"orderCount"`              // 订单数
	OrderPayPrice    int64     `gorm:"column:order_pay_price" json:"orderPayPrice"`       // 订单支付金额
	OrderRefundPrice int64     `gorm:"column:order_refund_price" json:"orderRefundPrice"` // 订单退款金额
	AfterSaleCount   int64     `gorm:"column:after_sale_count" json:"afterSaleCount"`     // 售后数
	AfterSalePrice   int64     `gorm:"column:after_sale_price" json:"afterSalePrice"`     // 售后金额
	BrokeragePrice   int64     `gorm:"column:brokerage_price" json:"brokeragePrice"`      // 佣金金额
	CreateTime        time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdateTime        time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
}

// ProductStatistics 商品统计数据对象
type ProductStatistics struct {
	ID             int64     `gorm:"primaryKey;column:id" json:"id"`
	StatisticsTime time.Time `gorm:"column:statistics_time;index" json:"statisticsTime"`
	SpuID          int64     `gorm:"column:spu_id;index" json:"spuId"`
	BuyCount       int64     `gorm:"column:buy_count" json:"buyCount"`           // 购买数
	BuyPrice       int64     `gorm:"column:buy_price" json:"buyPrice"`           // 购买金额
	BrowseCount    int64     `gorm:"column:browse_count" json:"browseCount"`     // 浏览数
	FavoriteCount  int64     `gorm:"column:favorite_count" json:"favoriteCount"` // 收藏数
	CommentCount   int64     `gorm:"column:comment_count" json:"commentCount"`   // 评价数
	CreateTime      time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdateTime      time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
}

// MemberStatistics 会员统计数据对象
type MemberStatistics struct {
	ID              int64     `gorm:"primaryKey;column:id" json:"id"`
	StatisticsTime  time.Time `gorm:"column:statistics_time;index" json:"statisticsTime"`
	RegisterCount   int64     `gorm:"column:register_count" json:"registerCount"`      // 注册数
	VisitUserCount  int64     `gorm:"column:visit_user_count" json:"visitUserCount"`   // 访客数
	OrderUserCount  int64     `gorm:"column:order_user_count" json:"orderUserCount"`   // 下单用户数
	PayUserCount    int64     `gorm:"column:pay_user_count" json:"payUserCount"`       // 支付用户数
	TotalUserCount  int64     `gorm:"column:total_user_count" json:"totalUserCount"`   // 总用户数
	ActiveUserCount int64     `gorm:"column:active_user_count" json:"activeUserCount"` // 活跃用户数
	CreateTime       time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdateTime       time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
}

// PayStatistics 支付统计数据对象
type PayStatistics struct {
	ID             int64     `gorm:"primaryKey;column:id" json:"id"`
	StatisticsTime time.Time `gorm:"column:statistics_time;index" json:"statisticsTime"`
	RechargePrice  int64     `gorm:"column:recharge_price" json:"rechargePrice"` // 充值金额
	CreateTime      time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdateTime      time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
}
