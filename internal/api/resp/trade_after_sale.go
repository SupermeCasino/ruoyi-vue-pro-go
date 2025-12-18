package resp

import "time"

type AppAfterSaleResp struct {
	ID               int64                            `json:"id"`
	No               string                           `json:"no"`
	Status           int                              `json:"status"`
	Way              int                              `json:"way"`
	Type             int                              `json:"type"`
	ApplyReason      string                           `json:"applyReason"`
	ApplyDescription string                           `json:"applyDescription"`
	ApplyPicURLs     []string                         `json:"applyPicUrls"`
	CreateTime       time.Time                        `json:"createTime"`
	UpdateTime       time.Time                        `json:"updateTime"`
	OrderID          int64                            `json:"orderId"`
	OrderNo          string                           `json:"orderNo"`
	OrderItemID      int64                            `json:"orderItemId"`
	SpuID            int64                            `json:"spuId"`
	SpuName          string                           `json:"spuName"`
	SkuID            int64                            `json:"skuId"`
	Properties       []ProductPropertyValueDetailResp `json:"properties"`
	PicURL           string                           `json:"picUrl"`
	Count            int                              `json:"count"`
	AuditReason      string                           `json:"auditReason"`
	RefundPrice      int                              `json:"refundPrice"`
	RefundTime       *time.Time                       `json:"refundTime"`
	LogisticsID      int64                            `json:"logisticsId"`
	LogisticsNo      string                           `json:"logisticsNo"`
	DeliveryTime     *time.Time                       `json:"deliveryTime"`
	ReceiveTime      *time.Time                       `json:"receiveTime"`
	ReceiveReason    string                           `json:"receiveReason"`
	AuditTime        *time.Time                       `json:"auditTime"`
}

type AppAfterSaleLogResp struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createTime"`
}

// TradeAfterSaleDetailResp 售后订单详情响应 (Admin) - 严格对齐 Java AfterSaleDetailRespVO
type TradeAfterSaleDetailResp struct {
	// 基础字段
	ID               int64      `json:"id"`
	No               string     `json:"no"`
	Status           int        `json:"status"`
	Type             int        `json:"type"`
	Way              int        `json:"way"`
	UserID           int64      `json:"userId"`
	ApplyReason      string     `json:"applyReason"`
	ApplyDescription string     `json:"applyDescription"`
	ApplyPicURLs     []string   `json:"applyPicUrls"`
	OrderID          int64      `json:"orderId"`
	OrderNo          string     `json:"orderNo"`
	OrderItemID      int64      `json:"orderItemId"`
	SpuID            int64      `json:"spuId"`
	SpuName          string     `json:"spuName"`
	SkuID            int64      `json:"skuId"`
	PicURL           string     `json:"picUrl"`
	Count            int        `json:"count"`
	RefundPrice      int        `json:"refundPrice"`
	AuditTime        *time.Time `json:"auditTime"`
	AuditUserID      int64      `json:"auditUserId"` // 新增
	AuditReason      string     `json:"auditReason"`
	PayRefundID      int64      `json:"payRefundId"` // 新增
	RefundTime       *time.Time `json:"refundTime"`
	LogisticsID      int64      `json:"logisticsId"`   // 新增
	LogisticsNo      string     `json:"logisticsNo"`   // 新增
	DeliveryTime     *time.Time `json:"deliveryTime"`  // 新增
	ReceiveTime      *time.Time `json:"receiveTime"`   // 新增
	ReceiveReason    string     `json:"receiveReason"` // 新增
	CreateTime       time.Time  `json:"createTime"`
	// 嵌套对象
	Order     *TradeOrderBase     `json:"order"`     // 新增
	OrderItem *AfterSaleOrderItem `json:"orderItem"` // 新增
	User      *MemberUserResp     `json:"user"`      // 新增
	Logs      []AfterSaleLogResp  `json:"logs"`      // 新增
}

// AfterSaleOrderItem 售后订单项
type AfterSaleOrderItem struct {
	TradeOrderItemBase
	Properties []ProductPropertyValueDetailResp `json:"properties"`
}

// AfterSaleLogResp 售后日志
type AfterSaleLogResp struct {
	ID           int64     `json:"id"`
	AfterSaleID  int64     `json:"afterSaleId"`
	BeforeStatus int       `json:"beforeStatus"`
	AfterStatus  int       `json:"afterStatus"`
	OperateType  int       `json:"operateType"`
	UserType     int       `json:"userType"`
	UserID       int64     `json:"userId"`
	Content      string    `json:"content"`
	CreateTime   time.Time `json:"createTime"`
}

// ProductPropertyValueDetailResp 商品属性值详情
type ProductPropertyValueDetailResp struct {
	PropertyID   int64  `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	ValueID      int64  `json:"valueId"`
	ValueName    string `json:"valueName"`
}

// AfterSalePageItemResp 售后分页项
type AfterSalePageItemResp struct {
	ID          int64           `json:"id"`
	No          string          `json:"no"`
	Status      int             `json:"status"`
	Type        int             `json:"type"`
	Way         int             `json:"way"`
	UserID      int64           `json:"userId"`
	ApplyReason string          `json:"applyReason"`
	SpuName     string          `json:"spuName"`
	PicURL      string          `json:"picUrl"`
	Count       int             `json:"count"`
	RefundPrice int             `json:"refundPrice"`
	CreateTime  time.Time       `json:"createTime"`
	User        *MemberUserResp `json:"user"`
}
