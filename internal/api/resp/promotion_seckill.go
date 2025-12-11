package resp

import "time"

// ========== Seckill Config ==========

// SeckillConfigResp 秒杀时段 Response
type SeckillConfigResp struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	StartTime     string    `json:"startTime"`
	EndTime       string    `json:"endTime"`
	SliderPicUrls []string  `json:"sliderPicUrls"`
	Status        int       `json:"status"`
	CreateTime    time.Time `json:"createTime"`
}

// SeckillConfigSimpleResp 秒杀时段精简 Response
type SeckillConfigSimpleResp struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// ========== Seckill Activity ==========

// SeckillProductResp 秒杀商品 Response
type SeckillProductResp struct {
	ID           int64 `json:"id"`
	ActivityID   int64 `json:"activityId"`
	SpuID        int64 `json:"spuId"`
	SkuID        int64 `json:"skuId"`
	SeckillPrice int   `json:"seckillPrice"` // 分
	Stock        int   `json:"stock"`
}

// SeckillActivityResp 秒杀活动 Response
type SeckillActivityResp struct {
	ID               int64     `json:"id"`
	SpuID            int64     `json:"spuId"`
	Name             string    `json:"name"`
	Status           int       `json:"status"`
	Remark           string    `json:"remark"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	Sort             int       `json:"sort"`
	ConfigIds        []int64   `json:"configIds"`
	TotalLimitCount  int       `json:"totalLimitCount"`
	SingleLimitCount int       `json:"singleLimitCount"`
	Stock            int       `json:"stock"`
	TotalStock       int       `json:"totalStock"`
	CreateTime       time.Time `json:"createTime"`
}

// SeckillActivityDetailResp 秒杀活动详情 Response
type SeckillActivityDetailResp struct {
	SeckillActivityResp
	Products []SeckillProductResp `json:"products"`
}
