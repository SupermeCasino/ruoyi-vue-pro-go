package model

// SystemUserPost 用户和岗位关联表
type SystemUserPost struct {
	ID     int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID int64 `gorm:"column:user_id;not null" json:"userId"`
	PostID int64 `gorm:"column:post_id;not null" json:"postId"`
	TenantBaseDO
}

func (SystemUserPost) TableName() string {
	return "system_user_post"
}
