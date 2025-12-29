package resp

import "time"

// TenantPackageResp 租户套餐响应 (对齐 Java: TenantPackageRespVO)
type TenantPackageResp struct {
	ID         int64     `json:"id"`         // 套餐编号
	Name       string    `json:"name"`       // 套餐名
	Status     int       `json:"status"`     // 状态
	Remark     string    `json:"remark"`     // 备注
	MenuIds    []int64   `json:"menuIds"`    // 关联菜单ID
	CreateTime time.Time `json:"createTime"` // 创建时间
}
