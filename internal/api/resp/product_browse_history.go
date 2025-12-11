package resp

import "time"

// ProductBrowseHistoryResp (Admin)
type ProductBrowseHistoryResp struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userId"`
	SpuID      int64     `json:"spuId"`
	SpuName    string    `json:"spuName,omitempty"`
	PicURL     string    `json:"picUrl,omitempty"`
	Price      int64     `json:"price,omitempty"`
	SalesCount int       `json:"salesCount,omitempty"`
	Stock      int       `json:"stock,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

// AppProductBrowseHistoryResp (App)
type AppProductBrowseHistoryResp struct {
	ID         int64     `json:"id"`
	SpuID      int64     `json:"spuId"`
	SpuName    string    `json:"spuName,omitempty"`
	PicURL     string    `json:"picUrl,omitempty"`
	Price      int64     `json:"price,omitempty"`
	SalesCount int       `json:"salesCount,omitempty"`
	Stock      int       `json:"stock,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}
