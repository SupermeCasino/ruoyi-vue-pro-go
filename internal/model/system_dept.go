package model

// SystemDept 部门表
type SystemDept struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"column:name;not null;default:''" json:"name"`
	ParentID     int64  `gorm:"column:parent_id;not null;default:0" json:"parentId"`
	Sort         int32  `gorm:"column:sort;not null;default:0" json:"sort"`
	LeaderUserID int64  `gorm:"column:leader_user_id;default:0" json:"leaderUserId"`
	Phone        string `gorm:"column:phone;default:''" json:"phone"`
	Email        string `gorm:"column:email;default:''" json:"email"`
	Status       int32  `gorm:"column:status;not null;default:0" json:"status"` // 0:开启, 1:禁用
	TenantBaseDO
}

func (SystemDept) TableName() string {
	return "system_dept"
}

// SystemPost 岗位表
type SystemPost struct {
	ID     int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"column:name;not null;default:''" json:"name"`
	Code   string `gorm:"column:code;not null;default:''" json:"code"`
	Sort   int32  `gorm:"column:sort;not null;default:0" json:"sort"`
	Status int32  `gorm:"column:status;not null;default:0" json:"status"` // 0:开启, 1:禁用
	Remark string `gorm:"column:remark;default:''" json:"remark"`
	TenantBaseDO
}

func (SystemPost) TableName() string {
	return "system_post"
}
