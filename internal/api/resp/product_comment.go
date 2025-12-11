package resp

import "time"

// ProductCommentResp 商品评价 Repsonse
type ProductCommentResp struct {
	ID                int64                    `json:"id"`
	UserID            int64                    `json:"userId"`
	UserNickname      string                   `json:"userNickname"`
	UserAvatar        string                   `json:"userAvatar"`
	Anonymous         bool                     `json:"anonymous"`
	OrderID           int64                    `json:"orderId"`
	OrderItemID       int64                    `json:"orderItemId"`
	SpuID             int64                    `json:"spuId"`
	SpuName           string                   `json:"spuName"`
	SkuID             int64                    `json:"skuId"`
	SkuPicURL         string                   `json:"skuPicUrl"`
	SkuProperties     []ProductSkuPropertyResp `json:"skuProperties"`
	Visible           bool                     `json:"visible"`
	Scores            int                      `json:"scores"`
	DescriptionScores int                      `json:"descriptionScores"`
	BenefitScores     int                      `json:"benefitScores"`
	Content           string                   `json:"content"`
	PicURLs           []string                 `json:"picUrls"`
	ReplyStatus       bool                     `json:"replyStatus"`
	ReplyUserID       int64                    `json:"replyUserId"`
	ReplyContent      string                   `json:"replyContent"`
	ReplyTime         *time.Time               `json:"replyTime"`
	CreatedAt         time.Time                `json:"createTime"`
}

// AppProductCommentResp 商品评价 App Response
type AppProductCommentResp struct {
	ID            int64                    `json:"id"`
	UserNickname  string                   `json:"userNickname"`
	UserAvatar    string                   `json:"userAvatar"`
	Scores        int                      `json:"scores"`
	Content       string                   `json:"content"`
	PicURLs       []string                 `json:"picUrls"`
	ReplyContent  string                   `json:"replyContent"`
	SkuProperties []ProductSkuPropertyResp `json:"skuProperties"`
	CreatedAt     time.Time                `json:"createTime"`
}
