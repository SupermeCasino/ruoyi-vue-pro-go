package resp

import "time"

// PromotionBannerResp 响应
type PromotionBannerResp struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	PicURL    string    `json:"picUrl"`
	Url       string    `json:"url"`
	Status    int       `json:"status"`
	Sort      int       `json:"sort"`
	Position  int       `json:"position"`
	Memo      string    `json:"memo"`
	CreateTime time.Time `json:"createTime"`
}
