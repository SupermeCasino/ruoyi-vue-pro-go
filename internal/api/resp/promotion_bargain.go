package resp

import "time"

// BargainActivityResp 砍价活动 Response (Admin)
type BargainActivityResp struct {
	ID                int64     `json:"id"`
	SpuID             int64     `json:"spuId"`
	SkuID             int64     `json:"skuId"`
	Name              string    `json:"name"`
	StartTime         time.Time `json:"startTime"`
	EndTime           time.Time `json:"endTime"`
	BargainFirstPrice int       `json:"bargainFirstPrice"`
	BargainMinPrice   int       `json:"bargainMinPrice"`
	Stock             int       `json:"stock"`
	TotalStock        int       `json:"totalStock"`
	HelpMaxCount      int       `json:"helpMaxCount"`
	BargainCount      int       `json:"bargainCount"`
	TotalLimitCount   int       `json:"totalLimitCount"`
	RandomMinPrice    int       `json:"randomMinPrice"`
	RandomMaxPrice    int       `json:"randomMaxPrice"`
	Status            int       `json:"status"`
	Sort              int       `json:"sort"`
	Remark            string    `json:"remark"`
	CreatedAt         time.Time `json:"createTime"`
}

// BargainActivityPageItemResp 砍价活动分页 Response (Admin)
type BargainActivityPageItemResp struct {
	BargainActivityResp
	SpuName                string `json:"spuName"`
	PicUrl                 string `json:"picUrl"`
	MarketPrice            int    `json:"marketPrice"`
	RecordUserCount        int    `json:"recordUserCount"`
	RecordSuccessUserCount int    `json:"recordSuccessUserCount"`
	HelpUserCount          int    `json:"helpUserCount"`
}

// BargainRecordResp 砍价记录 Response (Admin)
type BargainRecordResp struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"userId"`
	UserNickname      string    `json:"userNickname"`
	UserAvatar        string    `json:"userAvatar"`
	ActivityID        int64     `json:"activityId"`
	ActivityName      string    `json:"activityName"`
	SpuID             int64     `json:"spuId"`
	SkuID             int64     `json:"skuId"`
	BargainFirstPrice int       `json:"bargainFirstPrice"`
	BargainPrice      int       `json:"bargainPrice"`
	Status            int       `json:"status"`
	EndTime           time.Time `json:"endTime"`
	OrderID           int64     `json:"orderId"`
	CreatedAt         time.Time `json:"createTime"`
}

// BargainHelpResp 砍价助力 Response (Admin)
type BargainHelpResp struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	UserNickname string    `json:"userNickname"`
	UserAvatar   string    `json:"userAvatar"`
	ActivityID   int64     `json:"activityId"`
	RecordID     int64     `json:"recordId"`
	ReducePrice  int       `json:"reducePrice"`
	CreatedAt    time.Time `json:"createTime"`
}

// ========== App 端 Response ==========

// AppBargainActivityRespVO App 砍价活动 Response (匹配 Java AppBargainActivityRespVO)
type AppBargainActivityRespVO struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	SpuID           int64     `json:"spuId"`
	SkuID           int64     `json:"skuId"`
	Stock           int       `json:"stock"`
	PicUrl          string    `json:"picUrl"`
	MarketPrice     int       `json:"marketPrice"`
	BargainMinPrice int       `json:"bargainMinPrice"`
}

// AppBargainActivityDetailRespVO App 砍价活动详情 Response (匹配 Java AppBargainActivityDetailRespVO)
type AppBargainActivityDetailRespVO struct {
	AppBargainActivityRespVO
	BargainFirstPrice int    `json:"bargainFirstPrice"`
	HelpMaxCount      int    `json:"helpMaxCount"`
	BargainCount      int    `json:"bargainCount"`
	TotalLimitCount   int    `json:"totalLimitCount"`
	RandomMinPrice    int    `json:"randomMinPrice"`
	RandomMaxPrice    int    `json:"randomMaxPrice"`
	SuccessUserCount  int    `json:"successUserCount"`
	Remark            string `json:"remark"`
}

// AppBargainRecordRespVO App 砍价记录 Response (匹配 Java AppBargainRecordRespVO)
type AppBargainRecordRespVO struct {
	ID           int64     `json:"id"`
	SpuID        int64     `json:"spuId"`
	SkuID        int64     `json:"skuId"`
	ActivityID   int64     `json:"activityId"`
	Status       int       `json:"status"`
	BargainPrice int       `json:"bargainPrice"`
	ActivityName string    `json:"activityName"`
	EndTime      time.Time `json:"endTime"`
	PicUrl       string    `json:"picUrl"`
	OrderID      *int64    `json:"orderId,omitempty"`
	PayStatus    *bool     `json:"payStatus,omitempty"`
	PayOrderID   *int64    `json:"payOrderId,omitempty"`
}

// AppBargainRecordSummaryRespVO App 砍价记录概要 Response (匹配 Java AppBargainRecordSummaryRespVO)
type AppBargainRecordSummaryRespVO struct {
	SuccessUserCount int                               `json:"successUserCount"`
	SuccessList      []AppBargainRecordSummaryRecordVO `json:"successList"`
}

// AppBargainRecordSummaryRecordVO 概要中的记录 (匹配 Java AppBargainRecordSummaryRespVO.Record)
type AppBargainRecordSummaryRecordVO struct {
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	ActivityName string `json:"activityName"`
}

// AppBargainRecordDetailRespVO App 砍价记录详情 Response (匹配 Java AppBargainRecordDetailRespVO)
type AppBargainRecordDetailRespVO struct {
	ID                int64     `json:"id,omitempty"`
	ActivityID        int64     `json:"activityId"`
	UserID            int64     `json:"userId,omitempty"`
	SpuID             int64     `json:"spuId,omitempty"`
	SkuID             int64     `json:"skuId,omitempty"`
	BargainFirstPrice int       `json:"bargainFirstPrice,omitempty"`
	BargainPrice      int       `json:"bargainPrice,omitempty"`
	Status            int       `json:"status,omitempty"`
	OrderID           *int64    `json:"orderId,omitempty"`
	EndTime           time.Time `json:"endTime,omitempty"`
	ExpireTime        time.Time `json:"expireTime,omitempty"`
	HelpAction        *int      `json:"helpAction,omitempty"`
	PayStatus         *bool     `json:"payStatus,omitempty"`
	PayOrderID        *int64    `json:"payOrderId,omitempty"`
}

// AppBargainHelpRespVO App 砍价助力 Response (匹配 Java AppBargainHelpRespVO)
type AppBargainHelpRespVO struct {
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	ReducePrice int       `json:"reducePrice"`
	CreateTime  time.Time `json:"createTime"`
}

// 保留旧的别名以兼容
type AppBargainActivityResp = AppBargainActivityRespVO
type AppBargainActivityDetailResp = AppBargainActivityDetailRespVO
type AppBargainRecordResp = AppBargainRecordRespVO
type AppBargainRecordSummaryResp = AppBargainRecordSummaryRespVO
type AppBargainRecordDetailResp = AppBargainRecordDetailRespVO
type AppBargainHelpResp = AppBargainHelpRespVO
