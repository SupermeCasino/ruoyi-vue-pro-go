package resp

import (
	"time"
)

// CombinationProductRespVO 拼团商品 Response VO
type CombinationProductRespVO struct {
	SpuID             int64     `json:"spuId"`
	SkuID             int64     `json:"skuId"`
	CombinationPrice  int       `json:"combinationPrice"`
	ActivityStatus    int       `json:"activityStatus"`
	ActivityStartTime time.Time `json:"activityStartTime"`
	ActivityEndTime   time.Time `json:"activityEndTime"`
}

// CombinationActivityRespVO 拼团活动 Response VO (Admin)
type CombinationActivityRespVO struct {
	ID               int64                      `json:"id"`
	Name             string                     `json:"name"`
	SpuID            int64                      `json:"spuId"`
	TotalLimitCount  int                        `json:"totalLimitCount"`
	SingleLimitCount int                        `json:"singleLimitCount"`
	StartTime        time.Time                  `json:"startTime"`
	EndTime          time.Time                  `json:"endTime"`
	UserSize         int                        `json:"userSize"`
	VirtualGroup     bool                       `json:"virtualGroup"`
	LimitDuration    int                        `json:"limitDuration"`
	Status           int                        `json:"status"`
	Products         []CombinationProductRespVO `json:"products"`
	CreateTime       time.Time                  `json:"createTime"`
	// SPU 字段
	SpuName          string `json:"spuName"`          // 商品名称
	PicUrl           string `json:"picUrl"`           // 商品主图
	MarketPrice      int    `json:"marketPrice"`      // 商品市场价
	CombinationPrice int    `json:"combinationPrice"` // 拼团最低价
}

// CombinationActivityPageItemRespVO 拼团活动分页项 Response VO (Admin page)
type CombinationActivityPageItemRespVO struct {
	CombinationActivityRespVO
	// 统计字段
	GroupCount        int `json:"groupCount"`        // 开团组数
	GroupSuccessCount int `json:"groupSuccessCount"` // 成团组数
	RecordCount       int `json:"recordCount"`       // 购买次数
}

// AppCombinationActivityRespVO (Simple list item)
type AppCombinationActivityRespVO struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	UserSize         int    `json:"userSize"`
	SpuID            int64  `json:"spuId"`
	SpuName          string `json:"spuName"`
	PicUrl           string `json:"picUrl"`
	MarketPrice      int    `json:"marketPrice"`
	CombinationPrice int    `json:"combinationPrice"`
}

// AppCombinationActivityDetailRespVO 拼团活动详情 (App)
type AppCombinationActivityDetailRespVO struct {
	ID               int64                                 `json:"id"`
	Name             string                                `json:"name"`
	Status           int                                   `json:"status"`
	StartTime        *time.Time                            `json:"startTime"`
	EndTime          *time.Time                            `json:"endTime"`
	UserSize         int                                   `json:"userSize"`
	SuccessCount     int                                   `json:"successCount"` // 成功的拼团数量
	SpuID            int64                                 `json:"spuId"`
	TotalLimitCount  int                                   `json:"totalLimitCount"`
	SingleLimitCount int                                   `json:"singleLimitCount"`
	Products         []AppCombinationActivityDetailProduct `json:"products"`
}

// AppCombinationActivityDetailProduct 拼团活动详情商品
type AppCombinationActivityDetailProduct struct {
	SkuID            int64 `json:"skuId"`
	CombinationPrice int   `json:"combinationPrice"`
}

// AppCombinationRecordRespVO
type AppCombinationRecordRespVO struct {
	ID               int64     `json:"id"`
	ActivityID       int64     `json:"activityId"`
	Nickname         string    `json:"nickname"`
	Avatar           string    `json:"avatar"`
	ExpireTime       time.Time `json:"expireTime"`
	UserSize         int       `json:"userSize"`
	UserCount        int       `json:"userCount"`
	Status           int       `json:"status"`
	OrderID          int64     `json:"orderId"`
	SpuName          string    `json:"spuName"`
	PicUrl           string    `json:"picUrl"`
	Count            int       `json:"count"`
	CombinationPrice int       `json:"combinationPrice"`
}

// AppCombinationRecordDetailRespVO
type AppCombinationRecordDetailRespVO struct {
	HeadRecord    AppCombinationRecordRespVO   `json:"headRecord"`
	MemberRecords []AppCombinationRecordRespVO `json:"memberRecords"`
	OrderID       int64                        `json:"orderId"`
}

// AppCombinationRecordSummaryRespVO
type AppCombinationRecordSummaryRespVO struct {
	UserCount int64    `json:"userCount"`
	Avatars   []string `json:"avatars"`
}

// CombinationRecordPageItemRespVO 拼团记录 Admin 分页 VO
type CombinationRecordPageItemRespVO struct {
	ID               int64     `json:"id"`
	ActivityID       int64     `json:"activityId"`
	ActivityName     string    `json:"activityName"`
	SpuID            int64     `json:"spuId"`
	SpuName          string    `json:"spuName"`
	PicUrl           string    `json:"picUrl"`
	UserID           int64     `json:"userId"`
	Nickname         string    `json:"nickname"`
	Avatar           string    `json:"avatar"`
	UserCount        int       `json:"userCount"`
	UserSize         int       `json:"userSize"`
	Status           int       `json:"status"`
	CombinationPrice int       `json:"combinationPrice"`
	HeadID           int64     `json:"headId"`
	OrderID          int64     `json:"orderId"`
	VirtualGroup     bool      `json:"virtualGroup"`
	ExpireTime       time.Time `json:"expireTime"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	CreateTime       time.Time `json:"createTime"`
}

// CombinationRecordSummaryVO 拼团记录概要 VO (Admin)
type CombinationRecordSummaryVO struct {
	UserCount         int64 `json:"userCount"`         // 拼团用户参与数量
	SuccessCount      int64 `json:"successCount"`      // 成团数量
	VirtualGroupCount int64 `json:"virtualGroupCount"` // 虚拟成团数量
}
