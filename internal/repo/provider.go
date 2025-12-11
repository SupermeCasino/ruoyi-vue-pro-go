package repo

import (
	"backend-go/internal/repo/query"

	"gorm.io/gorm"
)

// NewQuery 适配 query.Use，屏蔽 opts 变长参数，方便 Wire 注入
func NewQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}
