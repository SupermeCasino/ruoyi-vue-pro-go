package req

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// ========== Seckill Config ==========

// SeckillConfigCreateReq 创建秒杀时段 Request
type SeckillConfigCreateReq struct {
	Name          string   `json:"name" binding:"required"`
	StartTime     string   `json:"startTime" binding:"required"` // HH:mm:ss
	EndTime       string   `json:"endTime" binding:"required"`
	SliderPicUrls []string `json:"sliderPicUrls"`
	Status        int      `json:"status"` // Can create with status? Java CreateReqVO extends BaseVO which might have Status? Or default.
	// Java BaseVO has name, startTime, endTime, sliderPicUrls, status.
}

// SeckillConfigUpdateReq 更新秒杀时段 Request
type SeckillConfigUpdateReq struct {
	ID            int64    `json:"id" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	StartTime     string   `json:"startTime" binding:"required"`
	EndTime       string   `json:"endTime" binding:"required"`
	SliderPicUrls []string `json:"sliderPicUrls"`
	Status        int      `json:"status"`
}

// SeckillConfigUpdateStatusReq 更新秒杀时段状态 Request
type SeckillConfigUpdateStatusReq struct {
	ID     int64 `json:"id" binding:"required"`
	Status *int  `json:"status" binding:"required"`
}

// SeckillConfigPageReq 分页 Request
type SeckillConfigPageReq struct {
	PageNo   int    `form:"pageNo"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
	Status   *int   `form:"status"`
}

// ========== Seckill Activity ==========

// SeckillProductBaseVO 秒杀商品 Base VO
type SeckillProductBaseVO struct {
	SkuID        int64 `json:"skuId" binding:"required"`
	SeckillPrice int   `json:"seckillPrice" binding:"required,min=0"` // 分
	Stock        int   `json:"stock" binding:"required,min=0"`
}

// SeckillActivityCreateReq 创建秒杀活动 Request
type SeckillActivityCreateReq struct {
	SpuID            int64                  `json:"spuId" binding:"required"`
	Name             string                 `json:"name" binding:"required"`
	Remark           string                 `json:"remark"`
	StartTime        time.Time              `json:"startTime" binding:"required"`
	EndTime          time.Time              `json:"endTime" binding:"required"`
	Sort             int                    `json:"sort" binding:"required"`
	ConfigIds        []int64                `json:"configIds" binding:"required"`
	TotalLimitCount  int                    `json:"totalLimitCount"`
	SingleLimitCount int                    `json:"singleLimitCount"`
	Products         []SeckillProductBaseVO `json:"products" binding:"required,dive"`
}

// SeckillActivityUpdateReq 更新秒杀活动 Request
type SeckillActivityUpdateReq struct {
	ID               int64                  `json:"id" binding:"required"`
	SpuID            int64                  `json:"spuId" binding:"required"`
	Name             string                 `json:"name" binding:"required"`
	Remark           string                 `json:"remark"`
	StartTime        time.Time              `json:"startTime" binding:"required"`
	EndTime          time.Time              `json:"endTime" binding:"required"`
	Sort             int                    `json:"sort" binding:"required"`
	ConfigIds        []int64                `json:"configIds" binding:"required"`
	TotalLimitCount  int                    `json:"totalLimitCount"`
	SingleLimitCount int                    `json:"singleLimitCount"`
	Products         []SeckillProductBaseVO `json:"products" binding:"required,dive"`
}

// SeckillActivityPageReq 分页 Request VO
type SeckillActivityPageReq struct {
	pagination.PageParam
	Name       string       `form:"name"`
	Status     *int         `form:"status"`
	ConfigID   *int64       `form:"configId"`
	CreateTime []*time.Time `form:"createTime"`
}
