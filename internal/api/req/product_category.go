package req

// ProductCategoryCreateReq 创建商品分类请求
type ProductCategoryCreateReq struct {
	ParentID    int64  `json:"parentId" binding:"min=0"` // 允许为 0 (根节点)
	Name        string `json:"name" binding:"required,max=255"`
	PicURL      string `json:"picUrl"`
	Sort        int32  `json:"sort" binding:"min=0"`
	Status      int32  `json:"status" binding:"required,oneof=0 1"`
	Description string `json:"description"` // 分类描述
}

// ProductCategoryUpdateReq 更新商品分类请求
type ProductCategoryUpdateReq struct {
	ID          int64  `json:"id" binding:"required"`
	ParentID    int64  `json:"parentId" binding:"min=0"`
	Name        string `json:"name" binding:"required,max=255"`
	PicURL      string `json:"picUrl"`
	Sort        int32  `json:"sort" binding:"min=0"`
	Status      int32  `json:"status" binding:"required,oneof=0 1"`
	Description string `json:"description"` // 分类描述
}

// ProductCategoryListReq 列表查询请求
type ProductCategoryListReq struct {
	Name      string  `form:"name"`
	ParentID  *int64  `form:"parentId"`
	ParentIDs []int64 `form:"parentIds"` // 父分类编号数组
	Status    *int32  `form:"status"`
}
