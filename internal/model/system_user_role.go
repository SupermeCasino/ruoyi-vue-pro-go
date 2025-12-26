package model

// SystemUserRole 用户和角色关联表
type SystemUserRole struct {
	ID     int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID int64 `gorm:"column:user_id;not null" json:"userId"`
	RoleID int64 `gorm:"column:role_id;not null" json:"roleId"`
	TenantBaseDO
}

func (SystemUserRole) TableName() string {
	return "system_user_role"
}
