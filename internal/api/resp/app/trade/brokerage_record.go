package trade

import "time"

type AppBrokerageRecordRespVO struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	BizType     int       `json:"bizType"`
	BizID       string    `json:"bizId"`
	Price       int       `json:"price"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	Total       int       `json:"total"` // TotalPrice in model
	CreatedAt   time.Time `json:"createTime"`
	StatusName  string    `json:"statusName"`
}

type AppBrokerageProductPriceRespVO struct {
	BrokerageEnabled bool `json:"brokerageEnabled"` // 是否开启分销
	BrokeragePrice   int  `json:"brokeragePrice"`   // 分销金额
}
