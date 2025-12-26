package model

// SystemRoleMenu 角色和菜单关联表
type SystemRoleMenu struct {
	ID     int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID int64 `gorm:"column:role_id;not null" json:"roleId"`
	MenuID int64 `gorm:"column:menu_id;not null" json:"menuId"`
	TenantBaseDO
}

func (SystemRoleMenu) TableName() string {
	return "system_role_menu"
}
