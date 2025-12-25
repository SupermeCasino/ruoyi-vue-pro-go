package req

// AppSeckillActivityDetailReq App 端 - 秒杀活动详情请求
type AppSeckillActivityDetailReq struct {
	ID int64 `form:"id" binding:"required"` // 活动 ID
}

// AppSeckillActivityPageReq App 端 - 秒杀活动分页请求 (对齐 Java: AppSeckillActivityPageReqVO)
type AppSeckillActivityPageReq struct {
	PageNo   int    `form:"pageNo" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"`
	ConfigID *int64 `form:"configId"` // 秒杀时段编号
}
