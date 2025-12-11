package resp

import "time"

type AppAfterSaleResp struct {
	ID               int64     `json:"id"`
	No               string    `json:"no"`
	Status           int       `json:"status"`
	Way              int       `json:"way"`
	Type             int       `json:"type"`
	ApplyReason      string    `json:"applyReason"`
	ApplyDescription string    `json:"applyDescription"`
	ApplyPicURLs     []string  `json:"applyPicUrls"`
	OrderNo          string    `json:"orderNo"`
	SpuName          string    `json:"spuName"`
	PicURL           string    `json:"picUrl"`
	Count            int       `json:"count"`
	RefundPrice      int       `json:"refundPrice"`
	AuditTime        time.Time `json:"auditTime"`
	AuditReason      string    `json:"auditReason"`
	CreateTime       time.Time `json:"createTime"`
}

type AppAfterSaleLogResp struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createTime"`
}

// TradeAfterSaleDetailResp 售后订单详情响应 (Admin)
type TradeAfterSaleDetailResp struct {
	ID                int64     `json:"id"`
	No                string    `json:"no"`
	Status            int       `json:"status"`
	Way               int       `json:"way"`
	Type              int       `json:"type"`
	UserID            int64     `json:"userId"`
	ApplyReason       string    `json:"applyReason"`
	ApplyDescription  string    `json:"applyDescription"`
	ApplyPicURLs      []string  `json:"applyPicUrls"`
	OrderID           int64     `json:"orderId"`
	OrderNo           string    `json:"orderNo"`
	OrderItemID       int64     `json:"orderItemId"`
	OrderPayPrice     int       `json:"orderPayPrice"`
	OrderItemPayPrice int       `json:"orderItemPayPrice"`
	SpuID             int64     `json:"spuId"`
	SpuName           string    `json:"spuName"`
	SkuID             int64     `json:"skuId"`
	PicURL            string    `json:"picUrl"`
	Count             int       `json:"count"`
	RefundPrice       int       `json:"refundPrice"`
	AuditTime         time.Time `json:"auditTime"`
	AuditReason       string    `json:"auditReason"`
	RefundTime        time.Time `json:"refundTime"`
	CreateTime        time.Time `json:"createTime"`
}
