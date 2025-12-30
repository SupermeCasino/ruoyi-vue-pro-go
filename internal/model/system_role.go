package model

// SystemRole 角色表
type SystemRole struct {
	ID               int64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name             string           `gorm:"column:name;not null" json:"name"`
	Code             string           `gorm:"column:code;not null" json:"code"`
	Sort             int32            `gorm:"column:sort" json:"sort"`
	DataScope        int32            `gorm:"column:data_scope;not null;default:1" json:"dataScope"`
	DataScopeDeptIds Int64ListFromCSV `gorm:"column:data_scope_dept_ids" json:"dataScopeDeptIds"` // 采用 CSV 适配器以兼容空值策略，避免产生 NULL
	Status           int32            `gorm:"column:status;not null" json:"status"`
	Type             int32            `gorm:"column:type;not null;default:1" json:"type"` // 角色类型(1:内置角色 2:自定义角色)
	Remark           string           `gorm:"column:remark" json:"remark"`
	TenantBaseDO
}

func (SystemRole) TableName() string {
	return "system_role"
}
