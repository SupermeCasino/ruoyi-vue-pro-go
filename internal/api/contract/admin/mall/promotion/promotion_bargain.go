package promotion

import "time"

// BargainActivityCreateReq 创建砍价活动 Request
type BargainActivityCreateReq struct {
	SpuID             int64  `json:"spuId" binding:"required"`
	SkuID             int64  `json:"skuId" binding:"required"`
	Name              string `json:"name" binding:"required"`
	StartTime         string `json:"startTime" binding:"required"` // 格式: "2006-01-02 15:04:05"
	EndTime           string `json:"endTime" binding:"required"`   // 格式: "2006-01-02 15:04:05"
	BargainFirstPrice int    `json:"bargainFirstPrice" binding:"min=0"`
	BargainMinPrice   int    `json:"bargainMinPrice" binding:"min=0"`
	Stock             int    `json:"stock" binding:"min=0"`
	TotalStock        int    `json:"totalStock" binding:"min=0"`
	HelpMaxCount      int    `json:"helpMaxCount" binding:"min=1"`
	BargainCount      int    `json:"bargainCount" binding:"min=1"`
	TotalLimitCount   int    `json:"totalLimitCount" binding:"min=0"`
	RandomMinPrice    int    `json:"randomMinPrice" binding:"min=0"`
	RandomMaxPrice    int    `json:"randomMaxPrice" binding:"min=0"`
	Sort              int    `json:"sort"`
	Remark            string `json:"remark"`
}

// BargainActivityUpdateReq 更新砍价活动 Request
type BargainActivityUpdateReq struct {
	ID int64 `json:"id" binding:"required"`
	BargainActivityCreateReq
}

// BargainActivityPageReq 获得砍价活动分页 Request
type BargainActivityPageReq struct {
	PageNo   int    `json:"pageNo" form:"pageNo,default=1"`
	PageSize int    `json:"pageSize" form:"pageSize,default=10"`
	Name     string `json:"name" form:"name"`
	Status   *int   `json:"status" form:"status"`
}

// BargainRecordPageReq 获得砍价记录分页 Request
type BargainRecordPageReq struct {
	PageNo    int         `json:"pageNo" form:"pageNo,default=1"`
	PageSize  int         `json:"pageSize" form:"pageSize,default=10"`
	Status    *int        `json:"status" form:"status"`
	Name      string      `json:"name" form:"name"`                                             // User nickname? Or activity name? Usually filtered by activity or user.
	DateRange []time.Time `json:"dateRange" form:"dateRange" time_format:"2006-01-02 15:04:05"` // Creation Time Range
	// Actually in Admin, usually filter by Status, Period, User.
}

// BargainHelpPageReq 获得砍价助力分页 Request
type BargainHelpPageReq struct {
	PageNo     int    `json:"pageNo" form:"pageNo,default=1"`
	PageSize   int    `json:"pageSize" form:"pageSize,default=10"`
	Name       string `json:"name" form:"name"` // User nickname
	ActivityID int64  `json:"activityId" form:"activityId"`
	RecordID   int64  `json:"recordId" form:"recordId"`
}

// AppBargainRecordCreateReq App 创建砍价记录 Request
type AppBargainRecordCreateReq struct {
	ActivityID int64 `json:"activityId" binding:"required"`
}

// AppBargainHelpCreateReq App 创建砍价助力 Request
type AppBargainHelpCreateReq struct {
	RecordID int64 `json:"recordId" binding:"required"`
}

// BargainActivityResp 砍价活动 Response
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
	CreateTime        time.Time `json:"createTime"`
}

// BargainActivityPageItemResp 砍价活动分页项 Response
type BargainActivityPageItemResp struct {
	BargainActivityResp
	SpuName                string `json:"spuName"`
	PicUrl                 string `json:"picUrl"`
	MarketPrice            int    `json:"marketPrice"`
	RecordUserCount        int    `json:"recordUserCount"`
	RecordSuccessUserCount int    `json:"recordSuccessUserCount"`
	HelpUserCount          int    `json:"helpUserCount"`
}

// BargainHelpResp 砍价助力 Response
type BargainHelpResp struct {
	ID           int64     `json:"id"`
	ActivityID   int64     `json:"activityId"`
	RecordID     int64     `json:"recordId"`
	UserID       int64     `json:"userId"`
	UserNickname string    `json:"userNickname"`
	UserAvatar   string    `json:"userAvatar"`
	ReducePrice  int       `json:"reducePrice"`
	CreateTime   time.Time `json:"createTime"`
}

// BargainRecordResp 砍价记录 Response
type BargainRecordResp struct {
	ID                int64     `json:"id"`
	ActivityID        int64     `json:"activityId"`
	ActivityName      string    `json:"activityName"`
	SpuID             int64     `json:"spuId"`
	SpuName           string    `json:"spuName"`
	PicUrl            string    `json:"picUrl"`
	SkuID             int64     `json:"skuId"`
	UserID            int64     `json:"userId"`
	UserNickname      string    `json:"userNickname"`
	UserAvatar        string    `json:"userAvatar"`
	BargainFirstPrice int       `json:"bargainFirstPrice"`
	BargainCurPrice   int       `json:"bargainCurPrice"`
	BargainMinPrice   int       `json:"bargainMinPrice"`
	Status            int       `json:"status"`
	OrderID           *int64    `json:"orderId"`
	CreateTime        time.Time `json:"createTime"`
	EndTime           time.Time `json:"endTime"`
}
