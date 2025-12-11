package req

import "time"

// BargainActivityCreateReq 创建砍价活动 Request
type BargainActivityCreateReq struct {
	SpuID             int64     `json:"spuId" binding:"required"`
	SkuID             int64     `json:"skuId" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	StartTime         time.Time `json:"startTime" binding:"required"`
	EndTime           time.Time `json:"endTime" binding:"required"`
	BargainFirstPrice int       `json:"bargainFirstPrice" binding:"required,min=0"`
	BargainMinPrice   int       `json:"bargainMinPrice" binding:"required,min=0"`
	Stock             int       `json:"stock" binding:"required,min=0"`
	TotalStock        int       `json:"totalStock" binding:"required,min=0"`
	HelpMaxCount      int       `json:"helpMaxCount" binding:"required,min=1"`
	BargainCount      int       `json:"bargainCount" binding:"required,min=1"`
	TotalLimitCount   int       `json:"totalLimitCount" binding:"required,min=0"`
	RandomMinPrice    int       `json:"randomMinPrice" binding:"required,min=0"`
	RandomMaxPrice    int       `json:"randomMaxPrice" binding:"required,min=0"`
	Sort              int       `json:"sort"`
	Remark            string    `json:"remark"`
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
