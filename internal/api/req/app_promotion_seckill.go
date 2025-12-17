package req

// AppSeckillActivityPageReq App 端 - 秒杀活动分页请求 (对齐 Java: AppSeckillActivityPageReqVO)
type AppSeckillActivityPageReq struct {
	PageNo   int    `form:"pageNo" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"`
	ConfigID *int64 `form:"configId"` // 秒杀时段编号
}
