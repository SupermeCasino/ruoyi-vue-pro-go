package model

// SystemTenantPackage 租户套餐表
type SystemTenantPackage struct {
	ID      int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name    string  `gorm:"column:name;not null" json:"name"`
	Status  int32   `gorm:"column:status;not null;default:0" json:"status"`
	MenuIDs []int64 `gorm:"column:menu_ids;serializer:json;type:text" json:"menuIds"` // JSON 数组存储
	Remark  string  `gorm:"column:remark;default:''" json:"remark"`
	BaseDO
}

func (SystemTenantPackage) TableName() string {
	return "system_tenant_package"
}
