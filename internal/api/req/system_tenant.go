package req

// TenantPackagePageReq 租户套餐分页请求 (对齐 Java: TenantPackagePageReqVO)
type TenantPackagePageReq struct {
	PageNo     int      `form:"pageNo,default=1"`
	PageSize   int      `form:"pageSize,default=10"`
	Name       string   `form:"name"`       // 套餐名（模糊查询）
	Status     *int     `form:"status"`     // 状态
	Remark     string   `form:"remark"`     // 备注（模糊查询）
	CreateTime []string `form:"createTime"` // 创建时间范围
}

// TenantPackageSaveReq 租户套餐保存请求 (对齐 Java: TenantPackageSaveReqVO)
type TenantPackageSaveReq struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name" binding:"required"`
	Status  *int    `json:"status" binding:"required"`
	Remark  string  `json:"remark"`
	MenuIDs []int64 `json:"menuIds" binding:"required"`
}
