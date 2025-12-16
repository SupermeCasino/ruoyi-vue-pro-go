package resp

// AppCategoryResp 用户 APP - 商品分类 Response VO
type AppCategoryResp struct {
	// 分类编号
	ID int64 `json:"id"`
	// 父分类编号
	ParentID int64 `json:"parentId"`
	// 分类名称
	Name string `json:"name"`
	// 分类图片
	PicURL string `json:"picUrl"`
}
