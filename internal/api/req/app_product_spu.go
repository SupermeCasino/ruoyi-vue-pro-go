package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// AppProductSpuPageReq 商品 SPU 分页 Request VO
type AppProductSpuPageReq struct {
	pagination.PageParam
	// 分类编号
	CategoryID *int64 `form:"categoryId"`
	// 关键字
	Keyword *string `form:"keyword"`
	// 排序字段
	SortField *string `form:"sortField"` // price, sales_count
	// 排序方式
	SortAsc *bool `form:"sortAsc"`
	// 推荐类型
	RecommendType *string `form:"recommendType"` // new, hot, best
}
