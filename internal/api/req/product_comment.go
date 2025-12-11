package req

// ProductCommentPageReq 商品评价分页 Request
type ProductCommentPageReq struct {
	PageNo       int      `form:"pageNo" binding:"required,min=1"`
	PageSize     int      `form:"pageSize" binding:"required,min=1,max=100"`
	UserNickname string   `form:"userNickname"`
	OrderID      int64    `form:"orderId"`
	SpuID        int64    `form:"spuId"`
	SpuName      string   `form:"spuName"`
	Scores       int      `form:"scores"`
	ReplyStatus  *bool    `form:"replyStatus"`
	CreateTime   []string `form:"createTime[]"` // time range
}

// ProductCommentUpdateVisibleReq 更新评论可见性 Request
type ProductCommentUpdateVisibleReq struct {
	ID      int64 `json:"id" binding:"required"`
	Visible bool  `json:"visible" binding:"required"`
}

// ProductCommentReplyReq 商家回复 Request
type ProductCommentReplyReq struct {
	ID      int64  `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// AppProductCommentCreateReq 添加自评 Request (App端)
type AppProductCommentCreateReq struct {
	OrderItemID       int64    `json:"orderItemId" binding:"required"`
	Anonymous         bool     `json:"anonymous"`
	Content           string   `json:"content" binding:"required"`
	PicURLs           []string `json:"picUrls"`
	Scores            int      `json:"scores" binding:"required,min=1,max=5"`
	DescriptionScores int      `json:"descriptionScores" binding:"required,min=1,max=5"`
	BenefitScores     int      `json:"benefitScores" binding:"required,min=1,max=5"`
}

// ProductCommentCreateReq 商品评价创建 Request (后台)
type ProductCommentCreateReq struct {
	UserID            int64    `json:"userId" binding:"required"`
	OrderItemID       int64    `json:"orderItemId" binding:"required"`
	UserNickname      string   `json:"userNickname" binding:"required"`
	UserAvatar        string   `json:"userAvatar" binding:"required"`
	SkuID             int64    `json:"skuId" binding:"required"`
	DescriptionScores int      `json:"descriptionScores" binding:"required,min=1,max=5"`
	BenefitScores     int      `json:"benefitScores" binding:"required,min=1,max=5"`
	Content           string   `json:"content" binding:"required"`
	PicURLs           []string `json:"picUrls"`
}

// AppProductCommentPageReq 商品评价分页 Request (App端)
type AppProductCommentPageReq struct {
	PageNo   int   `form:"pageNo" binding:"required,min=1"`
	PageSize int   `form:"pageSize" binding:"required,min=1,max=100"`
	SpuID    int64 `form:"spuId" binding:"required"`
	Type     int   `form:"type"` // 0: 全部, 1: 好评, 2: 中评, 3: 差评, 4: 有图
}
